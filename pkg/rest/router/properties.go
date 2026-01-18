package router

import (
	"github.com/go-chi/chi/v5"
	propertyHandler "github.com/your-org/jvairv2/pkg/rest/handler/property"
)

// RegisterPropertyRoutes registra las rutas de propiedades
func RegisterPropertyRoutes(r chi.Router, handler *propertyHandler.Handler) {
	r.Route("/properties", func(r chi.Router) {
		handler.RegisterRoutes(r)
	})
}
