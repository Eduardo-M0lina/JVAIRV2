package property_equipment

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete godoc
// @Summary Delete property equipment
// @Description Delete a property equipment entry (hard delete)
// @Tags PropertyEquipment
// @Param propertyId path int true "Property ID"
// @Param id path int true "Equipment ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties/{propertyId}/equipment/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.useCase.Delete(r.Context(), id, propertyID); err != nil {
		if err.Error() == "property equipment not found" {
			response.Error(w, http.StatusNotFound, "Property equipment not found")
			return
		}

		if err.Error() == "equipment does not belong to this property" {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to delete property equipment",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to delete property equipment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
