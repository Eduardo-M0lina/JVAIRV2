package warranty_claim_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/warranty_claim_status"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase warranty_claim_status.Service
}

func NewHandler(useCase warranty_claim_status.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/warranty-claim-statuses", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

type CreateWarrantyClaimStatusRequest struct {
	Label    string  `json:"label" example:"Pending"`
	Class    *string `json:"class,omitempty" example:"badge-warning"`
	Order    int     `json:"order" example:"1"`
	IsActive *bool   `json:"isActive,omitempty" example:"true"`
}

type UpdateWarrantyClaimStatusRequest struct {
	Label    string  `json:"label" example:"Pending"`
	Class    *string `json:"class,omitempty" example:"badge-warning"`
	Order    int     `json:"order" example:"1"`
	IsActive *bool   `json:"isActive,omitempty" example:"true"`
}

type WarrantyClaimStatusResponse struct {
	ID        int64   `json:"id" example:"1"`
	Label     string  `json:"label" example:"Pending"`
	Class     *string `json:"class,omitempty" example:"badge-warning"`
	Order     int     `json:"order" example:"1"`
	IsActive  bool    `json:"isActive" example:"true"`
	CreatedAt string  `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt string  `json:"updatedAt,omitempty" example:"2024-01-15T10:30:00Z"`
}

func toResponse(wcs *warranty_claim_status.WarrantyClaimStatus) WarrantyClaimStatusResponse {
	resp := WarrantyClaimStatusResponse{
		ID:       wcs.ID,
		Label:    wcs.Label,
		Class:    wcs.Class,
		Order:    wcs.Order,
		IsActive: wcs.IsActive,
	}
	if wcs.CreatedAt != nil {
		resp.CreatedAt = wcs.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if wcs.UpdatedAt != nil {
		resp.UpdatedAt = wcs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return resp
}

// @Summary Listar estados de reclamo de garantía
// @Description Obtiene una lista paginada de estados de reclamo de garantía
// @Tags WarrantyClaimStatuses
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-statuses [get]
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
		response.Error(w, http.StatusInternalServerError, "Error al listar estados de reclamo de garantía")
		return
	}

	items := make([]WarrantyClaimStatusResponse, len(statuses))
	for i, wcs := range statuses {
		items[i] = toResponse(wcs)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// @Summary Crear estado de reclamo de garantía
// @Description Crea un nuevo estado de reclamo de garantía
// @Tags WarrantyClaimStatuses
// @Accept json
// @Produce json
// @Param warrantyClaimStatus body CreateWarrantyClaimStatusRequest true "Datos del estado de reclamo"
// @Success 201 {object} WarrantyClaimStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-statuses [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWarrantyClaimStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	wcs := &warranty_claim_status.WarrantyClaimStatus{
		Label:    req.Label,
		Class:    req.Class,
		Order:    req.Order,
		IsActive: isActive,
	}

	if err := h.useCase.Create(r.Context(), wcs); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(wcs))
}

// @Summary Obtener estado de reclamo de garantía
// @Description Obtiene un estado de reclamo de garantía por su ID
// @Tags WarrantyClaimStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de reclamo"
// @Success 200 {object} WarrantyClaimStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-statuses/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	wcs, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == warranty_claim_status.ErrWarrantyClaimStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de reclamo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener estado de reclamo de garantía")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(wcs))
}

// @Summary Actualizar estado de reclamo de garantía
// @Description Actualiza un estado de reclamo de garantía existente
// @Tags WarrantyClaimStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de reclamo"
// @Param warrantyClaimStatus body UpdateWarrantyClaimStatusRequest true "Datos del estado de reclamo"
// @Success 200 {object} WarrantyClaimStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-statuses/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateWarrantyClaimStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	wcs := &warranty_claim_status.WarrantyClaimStatus{
		ID:       id,
		Label:    req.Label,
		Class:    req.Class,
		Order:    req.Order,
		IsActive: isActive,
	}

	if err := h.useCase.Update(r.Context(), wcs); err != nil {
		if err == warranty_claim_status.ErrWarrantyClaimStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de reclamo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(wcs))
}

// @Summary Eliminar estado de reclamo de garantía
// @Description Elimina un estado de reclamo de garantía. No se puede eliminar si tiene reclamos asociados
// @Tags WarrantyClaimStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de reclamo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-statuses/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == warranty_claim_status.ErrWarrantyClaimStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de reclamo de garantía no encontrado")
			return
		}
		if err == warranty_claim_status.ErrWarrantyClaimStatusInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar estado de reclamo de garantía")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
