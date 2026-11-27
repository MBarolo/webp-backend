package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mbarolo/webp_backend/controller"
	"github.com/mbarolo/webp_backend/middleware"
)

func InitUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/profile", controller.GetProfile)
		r.Patch("/profile", controller.UpdateProfile)
	})
}
