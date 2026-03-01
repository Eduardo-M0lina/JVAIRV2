package invoice

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
	domain "github.com/your-org/jvairv2/pkg/domain/invoice"
)

func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func TestHandler_List_Success(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	now := time.Now()
	mockSvc.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*domain.Invoice{{ID: 1, JobID: 1, InvoiceNumber: "INV-001", Total: 500.0, CreatedAt: &now}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/invoices?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create_Success(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*invoice.Invoice")).Return(nil)

	body := `{"jobId":1,"invoiceNumber":"INV-001","total":500.0}`
	req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/invoices", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	now := time.Now()
	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(&domain.Invoice{
		ID: 1, JobID: 1, InvoiceNumber: "INV-001", Total: 500.0, CreatedAt: &now,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/invoices/1", nil)
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
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/invoices/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/invoices/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/invoices/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
