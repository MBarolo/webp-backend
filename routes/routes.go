package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	InitUserRoutes(api)
	InitWebsocketRoutes(api)
}
