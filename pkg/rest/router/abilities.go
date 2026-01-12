package router

import (
	"github.com/go-chi/chi/v5"
	abilityHandler "github.com/your-org/jvairv2/pkg/rest/handler/ability"
)

// RegisterAbilityRoutes registra las rutas de abilities
func RegisterAbilityRoutes(r chi.Router, handler *abilityHandler.Handler) {
	r.Route("/abilities", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Get("/{id}", handler.Get)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})
}
