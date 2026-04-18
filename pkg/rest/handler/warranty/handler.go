package warranty

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	domainWarranty "github.com/angumol/jvairv2/pkg/domain/warranty"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

// Handler maneja las peticiones HTTP para warranties
type Handler struct {
	useCase domainWarranty.Service
}

// NewHandler crea una nueva instancia del handler de warranties
func NewHandler(useCase domainWarranty.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/warranties", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

const timeFormat = "2006-01-02T15:04:05Z07:00"
const dateFormat = "01-02-2006"

// CreateWarrantyRequest representa la solicitud para crear una garantía
type CreateWarrantyRequest struct {
	WarrantyNumber   string  `json:"warrantyNumber" example:"WRN-2024-001"`
	JobID            int64   `json:"jobId" example:"1"`
	WarrantyTypeID   int64   `json:"warrantyTypeId" example:"1"`
	WarrantyStatusID int64   `json:"warrantyStatusId" example:"1"`
	DateSubmitted    *string `json:"dateSubmitted,omitempty" example:"01-15-2024"`
	AgreementNumber  *string `json:"agreementNumber,omitempty" example:"AGR-001"`
	AuditDone        *bool   `json:"auditDone,omitempty" example:"false"`
	Notes            *string `json:"notes,omitempty" example:"Warranty notes"`
}

// UpdateWarrantyRequest representa la solicitud para actualizar una garantía
type UpdateWarrantyRequest struct {
	WarrantyNumber   string  `json:"warrantyNumber" example:"WRN-2024-001"`
	WarrantyTypeID   int64   `json:"warrantyTypeId" example:"1"`
	WarrantyStatusID int64   `json:"warrantyStatusId" example:"1"`
	DateSubmitted    *string `json:"dateSubmitted,omitempty" example:"01-15-2024"`
	AgreementNumber  *string `json:"agreementNumber,omitempty" example:"AGR-001"`
	AuditDone        *bool   `json:"auditDone,omitempty" example:"false"`
	Notes            *string `json:"notes,omitempty" example:"Warranty notes"`
}

// CustomerResponse representa la información del cliente en la respuesta
type CustomerResponse struct {
	Name string `json:"name" example:"John Doe"`
}

// PropertyResponse representa la información de la propiedad en la respuesta
type PropertyResponse struct {
	ID       int64            `json:"id" example:"1"`
	Address  string           `json:"address" example:"123 Main St, Atlanta, GA 30301"`
	Customer CustomerResponse `json:"customer"`
}

// JobResponse representa la información del trabajo en la respuesta
type JobResponse struct {
	ID             int64            `json:"id" example:"1"`
	Week           *int             `json:"week,omitempty" example:"42"`
	CompletionDate *string          `json:"completionDate,omitempty" example:"2024-01-15T10:30:00Z"`
	Property       PropertyResponse `json:"property"`
}

// WarrantyResponse representa la respuesta de una garantía
type WarrantyResponse struct {
	ID               int64        `json:"id" example:"1"`
	WarrantyNumber   string       `json:"warrantyNumber" example:"WRN-2024-001"`
	Job              *JobResponse `json:"job,omitempty"`
	WarrantyTypeID   int64        `json:"warrantyTypeId" example:"1"`
	WarrantyStatusID int64        `json:"warrantyStatusId" example:"1"`
	DateSubmitted    *string      `json:"dateSubmitted,omitempty" example:"01-15-2024"`
	AgreementNumber  *string      `json:"agreementNumber,omitempty" example:"AGR-001"`
	AuditDone        bool         `json:"auditDone" example:"false"`
	Notes            *string      `json:"notes,omitempty" example:"Warranty notes"`
	CreatedAt        string       `json:"createdAt,omitempty" example:"2024-01-15T10:30:00Z"`
	UpdatedAt        string       `json:"updatedAt,omitempty" example:"2024-01-15T10:30:00Z"`
}

func toWarrantyResponse(w *domainWarranty.Warranty) WarrantyResponse {
	resp := WarrantyResponse{
		ID:               w.ID,
		WarrantyNumber:   w.WarrantyNumber,
		WarrantyTypeID:   w.WarrantyTypeID,
		WarrantyStatusID: w.WarrantyStatusID,
		AuditDone:        w.AuditDone,
		AgreementNumber:  w.AgreementNumber,
		Notes:            w.Notes,
	}

	if w.DateSubmitted != nil {
		ds := w.DateSubmitted.Format(dateFormat)
		resp.DateSubmitted = &ds
	}
	if w.CreatedAt != nil {
		resp.CreatedAt = w.CreatedAt.Format(timeFormat)
	}
	if w.UpdatedAt != nil {
		resp.UpdatedAt = w.UpdatedAt.Format(timeFormat)
	}

	if w.Job != nil {
		jobResp := &JobResponse{
			ID:   w.Job.ID,
			Week: w.Job.WeekNumber,
			Property: PropertyResponse{
				ID:      w.Job.Property.ID,
				Address: w.Job.Property.Address,
				Customer: CustomerResponse{
					Name: w.Job.Property.Customer.Name,
				},
			},
		}
		if w.Job.CompletionDate != nil {
			completionDateStr := w.Job.CompletionDate.Format(timeFormat)
			jobResp.CompletionDate = &completionDateStr
		}
		resp.Job = jobResp
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
	if typeIDStr := r.URL.Query().Get("warrantyTypeId"); typeIDStr != "" {
		if id, err := strconv.ParseInt(typeIDStr, 10, 64); err == nil {
			filters["warranty_type_id"] = id
		}
	}
	if statusIDStr := r.URL.Query().Get("warrantyStatusId"); statusIDStr != "" {
		if id, err := strconv.ParseInt(statusIDStr, 10, 64); err == nil {
			filters["warranty_status_id"] = id
		}
	}
	if weekNumber := r.URL.Query().Get("weekNumber"); weekNumber != "" {
		filters["week_number"] = weekNumber
	}
	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters["sort"] = sort
	}
	if direction := r.URL.Query().Get("direction"); direction != "" {
		filters["direction"] = direction
	}

	return filters
}

// @Summary Listar garantías
// @Description Obtiene una lista paginada de garantías con filtros opcionales
// @Tags Warranties
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por número de garantía, número de acuerdo o notas"
// @Param jobId query int false "Filtrar por ID de trabajo"
// @Param warrantyTypeId query int false "Filtrar por tipo de garantía"
// @Param warrantyStatusId query int false "Filtrar por estado de garantía"
// @Param weekNumber query string false "Filtrar por número de semana del trabajo"
// @Param sort query string false "Campo de ordenamiento (warranty_number, date_submitted, created_at, week_number)"
// @Param direction query string false "Dirección de ordenamiento (asc, desc)"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	warranties, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar garantías")
		return
	}

	items := make([]WarrantyResponse, len(warranties))
	for i, wr := range warranties {
		items[i] = toWarrantyResponse(wr)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// @Summary Crear garantía
// @Description Crea una nueva garantía
// @Tags Warranties
// @Accept json
// @Produce json
// @Param warranty body CreateWarrantyRequest true "Datos de la garantía"
// @Success 201 {object} WarrantyResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWarrantyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	warranty := &domainWarranty.Warranty{
		WarrantyNumber:   req.WarrantyNumber,
		JobID:            req.JobID,
		WarrantyTypeID:   req.WarrantyTypeID,
		WarrantyStatusID: req.WarrantyStatusID,
		AgreementNumber:  req.AgreementNumber,
		Notes:            req.Notes,
	}

	if req.AuditDone != nil {
		warranty.AuditDone = *req.AuditDone
	}

	if req.DateSubmitted != nil && *req.DateSubmitted != "" {
		t, err := time.Parse(dateFormat, *req.DateSubmitted)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Formato de fecha inválido. Use MM-DD-YYYY")
			return
		}
		warranty.DateSubmitted = &t
	}

	if err := h.useCase.Create(r.Context(), warranty); err != nil {
		if err == domainWarranty.ErrInvalidJob {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if err == domainWarranty.ErrInvalidWarrantyType || err == domainWarranty.ErrInvalidWarrantyStatus {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toWarrantyResponse(warranty))
}

// @Summary Obtener garantía
// @Description Obtiene una garantía por su ID
// @Tags Warranties
// @Accept json
// @Produce json
// @Param id path int true "ID de la garantía"
// @Success 200 {object} WarrantyResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	warranty, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domainWarranty.ErrWarrantyNotFound {
			response.Error(w, http.StatusNotFound, "Garantía no encontrada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener garantía")
		return
	}

	response.JSON(w, http.StatusOK, toWarrantyResponse(warranty))
}

