package auth

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authDomain "github.com/your-org/jvairv2/pkg/domain/auth"
	"github.com/your-org/jvairv2/pkg/domain/user"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) CreateToken(ctx context.Context, u *user.User) (*authDomain.TokenDetails, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authDomain.TokenDetails), args.Error(1)
}

func (m *mockAuthService) ExtractTokenMetadata(ctx context.Context, tokenString string) (*authDomain.AccessDetails, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authDomain.AccessDetails), args.Error(1)
}

func (m *mockAuthService) ValidateToken(ctx context.Context, tokenString string) (bool, error) {
	args := m.Called(ctx, tokenString)
	return args.Bool(0), args.Error(1)
}

func (m *mockAuthService) StoreTokenDetails(ctx context.Context, userID int64, td *authDomain.TokenDetails) error {
	args := m.Called(ctx, userID, td)
	return args.Error(0)
}

func (m *mockAuthService) DeleteTokenDetails(ctx context.Context, accessUUID string) error {
	args := m.Called(ctx, accessUUID)
	return args.Error(0)
}

func (m *mockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*authDomain.TokenDetails, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authDomain.TokenDetails), args.Error(1)
}

func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/logout", h.Logout)
		r.Post("/refresh", h.RefreshToken)
	})
	return r
}

func TestHandler_Login_InvalidBody(t *testing.T) {
	mockUserRepo := new(user.MockRepository)
	mockAuthSvc := new(mockAuthService)
	uc := authDomain.NewUseCase(mockUserRepo, mockAuthSvc)
	h := NewHandler(uc)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Login_MissingFields(t *testing.T) {
	mockUserRepo := new(user.MockRepository)
	mockAuthSvc := new(mockAuthService)
	uc := authDomain.NewUseCase(mockUserRepo, mockAuthSvc)
	h := NewHandler(uc)
	router := setupRouter(h)

	body := `{"email":""}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(user.MockRepository)
	mockAuthSvc := new(mockAuthService)
	uc := authDomain.NewUseCase(mockUserRepo, mockAuthSvc)
	h := NewHandler(uc)
	router := setupRouter(h)

	mockUserRepo.On("VerifyCredentials", mock.Anything, "bad@test.com", "wrongpass").Return(nil, assert.AnError)

	body := `{"email":"bad@test.com","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockUserRepo.AssertExpectations(t)
}
