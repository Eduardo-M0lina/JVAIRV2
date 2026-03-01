package warranty_equipment

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	domainWE "github.com/your-org/jvairv2/pkg/domain/warranty_equipment"
)

// Mock implementation
type mockWarrantyEquipmentUseCase struct {
	equipment []*domainWE.WarrantyEquipment
}

func (m *mockWarrantyEquipmentUseCase) Create(ctx context.Context, we *domainWE.WarrantyEquipment) error {
	we.ID = int64(len(m.equipment) + 1)
	m.equipment = append(m.equipment, we)
	return nil
}

func (m *mockWarrantyEquipmentUseCase) GetByID(ctx context.Context, id int64) (*domainWE.WarrantyEquipment, error) {
	for _, e := range m.equipment {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, domainWE.ErrWarrantyEquipmentNotFound
}

func (m *mockWarrantyEquipmentUseCase) ListByWarrantyID(ctx context.Context, warrantyID int64) ([]*domainWE.WarrantyEquipment, error) {
	var result []*domainWE.WarrantyEquipment
	for _, e := range m.equipment {
		if e.WarrantyID == warrantyID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockWarrantyEquipmentUseCase) Update(ctx context.Context, we *domainWE.WarrantyEquipment) error {
	for i, e := range m.equipment {
		if e.ID == we.ID {
			m.equipment[i] = we
			return nil
		}
	}
	return domainWE.ErrWarrantyEquipmentNotFound
}

func (m *mockWarrantyEquipmentUseCase) Delete(ctx context.Context, id int64) error {
	for i, e := range m.equipment {
		if e.ID == id {
			m.equipment = append(m.equipment[:i], m.equipment[i+1:]...)
			return nil
		}
	}
	return domainWE.ErrWarrantyEquipmentNotFound
}

func (m *mockWarrantyEquipmentUseCase) CloneFromJobEquipment(ctx context.Context, warrantyID int64, jobID int64) error {
	// Mock implementation - just create a dummy equipment
	we := &domainWE.WarrantyEquipment{
		WarrantyID: warrantyID,
		Area:       "Cloned Area",
	}
	we.ID = int64(len(m.equipment) + 1)
	m.equipment = append(m.equipment, we)
	return nil
}

func setupHandler() *Handler {
	equipmentUC := &mockWarrantyEquipmentUseCase{}
	return NewHandler(equipmentUC)
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()

	// Create equipment first
	ctx := context.Background()
	equipment := &domainWE.WarrantyEquipment{
		WarrantyID: 1,
		Area:       "Main Floor",
	}
	_ = h.useCase.Create(ctx, equipment)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("warrantyId", "1")
	req := httptest.NewRequest("GET", "/warranties/1/equipment", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.List(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandler_Create_Success(t *testing.T) {
	h := setupHandler()

	body := `{
		"area": "Main Floor",
		"outdoorBrand": "Carrier",
		"outdoorModel": "24ACC636A003",
		"outdoorSerial": "1234567890"
	}`

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("warrantyId", "1")
	req := httptest.NewRequest("POST", "/warranties/1/equipment", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
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

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("warrantyId", "1")
	req := httptest.NewRequest("POST", "/warranties/1/equipment", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Create(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandler_Update_Success(t *testing.T) {
	h := setupHandler()

	// Create equipment first
	ctx := context.Background()
	equipment := &domainWE.WarrantyEquipment{
		WarrantyID: 1,
		Area:       "Main Floor",
	}
	_ = h.useCase.Create(ctx, equipment)

	body := `{
		"area": "Updated Area",
		"outdoorBrand": "Updated Brand"
	}`

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("warrantyId", "1")
	rctx.URLParams.Add("equipmentId", "1")
	req := httptest.NewRequest("PUT", "/warranties/1/equipment/1", strings.NewReader(body))
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

	// Create equipment first
	ctx := context.Background()
	equipment := &domainWE.WarrantyEquipment{
		WarrantyID: 1,
		Area:       "Main Floor",
	}
	_ = h.useCase.Create(ctx, equipment)

	// Set up chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("warrantyId", "1")
	rctx.URLParams.Add("equipmentId", "1")
	req := httptest.NewRequest("DELETE", "/warranties/1/equipment/1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	h.Delete(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
