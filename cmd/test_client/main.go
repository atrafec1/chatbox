package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("connected to server")

	// TCP reader/writer
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)

	// stdin reader
	stdin := bufio.NewReader(os.Stdin)

	// ---- RECEIVE LOOP (server → client)
	go func() {
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("disconnected from server")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	// ---- SEND LOOP (client → server)
	for {
		input, err := stdin.ReadString('\n')
		if err != nil {
			return
		}

		_, err = serverWriter.WriteString(input)
		if err != nil {
			return
		}
		serverWriter.Flush()
	}
}
