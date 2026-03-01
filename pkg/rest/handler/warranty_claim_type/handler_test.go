package warranty_claim_type

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	domainWCT "github.com/your-org/jvairv2/pkg/domain/warranty_claim_type"
)

// Mock implementation
type mockWarrantyClaimTypeUseCase struct {
	types []*domainWCT.WarrantyClaimType
}

func (m *mockWarrantyClaimTypeUseCase) Create(ctx context.Context, wct *domainWCT.WarrantyClaimType) error {
	wct.ID = int64(len(m.types) + 1)
	m.types = append(m.types, wct)
	return nil
}

func (m *mockWarrantyClaimTypeUseCase) GetByID(ctx context.Context, id int64) (*domainWCT.WarrantyClaimType, error) {
	for _, t := range m.types {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, domainWCT.ErrWarrantyClaimTypeNotFound
}

func (m *mockWarrantyClaimTypeUseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWCT.WarrantyClaimType, int, error) {
	return m.types, len(m.types), nil
}

func (m *mockWarrantyClaimTypeUseCase) Update(ctx context.Context, wct *domainWCT.WarrantyClaimType) error {
	for i, t := range m.types {
		if t.ID == wct.ID {
			m.types[i] = wct
			return nil
		}
	}
	return domainWCT.ErrWarrantyClaimTypeNotFound
}

func (m *mockWarrantyClaimTypeUseCase) Delete(ctx context.Context, id int64) error {
	for i, t := range m.types {
		if t.ID == id {
			m.types = append(m.types[:i], m.types[i+1:]...)
			return nil
		}
	}
	return domainWCT.ErrWarrantyClaimTypeNotFound
}

func setupHandler() *Handler {
	typeUC := &mockWarrantyClaimTypeUseCase{}
	return NewHandler(typeUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("GET", "/warranty-claim-types?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	h := setupHandler()

	// Create a type first
	ctx := context.Background()
	wct := &domainWCT.WarrantyClaimType{
		Label:       "Parts",
		LabelPlural: "Parts",
	}
	_ = h.useCase.Create(ctx, wct)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("GET", "/warranty-claim-types/1", nil)
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
	req := httptest.NewRequest("GET", "/warranty-claim-types/999", nil)
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
		"label": "Labor",
		"labelPlural": "Labor",
		"isActive": true
	}`

	req := httptest.NewRequest("POST", "/warranty-claim-types", strings.NewReader(body))
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

	req := httptest.NewRequest("POST", "/warranty-claim-types", strings.NewReader(body))
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

	// Create a type first
	ctx := context.Background()
	wct := &domainWCT.WarrantyClaimType{
		Label:       "Test Type",
		LabelPlural: "Test Types",
	}
	_ = h.useCase.Create(ctx, wct)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("DELETE", "/warranty-claim-types/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
