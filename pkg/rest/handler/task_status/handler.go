package task_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/task_status"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase task_status.Service
}

func NewHandler(useCase task_status.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/task-statuses", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateTaskStatusRequest representa la solicitud para crear un estado de tarea
type CreateTaskStatusRequest struct {
	Label    string  `json:"label" validate:"required"`
	Class    *string `json:"class,omitempty"`
	Order    int     `json:"order"`
	IsActive bool    `json:"isActive"`
}

// UpdateTaskStatusRequest representa la solicitud para actualizar un estado de tarea
type UpdateTaskStatusRequest struct {
	Label    string  `json:"label" validate:"required"`
	Class    *string `json:"class,omitempty"`
	Order    int     `json:"order"`
	IsActive bool    `json:"isActive"`
}

// TaskStatusResponse representa la respuesta de un estado de tarea
type TaskStatusResponse struct {
	ID        int64   `json:"id"`
	Label     string  `json:"label"`
	Class     *string `json:"class,omitempty"`
	Order     int     `json:"order"`
	IsActive  bool    `json:"isActive"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

func toResponse(s *task_status.TaskStatus) TaskStatusResponse {
	resp := TaskStatusResponse{
		ID:       s.ID,
		Label:    s.Label,
		Class:    s.Class,
		Order:    s.Order,
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

// List maneja la solicitud de listado de estados de tarea
// @Summary Listar estados de tarea
// @Description Obtiene una lista paginada de estados de tarea con filtros opcionales
// @Tags TaskStatuses
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label o class"
// @Param isActive query bool false "Filtrar por estado activo"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/task-statuses [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	statuses, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar estados de tarea")
		return
	}

	items := make([]TaskStatusResponse, len(statuses))
	for i, s := range statuses {
		items[i] = toResponse(s)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// Create maneja la solicitud de creación de un estado de tarea
// @Summary Crear estado de tarea
// @Description Crea un nuevo estado de tarea. El campo class acepta colores Bootstrap: blue, indigo, purple, pink, red, orange, yellow, green, teal, cyan, dark, light
// @Tags TaskStatuses
// @Accept json
// @Produce json
// @Param status body CreateTaskStatusRequest true "Datos del estado de tarea"
// @Success 201 {object} TaskStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/task-statuses [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	status := &task_status.TaskStatus{
		Label:    req.Label,
		Class:    req.Class,
		Order:    req.Order,
		IsActive: req.IsActive,
	}

	if err := h.useCase.Create(r.Context(), status); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(status))
}

// Get maneja la solicitud de obtención de un estado de tarea por ID
// @Summary Obtener estado de tarea
// @Description Obtiene un estado de tarea por su ID
// @Tags TaskStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de tarea"
// @Success 200 {object} TaskStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/task-statuses/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	status, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == task_status.ErrTaskStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de tarea no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener estado de tarea")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(status))
}

// Update maneja la solicitud de actualización de un estado de tarea
// @Summary Actualizar estado de tarea
// @Description Actualiza un estado de tarea existente
// @Tags TaskStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de tarea"
// @Param status body UpdateTaskStatusRequest true "Datos del estado de tarea"
// @Success 200 {object} TaskStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/task-statuses/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	status := &task_status.TaskStatus{
		ID:       id,
		Label:    req.Label,
		Class:    req.Class,
		Order:    req.Order,
		IsActive: req.IsActive,
	}

	if err := h.useCase.Update(r.Context(), status); err != nil {
		if err == task_status.ErrTaskStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de tarea no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(status))
}

// Delete maneja la solicitud de eliminación de un estado de tarea
// @Summary Eliminar estado de tarea
// @Description Elimina un estado de tarea. No se puede eliminar si tiene tareas asociadas
// @Tags TaskStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de tarea"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/task-statuses/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == task_status.ErrTaskStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de tarea no encontrado")
			return
		}
		if err == task_status.ErrTaskStatusInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar estado de tarea")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
