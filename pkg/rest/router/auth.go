package router

import (
	authHandler "github.com/angumol/jvairv2/pkg/rest/handler/auth"
	"github.com/go-chi/chi/v5"
)

// RegisterAuthRoutes registra las rutas de autenticación
func RegisterAuthRoutes(r chi.Router, handler *authHandler.Handler, passwordSecurityHandler *authHandler.PasswordSecurityHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
		r.Post("/refresh", handler.RefreshToken)
		r.Post("/forgot-password", passwordSecurityHandler.ForgotPassword)
		r.Post("/reset-password", passwordSecurityHandler.ResetPassword)
	})
}
