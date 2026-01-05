package router

import (
	"github.com/go-chi/chi/v5"
	userHandler "github.com/your-org/jvairv2/pkg/rest/handler/user"
)

// RegisterUserRoutes registra las rutas de usuarios
func RegisterUserRoutes(r chi.Router, handler *userHandler.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/{id}", handler.Get)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}
