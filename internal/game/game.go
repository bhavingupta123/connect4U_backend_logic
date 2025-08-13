package game

import (
	player "ludo_backend_refactored/internal/model/player"
	"sync"
)

type Game struct {
	Players      [2]*player.Player
	Board        *Board
	Turn         int
	GameOver     bool
	Mutex        sync.Mutex
	RematchVotes [2]bool
	RematchLock  sync.Mutex
	Stats        Service
}
