package middleware

import (
	"context"
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
		// Obtener token del encabezado Authorization
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

		// Validar token
		valid, err := m.authUseCase.ValidateToken(r.Context(), tokenString)
		if err != nil || !valid {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

		// Obtener usuario del token
		user, err := m.authUseCase.GetUserFromToken(r.Context(), tokenString)
		if err != nil {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}

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
