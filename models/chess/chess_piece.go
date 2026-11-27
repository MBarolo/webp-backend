package chess

import "github.com/mbarolo/test_back/models/common"

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
