package main

import (
	"log"
	"module/module"
	"module/module/data"
)

func main() {
	err := data.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db := data.DB
	server := module.NewAPIServer(":8080", db)
	server.Run()
}
