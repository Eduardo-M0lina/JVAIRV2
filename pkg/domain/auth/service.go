package auth

import (
	"context"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

// TokenDetails contiene información sobre el token JWT
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// AccessDetails contiene información extraída del token JWT
type AccessDetails struct {
	AccessUUID string
	UserID     string
	RoleID     string
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

// LoginRequest representa la solicitud de inicio de sesión
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse representa la respuesta de inicio de sesión
type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresAt    time.Time  `json:"expires_at"`
	User         *user.User `json:"user"`
}
