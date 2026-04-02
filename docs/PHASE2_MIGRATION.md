# Fase 2: Funcionalidades Pendientes de Migración (Laravel → Go)

> **Fecha**: Junio 2025
> **Referencia**: [MIGRATION_PLAN.md](./MIGRATION_PLAN.md) — Fase 1 completada (16 módulos CRUD migrados)

---

## Resumen Ejecutivo

La Fase 1 migró exitosamente los 16 módulos CRUD del proyecto Laravel a Go. Sin embargo, existen funcionalidades transversales y de negocio que aún no han sido migradas. Este documento detalla cada una, su estado actual en ambos proyectos, las variables de entorno necesarias y las recomendaciones de implementación.

---

## 1. ✅ Envío de Email (Mailgun) - COMPLETADO

### Estado en Laravel
El proyecto usa **Mailgun** como proveedor de email con 6 clases Mailable:

| Mailable | Archivo | Propósito | Template |
|----------|---------|-----------|----------|
| `Dispatch` | `app/Mail/Dispatch.php` | Email de dispatch a técnicos | `emails.dispatch` |
| `DispatchSupervisor` | `app/Mail/DispatchSupervisor.php` | Email de dispatch a supervisores (con adjuntos) | `emails.dispatch-supervisor` |
| `Invoice` | `app/Mail/Invoice.php` | Envío de factura al cliente | `emails.invoice` |
| `Quote` | `app/Mail/Quote.php` | Envío de cotización al cliente | `emails.quote` |
| `PayStub` | `app/Mail/PayStub.php` | Envío de recibo de pago al empleado | `emails.paystub` |
| `TaskNotification` | `app/Mail/TaskNotification.php` | Notificación de nueva tarea asignada | `emails.task-notification` |

### Endpoints que envían email (Laravel)

| Ruta | Controlador | Método |
|------|------------|--------|
| `POST jobs/{job}/dispatch` | `JobController` | `sendDispatchEmail` |
| `POST jobs/{job}/dispatch-supervisor` | `JobController` | `sendDispatchSupervisorEmail` |
| `POST invoices/{invoice}/email` | `InvoiceController` | `email` |
| `POST quotes/{quote}/email` | `QuoteController` | `email` |
| `POST payroll/{user}/paystub/email` | `PayrollController` | `emailPaystub` |
| `POST job/{job}/tasks/{task}/send-notification` | `JobTaskController` | `sendNotification` |

### Estado en Go
- **✅ MIGRADO**: CRUD de registros `job_emails` (crear, listar, eliminar registros en BD)
- **✅ MIGRADO**: La lógica real de envío de email via Mailgun/SMTP
- **✅ MIGRADO**: Templates de email (HTML/Markdown)

### Recomendación de implementación en Go
1. Crear paquete `pkg/infrastructure/email/` con interfaz `EmailSender`
2. Implementar `MailgunSender` usando el SDK de Mailgun para Go (`github.com/mailgun/mailgun-go/v4`)
3. Crear templates HTML en `templates/emails/`
4. Agregar endpoints de envío en los handlers existentes de `job`, `invoice`, `quote`, `job_task`
5. Registrar el envío en `job_emails` tras cada envío exitoso

### Variables de entorno necesarias
```env
MAIL_FROM_ADDRESS=info@wecoolatlanta.com
MAIL_FROM_NAME="WeCool Atlanta"
MAILGUN_DOMAIN=wecoolatlanta.com
MAILGUN_SECRET=<secret>
```

---

## 2. ✅ Envío de SMS (AWS SNS + Twilio) - COMPLETADO

### Estado en Laravel
Dos servicios de SMS coexisten:

| Servicio | Archivo | Proveedor | Configuración |
|----------|---------|-----------|---------------|
| `AWSSMSService` | `app/Services/Communications/AWSSMSService.php` | AWS SNS | Variables de entorno `AWS_*` |
| `TwilioService` | `app/Services/Communications/TwilioService.php` | Twilio | Tabla `settings` (campos `is_twilio_enabled`, `twilio_sid`, `twilio_auth_token`, `twilio_from_number`) |

