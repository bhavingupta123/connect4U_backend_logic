package server

import (
	"encoding/json"
	"ludo_backend_refactored/internal/game"
	"net/http"
)

func StatsHandler(service game.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerName := r.URL.Query().Get("player")
		if playerName == "" {
			http.Error(w, "Missing player name", http.StatusBadRequest)
			return
		}

		stats, err := service.GetStats(playerName) // call the stats service
		if err != nil {
			http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"results": stats}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
