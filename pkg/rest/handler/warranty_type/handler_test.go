package warranty_type

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	domainWT "github.com/your-org/jvairv2/pkg/domain/warranty_type"
)

// Mock implementation
type mockWarrantyTypeUseCase struct {
	types []*domainWT.WarrantyType
}

func (m *mockWarrantyTypeUseCase) Create(ctx context.Context, wt *domainWT.WarrantyType) error {
	wt.ID = int64(len(m.types) + 1)
	m.types = append(m.types, wt)
	return nil
}

func (m *mockWarrantyTypeUseCase) GetByID(ctx context.Context, id int64) (*domainWT.WarrantyType, error) {
	for _, t := range m.types {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, domainWT.ErrWarrantyTypeNotFound
}

func (m *mockWarrantyTypeUseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWT.WarrantyType, int, error) {
	return m.types, len(m.types), nil
}

func (m *mockWarrantyTypeUseCase) Update(ctx context.Context, wt *domainWT.WarrantyType) error {
	for i, t := range m.types {
		if t.ID == wt.ID {
			m.types[i] = wt
			return nil
		}
	}
	return domainWT.ErrWarrantyTypeNotFound
}

func (m *mockWarrantyTypeUseCase) Delete(ctx context.Context, id int64) error {
	for i, t := range m.types {
		if t.ID == id {
			m.types = append(m.types[:i], m.types[i+1:]...)
			return nil
		}
	}
	return domainWT.ErrWarrantyTypeNotFound
}

func setupHandler() *Handler {
	typeUC := &mockWarrantyTypeUseCase{}
	return NewHandler(typeUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("GET", "/warranty-types?page=1&pageSize=10", nil)
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
	wt := &domainWT.WarrantyType{
		Label:       "Parts",
		LabelPlural: "Parts",
		IsActive:    true,
	}
	_ = h.useCase.Create(ctx, wt)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("GET", "/warranty-types/1", nil)
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
	req := httptest.NewRequest("GET", "/warranty-types/999", nil)
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

	req := httptest.NewRequest("POST", "/warranty-types", strings.NewReader(body))
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

	req := httptest.NewRequest("POST", "/warranty-types", strings.NewReader(body))
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
	wt := &domainWT.WarrantyType{
		Label:       "Test Type",
		LabelPlural: "Test Types",
		IsActive:    true,
	}
	_ = h.useCase.Create(ctx, wt)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("DELETE", "/warranty-types/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
