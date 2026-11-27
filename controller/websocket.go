package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mbarolo/webp_backend/models"
	"github.com/mbarolo/webp_backend/utils"
)

var upgrader = websocket.Upgrader{}

func CreateGameHandler(c *gin.Context) {
	game, ok := models.Games[c.Param("game")]
	if !ok {
		utils.JsonResponse(c, http.StatusNotFound, "Game not found", nil)
		return
	}

	upgrade, 
}
