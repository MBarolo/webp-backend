package matches

import (
	"context"

	"github.com/looplab/fsm"
)

type MatchEvent = string

const (
	START      MatchEvent = "start"
	DISCONNECT MatchEvent = "disconnect"
	RECONNECT  MatchEvent = "reconnect"
	FINISH     MatchEvent = "finish"
)

type MatchStates = string

const (
	WAITING  MatchStates = "waiting"
	PLAYING  MatchStates = "playing"
	PAUSED   MatchStates = "paused"
	FINISHED MatchStates = "finished"
)

type MatchActions = string

const (
	ENTER_PAUSED = "enter_" + PAUSED
	LEAVE_PAUSED = "leave_" + PAUSED
)

func GetStates(m *Match) *fsm.FSM {
	states := fsm.NewFSM(
		string(WAITING),
		fsm.Events{
			{Name: START, Src: []string{WAITING}, Dst: PLAYING},
			{Name: DISCONNECT, Src: []string{PLAYING}, Dst: PAUSED},
			{Name: RECONNECT, Src: []string{PAUSED}, Dst: PLAYING},
			{Name: FINISH, Src: []string{PLAYING, WAITING, PAUSED}, Dst: FINISHED},
		},
		fsm.Callbacks{
			ENTER_PAUSED: func(ctx context.Context, e *fsm.Event) { m.startGraceTimer() },
			LEAVE_PAUSED: func(ctx context.Context, e *fsm.Event) { m.stopGraceTimer() },
		},
	)

	return states
}
