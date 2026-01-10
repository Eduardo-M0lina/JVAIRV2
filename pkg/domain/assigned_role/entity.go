package assigned_role

import (
	"time"
)

// AssignedRole representa la asignación de un rol a una entidad
type AssignedRole struct {
	ID         int64      // bigint unsigned NOT NULL AUTO_INCREMENT
	RoleID     int64      // bigint unsigned NOT NULL
	EntityID   int64      // bigint unsigned NOT NULL
	EntityType string     // varchar(191) NOT NULL
	Restricted bool       // tinyint(1) NOT NULL DEFAULT '0'
	Scope      *int       // int DEFAULT NULL
	CreatedAt  *time.Time // timestamp NULL DEFAULT NULL
	UpdatedAt  *time.Time // timestamp NULL DEFAULT NULL
}
