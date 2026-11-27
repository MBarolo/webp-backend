package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mbarolo/webp_backend/utils"
)

var upgrader = websocket.Upgrader{}

func WebsocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "Error upgrading websocket connection: "+err.Error(), nil)
		return
	}

	defer conn.Close()
}
