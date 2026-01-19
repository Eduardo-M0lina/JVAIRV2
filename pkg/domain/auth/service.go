package auth

import (
	"context"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// TokenDetails contiene información sobre el token JWT
type TokenDetails struct {
	AccessToken  string // @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	RefreshToken string // @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	AccessUUID   string // @example "f8776176-9586-4e3c-a767-c011f4d178f8"
	RefreshUUID  string // @example "a9776176-9586-4e3c-a767-c011f4d178f9"
	AtExpires    int64  // @example 1625097600
	RtExpires    int64  // @example 1625184000
}

// AccessDetails contiene información extraída del token JWT
type AccessDetails struct {
	AccessUUID string
	UserID     string
	RoleID     string
}

// RefreshResponse representa la respuesta de refresco de token
type RefreshResponse struct {
	AccessToken  string    `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string    `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt    time.Time `json:"expiresAt" example:"2023-01-01T00:00:00Z"`
}

// LoginRequest representa la solicitud de inicio de sesión
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"admin@example.com"`
	Password string `json:"password" validate:"required" example:"admin123"`
}

// UserResponse representa la respuesta de usuario en el login
type UserResponse struct {
	ID               int64      `json:"id"`
	RoleID           *string    `json:"roleId,omitempty"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	EmailVerifiedAt  *time.Time `json:"emailVerifiedAt,omitempty"`
	IsChangePassword bool       `json:"isChangePassword"`
	IsActive         bool       `json:"isActive"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
	RoleName         *string    `json:"roleName,omitempty"`
	RoleTitle        *string    `json:"roleTitle,omitempty"`
}

// LoginResponse representa la respuesta de inicio de sesión
type LoginResponse struct {
	AccessToken  string        `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string        `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt    time.Time     `json:"expiresAt" example:"2023-01-01T00:00:00Z"`
	User         *UserResponse `json:"user"`
}

// Service define las operaciones relacionadas con la autenticación
type Service interface {
	// Generar tokens JWT para un usuario
	CreateToken(ctx context.Context, user *user.User) (*TokenDetails, error)

	// Extraer información de un token JWT
	ExtractTokenMetadata(ctx context.Context, tokenString string) (*AccessDetails, error)

	// Verificar si un token JWT es válido
	ValidateToken(ctx context.Context, tokenString string) (bool, error)

	// Almacenar información del token en caché/redis
	StoreTokenDetails(ctx context.Context, userID int64, td *TokenDetails) error

	// Eliminar información del token de caché/redis (logout)
	DeleteTokenDetails(ctx context.Context, accessUUID string) error

	// Refrescar token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenDetails, error)
}
