package stripe

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

// PaymentIntentResult contiene el resultado de crear un PaymentIntent
type PaymentIntentResult struct {
	ID            string `json:"id"`
	ClientSecret  string `json:"clientSecret"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	InvoiceNumber string `json:"invoiceNumber"`
	PublicKey     string `json:"publicKey"`
}

// CreatePaymentIntent crea un PaymentIntent en Stripe para una factura
func (c *Client) CreatePaymentIntent(ctx context.Context, invoiceNumber string, balanceUSD float64) (*PaymentIntentResult, error) {
	// Convertir a centavos (Stripe usa la unidad más pequeña de la moneda)
	amountCents := int64(math.Round(balanceUSD * 100))

	if amountCents <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
	}
	params.AddMetadata("invoice_number", invoiceNumber)
	params.Context = ctx

	pi, err := paymentintent.New(params)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create Stripe PaymentIntent",
			slog.String("invoiceNumber", invoiceNumber),
			slog.Float64("amount", balanceUSD),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	slog.InfoContext(ctx, "Stripe PaymentIntent created",
		slog.String("paymentIntentId", pi.ID),
		slog.String("invoiceNumber", invoiceNumber),
		slog.Int64("amountCents", amountCents))

	return &PaymentIntentResult{
		ID:            pi.ID,
		ClientSecret:  pi.ClientSecret,
		Amount:        amountCents,
		Currency:      string(pi.Currency),
		InvoiceNumber: invoiceNumber,
		PublicKey:     c.publicKey,
	}, nil
}
