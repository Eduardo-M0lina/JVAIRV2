package account

// UpdateProfileRequest representa la solicitud para actualizar el perfil del usuario
type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// ChangePasswordRequest representa la solicitud para cambiar la contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// ProfileResponse representa la respuesta del perfil del usuario
type ProfileResponse struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	RoleID           *string `json:"role_id,omitempty"`
	RoleName         *string `json:"role_name,omitempty"`
	RoleTitle        *string `json:"role_title,omitempty"`
	IsActive         bool    `json:"is_active"`
	IsChangePassword bool    `json:"is_change_password"`
}
