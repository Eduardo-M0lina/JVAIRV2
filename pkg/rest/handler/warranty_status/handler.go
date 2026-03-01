package warranty_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/warranty_status"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase warranty_status.Service
}

func NewHandler(useCase warranty_status.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/warranty-statuses", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

type CreateWarrantyStatusRequest struct {
	Label    string  `json:"label" example:"Pending"`
	Class    *string `json:"class,omitempty" example:"badge-warning"`
	Order    int     `json:"order" example:"1"`
	IsActive *bool   `json:"isActive,omitempty" example:"true"`
}

type UpdateWarrantyStatusRequest struct {
	Label    string  `json:"label" example:"Pending"`
	Class    *string `json:"class,omitempty" example:"badge-warning"`
	Order    int     `json:"order" example:"1"`
	IsActive *bool   `json:"isActive,omitempty" example:"true"`
}

type WarrantyStatusResponse struct {
	ID        int64   `json:"id" example:"1"`
	Label     string  `json:"label" example:"Pending"`
	Class     *string `json:"class,omitempty" example:"badge-warning"`
	Order     int     `json:"order" example:"1"`
	IsActive  bool    `json:"isActive" example:"true"`
	CreatedAt string  `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt string  `json:"updatedAt,omitempty" example:"2024-01-15T10:30:00Z"`
}

func toResponse(ws *warranty_status.WarrantyStatus) WarrantyStatusResponse {
	resp := WarrantyStatusResponse{
		ID:       ws.ID,
		Label:    ws.Label,
		Class:    ws.Class,
		Order:    ws.Order,
		IsActive: ws.IsActive,
	}
	if ws.CreatedAt != nil {
		resp.CreatedAt = ws.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if ws.UpdatedAt != nil {
		resp.UpdatedAt = ws.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return resp
}

// @Summary Listar estados de garantía
// @Description Obtiene una lista paginada de estados de garantía
// @Tags WarrantyStatuses
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-statuses [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := make(map[string]interface{})
	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	statuses, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar estados de garantía")
		return
	}

	items := make([]WarrantyStatusResponse, len(statuses))
	for i, ws := range statuses {
		items[i] = toResponse(ws)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// @Summary Crear estado de garantía
// @Description Crea un nuevo estado de garantía
// @Tags WarrantyStatuses
// @Accept json
// @Produce json
// @Param warrantyStatus body CreateWarrantyStatusRequest true "Datos del estado de garantía"
// @Success 201 {object} WarrantyStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-statuses [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWarrantyStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	ws := &warranty_status.WarrantyStatus{
		Label:    req.Label,
		Class:    req.Class,
		Order:    req.Order,
		IsActive: isActive,
	}

	if err := h.useCase.Create(r.Context(), ws); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(ws))
}

// @Summary Obtener estado de garantía
// @Description Obtiene un estado de garantía por su ID
// @Tags WarrantyStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de garantía"
// @Success 200 {object} WarrantyStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-statuses/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	ws, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == warranty_status.ErrWarrantyStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener estado de garantía")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(ws))
}

// @Summary Actualizar estado de garantía
// @Description Actualiza un estado de garantía existente
// @Tags WarrantyStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de garantía"
// @Param warrantyStatus body UpdateWarrantyStatusRequest true "Datos del estado de garantía"
// @Success 200 {object} WarrantyStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-statuses/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateWarrantyStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	ws := &warranty_status.WarrantyStatus{
		ID:       id,
		Label:    req.Label,
		Class:    req.Class,
		Order:    req.Order,
		IsActive: isActive,
	}

	if err := h.useCase.Update(r.Context(), ws); err != nil {
		if err == warranty_status.ErrWarrantyStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(ws))
}

// @Summary Eliminar estado de garantía
// @Description Elimina un estado de garantía. No se puede eliminar si tiene garantías asociadas
// @Tags WarrantyStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de garantía"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-statuses/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == warranty_status.ErrWarrantyStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de garantía no encontrado")
			return
		}
		if err == warranty_status.ErrWarrantyStatusInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar estado de garantía")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
