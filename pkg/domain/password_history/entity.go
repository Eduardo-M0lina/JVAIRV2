package password_history

import "time"

// PasswordHistory representa un registro de historial de contraseñas
type PasswordHistory struct {
	ID        int64      // bigint unsigned NOT NULL AUTO_INCREMENT
	UserID    int64      // bigint unsigned NOT NULL
	Password  string     // varchar(255) NOT NULL (hash)
	CreatedAt *time.Time // timestamp NULL DEFAULT NULL
	UpdatedAt *time.Time // timestamp NULL DEFAULT NULL
}
