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

type TCPServer struct {
	address  string
	listener net.Listener
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
	var client Client

	client.Conn = conn
	client.Reader = bufio.NewReader(conn)
	client.Writer = bufio.NewWriter(conn)

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
		fmt.Println(msg)
	}
}

func promptUsername() {

}
