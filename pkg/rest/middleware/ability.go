package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// abilityKey es la clave para almacenar las habilidades del usuario en el contexto
type abilityKey struct{}

// WithAbilities agrega las habilidades del usuario al contexto
func WithAbilities(userUseCase *user.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("Iniciando carga de habilidades",
				"method", r.Method,
				"path", r.URL.Path,
			)

			// Obtener el usuario del contexto
			userCtx := r.Context().Value(UserContextKey)
			if userCtx == nil {
				slog.Warn("No hay usuario en el contexto",
					"path", r.URL.Path,
				)
				next.ServeHTTP(w, r)
				return
			}

			// Convertir el usuario a un objeto User
			u, ok := userCtx.(*user.User)
			if !ok {
				slog.Error("No se pudo convertir el usuario del contexto",
					"type", userCtx,
				)
				next.ServeHTTP(w, r)
				return
			}

			// Obtener las habilidades del usuario desde la base de datos
			abilities, err := userUseCase.GetUserAbilities(r.Context(), strconv.FormatInt(u.ID, 10))
			if err != nil {
				slog.Error("Error al obtener habilidades del usuario",
					"user_id", u.ID,
					"email", u.Email,
					"error", err,
				)
				next.ServeHTTP(w, r)
				return
			}

			slog.Info("Habilidades cargadas correctamente",
				"user_id", u.ID,
				"email", u.Email,
				"abilities_count", len(abilities),
			)

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
	slog.Debug("Verificando permiso", "ability", ability)

	// Obtener las habilidades del contexto
	abilities, ok := ctx.Value(abilityKey{}).([]string)
	if !ok {
		slog.Warn("No se encontraron habilidades en el contexto")
		return false
	}

	// Verificar si el usuario tiene la habilidad "*" (superadmin)
	for _, a := range abilities {
		if a == "*" {
			slog.Info("Acceso concedido por habilidad wildcard",
				"requested_ability", ability,
			)
			return true
		}
	}

	// Verificar si el usuario tiene la habilidad específica
	for _, a := range abilities {
		if a == ability {
			slog.Info("Acceso concedido",
				"ability", ability,
			)
			return true
		}
	}

	slog.Warn("Acceso denegado - habilidad no encontrada",
		"ability", ability,
		"available_abilities", abilities,
	)
	return false
}

// RequireAbility verifica si el usuario tiene una habilidad específica
// y devuelve un error 403 si no la tiene
func RequireAbility(ability string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasAbility(r.Context(), ability) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
