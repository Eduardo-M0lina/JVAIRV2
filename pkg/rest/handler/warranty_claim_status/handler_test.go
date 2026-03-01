package warranty_claim_status

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	domainWCS "github.com/your-org/jvairv2/pkg/domain/warranty_claim_status"
)

// Mock implementation
type mockWarrantyClaimStatusUseCase struct {
	statuses []*domainWCS.WarrantyClaimStatus
}

func (m *mockWarrantyClaimStatusUseCase) Create(ctx context.Context, wcs *domainWCS.WarrantyClaimStatus) error {
	wcs.ID = int64(len(m.statuses) + 1)
	m.statuses = append(m.statuses, wcs)
	return nil
}

func (m *mockWarrantyClaimStatusUseCase) GetByID(ctx context.Context, id int64) (*domainWCS.WarrantyClaimStatus, error) {
	for _, s := range m.statuses {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, domainWCS.ErrWarrantyClaimStatusNotFound
}

func (m *mockWarrantyClaimStatusUseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWCS.WarrantyClaimStatus, int, error) {
	return m.statuses, len(m.statuses), nil
}

func (m *mockWarrantyClaimStatusUseCase) Update(ctx context.Context, wcs *domainWCS.WarrantyClaimStatus) error {
	for i, s := range m.statuses {
		if s.ID == wcs.ID {
			m.statuses[i] = wcs
			return nil
		}
	}
	return domainWCS.ErrWarrantyClaimStatusNotFound
}

func (m *mockWarrantyClaimStatusUseCase) Delete(ctx context.Context, id int64) error {
	for i, s := range m.statuses {
		if s.ID == id {
			m.statuses = append(m.statuses[:i], m.statuses[i+1:]...)
			return nil
		}
	}
	return domainWCS.ErrWarrantyClaimStatusNotFound
}

func setupHandler() *Handler {
	statusUC := &mockWarrantyClaimStatusUseCase{}
	return NewHandler(statusUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("GET", "/warranty-claim-statuses?page=1&pageSize=10", nil)
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
	status := &domainWCS.WarrantyClaimStatus{
		Label:    "Pending",
		Order:    1,
		IsActive: true,
	}
	_ = h.useCase.Create(ctx, status)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("GET", "/warranty-claim-statuses/1", nil)
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
	req := httptest.NewRequest("GET", "/warranty-claim-statuses/999", nil)
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
		"label": "Approved",
		"order": 2,
		"isActive": true
	}`

	req := httptest.NewRequest("POST", "/warranty-claim-statuses", strings.NewReader(body))
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

	req := httptest.NewRequest("POST", "/warranty-claim-statuses", strings.NewReader(body))
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
	status := &domainWCS.WarrantyClaimStatus{
		Label:    "Test Status",
		Order:    1,
		IsActive: true,
	}
	_ = h.useCase.Create(ctx, status)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("DELETE", "/warranty-claim-statuses/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
