# Prompt de Migración: Módulo 15 — Alertas (Alerts)

```
Necesito migrar el módulo 15: Alertas (Alerts) del proyecto JVAIR (Laravel) a JVAIRV2 (Go).

## Contexto del Proyecto

### Proyecto origen (Laravel - código fuente a analizar)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIR`
- **Archivos clave a analizar**:
  - Modelo: `app/Models/Alert.php`
  - Controlador: `app/Http/Controllers/AlertController.php`
  - Request: `app/Http/Requests/Alerts/IndexRequest.php`
  - Rutas: `routes/web.php` (buscar sección "Alerts" — NO tiene rutas en api.php)
  - Creación programática en: `app/Http/Controllers/JobActivityController.php`, `app/Http/Controllers/JobController.php`
  - Componente de vista: `app/View/Components/Alerts/Table.php`

### Proyecto destino (Go - donde se implementa)
- **Ruta**: `/Users/eduardo/projects/jvair/JVAIRV2`
- **Base de datos de producción (estructura)**: `scripts/db_structure.sql` — SIEMPRE confiar en este archivo para la estructura real de las tablas
- **Plan de migración**: `docs/MIGRATION_PLAN.md`
- **Datos de producción**: `scripts/Data_prd.sql`

### Tabla involucrada (estructura confirmada desde db_structure.sql)

**alerts:**
```sql
id                bigint unsigned NOT NULL AUTO_INCREMENT,
user_id           bigint unsigned DEFAULT NULL,  -- nullable, FK → users
alert_type        varchar(191) NOT NULL,         -- valores conocidos: 'system'
entity_id         bigint unsigned NOT NULL,      -- ID de la entidad relacionada (polymorphic)
entity_type       varchar(191) NOT NULL,         -- tipo de entidad: 'job', 'call_log'
message_level     varchar(191) NOT NULL,         -- valores conocidos: 'info', 'danger'
message           varchar(191) NOT NULL,
is_read           tinyint(1) NOT NULL DEFAULT 0,
created_at        timestamp NULL DEFAULT NULL,
updated_at        timestamp NULL DEFAULT NULL,
PRIMARY KEY (id),
KEY alerts_user_id_foreign (user_id),
KEY alert_entity_index (entity_id, entity_type),
KEY alerts_alert_type_index (alert_type),
KEY alerts_is_read_index (is_read)
-- NO tiene deleted_at (sin soft delete). ~42k registros en producción.
```

### Dependencias con módulos ya migrados
- `users` (módulo 1) — user_id es FK nullable hacia users (una alerta puede o no estar asignada a un usuario)
- `jobs` (módulo 8) — la entidad relacionada más común es entity_type='job', entity_id=job.id

## Análisis del comportamiento original (Laravel)

El módulo de alertas en Laravel opera de forma diferente al CRUD estándar:
- **NO existe endpoint de creación pública**: las alertas se crean **programáticamente** desde otros controladores (JobController, JobActivityController) cuando ocurren eventos del sistema.
- **NO existe Update completo**: la única "modificación" es `markAsRead()` (is_read = true).
- **NO existe Delete desde la UI**: no hay endpoint de borrado en Laravel.
- Las rutas son todas web (no API): `GET alerts/`, `GET alerts/{id}/open`, `GET alerts/mark-call-log/{job}`.
- `openAlert()`: marca la alerta como leída y redirige a la entidad (job) — funcionalidad web, NO migrar la redirección.
- `markCallLogAsRead()`: desactiva `call_attempted` en el job relacionado — NO migrar esta lógica aquí.

## Arquitectura y Estándares del Proyecto JVAIRV2

### Estructura de archivos del módulo
```
pkg/domain/alert/
├── entity.go       # Struct Alert, validaciones
├── errors.go       # var ErrNotFound, ErrUserNotFound = errors.New(...)
├── repository.go   # Interface Repository
├── usecase.go      # Service interface, UseCase struct, NewUseCase(), UserExistsChecker
├── create.go       # Método Create (uso interno por otros handlers)
├── get_by_id.go    # Método GetByID
├── list.go         # Método List (paginado, filtros)
├── mark_as_read.go # Método MarkAsRead (actualiza is_read=true para un alert)
├── mark_all_read.go# Método MarkAllRead (actualiza is_read=true para todos los del user)
└── delete.go       # Método Delete (hard delete)

pkg/repository/mysql/alert/
├── repository.go
├── create.go
├── get_by_id.go
├── list.go         # Con filtros: user_id, is_read, alert_type, entity_type, paginación
├── mark_as_read.go # UPDATE alerts SET is_read=1, updated_at=NOW() WHERE id=?
├── mark_all_read.go# UPDATE alerts SET is_read=1, updated_at=NOW() WHERE user_id=? AND is_read=0
├── delete.go       # Hard delete: DELETE FROM alerts WHERE id=?
├── unread_count.go # SELECT COUNT(*) WHERE user_id=? AND is_read=0
└── adapters.go     # UserExistsChecker (si se valida user_id en Create)

pkg/rest/handler/alert/
└── handler.go      # Handler struct, NewHandler(), RegisterRoutes(), DTOs, todos los métodos HTTP con Swagger
```

