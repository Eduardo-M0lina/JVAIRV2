package property

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/your-org/jvairv2/pkg/domain/property"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Create godoc
// @Summary Create property
// @Description Create a new property
// @Tags Properties
// @Accept json
// @Produce json
// @Param property body CreatePropertyRequest true "Property data"
// @Success 201 {object} PropertyResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	prop := &property.Property{
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

	if err := h.useCase.Create(r.Context(), prop); err != nil {
		if err.Error() == "invalid customer_id" || err.Error() == "customer is deleted" {
			slog.WarnContext(r.Context(), "Invalid customer",
				slog.String("error", err.Error()))
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(r.Context(), "Failed to create property",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to create property")
		return
	}

	response.JSON(w, http.StatusCreated, toPropertyResponse(prop))
}
