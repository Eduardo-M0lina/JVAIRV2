package middleware

import (
	"context"
	"net/http"
)

// abilityKey es la clave para almacenar las habilidades del usuario en el contexto
type abilityKey struct{}

// WithAbilities agrega las habilidades del usuario al contexto
func WithAbilities(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// En un caso real, aquí se obtendrían las habilidades del usuario
		// desde el token JWT o desde la base de datos
		abilities := []string{
			"create_user",
			"view_user",
			"update_user",
			"delete_user",
			"list_users",
			"view_user_roles",
			"view_user_abilities",
		}

		// Agregar las habilidades al contexto
		ctx := context.WithValue(r.Context(), abilityKey{}, abilities)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HasAbility verifica si el usuario tiene una habilidad específica
func HasAbility(ctx context.Context, ability string) bool {
	// Obtener las habilidades del contexto
	abilities, ok := ctx.Value(abilityKey{}).([]string)
	if !ok {
		return false
	}

	// Verificar si el usuario tiene la habilidad
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
