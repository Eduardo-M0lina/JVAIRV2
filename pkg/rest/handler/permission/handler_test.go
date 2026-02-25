package permission

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	abilityDomain "github.com/your-org/jvairv2/pkg/domain/ability"
	domain "github.com/your-org/jvairv2/pkg/domain/permission"
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
	r.Route("/permissions", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Get("/ability/{abilityId}", h.GetByAbility)
		r.Get("/entity/{entityType}/{entityId}", h.GetByEntity)
		r.Get("/check/{abilityId}/{entityType}/{entityId}", h.Exists)
	})
	return r
}

func newHandler() (*Handler, *domain.MockRepository, *abilityDomain.MockRepository) {
	mockRepo := new(domain.MockRepository)
	mockAbilityRepo := new(abilityDomain.MockRepository)
	uc := domain.NewUseCase(mockRepo, mockAbilityRepo)
	handler := NewHandler(uc)
	return handler, mockRepo, mockAbilityRepo
}

func TestHandler_List_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	mockRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*domain.Permission{{ID: 1, AbilityID: 1, EntityID: 1, EntityType: "user"}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/permissions?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Create_Success(t *testing.T) {
	h, mockRepo, mockAbilityRepo := newHandler()
	router := setupRouter(h)

	mockAbilityRepo.On("GetByID", mock.Anything, int64(1)).Return(&abilityDomain.Ability{ID: 1, Name: "manage_users"}, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*permission.Permission")).Return(nil)

	body := `{"abilityId":1,"entityId":1,"entityType":"user","forbidden":false}`
	req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Permission{
		ID: 1, AbilityID: 1, EntityID: 1, EntityType: "user",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/permissions/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/permissions/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/permissions/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/permissions/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
