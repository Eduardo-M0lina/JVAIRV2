package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/your-org/jvairv2/pkg/domain/auth"
	"github.com/your-org/jvairv2/pkg/domain/user"
)

var (
	ErrInvalidToken = errors.New("token inválido")
	ErrExpiredToken = errors.New("token expirado")
)

// JWTService implementa la interfaz auth.Service para JWT
type JWTService struct {
	accessSecret  string
	refreshSecret string
	accessExp     time.Duration
	refreshExp    time.Duration
	tokenStore    TokenStore
}

// TokenStore define la interfaz para almacenar y recuperar tokens
type TokenStore interface {
	StoreToken(ctx context.Context, userID int64, tokenID string, expiration time.Duration) error
	DeleteToken(ctx context.Context, tokenID string) error
	CheckToken(ctx context.Context, tokenID string) (bool, error)
}

// NewJWTService crea una nueva instancia del servicio JWT
func NewJWTService(accessSecret, refreshSecret string, accessExp, refreshExp time.Duration, tokenStore TokenStore) *JWTService {
	return &JWTService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExp:     accessExp,
		refreshExp:    refreshExp,
		tokenStore:    tokenStore,
	}
}

// CreateToken genera tokens JWT para un usuario
func (s *JWTService) CreateToken(ctx context.Context, u *user.User) (*auth.TokenDetails, error) {
	td := &auth.TokenDetails{
		AtExpires:   time.Now().Add(s.accessExp).Unix(),
		RtExpires:   time.Now().Add(s.refreshExp).Unix(),
		AccessUUID:  fmt.Sprintf("%d-%d", u.ID, time.Now().Unix()),
		RefreshUUID: fmt.Sprintf("%d-%d-refresh", u.ID, time.Now().Unix()),
	}

	// Crear token de acceso
	atClaims := jwt.MapClaims{
		"user_id":     u.ID,
		"role_id":     u.RoleID,
		"access_uuid": td.AccessUUID,
		"exp":         td.AtExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, err
	}
	td.AccessToken = accessToken

	// Crear token de refresco
	rtClaims := jwt.MapClaims{
		"user_id":      u.ID,
		"refresh_uuid": td.RefreshUUID,
		"exp":          td.RtExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return nil, err
	}
	td.RefreshToken = refreshToken

	// Almacenar tokens en caché/redis
	err = s.StoreTokenDetails(ctx, u.ID, td)
	if err != nil {
		return nil, err
	}

	return td, nil
}

// ExtractTokenMetadata extrae información del token JWT
func (s *JWTService) ExtractTokenMetadata(ctx context.Context, tokenString string) (*auth.AccessDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	userID := int64(userIDFloat)

	roleID, ok := claims["role_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Verificar si el token existe en la caché/redis
	exists, err := s.tokenStore.CheckToken(ctx, accessUUID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrInvalidToken
	}

	return &auth.AccessDetails{
		AccessUUID: accessUUID,
		UserID:     fmt.Sprintf("%d", userID),
		RoleID:     roleID,
	}, nil
}

// ValidateToken verifica si un token JWT es válido
func (s *JWTService) ValidateToken(ctx context.Context, tokenString string) (bool, error) {
	_, err := s.ExtractTokenMetadata(ctx, tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}

// StoreTokenDetails almacena información del token en caché/redis
func (s *JWTService) StoreTokenDetails(ctx context.Context, userID int64, td *auth.TokenDetails) error {
	// Almacenar token de acceso
	atExpiration := time.Until(time.Unix(td.AtExpires, 0))
	err := s.tokenStore.StoreToken(ctx, userID, td.AccessUUID, atExpiration)
	if err != nil {
		return err
	}

	// Almacenar token de refresco
	rtExpiration := time.Until(time.Unix(td.RtExpires, 0))
	err = s.tokenStore.StoreToken(ctx, userID, td.RefreshUUID, rtExpiration)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTokenDetails elimina información del token de caché/redis
func (s *JWTService) DeleteTokenDetails(ctx context.Context, accessUUID string) error {
	return s.tokenStore.DeleteToken(ctx, accessUUID)
}

// RefreshToken refresca un token JWT
func (s *JWTService) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenDetails, error) {
	// Validar token de refresco
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	refreshUUID, ok := claims["refresh_uuid"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	userID := int64(userIDFloat)

	// Verificar si el token existe en la caché/redis
	exists, err := s.tokenStore.CheckToken(ctx, refreshUUID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrInvalidToken
	}

	// Eliminar el token de refresco actual
	err = s.tokenStore.DeleteToken(ctx, refreshUUID)
	if err != nil {
		return nil, err
	}

	// Crear un nuevo usuario temporal para generar nuevos tokens
	u := &user.User{
		ID: userID,
	}

	// Generar nuevos tokens
	return s.CreateToken(ctx, u)
}