// @Summary Actualizar garantía
// @Description Actualiza una garantía existente
// @Tags Warranties
// @Accept json
// @Produce json
// @Param id path int true "ID de la garantía"
// @Param warranty body UpdateWarrantyRequest true "Datos de la garantía"
// @Success 200 {object} WarrantyResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateWarrantyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	warranty := &domainWarranty.Warranty{
		ID:               id,
		WarrantyNumber:   req.WarrantyNumber,
		WarrantyTypeID:   req.WarrantyTypeID,
		WarrantyStatusID: req.WarrantyStatusID,
		AgreementNumber:  req.AgreementNumber,
		Notes:            req.Notes,
	}

	if req.AuditDone != nil {
		warranty.AuditDone = *req.AuditDone
	}

	if req.DateSubmitted != nil && *req.DateSubmitted != "" {
		t, err := time.Parse(dateFormat, *req.DateSubmitted)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Formato de fecha inválido. Use MM-DD-YYYY")
			return
		}
		warranty.DateSubmitted = &t
	}

	if err := h.useCase.Update(r.Context(), warranty); err != nil {
		if err == domainWarranty.ErrWarrantyNotFound {
			response.Error(w, http.StatusNotFound, "Garantía no encontrada")
			return
		}
		if err == domainWarranty.ErrWarrantyDeleted {
			response.Error(w, http.StatusGone, "Garantía eliminada")
			return
		}
		if err == domainWarranty.ErrInvalidWarrantyType || err == domainWarranty.ErrInvalidWarrantyStatus {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toWarrantyResponse(warranty))
}

// @Summary Eliminar garantía
// @Description Elimina una garantía (soft delete)
// @Tags Warranties
// @Accept json
// @Produce json
// @Param id path int true "ID de la garantía"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranties/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == domainWarranty.ErrWarrantyNotFound {
			response.Error(w, http.StatusNotFound, "Garantía no encontrada")
			return
		}
		if err == domainWarranty.ErrWarrantyDeleted {
			response.Error(w, http.StatusGone, "Garantía ya eliminada")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar garantía")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
