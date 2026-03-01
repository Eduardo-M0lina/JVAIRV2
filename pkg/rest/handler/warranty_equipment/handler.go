package warranty_equipment

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	domainWE "github.com/your-org/jvairv2/pkg/domain/warranty_equipment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las peticiones HTTP para equipos de garantía
type Handler struct {
	useCase domainWE.Service
}

// NewHandler crea una nueva instancia del handler
func NewHandler(useCase domainWE.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler como sub-recurso de warranties
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/warranties/{warrantyId}/equipment", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Put("/{equipmentId}", h.Update)
		r.Delete("/{equipmentId}", h.Delete)
	})
}

const timeFormat = "2006-01-02T15:04:05Z07:00"
const dateFormat = "2006-01-02"

// WarrantyEquipmentRequest representa la solicitud para crear/actualizar un equipo de garantía
type WarrantyEquipmentRequest struct {
	Area                string  `json:"area" example:"Main Floor"`
	OutdoorBrand        *string `json:"outdoorBrand,omitempty" example:"Carrier"`
	OutdoorModel        *string `json:"outdoorModel,omitempty" example:"24ACC636A003"`
	OutdoorSerial       *string `json:"outdoorSerial,omitempty" example:"1234567890"`
	OutdoorInstalled    *string `json:"outdoorInstalled,omitempty" example:"2024-01-15"`
	FurnaceBrand        *string `json:"furnaceBrand,omitempty" example:"Lennox"`
	FurnaceModel        *string `json:"furnaceModel,omitempty" example:"ML180UH"`
	FurnaceSerial       *string `json:"furnaceSerial,omitempty" example:"0987654321"`
	FurnaceInstalled    *string `json:"furnaceInstalled,omitempty" example:"2024-01-15"`
	EvaporatorBrand     *string `json:"evaporatorBrand,omitempty" example:"Goodman"`
	EvaporatorModel     *string `json:"evaporatorModel,omitempty" example:"CAPF4961D6"`
	EvaporatorSerial    *string `json:"evaporatorSerial,omitempty" example:"1122334455"`
	EvaporatorInstalled *string `json:"evaporatorInstalled,omitempty" example:"2024-01-15"`
	AirHandlerBrand     *string `json:"airHandlerBrand,omitempty" example:"Trane"`
	AirHandlerModel     *string `json:"airHandlerModel,omitempty" example:"GAM5A0C48M41SB"`
	AirHandlerSerial    *string `json:"airHandlerSerial,omitempty" example:"5566778899"`
	AirHandlerInstalled *string `json:"airHandlerInstalled,omitempty" example:"2024-01-15"`
}

// WarrantyEquipmentResponse representa la respuesta de un equipo de garantía
type WarrantyEquipmentResponse struct {
	ID                  int64   `json:"id" example:"1"`
	WarrantyID          int64   `json:"warrantyId" example:"1"`
	Area                string  `json:"area" example:"Main Floor"`
	OutdoorBrand        *string `json:"outdoorBrand,omitempty"`
	OutdoorModel        *string `json:"outdoorModel,omitempty"`
	OutdoorSerial       *string `json:"outdoorSerial,omitempty"`
	OutdoorInstalled    *string `json:"outdoorInstalled,omitempty"`
	FurnaceBrand        *string `json:"furnaceBrand,omitempty"`
	FurnaceModel        *string `json:"furnaceModel,omitempty"`
	FurnaceSerial       *string `json:"furnaceSerial,omitempty"`
	FurnaceInstalled    *string `json:"furnaceInstalled,omitempty"`
	EvaporatorBrand     *string `json:"evaporatorBrand,omitempty"`
	EvaporatorModel     *string `json:"evaporatorModel,omitempty"`
	EvaporatorSerial    *string `json:"evaporatorSerial,omitempty"`
	EvaporatorInstalled *string `json:"evaporatorInstalled,omitempty"`
	AirHandlerBrand     *string `json:"airHandlerBrand,omitempty"`
	AirHandlerModel     *string `json:"airHandlerModel,omitempty"`
	AirHandlerSerial    *string `json:"airHandlerSerial,omitempty"`
	AirHandlerInstalled *string `json:"airHandlerInstalled,omitempty"`
	CreatedAt           string  `json:"createdAt,omitempty"`
	UpdatedAt           string  `json:"updatedAt,omitempty"`
}

