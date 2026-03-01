package property_equipment

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
	domain "github.com/your-org/jvairv2/pkg/domain/property"
	pe "github.com/your-org/jvairv2/pkg/domain/property_equipment"
)

func strPtr(s string) *string { return &s }

func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func newHandler() (*Handler, *pe.MockRepository, *domain.MockRepository) {
	mockRepo := new(pe.MockRepository)
	mockPropRepo := new(domain.MockRepository)
	uc := pe.NewUseCase(mockRepo, mockPropRepo)
	handler := NewHandler(uc)
	return handler, mockRepo, mockPropRepo
}

func TestHandler_Create_Success(t *testing.T) {
	h, mockRepo, mockPropRepo := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockPropRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Property{
		ID: 1, Street: "123 Main St", CreatedAt: &now,
	}, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*property_equipment.PropertyEquipment")).Return(nil)

	body := `{"area":"Main Floor","outdoorBrand":"Carrier"}`
	req := httptest.NewRequest(http.MethodPost, "/properties/1/equipment", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
	mockPropRepo.AssertExpectations(t)
}

func TestHandler_Create_InvalidPropertyID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	body := `{"area":"Main Floor"}`
	req := httptest.NewRequest(http.MethodPost, "/properties/abc/equipment", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/properties/1/equipment", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_List_Success(t *testing.T) {
	h, mockRepo, mockPropRepo := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockPropRepo.On("GetByID", mock.Anything, int64(1)).Return(&domain.Property{
		ID: 1, Street: "123 Main St", CreatedAt: &now,
	}, nil)
	mockRepo.On("List", mock.Anything, int64(1)).Return([]*pe.PropertyEquipment{
		{ID: 1, PropertyID: 1, Area: strPtr("Main Floor"), CreatedAt: &now},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/properties/1/equipment", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var items []PropertyEquipmentResponse
	err := json.NewDecoder(rr.Body).Decode(&items)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
	mockPropRepo.AssertExpectations(t)
}

func TestHandler_List_InvalidPropertyID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/properties/abc/equipment", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockRepo.On("GetByID", mock.Anything, int64(5)).Return(&pe.PropertyEquipment{
		ID: 5, PropertyID: 1, Area: strPtr("Basement"), CreatedAt: &now,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/properties/1/equipment/5", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp PropertyEquipmentResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), resp.ID)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/properties/1/equipment/abc", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/properties/1/equipment/999", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Contains(t, []int{http.StatusNotFound, http.StatusInternalServerError}, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	h, mockRepo, mockPropRepo := newHandler()
	router := setupRouter(h)

	now := time.Now()
	existing := &pe.PropertyEquipment{
		ID: 1, PropertyID: 1, Area: strPtr("Main Floor"), CreatedAt: &now,
	}
	mockProp := &domain.Property{ID: 1, Street: "123 Main St", CreatedAt: &now}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existing, nil).Once()
	mockPropRepo.On("GetByID", mock.Anything, int64(1)).Return(mockProp, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*property_equipment.PropertyEquipment")).Return(nil)
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&pe.PropertyEquipment{
		ID: 1, PropertyID: 1, Area: strPtr("Second Floor"), CreatedAt: &now,
	}, nil).Once()

	body := `{"area":"Second Floor"}`
	req := httptest.NewRequest(http.MethodPut, "/properties/1/equipment/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
	mockPropRepo.AssertExpectations(t)
}

func TestHandler_Update_InvalidPropertyID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	body := `{"area":"Main Floor"}`
	req := httptest.NewRequest(http.MethodPut, "/properties/abc/equipment/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Update_InvalidEquipmentID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	body := `{"area":"Main Floor"}`
	req := httptest.NewRequest(http.MethodPut, "/properties/1/equipment/abc", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&pe.PropertyEquipment{
		ID: 1, PropertyID: 1, Area: strPtr("Main Floor"), CreatedAt: &now,
	}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/properties/1/equipment/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Delete_InvalidPropertyID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/properties/abc/equipment/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_InvalidEquipmentID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/properties/1/equipment/abc", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
