package repo

import "ludo_backend_refactored/internal/model"

type Repository interface {
	SaveResult(result model.MatchResult) error
	GetStatsForPlayer(playerName string) ([]model.MatchResult, error)
}