**Endpoint principal**: `POST jobs/{job}/dispatch-sms` → `JobController@sendDispatchSMS`

**Flujo**: El controlador usa `AWSSMSService` para enviar SMS y registra el envío en `job_sms` + `job_activity_logs`.

### Estado en Go
- **✅ MIGRADO**: CRUD de registros `job_sms` (crear, listar, eliminar registros en BD)
- **✅ MIGRADO**: La tabla `settings` ya incluye campos Twilio
- **✅ MIGRADO**: La lógica real de envío de SMS via AWS SNS o Twilio

### Recomendación de implementación en Go
1. Crear paquete `pkg/infrastructure/sms/` con interfaz `SMSSender`
2. Implementar `AWSSNSSender` usando AWS SDK para Go (`github.com/aws/aws-sdk-go-v2/service/sns`)
3. Opcionalmente implementar `TwilioSender` usando `github.com/twilio/twilio-go`
4. Agregar endpoint `POST /api/v1/jobs/{jobId}/dispatch-sms` en el handler de jobs
5. Leer configuración de Twilio desde la tabla `settings` (ya tiene los campos)
6. Registrar el envío en `job_sms` y `job_activity_logs`

### Variables de entorno necesarias
```env
# AWS SNS (las credenciales AWS ya existen para S3, se reutilizan)
AWS_DEFAULT_REGION=us-east-1  # Ya existe como S3_REGION

# Twilio (se lee de la tabla settings, no de env)
# No requiere variables de entorno adicionales
```

---

## 3. ✅ Seguridad de Contraseñas - COMPLETADO

### Estado en Laravel
Sistema completo de gestión de contraseñas:

| Funcionalidad | Controlador | Descripción |
|---------------|------------|-------------|
| Forgot Password | `ForgotPasswordController` | Envía link de reset por email |
| Reset Password | `ResetPasswordController` | Resetea contraseña con token |
| Enforce Reset | `PasswordSecurityController` | Fuerza cambio si expiró |
| Password History | Modelo `PasswordHistory` | Previene reutilización de contraseñas anteriores |

**Políticas configurables** (tabla `settings`):
- `is_enforce_routine_password_reset` — Activar/desactivar rotación obligatoria
- `password_expire_days` — Días hasta expiración (default: 30)
- `password_history_count` — Cantidad de contraseñas a recordar (default: 10)
- `password_minimum_length` — Longitud mínima (default: 8)
- `password_age` — Cantidad de contraseñas anteriores a verificar (default: 5)
- `password_include_numbers` — Requerir números (default: true)
- `password_include_symbols` — Requerir símbolos (default: true)

**Tablas involucradas**:
- `password_history` — Historial de contraseñas por usuario
- `password_resets` — Tokens temporales de reset

### Estado en Go
- **✅ MIGRADO**: Login, Logout, Refresh Token (JWT)
- **✅ MIGRADO**: Tabla `settings` con campos de políticas de contraseña
- **✅ MIGRADO**: Forgot password (envío de email con token)
- **✅ MIGRADO**: Reset password (validar token + cambiar contraseña)
- **✅ MIGRADO**: Enforce password reset (verificación en login)
- **✅ MIGRADO**: Password history (modelo, repositorio, validación)

### Recomendación de implementación en Go
1. Crear dominio `pkg/domain/password_history/` (entity, repository, usecase)
2. Crear dominio `pkg/domain/password_reset/` (entity, repository, usecase)
3. Agregar endpoints en auth routes:
   - `POST /auth/forgot-password` — Envía email con token
   - `POST /auth/reset-password` — Resetea con token
   - `POST /auth/change-password` — Cambio forzado
4. Modificar `Login` en auth usecase para verificar expiración de contraseña
5. Validar contra `password_history` en todos los flujos de cambio de contraseña
6. **Depende de**: Servicio de email (Sección 1)

---

## 4. ✅ Dashboard - COMPLETADO

### Estado en Laravel
Controlador `DashboardController` con dos vistas:

**Vista Admin** (`dashboard`):
- Jobs pendientes de dispatch (`dispatched = false`, `closed = false`)
- Ordenados por fecha de recepción

