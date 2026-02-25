package supervisor

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/your-org/jvairv2/pkg/domain/supervisor"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Create godoc
// @Summary Create a new supervisor
// @Description Create a new supervisor associated with a customer
// @Tags supervisors
// @Accept json
// @Produce json
// @Param supervisor body CreateSupervisorRequest true "Supervisor information"
// @Success 201 {object} SupervisorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /supervisors [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSupervisorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(r.Context(), "Failed to decode request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	s := &supervisor.Supervisor{
		CustomerID: req.CustomerID,
		Name:       req.Name,
		Phone:      req.Phone,
		Email:      req.Email,
	}

	// Validar campos requeridos usando el método de la entidad
	if err := s.Validate(); err != nil {
		slog.WarnContext(r.Context(), "Validation failed",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.useCase.Create(r.Context(), s); err != nil {
		slog.ErrorContext(r.Context(), "Failed to create supervisor",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toSupervisorResponse(s))
}
