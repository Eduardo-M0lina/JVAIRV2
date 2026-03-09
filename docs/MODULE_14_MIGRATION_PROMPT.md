# Prompt de Migración: Módulo 14 — Comunicaciones (Email/SMS)

```
Necesito migrar el módulo 14: Comunicaciones (Email/SMS) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo(s): `app/Models/EmailTemplate.php`, `app/Models/SmsTemplate.php`, `app/Models/JobEmail.php`, `app/Models/JobSms.php`
  - Controlador(es): `app/Http/Controllers/EmailTemplateController.php`, `app/Http/Controllers/SmsTemplateController.php`, `app/Http/Controllers/JobEmailController.php`, `app/Http/Controllers/JobSmsController.php`
  - Request(s): buscar en `app/Http/Requests/` carpetas relacionadas con EmailTemplate, SmsTemplate, JobEmail, JobSms
  - Rutas: `routes/web.php` y `routes/api.php` (buscar rutas de email-templates, sms-templates, job emails, job sms)

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql` — SIEMPRE confiar en este archivo para la estructura real de las tablas
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tablas involucradas (estructura confirmada desde db_structure.sql)

**email_templates:**
```sql
id, label varchar(191), subject varchar(191), body text,
is_active tinyint(1) DEFAULT 1, created_at, updated_at
-- NO tiene deleted_at (sin soft delete)
```

**sms_templates:**
```sql
id, label varchar(191), message text,
is_active tinyint(1) DEFAULT 1, created_at, updated_at
-- NO tiene deleted_at (sin soft delete)
```

**job_emails:**
```sql
id, job_id bigint unsigned (FK → jobs), recipients blob,
type varchar(255), created_at, updated_at
-- NO tiene deleted_at (sin soft delete). Sin campo message (message va en template)
```

**job_sms:**
```sql
id, job_id bigint unsigned (FK → jobs), recipients blob,
type varchar(191), message text, created_at, updated_at
-- NO tiene deleted_at (sin soft delete)
```

### Dependencias con módulos ya migrados
- `jobs` (módulo 8) — job_emails y job_sms pertenecen a un job (job_id, FK obligatoria)
- `users` (módulo 1) — verificar si hay user_id en los controladores (no está en la tabla pero puede usarse en lógica de negocio)

## Arquitectura y Estándares del Proyecto JVAIRV2

### Estructura de archivos por módulo
Cada módulo sigue esta estructura estricta:

```
pkg/domain/{{modulo}}/
├── entity.go          # Struct, Filters, validaciones
├── errors.go          # var ErrXxx = errors.New(...)
├── repository.go      # Interface Repository
├── usecase.go         # Service interface, UseCase struct, NewUseCase(), checker interfaces
├── create.go          # Método Create
├── get_by_id.go       # Método GetByID
├── list.go            # Método List
├── update.go          # Método Update (solo aplica si la entidad lo soporta)
└── delete.go          # Método Delete (hard delete, sin soft delete para estas tablas)

pkg/repository/mysql/{{modulo}}/
├── repository.go      # Repository struct, NewRepository()
├── create.go
├── get_by_id.go
├── list.go
├── update.go          # Si aplica
├── delete.go          # Hard delete (DELETE FROM, NO soft delete)
└── adapters.go        # JobExistsChecker para job_emails y job_sms

pkg/rest/handler/{{modulo}}/
└── handler.go         # Handler struct, NewHandler(), DTOs, Swagger annotations, todos los métodos
```

### Patrones obligatorios
1. **Sin soft deletes**: Estas tablas NO tienen deleted_at. DELETE = hard delete (`DELETE FROM tabla WHERE id = ?`). Las queries List NO filtran `deleted_at IS NULL`.
2. **recipients es blob**: Almacena JSON (ej: `["email@a.com","email@b.com"]`). En Go usar `[]byte` en el repository y deserializar a `[]string` en el dominio/handler.
3. **Paginación**: List retorna `([]*Entity, int64, error)`. Query params: `page` (default 1), `limit` (default 15).
4. **Errores de dominio**: var en errors.go, mapeados a HTTP status en handler.
5. **Validación de FKs**: UseCase valida job_id con JobExistsChecker.
6. **Imports**: `github.com/your-org/jvairv2/...`
7. **Router**: Chi (`github.com/go-chi/chi/v5`), usar `chi.URLParam(r, "jobId")` para path params.
8. **json.Encode**: SIEMPRE asignar `_ = json.NewEncoder(w).Encode(x)` (errcheck).
9. **rows.Close**: SIEMPRE usar `defer func() { _ = rows.Close() }()` (errcheck).
10. **Swagger**: Annotations completas en cada método handler. Al finalizar regenerar con `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs`

