# Template de Prompt para Migración de Módulos JVAIR → JVAIRV2

## Instrucciones de Uso

Copia el template de abajo, reemplaza los valores entre `{{...}}` con los datos del módulo que vas a migrar, y envíalo como prompt.

---

## Template

```
Necesito migrar el módulo {{NUMERO_MODULO}}: {{NOMBRE_MODULO}} del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/{{Modelo}}.php`
  - Controlador(es): `app/Http/Controllers/{{Controlador}}Controller.php`
  - Request(s): `app/Http/Requests/{{Modelo}}/StoreRequest.php`, `UpdateRequest.php`
  - Rutas: `routes/api.php` (buscar rutas de {{recurso}})

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql` — SIEMPRE confiar en este archivo para la estructura real de las tablas
- **Plan de migración**: `docs/MIGRATION_PLAN.md` — referencia de dependencias entre módulos
- **Datos de producción**: `scripts/Data_prd.sql` — datos reales para referencia

### Tablas involucradas
{{LISTA_DE_TABLAS}}

### Dependencias con módulos ya migrados
{{LISTA_DE_DEPENDENCIAS}}

## Arquitectura y Estándares del Proyecto JVAIRV2

### Estructura de archivos por módulo
Cada módulo sigue esta estructura estricta:

```
pkg/domain/{{modulo}}/
├── entity.go          # Struct de la entidad, Filters, validaciones
├── errors.go          # Errores del dominio (var ErrXxx = errors.New(...))
├── repository.go      # Interface Repository
├── usecase.go         # UseCase struct, Service interface, NewUseCase(), checker interfaces
├── create.go          # Método Create del UseCase
├── get_by_id.go       # Método GetByID del UseCase
├── list.go            # Método List del UseCase
├── update.go          # Método Update del UseCase
├── delete.go          # Método Delete del UseCase (soft delete)
├── {{extra}}.go       # Métodos adicionales (ej: close.go, duplicate.go)
├── mock.go            # Mocks para tests (MockRepository, MockCheckers)
└── usecase_test.go    # Tests unitarios del UseCase

pkg/repository/mysql/{{modulo}}/
├── repository.go      # Repository struct, NewRepository(), *sql.DB
├── create.go          # Implementación SQL de Create
├── get_by_id.go       # Implementación SQL de GetByID
├── list.go            # Implementación SQL de List (paginación, filtros, sorting)
├── update.go          # Implementación SQL de Update
├── delete.go          # Implementación SQL de Delete (soft delete)
├── {{extra}}.go       # Métodos adicionales
└── adapters.go        # (si aplica) Adapters/Checkers para validar FKs de otros módulos

pkg/rest/handler/{{modulo}}/
├── handler.go         # Handler struct, NewHandler(), RegisterRoutes(), DTOs (Request/Response), helpers
├── create.go          # Handler HTTP Create con Swagger annotations
├── get.go             # Handler HTTP Get con Swagger annotations
├── list.go            # Handler HTTP List con Swagger annotations
├── update.go          # Handler HTTP Update con Swagger annotations
├── delete.go          # Handler HTTP Delete con Swagger annotations
└── {{extra}}.go       # Handlers adicionales
```

### Patrones obligatorios
1. **Soft deletes**: Todas las entidades usan `deleted_at` (timestamp NULL). DELETE = UPDATE SET deleted_at = NOW(). Todas las queries filtran `deleted_at IS NULL`.
2. **Paginación**: List retorna `([]Entity, int64, error)` donde int64 es el total. Query params: `page` (default 1), `limit` (default 15).
3. **Filtros**: Se parsean de query params en el handler y se pasan como `map[string]interface{}` al repository.
4. **Errores de dominio**: Variables `var` en `errors.go`, se mapean a HTTP status codes en el handler.
5. **Validación de FKs**: El UseCase valida existencia de entidades relacionadas usando interfaces Checker (ej: `JobCategoryChecker`, `PropertyChecker`). Los Checkers se implementan como Adapters en el repository del módulo que los consume.
6. **Imports del proyecto**: `github.com/your-org/jvairv2/...`
7. **Router**: Chi (`github.com/go-chi/chi/v5`)
8. **Response helper**: `github.com/your-org/jvairv2/pkg/rest/response` — usar `response.JSON(w, status, data)` y `response.Error(w, status, msg)`
9. **Logging**: `log/slog` con `slog.ErrorContext`, `slog.InfoContext`, `slog.WarnContext`
10. **Tests**: `testify/assert` y `testify/mock`. Mockear Repository y todos los Checkers.
11. **Swagger**: Annotations en cada handler method (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`). Al finalizar el módulo, regenerar docs con `swag init`

### Integración (después de crear el módulo)
1. **`cmd/api/container.go`**: Agregar imports, inicializar repository, checkers/adapters, use case y handler. Pasar handler al router.
2. **`pkg/rest/router/router.go`**: Agregar import del handler, parámetro en `New()`, y registrar rutas con `handler.RegisterRoutes(r)`.

