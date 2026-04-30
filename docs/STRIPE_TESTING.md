# Guía de Pruebas - Integración con Stripe

> **Documento**: Guía completa para probar la integración de pagos con Stripe
> **Última actualización**: Junio 2025

---

## Tabla de Contenidos

1. [Configuración Inicial](#configuración-inicial)
2. [Pruebas Locales con Stripe CLI](#pruebas-locales-con-stripe-cli)
3. [Pruebas Manuales](#pruebas-manuales)
4. [Casos de Prueba](#casos-de-prueba)
5. [Tarjetas de Prueba](#tarjetas-de-prueba)
6. [Troubleshooting](#troubleshooting)

> 📘 **Guía detallada de configuración del webhook**: Ver [`STRIPE_WEBHOOK_SETUP.md`](./STRIPE_WEBHOOK_SETUP.md)

---

## Configuración Inicial

### 1. Obtener Credenciales de Stripe

#### Modo Test (Desarrollo)

1. Crear cuenta en [Stripe Dashboard](https://dashboard.stripe.com/register)
2. Activar **modo de prueba** (toggle en la esquina superior izquierda)
3. Ir a **Developers → API keys**
4. Copiar:
   - **Publishable key** (pk_test_...)
   - **Secret key** (sk_test_...)

#### Modo Live (Producción)

⚠️ **Solo después de completar todas las pruebas**

1. Completar activación de cuenta en Stripe
2. Cambiar a **modo live**
3. Obtener claves de producción (pk_live_..., sk_live_...)

### 2. Configurar Webhook en Stripe Dashboard

**⚠️ IMPORTANTE**: Stripe necesita saber la URL de tu webhook para enviar eventos.

#### Para Desarrollo Local

**Opción recomendada**: Usar Stripe CLI (ver sección siguiente), que automáticamente configura el webhook temporal.

#### Para Producción o Testing con URL Pública

1. Ir a [Stripe Dashboard → Developers → Webhooks](https://dashboard.stripe.com/test/webhooks)
2. Click en **"Add endpoint"**
3. Configurar:
   - **Endpoint URL**: `https://tu-dominio.com/webhooks/stripe`
     - Ejemplo desarrollo: `https://dev.jvair.com/webhooks/stripe`
     - Ejemplo producción: `https://app.jvair.com/webhooks/stripe`
   - **Description**: "JVAIR Invoice Payments"
   - **Events to send**: Seleccionar solo `payment_intent.succeeded`
4. Click en **"Add endpoint"**
5. Copiar el **Signing secret** (whsec_...)
6. Agregar a `configs/app.env` como `STRIPE_WEBHOOK_SECRET`

**Captura de pantalla de configuración**:
```
┌─────────────────────────────────────────────┐
│ Add endpoint                                │
├─────────────────────────────────────────────┤
│ Endpoint URL *                              │
│ https://app.jvair.com/webhooks/stripe       │
│                                             │
│ Description                                 │
│ JVAIR Invoice Payments                      │
│                                             │
│ Listen to                                   │
│ ○ Events on your account                    │
│ ○ Events on Connected accounts              │
│                                             │
│ Select events to listen to                  │
│ ☑ payment_intent.succeeded                  │
│                                             │
│ [Add endpoint]                              │
└─────────────────────────────────────────────┘
```

**Verificar configuración**:
- En la lista de webhooks, debe aparecer tu endpoint
- Status: **Enabled** ✅
- Signing secret visible al hacer click en el endpoint

### 3. Configurar Variables de Entorno

**Archivo**: `configs/app.env`

```env
# === Stripe (Test Mode) ===
STRIPE_SECRET_KEY=sk_test_51AbCdEfGhIjKlMnOpQrStUvWxYz1234567890
STRIPE_PUBLIC_KEY=pk_test_51AbCdEfGhIjKlMnOpQrStUvWxYz1234567890
STRIPE_WEBHOOK_SECRET=whsec_1234567890abcdefghijklmnopqrstuvwxyz
```

⚠️ **Nota**: El `STRIPE_WEBHOOK_SECRET` es diferente para cada webhook endpoint que configures.

### 4. Iniciar el Servidor

```bash
cd /Users/eduardo/projects/jvair/JVAIRV2
go run cmd/api/main.go
```

**Verificar**:
```bash
curl http://localhost:8090/health
# Debe retornar: {"status":"ok"}
```

---

## Pruebas Locales con Stripe CLI

### Instalación de Stripe CLI

#### macOS (Homebrew)
```bash
brew install stripe/stripe-cli/stripe
```

#### Linux
```bash
wget https://github.com/stripe/stripe-cli/releases/download/v1.19.0/stripe_1.19.0_linux_x86_64.tar.gz
tar -xvf stripe_1.19.0_linux_x86_64.tar.gz
sudo mv stripe /usr/local/bin/
```

#### Windows
```powershell
scoop bucket add stripe https://github.com/stripe/scoop-stripe-cli.git
scoop install stripe
```

### Autenticación

```bash
stripe login
```

Esto abrirá el navegador para autorizar el CLI.

### Reenviar Webhooks a Localhost

```bash
stripe listen --forward-to localhost:8090/webhooks/stripe
```

**Output esperado**:
```
> Ready! Your webhook signing secret is whsec_1234567890abcdefghijklmnopqrstuvwxyz (^C to quit)
```

⚠️ **Importante**: Copiar el `webhook signing secret` y agregarlo a `configs/app.env` como `STRIPE_WEBHOOK_SECRET`

### Enviar Evento de Prueba

En otra terminal:

```bash
stripe trigger payment_intent.succeeded
```

**Verificar en logs del servidor**:
```
INFO Stripe webhook received eventType=payment_intent.succeeded
INFO Stripe payment recorded successfully invoiceId=123 invoiceNumber=INV-001 chargeId=ch_xxx amount=100.00
```

---

## Pruebas Manuales

### Flujo Completo: Pagar una Factura

#### Paso 1: Crear una Factura de Prueba

**Opción A: Desde la BD directamente**

```sql
-- Crear cliente
INSERT INTO customers (name, email, phone, created_at, updated_at)
VALUES ('John Doe', 'john@example.com', '555-1234', NOW(), NOW());

-- Crear propiedad
INSERT INTO properties (customer_id, street, city, state, zip, created_at, updated_at)
VALUES (LAST_INSERT_ID(), '123 Main St', 'Atlanta', 'GA', '30301', NOW(), NOW());

-- Crear job
INSERT INTO jobs (property_id, work_order, received_date, created_at, updated_at)
VALUES (LAST_INSERT_ID(), 'WO-001', NOW(), NOW(), NOW());

-- Crear factura
INSERT INTO invoices (
    job_id,
    invoice_number,
    description,
    total,
    allow_online_payments,
    created_at,
    updated_at
)
VALUES (
    LAST_INSERT_ID(),
    'INV-TEST-001',
    'Test invoice for Stripe integration',
    150.00,
    1,  -- Permitir pagos online
    NOW(),
    NOW()
);
```

**Opción B: Usar el API**

```bash
# 1. Login
curl -X POST http://localhost:8090/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@jvair.com",
    "password": "your_password"
  }'

# Guardar el token de la respuesta
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 2. Crear factura (requiere customer, property, job previos)
curl -X POST http://localhost:8090/api/v1/invoices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jobId": 1,
    "invoiceNumber": "INV-TEST-001",
    "description": "Test invoice",
    "total": 150.00,
    "allowOnlinePayments": true
  }'
```

#### Paso 2: Obtener PaymentIntent

```bash
# Sin autenticación (endpoint público)
curl http://localhost:8090/api/v1/invoices/1/payment-intent
```

**Respuesta esperada**:
```json
{
  "id": "pi_3QRSTuvwxyz123456",
  "clientSecret": "pi_3QRSTuvwxyz123456_secret_abcdefghijklmnop",
  "amount": 15000,
  "currency": "usd",
  "invoiceNumber": "INV-TEST-001",
  "publicKey": "pk_test_51AbCdEfGhIjKlMnOpQrStUvWxYz1234567890"
}
```

#### Paso 3: Simular Pago desde Frontend

**Opción A: Usar Stripe Dashboard (Test Mode)**

1. Ir a [Stripe Dashboard → Payments](https://dashboard.stripe.com/test/payments)
2. Buscar el PaymentIntent por ID (`pi_3QRSTuvwxyz123456`)
3. Click en "Confirm" para simular pago exitoso

**Opción B: Usar HTML de Prueba**

Crear archivo `test_payment.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Test Stripe Payment</title>
    <script src="https://js.stripe.com/v3/"></script>
</head>
<body>
    <h1>Test Payment for Invoice INV-TEST-001</h1>

    <form id="payment-form">
        <div id="card-element"></div>
        <button type="submit">Pay $150.00</button>
        <div id="error-message"></div>
        <div id="success-message"></div>
    </form>

    <script>
        // Reemplazar con tu STRIPE_PUBLIC_KEY
        const stripe = Stripe('pk_test_51AbCdEfGhIjKlMnOpQrStUvWxYz1234567890');

        // Reemplazar con el clientSecret del paso 2
        const clientSecret = 'pi_3QRSTuvwxyz123456_secret_abcdefghijklmnop';

        const elements = stripe.elements();
        const cardElement = elements.create('card');
        cardElement.mount('#card-element');

        const form = document.getElementById('payment-form');
        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const {error, paymentIntent} = await stripe.confirmCardPayment(clientSecret, {
                payment_method: {
                    card: cardElement,
                    billing_details: {
                        name: 'John Doe',
                        email: 'john@example.com'
                    }
                }
            });

            if (error) {
                document.getElementById('error-message').textContent = error.message;
            } else if (paymentIntent.status === 'succeeded') {
                document.getElementById('success-message').textContent = 'Payment succeeded!';
                console.log('PaymentIntent:', paymentIntent);
            }
        });
    </script>
</body>
</html>
```

Abrir en navegador y usar tarjeta de prueba: `4242 4242 4242 4242`

#### Paso 4: Verificar Webhook Recibido

**En logs del servidor**:
```
INFO Stripe webhook received eventType=payment_intent.succeeded
INFO Stripe payment recorded successfully invoiceId=1 invoiceNumber=INV-TEST-001 chargeId=ch_3QRSTuvwxyz789012 amount=150.00
```

#### Paso 5: Verificar Pago Registrado

```bash
curl -X GET http://localhost:8090/api/v1/invoices/1/payments \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta esperada**:
```json
{
  "data": [
    {
      "id": 1,
      "invoiceId": 1,
      "paymentProcessor": "Stripe",
      "paymentId": "ch_3QRSTuvwxyz789012",
      "amount": 150.00,
      "notes": "Posted by John Doe (john@example.com) using card ending in 4242",
      "createdAt": "2025-06-15T10:30:00Z"
    }
  ],
  "total": 1
}
```

---

## Casos de Prueba

### ✅ Caso 1: Pago Exitoso

**Precondiciones**:
- Factura con `allow_online_payments = true`
- Balance > 0

**Pasos**:
1. Crear PaymentIntent
2. Confirmar pago con tarjeta válida
3. Verificar webhook recibido
4. Verificar pago registrado en `invoice_payments`

**Resultado esperado**: Pago registrado correctamente

---

### ❌ Caso 2: Factura No Permite Pagos Online

**Precondiciones**:
- Factura con `allow_online_payments = false`

**Request**:
```bash
curl http://localhost:8090/api/v1/invoices/1/payment-intent
```

**Respuesta esperada**:
```json
{
  "error": "Online payments not allowed for this invoice"
}
```

**Status**: `400 Bad Request`

---

### ❌ Caso 3: Factura Ya Pagada

**Precondiciones**:
- Factura con balance = 0 (ya tiene pagos que cubren el total)

**Request**:
```bash
curl http://localhost:8090/api/v1/invoices/1/payment-intent
```

**Respuesta esperada**:
```json
{
  "error": "Invoice is already paid"
}
```

**Status**: `400 Bad Request`

---

### ❌ Caso 4: Factura No Existe

**Request**:
```bash
curl http://localhost:8090/api/v1/invoices/99999/payment-intent
```

**Respuesta esperada**:
```json
{
  "error": "Factura no encontrada"
}
```

**Status**: `404 Not Found`

---

### ❌ Caso 5: Tarjeta Declinada

**Pasos**:
1. Crear PaymentIntent válido
2. Usar tarjeta de prueba declinada: `4000 0000 0000 0002`

**Resultado esperado**:
- Error en frontend: "Your card was declined"
- No se crea registro en `invoice_payments`
- No se envía webhook

---

### ❌ Caso 6: Tarjeta Requiere Autenticación 3D Secure

**Pasos**:
1. Crear PaymentIntent válido
2. Usar tarjeta de prueba: `4000 0027 6000 3184`

**Resultado esperado**:
- Stripe muestra modal de autenticación
- Usuario debe completar 3D Secure
- Pago se completa tras autenticación

---

### ⚠️ Caso 7: Webhook Duplicado (Idempotencia)

**Escenario**: Stripe envía el mismo webhook dos veces

**Pasos**:
1. Completar pago exitoso
2. Reenviar el mismo webhook manualmente:
```bash
stripe events resend evt_1234567890
```

**Resultado actual**: ⚠️ Se crea pago duplicado (bug conocido)

**Resultado esperado**: Detectar pago duplicado y rechazar

**Solución pendiente**: Implementar idempotencia basada en `payment_id`

---

## Tarjetas de Prueba

### Tarjetas Exitosas

| Número | Descripción |
|--------|-------------|
| `4242 4242 4242 4242` | Visa - Pago exitoso |
| `5555 5555 5555 4444` | Mastercard - Pago exitoso |
| `3782 822463 10005` | American Express - Pago exitoso |

### Tarjetas con Errores

| Número | Error |
|--------|-------|
| `4000 0000 0000 0002` | Tarjeta declinada (generic) |
| `4000 0000 0000 9995` | Fondos insuficientes |
| `4000 0000 0000 0069` | Tarjeta expirada |
| `4000 0000 0000 0127` | CVC incorrecto |

### Tarjetas con 3D Secure

| Número | Comportamiento |
|--------|----------------|
| `4000 0027 6000 3184` | Requiere autenticación |
| `4000 0025 0000 3155` | Autenticación requerida + declinada |

**Datos adicionales** (para todas las tarjetas):
- **Fecha de expiración**: Cualquier fecha futura (ej: 12/25)
- **CVC**: Cualquier 3 dígitos (ej: 123)
- **ZIP**: Cualquier código postal (ej: 12345)

---

## Troubleshooting

### Error: "STRIPE_SECRET_KEY is required"

**Causa**: Variable de entorno no configurada

**Solución**:
```bash
# Verificar archivo configs/app.env
cat configs/app.env | grep STRIPE

# Debe mostrar:
# STRIPE_SECRET_KEY=sk_test_...
# STRIPE_PUBLIC_KEY=pk_test_...
```

---

### Error: "failed to create payment intent: Invalid API Key"

**Causa**: API key inválida o expirada

**Solución**:
1. Verificar que estás en modo test en Stripe Dashboard
2. Regenerar API keys en Developers → API keys
3. Actualizar `configs/app.env`
4. Reiniciar servidor

---

### Webhook no se recibe

**Síntomas**: Pago exitoso en Stripe pero no se registra en `invoice_payments`

**Diagnóstico**:
```bash
# Verificar que Stripe CLI está corriendo
stripe listen --forward-to localhost:8090/webhooks/stripe

# Verificar logs del servidor
tail -f logs/app.log | grep "Stripe webhook"
```

**Causas comunes**:
1. Stripe CLI no está corriendo
2. URL del webhook incorrecta
3. Servidor no está corriendo en el puerto esperado
4. Firewall bloqueando conexión

---

### Error: "Invoice not found" en webhook

**Causa**: El `invoice_number` en metadata no existe en la BD

**Diagnóstico**:
```bash
# Ver el evento en Stripe Dashboard
stripe events list --limit 1

# Verificar metadata
stripe events retrieve evt_1234567890

# Buscar factura en BD
mysql -u root -p jvair -e "SELECT * FROM invoices WHERE invoice_number = 'INV-TEST-001';"
```

**Solución**: Verificar que el PaymentIntent se creó con el `invoice_number` correcto

---

### Pago duplicado

**Síntomas**: Mismo pago aparece dos veces en `invoice_payments`

**Causa**: Webhook procesado múltiples veces (no hay idempotencia)

**Solución temporal**:
```sql
-- Eliminar pago duplicado
DELETE FROM invoice_payments
WHERE id = <id_del_duplicado>;
```

**Solución permanente**: Implementar idempotencia (ver roadmap en STRIPE_INTEGRATION.md)

---

## Checklist de Pruebas Pre-Producción

Antes de activar en producción, verificar:

- [ ] Todas las tarjetas de prueba funcionan correctamente
- [ ] Webhooks se reciben y procesan correctamente
- [ ] Pagos se registran en `invoice_payments`
- [ ] Facturas ya pagadas rechazan nuevos PaymentIntents
- [ ] Facturas sin `allow_online_payments` rechazan PaymentIntents
- [ ] Logs no muestran errores
- [ ] Variables de entorno de producción configuradas
- [ ] Webhook signature verification implementada (CRÍTICO)
- [ ] Stripe Dashboard configurado con URL de producción
- [ ] Notificaciones por email funcionan (si aplica)

---

## Monitoreo en Producción

### Métricas Clave

1. **Tasa de éxito de pagos**: `(pagos exitosos / intentos) * 100`
2. **Tiempo promedio de procesamiento**: Desde PaymentIntent hasta registro en BD
3. **Webhooks fallidos**: Revisar en Stripe Dashboard → Webhooks → Attempts

### Logs a Monitorear

```bash
# Pagos exitosos
grep "Stripe payment recorded successfully" logs/app.log

# Errores en webhooks
grep "Failed to create invoice payment from Stripe webhook" logs/app.log

# Facturas no encontradas
grep "Invoice not found for Stripe webhook" logs/app.log
```

### Alertas Recomendadas

- ⚠️ Webhook con tasa de error > 5%
- ⚠️ Tiempo de procesamiento > 5 segundos
- 🚨 Pago registrado pero factura no actualizada

---

## Recursos Adicionales

- [Stripe Testing Guide](https://stripe.com/docs/testing)
- [Stripe CLI Docs](https://stripe.com/docs/stripe-cli)
- [Payment Intents API](https://stripe.com/docs/api/payment_intents)
- [Webhooks Best Practices](https://stripe.com/docs/webhooks/best-practices)
