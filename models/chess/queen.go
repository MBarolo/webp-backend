package chess

import "github.com/mbarolo/test_back/models/common"

var _ common.PieceInterface = (*Queen)(nil)

type Queen struct {
	ChessPiece
}

var queenDirs = []common.Position{
	{X: -1, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1},
	{X: 0, Y: -1}, {X: 0, Y: -0}, {X: 0, Y: 1},
	{X: 1, Y: -1}, {X: 1, Y: -0}, {X: 1, Y: 1},
}

func (q *Queen) GetLegalMoves(b *common.Board) []*common.Position {
	return q.Slide(b, queenDirs, 0)
}
