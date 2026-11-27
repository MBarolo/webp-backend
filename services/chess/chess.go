package chess

import (
	"github.com/mbarolo/webp_backend/models/chess"
	"github.com/mbarolo/webp_backend/models/common"
)

type pieceFactory func(chess.ChessPiece) common.PieceInterface

func rook(c chess.ChessPiece) common.PieceInterface   { return &chess.Rook{ChessPiece: c} }
func knight(c chess.ChessPiece) common.PieceInterface { return &chess.Knight{ChessPiece: c} }
func bishop(c chess.ChessPiece) common.PieceInterface { return &chess.Bishop{ChessPiece: c} }
func queen(c chess.ChessPiece) common.PieceInterface  { return &chess.Queen{ChessPiece: c} }
func king(c chess.ChessPiece) common.PieceInterface   { return &chess.King{ChessPiece: c} }

var backRank = [common.CHESS_BOARD_SIZE]pieceFactory{
	rook, knight, bishop, queen, king, bishop, knight, rook,
}

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
