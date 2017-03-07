package api

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

type WSPool struct {
	clients   map[*websocket.Conn]struct{}
	register  chan *websocket.Conn
	broadcast chan interface{}
}

func NewWSPool() *WSPool {
	return &WSPool{
		clients:   make(map[*websocket.Conn]struct{}),
		register:  make(chan *websocket.Conn),
		broadcast: make(chan interface{}),
	}
}

func (ws *WSPool) Run() {
	for {
		select {
		case c := <-ws.register:
			ws.clients[c] = struct{}{}

		case i := <-ws.broadcast:
			msg, err := json.Marshal(i)
			if err != nil {
				continue
			}

			jsonString := string(msg)

			for c := range ws.clients {
				err := websocket.Message.Send(c, jsonString)
				if err != nil {
					delete(ws.clients, c)
				}
			}

		}
	}
}
