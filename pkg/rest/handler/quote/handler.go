package quote

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// Handler maneja las peticiones HTTP para cotizaciones
type Handler struct {
	useCase domainQuote.Service
}

// NewHandler crea una nueva instancia del handler de cotizaciones
func NewHandler(useCase domainQuote.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/quotes", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateQuoteRequest representa la solicitud para crear una cotización
type CreateQuoteRequest struct {
	JobID         int64   `json:"jobId"`
	QuoteNumber   string  `json:"quoteNumber"`
	QuoteStatusID int64   `json:"quoteStatusId"`
	Amount        float64 `json:"amount"`
	Description   *string `json:"description,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// UpdateQuoteRequest representa la solicitud para actualizar una cotización
type UpdateQuoteRequest struct {
	JobID         int64   `json:"jobId"`
	QuoteNumber   string  `json:"quoteNumber"`
	QuoteStatusID int64   `json:"quoteStatusId"`
	Amount        float64 `json:"amount"`
	Description   *string `json:"description,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// QuoteResponse representa la respuesta de una cotización
type QuoteResponse struct {
	ID            int64   `json:"id"`
	JobID         int64   `json:"jobId"`
	QuoteNumber   string  `json:"quoteNumber"`
	QuoteStatusID int64   `json:"quoteStatusId"`
	Amount        float64 `json:"amount"`
	Description   *string `json:"description,omitempty"`
	Notes         *string `json:"notes,omitempty"`
	CreatedAt     string  `json:"createdAt,omitempty"`
	UpdatedAt     string  `json:"updatedAt,omitempty"`
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

func toQuoteResponse(q *domainQuote.Quote) QuoteResponse {
	resp := QuoteResponse{
		ID:            q.ID,
		JobID:         q.JobID,
		QuoteNumber:   q.QuoteNumber,
		QuoteStatusID: q.QuoteStatusID,
		Amount:        q.Amount,
		Description:   q.Description,
		Notes:         q.Notes,
	}

	if q.CreatedAt != nil {
		resp.CreatedAt = q.CreatedAt.Format(timeFormat)
	}
	if q.UpdatedAt != nil {
		resp.UpdatedAt = q.UpdatedAt.Format(timeFormat)
	}

	return resp
}

func parseFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	if jobIDStr := r.URL.Query().Get("jobId"); jobIDStr != "" {
		if id, err := strconv.ParseInt(jobIDStr, 10, 64); err == nil {
			filters["job_id"] = id
		}
	}

	if quoteStatusIDStr := r.URL.Query().Get("quoteStatusId"); quoteStatusIDStr != "" {
		if id, err := strconv.ParseInt(quoteStatusIDStr, 10, 64); err == nil {
			filters["quote_status_id"] = id
		}
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters["sort"] = sort
	}

	if direction := r.URL.Query().Get("direction"); direction != "" {
		filters["direction"] = direction
	}

	return filters
}

// List maneja la solicitud de listado de cotizaciones
// @Summary Listar cotizaciones
// @Description Obtiene una lista paginada de cotizaciones con filtros opcionales
// @Tags Quotes
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(15)
// @Param search query string false "Búsqueda por número de cotización"
// @Param jobId query int false "Filtrar por trabajo"
// @Param quoteStatusId query int false "Filtrar por estado de cotización"
// @Param sort query string false "Campo de ordenamiento (quote_number, amount, created_at)"
// @Param direction query string false "Dirección de ordenamiento (ASC, DESC)"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quotes [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	quotes, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar cotizaciones")
		return
	}

	items := make([]QuoteResponse, len(quotes))
	for i, q := range quotes {
		items[i] = toQuoteResponse(q)
	}

	response.Paginated(w, items, page, pageSize, int(total))
}

// Create maneja la solicitud de creación de una cotización
// @Summary Crear cotización
// @Description Crea una nueva cotización para un trabajo
// @Tags Quotes
// @Accept json
// @Produce json
// @Param quote body CreateQuoteRequest true "Datos de la cotización"
// @Success 201 {object} QuoteResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quotes [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	q := &domainQuote.Quote{
		JobID:         req.JobID,
		QuoteNumber:   req.QuoteNumber,
		QuoteStatusID: req.QuoteStatusID,
		Amount:        req.Amount,
		Description:   req.Description,
		Notes:         req.Notes,
	}

	if err := h.useCase.Create(r.Context(), q); err != nil {
		switch err {
		case domainQuote.ErrInvalidJob,
			domainQuote.ErrInvalidQuoteStatus:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			if err.Error() == "quote_number is required" ||
				err.Error() == "job_id is required" ||
				err.Error() == "quote_status_id is required" {
				response.Error(w, http.StatusBadRequest, err.Error())
			} else {
				response.Error(w, http.StatusInternalServerError, "Error al crear cotización")
			}
		}
		return
	}

	response.JSON(w, http.StatusCreated, toQuoteResponse(q))
}

// Get maneja la solicitud de obtención de una cotización por ID
// @Summary Obtener cotización
// @Description Obtiene una cotización por su ID
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path int true "ID de la cotización"
// @Success 200 {object} QuoteResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quotes/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	q, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domainQuote.ErrQuoteNotFound {
			response.Error(w, http.StatusNotFound, "Cotización no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener cotización")
		return
	}

	response.JSON(w, http.StatusOK, toQuoteResponse(q))
}

// Update maneja la solicitud de actualización de una cotización
// @Summary Actualizar cotización
// @Description Actualiza una cotización existente
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path int true "ID de la cotización"
// @Param quote body UpdateQuoteRequest true "Datos de la cotización"
// @Success 200 {object} QuoteResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quotes/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateQuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	q := &domainQuote.Quote{
		ID:            id,
		JobID:         req.JobID,
		QuoteNumber:   req.QuoteNumber,
		QuoteStatusID: req.QuoteStatusID,
		Amount:        req.Amount,
		Description:   req.Description,
		Notes:         req.Notes,
	}

	if err := h.useCase.Update(r.Context(), q); err != nil {
		switch err {
		case domainQuote.ErrQuoteNotFound:
			response.Error(w, http.StatusNotFound, "Cotización no encontrada")
		case domainQuote.ErrInvalidJob,
			domainQuote.ErrInvalidQuoteStatus:
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			if err.Error() == "quote_number is required" ||
				err.Error() == "job_id is required" ||
				err.Error() == "quote_status_id is required" {
				response.Error(w, http.StatusBadRequest, err.Error())
			} else {
				response.Error(w, http.StatusInternalServerError, "Error al actualizar cotización")
			}
		}
		return
	}

	response.JSON(w, http.StatusOK, toQuoteResponse(q))
}

// Delete maneja la solicitud de eliminación de una cotización
// @Summary Eliminar cotización
// @Description Elimina una cotización (soft delete)
// @Tags Quotes
// @Accept json
// @Produce json
// @Param id path int true "ID de la cotización"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/quotes/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == domainQuote.ErrQuoteNotFound {
			response.Error(w, http.StatusNotFound, "Cotización no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar cotización")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
