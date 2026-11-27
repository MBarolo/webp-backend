package chess

import "github.com/mbarolo/webp_backend/models/common"

var _ common.PieceInterface = (*Bishop)(nil)

type Bishop struct {
	ChessPiece
}

var bishopDirs = []common.Position{
	{X: -1, Y: -1}, {X: -1, Y: 1}, {X: 1, Y: -1}, {X: 1, Y: 1},
}

func (bi *Bishop) GetLegalMoves(b *common.Board) []*common.Position {
	return bi.Slide(b, bishopDirs, 0)
}
