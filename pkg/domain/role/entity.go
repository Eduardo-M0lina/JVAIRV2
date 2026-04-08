package role

import (
	"time"
)

// Role representa la entidad de dominio para un rol de usuario
type Role struct {
	ID        int64      // bigint unsigned NOT NULL AUTO_INCREMENT
	Name      string     // varchar(191) NOT NULL
	Title     *string    // varchar(191) DEFAULT NULL
	Scope     *int       // int DEFAULT NULL
	CreatedAt *time.Time // timestamp NULL DEFAULT NULL
	UpdatedAt *time.Time // timestamp NULL DEFAULT NULL
}