## Entregables Requeridos

### 1. Análisis previo
- Crear `docs/POINT_{{NUMERO}}_{{NOMBRE}}_ANALYSIS.md` con:
  - Estructura de tabla(s) desde `db_structure.sql`
  - Análisis de funcionalidad original (modelo, controlador, validaciones, relaciones)
  - Plan de ejecución por fases
  - Decisiones de diseño (qué se migra y qué NO se migra en este punto)
  - Filtros soportados en List

### 2. Código del módulo
- Todos los archivos listados en la estructura de arriba
- Seguir exactamente los patrones de módulos existentes (usar `customer` o `property` como referencia para módulos simples, `job` para módulos con dependencias complejas)

### 3. Tests unitarios
- `usecase_test.go` con cobertura de:
  - Caso exitoso para cada método
  - Errores de validación
  - Errores de FKs inválidas
  - Errores de not found
  - Errores de repository

### 4. Documentación Swagger
- Annotations completas en cada handler (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)
- **Al finalizar la migración del módulo**, ejecutar `swag init -g cmd/api/main.go -o docs/swagger` para regenerar la documentación Swagger completa con los nuevos endpoints
- Verificar que la documentación generada incluya correctamente todos los endpoints del módulo recién migrado
- Confirmar que los endpoints previos no se hayan perdido en la regeneración

### 5. Colección Postman
- Crear `docs/postman_{{modulo}}_collection.json` **siguiendo el mismo estándar y formato de las colecciones existentes** (ver como referencia: `docs/postman_customer_collection.json`, `docs/postman_jobs_collection.json`, u otra colección ya generada)
- La colección debe incluir:
  - Todas las operaciones CRUD del módulo
  - Variables: `{{baseUrl}}`, `{{accessToken}}`
  - Body de ejemplo con datos realistas basados en `scripts/Data_prd.sql`
  - Descripción clara en cada request
  - Headers correctos (Content-Type, Authorization Bearer)
  - Ejemplos de query params para filtros y paginación en los endpoints List
- Si el módulo tiene sub-recursos o catálogos, incluir requests para cada uno
- La colección debe estar lista para importar en Postman y ejecutar pruebas inmediatamente

### 6. Verificación final
- Ejecutar `go build ./...` para verificar compilación
- Ejecutar `go test ./pkg/domain/{{modulo}}/...` para verificar tests
- Ejecutar pre-commit hooks si están configurados

## Operaciones CRUD esperadas
{{LISTA_DE_OPERACIONES}}

## Notas adicionales
{{NOTAS_ADICIONALES}}
```

---

## Ejemplo de uso completo (Módulo 9: Cotizaciones)

```
Necesito migrar el módulo 9: Cotizaciones (Quotes) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/Quote.php`, `app/Models/QuoteStatus.php`
  - Controlador(es): `app/Http/Controllers/QuoteController.php`
  - Request(s): `app/Http/Requests/Quote/StoreRequest.php`, `UpdateRequest.php`
  - Rutas: `routes/api.php` (buscar rutas de quotes)

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql`
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tablas involucradas
- `quotes` — tabla principal
- `quote_statuses` — catálogo de estados de cotización

### Dependencias con módulos ya migrados
- `jobs` (módulo 8) — quotes pertenecen a un job
- `users` (módulo 1) — usuario que crea la cotización
- `properties` (módulo 6) — propiedad asociada al job

## Operaciones CRUD esperadas
- **List**: GET /api/v1/quotes — con filtros por jobId, status, paginación
- **Create**: POST /api/v1/quotes — crear cotización para un job
- **Get**: GET /api/v1/quotes/{id} — obtener cotización por ID
- **Update**: PUT /api/v1/quotes/{id} — actualizar cotización
- **Delete**: DELETE /api/v1/quotes/{id} — soft delete
- **QuoteStatuses CRUD**: GET/POST/PUT/DELETE /api/v1/quote-statuses (catálogo)

## Notas adicionales
- Las cotizaciones tienen un monto (amount) y un status
- Revisar si hay lógica de aprobación/rechazo en el controlador original
- El catálogo quote_statuses ya tiene datos en Data_prd.sql

## Entregables de documentación al finalizar

### Actualización de Swagger
- Asegurar annotations completas en cada handler (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)
- Ejecutar `swag init -g cmd/api/main.go -o docs/swagger` para regenerar la documentación Swagger con los nuevos endpoints de quotes y quote-statuses
- Verificar que los endpoints nuevos aparezcan correctamente y que los endpoints previos no se hayan perdido

