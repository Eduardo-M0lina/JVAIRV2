package permission

import (
	"time"
)

// Permission representa un permiso en el sistema
type Permission struct {
	ID         int64      // bigint unsigned NOT NULL AUTO_INCREMENT
	AbilityID  int64      // bigint unsigned NOT NULL
	EntityID   int64      // bigint unsigned NOT NULL
	EntityType string     // varchar(191) NOT NULL
	Forbidden  bool       // tinyint(1) NOT NULL DEFAULT '0'
	Conditions *string    // json DEFAULT NULL
	CreatedAt  *time.Time // timestamp NULL DEFAULT NULL
	UpdatedAt  *time.Time // timestamp NULL DEFAULT NULL
}
