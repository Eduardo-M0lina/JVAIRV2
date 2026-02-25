package customer

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
	domainCustomer "github.com/your-org/jvairv2/pkg/domain/customer"
)

func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func TestHandler_List_Success(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	now := time.Now()
	mockSvc.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*domainCustomer.Customer{{ID: 1, Name: "ACME Corp", CreatedAt: &now}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/customers?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create_Success(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*customer.Customer")).Return(nil)

	body := `{"name":"ACME Corp","workflowId":1}`
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	now := time.Now()
	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(&domainCustomer.Customer{
		ID: 1, Name: "ACME Corp", CreatedAt: &now,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/customers/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), resp["id"])
	mockSvc.AssertExpectations(t)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/customers/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	mockSvc.On("GetByID", mock.Anything, int64(999)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/customers/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Contains(t, []int{http.StatusNotFound, http.StatusInternalServerError}, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	now := time.Now()
	mockSvc.On("Update", mock.Anything, mock.AnythingOfType("*customer.Customer")).Return(nil)
	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(&domainCustomer.Customer{
		ID: 1, Name: "ACME Updated", CreatedAt: &now,
	}, nil)

	body := `{"name":"ACME Updated","workflowId":1}`
	req := httptest.NewRequest(http.MethodPut, "/customers/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Update_InvalidID(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	body := `{"name":"ACME"}`
	req := httptest.NewRequest(http.MethodPut, "/customers/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/customers/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	mockSvc := new(domainCustomer.MockService)
	h := NewHandler(mockSvc, nil)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/customers/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
