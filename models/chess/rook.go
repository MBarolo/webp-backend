package chess

import "github.com/mbarolo/webp_backend/models/common"

var _ common.PieceInterface = (*Rook)(nil)

type Rook struct {
	ChessPiece
}

var rookDirs = []common.Position{
	{X: 1, Y: 0}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1},
}

func (r *Rook) GetLegalMoves(b *common.Board) []*common.Position {
	return r.Slide(b, rookDirs, 0)
}
