package chess

import "github.com/mbarolo/test_back/models/common"

var _ common.PieceInterface = (*Pawn)(nil)

type Pawn struct {
	ChessPiece
}

func (p *Pawn) GetLegalMoves(b *common.Board) []*common.Position {
	var dir, startRow int8

	switch p.Kind {
	case WHITE:
		dir, startRow = 1, 1
	case BLACK:
		dir, startRow = -1, b.Size-2
	}

	steps := int8(1)
	if p.Position.Y == startRow {
		steps = 2
	}

	var moves []*common.Position
	for i := int8(1); i <= steps; i++ {
		next := common.Position{
			X: p.Position.X,
			Y: p.Position.Y + dir*i,
		}
		if !b.InBounds(next) {
			break
		}
		moves = append(moves, &next)
	}

	return moves
}
