package password_reset

import "time"

// PasswordReset representa un token de reseteo de contraseña
type PasswordReset struct {
	Email     string     // varchar(255) NOT NULL
	Token     string     // varchar(255) NOT NULL
	CreatedAt *time.Time // timestamp NULL DEFAULT NULL
}
