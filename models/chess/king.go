package chess

import "github.com/mbarolo/webp_backend/models/common"

var _ common.PieceInterface = (*King)(nil)

var kingDirs = []common.Position{
	{X: -1, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1},
	{X: 0, Y: -1}, {X: 0, Y: 1},
	{X: 1, Y: -1}, {X: 1, Y: 0}, {X: 1, Y: 1},
}

type King struct {
	ChessPiece
}

func (k *King) GetLegalMoves(b *common.Board) []*common.Position {
	return k.Slide(b, kingDirs, 1)
}
