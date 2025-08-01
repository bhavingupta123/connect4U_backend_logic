package main

import (
	"fmt"
	"log"
	"ludo_backend_refactored/internal/server"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", server.HandleWebSocket)
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