### Colección Postman
- Crear `docs/postman_quotes_collection.json` y `docs/postman_quote_statuses_collection.json` siguiendo el mismo estándar y formato de las colecciones existentes (ver como referencia: `docs/postman_customer_collection.json`, `docs/postman_jobs_collection.json`)
- Incluir todas las operaciones CRUD, variables `{{baseUrl}}` y `{{accessToken}}`, body de ejemplo con datos realistas de `Data_prd.sql`, headers correctos, y ejemplos de query params para filtros y paginación
- Las colecciones deben estar listas para importar en Postman y ejecutar pruebas inmediatamente
```

---

## Ejemplo de uso completo (Módulo 5: Supervisores)

```
Necesito migrar el módulo 5: Supervisores (Supervisors) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/Supervisor.php`
  - Controlador(es): `app/Http/Controllers/SupervisorController.php`
  - Request(s): `app/Http/Requests/Supervisors/StoreRequest.php`, `UpdateRequest.php`, `IndexRequest.php`, `EditRequest.php`, `CreateRequest.php`, `DeleteRequest.php`
  - Policy: `app/Policies/SupervisorPolicy.php`
  - Rutas: `routes/web.php` (buscar rutas de supervisors)
  - Relación con Customer: `app/Http/Controllers/CustomerController.php` (método `supervisors()`)
  - Relación con Job: `app/Models/Job.php` (campo `supervisor_ids`, método `supervisorEmails()`)
  - Dispatch a supervisores: `app/Http/Controllers/JobController.php` (método `sendDispatchSupervisorEmail()`)

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql`
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tablas involucradas
- `supervisors` — tabla principal de supervisores (contactos de clientes)

### Dependencias con módulos ya migrados
- `customers` (módulo 4) — cada supervisor pertenece a un customer (customer_id, FK obligatoria)

### Módulos que dependen de supervisors
- `jobs` (módulo 8) — los jobs referencian supervisores mediante el campo CSV `supervisor_ids` en la tabla `jobs`. El método `supervisorEmails()` en Job resuelve esos IDs para obtener emails.

## Operaciones CRUD esperadas

### Supervisors
- **List**: GET /api/v1/supervisors — con filtros por customerId, search (name, phone, email), paginación. Aplicar scope `forCurrentUser()` (si el usuario tiene permiso `job_view_user_only`, solo ve supervisores de clientes que tengan jobs asignados a él)
- **Create**: POST /api/v1/supervisors — crear supervisor asociado a un customer
- **Get**: GET /api/v1/supervisors/{id} — obtener supervisor por ID (incluir datos del customer)
- **Update**: PUT /api/v1/supervisors/{id} — actualizar supervisor
- **Delete**: DELETE /api/v1/supervisors/{id} — soft delete

### Supervisors por Customer (vista filtrada)
- **List**: GET /api/v1/customers/{customerId}/supervisors — listar supervisores de un cliente específico (con mismos filtros de search y paginación)

## Notas adicionales
- Entidad simple con campos: customer_id (FK requerida), name (requerido, string), phone (opcional), email (opcional)
- Usa SoftDeletes (campo `deleted_at`)
- Tiene scope `scopeForCurrentUser()` que restringe visibilidad: si el usuario tiene el permiso `job_view_user_only`, solo ve supervisores de clientes que tengan jobs donde `user_id` sea el usuario actual. Implementar esta lógica como filtro en el UseCase/Handler
- Searchable por: name, phone, email
- Permisos (Policy): `supervisor_view` para ver, `supervisor_edit` para crear/editar, `supervisor_delete` para eliminar. Nota: `viewAny()` tiene `return true` hardcodeado (bug en Laravel, cualquier usuario autenticado puede ver el listado)
- No se puede editar un supervisor que esté soft-deleted (validación en Policy)
- La relación Job→Supervisor es **indirecta**: el campo `jobs.supervisor_ids` almacena IDs como string CSV (ej: "1,5,12"). En Go, evaluar si:
  - (a) Normalizar con tabla pivot `job_supervisors` (mejor diseño) — requiere migración de datos
  - (b) Mantener como string CSV por compatibilidad con datos existentes
- El método `Job.supervisorEmails()` parsea `supervisor_ids`, busca los supervisores y retorna sus emails concatenados. Se usa en el flujo de dispatch email a supervisores. NO migrar esta lógica ahora, se integrará cuando se trabaje el módulo de comunicaciones
- Hay 173 registros en producción
- El controlador carga la lista de customers para el dropdown de asignación al crear/editar

## Entregables de documentación al finalizar

### Actualización de Swagger
- Asegurar annotations completas en cada handler (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)
- Ejecutar `swag init -g cmd/api/main.go -o docs/swagger` para regenerar la documentación Swagger con los nuevos endpoints de supervisors
- Verificar que los endpoints nuevos aparezcan correctamente y que los endpoints previos no se hayan perdido

