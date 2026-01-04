package main

import (
	"chatbox/server"
	"log"
)

func main() {
	port := "9090"
	if err := server.StartServer(port); err != nil {
		log.Fatal(err)
	}
}
