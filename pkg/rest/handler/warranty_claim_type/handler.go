package warranty_claim_type

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/warranty_claim_type"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase warranty_claim_type.Service
}

func NewHandler(useCase warranty_claim_type.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/warranty-claim-types", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

type CreateWarrantyClaimTypeRequest struct {
	Label       string `json:"label" example:"Parts"`
	LabelPlural string `json:"labelPlural" example:"Parts"`
}

type UpdateWarrantyClaimTypeRequest struct {
	Label       string `json:"label" example:"Parts"`
	LabelPlural string `json:"labelPlural" example:"Parts"`
}

type WarrantyClaimTypeResponse struct {
	ID          int64  `json:"id" example:"1"`
	Label       string `json:"label" example:"Parts"`
	LabelPlural string `json:"labelPlural" example:"Parts"`
	CreatedAt   string `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   string `json:"updatedAt,omitempty" example:"2024-01-15T10:30:00Z"`
}

func toResponse(wct *warranty_claim_type.WarrantyClaimType) WarrantyClaimTypeResponse {
	resp := WarrantyClaimTypeResponse{
		ID:          wct.ID,
		Label:       wct.Label,
		LabelPlural: wct.LabelPlural,
	}
	if wct.CreatedAt != nil {
		resp.CreatedAt = wct.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if wct.UpdatedAt != nil {
		resp.UpdatedAt = wct.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return resp
}

// @Summary Listar tipos de reclamo de garantía
// @Description Obtiene una lista paginada de tipos de reclamo de garantía
// @Tags WarrantyClaimTypes
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-types [get]
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
		response.Error(w, http.StatusInternalServerError, "Error al listar tipos de reclamo de garantía")
		return
	}

	items := make([]WarrantyClaimTypeResponse, len(types))
	for i, wct := range types {
		items[i] = toResponse(wct)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// @Summary Crear tipo de reclamo de garantía
// @Description Crea un nuevo tipo de reclamo de garantía
// @Tags WarrantyClaimTypes
// @Accept json
// @Produce json
// @Param warrantyClaimType body CreateWarrantyClaimTypeRequest true "Datos del tipo de reclamo"
// @Success 201 {object} WarrantyClaimTypeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-types [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWarrantyClaimTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	wct := &warranty_claim_type.WarrantyClaimType{
		Label:       req.Label,
		LabelPlural: req.LabelPlural,
	}

	if err := h.useCase.Create(r.Context(), wct); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(wct))
}

// @Summary Obtener tipo de reclamo de garantía
// @Description Obtiene un tipo de reclamo de garantía por su ID
// @Tags WarrantyClaimTypes
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de reclamo"
// @Success 200 {object} WarrantyClaimTypeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-types/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	wct, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == warranty_claim_type.ErrWarrantyClaimTypeNotFound {
			response.Error(w, http.StatusNotFound, "Tipo de reclamo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener tipo de reclamo de garantía")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(wct))
}

// @Summary Actualizar tipo de reclamo de garantía
// @Description Actualiza un tipo de reclamo de garantía existente
// @Tags WarrantyClaimTypes
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de reclamo"
// @Param warrantyClaimType body UpdateWarrantyClaimTypeRequest true "Datos del tipo de reclamo"
// @Success 200 {object} WarrantyClaimTypeResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-types/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateWarrantyClaimTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	wct := &warranty_claim_type.WarrantyClaimType{
		ID:          id,
		Label:       req.Label,
		LabelPlural: req.LabelPlural,
	}

	if err := h.useCase.Update(r.Context(), wct); err != nil {
		if err == warranty_claim_type.ErrWarrantyClaimTypeNotFound {
			response.Error(w, http.StatusNotFound, "Tipo de reclamo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(wct))
}

// @Summary Eliminar tipo de reclamo de garantía
// @Description Elimina un tipo de reclamo de garantía. No se puede eliminar si tiene reclamos asociados
// @Tags WarrantyClaimTypes
// @Accept json
// @Produce json
// @Param id path int true "ID del tipo de reclamo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claim-types/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == warranty_claim_type.ErrWarrantyClaimTypeNotFound {
			response.Error(w, http.StatusNotFound, "Tipo de reclamo de garantía no encontrado")
			return
		}
		if err == warranty_claim_type.ErrWarrantyClaimTypeInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar tipo de reclamo de garantía")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
