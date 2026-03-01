package warranty

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	domainWarranty "github.com/your-org/jvairv2/pkg/domain/warranty"
)

// Mock implementations
type mockWarrantyUseCase struct {
	warranties []*domainWarranty.Warranty
}

func (m *mockWarrantyUseCase) Create(ctx context.Context, warranty *domainWarranty.Warranty) error {
	warranty.ID = int64(len(m.warranties) + 1)
	m.warranties = append(m.warranties, warranty)
	return nil
}

func (m *mockWarrantyUseCase) GetByID(ctx context.Context, id int64) (*domainWarranty.Warranty, error) {
	for _, w := range m.warranties {
		if w.ID == id {
			return w, nil
		}
	}
	return nil, domainWarranty.ErrWarrantyNotFound
}

func (m *mockWarrantyUseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWarranty.Warranty, int, error) {
	return m.warranties, len(m.warranties), nil
}

func (m *mockWarrantyUseCase) Update(ctx context.Context, warranty *domainWarranty.Warranty) error {
	for i, w := range m.warranties {
		if w.ID == warranty.ID {
			m.warranties[i] = warranty
			return nil
		}
	}
	return domainWarranty.ErrWarrantyNotFound
}

func (m *mockWarrantyUseCase) Delete(ctx context.Context, id int64) error {
	for i, w := range m.warranties {
		if w.ID == id {
			m.warranties = append(m.warranties[:i], m.warranties[i+1:]...)
			return nil
		}
	}
	return domainWarranty.ErrWarrantyNotFound
}

func setupHandler() *Handler {
	warrantyUC := &mockWarrantyUseCase{}
	return NewHandler(warrantyUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("GET", "/warranties?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	h := setupHandler()

	// Create a warranty first
	ctx := context.Background()
	warranty := &domainWarranty.Warranty{
		WarrantyNumber:   "WARR-001",
		JobID:            1,
		WarrantyTypeID:   1,
		WarrantyStatusID: 1,
	}
	_ = h.useCase.Create(ctx, warranty)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("GET", "/warranties/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Get(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_NotFound(t *testing.T) {
	h := setupHandler()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req := httptest.NewRequest("GET", "/warranties/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Get(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestHandler_Create_Success(t *testing.T) {
	h := setupHandler()

	body := `{
		"warrantyNumber": "WARR-002",
		"jobId": 1,
		"warrantyTypeId": 1,
		"warrantyStatusId": 1,
		"notes": "Test warranty notes"
	}`

	req := httptest.NewRequest("POST", "/warranties", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h := setupHandler()

	body := `invalid json`

	req := httptest.NewRequest("POST", "/warranties", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandler_Update_Success(t *testing.T) {
	h := setupHandler()

	// Create a warranty first
	ctx := context.Background()
	warranty := &domainWarranty.Warranty{
		WarrantyNumber:   "WARR-001",
		JobID:            1,
		WarrantyTypeID:   1,
		WarrantyStatusID: 1,
	}
	_ = h.useCase.Create(ctx, warranty)

	body := `{
		"warrantyNumber": "WARR-003",
		"warrantyTypeId": 2,
		"warrantyStatusId": 2,
		"notes": "Updated warranty notes"
	}`

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("PUT", "/warranties/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Update(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	h := setupHandler()

	// Create a warranty first
	ctx := context.Background()
	warranty := &domainWarranty.Warranty{
		WarrantyNumber:   "WARR-001",
		JobID:            1,
		WarrantyTypeID:   1,
		WarrantyStatusID: 1,
	}
	_ = h.useCase.Create(ctx, warranty)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("DELETE", "/warranties/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_NotFound(t *testing.T) {
	h := setupHandler()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "999")
	req := httptest.NewRequest("DELETE", "/warranties/999", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}
