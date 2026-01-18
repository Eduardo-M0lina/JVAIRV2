package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/your-org/jvairv2/pkg/domain/auth"
)

// Tipo personalizado para la clave del contexto
type contextKey string

// Constantes para las claves del contexto
const (
	UserContextKey contextKey = "user_context_key"
)

// AuthMiddleware es un middleware para autenticar solicitudes HTTP
type AuthMiddleware struct {
	authUseCase *auth.UseCase
}

// NewAuthMiddleware crea una nueva instancia del middleware de autenticación
func NewAuthMiddleware(authUseCase *auth.UseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

// Authenticate es un middleware que verifica que el usuario esté autenticado
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Iniciando autenticación",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		// Obtener token del encabezado Authorization
		tokenString := extractToken(r)
		if tokenString == "" {
			slog.Warn("Token no encontrado en encabezado Authorization",
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}
		slog.Debug("Token extraído del encabezado",
			"token_prefix", tokenString[:min(20, len(tokenString))],
		)

		// Validar token
		valid, err := m.authUseCase.ValidateToken(r.Context(), tokenString)
		if err != nil || !valid {
			slog.Error("Token inválido",
				"valid", valid,
				"error", err,
				"path", r.URL.Path,
			)
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}
		slog.Debug("Token validado correctamente")

		// Obtener usuario del token
		user, err := m.authUseCase.GetUserFromToken(r.Context(), tokenString)
		if err != nil {
			slog.Error("No se pudo obtener usuario del token",
				"error", err,
				"path", r.URL.Path,
			)
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}
		slog.Info("Usuario autenticado correctamente",
			"user_id", user.ID,
			"email", user.Email,
			"name", user.Name,
		)

		// Agregar usuario al contexto usando la clave personalizada
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken extrae el token JWT del encabezado Authorization
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// min retorna el menor de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
