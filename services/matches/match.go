package matches

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/looplab/fsm"
	"github.com/mbarolo/webp_backend/models/common"
)

const TIMEOUT_PERIOD = 1 * time.Minute

type Match struct {
	mu         sync.Mutex
	Code       string
	Game       common.Game
	Players    []*Player
	State      *fsm.FSM
	graceTimer *time.Timer
}

func (m *Match) startGraceTimer() {
	m.graceTimer = time.AfterFunc(TIMEOUT_PERIOD, m.onTimeout)
}

func (m *Match) stopGraceTimer() {
	if m.graceTimer != nil {
		m.graceTimer.Stop()
		m.graceTimer = nil
	}
}

func (m *Match) onTimeout() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.State.Current() != PAUSED {
		return
	}

	result := m.Game.OnAbandon(m.connectedSeats())
	m.State.Event(context.Background(), FINISH)
	m.broadcast(gin.H{"type": "finished", "winner": result.Winner, "reason": "abandon"})
}

func (m *Match) connectedSeats() []int {
	var seats []int
	for seat, player := range m.Players {
		if player.Conn != nil {
			seats = append(seats, seat)
		}
	}

	return seats
}

func (m *Match) broadcast(v any) {
	for _, p := range m.Players {
		m.send(p, v)
	}
}

func (m *Match) send(p *Player, v any) {
	if p.Conn == nil {
		return
	}

	p.Conn.SetWriteDeadline(time.Now())
}
