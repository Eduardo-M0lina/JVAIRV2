# Punto 7: Módulo de Trabajos (Jobs) - Análisis y Plan de Ejecución

## 1. Análisis de la Estructura de Base de Datos

### Tabla `jobs`
```sql
CREATE TABLE `jobs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `work_order` varchar(255) DEFAULT NULL,
  `date_received` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `job_category_id` bigint unsigned NOT NULL,
  `job_priority_id` bigint unsigned NOT NULL,
  `job_status_id` bigint unsigned NOT NULL,
  `technician_job_status_id` bigint unsigned DEFAULT NULL,
  `workflow_id` bigint unsigned NOT NULL,
  `property_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `supervisor_ids` varchar(191) DEFAULT NULL,
  `dispatch_date` timestamp NULL DEFAULT NULL,
  `completion_date` timestamp NULL DEFAULT NULL,
  `week_number` int DEFAULT NULL,
  `route_number` int DEFAULT NULL,
  `scheduled_time_type` varchar(191) DEFAULT NULL,
  `scheduled_time` varchar(191) DEFAULT NULL,
  `internal_job_notes` text,
  `quick_notes` text,
  `job_report` text,
  `installation_due_date` timestamp NULL DEFAULT NULL,
  `cage_required` tinyint(1) NOT NULL DEFAULT '0',
  `warranty_claim` tinyint(1) NOT NULL DEFAULT '0',
  `warranty_registration` tinyint(1) NOT NULL DEFAULT '0',
  `job_sales_price` decimal(8,2) DEFAULT NULL,
  `money_turned_in` decimal(8,2) DEFAULT NULL,
  `closed` tinyint(1) NOT NULL DEFAULT '0',
  `dispatch_notes` text,
  `call_logs` text,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `due_date` timestamp NULL DEFAULT NULL,
  `call_attempted` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  -- FK: job_category_id, job_priority_id, job_status_id, workflow_id, property_id, user_id, technician_job_status_id
);
```

### Tabla `job_status_workflow` (tabla pivote)
```sql
CREATE TABLE `job_status_workflow` (
  `job_status_id` bigint unsigned NOT NULL,
  `workflow_id` bigint unsigned NOT NULL,
  `order` int NOT NULL DEFAULT '0',
  -- FK: job_status_id -> job_statuses, workflow_id -> workflows
);
```

## 2. Análisis de Funcionalidad Original (JVAIR Laravel)

### Modelo Job (Job.php)
- **SoftDeletes**: Usa `deleted_at` para eliminación lógica
- **Fillable**: Todos los campos editables del job
- **Relaciones**: property, customer (through property), category, priority, status, tech_status, workflow, user, emails, sms_messages, logs, activity_notes, residents, invoices, payments, tasks, visits, rates, claims, quotes, equipment, warranties, files
- **Scopes**: `forCurrentUser` (filtra por user_id si tiene permiso `job_view_user_only`), `open` (closed=0)
- **Filtros**: `isClosedFilter`, `userIdFilter`
- **Sorting**: `statusSort` (ordena por workflow status order)
- **Searchable**: work_order, property.property_code, property.street, property.city, property.state, property.zip, customer.name, customer.email, customer.phone, etc.

### Controlador JobController
- **index**: Lista jobs con filtros, paginación, búsqueda. Filtra por categoría y usuario.
- **store**: Crea job con property (nueva o existente), asigna workflow y status inicial del customer, copia equipment de property, crea activity log.
- **update**: Actualiza job con validación de permisos por campo. Lógica de tech_status -> job_status automática. Crea alerts para técnicos. Crea activity log.
- **destroy**: Soft delete del job.
- **close**: Cierra job, actualiza status, opcionalmente actualiza equipment de property y crea warranty.
- **sendDispatchEmail**: Envía email de dispatch.
- **sendDispatchSupervisorEmail**: Envía email a supervisor.
- **sendDispatchSMS**: Envía SMS de dispatch.
- **logAttemptedCall**: Registra intento de llamada.

### Validación (StoreRequest)
- `property.customer_id`: required, exists
- `job.property_id`: nullable, required_without:property.street
- `job.work_order`: required
- `job.job_category_id`: required, exists
- `job.job_priority_id`: required, exists
- `job.date_received`: nullable, date
- `job.dispatch_notes`: nullable, string
- `job.due_date`: nullable, date
- `job.user_id`: nullable, exists
- `job.dispatch_date`: nullable, date
- `job.week_number`: nullable, numeric

### Validación (UpdateRequest)
- Todos los campos son nullable excepto los que se validan por permisos en el controlador.

## 3. Plan de Ejecución

### Fase 1: Domain Entity
- Crear `pkg/domain/job/entity.go` con struct Job y validación
- Crear `pkg/domain/job/errors.go` con errores del dominio
- Crear `pkg/domain/job/repository.go` con interface Repository
- Crear `pkg/domain/job/usecase.go` con UseCase y Service interface

### Fase 2: Use Cases
- `Create`: Valida campos, verifica FKs, asigna workflow/status inicial
- `GetByID`: Obtiene job por ID (con soft delete check)
- `List`: Lista paginada con filtros (search, closed, category, status, priority, user, property, customer, date range, sort)
- `Update`: Actualiza job, lógica tech_status -> job_status
- `Delete`: Soft delete
- `Close`: Cierra job con status

### Fase 3: MySQL Repository
- CRUD completo con soft deletes
- List con JOINs para búsqueda en property/customer
- Sorting por status workflow order

### Fase 4: REST Handler
- Handlers separados por archivo (list, create, get, update, delete, close)
- DTOs para request/response
- Filtros de query params

### Fase 5: Integración
- Registrar en container.go y router.go

### Fase 6: Tests, Swagger y Postman

## 4. Decisiones de Diseño

### Alcance para esta migración (fiel al original)
Se migran las operaciones CRUD principales del Job:
1. **List** (con filtros, búsqueda, paginación, sorting)
2. **Create** (con asignación automática de workflow/status)
3. **Get** (por ID, con relaciones básicas)
4. **Update** (con lógica tech_status -> job_status)
5. **Delete** (soft delete)
6. **Close** (cerrar job con status)

### NO se migran en este punto (serán módulos separados):
- Dispatch email/SMS (módulo de comunicaciones)
- Job Equipment (punto 10)
- Job Tasks (punto 16 del JVAIR_V1_ENDPOINTS)
- Job Visits/Reports
- Job Rates
- Job Activities/Notes
- Invoices/Quotes (puntos 8-9)
- Warranties

### Filtros soportados en List
- `search`: Búsqueda en work_order, property (street, city, state, zip, property_code), customer (name)
- `closed` / `isClosed`: Filtro por estado cerrado (default: open)
- `jobCategoryId`: Filtro por categoría
- `jobStatusId`: Filtro por estado
- `jobPriorityId`: Filtro por prioridad
- `userId`: Filtro por técnico asignado (soporta "unassigned")
- `propertyId`: Filtro por propiedad
- `workflowId`: Filtro por workflow
- `sort` / `direction`: Ordenamiento (default: created_at DESC)
