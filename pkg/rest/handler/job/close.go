package job

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Close maneja la solicitud de cierre de un job
// @Summary Cerrar trabajo
// @Description Cierra un trabajo estableciendo closed=true y actualizando el job_status_id
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path int true "ID del trabajo"
// @Param close body CloseJobRequest true "Datos para cerrar el trabajo"
// @Success 200 {object} JobResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs/{id}/close [put]
// @Security BearerAuth
func (h *Handler) Close(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req CloseJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(r.Context(), "Invalid request body",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	if err := h.useCase.Close(r.Context(), id, req.JobStatusID); err != nil {
		switch err {
		case domainJob.ErrJobNotFound:
			response.Error(w, http.StatusNotFound, "Trabajo no encontrado")
		case domainJob.ErrJobAlreadyClosed:
			response.Error(w, http.StatusConflict, "El trabajo ya está cerrado")
		case domainJob.ErrInvalidJobStatus:
			response.Error(w, http.StatusBadRequest, "Estado de trabajo inválido")
		default:
			slog.ErrorContext(r.Context(), "Failed to close job",
				slog.Int64("id", id),
				slog.String("error", err.Error()))
			response.Error(w, http.StatusInternalServerError, "Error al cerrar trabajo")
		}
		return
	}

	// Re-fetch para obtener datos actualizados
	updated, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		response.JSON(w, http.StatusOK, map[string]string{"message": "Trabajo cerrado exitosamente"})
		return
	}

	response.JSON(w, http.StatusOK, toJobResponse(updated))
}
