package technician_job_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/technician_job_status"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase technician_job_status.Service
}

func NewHandler(useCase technician_job_status.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/technician-job-statuses", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateTechnicianJobStatusRequest representa la solicitud para crear un estado de técnico de trabajo
type CreateTechnicianJobStatusRequest struct {
	Label       string  `json:"label" validate:"required"`
	Class       *string `json:"class,omitempty"`
	JobStatusID *int64  `json:"jobStatusId,omitempty"`
	IsActive    bool    `json:"isActive"`
}

// UpdateTechnicianJobStatusRequest representa la solicitud para actualizar un estado de técnico de trabajo
type UpdateTechnicianJobStatusRequest struct {
	Label       string  `json:"label" validate:"required"`
	Class       *string `json:"class,omitempty"`
	JobStatusID *int64  `json:"jobStatusId,omitempty"`
	IsActive    bool    `json:"isActive"`
}

// TechnicianJobStatusResponse representa la respuesta de un estado de técnico de trabajo
type TechnicianJobStatusResponse struct {
	ID          int64   `json:"id"`
	Label       string  `json:"label"`
	Class       *string `json:"class,omitempty"`
	JobStatusID *int64  `json:"jobStatusId,omitempty"`
	IsActive    bool    `json:"isActive"`
	CreatedAt   string  `json:"createdAt,omitempty"`
	UpdatedAt   string  `json:"updatedAt,omitempty"`
}

func toResponse(s *technician_job_status.TechnicianJobStatus) TechnicianJobStatusResponse {
	resp := TechnicianJobStatusResponse{
		ID:          s.ID,
		Label:       s.Label,
		Class:       s.Class,
		JobStatusID: s.JobStatusID,
		IsActive:    s.IsActive,
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

	if jobStatusIDStr := r.URL.Query().Get("jobStatusId"); jobStatusIDStr != "" {
		if jobStatusID, err := strconv.ParseInt(jobStatusIDStr, 10, 64); err == nil {
			filters["job_status_id"] = jobStatusID
		}
	}

	return filters
}

// List maneja la solicitud de listado de estados de técnico de trabajo
// @Summary Listar estados de técnico de trabajo
// @Description Obtiene una lista paginada de estados de técnico de trabajo con filtros opcionales
// @Tags TechnicianJobStatuses
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label o class"
// @Param isActive query bool false "Filtrar por estado activo"
// @Param jobStatusId query int false "Filtrar por ID de estado de trabajo"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/technician-job-statuses [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	statuses, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar estados de técnico de trabajo")
		return
	}

	items := make([]TechnicianJobStatusResponse, len(statuses))
	for i, s := range statuses {
		items[i] = toResponse(s)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// Create maneja la solicitud de creación de un estado de técnico de trabajo
// @Summary Crear estado de técnico de trabajo
// @Description Crea un nuevo estado de técnico de trabajo. El campo class acepta colores Bootstrap. El campo jobStatusId es opcional y referencia un job_status existente
// @Tags TechnicianJobStatuses
// @Accept json
// @Produce json
// @Param status body CreateTechnicianJobStatusRequest true "Datos del estado de técnico de trabajo"
// @Success 201 {object} TechnicianJobStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/technician-job-statuses [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTechnicianJobStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	status := &technician_job_status.TechnicianJobStatus{
		Label:       req.Label,
		Class:       req.Class,
		JobStatusID: req.JobStatusID,
		IsActive:    req.IsActive,
	}

	if err := h.useCase.Create(r.Context(), status); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(status))
}

// Get maneja la solicitud de obtención de un estado de técnico de trabajo por ID
// @Summary Obtener estado de técnico de trabajo
// @Description Obtiene un estado de técnico de trabajo por su ID
// @Tags TechnicianJobStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de técnico de trabajo"
// @Success 200 {object} TechnicianJobStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/technician-job-statuses/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	status, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == technician_job_status.ErrTechnicianJobStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de técnico de trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener estado de técnico de trabajo")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(status))
}

// Update maneja la solicitud de actualización de un estado de técnico de trabajo
// @Summary Actualizar estado de técnico de trabajo
// @Description Actualiza un estado de técnico de trabajo existente
// @Tags TechnicianJobStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de técnico de trabajo"
// @Param status body UpdateTechnicianJobStatusRequest true "Datos del estado de técnico de trabajo"
// @Success 200 {object} TechnicianJobStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/technician-job-statuses/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateTechnicianJobStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	status := &technician_job_status.TechnicianJobStatus{
		ID:          id,
		Label:       req.Label,
		Class:       req.Class,
		JobStatusID: req.JobStatusID,
		IsActive:    req.IsActive,
	}

	if err := h.useCase.Update(r.Context(), status); err != nil {
		if err == technician_job_status.ErrTechnicianJobStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de técnico de trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(status))
}

// Delete maneja la solicitud de eliminación de un estado de técnico de trabajo
// @Summary Eliminar estado de técnico de trabajo
// @Description Elimina un estado de técnico de trabajo
// @Tags TechnicianJobStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de técnico de trabajo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/technician-job-statuses/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == technician_job_status.ErrTechnicianJobStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de técnico de trabajo no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar estado de técnico de trabajo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
