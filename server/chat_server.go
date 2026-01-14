package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"chatbox/database"
	"chatbox/domain"

	"gorm.io/gorm"
)

// ChatServer represents the main server
type ChatServer struct {
	Sessions map[uint]*Session
	Users    map[uint]*User
	Groups   map[uint]*Group
	inbox    chan *Message
	DB       *gorm.DB
	Address  string
	Listener net.Listener
	mu       sync.Mutex
}

// StartServer starts listening on the given port with the database attached
func StartServer(port string, db *gorm.DB) error {
	server := &ChatServer{
		Sessions: make(map[uint]*Session),
		Users:    make(map[uint]*User),
		Groups:   make(map[uint]*Group),
		inbox:    make(chan *Message),
		DB:       db,
		Address:  ":" + port,
	}

	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		return fmt.Errorf("error creating listener for port %s: %w", port, err)
	}
	server.Listener = listener

	fmt.Println("Server listening on", port)
	go server.run()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("could not accept connection:", err)
			continue
		}
		go server.handleConnection(conn)
	}
}

func (s *ChatServer) run() {
	for msg := range s.inbox {
		if err := s.routeMessage(msg); err != nil {
			fmt.Println("route error:", err)
		}
	}
}

func (server *ChatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	client := NewClient(conn)

	session, err := server.onboardUser(client)
	if err != nil {
		fmt.Println("could not onboard user:", err)
		return
	}
	group, err := database.GetGroupByID(server.DB, session.User.GroupID)
	if err != nil {
		fmt.Println("could not get user group:", err)
		return
	}

	// Get existing group or create once
	memoryGroup := server.getOrCreateGroup(group.ID, group.Name)
	memoryGroup.Add(session)
	fmt.Printf("User %v connected and added to group %v\n", session.User.Name, group.Name)

	if err := server.IOLoop(session); err != nil {
		fmt.Println("error in IO loop:", err)
	}

	// Cleanup on disconnect
	memoryGroup.Remove(session)
}

func (server *ChatServer) getOrCreateGroup(id uint, name string) *Group {
	server.mu.Lock()
	defer server.mu.Unlock()
	if group, exists := server.Groups[id]; exists {
		return group
	}
	group := NewGroup(id, name) // calls go g.distributeMessages()
	server.Groups[id] = group
	return group
}

func (server *ChatServer) IOLoop(s *Session) error {
	for {
		msg, err := s.ReadMsg()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("connection ended for: %v\n", s.User.Name)
				return nil
			}
			return fmt.Errorf("could not read message: %w", err)
		}

		if IsCommand(msg) {
			server.handleCommand(msg, s)
		}

		m := &Message{
			GroupID:  s.User.GroupID,
			UserID:   s.User.id,
			Username: s.User.Name,
			Content:  msg, // raw content, not formatted yet
		}
		server.saveMessage(m)

		server.inbox <- m
	}
}

func (s *ChatServer) onboardUser(client *Client) (*Session, error) {
	if err := client.SendMessage("Welcome to chatbox!"); err != nil {
		return nil, err
	}
	user, err := s.authenticateUser(client)
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate user: %w", err)
	}
	session := NewSession(s, user, client)
	session.SendMsg("Now logged in as: " + user.Name)
	return session, nil
}

// Messaging
func (s *ChatServer) handleCommand(raw_command string, session *Session) {
	cmd, err := ParseCommand(raw_command)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCommandArgs):
			session.SendMsg("invalid command arguments")
		case errors.Is(err, domain.ErrNotEnoughArguments):
			session.SendMsg("not enough arguments provided")
		case errors.Is(err, domain.ErrTooManyArguments):
			session.SendMsg("too many arguments provided")
		default:
			session.SendMsg(fmt.Sprintf("error: %v", err))
		}
		return
	}
	if err := session.ExecuteCommand([]string{}, cmd); err != nil {
		session.SendMsg(fmt.Sprintf("command error: %v", err))
	}
}

func (s *ChatServer) saveMessage(m *Message) {
	database.SaveMessage(s.DB, m.Content, m.UserID, m.GroupID)
}

func (s *ChatServer) getGroup(groupID uint) (*Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	group, exists := s.Groups[groupID]
	if !exists {
		return nil, fmt.Errorf("failed to get users group")
	}
	return group, nil
}

func (s *ChatServer) routeMessage(m *Message) error {
	group, err := s.getGroup(m.GroupID)
	if err != nil {
		return err
	}
	for _, id := range group.members {
		fmt.Printf("User: %v", id.User.Name)
	}
	fmt.Printf("Group members: %v", group.members)
	group.Enqueue(m)
	return nil
}
