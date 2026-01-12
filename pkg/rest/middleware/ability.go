package middleware

import (
	"context"
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
			// Obtener el usuario del contexto
			userCtx := r.Context().Value(UserContextKey)
			if userCtx == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Convertir el usuario a un objeto User
			u, ok := userCtx.(*user.User)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			// Obtener las habilidades del usuario desde la base de datos
			abilities, err := userUseCase.GetUserAbilities(r.Context(), strconv.FormatInt(u.ID, 10))
			if err != nil {
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
