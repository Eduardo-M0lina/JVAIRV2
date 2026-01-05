package user

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/ability"
	"github.com/your-org/jvairv2/pkg/domain/role"
)

// Repository define las operaciones de persistencia para usuarios
type Repository interface {
	// Obtener un usuario por ID
	GetByID(ctx context.Context, id string) (*User, error)

	// Obtener un usuario por email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Crear un nuevo usuario
	Create(ctx context.Context, user *User) error

	// Actualizar un usuario existente
	Update(ctx context.Context, user *User) error

	// Eliminar un usuario (soft delete)
	Delete(ctx context.Context, id string) error

	// Listar usuarios con paginación y filtros
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*User, int, error)

	// Verificar credenciales de usuario
	VerifyCredentials(ctx context.Context, email, password string) (*User, error)

	// Obtener roles de un usuario
	GetUserRoles(ctx context.Context, userID string) ([]*role.Role, error)

	// Obtener habilidades de un usuario
	GetUserAbilities(ctx context.Context, userID string) ([]*ability.Ability, error)
}