### Colección Postman
- Crear `docs/postman_supervisors_collection.json` siguiendo el mismo estándar y formato de las colecciones existentes (ver como referencia: `docs/postman_customer_collection.json`, `docs/postman_jobs_collection.json`)
- Incluir todas las operaciones CRUD (supervisors y supervisors por customer), variables `{{baseUrl}}` y `{{accessToken}}`, body de ejemplo con datos realistas de `Data_prd.sql`, headers correctos, y ejemplos de query params para filtros y paginación
- La colección debe estar lista para importar en Postman y ejecutar pruebas inmediatamente
```

---

## Ejemplo de uso completo (Módulo 10: Facturas)

```
Necesito migrar el módulo 10: Facturas (Invoices) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/Invoice.php`, `app/Models/InvoicePayment.php`
  - Controlador(es): `app/Http/Controllers/InvoiceController.php`, `app/Http/Controllers/InvoicePaymentController.php`
  - Request(s): `app/Http/Requests/Invoices/StoreRequest.php`, `UpdateRequest.php`, `EmailRequest.php`, `PayRequest.php`
  - Request(s): `app/Http/Requests/InvoicePayments/StoreRequest.php`, `UpdateRequest.php`
  - Rutas: `routes/web.php` (buscar rutas de invoices — actualmente comentadas pero con estructura definida)

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql`
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tablas involucradas
- `invoices` — tabla principal de facturas
- `invoice_payments` — pagos asociados a facturas

### Dependencias con módulos ya migrados
- `jobs` (módulo 8) — cada factura pertenece a un job (job_id)
- `users` (módulo 1) — usuario que crea/gestiona la factura
- `properties` (módulo 6) — propiedad asociada al job de la factura
- `customers` (módulo 4) — cliente asociado a la propiedad

## Operaciones CRUD esperadas

### Invoices
- **List**: GET /api/v1/invoices — con filtros por jobId, status (paid/unpaid), search, paginación
- **Create**: POST /api/v1/invoices — crear factura para un job
- **Get**: GET /api/v1/invoices/{id} — obtener factura por ID (incluir payments y balance calculado)
- **Update**: PUT /api/v1/invoices/{id} — actualizar factura
- **Delete**: DELETE /api/v1/invoices/{id} — soft delete

### Invoice Payments (sub-recurso de invoices)
- **List**: GET /api/v1/invoices/{invoiceId}/payments — listar pagos de una factura
- **Create**: POST /api/v1/invoices/{invoiceId}/payments — registrar pago
- **Get**: GET /api/v1/invoices/{invoiceId}/payments/{id} — obtener pago por ID
- **Update**: PUT /api/v1/invoices/{invoiceId}/payments/{id} — actualizar pago
- **Delete**: DELETE /api/v1/invoices/{invoiceId}/payments/{id} — soft delete

## Notas adicionales
- Las facturas tienen campos: invoice_number, total (decimal), description, allow_online_payments (boolean), notes
- El balance se calcula como: total - SUM(payments.amount). Incluir campo calculado `balance` en la respuesta
- Los pagos tienen: payment_processor, payment_id, amount (decimal), notes
- Hay lógica de filtrado por status: "paid" (balance <= 0) y "unpaid" (balance > 0) — implementar como filtro en List
- El método isPaid() verifica si balance == 0
- Hay integración con Stripe para pagos online (PayRequest, PaymentIntent) — NO migrar la integración con Stripe en este punto, solo el CRUD básico
- El método email() envía facturas por correo — NO migrar en este punto, se hará en el módulo 14
- Las rutas en Laravel están comentadas pero la estructura de controladores y requests está completa

## Entregables de documentación al finalizar

### Actualización de Swagger
- Asegurar annotations completas en cada handler (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)
- Ejecutar `swag init -g cmd/api/main.go -o docs/swagger` para regenerar la documentación Swagger con los nuevos endpoints de invoices e invoice payments
- Verificar que los endpoints nuevos aparezcan correctamente y que los endpoints previos no se hayan perdido

### Colección Postman
- Crear `docs/postman_invoices_collection.json` siguiendo el mismo estándar y formato de las colecciones existentes (ver como referencia: `docs/postman_customer_collection.json`, `docs/postman_jobs_collection.json`)
- Incluir todas las operaciones CRUD (invoices y invoice payments como sub-recurso), variables `{{baseUrl}}` y `{{accessToken}}`, body de ejemplo con datos realistas de `Data_prd.sql`, headers correctos, y ejemplos de query params para filtros y paginación
- La colección debe estar lista para importar en Postman y ejecutar pruebas inmediatamente
```

---

## Ejemplo de uso completo (Módulo 11: Equipos)

```
Necesito migrar el módulo 11: Equipos (Equipment) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/PropertyEquipment.php`, `app/Models/JobEquipment.php`
  - Controlador(es): `app/Http/Controllers/PropertyEquipmentController.php`, `app/Http/Controllers/JobEquipmentController.php`
  - Request(s): `app/Http/Requests/PropertyEquipment/StoreRequest.php`, `UpdateRequest.php`
  - Request(s): `app/Http/Requests/JobEquipment/StoreRequest.php`, `UpdateRequest.php`
  - Rutas: `routes/web.php` (buscar rutas de equipment dentro de properties y jobs)

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql`
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tablas involucradas
- `property_equipment` — equipos asociados a propiedades
- `job_equipment` — equipos asociados a trabajos

