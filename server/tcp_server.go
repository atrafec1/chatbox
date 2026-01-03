package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

//Start listener
//Start accept loop
//accept connections
//start handler for cennections
// make groups
//assign connections to groups

func StartServer(port string) error {
	var group Group
	group.members = make(map[net.Conn]*User)
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

		go handleConnection(conn, &group)
	}
}

func handleConnection(conn net.Conn, group *Group) {
	defer conn.Close()
	newUser := NewUser("Adam")
	group.AddMember(conn, newUser)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer writer.Flush()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("connection ended for: %v\n", conn)
				user, ok := group.Members[conn]
				if !ok {
					fmt.Println("could not find user for connection")
					break
				}
				group.RemoveMember(conn)
				fmt.Printf("removed member: %v\n", user)
				break
			}
			fmt.Println("could not read message: ", err)
		}
		group.relayMessage(group.Members[conn].Name + ": " + line)
	}
	group.RemoveMember(conn)
}
func (group *Group) relayMessage(message string) {
	group.Mu.Lock()
	defer group.Mu.Unlock()
	for conn := range group.Members {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message")
		}
	}

}
