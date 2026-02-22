package job_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_status"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase job_status.Service
}

func NewHandler(useCase job_status.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/job-statuses", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateJobStatusRequest representa la solicitud para crear un estado de trabajo
type CreateJobStatusRequest struct {
	Label    string  `json:"label" validate:"required"`
	Class    *string `json:"class,omitempty"`
	IsActive bool    `json:"isActive"`
}

// UpdateJobStatusRequest representa la solicitud para actualizar un estado de trabajo
type UpdateJobStatusRequest struct {
	Label    string  `json:"label" validate:"required"`
	Class    *string `json:"class,omitempty"`
	IsActive bool    `json:"isActive"`
}

// JobStatusResponse representa la respuesta de un estado de trabajo
type JobStatusResponse struct {
	ID        int64   `json:"id"`
	Label     string  `json:"label"`
	Class     *string `json:"class,omitempty"`
	IsActive  bool    `json:"isActive"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

func toResponse(s *job_status.JobStatus) JobStatusResponse {
	resp := JobStatusResponse{
		ID:       s.ID,
		Label:    s.Label,
		Class:    s.Class,
		IsActive: s.IsActive,
	}

	if s.CreatedAt != nil {
		resp.CreatedAt = s.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if s.UpdatedAt != nil {
		resp.UpdatedAt = s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return resp
}

func parseFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	if isActiveStr := r.URL.Query().Get("isActive"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			filters["is_active"] = isActive
		}
	}

	return filters
}

// List maneja la solicitud de listado de estados de trabajo
// @Summary Listar estados de trabajo
// @Description Obtiene una lista paginada de estados de trabajo con filtros opcionales
// @Tags JobStatuses
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label o class"
// @Param isActive query bool false "Filtrar por estado activo"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-statuses [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	statuses, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar estados de trabajo")
		return
	}

	items := make([]JobStatusResponse, len(statuses))
	for i, s := range statuses {
		items[i] = toResponse(s)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// Create maneja la solicitud de creación de un estado de trabajo
// @Summary Crear estado de trabajo
// @Description Crea un nuevo estado de trabajo. El campo class acepta colores Bootstrap: blue, indigo, purple, pink, red, orange, yellow, green, teal, cyan, dark, light
// @Tags JobStatuses
// @Accept json
// @Produce json
// @Param status body CreateJobStatusRequest true "Datos del estado de trabajo"
// @Success 201 {object} JobStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-statuses [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateJobStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	status := &job_status.JobStatus{
		Label:    req.Label,
		Class:    req.Class,
		IsActive: req.IsActive,
	}

	if err := h.useCase.Create(r.Context(), status); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(status))
}

// Get maneja la solicitud de obtención de un estado de trabajo por ID
// @Summary Obtener estado de trabajo
// @Description Obtiene un estado de trabajo por su ID
// @Tags JobStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de trabajo"
// @Success 200 {object} JobStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-statuses/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	status, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == job_status.ErrJobStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener estado de trabajo")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(status))
}

// Update maneja la solicitud de actualización de un estado de trabajo
// @Summary Actualizar estado de trabajo
// @Description Actualiza un estado de trabajo existente
// @Tags JobStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de trabajo"
// @Param status body UpdateJobStatusRequest true "Datos del estado de trabajo"
// @Success 200 {object} JobStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-statuses/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateJobStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	status := &job_status.JobStatus{
		ID:       id,
		Label:    req.Label,
		Class:    req.Class,
		IsActive: req.IsActive,
	}

	if err := h.useCase.Update(r.Context(), status); err != nil {
		if err == job_status.ErrJobStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(status))
}

// Delete maneja la solicitud de eliminación de un estado de trabajo
// @Summary Eliminar estado de trabajo
// @Description Elimina un estado de trabajo
// @Tags JobStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de trabajo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-statuses/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == job_status.ErrJobStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar estado de trabajo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
