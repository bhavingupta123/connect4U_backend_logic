package game

import (
	model "ludo_backend_refactored/internal/model/stat"
	"ludo_backend_refactored/internal/repo"
)

type Service interface {
	RecordMatch(player, opponent, result string) error
	GetStats(player string) ([]model.MatchResult, error)
}

type service struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) RecordMatch(player, opponent, result string) error {
	return s.repo.SaveResult(model.MatchResult{
		Player:   player,
		Opponent: opponent,
		Result:   result,
	})
}

func (s *service) GetStats(player string) ([]model.MatchResult, error) {
	return s.repo.GetStatsForPlayer(player)
}
