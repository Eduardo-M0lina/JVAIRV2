package invoice_payment

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
)

// Handler maneja las peticiones HTTP para invoice payments
type Handler struct {
	useCase domainPayment.Service
}

// NewHandler crea una nueva instancia del handler de invoice payments
func NewHandler(useCase domainPayment.Service) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registra las rutas del handler como sub-recurso de invoices
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/invoices/{invoiceId}/payments", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreatePaymentRequest representa la solicitud para crear un pago
type CreatePaymentRequest struct {
	PaymentProcessor string  `json:"paymentProcessor"`
	PaymentID        string  `json:"paymentId"`
	Amount           float64 `json:"amount"`
	Notes            string  `json:"notes"`
}

// UpdatePaymentRequest representa la solicitud para actualizar un pago
type UpdatePaymentRequest struct {
	PaymentProcessor *string  `json:"paymentProcessor,omitempty"`
	PaymentID        *string  `json:"paymentId,omitempty"`
	Amount           *float64 `json:"amount,omitempty"`
	Notes            *string  `json:"notes,omitempty"`
}

// PaymentResponse representa la respuesta de un pago
type PaymentResponse struct {
	ID               int64   `json:"id"`
	InvoiceID        int64   `json:"invoiceId"`
	PaymentProcessor string  `json:"paymentProcessor"`
	PaymentID        string  `json:"paymentId"`
	Amount           float64 `json:"amount"`
	Notes            string  `json:"notes"`
	CreatedAt        string  `json:"createdAt,omitempty"`
	UpdatedAt        string  `json:"updatedAt,omitempty"`
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

func toPaymentResponse(p *domainPayment.InvoicePayment) PaymentResponse {
	resp := PaymentResponse{
		ID:               p.ID,
		InvoiceID:        p.InvoiceID,
		PaymentProcessor: p.PaymentProcessor,
		PaymentID:        p.PaymentID,
		Amount:           p.Amount,
		Notes:            p.Notes,
	}

	if p.CreatedAt != nil {
		resp.CreatedAt = p.CreatedAt.Format(timeFormat)
	}
	if p.UpdatedAt != nil {
		resp.UpdatedAt = p.UpdatedAt.Format(timeFormat)
	}

	return resp
}

func parseInvoiceID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "invoiceId"), 10, 64)
}

func parseFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	if search := r.URL.Query().Get("search"); search != "" {
		filters["search"] = search
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		filters["sort"] = sort
	}

	if direction := r.URL.Query().Get("direction"); direction != "" {
		filters["direction"] = direction
	}

	return filters
}
