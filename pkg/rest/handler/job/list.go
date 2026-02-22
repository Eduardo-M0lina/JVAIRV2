package job

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List maneja la solicitud de listado de jobs
// @Summary Listar trabajos
// @Description Obtiene una lista paginada de trabajos con filtros opcionales. Por defecto muestra solo trabajos abiertos (closed=0)
// @Tags Jobs
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda en work_order, property, customer"
// @Param closed query string false "Filtrar por cerrado: 0, 1, all" default(0)
// @Param jobCategoryId query int false "Filtrar por categoría de trabajo"
// @Param jobStatusId query int false "Filtrar por estado de trabajo"
// @Param jobPriorityId query int false "Filtrar por prioridad de trabajo"
// @Param userId query string false "Filtrar por usuario asignado (ID o 'unassigned')"
// @Param propertyId query int false "Filtrar por propiedad"
// @Param workflowId query int false "Filtrar por workflow"
// @Param sort query string false "Campo de ordenamiento: work_order, date_received, created_at, due_date, status"
// @Param direction query string false "Dirección de ordenamiento: asc, desc" default(desc)
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/jobs [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	filters := parseFilters(r)

	jobs, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to list jobs",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al listar trabajos")
		return
	}

	items := make([]JobResponse, len(jobs))
	for i, j := range jobs {
		items[i] = toJobResponse(j)
	}

	totalPages := (total + pageSize - 1) / pageSize

	response.JSON(w, http.StatusOK, response.PaginatedResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}
