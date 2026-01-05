package server

import (
	"fmt"
	"io"
	"net"
	"sync"

	"chatbox/database"

	"gorm.io/gorm"
)

// ChatServer represents the main server
type ChatServer struct {
	Sessions map[string]*Session
	Users    map[string]*User
	Groups   map[string]*Group

	DB       *gorm.DB
	Address  string
	Listener net.Listener
	mu       sync.Mutex
}

// StartServer starts listening on the given port with the database attached
func StartServer(port string, db *gorm.DB) error {
	server := &ChatServer{
		Sessions: make(map[string]*Session),
		Users:    make(map[string]*User),
		Groups:   make(map[string]*Group),
		DB:       db,
		Address:  ":" + port,
	}

	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		return fmt.Errorf("error creating listener for port %s: %w", port, err)
	}
	server.Listener = listener

	fmt.Println("Server listening on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("could not accept connection:", err)
			continue
		}

		go server.handleConnection(conn)
	}
}

// handleConnection manages a single client connection
func (server *ChatServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	client := NewClient(conn)

	// Prompt username
	session, err := server.onboardUser(client)
	if err != nil {
		fmt.Println("could not onboard user:", err)
		return
	}
	for {
		msg, err := session.ReadMsg()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("connection ended for: %v\n", conn)
			} else {
				fmt.Println("could not read message: ", err)
			}
			break
		}

		fmt.Println("message:", msg)
		session.SendMsg("wow, Nice!")
	}
}

func (server *ChatServer) authenticateUser(c *Client) (*User, error) {
	username, err := server.promptUsername(c)
	if err != nil {
		return nil, err
	}

	userExists, err := database.UsernameExists(server.DB, username)
	if err != nil {
		return nil, err
	}
	password, err := server.promptPassword(c)
	if err != nil {
		return nil, err
	}

	if !userExists {
		return server.registerUser(c, username, password)
	} else {
		return server.loginUser(c, username, password)
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
	session := NewSession(user, client)
	session.SendMsg("Now logged in as: " + user.Name)
	return session, nil
}

func (s *ChatServer) promptUsername(c *Client) (string, error) {

	if err := c.SendMessage("Enter username: "); err != nil {
		return "", err
	}
	username, err := c.ReadMessage()
	if err != nil {
		return "", err
	}
	return username, nil
}

func (s *ChatServer) promptPassword(c *Client) (string, error) {
	if err := c.SendMessage("Password: "); err != nil {
		return "", err
	}
	password, err := c.ReadMessage()
	if err != nil {
		return "", err
	}
	return password, nil
}

func (s *ChatServer) registerUser(c *Client, username, password string) (*User, error) {
	user, err := database.RegisterUser(s.DB, username, password)
	if err != nil {
		return nil, err
	}
	return &User{
		id:   user.ID,
		Name: user.Username,
	}, nil
}

func (s *ChatServer) loginUser(c *Client, username, password string) (*User, error) {
	user, err := database.Login(s.DB, username, password)
	if err != nil {
		return nil, err
	}
	return &User{
		id:   user.ID,
		Name: user.Username,
	}, nil
}
