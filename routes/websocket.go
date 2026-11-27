package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mbarolo/webp_backend/controller"
)

func InitWebsocketRoutes(group *gin.RouterGroup) {
	group.GET("/ws", controller.WebsocketHandler)
}
