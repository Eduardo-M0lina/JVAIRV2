package property_equipment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List godoc
// @Summary List property equipment
// @Description Get a list of equipment for a property
// @Tags PropertyEquipment
// @Produce json
// @Param propertyId path int true "Property ID"
// @Success 200 {array} PropertyEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties/{propertyId}/equipment [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	propertyIDStr := chi.URLParam(r, "propertyId")
	propertyID, err := strconv.ParseInt(propertyIDStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid property ID",
			slog.String("propertyId", propertyIDStr))
		response.Error(w, http.StatusBadRequest, "Invalid property ID")
		return
	}

	equipment, err := h.useCase.List(r.Context(), propertyID)
	if err != nil {
		if err.Error() == "invalid property_id" || err.Error() == "property is deleted" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to list property equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to list property equipment")
		return
	}

	items := make([]PropertyEquipmentResponse, len(equipment))
	for i, e := range equipment {
		items[i] = toResponse(e)
	}

	response.JSON(w, http.StatusOK, items)
}
