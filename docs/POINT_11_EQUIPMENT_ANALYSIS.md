# Punto 11: Equipos (Equipment) — Análisis de Migración

## 1. Estructura de Tablas (desde db_structure.sql)

### property_equipment
```sql
CREATE TABLE `property_equipment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `property_id` bigint unsigned NOT NULL,
  `area` varchar(255) DEFAULT NULL,
  `outdoor_brand` varchar(255) DEFAULT NULL,
  `outdoor_model` varchar(255) DEFAULT NULL,
  `outdoor_serial` varchar(255) DEFAULT NULL,
  `outdoor_installed` timestamp NULL DEFAULT NULL,
  `furnace_brand` varchar(255) DEFAULT NULL,
  `furnace_model` varchar(255) DEFAULT NULL,
  `furnace_serial` varchar(255) DEFAULT NULL,
  `furnace_installed` timestamp NULL DEFAULT NULL,
  `evaporator_brand` varchar(255) DEFAULT NULL,
  `evaporator_model` varchar(255) DEFAULT NULL,
  `evaporator_serial` varchar(255) DEFAULT NULL,
  `evaporator_installed` timestamp NULL DEFAULT NULL,
  `air_handler_brand` varchar(255) DEFAULT NULL,
  `air_handler_model` varchar(255) DEFAULT NULL,
  `air_handler_serial` varchar(255) DEFAULT NULL,
  `air_handler_installed` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT FK property_id -> properties(id)
);
```

**Notas**: NO tiene `deleted_at`. 27 registros en producción.

### job_equipment
```sql
CREATE TABLE `job_equipment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `job_id` bigint unsigned NOT NULL,
  `type` varchar(255) NOT NULL DEFAULT 'current',
  `area` varchar(255) DEFAULT NULL,
  -- mismos campos HVAC que property_equipment --
  `outdoor_brand` ... `outdoor_installed`,
  `furnace_brand` ... `furnace_installed`,
  `evaporator_brand` ... `evaporator_installed`,
  `air_handler_brand` ... `air_handler_installed`,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT FK job_id -> jobs(id)
);
```

**Notas**: NO tiene `deleted_at`. Campo `type` con default 'current'. 3552 registros en producción.

## 2. Análisis de Funcionalidad Original (Laravel)

### PropertyEquipment
- **Modelo**: `PropertyEquipment.php` — No usa SoftDeletes. Usa FilterSortSearchable, Paged.
- **Fillable**: property_id, area, y todos los campos HVAC (16 campos)
- **Casts**: *_installed → Date
- **Validación (StoreRequest/UpdateRequest)**: `area` requerido, resto opcional
- **Controlador**: Solo store, update, destroy (no tiene index ni show independiente)
- **Accessor**: `getOutdoorUnitAttribute()` concatena outdoor_brand + outdoor_model + " | S/N " + outdoor_serial
- **Método especial**: `cloneFromJobEquipment()` — clona equipos de un job a una propiedad (upsert por area)
- **Relación**: belongsTo Property

### JobEquipment
- **Modelo**: `JobEquipment.php` — No usa SoftDeletes. Usa FilterSortSearchable, Paged.
- **Fillable**: job_id, area, type, y todos los campos HVAC (17 campos)
- **Casts**: *_installed → Date
- **Validación**: `area` requerido, `type` requerido + in:current,new
- **Controlador**: Solo store, update, destroy
- **Scopes**: `scopeCurrent()` (type='current'), `scopeNew()` (type='new')
- **Relación**: belongsTo Job

### Rutas (web.php)
- Property Equipment: `properties/{property}/equipment/` → create, update, destroy (anidadas bajo properties)
- Job Equipment: `jobs/{job}/equipment/` → create, update, destroy (anidadas bajo jobs)

## 3. Decisiones de Diseño

### Qué se migra
- CRUD completo para PropertyEquipment (List, Create, Get, Update, Delete)
- CRUD completo para JobEquipment (List, Create, Get, Update, Delete)
- Filtro por `type` (current/new) en JobEquipment List
- Campo calculado `outdoorUnit` en la respuesta

### Qué NO se migra en este punto
- `cloneFromJobEquipment()` — se evaluará como endpoint adicional en una fase posterior o cuando se integre con el flujo de cierre de jobs
- Lógica de searchable (no necesario por ahora, los equipos se listan por parent)

### Decisiones clave
1. **Sin soft deletes**: Ambas tablas NO tienen `deleted_at` en la DB. Se implementa DELETE real (hard delete).
2. **Sub-recursos**: Los equipos son sub-recursos. Las rutas se anidan bajo `/properties/{propertyId}/equipment` y `/jobs/{jobId}/equipment`.
3. **Dos módulos separados**: `property_equipment` y `job_equipment` como dominios independientes, ya que tienen entidades padre diferentes y campos ligeramente distintos (type en job_equipment).
4. **Validación de FK**: Validar que la property/job exista antes de crear/actualizar equipment.

## 4. Plan de Ejecución

### Fase 1: PropertyEquipment
1. Domain: entity.go, repository.go, usecase.go, create.go, get_by_id.go, list.go, update.go, delete.go, mock.go, usecase_test.go
2. Repository: repository.go, create.go, get_by_id.go, list.go, update.go, delete.go
3. Handler: handler.go, create.go, get.go, list.go, update.go, delete.go

### Fase 2: JobEquipment
1. Domain: misma estructura que PropertyEquipment + campo type
2. Repository: misma estructura + filtro por type
3. Handler: misma estructura + filtro por type en List

### Fase 3: Integración
1. container.go — registrar repos, checkers, use cases y handlers
2. router.go — registrar rutas

### Fase 4: Documentación
1. Swagger annotations
2. Colecciones Postman
3. Verificación final

## 5. Endpoints

### Property Equipment
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/properties/{propertyId}/equipment | Listar equipos de una propiedad |
| POST | /api/v1/properties/{propertyId}/equipment | Crear equipo en propiedad |
| GET | /api/v1/properties/{propertyId}/equipment/{id} | Obtener equipo por ID |
| PUT | /api/v1/properties/{propertyId}/equipment/{id} | Actualizar equipo |
| DELETE | /api/v1/properties/{propertyId}/equipment/{id} | Eliminar equipo (hard delete) |

### Job Equipment
| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/jobs/{jobId}/equipment | Listar equipos de un job (filtro por type) |
| POST | /api/v1/jobs/{jobId}/equipment | Crear equipo en job |
| GET | /api/v1/jobs/{jobId}/equipment/{id} | Obtener equipo por ID |
| PUT | /api/v1/jobs/{jobId}/equipment/{id} | Actualizar equipo |
| DELETE | /api/v1/jobs/{jobId}/equipment/{id} | Eliminar equipo (hard delete) |
