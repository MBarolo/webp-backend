package services

import (
	"github.com/mbarolo/webp_backend/models/common"
	"github.com/mbarolo/webp_backend/services/chess"
)

type Games = string

const (
	CHESS Games = "chess"
)

var GameRegistry = map[string]func() common.Game{CHESS: chess.NewGame}
