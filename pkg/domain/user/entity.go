package user

import (
	"time"
)

// User representa la entidad de dominio para un usuario
type User struct {
	ID              int64      // bigint
	Name            string     // varchar
	Email           string     // varchar
	Password        string     // varchar
	RoleID          string     // varchar
	EmailVerifiedAt *time.Time // timestamp
	RememberToken   *string    // varchar
	CreatedAt       *time.Time // timestamp
	UpdatedAt       *time.Time // timestamp
	DeletedAt       *time.Time // timestamp
	// Campos virtuales (no existen en la base de datos)
	IsActive bool // Calculado basado en DeletedAt == nil
}

// Role representa la entidad de dominio para un rol de usuario
type Role struct {
	ID        string     // varchar
	Label     string     // varchar
	CreatedAt *time.Time // timestamp
	UpdatedAt *time.Time // timestamp
}

// Ability representa una capacidad o permiso en el sistema
type Ability struct {
	ID          string
	Name        string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