**Vista Técnico** (`technician-dashboard`):
- Jobs dispatched asignados al usuario actual
- Jobs urgentes (prioridad alta)

### Estado en Go
- **✅ MIGRADO**: Endpoint `GET /api/v1/dashboard` implementado
  - Detecta automáticamente el rol del usuario (admin vs técnico) desde el JWT
  - **Admin**: Jobs pendientes de dispatch + urgentes + estadísticas generales
  - **Técnico**: Jobs dispatched asignados + urgentes + estadísticas personales
  - Jobs enriquecidos con JOINs (property, customer, technician, category, priority, status)
  - Parámetro `?limit=N` para controlar cantidad de jobs por sección

### ✅ Archivos creados
- `pkg/domain/dashboard/entity.go` — Entidades del dominio
- `pkg/domain/dashboard/repository.go` — Interfaz del repositorio
- `pkg/domain/dashboard/usecase.go` — Lógica de negocio
- `pkg/repository/mysql/dashboard/repository.go` — Implementación MySQL
- `pkg/rest/handler/dashboard/handler.go` — Handler HTTP

### Propuesta de mejora
- Ver `docs/DASHBOARD_ENHANCEMENT_PROPOSAL.md` para plan de enriquecimiento con datos de invoices, quotes, tasks, warranties, alerts, etc.

---

## 5. ✅ Payroll - COMPLETADO

### Estado en Laravel
Controlador `PayrollController` con las siguientes funcionalidades:

| Método | Ruta | Descripción |
|--------|------|-------------|
| `index` | `GET payroll` | Lista de payroll (basada en `job_rates`) |
| `markPaid` | `PUT payroll/{user}/mark-paid` | Marcar rates como pagados |
| `markHeld` | `PUT payroll/{user}/mark-held` | Marcar rates como retenidos |
| `payStub` | `GET payroll/{user}/paystub` | Ver recibo de pago |
| `emailPaystub` | `POST payroll/{user}/paystub/email` | Enviar recibo por email |

**Nota**: Payroll se basa en las tablas `job_rates` y `job_rate_statuses`, que ya están migradas como CRUD.

### Estado en Go
- **✅ MIGRADO**: CRUD de `job_rates` y `job_rate_statuses`
- **✅ MIGRADO**: Lógica de negocio completa de payroll (mark paid/held, paystub view, paystub email)

### ✅ Implementación completada
1. ✅ **Dominio**: `pkg/domain/payroll/` con entities, repository y usecase
2. ✅ **Repositorio MySQL**: `pkg/repository/mysql/payroll/` con implementación completa
3. ✅ **Handler HTTP**: `pkg/rest/handler/payroll/` con todos los endpoints
4. ✅ **Endpoints implementados**:
   - `GET /api/v1/payroll` — Lista agrupada por usuario
   - `GET /api/v1/payroll/{userId}` — Rates de usuario específico
   - `PUT /api/v1/payroll/{userId}/mark-paid` — Marcar como pagado
   - `PUT /api/v1/payroll/{userId}/mark-held` — Marcar como retenido
   - `GET /api/v1/payroll/{userId}/paystub` — Ver recibo
   - `POST /api/v1/payroll/{userId}/paystub/email` — Enviar recibo por email
5. ✅ **Integración**: Servicio de email integrado para envío de paystub
6. ✅ **Documentación**: Swagger regenerado con todos los endpoints

### ✅ Archivos creados
- `pkg/domain/payroll/entity.go` — Entidades del dominio
- `pkg/domain/payroll/repository.go` — Interface del repositorio
- `pkg/domain/payroll/usecase.go` — Lógica de negocio
- `pkg/repository/mysql/payroll/repository.go` — Implementación MySQL
- `pkg/repository/mysql/payroll/list_payroll_users.go` — Listado de usuarios con rates
- `pkg/repository/mysql/payroll/get_user_rates.go` — Rates por usuario
- `pkg/repository/mysql/payroll/mark_rates.go` — Marcar estados
- `pkg/repository/mysql/payroll/get_paystub.go` — Generación de paystub
- `pkg/repository/mysql/payroll/adapters.go` — Validaciones
- `pkg/rest/handler/payroll/handler.go` — Handler HTTP con 6 endpoints

