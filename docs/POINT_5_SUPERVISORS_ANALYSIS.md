# Módulo 5: Supervisores (Supervisors) - Análisis

## Estructura de tabla (desde db_structure.sql)

```sql
CREATE TABLE `supervisors` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` bigint unsigned NOT NULL,
  `name` varchar(191) NOT NULL,
  `phone` varchar(191) DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `supervisors_customer_id_foreign` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`)
);
```

## Análisis de funcionalidad original (Laravel)

### Modelo: Supervisor.php
- **fillable**: id, customer_id, name, phone, email
- **SoftDeletes**: Sí (deleted_at)
- **Traits**: FilterSortSearchable, Paged
- **Searchable**: name, phone, email
- **Relación**: belongsTo Customer
- **Scope**: `scopeForCurrentUser()` — si el usuario tiene permiso `job_view_user_only`, filtra supervisores de clientes que tengan jobs donde `user_id` sea el usuario actual

### Controlador: SupervisorController.php
- CRUD completo: index, create, store, edit, update, destroy
- index: usa `forCurrentUser()` + `filterSortSearch()` + `paged()`
- create/edit: carga lista de customers para dropdown

### Policy: SupervisorPolicy.php
- viewAny: `return true` (bug — hardcoded)
- view: `supervisor_view`
- create: `supervisor_edit`
- update: `supervisor_edit` (+ no puede estar soft-deleted)
- delete: `supervisor_delete`

### Relación con Customer
- Customer tiene `hasMany(Supervisor::class)`
- CustomerController tiene método `supervisors()` que lista supervisores de un customer

### Relación con Job
- Campo `jobs.supervisor_ids` (varchar CSV, ej: "1,5,12")
- Método `Job::supervisorEmails()` parsea CSV y obtiene emails
- Se usa en dispatch email a supervisores
- **NO se migra en este módulo** — se integrará con comunicaciones

## Plan de ejecución

### Fase 1: Dominio
- entity.go — Supervisor struct, Validate()
- repository.go — Interface Repository
- usecase.go — UseCase struct, Service interface, NewUseCase()
- create.go, get_by_id.go, list.go, update.go, delete.go
- mock.go, usecase_test.go

### Fase 2: Repository MySQL
- repository.go, create.go, get_by_id.go, list.go, update.go, delete.go

### Fase 3: Handler REST
- handler.go (DTOs, RegisterRoutes, helpers)
- create.go, get.go, list.go, update.go, delete.go

### Fase 4: Integración
- container.go — inicializar repo, usecase, handler
- router.go — registrar rutas
- Customer handler — agregar ruta /customers/{id}/supervisors

## Decisiones de diseño

1. **Validación de FK customer_id**: El UseCase valida que el customer existe usando customer.Repository
2. **Scope forCurrentUser**: NO se implementa en esta fase (requiere contexto de usuario autenticado y permisos). Se dejará como filtro futuro
3. **Relación con Jobs (supervisor_ids CSV)**: NO se migra. Se mantiene la referencia para el módulo de comunicaciones
4. **Sub-ruta en customers**: Se agrega GET /customers/{id}/supervisors como ruta en el handler de supervisors (no en customer handler para evitar dependencia circular)

## Filtros soportados en List
- `customer_id` (int64) — filtrar por cliente
- `search` (string) — buscar en name, phone, email
