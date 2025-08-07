package main

import (
	"encoding/json"
	"fmt"
	"log"
	"ludo_backend_refactored/internal/game"
	"ludo_backend_refactored/internal/repo"
	"ludo_backend_refactored/internal/server"
	"net/http"
)

func main() {
	uri := "mongodb+srv://admin:admin@cluster0.5tfshir.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	repo, err := repo.NewMongoRepository(uri)
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}
	service := game.NewService(repo)
	server.SetStatsService(service)

	http.HandleFunc("/ws", server.NewWebSocketHandler(service))

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		playerName := r.URL.Query().Get("player")
		if playerName == "" {
			http.Error(w, "Missing player name", http.StatusBadRequest)
			return
		}

		stats, err := service.GetStats(playerName)
		if err != nil {
			http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"results": stats,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
