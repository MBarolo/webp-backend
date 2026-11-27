package sessions

import "github.com/gorilla/websocket"

type Player struct {
	Conn *websocket.Conn
}
