package warranty_type

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/angumol/jvairv2/pkg/domain/warranty_type"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	useCase warranty_type.Service
}

func NewHandler(useCase warranty_type.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/warranty-types", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateWarrantyTypeRequest representa la solicitud para crear un tipo de garantía
type CreateWarrantyTypeRequest struct {
	Label       string `json:"label" example:"Parts"`
	LabelPlural string `json:"labelPlural" example:"Parts"`
	IsActive    *bool  `json:"isActive,omitempty" example:"true"`
}

// UpdateWarrantyTypeRequest representa la solicitud para actualizar un tipo de garantía
type UpdateWarrantyTypeRequest struct {
	Label       string `json:"label" example:"Parts"`
	LabelPlural string `json:"labelPlural" example:"Parts"`
	IsActive    *bool  `json:"isActive,omitempty" example:"true"`
}

// WarrantyTypeResponse representa la respuesta de un tipo de garantía
type WarrantyTypeResponse struct {
	ID          int64  `json:"id" example:"1"`
	Label       string `json:"label" example:"Parts"`
	LabelPlural string `json:"labelPlural" example:"Parts"`
	IsActive    bool   `json:"isActive" example:"true"`
	CreatedAt   string `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   string `json:"updatedAt,omitempty" example:"2024-01-15T10:30:00Z"`
}

func toResponse(wt *warranty_type.WarrantyType) WarrantyTypeResponse {
	resp := WarrantyTypeResponse{
		ID:          wt.ID,
		Label:       wt.Label,
		LabelPlural: wt.LabelPlural,
		IsActive:    wt.IsActive,
	}

	if wt.CreatedAt != nil {
		resp.CreatedAt = wt.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if wt.UpdatedAt != nil {
		resp.UpdatedAt = wt.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return resp
}

// List maneja la solicitud de listado de tipos de garantía
// @Summary Listar tipos de garantía
// @Description Obtiene una lista paginada de tipos de garantía
// @Tags WarrantyTypes
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-types [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := make(map[string]interface{})
	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	types, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar tipos de garantía")
		return
	}

	items := make([]WarrantyTypeResponse, len(types))
	for i, wt := range types {
		items[i] = toResponse(wt)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// Create maneja la solicitud de creación de un tipo de garantía
// @Summary Crear tipo de garantía
// @Description Crea un nuevo tipo de garantía
// @Tags WarrantyTypes
// @Accept json
// @Produce json
// @Param warrantyType body CreateWarrantyTypeRequest true "Datos del tipo de garantía"
// @Success 201 {object} WarrantyTypeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-types [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWarrantyTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	wt := &warranty_type.WarrantyType{
		Label:       req.Label,
		LabelPlural: req.LabelPlural,
		IsActive:    isActive,
	}

	if err := h.useCase.Create(r.Context(), wt); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(wt))
}

// Get maneja la solicitud de obtención de un tipo de garantía por ID
// @Summary Obtener tipo de garantía
// @Description Obtiene un tipo de garantía por su ID
// @Tags WarrantyTypes
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de garantía"
// @Success 200 {object} WarrantyTypeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-types/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	wt, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == warranty_type.ErrWarrantyTypeNotFound {
			response.Error(w, http.StatusNotFound, "Tipo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener tipo de garantía")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(wt))
}

// Update maneja la solicitud de actualización de un tipo de garantía
// @Summary Actualizar tipo de garantía
// @Description Actualiza un tipo de garantía existente
// @Tags WarrantyTypes
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de garantía"
// @Param warrantyType body UpdateWarrantyTypeRequest true "Datos del tipo de garantía"
// @Success 200 {object} WarrantyTypeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-types/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateWarrantyTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	wt := &warranty_type.WarrantyType{
		ID:          id,
		Label:       req.Label,
		LabelPlural: req.LabelPlural,
		IsActive:    isActive,
	}

	if err := h.useCase.Update(r.Context(), wt); err != nil {
		if err == warranty_type.ErrWarrantyTypeNotFound {
			response.Error(w, http.StatusNotFound, "Tipo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(wt))
}

// Delete maneja la solicitud de eliminación de un tipo de garantía
// @Summary Eliminar tipo de garantía
// @Description Elimina un tipo de garantía. No se puede eliminar si tiene garantías asociadas
// @Tags WarrantyTypes
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de garantía"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-types/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == warranty_type.ErrWarrantyTypeNotFound {
			response.Error(w, http.StatusNotFound, "Tipo de garantía no encontrado")
			return
		}
		if err == warranty_type.ErrWarrantyTypeInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar tipo de garantía")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
