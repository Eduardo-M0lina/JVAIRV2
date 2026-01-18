package router

import (
	"github.com/go-chi/chi/v5"
	customerHandler "github.com/your-org/jvairv2/pkg/rest/handler/customer"
)

func RegisterCustomerRoutes(r chi.Router, handler *customerHandler.Handler) {
	handler.RegisterRoutes(r)
}
