package router

import (
	abilityHandler "github.com/angumol/jvairv2/pkg/rest/handler/ability"
	"github.com/go-chi/chi/v5"
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