---

## 6. ✅ Búsqueda Global (Search) - COMPLETADO

### Estado en Laravel
Controlador `SearchController` que busca en múltiples entidades simultáneamente:

**Entidades buscadas**:
- Jobs (por work_order, notes, quick_notes)
- Customers (por name, email, phone)
- Properties (por street, city, zip)
- Invoices (por invoice_number)
- Quotes (por quote_number)
- Warranties (por warranty_number)
- Warranty Claims (por internal_claim_number, claim_number)
- Users (por name, email)

### Estado en Go
- **✅ MIGRADO**: Endpoint de búsqueda unificada cross-entity
- **✅ MIGRADO**: Búsqueda en paralelo usando goroutines
- **✅ MIGRADO**: Respeta permiso `job_view_user_only` para filtrar resultados

### ✅ Implementación completada
1. **Dominio**: `pkg/domain/search/` con entities, repository interface y usecase
2. **Repositorio MySQL**: `pkg/repository/mysql/search/` con queries optimizadas
3. **Handler HTTP**: `pkg/rest/handler/search/` con endpoint documentado
4. **Endpoint**: `GET /api/v1/search?q={query}&limit={limit}`

### ✅ Archivos creados
- `pkg/domain/search/entity.go` — Entidades de respuesta por tipo
- `pkg/domain/search/repository.go` — Interface del repositorio
- `pkg/domain/search/usecase.go` — Lógica de negocio con goroutines paralelas
- `pkg/repository/mysql/search/repository.go` — Implementación MySQL con 8 métodos de búsqueda
- `pkg/rest/handler/search/handler.go` — Handler HTTP con documentación Swagger

### Características
- Búsqueda en 8 entidades simultáneamente (jobs, customers, properties, invoices, quotes, warranties, warranty_claims, users)
- Ejecución paralela con goroutines para máximo rendimiento
- Límite configurable por entidad (default: 10, max: 50)
- Respeta el permiso `job_view_user_only` para técnicos
- Resultados agrupados por tipo de entidad con campos relevantes

---

## 7. Exportación a Excel

### Estado en Laravel
Clase `OpenJobsExport` usando `Maatwebsite/Excel`:

**Datos exportados**: Jobs abiertos (`closed = false`) con 18 columnas (fecha, dirección, W/O, precio, cliente, tipo, técnico, estado, prioridad, etc.)

**Ruta**: `GET jobs/export` → `JobController@export`

### Estado en Go
- **NO migrado**: No existe funcionalidad de exportación

### Recomendación de implementación en Go
1. Usar librería `github.com/xuri/excelize/v2` para generación de Excel
2. Agregar endpoint `GET /api/v1/jobs/export` en el handler de jobs
3. Reutilizar el repositorio de `job` con filtro `closed = false`
4. Generar y devolver archivo `.xlsx` como response con headers adecuados

---

## 8. Integración con Stripe (Pagos Online)

### Estado en Laravel

**Pago de facturas online** (`InvoiceController`):
- `GET invoices/{invoice}/pay` — Página de pago (pública, sin auth)
- Crea `PaymentIntent` via Stripe SDK
- Usa `STRIPE_SECRET_KEY` y `STRIPE_PUBLIC_KEY`

**Webhook** (`StripeWebhookController`):
- `POST stripe/webhook` — Recibe eventos de Stripe
- Procesa `payment_intent.succeeded`
- Crea registro en `invoice_payments` con datos del pago

### Estado en Go
- **Migrado**: CRUD de `invoice_payments`
- **NO migrado**: Creación de PaymentIntent
- **NO migrado**: Webhook de Stripe
- **NO migrado**: Página/endpoint público de pago

### Recomendación de implementación en Go
1. Crear paquete `pkg/infrastructure/stripe/`
2. Usar SDK oficial: `github.com/stripe/stripe-go/v76`
3. Endpoints:
   - `GET /api/v1/invoices/{id}/payment-intent` — Crear PaymentIntent (público)
   - `POST /webhooks/stripe` — Procesar eventos de Stripe (público, verificar firma)
