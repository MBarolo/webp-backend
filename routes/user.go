package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mbarolo/webp_backend/controller"
	"github.com/mbarolo/webp_backend/middleware"
)

func InitUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")

	users.GET("/profile", middleware.AuthMiddleware, controller.GetProfile)
	users.PATCH("/profile", middleware.AuthMiddleware, controller.UpdateProfile)
}
