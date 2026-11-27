package matches

import (
	"sync"
)

type MatchRegistry struct {
	mu      sync.Mutex
	matches map[string]*Match
}
