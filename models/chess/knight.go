package chess

import "github.com/mbarolo/webp_backend/models/common"

var _ common.PieceInterface = (*Knight)(nil)

type Knight struct {
	ChessPiece
}

func (k *Knight) GetLegalMoves(b *common.Board) []*common.Position {
	x := []int8{
		2, 1, -1, -2, -2, -1, 1, 2,
	}
	y := []int8{
		1, 2, 2, 1, -1, -2, -2, -1,
	}

	var moves []*common.Position
	for idx, _ := range x {
		next := &common.Position{
			X: k.Position.X + x[idx],
			Y: k.Position.Y + y[idx],
		}
		if b.InBounds(*next) {
			moves = append(moves, next)
		}
	}

	return moves
}
