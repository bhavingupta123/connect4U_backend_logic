package game

import "ludo_backend_refactored/internal/player"

func (g *Game) HandleRematchRequest(p *player.Player) {
	g.RematchLock.Lock()
	defer g.RematchLock.Unlock()

	g.RematchVotes[p.ID] = true
	other := g.Players[1-p.ID]
	if other != nil && g.RematchVotes[1-p.ID] {
		g.Players[0].InRematch = false
		g.Players[1].InRematch = false
		g.ResetGame()
	} else if other != nil && !other.Disconnected {
		other.Conn.WriteJSON(map[string]interface{}{
			"type":    "rematch_offer",
			"message": "Opponent wants rematch?",
		})
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
