package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mbarolo/webp_backend/controller"
)

func InitWebsocketRoutes(r *gin.RouterGroup) {
	ws := r.Group("/ws/:game")
	ws.GET("/new", controller.CreateGameHandler)
	ws.GET("/join/:code", controller.JoinGameHandler)
}
