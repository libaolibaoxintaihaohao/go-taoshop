package main

import (
	"log"

	"taoshop/internal/app"
)

func main() {
	server, err := app.NewServer()
	if err != nil {
		log.Fatalf("bootstrap server: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