### Dependencias con módulos ya migrados
- `properties` (módulo 6) — property_equipment pertenece a una propiedad (property_id)
- `jobs` (módulo 8) — job_equipment pertenece a un job (job_id)

## Operaciones CRUD esperadas

### Property Equipment (sub-recurso de properties)
- **List**: GET /api/v1/properties/{propertyId}/equipment — listar equipos de una propiedad
- **Create**: POST /api/v1/properties/{propertyId}/equipment — crear equipo en propiedad
- **Get**: GET /api/v1/properties/{propertyId}/equipment/{id} — obtener equipo por ID
- **Update**: PUT /api/v1/properties/{propertyId}/equipment/{id} — actualizar equipo
- **Delete**: DELETE /api/v1/properties/{propertyId}/equipment/{id} — soft delete

### Job Equipment (sub-recurso de jobs)
- **List**: GET /api/v1/jobs/{jobId}/equipment — listar equipos de un job
- **Create**: POST /api/v1/jobs/{jobId}/equipment — crear equipo en job
- **Get**: GET /api/v1/jobs/{jobId}/equipment/{id} — obtener equipo por ID
- **Update**: PUT /api/v1/jobs/{jobId}/equipment/{id} — actualizar equipo
- **Delete**: DELETE /api/v1/jobs/{jobId}/equipment/{id} — soft delete

## Notas adicionales
- Ambas tablas comparten estructura similar de campos HVAC: area, outdoor_brand/model/serial/installed, furnace_brand/model/serial/installed, evaporator_brand/model/serial/installed, air_handler_brand/model/serial/installed
- Los campos *_installed son fechas (Date cast)
- job_equipment tiene un campo adicional `type` (valores: "current", "new") con scopes scopeCurrent() y scopeNew()
- property_equipment NO tiene soft deletes; job_equipment tampoco tiene SoftDeletes en el modelo pero verificar en db_structure.sql
- Existe lógica de clonación: `PropertyEquipment::cloneFromJobEquipment()` copia equipos de un job a la propiedad — evaluar si migrar como endpoint adicional o como lógica interna del UseCase
- PropertyEquipment tiene un accessor `getOutdoorUnitAttribute()` que concatena brand + model + serial — implementar como campo calculado en la respuesta si aplica
- Los equipos son sub-recursos anidados, NO tienen listado global independiente

## Entregables de documentación al finalizar

### Actualización de Swagger
- Asegurar annotations completas en cada handler (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)
- Ejecutar `swag init -g cmd/api/main.go -o docs/swagger` para regenerar la documentación Swagger con los nuevos endpoints de property-equipment y job-equipment
- Verificar que los endpoints nuevos aparezcan correctamente y que los endpoints previos no se hayan perdido

### Colección Postman
- Crear `docs/postman_property_equipment_collection.json` y `docs/postman_job_equipment_collection.json` siguiendo el mismo estándar y formato de las colecciones existentes (ver como referencia: `docs/postman_customer_collection.json`, `docs/postman_jobs_collection.json`)
- Incluir todas las operaciones CRUD de ambos sub-recursos, variables `{{baseUrl}}` y `{{accessToken}}`, body de ejemplo con datos realistas de `Data_prd.sql`, headers correctos
- Las colecciones deben estar listas para importar en Postman y ejecutar pruebas inmediatamente
```

---

## Ejemplo de uso completo (Módulo 12: Garantías)

```
Necesito migrar el módulo 12: Garantías (Warranties) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/Warranty.php`, `app/Models/WarrantyClaim.php`, `app/Models/WarrantyEquipment.php`, `app/Models/WarrantyClaimType.php`, `app/Models/WarrantyClaimStatus.php`, `app/Models/WarrantyStatus.php`, `app/Models/WarrantyType.php`
  - Controlador(es): `app/Http/Controllers/WarrantyController.php`, `app/Http/Controllers/WarrantyClaimController.php`, `app/Http/Controllers/WarrantyEquipmentController.php`, `app/Http/Controllers/WarrantyStatusController.php`, `app/Http/Controllers/WarrantyTypeController.php`, `app/Http/Controllers/WarrantyClaimStatusController.php`
  - Request(s): `app/Http/Requests/Warranties/StoreRequest.php`, `UpdateRequest.php`
  - Request(s): `app/Http/Requests/WarrantyClaims/StoreRequest.php`, `UpdateRequest.php`
  - Request(s): `app/Http/Requests/WarrantyEquipment/StoreRequest.php`, `UpdateRequest.php`
  - Request(s): `app/Http/Requests/WarrantyStatuses/`, `app/Http/Requests/WarrantyTypes/`, `app/Http/Requests/WarrantyClaimStatuses/`
  - Rutas: `routes/web.php` (buscar rutas de warranties, warranty-claims, y catálogos en settings)

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql`
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tablas involucradas
- `warranties` — tabla principal de garantías
- `warranty_claims` — reclamos de garantía
- `warranty_equipment` — equipos asociados a garantías
- `warranty_types` — catálogo de tipos de garantía
- `warranty_statuses` — catálogo de estados de garantía
- `warranty_claim_types` — catálogo de tipos de reclamo
- `warranty_claim_statuses` — catálogo de estados de reclamo

