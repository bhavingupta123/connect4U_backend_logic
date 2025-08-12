package main

import (
	"log"
	application "ludo_backend_refactored/internal/app"
)

func main() {
	if err := application.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
