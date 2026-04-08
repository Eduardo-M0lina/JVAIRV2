package router

import (
	propertyHandler "github.com/angumol/jvairv2/pkg/rest/handler/property"
	"github.com/go-chi/chi/v5"
)

// RegisterPropertyRoutes registra las rutas de propiedades
func RegisterPropertyRoutes(r chi.Router, handler *propertyHandler.Handler) {
	r.Route("/properties", func(r chi.Router) {
		handler.RegisterRoutes(r)
	})
}
