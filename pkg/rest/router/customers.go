package router

import (
	customerHandler "github.com/angumol/jvairv2/pkg/rest/handler/customer"
	"github.com/go-chi/chi/v5"
)

func RegisterCustomerRoutes(r chi.Router, handler *customerHandler.Handler) {
	handler.RegisterRoutes(r)
}