4. Registrar pagos en `invoice_payments` via repositorio existente

### Variables de entorno necesarias
```env
STRIPE_SECRET_KEY=<secret>
STRIPE_PUBLIC_KEY=<public>
STRIPE_WEBHOOK_SECRET=<webhook_secret>
```

---

## 9. ✅ Gestión de Cuenta (Self-Service) - COMPLETADO

### Estado en Laravel
Controlador `AccountController`:
- `GET account` — Ver perfil
- `PUT account` — Actualizar nombre/email
- `POST account/sidebar` — Toggle sidebar (UI-specific)

### Estado en Go
- **✅ MIGRADO**: Endpoints de gestión de cuenta implementados
- **✅ MIGRADO**: Validación de contraseñas con historial y políticas de seguridad
- **✅ MIGRADO**: Validación de email duplicado

### ✅ Implementación completada
1. **Dominio**: `pkg/domain/account/` con entities y usecase
2. **Handler HTTP**: `pkg/rest/handler/account/` con endpoints documentados
3. **Endpoints implementados**:
   - `GET /api/v1/account` — Obtener perfil del usuario autenticado
   - `PUT /api/v1/account` — Actualizar nombre y email del usuario
   - `PUT /api/v1/account/password` — Cambiar contraseña con validaciones de seguridad
4. **Validaciones implementadas**:
   - Verificación de contraseña actual
   - Validación de requisitos de seguridad (longitud, números, símbolos)
   - Verificación contra historial de contraseñas
   - Validación de email duplicado
5. **Integración**: Usa repositorios existentes de `user`, `password_history` y `settings`
6. **Documentación**: Swagger actualizado con todos los endpoints
7. **Testing**: Colección Postman creada en `docs/postman_account_collection.json`

### ✅ Archivos creados
- `pkg/domain/account/entity.go` — DTOs de request/response
- `pkg/domain/account/usecase.go` — Lógica de negocio con validaciones
- `pkg/rest/handler/account/handler.go` — Handler HTTP con 3 endpoints

### Nota
El endpoint `POST account/sidebar` de Laravel no fue migrado ya que es específico de la UI del frontend Laravel y no aplica para una API REST stateless.

---

## 10. Análisis de Variables de Entorno

### Variables del JSON de Producción (`WeCoolAtlanta-Prod-revision74.json`)

#### Variables que NECESITAN migración

| Variable | Valor Producción | Uso | Prioridad |
|----------|-----------------|-----|-----------|
| `MAIL_FROM_ADDRESS` | `info@wecoolatlanta.com` | Remitente de emails | Alta |
| `MAIL_FROM_NAME` | `WeCool Atlanta` | Nombre del remitente | Alta |
| `MAILGUN_DOMAIN` | `wecoolatlanta.com` | Dominio Mailgun | Alta |
| `MAILGUN_SECRET` | `key-***` | API Key Mailgun | Alta |
| `STRIPE_SECRET_KEY` | `sk_live_***` | API Stripe server-side | Media |
| `STRIPE_PUBLIC_KEY` | `pk_live_***` | API Stripe client-side | Media |
| `MAIL_HOST` | `smtp.mailgun.org` | Fallback SMTP host | Baja |
| `MAIL_PORT` | `587` | Fallback SMTP port | Baja |
| `MAIL_USERNAME` | `postmaster@...` | Fallback SMTP user | Baja |
| `MAIL_PASSWORD` | `***` | Fallback SMTP pass | Baja |

#### Variables que NO necesitan migración

