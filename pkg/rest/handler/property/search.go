package property

import (
	"log/slog"
	"net/http"

	"github.com/angumol/jvairv2/pkg/rest/middleware"
	"github.com/angumol/jvairv2/pkg/rest/response"
)

// SearchByAddress godoc
// @Summary Search properties by address
// @Description Search properties by street address (minimum 3 characters)
// @Tags Properties
// @Produce json
// @Param address query string true "Address search term (minimum 3 characters)"
// @Success 200 {array} PropertyResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/properties/search [get]
// @Security BearerAuth
func (h *Handler) SearchByAddress(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasAbility(r.Context(), "property_view") {
		response.Error(w, http.StatusForbidden, "No tiene permisos para buscar properties")
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		response.Error(w, http.StatusBadRequest, "address parameter is required")
		return
	}

	if len(address) < 3 {
		response.Error(w, http.StatusBadRequest, "address must be at least 3 characters long")
		return
	}

	properties, err := h.useCase.SearchByAddress(r.Context(), address)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to search properties by address",
			slog.String("error", err.Error()),
			slog.String("address", address))
		response.Error(w, http.StatusInternalServerError, "Failed to search properties")
		return
	}

	items := make([]PropertyResponse, len(properties))
	for i, p := range properties {
		items[i] = toPropertyResponse(p)
	}

	response.JSON(w, http.StatusOK, items)
}
