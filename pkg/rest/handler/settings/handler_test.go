package settings

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domain "github.com/your-org/jvairv2/pkg/domain/settings"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
)

func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.WithTestAbilities(r.Context(), []string{"*"})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Route("/settings", func(r chi.Router) {
		r.Get("/", h.Get)
		r.Put("/", h.Update)
	})
	return r
}

func newHandler() (*Handler, *domain.MockRepository) {
	mockRepo := new(domain.MockRepository)
	uc := domain.NewUseCase(mockRepo)
	handler := NewHandler(uc)
	return handler, mockRepo
}

func TestHandler_Get_Success(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("Get", mock.Anything).Return(&domain.Settings{
		ID: 1, IsTwilioEnabled: false,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Update_InvalidBody(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("Get", mock.Anything).Return(&domain.Settings{
		ID: 1, IsTwilioEnabled: false,
	}, nil)

	req := httptest.NewRequest(http.MethodPut, "/settings", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
