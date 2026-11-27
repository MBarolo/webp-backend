package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mbarolo/test_back/utils"
)

func InitRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		InitUserRoutes(r)

		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			utils.JsonResponse(w, http.StatusOK, "ok", nil)
		})
	})
}
