package property

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get godoc
// @Summary Get property by ID
// @Description Get a property by its ID
// @Tags Properties
// @Produce json
// @Param id path int true "Property ID"
// @Success 200 {object} PropertyResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid property ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid property ID")
		return
	}

	property, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "property not found" {
			response.Error(w, http.StatusNotFound, "Property not found")
			return
		}

		slog.ErrorContext(r.Context(), "Failed to get property",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to get property")
		return
	}

	response.JSON(w, http.StatusOK, toPropertyResponse(property))
}
