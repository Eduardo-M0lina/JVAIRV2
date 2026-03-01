package warranty_status

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	domainWS "github.com/your-org/jvairv2/pkg/domain/warranty_status"
)

// Mock implementation
type mockWarrantyStatusUseCase struct {
	statuses []*domainWS.WarrantyStatus
}

func (m *mockWarrantyStatusUseCase) Create(ctx context.Context, status *domainWS.WarrantyStatus) error {
	status.ID = int64(len(m.statuses) + 1)
	m.statuses = append(m.statuses, status)
	return nil
}

func (m *mockWarrantyStatusUseCase) GetByID(ctx context.Context, id int64) (*domainWS.WarrantyStatus, error) {
	for _, s := range m.statuses {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, domainWS.ErrWarrantyStatusNotFound
}

func (m *mockWarrantyStatusUseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWS.WarrantyStatus, int, error) {
	return m.statuses, len(m.statuses), nil
}

func (m *mockWarrantyStatusUseCase) Update(ctx context.Context, status *domainWS.WarrantyStatus) error {
	for i, s := range m.statuses {
		if s.ID == status.ID {
			m.statuses[i] = status
			return nil
		}
	}
	return domainWS.ErrWarrantyStatusNotFound
}

func (m *mockWarrantyStatusUseCase) Delete(ctx context.Context, id int64) error {
	for i, s := range m.statuses {
		if s.ID == id {
			m.statuses = append(m.statuses[:i], m.statuses[i+1:]...)
			return nil
		}
	}
	return domainWS.ErrWarrantyStatusNotFound
}

func setupHandler() *Handler {
	statusUC := &mockWarrantyStatusUseCase{}
	return NewHandler(statusUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("GET", "/warranty-statuses?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	h := setupHandler()

	// Create a status first
	ctx := context.Background()
	status := &domainWS.WarrantyStatus{
		Label:    "Active",
		Order:    1,
		IsActive: true,
	}
	_ = h.useCase.Create(ctx, status)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("GET", "/warranty-statuses/1", nil)
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
	req := httptest.NewRequest("GET", "/warranty-statuses/999", nil)
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
		"name": "Pending",
		"description": "Pending warranty status"
	}`

	req := httptest.NewRequest("POST", "/warranty-statuses", strings.NewReader(body))
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

	req := httptest.NewRequest("POST", "/warranty-statuses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	h := setupHandler()

	// Create a status first
	ctx := context.Background()
	status := &domainWS.WarrantyStatus{
		Label:    "Test Status",
		Order:    1,
		IsActive: true,
	}
	_ = h.useCase.Create(ctx, status)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("DELETE", "/warranty-statuses/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
