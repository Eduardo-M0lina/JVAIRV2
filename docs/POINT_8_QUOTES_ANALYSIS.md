# Punto 8: Módulo de Cotizaciones (Quotes) - Análisis y Plan de Ejecución

## 1. Análisis de la Estructura de Base de Datos

### Tabla `quote_statuses`
```sql
CREATE TABLE `quote_statuses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `label` varchar(255) NOT NULL,
  `class` varchar(255) DEFAULT NULL,
  `order` int NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
);
```
**Nota**: No tiene `deleted_at` — no usa soft deletes.

### Tabla `quotes`
```sql
CREATE TABLE `quotes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job_id` bigint unsigned NOT NULL,
  `quote_number` varchar(255) NOT NULL,
  `quote_status_id` bigint unsigned NOT NULL,
  `amount` decimal(8,2) NOT NULL DEFAULT '0.00',
  `description` text,
  `notes` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  FK: job_id -> jobs(id), quote_status_id -> quote_statuses(id)
);
```

## 2. Análisis de Funcionalidad Original (JVAIR Laravel)

### Modelo Quote
- **SoftDeletes**: Sí (`deleted_at`)
- **Fillable**: `id`, `job_id`, `quote_number`, `quote_status_id`, `amount`, `description`, `notes`
- **Relaciones**: `status` (belongsTo QuoteStatus), `job` (belongsTo Job)
- **Searchable**: `quote_number`, `job.work_order`, `job.property.*`, `job.property.customer.*`

### Modelo QuoteStatus
- **SoftDeletes**: No
- **Fillable**: `id`, `label`, `class`, `order`
- Sin relaciones definidas

### Controlador QuoteController
- **index**: Lista quotes con filtros, paginación, búsqueda
- **store**: Crea quote con validación (quote_number, job_id, quote_status_id, amount requeridos)
- **update**: Actualiza quote con misma validación
- **destroy**: Soft delete
- **email**: Envía cotización por email (NO se migra en este punto)

### Validación (StoreRequest / UpdateRequest)
- `quote_number`: required, string
- `job_id`: required, exists:jobs,id
- `quote_status_id`: required, exists:quote_statuses,id
- `description`: nullable, string
- `amount`: required, numeric
- `notes`: nullable, string

## 3. Plan de Ejecución

### Fase 1: QuoteStatus (catálogo)
- `pkg/domain/quote_status/` — entity, errors, repository, usecase, CRUD
- `pkg/repository/mysql/quote_status/` — repository, CRUD
- `pkg/rest/handler/quote_status/` — handler con Swagger

### Fase 2: Quote (entidad principal)
- `pkg/domain/quote/` — entity, errors, repository, usecase, CRUD, mock, tests
- `pkg/repository/mysql/quote/` — repository, CRUD, adapters (JobChecker, QuoteStatusChecker)
- `pkg/rest/handler/quote/` — handler con Swagger

### Fase 3: Integración
- Registrar en `container.go` y `router.go`

### Fase 4: Tests, Swagger y Postman

## 4. Decisiones de Diseño

### Alcance para esta migración
1. **QuoteStatuses CRUD** (catálogo simple)
2. **Quotes CRUD** (con validación de FKs: job, quote_status)
3. **List con filtros**: search, jobId, quoteStatusId, paginación
4. **Soft delete** para quotes

### NO se migra en este punto
- Envío de cotización por email (módulo de comunicaciones)
- Búsqueda profunda en property/customer (simplificado a quote_number y job.work_order)

### Filtros soportados en List
- `search`: Búsqueda en quote_number
- `jobId`: Filtro por trabajo
- `quoteStatusId`: Filtro por estado
- `sort` / `direction`: Ordenamiento (default: created_at DESC)
