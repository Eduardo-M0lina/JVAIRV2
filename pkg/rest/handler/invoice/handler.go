package invoice

import (
	"net/http"
	"strconv"

	domainEmail "github.com/angumol/jvairv2/pkg/domain/email"
	domainInvoice "github.com/angumol/jvairv2/pkg/domain/invoice"
	"github.com/go-chi/chi/v5"
)

// Handler maneja las peticiones HTTP para invoices
type Handler struct {
	useCase      domainInvoice.Service
	emailService domainEmail.Service
}

// NewHandler crea una nueva instancia del handler de invoices
func NewHandler(useCase domainInvoice.Service, emailService domainEmail.Service) *Handler {
	return &Handler{
		useCase:      useCase,
		emailService: emailService,
	}
}

// RegisterRoutes registra las rutas del handler
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/invoices", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
		r.Post("/{id}/email", h.SendEmail)
	})
}

// CreateInvoiceRequest representa la solicitud para crear una factura
type CreateInvoiceRequest struct {
	JobID               int64   `json:"jobId"`
	InvoiceNumber       string  `json:"invoiceNumber"`
	Total               float64 `json:"total"`
	Description         *string `json:"description,omitempty"`
	AllowOnlinePayments *bool   `json:"allowOnlinePayments,omitempty"`
	Notes               *string `json:"notes,omitempty"`
}

// UpdateInvoiceRequest representa la solicitud para actualizar una factura
type UpdateInvoiceRequest struct {
	JobID               *int64   `json:"jobId,omitempty"`
	InvoiceNumber       *string  `json:"invoiceNumber,omitempty"`
	Total               *float64 `json:"total,omitempty"`
	Description         *string  `json:"description,omitempty"`
	AllowOnlinePayments *bool    `json:"allowOnlinePayments,omitempty"`
	Notes               *string  `json:"notes,omitempty"`
}

// InvoiceResponse representa la respuesta de una factura
type InvoiceResponse struct {
	ID                  int64    `json:"id"`
	JobID               int64    `json:"jobId"`
	InvoiceNumber       string   `json:"invoiceNumber"`
	Total               float64  `json:"total"`
	Description         *string  `json:"description,omitempty"`
	AllowOnlinePayments bool     `json:"allowOnlinePayments"`
	Notes               *string  `json:"notes,omitempty"`
	Balance             *float64 `json:"balance,omitempty"`
	CreatedAt           string   `json:"createdAt,omitempty"`
	UpdatedAt           string   `json:"updatedAt,omitempty"`
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

func toInvoiceResponse(inv *domainInvoice.Invoice) InvoiceResponse {
	resp := InvoiceResponse{
		ID:                  inv.ID,
		JobID:               inv.JobID,
		InvoiceNumber:       inv.InvoiceNumber,
		Total:               inv.Total,
		Description:         inv.Description,
		AllowOnlinePayments: inv.AllowOnlinePayments,
		Notes:               inv.Notes,
		Balance:             inv.Balance,
	}

	if inv.CreatedAt != nil {
		resp.CreatedAt = inv.CreatedAt.Format(timeFormat)
	}
	if inv.UpdatedAt != nil {
		resp.UpdatedAt = inv.UpdatedAt.Format(timeFormat)
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

	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters["sort"] = sort
	}

	if direction := r.URL.Query().Get("direction"); direction != "" {
		filters["direction"] = direction
	}

	return filters
}
