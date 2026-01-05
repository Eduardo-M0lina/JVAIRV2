package router

import (
	"github.com/go-chi/chi/v5"
	roleHandler "github.com/your-org/jvairv2/pkg/rest/handler/role"
)

// RegisterRoleRoutes registra las rutas de roles
func RegisterRoleRoutes(r chi.Router, handler *roleHandler.Handler) {
	r.Route("/roles", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/{id}", handler.Get)
	})
}
