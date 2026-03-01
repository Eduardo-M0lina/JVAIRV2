package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	assignedRoleDomain "github.com/your-org/jvairv2/pkg/domain/assigned_role"
	roleDomain "github.com/your-org/jvairv2/pkg/domain/role"
	domain "github.com/your-org/jvairv2/pkg/domain/user"
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
	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
	return r
}

func newHandler() (*Handler, *domain.MockRepository) {
	mockRepo := new(domain.MockRepository)
	mockAssignedRoleRepo := new(assignedRoleDomain.MockRepository)
	mockRoleRepo := new(roleDomain.MockRepository)
	uc := domain.NewUseCase(mockRepo, mockAssignedRoleRepo, mockRoleRepo)
	handler := NewHandler(uc)
	return handler, mockRepo
}

func TestHandler_List_Success(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*domain.User{{ID: 1, Name: "Admin", Email: "admin@test.com"}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("GetByID", mock.Anything, "999").Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Contains(t, []int{http.StatusNotFound, http.StatusInternalServerError}, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockRepo.AssertExpectations(t)
}
