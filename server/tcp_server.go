package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

//Start listener
//Start accept loop
//accept connections
//start handler for cennections
// make groups
//assign connections to groups

type TCPServer struct {
	Sessions map[string]*Session
	Users    map[string]*User
	Groups   map[string]*Group

	address  string
	listener net.Listener
	mu       sync.Mutex
}

func StartServer(port string) error {
	address := ":" + port

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("error creating listener for port %s: %w", port, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("could not accept connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	client := NewClient(conn)

	for {
		msg, err := client.ReadMessage()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("connection ended for: %v\n", conn)
			} else {
				fmt.Println("could not read message: ", err)
			}
			break
		}
		fmt.Println("message:", msg)
		client.SendMessage("thank you for your message! - server")
	}
}

func (server *TCPServer) UserExists(username string) bool {
	server.mu.Lock()
	defer server.mu.Unlock()
	_, exists := server.Users[username]
	return exists
}

func (server *TCPServer) AddUser(user *User) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.Users[user.Name] = user
}

func (server *TCPServer) NewUser(name string) *User {
	return &User{
		Name: name,
		id:   "randomIDForNow",
		mu:   sync.RWMutex{},
	}
}

func (server *TCPServer) promptUsername(client *Client) {
	client.SendMessage("Enter your username: ")
	if username, err := client.ReadMessage(); err == nil {
		if server.UserExists(username) {
			client.SendMessage("Username already exists. Try again. \n")
			server.promptUsername(client)
		}

	}
}