### Dependencias con módulos ya migrados
- `jobs` (módulo 8) — warranties y warranty_claims pertenecen a un job (job_id)
- `equipment` (módulo 11) — warranty_equipment comparte estructura con job_equipment, y existe clonación desde job_equipment
- `users` (módulo 1) — usuario que gestiona garantías
- `properties` (módulo 6) — propiedad asociada al job

## Operaciones CRUD esperadas

### Warranties
- **List**: GET /api/v1/warranties — con filtros por jobId, warrantyTypeId, warrantyStatusId, week, search, paginación
- **Create**: POST /api/v1/warranties — crear garantía para un job
- **Get**: GET /api/v1/warranties/{id} — obtener garantía por ID (incluir equipment, type, status)
- **Update**: PUT /api/v1/warranties/{id} — actualizar garantía
- **Delete**: DELETE /api/v1/warranties/{id} — soft delete

### Warranty Claims
- **List**: GET /api/v1/warranty-claims — con filtros por jobId, claimTypeId, claimStatusId, week, search, paginación
- **Create**: POST /api/v1/warranty-claims — crear reclamo
- **Get**: GET /api/v1/warranty-claims/{id} — obtener reclamo por ID
- **Update**: PUT /api/v1/warranty-claims/{id} — actualizar reclamo
- **Delete**: DELETE /api/v1/warranty-claims/{id} — soft delete

### Warranty Equipment (sub-recurso de warranties)
- **List**: GET /api/v1/warranties/{warrantyId}/equipment — listar equipos de garantía
- **Create**: POST /api/v1/warranties/{warrantyId}/equipment — crear equipo en garantía
- **Update**: PUT /api/v1/warranties/{warrantyId}/equipment/{id} — actualizar equipo
- **Delete**: DELETE /api/v1/warranties/{warrantyId}/equipment/{id} — soft delete

### Catálogos (CRUD simple para cada uno)
- **WarrantyTypes**: GET/POST/PUT/DELETE /api/v1/warranty-types
- **WarrantyStatuses**: GET/POST/PUT/DELETE /api/v1/warranty-statuses
- **WarrantyClaimTypes**: GET/POST/PUT/DELETE /api/v1/warranty-claim-types
- **WarrantyClaimStatuses**: GET/POST/PUT/DELETE /api/v1/warranty-claim-statuses

## Notas adicionales
- Warranty tiene campos: warranty_number, job_id, warranty_type_id, warranty_status_id, date_submitted (fecha), agreement_number, audit_done (boolean), notes
- WarrantyClaim tiene muchos campos de partes: warranty_part, manufacturer, model_number, part_number, replacement_part_number, part_distributor, part_invoice_number, old_part_serial_number, new_part_serial_number, serial, claim_number, notes
- WarrantyEquipment comparte la misma estructura HVAC que PropertyEquipment y JobEquipment (area, outdoor_*, furnace_*, evaporator_*, air_handler_*)
- Existe `WarrantyEquipment::cloneFromJobEquipment()` que copia equipos del job a la garantía — evaluar si migrar como endpoint o lógica interna
- Warranties y Claims tienen filtro por week_number (a través del job relacionado) con soporte para "unassigned"
- Warranties y Claims tienen sorting por week_number (join con jobs)
- Los catálogos (warranty_types, warranty_statuses, warranty_claim_types, warranty_claim_statuses) son CRUDs simples con campos: label, class (para CSS)
- Este es un módulo grande — considerar dividir la implementación en fases: 1) Catálogos, 2) Warranties + Equipment, 3) Claims

## Entregables de documentación al finalizar

### Actualización de Swagger
- Asegurar annotations completas en cada handler (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)
- Ejecutar `swag init -g cmd/api/main.go -o docs/swagger` para regenerar la documentación Swagger con los nuevos endpoints de warranties, warranty-claims, warranty-equipment y todos los catálogos
- Verificar que los endpoints nuevos aparezcan correctamente y que los endpoints previos no se hayan perdido