### Integración (después de crear el módulo)
1. **`cmd/api/container.go`**: Agregar imports, inicializar repository, adapters/checkers, use case y handler.
2. **`pkg/rest/router/router.go`**: Agregar parámetro en `New()` y registrar rutas con `handler.RegisterRoutes(r)`.

## Operaciones CRUD esperadas

### Email Templates (catálogo independiente)
- **List**: GET /api/v1/email-templates — listar templates (filtro por is_active, paginación)
- **Create**: POST /api/v1/email-templates — crear template
- **Get**: GET /api/v1/email-templates/{id} — obtener por ID
- **Update**: PUT /api/v1/email-templates/{id} — actualizar template
- **Delete**: DELETE /api/v1/email-templates/{id} — hard delete

### SMS Templates (catálogo independiente)
- **List**: GET /api/v1/sms-templates — listar templates (filtro por is_active, paginación)
- **Create**: POST /api/v1/sms-templates — crear template
- **Get**: GET /api/v1/sms-templates/{id} — obtener por ID
- **Update**: PUT /api/v1/sms-templates/{id} — actualizar template
- **Delete**: DELETE /api/v1/sms-templates/{id} — hard delete

### Job Emails (sub-recurso de jobs — log de emails enviados)
- **List**: GET /api/v1/jobs/{jobId}/emails — listar emails enviados de un job (paginación)
- **Create**: POST /api/v1/jobs/{jobId}/emails — registrar email enviado
- **Delete**: DELETE /api/v1/jobs/{jobId}/emails/{id} — hard delete

### Job SMS (sub-recurso de jobs — log de SMS enviados)
- **List**: GET /api/v1/jobs/{jobId}/sms — listar SMS enviados de un job (paginación)
- **Create**: POST /api/v1/jobs/{jobId}/sms — registrar SMS enviado
- **Delete**: DELETE /api/v1/jobs/{jobId}/sms/{id} — hard delete

## Notas adicionales
- **NO migrar el envío real** (SMTP, Twilio, etc.) en este punto — solo el CRUD/logging de registros enviados. El envío se integrará en una fase posterior.
- **email_templates** usa `subject` + `body` (HTML). `body` puede contener variables placeholder como `{{job_number}}`, `{{technician_name}}` — almacenar como texto plano, sin procesar templates.
- **sms_templates** solo tiene `message` (texto plano con posibles placeholders).
- **job_emails.recipients** y **job_sms.recipients** son `blob` en MySQL → en Go mapear como `[]byte` en SQL Scan y convertir a `[]string` (JSON array) en la entidad de dominio. En el handler, recibir como `[]string` en el request body.
- **job_emails** NO tiene campo `message` propio — el mensaje va referenciado por el template o enviado directamente; almacenar solo destinatarios y tipo.
- **job_sms** SÍ tiene campo `message` propio además de recipients y type.
- `type` en job_emails y job_sms indica el tipo de comunicación (ej: "dispatch", "invoice", "notification") — almacenar como string sin validación de enum.
- **Sin soft delete** en ninguna de las 4 tablas → implementar hard delete, sin campo deleted_at, sin filtro `deleted_at IS NULL` en queries.
- Los templates tienen 20-27 registros en producción (AUTO_INCREMENT=21 y 28) — relativamente pocos.
- job_emails tiene ~60k registros, job_sms ~10k — son tablas de auditoría de alto volumen.

## Entregables de documentación al finalizar

### Documento de análisis
- Crear `docs/MODULE_14_MIGRATION_COMPLETE.md` con resumen de: entidades, endpoints implementados, decisiones de diseño, pasos de integración.

### Actualización de Swagger
- Annotations completas en cada handler method.
- Ejecutar `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs`
- Verificar que todos los endpoints nuevos aparezcan en http://localhost:8090/swagger/index.html

### Colecciones Postman
Crear siguiendo el formato de colecciones existentes (`docs/postman_jobs_collection.json`):
- `docs/postman_email_templates_collection.json`
- `docs/postman_sms_templates_collection.json`
- `docs/postman_job_emails_collection.json`
- `docs/postman_job_sms_collection.json`

Requisitos de cada colección:
- Variables de colección: solo IDs (templateId, jobId, emailId, smsId) — **NO** incluir baseUrl ni accessToken en variables de colección (vienen del environment)
- URLs sin prefijo `/api/v1/` repetido: usar `{{baseUrl}}/email-templates`, NO `{{baseUrl}}/api/v1/email-templates`
- Auth: Bearer token con `{{accessToken}}`
- Bodies de ejemplo con datos realistas
- Query params de paginación en endpoints List

### Verificación final
- `go build ./...` debe compilar sin errores
- `pre-commit run --all-files` debe pasar todos los hooks
```
