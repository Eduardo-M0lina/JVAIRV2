package sms_template

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/domain/sms_template"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	useCase sms_template.Service
}

func NewHandler(useCase sms_template.Service) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/sms-templates", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

type CreateRequest struct {
	Label    string `json:"label"`
	Message  string `json:"message"`
	IsActive bool   `json:"isActive"`
}

type UpdateRequest struct {
	Label    string `json:"label"`
	Message  string `json:"message"`
	IsActive bool   `json:"isActive"`
}

type ItemResponse struct {
	ID        int64  `json:"id"`
	Label     string `json:"label"`
	Message   string `json:"message"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

func toResponse(item *sms_template.SMSTemplate) ItemResponse {
	resp := ItemResponse{ID: item.ID, Label: item.Label, Message: item.Message, IsActive: item.IsActive}
	if item.CreatedAt != nil {
		resp.CreatedAt = item.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if item.UpdatedAt != nil {
		resp.UpdatedAt = item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
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

// List godoc
// @Summary Listar plantillas de SMS
// @Description Obtiene una lista paginada de plantillas de SMS con filtros opcionales
// @Tags SMSTemplates
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(15)
// @Param search query string false "Búsqueda por label o message"
// @Param isActive query bool false "Filtrar por estado activo"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/sms-templates [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	items, total, err := h.useCase.List(r.Context(), parseFilters(r), page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar plantillas de SMS")
		return
	}
	result := make([]ItemResponse, len(items))
	for i, item := range items {
		result[i] = toResponse(item)
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 15
	}
	response.Paginated(w, result, page, pageSize, total)
}

// Create godoc
// @Summary Crear plantilla de SMS
// @Description Crea una nueva plantilla de SMS
// @Tags SMSTemplates
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Datos de la plantilla de SMS"
// @Success 201 {object} ItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/sms-templates [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}
	item := &sms_template.SMSTemplate{Label: req.Label, Message: req.Message, IsActive: req.IsActive}
	if err := h.useCase.Create(r.Context(), item); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, toResponse(item))
}

// Get godoc
// @Summary Obtener plantilla de SMS
// @Description Obtiene una plantilla de SMS por su ID
// @Tags SMSTemplates
// @Accept json
// @Produce json
// @Param id path int true "ID de la plantilla de SMS"
// @Success 200 {object} ItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/sms-templates/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}
	item, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == sms_template.ErrSMSTemplateNotFound {
			response.Error(w, http.StatusNotFound, "Plantilla de SMS no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener plantilla de SMS")
		return
	}
	response.JSON(w, http.StatusOK, toResponse(item))
}

// Update godoc
// @Summary Actualizar plantilla de SMS
// @Description Actualiza una plantilla de SMS existente
// @Tags SMSTemplates
// @Accept json
// @Produce json
// @Param id path int true "ID de la plantilla de SMS"
// @Param request body UpdateRequest true "Datos de la plantilla de SMS"
// @Success 200 {object} ItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/sms-templates/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}
	item := &sms_template.SMSTemplate{ID: id, Label: req.Label, Message: req.Message, IsActive: req.IsActive}
	if err := h.useCase.Update(r.Context(), item); err != nil {
		if err == sms_template.ErrSMSTemplateNotFound {
			response.Error(w, http.StatusNotFound, "Plantilla de SMS no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	updated, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener plantilla de SMS actualizada")
		return
	}
	response.JSON(w, http.StatusOK, toResponse(updated))
}

// Delete godoc
// @Summary Eliminar plantilla de SMS
// @Description Elimina una plantilla de SMS
// @Tags SMSTemplates
// @Accept json
// @Produce json
// @Param id path int true "ID de la plantilla de SMS"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/sms-templates/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}
	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == sms_template.ErrSMSTemplateNotFound {
			response.Error(w, http.StatusNotFound, "Plantilla de SMS no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar plantilla de SMS")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
