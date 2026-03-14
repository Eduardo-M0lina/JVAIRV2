# Module 14 Migration Complete

## Resumen

Se migró el módulo 14 de **Comunicaciones (Email/SMS)** desde Laravel hacia JVAIRV2 en Go.

Módulos implementados:

- `email_templates`
- `sms_templates`
- `job_emails`
- `job_sms`

## Entidades implementadas

### Email Templates

Tabla:

- `email_templates`

Campos soportados:

- `id`
- `label`
- `subject`
- `body`
- `is_active`
- `created_at`
- `updated_at`

Operaciones implementadas:

- List
- Create
- GetByID
- Update
- Delete

### SMS Templates

Tabla:

- `sms_templates`

Campos soportados:

- `id`
- `label`
- `message`
- `is_active`
- `created_at`
- `updated_at`

Operaciones implementadas:

- List
- Create
- GetByID
- Update
- Delete

### Job Emails

Tabla:

- `job_emails`

Campos soportados:

- `id`
- `job_id`
- `recipients`
- `type`
- `created_at`
- `updated_at`

Operaciones implementadas:

- List por `jobId`
- Create
- Delete

### Job SMS

Tabla:

- `job_sms`

Campos soportados:

- `id`
- `job_id`
- `recipients`
- `type`
- `message`
- `created_at`
- `updated_at`

Operaciones implementadas:

- List por `jobId`
- Create
- Delete

## Endpoints implementados

### Email Templates

- `GET /api/v1/email-templates`
- `POST /api/v1/email-templates`
- `GET /api/v1/email-templates/{id}`
- `PUT /api/v1/email-templates/{id}`
- `DELETE /api/v1/email-templates/{id}`

### SMS Templates

- `GET /api/v1/sms-templates`
- `POST /api/v1/sms-templates`
- `GET /api/v1/sms-templates/{id}`
- `PUT /api/v1/sms-templates/{id}`
- `DELETE /api/v1/sms-templates/{id}`

### Job Emails

- `GET /api/v1/jobs/{jobId}/emails`
- `POST /api/v1/jobs/{jobId}/emails`
- `DELETE /api/v1/jobs/{jobId}/emails/{id}`

### Job SMS

- `GET /api/v1/jobs/{jobId}/sms`
- `POST /api/v1/jobs/{jobId}/sms`
- `DELETE /api/v1/jobs/{jobId}/sms/{id}`

## Decisiones de diseño

### 1. Sin soft delete

Las 4 tablas fueron implementadas con **hard delete**, porque la estructura real de producción no contiene `deleted_at`.

### 2. Validación de FK de jobs

Los módulos `job_emails` y `job_sms` validan que el `job_id` exista usando adapters `JobExistsChecker`.

### 3. Manejo de `recipients`

En base de datos, `recipients` es `blob`.

En JVAIRV2 se decidió:

- recibir `recipients` como `[]string` en los requests JSON
- almacenar internamente como **JSON array serializado** en la columna `blob`
- leer con compatibilidad para datos legacy que puedan venir en formato CSV

Esto permite compatibilidad con datos previos de Laravel y deja el formato nuevo más consistente.

### 4. No se migró envío real

No se implementó envío real por SMTP, Mailgun, Twilio o AWS SNS.

Esta migración cubre únicamente:

- catálogos de plantillas
- logging/auditoría de emails enviados
- logging/auditoría de SMS enviados

## Integración realizada

Se actualizaron:

- `cmd/api/container.go`
- `pkg/rest/router/router.go`

Para registrar:

- repositories
- use cases
- handlers
- rutas del módulo 14

## Verificación

Se verificó compilación con:

```bash
go build ./...
```

Resultado:

- compilación exitosa

## Pendiente sugerido

Para completar la parte operativa del módulo en una fase posterior:

- integración con proveedor real de email
- integración con proveedor real de SMS
- trazabilidad de template usado al enviar
- plantillas con renderizado de placeholders
- permisos específicos por endpoint si se requiere granularidad adicional
