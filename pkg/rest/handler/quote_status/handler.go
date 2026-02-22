package quote_status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/jvairv2/pkg/domain/quote_status"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

type Handler struct {
	useCase quote_status.Service
}

func NewHandler(useCase quote_status.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/quote-statuses", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateQuoteStatusRequest representa la solicitud para crear un estado de cotización
type CreateQuoteStatusRequest struct {
	Label string  `json:"label" validate:"required"`
	Class *string `json:"class,omitempty"`
	Order int     `json:"order"`
}

// UpdateQuoteStatusRequest representa la solicitud para actualizar un estado de cotización
type UpdateQuoteStatusRequest struct {
	Label string  `json:"label" validate:"required"`
	Class *string `json:"class,omitempty"`
	Order int     `json:"order"`
}

// QuoteStatusResponse representa la respuesta de un estado de cotización
type QuoteStatusResponse struct {
	ID        int64   `json:"id"`
	Label     string  `json:"label"`
	Class     *string `json:"class,omitempty"`
	Order     int     `json:"order"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

func toResponse(qs *quote_status.QuoteStatus) QuoteStatusResponse {
	resp := QuoteStatusResponse{
		ID:    qs.ID,
		Label: qs.Label,
		Class: qs.Class,
		Order: qs.Order,
	}

	if qs.CreatedAt != nil {
		resp.CreatedAt = qs.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if qs.UpdatedAt != nil {
		resp.UpdatedAt = qs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return resp
}

// List maneja la solicitud de listado de estados de cotización
// @Summary Listar estados de cotización
// @Description Obtiene una lista paginada de estados de cotización
// @Tags QuoteStatuses
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por label"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quote-statuses [get]
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
		response.Error(w, http.StatusInternalServerError, "Error al listar estados de cotización")
		return
	}

	items := make([]QuoteStatusResponse, len(statuses))
	for i, qs := range statuses {
		items[i] = toResponse(qs)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// Create maneja la solicitud de creación de un estado de cotización
// @Summary Crear estado de cotización
// @Description Crea un nuevo estado de cotización
// @Tags QuoteStatuses
// @Accept json
// @Produce json
// @Param quoteStatus body CreateQuoteStatusRequest true "Datos del estado de cotización"
// @Success 201 {object} QuoteStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quote-statuses [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateQuoteStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	qs := &quote_status.QuoteStatus{
		Label: req.Label,
		Class: req.Class,
		Order: req.Order,
	}

	if err := h.useCase.Create(r.Context(), qs); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toResponse(qs))
}

// Get maneja la solicitud de obtención de un estado de cotización por ID
// @Summary Obtener estado de cotización
// @Description Obtiene un estado de cotización por su ID
// @Tags QuoteStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de cotización"
// @Success 200 {object} QuoteStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quote-statuses/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	qs, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == quote_status.ErrQuoteStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de cotización no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener estado de cotización")
		return
	}

	response.JSON(w, http.StatusOK, toResponse(qs))
}

// Update maneja la solicitud de actualización de un estado de cotización
// @Summary Actualizar estado de cotización
// @Description Actualiza un estado de cotización existente
// @Tags QuoteStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de cotización"
// @Param quoteStatus body UpdateQuoteStatusRequest true "Datos del estado de cotización"
// @Success 200 {object} QuoteStatusResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quote-statuses/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateQuoteStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	qs := &quote_status.QuoteStatus{
		ID:    id,
		Label: req.Label,
		Class: req.Class,
		Order: req.Order,
	}

	if err := h.useCase.Update(r.Context(), qs); err != nil {
		if err == quote_status.ErrQuoteStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de cotización no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toResponse(qs))
}

// Delete maneja la solicitud de eliminación de un estado de cotización
// @Summary Eliminar estado de cotización
// @Description Elimina un estado de cotización. No se puede eliminar si tiene cotizaciones asociadas
// @Tags QuoteStatuses
// @Accept json
// @Produce json
// @Param id path int true "ID del estado de cotización"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quote-statuses/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == quote_status.ErrQuoteStatusNotFound {
			response.Error(w, http.StatusNotFound, "Estado de cotización no encontrado")
			return
		}
		if err == quote_status.ErrQuoteStatusInUse {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar estado de cotización")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
