package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/domain/user"
)

// abilityKey es la clave para almacenar las habilidades del usuario en el contexto
type abilityKey struct{}

// WithAbilities agrega las habilidades del usuario al contexto
func WithAbilities(userUseCase *user.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Obtener el usuario del contexto
			userCtx := r.Context().Value(UserContextKey)
			if userCtx == nil {
				slog.Error("No hay usuario en el contexto - esto no debería pasar si AuthMiddleware funciona correctamente",
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				http.Error(w, "Error interno - Usuario no encontrado en contexto", http.StatusInternalServerError)
				return
			}

			// Convertir el usuario a un objeto User
			u, ok := userCtx.(*user.User)
			if !ok {
				slog.Error("No se pudo convertir el usuario del contexto",
					"type", userCtx,
					"path", r.URL.Path,
				)
				http.Error(w, "Error interno - Tipo de usuario inválido", http.StatusInternalServerError)
				return
			}

			// Obtener las habilidades del usuario desde la base de datos
			abilities, err := userUseCase.GetUserAbilities(r.Context(), strconv.FormatInt(u.ID, 10))
			if err != nil {
				slog.Error("Error al obtener habilidades del usuario",
					"user_id", u.ID,
					"email", u.Email,
					"error", err,
					"path", r.URL.Path,
				)
				// No bloqueamos la petición si hay error en habilidades, solo registramos
				next.ServeHTTP(w, r)
				return
			}

			// Convertir las habilidades a un slice de strings
			abilityNames := make([]string, len(abilities))
			for i, a := range abilities {
				abilityNames[i] = a.Name
			}

			// Agregar las habilidades al contexto
			ctx := context.WithValue(r.Context(), abilityKey{}, abilityNames)

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// HasAbility verifica si el usuario tiene una habilidad específica
func HasAbility(ctx context.Context, ability string) bool {

	// Obtener las habilidades del contexto
	abilities, ok := ctx.Value(abilityKey{}).([]string)
	if !ok {
		slog.Error("No se encontraron habilidades en el contexto",
			"requested_ability", ability,
			"context_type", fmt.Sprintf("%T", ctx.Value(abilityKey{})),
		)
		return false
	}

	// Verificar si el usuario tiene la habilidad "*" (superadmin)
	for _, a := range abilities {
		if a == "*" {
			return true
		}
	}

	// Verificar si el usuario tiene la habilidad específica
	for _, a := range abilities {
		if a == ability {
			return true
		}
	}

	return false
}

// HasAbilityExact verifica si el usuario tiene una habilidad específica SIN considerar el wildcard "*"
// Útil para verificar restricciones como "job_view_user_only" donde el superadmin NO debe tener la restricción
func HasAbilityExact(ctx context.Context, ability string) bool {
	// Obtener las habilidades del contexto
	abilities, ok := ctx.Value(abilityKey{}).([]string)
	if !ok {
		return false
	}

	// Verificar si el usuario tiene la habilidad específica (sin considerar "*")
	for _, a := range abilities {
		if a == ability {
			return true
		}
	}

	return false
}

// RequireAbility verifica si el usuario tiene una habilidad específica
// y devuelve un error 403 si no la tiene
func RequireAbility(ability string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasAbility(r.Context(), ability) {
				slog.Warn("Acceso denegado - habilidad requerida no encontrada",
					"required_ability", ability,
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				http.Error(w, "Forbidden - Permiso insuficiente", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
