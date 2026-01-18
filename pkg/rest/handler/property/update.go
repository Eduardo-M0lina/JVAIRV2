package property

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/property"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Update godoc
// @Summary Update property
// @Description Update an existing property
// @Tags properties
// @Accept json
// @Produce json
// @Param id path int true "Property ID"
// @Param property body UpdatePropertyRequest true "Property data"
// @Success 200 {object} PropertyResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /properties/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid property ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid property ID")
		return
	}

	var req UpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	prop := &property.Property{
		ID:           id,
		CustomerID:   req.CustomerID,
		PropertyCode: req.PropertyCode,
		Street:       req.Street,
		City:         req.City,
		State:        req.State,
		Zip:          req.Zip,
		Notes:        req.Notes,
	}

	// Validar campos requeridos usando el método de la entidad
	if err := prop.Validate(); err != nil {
		slog.WarnContext(r.Context(), "Validation failed",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.useCase.Update(r.Context(), prop); err != nil {
		if err.Error() == "property not found" {
			response.Error(w, http.StatusNotFound, "Property not found")
			return
		}

		if err.Error() == "cannot update deleted property" ||
			err.Error() == "invalid customer_id" ||
			err.Error() == "customer is deleted" {
			slog.WarnContext(r.Context(), "Invalid update request",
				slog.String("error", err.Error()))
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to update property",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to update property")
		return
	}

	updatedProperty, _ := h.useCase.GetByID(r.Context(), id)
	response.JSON(w, http.StatusOK, toPropertyResponse(updatedProperty))
}
