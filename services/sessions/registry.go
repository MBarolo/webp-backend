package sessions

import (
	"sync"

	"github.com/mbarolo/webp_backend/models/common"
	"github.com/mbarolo/webp_backend/utils"
)

const CODE_LENGTH = 6

type SessionRegistry struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

var Registry = &SessionRegistry{sessions: map[string]*Session{}}

func (r *SessionRegistry) Create(game common.Game) *Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	for {
		code := utils.GenerateRandomString(CODE_LENGTH)
		if _, exists := r.sessions[code]; !exists {
			s := NewSession(code, game)
			r.sessions[code] = s
			return s
		}
	}
}

func (r *SessionRegistry) Find(code string) (*Session, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.sessions[code]
	return s, ok
}

func (r *SessionRegistry) Delete(code string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, code)
}
