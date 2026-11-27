package common

import "github.com/mbarolo/webp_backend/models/ws"

type Result struct {
	State    any
	Finished bool
	Winner   int
}

type Game interface {
	MaxPlayers() int
	State() any
	HandleMessage(seat int, msg ws.Envelope) (Result, error)
	OnAbandon(connectedSeats []int) Result
}
