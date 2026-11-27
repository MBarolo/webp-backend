package sessions

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/looplab/fsm"
	"github.com/mbarolo/webp_backend/models/common"
	"github.com/mbarolo/webp_backend/models/ws"
)

const TIMEOUT_PERIOD = 1 * time.Minute
const WRITE_TIMEOUT = 5 * time.Second
const PONG_WAIT = 15 * time.Second
const PING_PERIOD = 10 * time.Second

type Session struct {
	mu         sync.Mutex
	Code       string
	Game       common.Game
	Players    []*Player
	State      *fsm.FSM
	graceTimer *time.Timer
}

func NewSession(code string, game common.Game) *Session {
	s := &Session{Code: code, Game: game}
	s.State = GetStates(s)

	return s
}

func (s *Session) seatOf(p *Player) int {
	for seat, other := range s.Players {
		if other == p {
			return seat
		}
	}
	return -1
}

func (s *Session) freeSeat() int {
	for seat, p := range s.Players {
		if p.Conn == nil {
			return seat
		}
	}
	return -1
}

func (s *Session) connectedSeats() []int {
	var seats []int
	for seat, player := range s.Players {
		if player.Conn != nil {
			seats = append(seats, seat)
		}
	}

	return seats
}

func (s *Session) CanJoin() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.State.Current() {
	case WAITING:
		if len(s.Players) < s.Game.MaxPlayers() {
			return nil
		}
	case PAUSED:
		return nil
	}
	return ErrSessionFull
}

func (s *Session) Join(conn *websocket.Conn) (*Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.State.Current() {
	case WAITING:
		if len(s.Players) >= s.Game.MaxPlayers() {
			return nil, ErrSessionFull
		}
		p := &Player{Conn: conn}
		s.Players = append(s.Players, p)
		s.send(p, gin.H{"type": "joined", "code": s.Code, "seat": len(s.Players) - 1})

		if len(s.Players) == s.Game.MaxPlayers() {
			s.State.Event(context.Background(), START)
			s.broadcast(gin.H{"type": "game_started", "state": s.Game.State()})
		}
		return p, nil

	case PAUSED:
		seat := s.freeSeat()
		p := s.Players[seat]
		p.Conn = conn

		if s.freeSeat() == -1 {
			s.State.Event(context.Background(), RECONNECT)
			s.broadcast(gin.H{"type": "player_reconnected", "seat": seat})
		}
		s.send(p, gin.H{"type": "state", "seat": seat, "data": s.Game.State()})
		return p, nil
	}

	return nil, ErrSessionFull
}

func (s *Session) Ping(p *Player, conn *websocket.Conn) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p.Conn != conn {
		return errors.New("connection replaced")
	}
	conn.SetWriteDeadline(time.Now().Add(WRITE_TIMEOUT))
	return conn.WriteMessage(websocket.PingMessage, nil)
}

func (s *Session) Disconnect(p *Player, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p.Conn != conn {
		return
	}
	p.Conn = nil

	switch s.State.Current() {
	case WAITING:
		// player who created game left first
		s.State.Event(context.Background(), FINISH)
	case PLAYING:
		// player leaves and timer is activated
		s.State.Event(context.Background(), DISCONNECT)
		s.broadcast(gin.H{"type": "player_disconnected", "seat": s.seatOf(p), "timeout": TIMEOUT_PERIOD.Seconds()})
	}
}

func (s *Session) startGraceTimer() {
	s.graceTimer = time.AfterFunc(TIMEOUT_PERIOD, s.onTimeout)
}

func (s *Session) stopGraceTimer() {
	if s.graceTimer != nil {
		s.graceTimer.Stop()
		s.graceTimer = nil
	}
}

func (s *Session) onTimeout() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State.Current() != PAUSED {
		return
	}

	result := s.Game.OnAbandon(s.connectedSeats())
	s.State.Event(context.Background(), FINISH)
	s.broadcast(gin.H{"type": "finished", "winner": result.Winner, "reason": "abandon"})
}

func (s *Session) send(p *Player, v any) {
	if p.Conn == nil {
		return
	}

	p.Conn.SetWriteDeadline(time.Now().Add(TIMEOUT_PERIOD))
	if err := p.Conn.WriteJSON(v); err != nil {
		log.Printf("session %s: error sending to seat %d: %v", s.Code, s.seatOf(p), err)
	}
}

func (s *Session) broadcast(v any) {
	for _, p := range s.Players {
		s.send(p, v)
	}
}

func (s *Session) HandleMessage(p *Player, raw []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var env ws.Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		s.send(p, gin.H{"type": "error", "message": "invalid message"})
		return
	}

	if s.State.Current() != PLAYING {
		s.send(p, gin.H{"type": "error", "message": "the match is not in progress"})
	}

	res, err := s.Game.HandleMessage(s.seatOf(p), env)
	if err != nil {
		s.send(p, gin.H{"type": "error", "message": err.Error()})
		return
	}

	s.broadcast(gin.H{"type": "state", "data": res.State})
	if res.Finished {
		s.broadcast(gin.H{"type": "finished", "winner": res.Winner, "reason": "game_over"})
		s.State.Event(context.Background(), FINISH)
	}
}
