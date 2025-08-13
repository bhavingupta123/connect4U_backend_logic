package app

import (
	"fmt"
	consts "ludo_backend_refactored/internal/config"
	"ludo_backend_refactored/internal/game"
	"ludo_backend_refactored/internal/repo"
	"ludo_backend_refactored/internal/server"
	"net/http"
)

func Start() error {
	repository, err := repo.NewMongoRepository(consts.MongoURI)
	if err != nil {
		return fmt.Errorf("MongoDB connection failed: %w", err)
	}

	gameService := game.NewService(repository)
	server.SetStatsService(gameService)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", server.NewWebSocketHandler(gameService))
	mux.HandleFunc("/stats", server.StatsHandler(gameService))

	fmt.Println("Server started on :8080")
	return http.ListenAndServe(":8080", mux)
}
