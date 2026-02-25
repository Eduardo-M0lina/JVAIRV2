package job_equipment

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
	je "github.com/your-org/jvairv2/pkg/domain/job_equipment"
)

func strPtr(s string) *string { return &s }

func setupRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func newHandler() (*Handler, *je.MockRepository, *je.MockJobChecker) {
	mockRepo := new(je.MockRepository)
	mockJobChecker := new(je.MockJobChecker)
	uc := je.NewUseCase(mockRepo, mockJobChecker)
	handler := NewHandler(uc)
	return handler, mockRepo, mockJobChecker
}

func TestHandler_Create_Success(t *testing.T) {
	h, mockRepo, mockJobChecker := newHandler()
	router := setupRouter(h)

	mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*job_equipment.JobEquipment")).Return(nil)

	body := `{"type":"current","area":"Main Floor","outdoorBrand":"Carrier"}`
	req := httptest.NewRequest(http.MethodPost, "/jobs/1/equipment", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
	mockJobChecker.AssertExpectations(t)
}

func TestHandler_Create_InvalidJobID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	body := `{"type":"current","area":"Main Floor"}`
	req := httptest.NewRequest(http.MethodPost, "/jobs/abc/equipment", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/jobs/1/equipment", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_List_Success(t *testing.T) {
	h, mockRepo, mockJobChecker := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
	mockRepo.On("List", mock.Anything, int64(1), "").Return([]*je.JobEquipment{
		{ID: 1, JobID: 1, Type: "current", Area: strPtr("Main Floor"), CreatedAt: &now},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/1/equipment", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var items []JobEquipmentResponse
	err := json.NewDecoder(rr.Body).Decode(&items)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	mockRepo.AssertExpectations(t)
}

func TestHandler_List_WithTypeFilter(t *testing.T) {
	h, mockRepo, mockJobChecker := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
	mockRepo.On("List", mock.Anything, int64(1), "current").Return([]*je.JobEquipment{
		{ID: 1, JobID: 1, Type: "current", Area: strPtr("Main Floor"), CreatedAt: &now},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/1/equipment?type=current", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_List_InvalidJobID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/jobs/abc/equipment", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockRepo.On("GetByID", mock.Anything, int64(5)).Return(&je.JobEquipment{
		ID: 5, JobID: 1, Type: "current", Area: strPtr("Basement"), CreatedAt: &now,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/1/equipment/5", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp JobEquipmentResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), resp.ID)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/jobs/1/equipment/abc", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Get_NotFound(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/jobs/1/equipment/999", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Contains(t, []int{http.StatusNotFound, http.StatusInternalServerError}, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	h, mockRepo, _ := newHandler()
	router := setupRouter(h)

	now := time.Now()
	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(&je.JobEquipment{
		ID: 1, JobID: 1, Type: "current", Area: strPtr("Main Floor"), CreatedAt: &now,
	}, nil)
	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/jobs/1/equipment/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandler_Delete_InvalidJobID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/jobs/abc/equipment/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_Delete_InvalidEquipmentID(t *testing.T) {
	h, _, _ := newHandler()
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodDelete, "/jobs/1/equipment/abc", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
