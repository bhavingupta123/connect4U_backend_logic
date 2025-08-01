package game

import (
	"ludo_backend_refactored/internal/player"
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
}

func (g *Game) ResetGame() {
	g.Turn = 0
	g.GameOver = false
	g.RematchVotes = [2]bool{false, false}
	g.Board.Reset()
	for idx, p := range g.Players {
		if p != nil && !p.Disconnected {
			p.Conn.WriteJSON(map[string]interface{}{
				"type":   "rematch_start",
				"player": idx,
			})
		}
	}
}
