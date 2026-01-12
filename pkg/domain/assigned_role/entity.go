package assigned_role

import (
	"time"
)

// AssignedRole representa la asignación de un rol a una entidad
type AssignedRole struct {
	ID               int64   // bigint unsigned NOT NULL AUTO_INCREMENT
	RoleID           int64   // bigint unsigned NOT NULL
	EntityID         int64   // bigint unsigned NOT NULL
	EntityType       string  // varchar(191) NOT NULL
	RestrictedToID   *int64  // bigint unsigned DEFAULT NULL
	RestrictedToType *string // varchar(191) DEFAULT NULL
	Scope            *int    // int DEFAULT NULL

	// Campos calculados para compatibilidad
	Restricted bool       `db:"-"` // Calculado: true si RestrictedToID != nil
	CreatedAt  *time.Time `db:"-"` // No existe en la tabla, solo para compatibilidad con handlers
	UpdatedAt  *time.Time `db:"-"` // No existe en la tabla, solo para compatibilidad con handlers
}
