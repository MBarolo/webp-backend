package chess

import "github.com/mbarolo/webp_backend/models/common"

// Standard will be:

// Whites will always advance towards top
// Blacks will always advance towards bottom

type Kind uint8

const (
	WHITE Kind = iota
	BLACK
)

type ChessPiece struct {
	common.Piece
	Kind Kind `json:"kind"`
}

func NewChessPiece(kind Kind, x, y int8) ChessPiece {
	return ChessPiece{
		Piece: common.Piece{Position: common.Position{X: x, Y: y}},
		Kind:  kind,
	}
}
