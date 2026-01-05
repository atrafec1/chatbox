package server

import (
	"fmt"
	"io"
	"net"
	"sync"

	"chatbox/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatServer represents the main server
type ChatServer struct {
	Sessions map[string]*Session
	Users    map[string]*database.User
	Groups   map[string]*database.Group

	DB       *gorm.DB
	Address  string
	Listener net.Listener
	mu       sync.Mutex
}

// StartServer starts listening on the given port with the database attached
func StartServer(port string, db *gorm.DB) error {
	server := &ChatServer{
		Sessions: make(map[string]*Session),
		Users:    make(map[string]*database.User),
		Groups:   make(map[string]*database.Group),
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
	username := server.promptUsername(client)
	user := server.newUser(username)
	session := NewSession(user, client)

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

func (s *ChatServer) newUser(username string) *User {
	id := uuid.NewString()
	return &User{
		id:   id,
		Name: username,
	}
}

// UserExists checks if a user exists in memory or DB
func (server *ChatServer) UserExists(username string) bool {
	server.mu.Lock()
	defer server.mu.Unlock()
	_, exists := server.Users[username]
	return exists
}

// AddUser adds a user to the server memory
func (server *ChatServer) AddUser(user *database.User) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.Users[user.Username] = user
}

//  ******  User Logic ********

func (server *ChatServer) loginUser(c *Client) (*User, error) {
	for {
		if err := c.SendMessage("Username: "); err != nil {
			return nil, err
		}

		username, err := c.ReadMessage()
		if err != nil {
			return nil, err
		}

		user := server.newUser(username)
		return user, nil
	}
}

func (s *ChatServer) registerUser(username string) (*database.User, error) {

	user, err := database.CreateUser(s.DB, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ChatServer) onboardUser(client *Client) (*Session, error) {
	if err := client.SendMessage("Welcome to chatbox!"); err != nil {
		return nil, err
	}

}