func toResponse(we *domainWE.WarrantyEquipment) WarrantyEquipmentResponse {
	resp := WarrantyEquipmentResponse{
		ID:               we.ID,
		WarrantyID:       we.WarrantyID,
		Area:             we.Area,
		OutdoorBrand:     we.OutdoorBrand,
		OutdoorModel:     we.OutdoorModel,
		OutdoorSerial:    we.OutdoorSerial,
		FurnaceBrand:     we.FurnaceBrand,
		FurnaceModel:     we.FurnaceModel,
		FurnaceSerial:    we.FurnaceSerial,
		EvaporatorBrand:  we.EvaporatorBrand,
		EvaporatorModel:  we.EvaporatorModel,
		EvaporatorSerial: we.EvaporatorSerial,
		AirHandlerBrand:  we.AirHandlerBrand,
		AirHandlerModel:  we.AirHandlerModel,
		AirHandlerSerial: we.AirHandlerSerial,
	}

	if we.OutdoorInstalled != nil {
		s := we.OutdoorInstalled.Format(dateFormat)
		resp.OutdoorInstalled = &s
	}
	if we.FurnaceInstalled != nil {
		s := we.FurnaceInstalled.Format(dateFormat)
		resp.FurnaceInstalled = &s
	}
	if we.EvaporatorInstalled != nil {
		s := we.EvaporatorInstalled.Format(dateFormat)
		resp.EvaporatorInstalled = &s
	}
	if we.AirHandlerInstalled != nil {
		s := we.AirHandlerInstalled.Format(dateFormat)
		resp.AirHandlerInstalled = &s
	}
	if we.CreatedAt != nil {
		resp.CreatedAt = we.CreatedAt.Format(timeFormat)
	}
	if we.UpdatedAt != nil {
		resp.UpdatedAt = we.UpdatedAt.Format(timeFormat)
	}

	return resp
}

func parseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse(dateFormat, *s)
	if err != nil {
		return nil
	}
	return &t
}

// @Summary Listar equipos de garantía
// @Description Obtiene la lista de equipos de una garantía
// @Tags WarrantyEquipment
// @Accept json
// @Produce json
// @Param warrantyId path int true "ID de la garantía"
// @Success 200 {array} WarrantyEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties/{warrantyId}/equipment [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	warrantyID, err := strconv.ParseInt(chi.URLParam(r, "warrantyId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de garantía inválido")
		return
	}

	equipment, err := h.useCase.ListByWarrantyID(r.Context(), warrantyID)
	if err != nil {
		if err == domainWE.ErrInvalidWarranty {
			response.Error(w, http.StatusNotFound, "Garantía no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al listar equipos de garantía")
		return
	}

	items := make([]WarrantyEquipmentResponse, len(equipment))
	for i, we := range equipment {
		items[i] = toResponse(we)
	}

	response.JSON(w, http.StatusOK, items)
}

