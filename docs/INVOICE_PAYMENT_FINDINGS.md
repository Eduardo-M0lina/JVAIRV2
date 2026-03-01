# Hallazgos: Invoice Payment en JVAIR (Proyecto Base Laravel)

## Resumen Ejecutivo

El proyecto base JVAIR implementa un **sistema dual de pagos** para invoices:

1. **Pagos online vía Stripe** (tarjeta de crédito) — orientado al cliente final
2. **Pagos manuales** (ACH / Check) — registrados por administradores

Ambos tipos de pago se almacenan en la misma tabla `invoice_payments`, diferenciados por el campo `payment_processor`.

> **Dato importante:** La tabla `invoice_payments` en producción está **vacía** — no hay registros de pagos reales en el dump de producción (`Data_prd.sql`).

---

## 1. Estructura de Base de Datos

### Tabla `invoices`

| Columna | Tipo | Notas |
|---------|------|-------|
| `id` | bigint unsigned (PK) | Auto-increment |
| `job_id` | bigint unsigned (FK → jobs) | Relación con trabajo |
| `invoice_number` | varchar(255) | Número de factura |
| `total` | decimal(8,2) | Total de la factura (default 0.00) |
| `description` | text (nullable) | Descripción |
| `allow_online_payments` | tinyint(1) | **Flag para habilitar pagos Stripe** (default 0/false) |
| `notes` | text (nullable) | Notas |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |
| `deleted_at` | timestamp | Soft delete |

### Tabla `invoice_payments`

| Columna | Tipo | Notas |
|---------|------|-------|
| `id` | bigint unsigned (PK) | Auto-increment |
| `invoice_id` | bigint unsigned (FK → invoices) | Relación con factura |
| `payment_processor` | varchar(255) | `"Stripe"`, `"ACH"`, o `"Check"` |
| `payment_id` | varchar(255) | ID externo del pago (Stripe charge ID o ref manual) |
| `amount` | decimal(8,2) | Monto del pago |
| `notes` | text | Notas del pago |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |
| `deleted_at` | timestamp | Soft delete |

---

## 2. Integración con Stripe (Pagos Online)

### SÍ existe integración real con Stripe. Estos son los componentes:

### 2.1. Dependencia PHP
- **Paquete:** `stripe/stripe-php` v7.67+ en `composer.json`

### 2.2. Configuración
- **Archivo:** `config/payments.php`
  ```php
  'stripe' => [
      'public_key' => env('STRIPE_PUBLIC_KEY'),
      'secret_key' => env('STRIPE_SECRET_KEY'),
  ]
  ```
- **Variables de entorno** en `.env`:
  ```
  STRIPE_PUBLIC_KEY=
  STRIPE_SECRET_KEY=
  ```
  > Ambas están **vacías** en el `.env` del repo — nunca se configuraron (o se removieron).

### 2.3. Inicialización del SDK
- **Archivo:** `app/Providers/AppServiceProvider.php` línea 43
  ```php
  Stripe::setApiKey(config('payments.stripe.secret_key'));
  ```
  Se inicializa **globalmente** en el boot de la aplicación.

### 2.4. Creación de PaymentIntent (Backend)
- **Archivo:** `app/Http/Controllers/InvoiceController.php` → método `pay()`
  ```php
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
  - Crea un `PaymentIntent` de Stripe con el balance pendiente en centavos
  - Solo si `allow_online_payments` es true y la factura no está pagada
  - Pasa el `invoice_number` como metadata para identificar el pago en el webhook

### 2.5. Frontend de Pago con Stripe Elements
- **Vista:** `resources/views/pages/invoices/pay.blade.php`
  - Layout público (`layouts/public`) — **no requiere autenticación**
  - Muestra detalles de la factura (número, work order, total, balance)
  - Integra Stripe Elements para capturar datos de tarjeta
  - Campos: Full Name, Email, Card Details (Stripe Elements)
  - Botón "Submit Payment" con el `client_secret` del PaymentIntent

- **JavaScript:** `resources/js/admin/stripe.js`
  - Carga Stripe.js (`https://js.stripe.com/v3/`) en el layout público
  - Inicializa `Stripe(window.stripe_public_key)` con la key pública desde config
  - Crea un elemento `card` de Stripe Elements
  - En submit, llama a `stripe.confirmCardPayment(secret, {...})`
  - Envía `billing_details` (name, email)
  - En éxito, oculta el formulario y muestra mensaje de confirmación

