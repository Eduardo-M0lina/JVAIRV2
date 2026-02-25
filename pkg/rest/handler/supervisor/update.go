package supervisor

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/supervisor"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Update godoc
// @Summary Update a supervisor
// @Description Update an existing supervisor with the provided information
// @Tags Supervisors
// @Accept json
// @Produce json
// @Param id path int true "Supervisor ID"
// @Param supervisor body UpdateSupervisorRequest true "Supervisor information"
// @Success 200 {object} SupervisorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/supervisors/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.WarnContext(r.Context(), "Invalid supervisor ID",
			slog.String("id", idStr))
		response.Error(w, http.StatusBadRequest, "Invalid supervisor ID")
		return
	}

	var req UpdateSupervisorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(r.Context(), "Failed to decode request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	s := &supervisor.Supervisor{
		ID:         id,
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

	if err := h.useCase.Update(r.Context(), s); err != nil {
		slog.ErrorContext(r.Context(), "Failed to update supervisor",
			slog.String("error", err.Error()),
			slog.Int64("supervisorId", id))
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedSupervisor, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to get updated supervisor",
			slog.String("error", err.Error()),
			slog.Int64("supervisorId", id))
		response.Error(w, http.StatusInternalServerError, "Failed to get updated supervisor")
		return
	}

	response.JSON(w, http.StatusOK, toSupervisorResponse(updatedSupervisor))
}
