# Integración con Stripe - Pagos Online de Facturas

> **Estado**: ✅ **COMPLETAMENTE MIGRADO**
> **Fecha de migración**: Fase 2
> **Versión SDK**: `stripe-go/v76`

---

## Resumen Ejecutivo

La integración con Stripe permite a los clientes pagar facturas online mediante tarjeta de crédito/débito. El sistema crea un **PaymentIntent** en Stripe y procesa los pagos mediante webhooks, registrando automáticamente los pagos en la tabla `invoice_payments`.

### Diferencias clave Laravel vs Go

| Aspecto | Laravel | Go (JVAIRV2) |
|---------|---------|--------------|
| **SDK** | `stripe/stripe-php` | `stripe/stripe-go/v76` |
| **Endpoint de pago** | `GET /invoices/{id}/pay` (Blade view) | `GET /api/v1/invoices/{id}/payment-intent` (JSON API) |
| **Webhook** | `POST /stripe/webhook` | `POST /webhooks/stripe` |
| **Autenticación** | Pública (sin auth) | Pública (sin auth) |
| **Frontend** | Blade + jQuery + Stripe.js | API REST (frontend separado) |
| **Verificación de firma** | ❌ No implementada | ❌ No implementada (pendiente) |

---

## Arquitectura de la Solución

### Flujo de Pago Completo

```
┌─────────────┐
│   Cliente   │
└──────┬──────┘
       │ 1. Solicita pagar factura
       ▼
┌─────────────────────────────────────┐
│  GET /api/v1/invoices/{id}/         │
│      payment-intent                 │
│                                     │
│  - Valida factura existe            │
│  - Verifica allow_online_payments   │
│  - Verifica no está pagada          │
│  - Calcula balance pendiente        │
└──────┬──────────────────────────────┘
       │ 2. Crea PaymentIntent
       ▼
┌─────────────────────────────────────┐
│         Stripe API                  │
│  stripe.PaymentIntent.Create()      │
│                                     │
│  Metadata:                          │
│    - invoice_number                 │
└──────┬──────────────────────────────┘
       │ 3. Retorna client_secret
       ▼
┌─────────────────────────────────────┐
│   Frontend (Stripe.js)              │
│  stripe.confirmCardPayment()        │
│                                     │
│  - Captura datos de tarjeta         │
│  - Envía a Stripe                   │
└──────┬──────────────────────────────┘
       │ 4. Pago procesado
       ▼
┌─────────────────────────────────────┐
│  Stripe envía webhook               │
│  POST /webhooks/stripe              │
│                                     │
│  Event: payment_intent.succeeded    │
└──────┬──────────────────────────────┘
       │ 5. Procesa webhook
       ▼
┌─────────────────────────────────────┐
│  Backend Go                         │
│  - Extrae invoice_number            │
│  - Busca factura                    │
│  - Crea InvoicePayment              │
│  - Registra en BD                   │
└─────────────────────────────────────┘
```

---

## Implementación en Go

### 1. Estructura de Archivos

```
pkg/infrastructure/stripe/
├── client.go           # Cliente principal de Stripe
├── payment_intent.go   # Creación de PaymentIntents
└── webhook.go          # Procesamiento de webhooks

pkg/rest/handler/stripe/
└── handler.go          # Endpoints HTTP públicos

pkg/domain/invoice_payment/
├── entity.go           # Entidad InvoicePayment
├── repository.go       # Interfaz del repositorio
├── usecase.go          # Lógica de negocio
└── create.go           # Caso de uso de creación
```

### 2. Configuración

#### Variables de Entorno (`configs/app.env`)

```env
# === Stripe ===
STRIPE_SECRET_KEY=sk_test_***
STRIPE_PUBLIC_KEY=pk_test_***
STRIPE_WEBHOOK_SECRET=whsec_***  # Opcional (no implementado aún)
```

#### Struct de Configuración (`configs/config.go`)

```go
type StripeConfig struct {
    SecretKey string
    PublicKey string
}
```

#### Configuración del Webhook en Stripe Dashboard

**⚠️ CRÍTICO**: Stripe necesita saber la URL de tu servidor para enviar webhooks.

**Pasos para configurar**:

1. **Ir a Stripe Dashboard**:
   - Test mode: https://dashboard.stripe.com/test/webhooks
   - Live mode: https://dashboard.stripe.com/webhooks

2. **Agregar endpoint**:
   - Click en "Add endpoint"
   - **Endpoint URL**: `https://tu-dominio.com/webhooks/stripe`
   - **Description**: "JVAIR Invoice Payments"
   - **Events to send**: Seleccionar `payment_intent.succeeded`
   - Click "Add endpoint"

