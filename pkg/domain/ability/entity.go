package ability

import (
	"time"
)

// Ability representa una capacidad o permiso en el sistema
type Ability struct {
	ID         int64      // bigint unsigned NOT NULL AUTO_INCREMENT
	Name       string     // varchar(191) NOT NULL
	Title      *string    // varchar(191) DEFAULT NULL
	EntityID   *int64     // bigint unsigned DEFAULT NULL
	EntityType *string    // varchar(191) DEFAULT NULL
	OnlyOwned  bool       // tinyint(1) NOT NULL DEFAULT '0'
	Options    *string    // json DEFAULT NULL
	Scope      *int       // int DEFAULT NULL
	CreatedAt  *time.Time // timestamp NULL DEFAULT NULL
	UpdatedAt  *time.Time // timestamp NULL DEFAULT NULL
}
