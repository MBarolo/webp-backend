package sessions

import (
	"context"

	"github.com/looplab/fsm"
)

type SessionEvent = string

const (
	START      SessionEvent = "start"
	DISCONNECT SessionEvent = "disconnect"
	RECONNECT  SessionEvent = "reconnect"
	FINISH     SessionEvent = "finish"
)

type SessionStates = string

const (
	WAITING  SessionStates = "waiting"
	PLAYING  SessionStates = "playing"
	PAUSED   SessionStates = "paused"
	FINISHED SessionStates = "finished"
)

type SessionActions = string

const (
	ENTER_PAUSED   = "enter_" + PAUSED
	LEAVE_PAUSED   = "leave_" + PAUSED
	ENTER_FINISHED = "enter_" + FINISHED
)

func GetStates(s *Session) *fsm.FSM {
	states := fsm.NewFSM(
		string(WAITING),
		fsm.Events{
			{Name: START, Src: []string{WAITING}, Dst: PLAYING},
			{Name: DISCONNECT, Src: []string{PLAYING}, Dst: PAUSED},
			{Name: RECONNECT, Src: []string{PAUSED}, Dst: PLAYING},
			{Name: FINISH, Src: []string{PLAYING, WAITING, PAUSED}, Dst: FINISHED},
		},
		fsm.Callbacks{
			ENTER_PAUSED:   func(ctx context.Context, e *fsm.Event) { s.startGraceTimer() },
			LEAVE_PAUSED:   func(ctx context.Context, e *fsm.Event) { s.stopGraceTimer() },
			ENTER_FINISHED: func(ctx context.Context, e *fsm.Event) { Registry.Delete(s.Code) },
		},
	)

	return states
}
