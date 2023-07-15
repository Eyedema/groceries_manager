package main

import (
	"log"
	"net/http"
	"ubaldo/api_server/internal/server"
)

func main() {
	s := server.NewServer()

	log.Println("Server is starting...")
	err := http.ListenAndServe(":8080", s)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
