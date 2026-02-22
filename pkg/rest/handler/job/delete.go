package job

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Delete maneja la solicitud de eliminación de un job (soft delete)
// @Summary Eliminar trabajo
// @Description Elimina un trabajo (soft delete)
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path int true "ID del trabajo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == domainJob.ErrJobNotFound {
			response.Error(w, http.StatusNotFound, "Trabajo no encontrado")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to delete job",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al eliminar trabajo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