- **Layout público:** `resources/views/layouts/public.blade.php`
  ```html
  <script>window.stripe_public_key="{{ config('payments.stripe.public_key') }}";</script>
  <script src="https://js.stripe.com/v3/"></script>
  ```

### 2.6. Webhook de Stripe (Confirmación Asíncrona)
- **Archivo:** `app/Http/Controllers/StripeWebhookController.php`
- **Ruta:** `POST /stripe/webhook` (sin autenticación, pública)
- **Flujo:**
  1. Recibe evento de Stripe
  2. Filtra solo `payment_intent.succeeded`
  3. Itera sobre los `charges` del PaymentIntent
  4. Extrae `invoice_number` del metadata del charge
  5. Busca la factura por `invoice_number`
  6. Crea un `InvoicePayment` con:
     - `payment_processor` = `"Stripe"`
     - `payment_id` = ID del charge de Stripe
     - `amount` = monto en dólares (charge.amount / 100)
     - `notes` = "Posted by [nombre] ([email]) using card ending in [last4]"

> **Observación de seguridad:** El webhook NO verifica la firma de Stripe (`Stripe-Signature` header). Usa `StripeEvent::constructFrom($request->all())` en lugar de `Webhook::constructEvent()`. Esto es una **vulnerabilidad** — cualquiera podría enviar eventos falsos al endpoint.

---

## 3. Pagos Manuales (ACH / Check)

### 3.1. Controller
- **Archivo:** `app/Http/Controllers/InvoicePaymentController.php`
- CRUD completo: create, store, edit, update, destroy
- **Solo para uso admin** (requiere autenticación)

### 3.2. Validación (Store/Update)
- **Archivo:** `app/Http/Requests/InvoicePayments/StoreRequest.php`
  ```php
  'payment_id' => ['required', 'string'],
  'payment_processor' => ['required', 'string', 'in:ACH,Check'],  // Solo ACH o Check
  'amount' => ['required', 'numeric'],
  'date' => ['nullable', 'date'],
  'notes' => ['nullable', 'string']
  ```
  > **Nota:** Los pagos manuales solo permiten `ACH` o `Check` como processor. `Stripe` NO es opción manual.

### 3.3. Vistas
- **Crear pago:** `resources/views/pages/invoices/payments/create.blade.php`
  - Campos: Payment ID (texto libre), Payment Method (select: ACH/Check), Amount, Date, Notes
  - El amount tiene max = balance de la factura
- **Editar pago:** `resources/views/pages/invoices/payments/edit.blade.php`
  - Mismos campos, con opción de eliminar

---

## 4. Políticas de Autorización

**Archivo:** `app/Policies/InvoicePaymentPolicy.php`

| Acción | Regla |
|--------|-------|
| viewAny | Solo admin |
| view | Solo admin |
| create | Solo admin |
| **update** | Solo admin **Y solo si el payment_processor NO es "Stripe"** |
| **delete** | Solo admin **Y solo si el payment_processor NO es "Stripe"** |
| restore | Solo admin |
| forceDelete | Solo admin |

> **Hallazgo importante:** Los pagos de Stripe son **inmutables** — no se pueden editar ni eliminar desde el panel admin. Solo se pueden modificar pagos manuales (ACH/Check).

---

## 5. Rutas

### Rutas públicas (sin auth):
| Método | Ruta | Acción |
|--------|------|--------|
| GET | `/invoices/{invoice}/pay` | Formulario de pago Stripe |
| POST | `/stripe/webhook` | Webhook de Stripe |

### Rutas admin (con auth) — **COMENTADAS/DESHABILITADAS:**
Las rutas de invoices y payments están **comentadas** en `routes/web.php` (líneas 324-342):
```php
// // Invoices
// Route::prefix('invoices')->name('invoices.')->group(function() {
//     ...
//     // Payments
//     Route::prefix('/{invoice}/payments')->name('payments.')->group(function() {
//         Route::get('/create', [InvoicePaymentController::class, 'create'])->name('create');
//         Route::post('/create', [InvoicePaymentController::class, 'store'])->name('store');
//         ...
//     });
// });
```

> **Hallazgo crítico:** Las rutas CRUD de invoices y payments para admin están **deshabilitadas** en producción. Solo están activas las rutas de pago online y webhook. Esto significa que **en producción el módulo de invoices/payments estaba desactivado** del panel admin.

---

## 6. Modelo Invoice — Cálculo de Balance

**Archivo:** `app/Models/Invoice.php`

