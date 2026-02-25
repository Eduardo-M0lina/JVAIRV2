package property_equipment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get godoc
// @Summary Get property equipment by ID
// @Description Get a property equipment entry by its ID
// @Tags property-equipment
// @Produce json
// @Param propertyId path int true "Property ID"
// @Param id path int true "Equipment ID"
// @Success 200 {object} PropertyEquipmentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /properties/{propertyId}/equipment/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid equipment ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	equipment, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "property equipment not found" {
			response.Error(w, http.StatusNotFound, "Property equipment not found")
			return
		}

		slog.ErrorContext(r.Context(), "Failed to get property equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to get property equipment")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(equipment))
}
