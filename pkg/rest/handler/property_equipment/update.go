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

// Update godoc
// @Summary Update property equipment
// @Description Update an existing property equipment entry
// @Tags property-equipment
// @Accept json
// @Produce json
// @Param propertyId path int true "Property ID"
// @Param id path int true "Equipment ID"
// @Param equipment body UpdatePropertyEquipmentRequest true "Equipment data"
// @Success 200 {object} PropertyEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /properties/{propertyId}/equipment/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	propertyIDStr := chi.URLParam(r, "propertyId")
	propertyID, err := strconv.ParseInt(propertyIDStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid property ID",
			slog.String("propertyId", propertyIDStr))
		response.Error(w, http.StatusBadRequest, "Invalid property ID")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid equipment ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	var req UpdatePropertyEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	eq := &domain.PropertyEquipment{
		ID:                  id,
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

	if err := h.useCase.Update(r.Context(), eq); err != nil {
		if err.Error() == "property equipment not found" {
			response.Error(w, http.StatusNotFound, "Property equipment not found")
			return
		}

		if err.Error() == "invalid property_id" ||
			err.Error() == "property is deleted" ||
			err.Error() == "equipment does not belong to this property" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to update property equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to update property equipment")
		return
	}

	updatedEq, _ := h.useCase.GetByID(r.Context(), id)
	response.JSON(w, http.StatusOK, toResponse(updatedEq))
}