### Colección Postman
- Crear colecciones Postman para cada entidad del módulo siguiendo el mismo estándar y formato de las colecciones existentes (ver como referencia: `docs/postman_customer_collection.json`, `docs/postman_jobs_collection.json`):
  - `docs/postman_warranties_collection.json` (incluir warranty equipment como sub-recurso)
  - `docs/postman_warranty_claims_collection.json`
  - `docs/postman_warranty_types_collection.json`
  - `docs/postman_warranty_statuses_collection.json`
  - `docs/postman_warranty_claim_types_collection.json`
  - `docs/postman_warranty_claim_statuses_collection.json`
- Incluir todas las operaciones CRUD, variables `{{baseUrl}}` y `{{accessToken}}`, body de ejemplo con datos realistas de `Data_prd.sql`, headers correctos, y ejemplos de query params para filtros y paginación
- Las colecciones deben estar listas para importar en Postman y ejecutar pruebas inmediatamente
```

---

## Ejemplo de uso completo (Módulo 13: Actividades y Comunicaciones)

```
Necesito migrar el módulo 13: Actividades y Comunicaciones (Activities/Logs) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/JobActivityLog.php`, `app/Models/JobVisit.php`, `app/Models/JobResident.php`, `app/Models/JobTask.php`, `app/Models/JobRate.php`, `app/Models/JobRateStatus.php`
  - Controlador(es): `app/Http/Controllers/JobActivityController.php`, `app/Http/Controllers/JobVisitController.php`, `app/Http/Controllers/JobResidentController.php`, `app/Http/Controllers/JobTaskController.php`, `app/Http/Controllers/JobRateController.php`
  - Request(s): `app/Http/Requests/JobActivities/StoreRequest.php`
  - Request(s): `app/Http/Requests/JobVisits/StoreRequest.php`, `UpdateRequest.php`
  - Request(s): `app/Http/Requests/JobResidents/StoreRequest.php`, `UpdateRequest.php`
  - Request(s): `app/Http/Requests/JobTasks/StoreRequest.php`, `UpdateRequest.php`
  - Request(s): `app/Http/Requests/JobRates/StoreRequest.php`, `UpdateRequest.php`
  - Rutas: `routes/web.php` (buscar rutas anidadas dentro de jobs: activities, visits, residents, tasks, rates)

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql`
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tablas involucradas
- `job_activity_logs` — registro de actividades/notas de un job
- `job_visits` — reportes de visitas a un job
- `job_residents` — residentes asociados a un job
- `job_tasks` — tareas asignadas en un job
- `job_rates` — tarifas/comisiones de técnicos por job
- `job_rate_statuses` — catálogo de estados de tarifas

### Dependencias con módulos ya migrados
- `jobs` (módulo 8) — todos los sub-recursos pertenecen a un job (job_id)
- `users` (módulo 1) — job_activity_logs, job_visits, job_tasks y job_rates tienen user_id
- `task_statuses` (módulo 7) — job_tasks referencia task_status_id

## Operaciones CRUD esperadas

### Job Activity Logs (sub-recurso de jobs)
- **List**: GET /api/v1/jobs/{jobId}/activities — listar actividades/notas del job
- **Create**: POST /api/v1/jobs/{jobId}/activities — crear nota/actividad
- **Delete**: DELETE /api/v1/jobs/{jobId}/activities/{id} — eliminar actividad (NO tiene update)

### Job Visits (sub-recurso de jobs)
- **List**: GET /api/v1/jobs/{jobId}/visits — listar visitas/reportes del job
- **Create**: POST /api/v1/jobs/{jobId}/visits — crear reporte de visita
- **Get**: GET /api/v1/jobs/{jobId}/visits/{id} — obtener visita por ID
- **Update**: PUT /api/v1/jobs/{jobId}/visits/{id} — actualizar visita
- **Delete**: DELETE /api/v1/jobs/{jobId}/visits/{id} — soft delete
- **Download**: GET /api/v1/jobs/{jobId}/visits/{id}/download — descargar reporte (evaluar si migrar)

### Job Residents (sub-recurso de jobs)
- **List**: GET /api/v1/jobs/{jobId}/residents — listar residentes del job
- **Create**: POST /api/v1/jobs/{jobId}/residents — crear residente
- **Update**: PUT /api/v1/jobs/{jobId}/residents/{id} — actualizar residente
- **Delete**: DELETE /api/v1/jobs/{jobId}/residents/{id} — soft delete

### Job Tasks (sub-recurso de jobs)
- **List**: GET /api/v1/jobs/{jobId}/tasks — listar tareas del job
- **List Global**: GET /api/v1/tasks — listar todas las tareas (vista global, ruta independiente)
- **Create**: POST /api/v1/jobs/{jobId}/tasks — crear tarea
- **Update**: PUT /api/v1/jobs/{jobId}/tasks/{id} — actualizar tarea
- **Delete**: DELETE /api/v1/jobs/{jobId}/tasks/{id} — soft delete
- **SendNotification**: PUT /api/v1/jobs/{jobId}/tasks/{id}/send-notification — evaluar si migrar

