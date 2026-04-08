package stripe

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/stripe/stripe-go/v76"
)

// ChargeData contiene los datos relevantes de un charge procesado
type ChargeData struct {
	ChargeID      string
	Amount        float64
	InvoiceNumber string
	CustomerName  string
	CustomerEmail string
	CardLast4     string
}

// ParseWebhookEvent parsea un evento de webhook de Stripe
// Replica el comportamiento de Laravel: StripeEvent::constructFrom($request->all())
func (c *Client) ParseWebhookEvent(r *http.Request) (*stripe.Event, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			slog.Warn("Failed to close request body", slog.String("error", closeErr.Error()))
		}
	}()

	var event stripe.Event
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook event: %w", err)
	}

	return &event, nil
}

// ExtractChargeFromPaymentIntent extrae los datos del charge de un PaymentIntent
// En stripe-go v76, PaymentIntent usa LatestCharge en lugar de Charges
func ExtractChargeFromPaymentIntent(event *stripe.Event) (*ChargeData, error) {
	if event.Type != "payment_intent.succeeded" {
		return nil, nil
	}

	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment intent: %w", err)
	}

	charge := pi.LatestCharge
	if charge == nil {
		// Fallback: extraer invoice_number de metadata del PaymentIntent directamente
		invoiceNumber := ""
		if pi.Metadata != nil {
			invoiceNumber = pi.Metadata["invoice_number"]
		}
		if invoiceNumber == "" {
			return nil, fmt.Errorf("no charge data and no invoice_number in payment intent metadata")
		}

		return &ChargeData{
			ChargeID:      pi.ID,
			Amount:        float64(pi.Amount) / 100.0,
			InvoiceNumber: invoiceNumber,
		}, nil
	}

	invoiceNumber := ""
	if charge.Metadata != nil {
		invoiceNumber = charge.Metadata["invoice_number"]
	}
	// Fallback a metadata del PaymentIntent
	if invoiceNumber == "" && pi.Metadata != nil {
		invoiceNumber = pi.Metadata["invoice_number"]
	}

	if invoiceNumber == "" {
		slog.Warn("Stripe charge missing invoice_number in metadata",
			slog.String("chargeId", charge.ID))
		return nil, fmt.Errorf("missing invoice_number in charge metadata")
	}

	cd := &ChargeData{
		ChargeID:      charge.ID,
		Amount:        float64(charge.Amount) / 100.0,
		InvoiceNumber: invoiceNumber,
	}

	if charge.BillingDetails != nil {
		cd.CustomerName = charge.BillingDetails.Name
		cd.CustomerEmail = charge.BillingDetails.Email
	}

	if charge.PaymentMethodDetails != nil &&
		charge.PaymentMethodDetails.Card != nil {
		cd.CardLast4 = charge.PaymentMethodDetails.Card.Last4
	}

	return cd, nil
}

// FormatPaymentNotes genera las notas del pago en el mismo formato que Laravel
func FormatPaymentNotes(cd *ChargeData) string {
	return fmt.Sprintf("Posted by %s (%s) using card ending in %s",
		cd.CustomerName, cd.CustomerEmail, cd.CardLast4)
}
