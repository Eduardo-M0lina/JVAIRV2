package invoice_payment

import (
	"log/slog"
	"net/http"
	"strconv"

	domainPayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
	"github.com/your-org/jvairv2/pkg/rest/response"
)

// List maneja la solicitud de listado de pagos de una factura
// @Summary Listar pagos de factura
// @Description Obtiene una lista paginada de pagos de una factura
// @Tags Invoice Payments
// @Accept json
// @Produce json
// @Param invoiceId path int true "ID de la factura"
// @Param page query int false "Número de página" default(1)
// @Param pageSize query int false "Tamaño de página" default(10)
// @Param search query string false "Búsqueda por payment_id"
// @Param sort query string false "Campo de ordenamiento: amount, created_at, payment_processor"
// @Param direction query string false "Dirección de ordenamiento: asc, desc" default(desc)
// @Success 200 {object} response.PaginatedResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/invoices/{invoiceId}/payments [get]
// @Security BearerAuth
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseInvoiceID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de factura inválido")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	filters := parseFilters(r)

	payments, total, err := h.useCase.ListByInvoiceID(r.Context(), invoiceID, filters, page, pageSize)
	if err != nil {
		if err == domainPayment.ErrInvalidInvoice {
			response.Error(w, http.StatusNotFound, "Factura no encontrada")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to list invoice payments",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al listar pagos")
		return
	}

	items := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		items[i] = toPaymentResponse(p)
	}

	totalPages := (total + pageSize - 1) / pageSize

	response.JSON(w, http.StatusOK, response.PaginatedResponse{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	})
}