### Job Rates (sub-recurso de jobs)
- **List**: GET /api/v1/jobs/{jobId}/rates — listar tarifas del job
- **Create**: POST /api/v1/jobs/{jobId}/rates — crear tarifa
- **Update**: PUT /api/v1/jobs/{jobId}/rates/{id} — actualizar tarifa
- **Delete**: DELETE /api/v1/jobs/{jobId}/rates/{id} — soft delete
- **CalculatePayment**: POST /api/v1/calculate-rate-payment — calcular pago de tarifa

### Job Rate Statuses (catálogo)
- **CRUD**: GET/POST/PUT/DELETE /api/v1/job-rate-statuses

## Notas adicionales
- JobActivityLog: campos simples (job_id, log, type, user_id). Solo Create y Delete, NO tiene Update. NO tiene SoftDeletes en el modelo pero verificar en db_structure.sql
- JobVisit: campos (job_id, user_id, viewable_by (JSON array de role IDs), date, report). Tiene scope `scopeViewable()` que filtra por roles del usuario autenticado — implementar como filtro en el handler. Tiene relación morphMany con Files
- JobResident: campos simples (job_id, name, mobile_phone, home_phone, email). Contactos de residentes en la propiedad del job
- JobTask: campos (job_id, user_id, task, task_status_id, due_date). Referencia a task_statuses (módulo 7). Tiene vista global de tareas además de la vista por job. Tiene acción sendNotification — evaluar si migrar en este punto o en módulo 14
- JobRate: campos monetarios (job_id, user_id, job_rate_status_id, sale_price, rate_percent, rate_flat, tech_parts, company_parts, parts_replaced, deduction, payment). Tiene método estático `calculatePayment()` con fórmula: ((salePrice - techParts - companyParts) * (commissionRate/100)) + flatRate - deduction. Tiene filtros por userId, paid status
- JobRateStatus: catálogo simple con campos (label, class)
- Este es un módulo grande con 6 sub-entidades — considerar dividir la implementación en fases:
  1) JobActivityLogs + JobResidents (más simples)
  2) JobTasks + JobRateStatuses
  3) JobVisits
  4) JobRates (más complejo por la lógica de cálculo)

## Entregables de documentación al finalizar

### Actualización de Swagger
- Asegurar annotations completas en cada handler (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)
- Ejecutar `swag init -g cmd/api/main.go -o docs/swagger` para regenerar la documentación Swagger con los nuevos endpoints de job-activities, job-visits, job-residents, job-tasks, job-rates y job-rate-statuses
- Verificar que los endpoints nuevos aparezcan correctamente y que los endpoints previos no se hayan perdido

### Colección Postman
- Crear colecciones Postman para cada sub-entidad del módulo siguiendo el mismo estándar y formato de las colecciones existentes (ver como referencia: `docs/postman_customer_collection.json`, `docs/postman_jobs_collection.json`):
  - `docs/postman_job_activities_collection.json`
  - `docs/postman_job_visits_collection.json`
  - `docs/postman_job_residents_collection.json`
  - `docs/postman_job_tasks_collection.json`
  - `docs/postman_job_rates_collection.json`
  - `docs/postman_job_rate_statuses_collection.json`
- Incluir todas las operaciones CRUD de cada sub-entidad, variables `{{baseUrl}}` y `{{accessToken}}`, body de ejemplo con datos realistas de `Data_prd.sql`, headers correctos, y ejemplos de query params para filtros y paginación
- Las colecciones deben estar listas para importar en Postman y ejecutar pruebas inmediatamente
```

---

## Referencia rápida: Módulos y su estado

| # | Módulo | Branch | Estado |
|---|--------|--------|--------|
| 1 | Usuarios y Roles | main | ✅ Migrado |
| 2 | Settings | main | ✅ Migrado |
| 3 | Workflows | main | ✅ Migrado |
| 4 | Customers | main | ✅ Migrado |
| 5 | Supervisors | — | ⏳ Pendiente |
| 6 | Properties | main | ✅ Migrado |
| 7 | Catálogos de Jobs | main | ✅ Migrado |
| 8 | Jobs | feature/jobs-all | ✅ Migrado |
| 9 | Quotes | — | ⏳ Pendiente |
| 10 | Invoices | — | ⏳ Pendiente |
| 11 | Equipment | — | ⏳ Pendiente |
| 12 | Warranties | — | ⏳ Pendiente |
| 13 | Activities/Logs | — | ⏳ Pendiente |
| 14 | Email/SMS | — | ⏳ Pendiente |
| 15 | Alerts | — | ⏳ Pendiente |