3. **Obtener Signing Secret**:
   - Click en el endpoint recién creado
   - Copiar el "Signing secret" (whsec_...)
   - Agregar a `configs/app.env` como `STRIPE_WEBHOOK_SECRET`

**Ejemplo de configuración**:
```
Endpoint URL: https://app.jvair.com/webhooks/stripe
Events: payment_intent.succeeded
Status: Enabled ✅
```

**Para desarrollo local**:
- Usar Stripe CLI: `stripe listen --forward-to localhost:8090/webhooks/stripe`
- Esto crea un webhook temporal automáticamente

### 3. Cliente de Stripe

**Archivo**: `pkg/infrastructure/stripe/client.go`

```go
type Client struct {
    secretKey string
    publicKey string
}

func NewClient(cfg *configs.StripeConfig) (*Client, error) {
    if cfg.SecretKey == "" {
        return nil, fmt.Errorf("STRIPE_SECRET_KEY is required")
    }

    // Configurar la API key global del SDK
    stripe.Key = cfg.SecretKey

    return &Client{
        secretKey: cfg.SecretKey,
        publicKey: cfg.PublicKey,
    }, nil
}
```

### 4. Creación de PaymentIntent

**Archivo**: `pkg/infrastructure/stripe/payment_intent.go`

```go
type PaymentIntentResult struct {
    ID            string  `json:"id"`
    ClientSecret  string  `json:"clientSecret"`
    Amount        int64   `json:"amount"`
    Currency      string  `json:"currency"`
    InvoiceNumber string  `json:"invoiceNumber"`
    PublicKey     string  `json:"publicKey"`
}

func (c *Client) CreatePaymentIntent(ctx context.Context, invoiceNumber string, balanceUSD float64) (*PaymentIntentResult, error) {
    // Convertir a centavos (Stripe usa la unidad más pequeña)
    amountCents := int64(math.Round(balanceUSD * 100))

    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(amountCents),
        Currency: stripe.String(string(stripe.CurrencyUSD)),
    }
    params.AddMetadata("invoice_number", invoiceNumber)
    params.Context = ctx

    pi, err := paymentintent.New(params)
    if err != nil {
        return nil, fmt.Errorf("failed to create payment intent: %w", err)
    }

    return &PaymentIntentResult{
        ID:            pi.ID,
        ClientSecret:  pi.ClientSecret,
        Amount:        amountCents,
        Currency:      string(pi.Currency),
        InvoiceNumber: invoiceNumber,
        PublicKey:     c.publicKey,
    }, nil
}
```

**Validaciones previas** (en el handler):
- ✅ Factura existe
- ✅ `allow_online_payments = true`
- ✅ Factura no está pagada (`balance > 0`)

### 5. Procesamiento de Webhooks

**Archivo**: `pkg/infrastructure/stripe/webhook.go`

#### Parseo del Evento

```go
func (c *Client) ParseWebhookEvent(r *http.Request) (*stripe.Event, error) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read request body: %w", err)
    }

    var event stripe.Event
    if err := json.Unmarshal(body, &event); err != nil {
        return nil, fmt.Errorf("failed to parse webhook event: %w", err)
    }

    return &event, nil
}
```

#### Extracción de Datos del Charge

```go
type ChargeData struct {
    ChargeID      string
    Amount        float64
    InvoiceNumber string
    CustomerName  string
    CustomerEmail string
    CardLast4     string
}

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
        // Fallback: extraer de metadata del PaymentIntent
        invoiceNumber := pi.Metadata["invoice_number"]
        return &ChargeData{
            ChargeID:      pi.ID,
            Amount:        float64(pi.Amount) / 100.0,
            InvoiceNumber: invoiceNumber,
        }, nil
    }

    // Extraer invoice_number de metadata
    invoiceNumber := charge.Metadata["invoice_number"]
    if invoiceNumber == "" && pi.Metadata != nil {
        invoiceNumber = pi.Metadata["invoice_number"]
    }

    return &ChargeData{
        ChargeID:      charge.ID,
        Amount:        float64(charge.Amount) / 100.0,
        InvoiceNumber: invoiceNumber,
        CustomerName:  charge.BillingDetails.Name,
        CustomerEmail: charge.BillingDetails.Email,
        CardLast4:     charge.PaymentMethodDetails.Card.Last4,
    }, nil
}
```

#### Formato de Notas

```go
func FormatPaymentNotes(cd *ChargeData) string {
    return fmt.Sprintf("Posted by %s (%s) using card ending in %s",
        cd.CustomerName, cd.CustomerEmail, cd.CardLast4)
}
```

### 6. Endpoints HTTP

**Archivo**: `pkg/rest/handler/stripe/handler.go`

#### Endpoint: Crear PaymentIntent

