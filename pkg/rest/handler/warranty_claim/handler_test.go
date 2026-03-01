package warranty_claim

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	domainWC "github.com/your-org/jvairv2/pkg/domain/warranty_claim"
)

// Mock implementations
type mockWarrantyClaimUseCase struct {
	claims []*domainWC.WarrantyClaim
}

func (m *mockWarrantyClaimUseCase) Create(ctx context.Context, claim *domainWC.WarrantyClaim) error {
	claim.ID = int64(len(m.claims) + 1)
	m.claims = append(m.claims, claim)
	return nil
}

func (m *mockWarrantyClaimUseCase) GetByID(ctx context.Context, id int64) (*domainWC.WarrantyClaim, error) {
	for _, c := range m.claims {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, domainWC.ErrWarrantyClaimNotFound
}

func (m *mockWarrantyClaimUseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domainWC.WarrantyClaim, int, error) {
	return m.claims, len(m.claims), nil
}

func (m *mockWarrantyClaimUseCase) Update(ctx context.Context, claim *domainWC.WarrantyClaim) error {
	for i, c := range m.claims {
		if c.ID == claim.ID {
			m.claims[i] = claim
			return nil
		}
	}
	return domainWC.ErrWarrantyClaimNotFound
}

func (m *mockWarrantyClaimUseCase) Delete(ctx context.Context, id int64) error {
	for i, c := range m.claims {
		if c.ID == id {
			m.claims = append(m.claims[:i], m.claims[i+1:]...)
			return nil
		}
	}
	return domainWC.ErrWarrantyClaimNotFound
}

func setupHandler() *Handler {
	claimUC := &mockWarrantyClaimUseCase{}
	return NewHandler(claimUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest("GET", "/warranty-claims?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	h := setupHandler()

	// Create a claim first
	ctx := context.Background()
	claim := &domainWC.WarrantyClaim{
		JobID:                 1,
		WarrantyClaimTypeID:   1,
		WarrantyClaimStatusID: 1,
	}
	_ = h.useCase.Create(ctx, claim)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("GET", "/warranty-claims/1", nil)
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
	req := httptest.NewRequest("GET", "/warranty-claims/999", nil)
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
		"jobId": 1,
		"warrantyClaimTypeId": 1,
		"warrantyClaimStatusId": 1,
		"notes": "Test claim notes"
	}`

	req := httptest.NewRequest("POST", "/warranty-claims", strings.NewReader(body))
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

	req := httptest.NewRequest("POST", "/warranty-claims", strings.NewReader(body))
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

	// Create a claim first
	ctx := context.Background()
	claim := &domainWC.WarrantyClaim{
		JobID:                 1,
		WarrantyClaimTypeID:   1,
		WarrantyClaimStatusID: 1,
	}
	_ = h.useCase.Create(ctx, claim)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req := httptest.NewRequest("DELETE", "/warranty-claims/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
