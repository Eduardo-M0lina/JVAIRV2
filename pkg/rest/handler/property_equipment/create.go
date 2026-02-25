package property_equipment

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domain "github.com/your-org/jvairv2/pkg/domain/property_equipment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Create godoc
// @Summary Create property equipment
// @Description Create a new equipment entry for a property
// @Tags PropertyEquipment
// @Accept json
// @Produce json
// @Param propertyId path int true "Property ID"
// @Param equipment body CreatePropertyEquipmentRequest true "Equipment data"
// @Success 201 {object} PropertyEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties/{propertyId}/equipment [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	propertyIDStr := chi.URLParam(r, "propertyId")
	propertyID, err := strconv.ParseInt(propertyIDStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid property ID",
			slog.String("propertyId", propertyIDStr))
		response.Error(w, http.StatusBadRequest, "Invalid property ID")
		return
	}

	var req CreatePropertyEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eq := &domain.PropertyEquipment{
		PropertyID:          propertyID,
		Area:                req.Area,
		OutdoorBrand:        req.OutdoorBrand,
		OutdoorModel:        req.OutdoorModel,
		OutdoorSerial:       req.OutdoorSerial,
		OutdoorInstalled:    parseTimePtr(req.OutdoorInstalled),
		FurnaceBrand:        req.FurnaceBrand,
		FurnaceModel:        req.FurnaceModel,
		FurnaceSerial:       req.FurnaceSerial,
		FurnaceInstalled:    parseTimePtr(req.FurnaceInstalled),
		EvaporatorBrand:     req.EvaporatorBrand,
		EvaporatorModel:     req.EvaporatorModel,
		EvaporatorSerial:    req.EvaporatorSerial,
		EvaporatorInstalled: parseTimePtr(req.EvaporatorInstalled),
		AirHandlerBrand:     req.AirHandlerBrand,
		AirHandlerModel:     req.AirHandlerModel,
		AirHandlerSerial:    req.AirHandlerSerial,
		AirHandlerInstalled: parseTimePtr(req.AirHandlerInstalled),
	}

	if err := eq.Validate(); err != nil {
		slog.WarnContext(r.Context(), "Validation failed",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.useCase.Create(r.Context(), eq); err != nil {
		if err.Error() == "invalid property_id" || err.Error() == "property is deleted" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to create property equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to create property equipment")
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(eq))
}
