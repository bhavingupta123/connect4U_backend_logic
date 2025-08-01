package server

import (
	"log"
	"ludo_backend_refactored/internal/config"
	"ludo_backend_refactored/internal/game"
	"ludo_backend_refactored/internal/player"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	matchQueue = make(chan *player.Player, 100)
	queueLock  sync.Mutex
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	p := &player.Player{Conn: conn}
	go InitPlayer(p)
}

func InitPlayer(p *player.Player) {
	var initMsg map[string]interface{}
	if err := p.Conn.ReadJSON(&initMsg); err != nil {
		log.Println("Failed to read init message:", err)
		p.Conn.Close()
		return
	}
	mode := initMsg["mode"].(string)
	if mode == "bot" {
		g := &game.Game{Players: [2]*player.Player{p, nil}, Board: game.NewBoard()}
		p.ID = 0
		startGame(g)
	} else {
		matchPlayer(p)
	}
}

func matchPlayer(p *player.Player) {
	select {
	case opp := <-matchQueue:
		if opp != nil {
			startGame(&game.Game{
				Players: [2]*player.Player{opp, p},
				Board:   game.NewBoard(),
			})
			return
		}
	default:
		queueLock.Lock()
		matchQueue <- p
		queueLock.Unlock()

		time.AfterFunc(config.WaitSec*time.Second, func() {
			queueLock.Lock()
			defer queueLock.Unlock()
			select {
			case opp := <-matchQueue:
				if opp == p {
					startGame(&game.Game{
						Players: [2]*player.Player{p, nil},
						Board:   game.NewBoard(),
					})
				} else {
					startGame(&game.Game{
						Players: [2]*player.Player{opp, p},
						Board:   game.NewBoard(),
					})
				}
			default:
			}
		})
	}
}

func startGame(g *game.Game) {
	p1 := g.Players[0]
	p1.ID = 0
	p1.Conn.WriteJSON(map[string]interface{}{"type": "match_found", "player": 0})

	if g.Players[1] != nil {
		p2 := g.Players[1]
		p2.ID = 1
		p2.Conn.WriteJSON(map[string]interface{}{"type": "match_found", "player": 1})
		go listen(g, p1)
		go listen(g, p2)
	} else {
		p1.Conn.WriteJSON(map[string]interface{}{"type": "bot_mode", "player": 0})
		go listen(g, p1)
		go botLoop(g)
	}
}

func listen(g *game.Game, p *player.Player) {
	for {
		var msg map[string]interface{}
		if err := p.Conn.ReadJSON(&msg); err != nil {
			p.Mutex.Lock()
			p.Disconnected = true
			p.Mutex.Unlock()
			return
		}
		g.Mutex.Lock()
		switch msg["type"] {
		case "move":
			col := int(msg["column"].(float64))
			if g.Turn != p.ID || g.GameOver || !g.Board.IsValidMove(col) {
				break
			}
			row := g.Board.ApplyMove(col, p.ID+1)
			sendToAll(g, map[string]interface{}{
				"type":   "move",
				"column": col,
				"player": p.ID,
			})

			if g.Board.CheckWin(col, row, p.ID+1) {
				g.GameOver = true
				sendToAll(g, map[string]interface{}{
					"type":   "game_over",
					"winner": p.ID,
				})
				break
			}
			g.Turn = 1 - p.ID
		case "rematch_request":
			p.InRematch = true
			g.HandleRematchRequest(p)
		case "rematch_cancelled":
			p.InRematch = false
			g.HandleRematchCancel(p)
		}
		g.Mutex.Unlock()
	}
}

func sendToAll(g *game.Game, msg map[string]interface{}) {
	for _, p := range g.Players {
		if p != nil && !p.Disconnected {
			p.Conn.WriteJSON(msg)
		}
	}
}

func botLoop(g *game.Game) {
	human := g.Players[0]

	for {
		time.Sleep(700 * time.Millisecond)

		g.Mutex.Lock()

		if g.GameOver || human.Disconnected {
			g.Mutex.Unlock()
			return
		}

		if g.Turn == 1 {
			// ✅ Check if HUMAN already won (human ID = 1)
			if g.Board.HasAnyWin(1) {
				g.GameOver = true
				human.Conn.WriteJSON(map[string]interface{}{
					"type":   "game_over",
					"winner": 0,
				})
				g.Mutex.Unlock()
				return
			}

			// ✅ Get bot move (bot internal ID = 2)
			col := g.BotBestMoveMiniMax()
			if !g.Board.IsValidMove(col) {
				for c := 0; c < game.Cols; c++ {
					if g.Board.IsValidMove(c) {
						col = c
						break
					}
				}
			}

			// ✅ Apply bot move with ID = 2
			row := g.Board.ApplyMove(col, 2)

			// ✅ Notify client using external ID = 1 (bot)
			human.Conn.WriteJSON(map[string]interface{}{
				"type":   "move",
				"column": col,
				"player": 1,
			})

			// ✅ Check win for bot (ID = 2)
			if g.Board.CheckWin(col, row, 2) {
				g.GameOver = true
				human.Conn.WriteJSON(map[string]interface{}{
					"type":   "game_over",
					"winner": 1,
				})
				g.Mutex.Unlock()
				return
			}

			g.Turn = 0
		}

		g.Mutex.Unlock()
	}
}
