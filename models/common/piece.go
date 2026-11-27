package common

type Piece struct {
	Position   Position    `json:"position"`
	LegalMoves []*Position `json:"legal_moves"`
}

type PieceInterface interface {
	GetPosition() Position
	GetLegalMoves(board *Board) []*Position
}

func (p *Piece) GetPosition() Position { return p.Position }

func (p *Piece) Slide(b *Board, dirs []Position, maxSteps int8) []*Position {
	var moves []*Position

	for _, d := range dirs {
		for i := int8(1); maxSteps == 0 || i <= maxSteps; i++ {
			next := Position{X: p.Position.X + d.X*i, Y: p.Position.Y + d.Y*i}
			if !b.InBounds(next) {
				break
			}
			moves = append(moves, &next)
		}
	}

	return moves
}
