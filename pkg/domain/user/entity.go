package user

import (
	"time"
)

// User representa la entidad de dominio para un usuario
type User struct {
	ID               int64      // bigint unsigned NOT NULL AUTO_INCREMENT
	RoleID           *string    // varchar(191) DEFAULT NULL
	Name             string     // varchar(255) NOT NULL
	Email            string     // varchar(255) NOT NULL
	EmailVerifiedAt  *time.Time // timestamp NULL DEFAULT NULL
	Password         string     // varchar(255) NOT NULL
	IsChangePassword bool       // tinyint(1) NOT NULL DEFAULT '0'
	RememberToken    *string    // varchar(100) DEFAULT NULL
	IsActive         bool       // tinyint(1) NOT NULL DEFAULT '1'
	CreatedAt        *time.Time // timestamp NULL DEFAULT NULL
	UpdatedAt        *time.Time // timestamp NULL DEFAULT NULL
	DeletedAt        *time.Time // timestamp NULL DEFAULT NULL
}
