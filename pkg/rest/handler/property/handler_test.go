package property

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	customerDomain "github.com/your-org/jvairv2/pkg/domain/customer"
	domain "github.com/your-org/jvairv2/pkg/domain/property"
)

func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/properties", func(r chi.Router) {
		h.RegisterRoutes(r)
	})
	return r
}

func newHandler() (*Handler, *domain.MockRepository, *customerDomain.MockRepository) {
	mockRepo := new(domain.MockRepository)
	mockCustRepo := new(customerDomain.MockRepository)
	uc := domain.NewUseCase(mockRepo, mockCustRepo)
	handler := NewHandler(uc)
	return handler, mockRepo, mockCustRepo
}

func TestHandler_List_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*domain.Property{{ID: 1, Street: "123 Main St", CreatedAt: &now}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/properties?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Create_Success(t *testing.T) {
	h, mockRepo, mockCustRepo := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockCustRepo.On("GetByID", mock.Anything, int64(1)).Return(&customerDomain.Customer{
		ID: 1, Name: "ACME", CreatedAt: &now,
	}, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*property.Property")).Return(nil)

	body := `{"customerId":1,"street":"123 Main St","city":"NY","state":"NY","zip":"10001"}`
	req := httptest.NewRequest(http.MethodPost, "/properties", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
	mockCustRepo.AssertExpectations(t)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/properties", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Property{
		ID: 1, Street: "123 Main St", CustomerID: 1, CreatedAt: &now,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/properties/1", nil)
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
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/properties/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Property{
		ID: 1, Street: "123 Main St", CreatedAt: &now,
	}, nil)
	mockRepo.On("HasJobs", mock.Anything, int64(1)).Return(false, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/properties/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/properties/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