```php
// Balance = total - sum(payments.amount)
public function getBalanceAttribute($attr)
{
    return $this->asDecimal($this->total - $this->payments->sum('amount'), 2);
}

// Factura pagada cuando balance == 0
public function isPaid(): bool
{
    return (float) $this->balance == 0;
}
```

- Soporta pagos parciales (el balance se calcula dinámicamente)
- Scopes para filtrar: `withBalance()` (facturas con saldo) y `withoutBalance()` (facturas pagadas)

---

## 7. Email de Invoice

**Archivo:** `app/Mail/Invoice.php` + `resources/views/emails/invoice.blade.php`

- Se envía por email al cliente con detalles de la factura
- Si `allow_online_payments` es true y hay balance pendiente → incluye botón **"Pay Online"** con link a `/invoices/{id}/pay`
- Si no permite pagos online → muestra "Please contact JVAIR directly to make a payment"

---

## 8. Hallazgos Clave para la Migración a JVAIRV2

### 8.1. ¿Qué existe?
- ✅ Integración completa con Stripe (PaymentIntent + Elements + Webhook)
- ✅ Sistema de pagos manuales (ACH/Check) con CRUD admin
- ✅ Tabla `invoice_payments` en la BD
- ✅ Políticas de autorización diferenciadas por tipo de pago
- ✅ Email de invoice con link de pago online

### 8.2. ¿Qué NO existe?
- ❌ **No hay datos en producción** — la tabla `invoice_payments` está vacía
- ❌ Las rutas de invoices/payments están **deshabilitadas** en `web.php`
- ❌ No hay verificación de firma en el webhook de Stripe (vulnerabilidad)
- ❌ No hay Service/UseCase layer — la lógica está directamente en controllers
- ❌ No hay logging ni auditoría de pagos
- ❌ No hay manejo de refunds (reembolsos)
- ❌ No hay manejo de pagos fallidos (solo `payment_intent.succeeded`)
- ❌ No hay retry logic ni idempotency en el webhook

### 8.3. Estado en JVAIRV2 (Go)
- ❌ **NO hay ninguna implementación** de invoice payments en JVAIRV2 todavía
- ❌ No hay modelos, handlers, repositorios ni rutas de payments en Go
- ❌ No hay integración con Stripe en el proyecto Go

### 8.4. Decisiones a tomar para JVAIRV2
1. **¿Se necesita Stripe?** — Dado que la tabla está vacía y las rutas estaban deshabilitadas, ¿el cliente realmente usa pagos online?
2. **¿Se migran solo pagos manuales?** — Podría ser suficiente un CRUD de pagos manuales (ACH/Check) inicialmente
3. **Si se implementa Stripe:** Usar la API actual (PaymentIntent) con verificación de firma de webhook, idempotency keys, y logging adecuado
4. **Procesadores de pago permitidos:** `Stripe`, `ACH`, `Check` — considerar si se necesitan más opciones

---

## 9. Archivos Relevantes del Proyecto Base

| Archivo | Descripción |
|---------|-------------|
| `app/Models/Invoice.php` | Modelo con balance, isPaid, relación payments |
| `app/Models/InvoicePayment.php` | Modelo de pago con fillable y relación invoice |
| `app/Http/Controllers/InvoiceController.php` | Controller con método `pay()` (Stripe PaymentIntent) |
| `app/Http/Controllers/InvoicePaymentController.php` | CRUD de pagos manuales |
| `app/Http/Controllers/StripeWebhookController.php` | Webhook para confirmar pagos Stripe |
| `app/Policies/InvoicePaymentPolicy.php` | Políticas (Stripe inmutable, manual editable) |
| `app/Http/Requests/InvoicePayments/StoreRequest.php` | Validación: solo ACH/Check manual |
| `app/Providers/AppServiceProvider.php` | Inicialización global de Stripe API key |
| `config/payments.php` | Config de Stripe (public_key, secret_key) |
| `resources/views/pages/invoices/pay.blade.php` | Vista pública de pago con Stripe Elements |
| `resources/js/admin/stripe.js` | JS frontend: Stripe Elements + confirmCardPayment |
| `resources/views/layouts/public.blade.php` | Layout público con Stripe.js SDK |
| `resources/views/pages/invoices/payments/create.blade.php` | Vista admin: crear pago manual |
| `resources/views/pages/invoices/payments/edit.blade.php` | Vista admin: editar pago manual |
| `resources/views/emails/invoice.blade.php` | Email con botón "Pay Online" |
| `routes/web.php` | Rutas (invoice CRUD **comentado**, pay y webhook activos) |
| `database/migrations/2020_12_01_173353_create_invoice_payments_table.php` | Migración de BD |
