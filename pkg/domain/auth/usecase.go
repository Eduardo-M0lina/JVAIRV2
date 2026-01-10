package auth

import (
	"context"
	"errors"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/user"
)

var (
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUserInactive       = errors.New("usuario inactivo")
	ErrInvalidToken       = errors.New("token inválido")
)

// UseCase define los casos de uso para autenticación
type UseCase struct {
	userRepo    user.Repository
	authService Service
}

// NewUseCase crea una nueva instancia del caso de uso de autenticación
func NewUseCase(userRepo user.Repository, authService Service) *UseCase {
	return &UseCase{
		userRepo:    userRepo,
		authService: authService,
	}
}

// Login autentica a un usuario y genera tokens JWT
func (uc *UseCase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Verificar credenciales
	u, err := uc.userRepo.VerifyCredentials(ctx, req.Email, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verificar si el usuario está activo
	if !u.IsActive {
		return nil, ErrUserInactive
	}

	// Generar tokens JWT
	td, err := uc.authService.CreateToken(ctx, u)
	if err != nil {
		return nil, err
	}

	// Crear respuesta
	resp := &LoginResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
		ExpiresAt:    time.Unix(td.AtExpires, 0),
		User:         u,
	}

	return resp, nil
}

// Logout cierra la sesión de un usuario
func (uc *UseCase) Logout(ctx context.Context, accessToken string) error {
	// Extraer metadata del token
	ad, err := uc.authService.ExtractTokenMetadata(ctx, accessToken)
	if err != nil {
		return err
	}

	// Eliminar token
	err = uc.authService.DeleteTokenDetails(ctx, ad.AccessUUID)
	if err != nil {
		return err
	}

	return nil
}

// RefreshToken refresca un token JWT
func (uc *UseCase) RefreshToken(ctx context.Context, refreshToken string) (*TokenDetails, error) {
	// Refrescar token
	td, err := uc.authService.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return td, nil
}

// ValidateToken valida un token JWT
func (uc *UseCase) ValidateToken(ctx context.Context, accessToken string) (bool, error) {
	// Validar token
	valid, err := uc.authService.ValidateToken(ctx, accessToken)
	if err != nil {
		return false, err
	}

	return valid, nil
}

// GetUserFromToken obtiene un usuario a partir de un token JWT
func (uc *UseCase) GetUserFromToken(ctx context.Context, accessToken string) (*user.User, error) {
	// Extraer metadata del token
	ad, err := uc.authService.ExtractTokenMetadata(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	// Obtener usuario
	u, err := uc.userRepo.GetByID(ctx, ad.UserID)
	if err != nil {
		return nil, err
	}

	return u, nil
}
