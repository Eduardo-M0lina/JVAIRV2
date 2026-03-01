package quote_status

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domain "github.com/your-org/jvairv2/pkg/domain/quote_status"
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

	mockSvc.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*domain.QuoteStatus{{ID: 1, Label: "Pending"}}, 1, nil)

	req := httptest.NewRequest(http.MethodGet, "/quote-statuses?page=1&pageSize=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Create_Success(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*quote_status.QuoteStatus")).Return(nil)

	body := `{"label":"Pending","order":1}`
	req := httptest.NewRequest(http.MethodPost, "/quote-statuses", bytes.NewBufferString(body))
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

	req := httptest.NewRequest(http.MethodPost, "/quote-statuses", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(&domain.QuoteStatus{ID: 1, Label: "Pending"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/quote-statuses/1", nil)
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

	req := httptest.NewRequest(http.MethodGet, "/quote-statuses/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.On("GetByID", mock.Anything, int64(999)).Return(nil, domain.ErrQuoteStatusNotFound)

	req := httptest.NewRequest(http.MethodGet, "/quote-statuses/999", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/quote-statuses/1", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	mockSvc := new(domain.MockService)
	h := NewHandler(mockSvc)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/quote-statuses/abc", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
