package job_priority

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_priority"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase job_priority.Service
}

func NewHandler(useCase job_priority.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/job-priorities", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateJobPriorityRequest representa la solicitud para crear una prioridad de trabajo
type CreateJobPriorityRequest struct {
	Label    string  `json:"label" validate:"required"`
	Order    int     `json:"order"`
	Class    *string `json:"class,omitempty"`
	IsActive bool    `json:"isActive"`
}

// UpdateJobPriorityRequest representa la solicitud para actualizar una prioridad de trabajo
type UpdateJobPriorityRequest struct {
	Label    string  `json:"label" validate:"required"`
	Order    int     `json:"order"`
	Class    *string `json:"class,omitempty"`
	IsActive bool    `json:"isActive"`
}

// JobPriorityResponse representa la respuesta de una prioridad de trabajo
type JobPriorityResponse struct {
	ID        int64   `json:"id"`
	Label     string  `json:"label"`
	Order     int     `json:"order"`
	Class     *string `json:"class,omitempty"`
	IsActive  bool    `json:"isActive"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

func toResponse(p *job_priority.JobPriority) JobPriorityResponse {
	resp := JobPriorityResponse{
		ID:       p.ID,
		Label:    p.Label,
		Order:    p.Order,
		Class:    p.Class,
		IsActive: p.IsActive,
	}

	if p.CreatedAt != nil {
		resp.CreatedAt = p.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if p.UpdatedAt != nil {
		resp.UpdatedAt = p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
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

// List maneja la solicitud de listado de prioridades de trabajo
// @Summary Listar prioridades de trabajo
// @Description Obtiene una lista paginada de prioridades de trabajo con filtros opcionales
// @Tags JobPriorities
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label o class"
// @Param isActive query bool false "Filtrar por estado activo"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-priorities [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	priorities, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar prioridades de trabajo")
		return
	}

	items := make([]JobPriorityResponse, len(priorities))
	for i, p := range priorities {
		items[i] = toResponse(p)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// Create maneja la solicitud de creación de una prioridad de trabajo
// @Summary Crear prioridad de trabajo
// @Description Crea una nueva prioridad de trabajo. El campo class acepta colores Bootstrap: blue, indigo, purple, pink, red, orange, yellow, green, teal, cyan, dark, light
// @Tags JobPriorities
// @Accept json
// @Produce json
// @Param priority body CreateJobPriorityRequest true "Datos de la prioridad de trabajo"
// @Success 201 {object} JobPriorityResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-priorities [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateJobPriorityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	priority := &job_priority.JobPriority{
		Label:    req.Label,
		Order:    req.Order,
		Class:    req.Class,
		IsActive: req.IsActive,
	}

	if err := h.useCase.Create(r.Context(), priority); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(priority))
}

// Get maneja la solicitud de obtención de una prioridad de trabajo por ID
// @Summary Obtener prioridad de trabajo
// @Description Obtiene una prioridad de trabajo por su ID
// @Tags JobPriorities
// @Accept json
// @Produce json
// @Param id path int true "ID de la prioridad de trabajo"
// @Success 200 {object} JobPriorityResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-priorities/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	priority, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == job_priority.ErrJobPriorityNotFound {
			response.Error(w, http.StatusNotFound, "Prioridad de trabajo no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener prioridad de trabajo")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(priority))
}

// Update maneja la solicitud de actualización de una prioridad de trabajo
// @Summary Actualizar prioridad de trabajo
// @Description Actualiza una prioridad de trabajo existente
// @Tags JobPriorities
// @Accept json
// @Produce json
// @Param id path int true "ID de la prioridad de trabajo"
// @Param priority body UpdateJobPriorityRequest true "Datos de la prioridad de trabajo"
// @Success 200 {object} JobPriorityResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-priorities/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateJobPriorityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	priority := &job_priority.JobPriority{
		ID:       id,
		Label:    req.Label,
		Order:    req.Order,
		Class:    req.Class,
		IsActive: req.IsActive,
	}

	if err := h.useCase.Update(r.Context(), priority); err != nil {
		if err == job_priority.ErrJobPriorityNotFound {
			response.Error(w, http.StatusNotFound, "Prioridad de trabajo no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(priority))
}

// Delete maneja la solicitud de eliminación de una prioridad de trabajo
// @Summary Eliminar prioridad de trabajo
// @Description Elimina una prioridad de trabajo. No se puede eliminar si tiene trabajos asociados
// @Tags JobPriorities
// @Accept json
// @Produce json
// @Param id path int true "ID de la prioridad de trabajo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-priorities/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == job_priority.ErrJobPriorityNotFound {
			response.Error(w, http.StatusNotFound, "Prioridad de trabajo no encontrada")
			return
		}
		if err == job_priority.ErrJobPriorityInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar prioridad de trabajo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