### Patrones obligatorios
1. **Sin soft delete**: NO hay deleted_at. DELETE = hard delete (`DELETE FROM alerts WHERE id = ?`). Las queries List NO filtran `deleted_at IS NULL`.
2. **user_id nullable**: En el dominio usar `*int64` (puntero) para user_id. En SQL usar `sql.NullInt64`.
3. **Paginación**: List retorna `([]*Alert, int64, error)`. Query params: `page` (default 1), `limit` (default 15).
4. **Errores de dominio**: var en errors.go, mapeados a HTTP status en handler.
5. **Imports**: `github.com/your-org/jvairv2/...`
6. **Router**: Chi (`github.com/go-chi/chi/v5`), usar `chi.URLParam(r, "id")` para path params.
7. **json.Encode**: SIEMPRE asignar `_ = json.NewEncoder(w).Encode(x)` (errcheck).
8. **rows.Close**: SIEMPRE usar `defer func() { _ = rows.Close() }()` (errcheck).
9. **Swagger**: Annotations completas en cada método. Al finalizar regenerar con:
   `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs`

### Integración (después de crear el módulo)
1. **`cmd/api/container.go`**: Agregar imports, inicializar repository, adapters, use case y handler.
2. **`pkg/rest/router/router.go`**: Agregar parámetro `alertHandler` en `New()` y registrar rutas con `alertHandler.RegisterRoutes(r)`.

## Operaciones esperadas

### Alerts
- **List**: GET /api/v1/alerts
  - Filtros: `userId` (int), `isRead` (bool: true/false), `alertType` (string), `entityType` (string)
  - Paginación: `page` (default 1), `limit` (default 15)
  - Ordenado por created_at DESC
- **Get**: GET /api/v1/alerts/{id} — obtener alerta por ID
- **Create**: POST /api/v1/alerts — crear alerta (uso interno del sistema; también exponer como endpoint REST para que otros módulos puedan crearlas)
- **Mark as Read**: PUT /api/v1/alerts/{id}/read — marcar una alerta individual como leída (is_read=true)
- **Mark All as Read**: PUT /api/v1/alerts/read-all — marcar todas las alertas no leídas del usuario autenticado como leídas
- **Unread Count**: GET /api/v1/alerts/unread-count — retorna `{"count": N}` con el total de alertas no leídas del usuario autenticado (útil para badge de UI)
- **Delete**: DELETE /api/v1/alerts/{id} — hard delete de una alerta

## Notas adicionales
- **user_id es nullable** en la tabla: en Go usar `*int64` en la entidad de dominio y `sql.NullInt64` en el repository para el Scan.
- **entity_type / entity_id** forman una relación polimórfica: `entity_type` indica el tipo de entidad ('job', 'call_log') y `entity_id` es su ID. Almacenar como strings/int64 simples, sin resolver la entidad en este módulo.
- **alert_type** conocido: 'system' — almacenar como string sin validación de enum para permitir extensión futura.
- **message_level** conocido: 'info', 'danger' — almacenar como string sin validación de enum.
- **No existe Update completo** (no se actualiza subject/body/etc.) — solo `MarkAsRead` y `MarkAllRead` son las únicas operaciones de escritura post-creación.
- **Sin soft delete** en ninguna operación — hard delete directo.
- Las alertas son creadas por eventos de otros módulos (job update, job activity log, call_log). Exponer un endpoint POST /api/v1/alerts permite que los handlers de otros módulos creen alertas via HTTP o internamente inyectando el service de alerts.
- Para **Mark All as Read** y **Unread Count**, extraer el user_id del JWT del usuario autenticado (disponible en el middleware de auth existente en el proyecto).
- **Ruta especial**: `PUT /api/v1/alerts/read-all` debe registrarse ANTES de `PUT /api/v1/alerts/{id}/read` en el router para que chi no interprete "read-all" como un ID.

## Entregables de documentación al finalizar

### Documento de resumen
- Crear `docs/MODULE_15_MIGRATION_COMPLETE.md` con: entidades, endpoints implementados, decisiones de diseño, pasos de integración realizados.

### Actualización de Swagger
- Annotations completas en cada handler method.
- Ejecutar `go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs`
- Verificar que todos los endpoints nuevos aparezcan en http://localhost:8090/swagger/index.html

### Colección Postman
Crear `docs/postman_alerts_collection.json` siguiendo el formato de colecciones existentes.

Requisitos:
- Variables de colección: solo `alertId` (int) — **NO** incluir `baseUrl` ni `accessToken` (vienen del environment)
- URLs sin prefijo `/api/v1/` repetido: usar `{{baseUrl}}/alerts`, NO `{{baseUrl}}/api/v1/alerts`
- Auth: Bearer `{{accessToken}}`
- Incluir todas las operaciones: List (con filtros), Get, Create, Mark as Read, Mark All as Read, Unread Count, Delete
- Body de ejemplo para Create:
  ```json
  {
    "userId": 146,
    "alertType": "system",
    "entityId": 36631,
    "entityType": "job",
    "messageLevel": "info",
    "message": "Job WO-36631 has been updated."
  }
  ```
- Query params en List: `page=1&limit=15&isRead=false&alertType=system&entityType=job`

### Verificación final
- `go build ./...` debe compilar sin errores
- `pre-commit run --all-files` debe pasar todos los hooks
```
