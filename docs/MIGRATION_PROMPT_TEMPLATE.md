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
11. **Swagger**: Annotations en cada handler method (`@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Router`, `@Security BearerAuth`)

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
- Annotations completas en cada handler
- Ejecutar `swag init` después de implementar

### 5. Colección Postman
- Crear `docs/postman_{{modulo}}_collection.json` con:
  - Todas las operaciones CRUD
  - Variables: `{{baseUrl}}`, `{{accessToken}}`
  - Body de ejemplo con datos realistas
  - Descripción en cada request

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

## Ejemplo de uso completo (Módulo 8: Cotizaciones)

```
Necesito migrar el módulo 8: Cotizaciones (Quotes) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

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
- `jobs` (módulo 7) — quotes pertenecen a un job
- `users` (módulo 1) — usuario que crea la cotización
- `properties` (módulo 5) — propiedad asociada al job

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
```

---

## Referencia rápida: Módulos y su estado

| # | Módulo | Branch | Estado |
|---|--------|--------|--------|
| 1 | Usuarios y Roles | main | ✅ Migrado |
| 2 | Settings | main | ✅ Migrado |
| 3 | Workflows | main | ✅ Migrado |
| 4 | Customers | main | ✅ Migrado |
| 5 | Properties | main | ✅ Migrado |
| 6 | Catálogos de Jobs | main | ✅ Migrado |
| 7 | Jobs | feature/jobs-all | ✅ Migrado |
| 8 | Quotes | — | ⏳ Pendiente |
| 9 | Invoices | — | ⏳ Pendiente |
| 10 | Equipment | — | ⏳ Pendiente |
| 11 | Warranties | — | ⏳ Pendiente |
| 12 | Activities/Logs | — | ⏳ Pendiente |
| 13 | Email/SMS | — | ⏳ Pendiente |
| 14 | Alerts | — | ⏳ Pendiente |
