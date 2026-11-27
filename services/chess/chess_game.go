package chess

import (
	"github.com/mbarolo/webp_backend/models/chess"
	"github.com/mbarolo/webp_backend/models/common"
	"github.com/mbarolo/webp_backend/models/ws"
)

var _ common.Game = (*ChessGame)(nil)

type ChessGame struct {
	board common.Board
	turn  chess.Kind
}

type ChessGameState struct {
	Board common.Board `json:"board"`
	Turn  chess.Kind   `json:"kind"`
}

type move struct {
	From common.Position `json:"from"`
	To   common.Position `json:"to"`
}

func NewGame() common.Game {
	return &ChessGame{board: InitializeChessBoard(), turn: chess.WHITE}
}

func (cg *ChessGame) MaxPlayers() int { return 2 }

func (cg *ChessGame) State() any { return ChessGameState{Board: cg.board, Turn: cg.turn} }

func (cg *ChessGame) HandleMessage(seat int, msg ws.Envelope) (common.Result, error) {
	return common.Result{}, nil
}

func (cg *ChessGame) OnAbandon(connectedSeats []int) common.Result {
	if len(connectedSeats) == 1 {
		return common.Result{Finished: true, Winner: connectedSeats[0]}
	}
	return common.Result{Finished: true, Winner: -1}
}
