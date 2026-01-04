package main

import (
	"chatbox/database"
	"log"
)

func main() {
	db, err := database.InitDB()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()

	log.Println("Database initialized successfully:", db)
}