```go
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
        response.Error(w, http.StatusInternalServerError, "Error al obtener factura")
        return
    }

    // Validaciones
    if !invoice.AllowOnlinePayments {
        response.Error(w, http.StatusBadRequest, "Online payments not allowed for this invoice")
        return
    }

    if invoice.IsPaid() {
        response.Error(w, http.StatusBadRequest, "Invoice is already paid")
        return
    }

    balance := *invoice.Balance
    if balance <= 0 {
        response.Error(w, http.StatusBadRequest, "Invoice has no outstanding balance")
        return
    }

    // Crear PaymentIntent en Stripe
    result, err := h.stripeClient.CreatePaymentIntent(r.Context(), invoice.InvoiceNumber, balance)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, "Error al crear payment intent")
        return
    }

    response.JSON(w, http.StatusOK, result)
}
```

#### Endpoint: Webhook

```go
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
    // Parsear evento
    event, err := h.stripeClient.ParseWebhookEvent(r)
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid webhook event")
        return
    }

    // Solo procesar payment_intent.succeeded
    if event.Type != "payment_intent.succeeded" {
        response.JSON(w, http.StatusOK, map[string]string{"message": "Event type not handled"})
        return
    }

    // Extraer datos del charge
    chargeData, err := infraStripe.ExtractChargeFromPaymentIntent(event)
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Failed to extract payment data")
        return
    }

    // Buscar factura por invoice_number
    invoice, err := h.invoiceRepo.GetByInvoiceNumber(r.Context(), chargeData.InvoiceNumber)
    if err != nil {
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
        response.Error(w, http.StatusInternalServerError, "Failed to record payment")
        return
    }

    response.JSON(w, http.StatusOK, map[string]string{"message": "Payment recorded successfully"})
}
```

### 7. Registro de Rutas

**Archivo**: `pkg/rest/router/router.go`

```go
// Rutas públicas de Stripe (sin autenticación)
if stripeHdlr != nil {
    r.Get("/api/v1/invoices/{id}/payment-intent", stripeHdlr.CreatePaymentIntent)
    r.Post("/webhooks/stripe", stripeHdlr.Webhook)
}
```

---

## Comparación con Laravel

### Endpoint de Pago

#### Laravel (Blade View)

```php
// routes/web.php
Route::get('invoices/{invoice}/pay', [InvoiceController::class, 'pay'])->name('invoices.pay');

// InvoiceController.php
public function pay(PayRequest $request, Invoice $invoice): Response
{
    $intent = null;

    if ($invoice->allow_online_payments && !$invoice->isPaid()) {
        $intent = PaymentIntent::create([
            'amount' => round((float) $invoice->balance * 100),
            'currency' => 'usd',
            'metadata' => [
                'invoice_number' => $invoice->invoice_number,
            ],
        ]);
    }

    return response()->view('pages.invoices.pay', compact('invoice', 'intent'));
}
```

**Retorna**: Vista Blade con formulario de pago integrado (Stripe.js)

#### Go (JSON API)

```go
// GET /api/v1/invoices/{id}/payment-intent
func (h *Handler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
    // ... validaciones ...

    result, err := h.stripeClient.CreatePaymentIntent(r.Context(), invoice.InvoiceNumber, balance)

    response.JSON(w, http.StatusOK, result)
}
```

**Retorna**: JSON con `clientSecret`, `publicKey`, etc. (frontend separado)

### Webhook

#### Laravel

```php
public function webhook(Request $request): Response
{
    try {
        $event = StripeEvent::constructFrom($request->all());

        if ($event->type == 'payment_intent.succeeded') {
            foreach ($event->data->object->charges->data as $charge) {
                $amount = round($charge->amount / 100, 2);
                $invoiceNumber = $charge->metadata->invoice_number;

                $invoice = Invoice::where('invoice_number', $invoiceNumber)->firstOrFail();

                InvoicePayment::create([
                    'invoice_id' => $invoice->id,
                    'amount' => $amount,
                    'payment_id' => $charge->id,
                    'payment_processor' => 'Stripe',
                    'notes' => 'Posted by ' . $charge->billing_details->name . '...',
                ]);
            }
        }

        return response('All events processed', 200);
    } catch (Exception $e) {
        return response($e->getMessage(), 500);
    }
}
```

**Diferencias**:
- Laravel itera sobre `charges->data` (array)
- Go usa `LatestCharge` (objeto único) del SDK v76

#### Go

```go
func (h *Handler) Webhook(w http.ResponseWriter, r *http.Request) {
    event, err := h.stripeClient.ParseWebhookEvent(r)
    // ...

    chargeData, err := infraStripe.ExtractChargeFromPaymentIntent(event)
    // ...

    invoice, err := h.invoiceRepo.GetByInvoiceNumber(r.Context(), chargeData.InvoiceNumber)
    // ...

    payment := &domainPayment.InvoicePayment{
        InvoiceID:        invoice.ID,
        Amount:           chargeData.Amount,
        PaymentID:        chargeData.ChargeID,
        PaymentProcessor: "Stripe",
        Notes:            infraStripe.FormatPaymentNotes(chargeData),
    }

    h.paymentUC.Create(r.Context(), payment)
}
```

