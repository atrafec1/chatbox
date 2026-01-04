package main

import (
	"chatbox/database"
	"chatbox/server"
	"log"
)

func main() {
	log.Println("Starting TCP server on port 9090")
	db, err := database.InitDB()

	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to make db instance: %v", err)
	}
	defer sqlDB.Close()
	server.StartServer("9090")

	log.Println("tcp server and database successfully initialized")

}
