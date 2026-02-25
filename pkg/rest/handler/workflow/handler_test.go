package workflow

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domain "github.com/your-org/jvairv2/pkg/domain/workflow"
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
	r.Route("/workflows", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Post("/{id}/duplicate", h.Duplicate)
	})
	return r
}

func newHandler() (*Handler, *domain.MockRepository) {
	mockRepo := new(domain.MockRepository)
	uc := domain.NewUseCase(mockRepo)
	handler := NewHandler(uc)
	return handler, mockRepo
}

func TestHandler_List_Success(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]domain.Workflow{{ID: 1, Name: "Default"}}, int64(1), nil)

	req := httptest.NewRequest(http.MethodGet, "/workflows?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Create_Success(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*workflow.Workflow")).Return(nil)
	mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(&domain.Workflow{ID: 1, Name: "Default"}, nil)
	mockRepo.On("GetWorkflowStatuses", mock.Anything, mock.Anything).Return([]domain.WorkflowStatus{}, nil)

	body := `{"name":"Default","isActive":true}`
	req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Workflow{ID: 1, Name: "Default"}, nil)
	mockRepo.On("GetWorkflowStatuses", mock.Anything, int64(1)).Return([]domain.WorkflowStatus{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/workflows/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), resp["id"])
	mockRepo.AssertExpectations(t)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	h, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/workflows/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	h, mockRepo := newHandler()
	router := setupRouter(h)

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Workflow{ID: 1, Name: "Default"}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/workflows/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	h, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/workflows/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