---

## Tabla `invoice_payments`

### Estructura

```sql
CREATE TABLE `invoice_payments` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `invoice_id` bigint(20) unsigned NOT NULL,
  `payment_processor` varchar(255) NOT NULL,
  `payment_id` varchar(255) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `notes` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `invoice_payments_invoice_id_foreign` (`invoice_id`),
  CONSTRAINT `invoice_payments_invoice_id_foreign` FOREIGN KEY (`invoice_id`) REFERENCES `invoices` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Ejemplo de Registro

```json
{
  "id": 1,
  "invoice_id": 123,
  "payment_processor": "Stripe",
  "payment_id": "ch_3QRSTuvwxyz123456",
  "amount": 1250.00,
  "notes": "Posted by John Doe (john@example.com) using card ending in 4242",
  "created_at": "2025-06-15T10:30:00Z",
  "updated_at": "2025-06-15T10:30:00Z",
  "deleted_at": null
}
```

---

## Seguridad

### ⚠️ Verificación de Firma de Webhook (PENDIENTE)

**Estado actual**: ❌ No implementada

**Riesgo**: Cualquiera puede enviar webhooks falsos al endpoint `/webhooks/stripe`

**Solución recomendada**:

```go
func (c *Client) VerifyWebhookSignature(r *http.Request, webhookSecret string) (*stripe.Event, error) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        return nil, err
    }

    signature := r.Header.Get("Stripe-Signature")

    event, err := webhook.ConstructEvent(body, signature, webhookSecret)
    if err != nil {
        return nil, fmt.Errorf("webhook signature verification failed: %w", err)
    }

    return &event, nil
}
```

**Configuración**:
1. Obtener `STRIPE_WEBHOOK_SECRET` del dashboard de Stripe
2. Agregar a `configs/app.env`
3. Modificar `ParseWebhookEvent` para usar `VerifyWebhookSignature`

### Otras Consideraciones

- ✅ Endpoints públicos (sin autenticación) - correcto para webhooks
- ✅ Validación de `allow_online_payments` antes de crear PaymentIntent
- ✅ Validación de balance > 0
- ⚠️ No hay rate limiting en endpoints públicos (considerar implementar)
- ⚠️ No hay idempotencia en webhook (podría procesar el mismo evento dos veces)

---

## Testing

Ver documento: [STRIPE_TESTING.md](./STRIPE_TESTING.md)

---

## Troubleshooting

### Error: "STRIPE_SECRET_KEY not configured"

**Causa**: Variable de entorno no configurada

**Solución**:
```bash
# configs/app.env
STRIPE_SECRET_KEY=sk_test_***
STRIPE_PUBLIC_KEY=pk_test_***
```

### Error: "Invoice not found" en webhook

**Causa**: El `invoice_number` en metadata no coincide con ninguna factura

**Debug**:
```bash
# Ver logs del webhook
grep "Invoice not found for Stripe webhook" logs/app.log
```

**Verificar**:
1. Que el PaymentIntent se creó con el `invoice_number` correcto
2. Que la factura existe en la BD

### Webhook no se ejecuta

**Causas posibles**:
1. Stripe no puede alcanzar el servidor (firewall, localhost)
2. URL del webhook mal configurada en Stripe dashboard
3. Evento no es `payment_intent.succeeded`

**Solución**:
- Usar **Stripe CLI** para testing local (ver STRIPE_TESTING.md)
- Verificar URL en Stripe Dashboard → Webhooks
- Revisar logs de Stripe Dashboard → Webhooks → Attempts

---

## Roadmap de Mejoras

### Prioridad Alta
- [ ] Implementar verificación de firma de webhook (`STRIPE_WEBHOOK_SECRET`)
- [ ] Agregar idempotencia en webhook (evitar pagos duplicados)

### Prioridad Media
- [ ] Rate limiting en endpoints públicos
- [ ] Retry logic en caso de fallo al registrar pago
- [ ] Notificación por email al cliente tras pago exitoso

### Prioridad Baja
- [ ] Soporte para múltiples monedas
- [ ] Refunds (reembolsos) desde la aplicación
- [ ] Dashboard de métricas de pagos

---

## Referencias

- [Stripe API Docs](https://stripe.com/docs/api)
- [Stripe Go SDK](https://github.com/stripe/stripe-go)
- [Payment Intents Guide](https://stripe.com/docs/payments/payment-intents)
- [Webhooks Guide](https://stripe.com/docs/webhooks)
