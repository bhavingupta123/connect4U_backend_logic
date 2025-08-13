package game

import (
	player "ludo_backend_refactored/internal/model/player"
	"time"
)

func (g *Game) HandleRematchRequest(p *player.Player) {
	g.RematchLock.Lock()
	defer g.RematchLock.Unlock()

	g.RematchVotes[p.ID] = true
	other := g.Players[1-p.ID]

	if other != nil && g.RematchVotes[1-p.ID] {
		g.Players[0].InRematch = false
		g.Players[1].InRematch = false
		g.ResetGame()
		return
	}

	if other != nil && !other.Disconnected {
		other.Conn.WriteJSON(map[string]interface{}{
			"type":    "rematch_offer",
			"message": "Opponent wants rematch?",
		})

		go func(requestingPlayer *player.Player) {
			time.Sleep(5 * time.Second)

			g.RematchLock.Lock()
			defer g.RematchLock.Unlock()

			if !g.RematchVotes[1-requestingPlayer.ID] {
				requestingPlayer.Conn.WriteJSON(map[string]interface{}{
					"type":    "rematch_declined",
					"message": "Opponent did not respond in time.",
				})
				requestingPlayer.InRematch = false
			}
		}(p)
	}
}

func (g *Game) HandleRematchCancel(p *player.Player) {
	p.InRematch = false
	g.RematchVotes[p.ID] = false
	other := g.Players[1-p.ID]
	if other != nil && !other.Disconnected {
		other.Conn.WriteJSON(map[string]interface{}{
			"type":    "rematch_declined",
			"message": "Rematch was not accepted.",
		})
	}
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
