package common

type Board struct {
	Size int8 `json:"size"`
}

func (b *Board) InBounds(pos Position) bool {
	return pos.X >= 0 && pos.X < b.Size && pos.Y >= 0 && pos.Y < b.Size
}
