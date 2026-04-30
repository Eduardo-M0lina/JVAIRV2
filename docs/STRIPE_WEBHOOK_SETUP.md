# Configuración del Webhook de Stripe - Guía Visual

## ¿Por qué necesito configurar el webhook?

Cuando un cliente paga una factura con Stripe, el pago se procesa en los servidores de Stripe, **no en tu servidor**. Para que tu aplicación se entere de que el pago fue exitoso, Stripe necesita **notificarte** enviando un webhook (una petición HTTP POST) a tu servidor.

```
┌──────────────┐                    ┌──────────────┐
│   Cliente    │                    │   Stripe     │
│  (Frontend)  │                    │  (Servidores)│
└──────┬───────┘                    └──────┬───────┘
       │                                   │
       │ 1. Ingresa tarjeta                │
       │    y confirma pago                │
       ├──────────────────────────────────>│
       │                                   │
       │                                   │ 2. Procesa pago
       │                                   │    (cargo a tarjeta)
       │                                   │
       │ 3. Respuesta: "Pago exitoso"      │
       │<──────────────────────────────────┤
       │                                   │
       │                                   │ 4. Envía webhook
       │                                   │    POST /webhooks/stripe
       │                                   │    {"type": "payment_intent.succeeded"}
       │                                   │
       │                                   ▼
       │                            ┌──────────────┐
       │                            │  Tu Servidor │
       │                            │   (Go API)   │
       │                            └──────┬───────┘
       │                                   │
       │                                   │ 5. Registra pago
       │                                   │    en invoice_payments
       │                                   ▼
       │                            ┌──────────────┐
       │                            │  Base de     │
       │                            │  Datos       │
       │                            └──────────────┘
```

**Sin el webhook configurado**: El pago se procesa en Stripe pero **nunca se registra en tu base de datos**.

---

## Configuración Paso a Paso

### Opción 1: Desarrollo Local (Stripe CLI)

**Ventajas**: Automático, no requiere URL pública

```bash
# 1. Instalar Stripe CLI
brew install stripe/stripe-cli/stripe

# 2. Autenticarse
stripe login

# 3. Reenviar webhooks a localhost
stripe listen --forward-to localhost:8090/webhooks/stripe
```

**Output**:
```
> Ready! Your webhook signing secret is whsec_abc123... (^C to quit)
```

**Copiar el signing secret** y agregarlo a `configs/app.env`:
```env
STRIPE_WEBHOOK_SECRET=whsec_abc123...
```

**Mantener el comando corriendo** mientras desarrollas.

---

### Opción 2: Servidor con URL Pública (Staging/Producción)

**Requisitos**:
- Servidor accesible desde internet
- URL con HTTPS (Stripe requiere SSL)

#### Paso 1: Ir a Stripe Dashboard

**Test Mode**: https://dashboard.stripe.com/test/webhooks

