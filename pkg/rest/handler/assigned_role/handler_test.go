package assigned_role

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domain "github.com/your-org/jvairv2/pkg/domain/assigned_role"
	roleDomain "github.com/your-org/jvairv2/pkg/domain/role"
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
	r.Route("/assigned-roles", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Assign)
		r.Get("/{id}", h.Get)
		r.Get("/entity/{entityType}/{entityId}", h.GetByEntity)
		r.Get("/check/{roleId}/{entityType}/{entityId}", h.HasRole)
		r.Delete("/revoke/{roleId}/{entityType}/{entityId}", h.Revoke)
	})
	return r
}

func newHandler() (*Handler, *domain.MockRepository, *roleDomain.MockRepository) {
	mockRepo := new(domain.MockRepository)
	mockRoleRepo := new(roleDomain.MockRepository)
	uc := domain.NewUseCase(mockRepo, mockRoleRepo)
	handler := NewHandler(uc)
	return handler, mockRepo, mockRoleRepo
}

func TestHandler_List_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	mockRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*domain.AssignedRole{{ID: 1, RoleID: 1, EntityID: 1, EntityType: "user"}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/assigned-roles?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Assign_Success(t *testing.T) {
	h, mockRepo, mockRoleRepo := newHandler()
	router := setupRouter(h)

	mockRoleRepo.On("GetByID", mock.Anything, int64(1)).Return(&roleDomain.Role{ID: 1, Name: "admin"}, nil)
	mockRepo.On("Assign", mock.Anything, mock.AnythingOfType("*assigned_role.AssignedRole")).Return(nil)

	body := `{"roleId":1,"entityId":1,"entityType":"user"}`
	req := httptest.NewRequest(http.MethodPost, "/assigned-roles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Assign_InvalidBody(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/assigned-roles", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.AssignedRole{
		ID: 1, RoleID: 1, EntityID: 1, EntityType: "user",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/assigned-roles/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/assigned-roles/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
