package assigned_role

// AssignedRole representa la asignación de un rol a una entidad
type AssignedRole struct {
	ID               int64   // bigint unsigned NOT NULL AUTO_INCREMENT
	RoleID           int64   // bigint unsigned NOT NULL
	EntityID         int64   // bigint unsigned NOT NULL
	EntityType       string  // varchar(191) NOT NULL
	RestrictedToID   *int64  // bigint unsigned DEFAULT NULL
	RestrictedToType *string // varchar(191) DEFAULT NULL
	Scope            *int    // int DEFAULT NULL

	// Campo calculado para compatibilidad
	Restricted bool `db:"-"` // Calculado: true si RestrictedToID != nil
}
