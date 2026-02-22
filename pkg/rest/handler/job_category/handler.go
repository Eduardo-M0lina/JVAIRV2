package job_category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/job_category"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase job_category.Service
}

func NewHandler(useCase job_category.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/job-categories", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateJobCategoryRequest representa la solicitud para crear una categoría de trabajo
type CreateJobCategoryRequest struct {
	Label       string `json:"label" validate:"required"`
	LabelPlural string `json:"labelPlural" validate:"required"`
	Type        string `json:"type" validate:"required"`
	IsActive    bool   `json:"isActive"`
}

// UpdateJobCategoryRequest representa la solicitud para actualizar una categoría de trabajo
type UpdateJobCategoryRequest struct {
	Label       string `json:"label" validate:"required"`
	LabelPlural string `json:"labelPlural" validate:"required"`
	Type        string `json:"type" validate:"required"`
	IsActive    bool   `json:"isActive"`
}

// JobCategoryResponse representa la respuesta de una categoría de trabajo
type JobCategoryResponse struct {
	ID          int64  `json:"id"`
	Label       string `json:"label"`
	LabelPlural string `json:"labelPlural"`
	Type        string `json:"type"`
	IsActive    bool   `json:"isActive"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

func toResponse(c *job_category.JobCategory) JobCategoryResponse {
	resp := JobCategoryResponse{
		ID:          c.ID,
		Label:       c.Label,
		LabelPlural: c.LabelPlural,
		Type:        c.Type,
		IsActive:    c.IsActive,
	}

	if c.CreatedAt != nil {
		resp.CreatedAt = c.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if c.UpdatedAt != nil {
		resp.UpdatedAt = c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
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

	if typ := r.URL.Query().Get("type"); typ != "" {
		filters["type"] = typ
	}

	return filters
}

// List maneja la solicitud de listado de categorías de trabajo
// @Summary Listar categorías de trabajo
// @Description Obtiene una lista paginada de categorías de trabajo con filtros opcionales
// @Tags JobCategories
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label o tipo"
// @Param isActive query bool false "Filtrar por estado activo"
// @Param type query string false "Filtrar por tipo"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-categories [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	categories, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar categorías de trabajo")
		return
	}

	items := make([]JobCategoryResponse, len(categories))
	for i, c := range categories {
		items[i] = toResponse(c)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// Create maneja la solicitud de creación de una categoría de trabajo
// @Summary Crear categoría de trabajo
// @Description Crea una nueva categoría de trabajo
// @Tags JobCategories
// @Accept json
// @Produce json
// @Param category body CreateJobCategoryRequest true "Datos de la categoría de trabajo"
// @Success 201 {object} JobCategoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-categories [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateJobCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	category := &job_category.JobCategory{
		Label:       req.Label,
		LabelPlural: req.LabelPlural,
		Type:        req.Type,
		IsActive:    req.IsActive,
	}

	if err := h.useCase.Create(r.Context(), category); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(category))
}

// Get maneja la solicitud de obtención de una categoría de trabajo por ID
// @Summary Obtener categoría de trabajo
// @Description Obtiene una categoría de trabajo por su ID
// @Tags JobCategories
// @Accept json
// @Produce json
// @Param id path int true "ID de la categoría de trabajo"
// @Success 200 {object} JobCategoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-categories/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	category, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == job_category.ErrJobCategoryNotFound {
			response.Error(w, http.StatusNotFound, "Categoría de trabajo no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener categoría de trabajo")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(category))
}

// Update maneja la solicitud de actualización de una categoría de trabajo
// @Summary Actualizar categoría de trabajo
// @Description Actualiza una categoría de trabajo existente
// @Tags JobCategories
// @Accept json
// @Produce json
// @Param id path int true "ID de la categoría de trabajo"
// @Param category body UpdateJobCategoryRequest true "Datos de la categoría de trabajo"
// @Success 200 {object} JobCategoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-categories/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateJobCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	category := &job_category.JobCategory{
		ID:          id,
		Label:       req.Label,
		LabelPlural: req.LabelPlural,
		Type:        req.Type,
		IsActive:    req.IsActive,
	}

	if err := h.useCase.Update(r.Context(), category); err != nil {
		if err == job_category.ErrJobCategoryNotFound {
			response.Error(w, http.StatusNotFound, "Categoría de trabajo no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(category))
}

// Delete maneja la solicitud de eliminación de una categoría de trabajo
// @Summary Eliminar categoría de trabajo
// @Description Elimina una categoría de trabajo. No se puede eliminar si tiene trabajos asociados
// @Tags JobCategories
// @Accept json
// @Produce json
// @Param id path int true "ID de la categoría de trabajo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/job-categories/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == job_category.ErrJobCategoryNotFound {
			response.Error(w, http.StatusNotFound, "Categoría de trabajo no encontrada")
			return
		}
		if err == job_category.ErrJobCategoryInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar categoría de trabajo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
