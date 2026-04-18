package warranty_claim

import (
	"encoding/json"
	"net/http"
	"strconv"

	domainWC "github.com/angumol/jvairv2/pkg/domain/warranty_claim"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

// Handler maneja las peticiones HTTP para warranty claims
type Handler struct {
	useCase domainWC.Service
}

// NewHandler crea una nueva instancia del handler
func NewHandler(useCase domainWC.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/warranty-claims", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

// CreateWarrantyClaimRequest representa la solicitud para crear un reclamo
type CreateWarrantyClaimRequest struct {
	InternalClaimNumber   string  `json:"internalClaimNumber" example:"ICN-2024-001"`
	WarrantyClaimTypeID   int64   `json:"warrantyClaimTypeId" example:"1"`
	WarrantyClaimStatusID int64   `json:"warrantyClaimStatusId" example:"1"`
	JobID                 int64   `json:"jobId" example:"1"`
	InvoiceNumber         *string `json:"invoiceNumber,omitempty" example:"INV-001"`
	WorkDone              *bool   `json:"workDone,omitempty" example:"false"`
	WarrantyPart          *string `json:"warrantyPart,omitempty" example:"Compressor"`
	Manufacturer          *string `json:"manufacturer,omitempty" example:"Carrier"`
	ModelNumber           *string `json:"modelNumber,omitempty" example:"24ACC636A003"`
	PartNumber            *string `json:"partNumber,omitempty" example:"PN-12345"`
	ReplacementPartNumber *string `json:"replacementPartNumber,omitempty" example:"PN-67890"`
	PartDistributor       *string `json:"partDistributor,omitempty" example:"Distributor"`
	PartInvoiceNumber     *string `json:"partInvoiceNumber,omitempty" example:"PI-001"`
	OldPartSerialNumber   *string `json:"oldPartSerialNumber,omitempty" example:"OLD-SN-001"`
	NewPartSerialNumber   *string `json:"newPartSerialNumber,omitempty" example:"NEW-SN-001"`
	EsaNumber             *string `json:"esaNumber,omitempty" example:"ESA-001"`
	Serial                *string `json:"serial,omitempty" example:"SER-001"`
	ClaimNumber           *string `json:"claimNumber,omitempty" example:"CLM-001"`
	Approved              *bool   `json:"approved,omitempty" example:"false"`
	PartsCreditReceived   *bool   `json:"partsCreditReceived,omitempty" example:"false"`
	LaborPaymentReceived  *bool   `json:"laborPaymentReceived,omitempty" example:"false"`
	Notes                 *string `json:"notes,omitempty" example:"Claim notes"`
}

// UpdateWarrantyClaimRequest representa la solicitud para actualizar un reclamo
type UpdateWarrantyClaimRequest struct {
	InternalClaimNumber   string  `json:"internalClaimNumber" example:"ICN-2024-001"`
	WarrantyClaimTypeID   int64   `json:"warrantyClaimTypeId" example:"1"`
	WarrantyClaimStatusID int64   `json:"warrantyClaimStatusId" example:"1"`
	InvoiceNumber         *string `json:"invoiceNumber,omitempty" example:"INV-001"`
	WorkDone              *bool   `json:"workDone,omitempty" example:"false"`
	WarrantyPart          *string `json:"warrantyPart,omitempty" example:"Compressor"`
	Manufacturer          *string `json:"manufacturer,omitempty" example:"Carrier"`
	ModelNumber           *string `json:"modelNumber,omitempty" example:"24ACC636A003"`
	PartNumber            *string `json:"partNumber,omitempty" example:"PN-12345"`
	ReplacementPartNumber *string `json:"replacementPartNumber,omitempty" example:"PN-67890"`
	PartDistributor       *string `json:"partDistributor,omitempty" example:"Distributor"`
	PartInvoiceNumber     *string `json:"partInvoiceNumber,omitempty" example:"PI-001"`
	OldPartSerialNumber   *string `json:"oldPartSerialNumber,omitempty" example:"OLD-SN-001"`
	NewPartSerialNumber   *string `json:"newPartSerialNumber,omitempty" example:"NEW-SN-001"`
	EsaNumber             *string `json:"esaNumber,omitempty" example:"ESA-001"`
	Serial                *string `json:"serial,omitempty" example:"SER-001"`
	ClaimNumber           *string `json:"claimNumber,omitempty" example:"CLM-001"`
	Approved              *bool   `json:"approved,omitempty" example:"false"`
	PartsCreditReceived   *bool   `json:"partsCreditReceived,omitempty" example:"false"`
	LaborPaymentReceived  *bool   `json:"laborPaymentReceived,omitempty" example:"false"`
	Notes                 *string `json:"notes,omitempty" example:"Claim notes"`
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

// WarrantyClaimResponse representa la respuesta de un reclamo
type WarrantyClaimResponse struct {
	ID                    int64        `json:"id" example:"1"`
	InternalClaimNumber   string       `json:"internalClaimNumber" example:"ICN-2024-001"`
	WarrantyClaimTypeID   int64        `json:"warrantyClaimTypeId" example:"1"`
	WarrantyClaimStatusID int64        `json:"warrantyClaimStatusId" example:"1"`
	Job                   *JobResponse `json:"job,omitempty"`
	InvoiceNumber         *string      `json:"invoiceNumber,omitempty"`
	WorkDone              bool         `json:"workDone" example:"false"`
	WarrantyPart          *string      `json:"warrantyPart,omitempty"`
	Manufacturer          *string      `json:"manufacturer,omitempty"`
	ModelNumber           *string      `json:"modelNumber,omitempty"`
	PartNumber            *string      `json:"partNumber,omitempty"`
	ReplacementPartNumber *string      `json:"replacementPartNumber,omitempty"`
	PartDistributor       *string      `json:"partDistributor,omitempty"`
	PartInvoiceNumber     *string      `json:"partInvoiceNumber,omitempty"`
	OldPartSerialNumber   *string      `json:"oldPartSerialNumber,omitempty"`
	NewPartSerialNumber   *string      `json:"newPartSerialNumber,omitempty"`
	EsaNumber             *string      `json:"esaNumber,omitempty"`
	Serial                *string      `json:"serial,omitempty"`
	ClaimNumber           *string      `json:"claimNumber,omitempty"`
	Approved              bool         `json:"approved" example:"false"`
	PartsCreditReceived   bool         `json:"partsCreditReceived" example:"false"`
	LaborPaymentReceived  bool         `json:"laborPaymentReceived" example:"false"`
	Notes                 *string      `json:"notes,omitempty"`
	CreatedAt             string       `json:"createdAt,omitempty"`
	UpdatedAt             string       `json:"updatedAt,omitempty"`
}

func toClaimResponse(wc *domainWC.WarrantyClaim) WarrantyClaimResponse {
	resp := WarrantyClaimResponse{
		ID:                    wc.ID,
		InternalClaimNumber:   wc.InternalClaimNumber,
		WarrantyClaimTypeID:   wc.WarrantyClaimTypeID,
		WarrantyClaimStatusID: wc.WarrantyClaimStatusID,
		InvoiceNumber:         wc.InvoiceNumber,
		WorkDone:              wc.WorkDone,
		WarrantyPart:          wc.WarrantyPart,
		Manufacturer:          wc.Manufacturer,
		ModelNumber:           wc.ModelNumber,
		PartNumber:            wc.PartNumber,
		ReplacementPartNumber: wc.ReplacementPartNumber,
		PartDistributor:       wc.PartDistributor,
		PartInvoiceNumber:     wc.PartInvoiceNumber,
		OldPartSerialNumber:   wc.OldPartSerialNumber,
		NewPartSerialNumber:   wc.NewPartSerialNumber,
		EsaNumber:             wc.EsaNumber,
		Serial:                wc.Serial,
		ClaimNumber:           wc.ClaimNumber,
		Approved:              wc.Approved,
		PartsCreditReceived:   wc.PartsCreditReceived,
		LaborPaymentReceived:  wc.LaborPaymentReceived,
		Notes:                 wc.Notes,
	}

	if wc.CreatedAt != nil {
		resp.CreatedAt = wc.CreatedAt.Format(timeFormat)
	}
	if wc.UpdatedAt != nil {
		resp.UpdatedAt = wc.UpdatedAt.Format(timeFormat)
	}

	if wc.Job != nil {
		jobResp := &JobResponse{
			ID:   wc.Job.ID,
			Week: wc.Job.WeekNumber,
			Property: PropertyResponse{
				ID:      wc.Job.Property.ID,
				Address: wc.Job.Property.Address,
				Customer: CustomerResponse{
					Name: wc.Job.Property.Customer.Name,
				},
			},
		}
		if wc.Job.CompletionDate != nil {
			completionDateStr := wc.Job.CompletionDate.Format(timeFormat)
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
	if typeIDStr := r.URL.Query().Get("warrantyClaimTypeId"); typeIDStr != "" {
		if id, err := strconv.ParseInt(typeIDStr, 10, 64); err == nil {
			filters["warranty_claim_type_id"] = id
		}
	}
	if statusIDStr := r.URL.Query().Get("warrantyClaimStatusId"); statusIDStr != "" {
		if id, err := strconv.ParseInt(statusIDStr, 10, 64); err == nil {
			filters["warranty_claim_status_id"] = id
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

// @Summary Listar reclamos de garantía
// @Description Obtiene una lista paginada de reclamos de garantía con filtros opcionales
// @Tags WarrantyClaims
// @Accept json
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por número de reclamo interno, número de reclamo o notas"
// @Param jobId query int false "Filtrar por ID de trabajo"
// @Param warrantyClaimTypeId query int false "Filtrar por tipo de reclamo"
// @Param warrantyClaimStatusId query int false "Filtrar por estado de reclamo"
// @Param weekNumber query string false "Filtrar por número de semana del trabajo"
// @Param sort query string false "Campo de ordenamiento (internal_claim_number, claim_number, created_at, week_number)"
// @Param direction query string false "Dirección de ordenamiento (asc, desc)"
// @Success 200 {object} response.PaginatedResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claims [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	filters := parseFilters(r)

	claims, total, err := h.useCase.List(r.Context(), filters, page, pageSize)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al listar reclamos de garantía")
		return
	}

	items := make([]WarrantyClaimResponse, len(claims))
	for i, wc := range claims {
		items[i] = toClaimResponse(wc)
	}

	response.Paginated(w, items, page, pageSize, total)
}

// @Summary Crear reclamo de garantía
// @Description Crea un nuevo reclamo de garantía
// @Tags WarrantyClaims
// @Accept json
// @Produce json
// @Param warrantyClaim body CreateWarrantyClaimRequest true "Datos del reclamo"
// @Success 201 {object} WarrantyClaimResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claims [post]
// @Security BearerAuth
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWarrantyClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	claim := &domainWC.WarrantyClaim{
		InternalClaimNumber:   req.InternalClaimNumber,
		WarrantyClaimTypeID:   req.WarrantyClaimTypeID,
		WarrantyClaimStatusID: req.WarrantyClaimStatusID,
		JobID:                 req.JobID,
		InvoiceNumber:         req.InvoiceNumber,
		WarrantyPart:          req.WarrantyPart,
		Manufacturer:          req.Manufacturer,
		ModelNumber:           req.ModelNumber,
		PartNumber:            req.PartNumber,
		ReplacementPartNumber: req.ReplacementPartNumber,
		PartDistributor:       req.PartDistributor,
		PartInvoiceNumber:     req.PartInvoiceNumber,
		OldPartSerialNumber:   req.OldPartSerialNumber,
		NewPartSerialNumber:   req.NewPartSerialNumber,
		EsaNumber:             req.EsaNumber,
		Serial:                req.Serial,
		ClaimNumber:           req.ClaimNumber,
		Notes:                 req.Notes,
	}

	if req.WorkDone != nil {
		claim.WorkDone = *req.WorkDone
	}
	if req.Approved != nil {
		claim.Approved = *req.Approved
	}
	if req.PartsCreditReceived != nil {
		claim.PartsCreditReceived = *req.PartsCreditReceived
	}
	if req.LaborPaymentReceived != nil {
		claim.LaborPaymentReceived = *req.LaborPaymentReceived
	}

	if err := h.useCase.Create(r.Context(), claim); err != nil {
		if err == domainWC.ErrInvalidJob || err == domainWC.ErrInvalidClaimType || err == domainWC.ErrInvalidClaimStatus {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toClaimResponse(claim))
}

// @Summary Obtener reclamo de garantía
// @Description Obtiene un reclamo de garantía por su ID
// @Tags WarrantyClaims
// @Accept json
// @Produce json
// @Param id path int true "ID del reclamo"
// @Success 200 {object} WarrantyClaimResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claims/{id} [get]
// @Security BearerAuth
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	claim, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domainWC.ErrWarrantyClaimNotFound {
			response.Error(w, http.StatusNotFound, "Reclamo de garantía no encontrado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al obtener reclamo de garantía")
		return
	}

	response.JSON(w, http.StatusOK, toClaimResponse(claim))
}

// @Summary Actualizar reclamo de garantía
// @Description Actualiza un reclamo de garantía existente
// @Tags WarrantyClaims
// @Accept json
// @Produce json
// @Param id path int true "ID del reclamo"
// @Param warrantyClaim body UpdateWarrantyClaimRequest true "Datos del reclamo"
// @Success 200 {object} WarrantyClaimResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claims/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req UpdateWarrantyClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Error al decodificar la solicitud")
		return
	}

	claim := &domainWC.WarrantyClaim{
		ID:                    id,
		InternalClaimNumber:   req.InternalClaimNumber,
		WarrantyClaimTypeID:   req.WarrantyClaimTypeID,
		WarrantyClaimStatusID: req.WarrantyClaimStatusID,
		InvoiceNumber:         req.InvoiceNumber,
		WarrantyPart:          req.WarrantyPart,
		Manufacturer:          req.Manufacturer,
		ModelNumber:           req.ModelNumber,
		PartNumber:            req.PartNumber,
		ReplacementPartNumber: req.ReplacementPartNumber,
		PartDistributor:       req.PartDistributor,
		PartInvoiceNumber:     req.PartInvoiceNumber,
		OldPartSerialNumber:   req.OldPartSerialNumber,
		NewPartSerialNumber:   req.NewPartSerialNumber,
		EsaNumber:             req.EsaNumber,
		Serial:                req.Serial,
		ClaimNumber:           req.ClaimNumber,
		Notes:                 req.Notes,
	}

	if req.WorkDone != nil {
		claim.WorkDone = *req.WorkDone
	}
	if req.Approved != nil {
		claim.Approved = *req.Approved
	}
	if req.PartsCreditReceived != nil {
		claim.PartsCreditReceived = *req.PartsCreditReceived
	}
	if req.LaborPaymentReceived != nil {
		claim.LaborPaymentReceived = *req.LaborPaymentReceived
	}

	if err := h.useCase.Update(r.Context(), claim); err != nil {
		if err == domainWC.ErrWarrantyClaimNotFound {
			response.Error(w, http.StatusNotFound, "Reclamo de garantía no encontrado")
			return
		}
		if err == domainWC.ErrWarrantyClaimDeleted {
			response.Error(w, http.StatusGone, "Reclamo de garantía eliminado")
			return
		}
		if err == domainWC.ErrInvalidClaimType || err == domainWC.ErrInvalidClaimStatus {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, toClaimResponse(claim))
}

// @Summary Eliminar reclamo de garantía
// @Description Elimina un reclamo de garantía (soft delete)
// @Tags WarrantyClaims
// @Accept json
// @Produce json
// @Param id path int true "ID del reclamo"
// @Success 204 "No Content"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/warranty-claims/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.useCase.Delete(r.Context(), id); err != nil {
		if err == domainWC.ErrWarrantyClaimNotFound {
			response.Error(w, http.StatusNotFound, "Reclamo de garantía no encontrado")
			return
		}
		if err == domainWC.ErrWarrantyClaimDeleted {
			response.Error(w, http.StatusGone, "Reclamo de garantía ya eliminado")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Error al eliminar reclamo de garantía")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