| Variable | Razón |
|----------|-------|
| `APP_KEY` | Cifrado específico de Laravel. Go usa JWT con su propia secret. |
| `APP_DEBUG` | Go tiene su propio manejo de debug/logging. |
| `SESSION_DRIVER=redis` | Go usa JWT stateless, no necesita sesiones server-side. |
| `SESSION_LIFETIME=120` | No aplica con JWT. |
| `CACHE_DRIVER=redis` | Solo necesario si se implementa caching en Go (futuro). |
| `BROADCAST_DRIVER=log` | No se usa activamente en producción. |
| `QUEUE_CONNECTION=sync` | Laravel ejecuta todo sincrónicamente. Go usa goroutines nativamente. |
| `LOG_CHANNEL=stack` | Go tiene su propio sistema de logging. |
| `LOG_LEVEL=debug` | Se configura independientemente en Go. |
| `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` | Solo si se implementa caching/rate-limiting en Go. |
| `DB_CONNECTION=mysql` | Go ya tiene `DB_DRIVER`. |
| `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD` | Ya migradas en `configs/app.env`. |
| `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION`, `AWS_BUCKET` | Ya migradas como `S3_*` en `configs/app.env`. |
| `FILESYSTEM_DRIVER=s3` | Go ya usa S3 directamente. |

#### Variables a migrar en `configs/app.env`

```env
# === Email (Mailgun) ===
MAIL_FROM_ADDRESS=info@wecoolatlanta.com
MAIL_FROM_NAME=WeCool Atlanta
MAILGUN_DOMAIN=wecoolatlanta.com
MAILGUN_SECRET=

# === Stripe ===
STRIPE_SECRET_KEY=
STRIPE_PUBLIC_KEY=
STRIPE_WEBHOOK_SECRET=
```

#### Cambios necesarios en `configs/config.go`

Agregar structs:
```go
type MailConfig struct {
    FromAddress   string `mapstructure:"MAIL_FROM_ADDRESS"`
    FromName      string `mapstructure:"MAIL_FROM_NAME"`
    MailgunDomain string `mapstructure:"MAILGUN_DOMAIN"`
    MailgunSecret string `mapstructure:"MAILGUN_SECRET"`
}

type StripeConfig struct {
    SecretKey     string `mapstructure:"STRIPE_SECRET_KEY"`
    PublicKey     string `mapstructure:"STRIPE_PUBLIC_KEY"`
    WebhookSecret string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
}
```

---

## Priorización Recomendada

### Fase 2A — Comunicaciones (Prioridad Alta)
1. **Servicio de Email** (Mailgun) — Base para múltiples funcionalidades
2. **Envío de SMS** (AWS SNS) — Dispatch de técnicos
3. **Seguridad de Contraseñas** — Forgot/Reset password (depende de email)

### Fase 2B — Funcionalidades de Negocio (Prioridad Media)
4. **Dashboard** — Vista principal de la aplicación
5. **Payroll** — Gestión de pagos a técnicos
6. **Gestión de Cuenta** — Self-service para usuarios

### Fase 2C — Integraciones y Extras (Prioridad Baja)
7. **Stripe Integration** — Pagos online de facturas
8. **Búsqueda Global** — Búsqueda unificada cross-entity
9. **Exportación Excel** — Reportes de jobs abiertos

---

## Dependencias entre Funcionalidades

```
Servicio de Email (1)
├── Seguridad de Contraseñas (3) — necesita enviar email de reset
├── Payroll (5) — necesita enviar paystub por email
└── Envío de Email en Jobs/Invoices/Quotes/Tasks

Servicio de SMS (2)
└── Dispatch SMS en Jobs

Dashboard (4) — independiente, usa repos existentes

Stripe (8) — independiente, pero necesita endpoint público

Search (6) — independiente, usa repos existentes

Export Excel (7) — independiente, usa repo de jobs
```

---

## Notas Adicionales

- **Rutas comentadas en Laravel**: Las rutas CRUD de invoices y quotes están comentadas en `routes/web.php` (líneas 324-353), pero los endpoints de **email** de invoice/quote sí estaban activos. Go ya tiene los CRUDs.
- **Twilio vs AWS SNS**: Laravel tiene ambos servicios pero el endpoint de dispatch SMS usa `AWSSMSService`. Las credenciales de Twilio se guardan en la tabla `settings` por si se quiere activar. Se recomienda migrar primero AWS SNS y dejar Twilio como segunda opción configurable.
- **Templates de email**: Los templates Blade de Laravel (`resources/views/emails/`) necesitan ser convertidos a templates HTML para Go. Se recomienda usar `html/template` de la librería estándar.