// @Summary Crear equipo de garantía
// @Description Crea un nuevo equipo de garantía
// @Tags WarrantyEquipment
// @Accept json
// @Produce json
// @Param warrantyId path int true "ID de la garantía"
// @Param equipment body WarrantyEquipmentRequest true "Datos del equipo"
// @Success 201 {object} WarrantyEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties/{warrantyId}/equipment [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	warrantyID, err := strconv.ParseInt(chi.URLParam(r, "warrantyId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de garantía inválido")
		return
	}

	var req WarrantyEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	equipment := &domainWE.WarrantyEquipment{
		WarrantyID:          warrantyID,
		Area:                req.Area,
		OutdoorBrand:        req.OutdoorBrand,
		OutdoorModel:        req.OutdoorModel,
		OutdoorSerial:       req.OutdoorSerial,
		OutdoorInstalled:    parseDate(req.OutdoorInstalled),
		FurnaceBrand:        req.FurnaceBrand,
		FurnaceModel:        req.FurnaceModel,
		FurnaceSerial:       req.FurnaceSerial,
		FurnaceInstalled:    parseDate(req.FurnaceInstalled),
		EvaporatorBrand:     req.EvaporatorBrand,
		EvaporatorModel:     req.EvaporatorModel,
		EvaporatorSerial:    req.EvaporatorSerial,
		EvaporatorInstalled: parseDate(req.EvaporatorInstalled),
		AirHandlerBrand:     req.AirHandlerBrand,
		AirHandlerModel:     req.AirHandlerModel,
		AirHandlerSerial:    req.AirHandlerSerial,
		AirHandlerInstalled: parseDate(req.AirHandlerInstalled),
	}

	if err := h.useCase.Create(r.Context(), equipment); err != nil {
		if err == domainWE.ErrInvalidWarranty {
			response.Error(w, http.StatusNotFound, "Garantía no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(equipment))
}

// @Summary Actualizar equipo de garantía
// @Description Actualiza un equipo de garantía existente
// @Tags WarrantyEquipment
// @Accept json
// @Produce json
// @Param warrantyId path int true "ID de la garantía"
// @Param equipmentId path int true "ID del equipo"
// @Param equipment body WarrantyEquipmentRequest true "Datos del equipo"
// @Success 200 {object} WarrantyEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties/{warrantyId}/equipment/{equipmentId} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	warrantyID, err := strconv.ParseInt(chi.URLParam(r, "warrantyId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de garantía inválido")
		return
	}

	equipmentID, err := strconv.ParseInt(chi.URLParam(r, "equipmentId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	var req WarrantyEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	equipment := &domainWE.WarrantyEquipment{
		ID:                  equipmentID,
		WarrantyID:          warrantyID,
		Area:                req.Area,
		OutdoorBrand:        req.OutdoorBrand,
		OutdoorModel:        req.OutdoorModel,
		OutdoorSerial:       req.OutdoorSerial,
		OutdoorInstalled:    parseDate(req.OutdoorInstalled),
		FurnaceBrand:        req.FurnaceBrand,
		FurnaceModel:        req.FurnaceModel,
		FurnaceSerial:       req.FurnaceSerial,
		FurnaceInstalled:    parseDate(req.FurnaceInstalled),
		EvaporatorBrand:     req.EvaporatorBrand,
		EvaporatorModel:     req.EvaporatorModel,
		EvaporatorSerial:    req.EvaporatorSerial,
		EvaporatorInstalled: parseDate(req.EvaporatorInstalled),
		AirHandlerBrand:     req.AirHandlerBrand,
		AirHandlerModel:     req.AirHandlerModel,
		AirHandlerSerial:    req.AirHandlerSerial,
		AirHandlerInstalled: parseDate(req.AirHandlerInstalled),
	}

	if err := h.useCase.Update(r.Context(), equipment); err != nil {
		if err == domainWE.ErrWarrantyEquipmentNotFound {
			response.Error(w, http.StatusNotFound, "Equipo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(equipment))
}

// @Summary Eliminar equipo de garantía
// @Description Elimina un equipo de garantía (hard delete)
// @Tags WarrantyEquipment
// @Accept json
// @Produce json
// @Param warrantyId path int true "ID de la garantía"
// @Param equipmentId path int true "ID del equipo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties/{warrantyId}/equipment/{equipmentId} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	equipmentID, err := strconv.ParseInt(chi.URLParam(r, "equipmentId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de equipo inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), equipmentID); err != nil {
		if err == domainWE.ErrWarrantyEquipmentNotFound {
			response.Error(w, http.StatusNotFound, "Equipo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar equipo de garantía")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
