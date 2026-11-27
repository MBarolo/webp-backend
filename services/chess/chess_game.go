package chess

import (
	"github.com/mbarolo/webp_backend/models/chess"
	"github.com/mbarolo/webp_backend/models/common"
	"github.com/mbarolo/webp_backend/models/ws"
)

type pieceFactory func(chess.ChessPiece) common.PieceInterface

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

var _ common.Game = (*ChessGame)(nil)

var backRank = [common.CHESS_BOARD_SIZE]pieceFactory{
	rook, knight, bishop, queen, king, bishop, knight, rook,
}

func rook(c chess.ChessPiece) common.PieceInterface   { return &chess.Rook{ChessPiece: c} }
func knight(c chess.ChessPiece) common.PieceInterface { return &chess.Knight{ChessPiece: c} }
func bishop(c chess.ChessPiece) common.PieceInterface { return &chess.Bishop{ChessPiece: c} }
func queen(c chess.ChessPiece) common.PieceInterface  { return &chess.Queen{ChessPiece: c} }
func king(c chess.ChessPiece) common.PieceInterface   { return &chess.King{ChessPiece: c} }
func InitializeChessBoard() common.Board {
	board := common.Board{Size: common.CHESS_BOARD_SIZE - 1}
	last := int8(common.CHESS_BOARD_SIZE - 1)

	for x := range int8(common.CHESS_BOARD_SIZE) {
		board.Pieces = append(board.Pieces,
			backRank[x](chess.NewChessPiece(chess.WHITE, x, 0)),
			&chess.Pawn{ChessPiece: chess.NewChessPiece(chess.WHITE, x, 1)},
			&chess.Pawn{ChessPiece: chess.NewChessPiece(chess.BLACK, x, last-1)},
			backRank[x](chess.NewChessPiece(chess.BLACK, x, last)),
		)
	}

	return board
}

func (cg *ChessGame) MaxPlayers() int { return 2 }
func (cg *ChessGame) State() any      { return ChessGameState{Board: cg.board, Turn: cg.turn} }
func (cg *ChessGame) HandleMessage(seat int, msg ws.Envelope) (common.Result, error) {
	return common.Result{}, nil
}
func (cg *ChessGame) OnAbandon(connectedSeats []int) common.Result {
	if len(connectedSeats) == 1 {
		return common.Result{Finished: true, Winner: connectedSeats[0]}
	}
	return common.Result{Finished: true, Winner: -1}
}

func NewGame() common.Game {
	return &ChessGame{board: InitializeChessBoard(), turn: chess.WHITE}
}
