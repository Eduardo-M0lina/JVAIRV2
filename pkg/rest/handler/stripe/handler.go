package stripe

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	domainInvoice "github.com/angumol/jvairv2/pkg/domain/invoice"
	domainPayment "github.com/angumol/jvairv2/pkg/domain/invoice_payment"
	infraStripe "github.com/angumol/jvairv2/pkg/infrastructure/stripe"
	"github.com/angumol/jvairv2/pkg/rest/response"
	"github.com/go-chi/chi/v5"
)

// Handler maneja los endpoints públicos de Stripe
type Handler struct {
	stripeClient *infraStripe.Client
	invoiceRepo  domainInvoice.Repository
	paymentUC    domainPayment.Service
}

// NewHandler crea una nueva instancia del handler de Stripe
func NewHandler(
	stripeClient *infraStripe.Client,
	invoiceRepo domainInvoice.Repository,
	paymentUC domainPayment.Service,
) *Handler {
	return &Handler{
		stripeClient: stripeClient,
		invoiceRepo:  invoiceRepo,
		paymentUC:    paymentUC,
	}
}

// CreatePaymentIntent crea un PaymentIntent para una factura
// @Summary Crear PaymentIntent para pago de factura
// @Description Crea un PaymentIntent en Stripe para el balance pendiente de una factura. Endpoint público (sin autenticación).
// @Tags Stripe
// @Produce json
// @Param id path int true "ID de la factura"
// @Success 200 {object} stripe.PaymentIntentResult
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /invoices/{id}/payment-intent [get]
func (h *Handler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de factura inválido")
		return
	}

	// Obtener factura con balance
	invoice, err := h.invoiceRepo.GetByID(r.Context(), invoiceID)
	if err != nil {
		if err == domainInvoice.ErrInvoiceNotFound {
			response.Error(w, http.StatusNotFound, "Factura no encontrada")
			return
		}
		slog.ErrorContext(r.Context(), "Failed to get invoice for payment intent",
			slog.Int64("invoiceId", invoiceID),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al obtener factura")
		return
	}

	// Verificar que permite pagos online
	if !invoice.AllowOnlinePayments {
		response.Error(w, http.StatusBadRequest, "Online payments not allowed for this invoice")
		return
	}

	// Verificar que no está pagada
	if invoice.IsPaid() {
		response.Error(w, http.StatusBadRequest, "Invoice is already paid")
		return
	}

	// Obtener balance
	balance := 0.0
	if invoice.Balance != nil {
		balance = *invoice.Balance
	}
	balance = math.Round(balance*100) / 100

	if balance <= 0 {
		response.Error(w, http.StatusBadRequest, "Invoice has no outstanding balance")
		return
	}

	// Crear PaymentIntent en Stripe
	result, err := h.stripeClient.CreatePaymentIntent(r.Context(), invoice.InvoiceNumber, balance)
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to create payment intent",
			slog.Int64("invoiceId", invoiceID),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Error al crear payment intent")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// Webhook procesa eventos de webhook de Stripe
// @Summary Procesar webhook de Stripe
// @Description Recibe y procesa eventos de Stripe (payment_intent.succeeded). Endpoint público (sin autenticación).
// @Tags Stripe
// @Accept json
// @Produce json
// @Success 200 {string} string "Event processed"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /webhooks/stripe [post]
func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	// Parsear y verificar el evento
	event, err := h.stripeClient.ParseWebhookEvent(r)
	if err != nil {
		slog.Error("Failed to parse Stripe webhook event",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Invalid webhook event")
		return
	}

	slog.Info("Stripe webhook received", slog.String("eventType", string(event.Type)))

	// Solo procesar payment_intent.succeeded
	if event.Type != "payment_intent.succeeded" {
		response.JSON(w, http.StatusOK, map[string]string{"message": "Event type not handled"})
		return
	}

	// Extraer datos del charge
	chargeData, err := infraStripe.ExtractChargeFromPaymentIntent(event)
	if err != nil {
		slog.Error("Failed to extract charge data from payment intent",
			slog.String("error", err.Error()))
		response.Error(w, http.StatusBadRequest, "Failed to extract payment data")
		return
	}

	if chargeData == nil {
		response.JSON(w, http.StatusOK, map[string]string{"message": "No charge data to process"})
		return
	}

	// Buscar factura por invoice_number
	invoice, err := h.invoiceRepo.GetByInvoiceNumber(r.Context(), chargeData.InvoiceNumber)
	if err != nil {
		slog.Error("Invoice not found for Stripe webhook",
			slog.String("invoiceNumber", chargeData.InvoiceNumber),
			slog.String("chargeId", chargeData.ChargeID),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Invoice not found")
		return
	}

	// Crear registro de pago
	payment := &domainPayment.InvoicePayment{
		InvoiceID:        invoice.ID,
		Amount:           chargeData.Amount,
		PaymentID:        chargeData.ChargeID,
		PaymentProcessor: "Stripe",
		Notes:            infraStripe.FormatPaymentNotes(chargeData),
	}

	if err := h.paymentUC.Create(r.Context(), payment); err != nil {
		slog.Error("Failed to create invoice payment from Stripe webhook",
			slog.String("chargeId", chargeData.ChargeID),
			slog.String("invoiceNumber", chargeData.InvoiceNumber),
			slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Failed to record payment")
		return
	}

	slog.Info("Stripe payment recorded successfully",
		slog.Int64("invoiceId", invoice.ID),
		slog.String("invoiceNumber", chargeData.InvoiceNumber),
		slog.String("chargeId", chargeData.ChargeID),
		slog.Float64("amount", chargeData.Amount))

	response.JSON(w, http.StatusOK, map[string]string{"message": "Payment recorded successfully"})
}
