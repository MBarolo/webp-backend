package matches

import "github.com/gorilla/websocket"

type Player struct {
	UserId int64
	Conn   *websocket.Conn
}