![Stripe Dashboard](https://i.imgur.com/example.png)

#### Paso 2: Agregar Endpoint

Click en **"Add endpoint"**

```
┌─────────────────────────────────────────────────────────┐
│ Add endpoint                                            │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ Endpoint URL *                                          │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ https://app.jvair.com/webhooks/stripe              │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Description (optional)                                  │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ JVAIR Invoice Payments                              │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Listen to                                               │
│ ○ Events on your account                                │
│ ○ Events on Connected accounts                          │
│                                                         │
│ Version                                                 │
│ Latest API version (2024-11-20)                         │
│                                                         │
│ Select events to listen to                              │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Search events...                                    │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ ☑ payment_intent.succeeded                              │
│   Occurs when a PaymentIntent has successfully         │
│   completed payment.                                    │
│                                                         │
│                                    [Add endpoint]       │
└─────────────────────────────────────────────────────────┘
```

**Valores a configurar**:
- **Endpoint URL**: `https://app.jvair.com/webhooks/stripe`
  - ⚠️ Debe ser HTTPS (no HTTP)
  - ⚠️ Debe ser accesible desde internet
  - ⚠️ No debe requerir autenticación
- **Description**: `JVAIR Invoice Payments`
- **Events**: Solo `payment_intent.succeeded`

#### Paso 3: Obtener Signing Secret

Después de crear el endpoint:

1. Click en el endpoint en la lista
2. En la sección **"Signing secret"**, click en **"Reveal"**
3. Copiar el secret (empieza con `whsec_...`)

```
┌─────────────────────────────────────────────────────────┐
│ Endpoint details                                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ URL: https://app.jvair.com/webhooks/stripe              │
│ Status: Enabled ✅                                       │
│                                                         │
│ Signing secret                                          │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ whsec_1234567890abcdefghijklmnopqrstuvwxyz         │ │
│ └─────────────────────────────────────────────────────┘ │
│ [Copy]                                                  │
│                                                         │
│ Events being sent                                       │
│ • payment_intent.succeeded                              │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

#### Paso 4: Configurar en el Servidor

Agregar a `configs/app.env`:

```env
STRIPE_WEBHOOK_SECRET=whsec_1234567890abcdefghijklmnopqrstuvwxyz
```

Reiniciar el servidor:
```bash
# Detener servidor actual (Ctrl+C)
# Iniciar de nuevo
go run cmd/api/main.go
```

---

## Verificar que Funciona

### Método 1: Enviar Evento de Prueba desde Stripe CLI

```bash
stripe trigger payment_intent.succeeded
```

**Verificar en logs del servidor**:
```
INFO Stripe webhook received eventType=payment_intent.succeeded
```

### Método 2: Hacer un Pago de Prueba

1. Crear PaymentIntent para una factura
2. Pagar con tarjeta de prueba `4242 4242 4242 4242`
3. Verificar que el pago se registró en `invoice_payments`

### Método 3: Ver Intentos en Stripe Dashboard

Ir a **Developers → Webhooks → [Tu endpoint] → Attempts**

Verás una lista de todos los webhooks enviados:

```
┌────────────────────────────────────────────────────────┐
│ Recent attempts                                        │
├────────────────────────────────────────────────────────┤
│ ✅ payment_intent.succeeded  200 OK  2ms  Just now     │
│ ✅ payment_intent.succeeded  200 OK  3ms  5 min ago    │
│ ❌ payment_intent.succeeded  500     2ms  1 hour ago   │
└────────────────────────────────────────────────────────┘
```

- ✅ **200 OK**: Webhook procesado correctamente
- ❌ **500/400**: Error en tu servidor (revisar logs)

---

## Troubleshooting

### ❌ Error: "No signature found"

**Causa**: El `STRIPE_WEBHOOK_SECRET` no está configurado o es incorrecto.

**Solución**:
1. Verificar que `STRIPE_WEBHOOK_SECRET` está en `configs/app.env`
2. Verificar que el secret coincide con el de Stripe Dashboard
3. Reiniciar el servidor

---

### ❌ Error: "Webhook signature verification failed"

**Causa**: El secret es incorrecto o el endpoint en Stripe Dashboard apunta a otra URL.

**Solución**:
1. Ir a Stripe Dashboard → Webhooks
2. Verificar que la URL del endpoint es correcta
3. Regenerar el signing secret si es necesario
4. Actualizar `configs/app.env` con el nuevo secret

---

### ❌ Webhook no llega (timeout)

**Causas posibles**:
- Servidor no está corriendo
- Firewall bloqueando la conexión
- URL incorrecta en Stripe Dashboard
- Servidor no es accesible desde internet (para producción)

**Diagnóstico**:
```bash
# Verificar que el servidor está corriendo
curl http://localhost:8090/health

# Verificar que el endpoint webhook responde
curl -X POST http://localhost:8090/webhooks/stripe \
  -H "Content-Type: application/json" \
  -d '{"type":"test"}'
```

**Para producción**:
```bash
# Verificar que el servidor es accesible desde internet
curl https://app.jvair.com/health
```

---

### ❌ Pago exitoso pero no se registra en BD

**Síntomas**: El pago aparece en Stripe Dashboard pero no en `invoice_payments`

**Diagnóstico**:
1. Verificar logs del servidor: `grep "Stripe webhook" logs/app.log`
2. Verificar en Stripe Dashboard → Webhooks → Attempts
3. Si el webhook muestra 200 OK pero no se registra, revisar logs de errores

**Causas comunes**:
- Invoice no existe (invoice_number incorrecto)
- Error en la BD (conexión, permisos)
- Error en el código (revisar logs)

---

## Checklist de Configuración

### Desarrollo Local
- [ ] Stripe CLI instalado
- [ ] `stripe login` ejecutado
- [ ] `stripe listen --forward-to localhost:8090/webhooks/stripe` corriendo
- [ ] `STRIPE_WEBHOOK_SECRET` copiado a `configs/app.env`
- [ ] Servidor corriendo en puerto 8090
- [ ] Webhook de prueba recibido correctamente

### Producción
- [ ] Endpoint agregado en Stripe Dashboard
- [ ] URL es HTTPS
- [ ] URL es accesible desde internet
- [ ] Solo evento `payment_intent.succeeded` seleccionado
- [ ] Signing secret copiado a `configs/app.env`
- [ ] Servidor reiniciado con nueva configuración
- [ ] Pago de prueba procesado correctamente
- [ ] Webhook muestra 200 OK en Stripe Dashboard

---

## Resumen

**Sin configurar el webhook**:
```
Cliente paga → Stripe procesa → ❌ Tu servidor nunca se entera
```

**Con webhook configurado**:
```
Cliente paga → Stripe procesa → ✅ Stripe notifica a tu servidor → ✅ Pago registrado en BD
```

**URLs importantes**:
- Test webhooks: https://dashboard.stripe.com/test/webhooks
- Live webhooks: https://dashboard.stripe.com/webhooks
- Documentación: https://stripe.com/docs/webhooks
