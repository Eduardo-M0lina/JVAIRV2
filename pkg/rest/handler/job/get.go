package job

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Get maneja la solicitud de obtención de un job por ID
// @Summary Obtener trabajo
// @Description Obtiene un trabajo por su ID
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path int true "ID del trabajo"
// @Success 200 {object} JobResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	j, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domainJob.ErrJobNotFound {
			response.Error(w, http.StatusNotFound, "Trabajo no encontrado")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get job",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al obtener trabajo")
		return
	}

	response.JSON(w, http.StatusOK, toJobResponse(j))
}
