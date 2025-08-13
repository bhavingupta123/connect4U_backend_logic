package repo

import model "ludo_backend_refactored/internal/model/stat"

type Repository interface {
	SaveResult(result model.MatchResult) error
	GetStatsForPlayer(playerName string) ([]model.MatchResult, error)
}
