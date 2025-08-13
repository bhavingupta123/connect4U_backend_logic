package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn         *websocket.Conn
	ID           int
	Disconnected bool
	InRematch    bool
	Mutex        sync.Mutex
	Name         string
}
