# Módulo de Garantías (Warranties)

Módulo migrado desde Laravel (JVAIR) a Go (JVAIRV2). Gestiona el ciclo de vida completo de garantías HVAC: creación, seguimiento de equipos, reclamos y sus catálogos asociados.

---

## Índice

1. [Visión general](#visión-general)
2. [Arquitectura](#arquitectura)
3. [Entidades y relaciones](#entidades-y-relaciones)
4. [Catálogos](#catálogos)
5. [Endpoints](#endpoints)
6. [Filtros, paginación y ordenamiento](#filtros-paginación-y-ordenamiento)
7. [Validaciones](#validaciones)
8. [Soft deletes](#soft-deletes)
9. [Estructura de archivos](#estructura-de-archivos)
10. [Postman](#postman)

---

## Visión general

El módulo permite a la empresa:

- **Registrar garantías** asociadas a trabajos (jobs) realizados, con tipo, estado, número de acuerdo y fecha de envío.
- **Gestionar equipos HVAC** vinculados a cada garantía (outdoor unit, furnace, evaporator, air handler), incluyendo clonado automático desde equipos de trabajo.
- **Registrar reclamos de garantía** con datos de piezas, fabricante, números de serie, estado de aprobación y pagos recibidos.
- **Administrar catálogos** de tipos y estados tanto para garantías como para reclamos.

Todos los endpoints requieren autenticación (`Bearer token`) y están disponibles bajo el prefijo `/api/v1`.

---

## Arquitectura

El módulo sigue la arquitectura por capas del proyecto:

```
Domain (entidad, validación, interfaz de repositorio, interfaz de servicio, use case)
   ↓
Repository (implementación MySQL con queries, adapters para FK checks)
   ↓
Handler (REST HTTP con chi, request/response structs, Swagger annotations)
   ↓
Router + Container (inyección de dependencias y registro de rutas)
```

**Patrones clave:**

- **Adapters**: Para verificar existencia de foreign keys (job, warranty_type, warranty_status, etc.) sin acoplar dominios.
- **Soft delete**: Warranties y Warranty Claims usan `deleted_at` en lugar de borrado físico.
- **Hard delete**: Warranty Equipment y catálogos se eliminan físicamente (con validación de dependencias en catálogos).
- **Logging**: Todos los use cases usan `log/slog` para trazabilidad.

---

## Entidades y relaciones

```
Jobs (existente)
 ├── Warranties            (N:1 con Job, soft delete)
 │    └── Warranty Equipment  (N:1 con Warranty, hard delete)
 └── Warranty Claims       (N:1 con Job, soft delete)

Warranty Types          → referenciado por Warranties
Warranty Statuses       → referenciado por Warranties
Warranty Claim Types    → referenciado por Warranty Claims
Warranty Claim Statuses → referenciado por Warranty Claims
```

---

## Catálogos

### Warranty Types (`/api/v1/warranty-types`)

Tipos de garantía (ej: "Parts", "Labor", "Parts & Labor").

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `label` | string | ✅ | Nombre del tipo |
| `labelPlural` | string | ✅ | Nombre en plural |
| `isActive` | bool | ❌ | Activo (default: `true`) |

**Endpoints:** `GET /` · `POST /` · `GET /{id}` · `PUT /{id}` · `DELETE /{id}`

> El DELETE valida que no existan garantías usando ese tipo antes de eliminar.

---

### Warranty Statuses (`/api/v1/warranty-statuses`)

Estados de garantía (ej: "Pending", "In Progress", "Completed").

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `label` | string | ✅ | Nombre del estado |
| `class` | string | ❌ | Clase CSS para UI (ej: `badge-warning`) |
| `order` | int | ✅ | Orden de visualización |
| `isActive` | bool | ❌ | Activo (default: `true`) |

**Endpoints:** `GET /` · `POST /` · `GET /{id}` · `PUT /{id}` · `DELETE /{id}`

> El DELETE valida que no existan garantías usando ese estado antes de eliminar.

---

### Warranty Claim Types (`/api/v1/warranty-claim-types`)

Tipos de reclamo de garantía (ej: "Parts Warranty", "Labor Warranty").

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `label` | string | ✅ | Nombre del tipo |
| `labelPlural` | string | ✅ | Nombre en plural |

**Endpoints:** `GET /` · `POST /` · `GET /{id}` · `PUT /{id}` · `DELETE /{id}`

> El DELETE valida que no existan reclamos usando ese tipo.

---

### Warranty Claim Statuses (`/api/v1/warranty-claim-statuses`)

Estados de reclamo (ej: "Pending Review", "Approved", "Denied").

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `label` | string | ✅ | Nombre del estado |
| `class` | string | ❌ | Clase CSS para UI |
| `order` | int | ✅ | Orden de visualización |
| `isActive` | bool | ❌ | Activo (default: `true`) |

**Endpoints:** `GET /` · `POST /` · `GET /{id}` · `PUT /{id}` · `DELETE /{id}`

> El DELETE valida que no existan reclamos usando ese estado.

---

## Endpoints

### Warranties (`/api/v1/warranties`)

Gestiona las garantías asociadas a trabajos (jobs).

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/warranties` | Lista paginada con filtros |
| `POST` | `/warranties` | Crear garantía |
| `GET` | `/warranties/{id}` | Obtener por ID |
| `PUT` | `/warranties/{id}` | Actualizar |
| `DELETE` | `/warranties/{id}` | Eliminar (soft delete) |

#### Campos del body (Create)

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `warrantyNumber` | string | ✅ | Número único de garantía |
| `jobId` | int64 | ✅ | ID del trabajo asociado |
| `warrantyTypeId` | int64 | ✅ | FK → warranty_types |
| `warrantyStatusId` | int64 | ✅ | FK → warranty_statuses |
| `dateSubmitted` | string | ❌ | Fecha de envío (formato `MM-DD-YYYY`) |
| `agreementNumber` | string | ❌ | Número de acuerdo |
| `auditDone` | bool | ❌ | Auditoría realizada (default: `false`) |
| `notes` | string | ❌ | Notas |

#### Campos del body (Update)

Igual que Create pero sin `jobId` (no se puede cambiar el trabajo asociado).

#### Ejemplo de request (Create)

```json
POST /api/v1/warranties
{
    "warrantyNumber": "WRN-2024-001",
    "jobId": 1,
    "warrantyTypeId": 1,
    "warrantyStatusId": 1,
    "dateSubmitted": "01-15-2024",
    "agreementNumber": "AGR-2024-001",
    "auditDone": false,
    "notes": "New warranty for HVAC system replacement"
}
```

---

### Warranty Equipment (`/api/v1/warranties/{warrantyId}/equipment`)

Sub-recurso de garantía. Gestiona los equipos HVAC asociados. Cada equipo tiene 4 unidades opcionales: outdoor, furnace, evaporator y air handler.

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/warranties/{warrantyId}/equipment` | Listar equipos de una garantía |
| `POST` | `/warranties/{warrantyId}/equipment` | Crear equipo |
| `PUT` | `/warranties/{warrantyId}/equipment/{equipmentId}` | Actualizar equipo |
| `DELETE` | `/warranties/{warrantyId}/equipment/{equipmentId}` | Eliminar equipo (hard delete) |

#### Campos del body

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `area` | string | ✅ | Área de instalación (ej: "Main Floor") |
| `outdoorBrand` | string | ❌ | Marca unidad exterior |
| `outdoorModel` | string | ❌ | Modelo unidad exterior |
| `outdoorSerial` | string | ❌ | Serial unidad exterior |
| `outdoorInstalled` | string | ❌ | Fecha instalación (`YYYY-MM-DD`) |
| `furnaceBrand` | string | ❌ | Marca horno |
| `furnaceModel` | string | ❌ | Modelo horno |
| `furnaceSerial` | string | ❌ | Serial horno |
| `furnaceInstalled` | string | ❌ | Fecha instalación (`YYYY-MM-DD`) |
| `evaporatorBrand` | string | ❌ | Marca evaporador |
| `evaporatorModel` | string | ❌ | Modelo evaporador |
| `evaporatorSerial` | string | ❌ | Serial evaporador |
| `evaporatorInstalled` | string | ❌ | Fecha instalación (`YYYY-MM-DD`) |
| `airHandlerBrand` | string | ❌ | Marca manejadora de aire |
| `airHandlerModel` | string | ❌ | Modelo manejadora de aire |
| `airHandlerSerial` | string | ❌ | Serial manejadora de aire |
| `airHandlerInstalled` | string | ❌ | Fecha instalación (`YYYY-MM-DD`) |

#### Ejemplo de request (Create)

```json
POST /api/v1/warranties/1/equipment
{
    "area": "Main Floor",
    "outdoorBrand": "Carrier",
    "outdoorModel": "24ACC636A003",
    "outdoorSerial": "1234567890",
    "outdoorInstalled": "2024-01-15",
    "furnaceBrand": "Lennox",
    "furnaceModel": "ML180UH",
    "furnaceSerial": "0987654321",
    "furnaceInstalled": "2024-01-15"
}
```

> **Nota:** El módulo también implementa lógica de clonado (`CloneFromJobEquipment`) que permite copiar los equipos de un trabajo directamente a la garantía.

---

### Warranty Claims (`/api/v1/warranty-claims`)

Gestiona los reclamos de garantía. Cada reclamo está asociado a un trabajo y contiene información detallada sobre la pieza reclamada, fabricante, números de serie y estado de aprobación/pago.

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/warranty-claims` | Lista paginada con filtros |
| `POST` | `/warranty-claims` | Crear reclamo |
| `GET` | `/warranty-claims/{id}` | Obtener por ID |
| `PUT` | `/warranty-claims/{id}` | Actualizar |
| `DELETE` | `/warranty-claims/{id}` | Eliminar (soft delete) |

#### Campos del body (Create)

| Campo | Tipo | Requerido | Descripción |
|-------|------|-----------|-------------|
| `internalClaimNumber` | string | ✅ | Número interno de reclamo |
| `warrantyClaimTypeId` | int64 | ✅ | FK → warranty_claim_types |
| `warrantyClaimStatusId` | int64 | ✅ | FK → warranty_claim_statuses |
| `jobId` | int64 | ✅ | FK → jobs |
| `invoiceNumber` | string | ❌ | Número de factura |
| `workDone` | bool | ❌ | Trabajo realizado (default: `false`) |
| `warrantyPart` | string | ❌ | Pieza en garantía |
| `manufacturer` | string | ❌ | Fabricante |
| `modelNumber` | string | ❌ | Número de modelo |
| `partNumber` | string | ❌ | Número de pieza |
| `replacementPartNumber` | string | ❌ | Número de pieza de reemplazo |
| `partDistributor` | string | ❌ | Distribuidor de piezas |
| `partInvoiceNumber` | string | ❌ | Número de factura de pieza |
| `oldPartSerialNumber` | string | ❌ | Serial de pieza original |
| `newPartSerialNumber` | string | ❌ | Serial de pieza nueva |
| `esaNumber` | string | ❌ | Número ESA |
| `serial` | string | ❌ | Serial del equipo |
| `claimNumber` | string | ❌ | Número de reclamo (del fabricante) |
| `approved` | bool | ❌ | Aprobado (default: `false`) |
| `partsCreditReceived` | bool | ❌ | Crédito de piezas recibido (default: `false`) |
| `laborPaymentReceived` | bool | ❌ | Pago de mano de obra recibido (default: `false`) |
| `notes` | string | ❌ | Notas |

#### Campos del body (Update)

Igual que Create pero sin `jobId` (no se puede cambiar el trabajo asociado).

#### Ejemplo de request (Create)

```json
POST /api/v1/warranty-claims
{
    "internalClaimNumber": "ICN-2024-001",
    "warrantyClaimTypeId": 1,
    "warrantyClaimStatusId": 1,
    "jobId": 1,
    "invoiceNumber": "INV-2024-001",
    "workDone": false,
    "warrantyPart": "Compressor",
    "manufacturer": "Carrier",
    "modelNumber": "24ACC636A003",
    "partNumber": "PN-12345",
    "replacementPartNumber": "PN-67890",
    "partDistributor": "HD Supply",
    "partInvoiceNumber": "PI-2024-001",
    "oldPartSerialNumber": "OLD-SN-001",
    "newPartSerialNumber": "NEW-SN-001",
    "esaNumber": "ESA-2024-001",
    "serial": "SER-2024-001",
    "claimNumber": "CLM-2024-001",
    "approved": false,
    "partsCreditReceived": false,
    "laborPaymentReceived": false,
    "notes": "Initial claim for compressor replacement"
}
```

---

## Filtros, paginación y ordenamiento

### Warranties — Query params

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `page` | int | Número de página (default: 1) |
| `pageSize` | int | Tamaño de página (default: 10) |
| `search` | string | Busca en `warranty_number`, `agreement_number`, `notes` |
| `jobId` | int64 | Filtrar por trabajo |
| `warrantyTypeId` | int64 | Filtrar por tipo de garantía |
| `warrantyStatusId` | int64 | Filtrar por estado de garantía |
| `weekNumber` | string | Filtrar por semana del trabajo (JOIN con `jobs`) |
| `sort` | string | `warranty_number`, `date_submitted`, `created_at`, `week_number` |
| `direction` | string | `ASC` o `DESC` |

### Warranty Claims — Query params

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `page` | int | Número de página (default: 1) |
| `pageSize` | int | Tamaño de página (default: 10) |
| `search` | string | Busca en `internal_claim_number`, `claim_number`, `notes` |
| `jobId` | int64 | Filtrar por trabajo |
| `warrantyClaimTypeId` | int64 | Filtrar por tipo de reclamo |
| `warrantyClaimStatusId` | int64 | Filtrar por estado de reclamo |
| `weekNumber` | string | Filtrar por semana del trabajo (JOIN con `jobs`) |
| `sort` | string | `internal_claim_number`, `claim_number`, `created_at`, `week_number` |
| `direction` | string | `ASC` o `DESC` |

### Catálogos — Query params

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `page` | int | Número de página |
| `pageSize` | int | Tamaño de página |
| `search` | string | Busca en `label` |
| `isActive` | bool | Filtrar por estado activo (solo en tipos y estados) |

### Formato de respuesta paginada

```json
{
    "data": [...],
    "page": 1,
    "pageSize": 15,
    "total": 42
}
```

---

## Validaciones

### Foreign keys

Antes de crear o actualizar, el use case verifica la existencia de las FK mediante **adapters**:

- **Warranties**: `jobId` → jobs, `warrantyTypeId` → warranty_types, `warrantyStatusId` → warranty_statuses
- **Warranty Claims**: `jobId` → jobs, `warrantyClaimTypeId` → warranty_claim_types, `warrantyClaimStatusId` → warranty_claim_statuses
- **Warranty Equipment**: `warrantyId` → warranties

Si la FK no existe, se retorna un error `400 Bad Request`.

### Campos requeridos

- **Warranty**: `warrantyNumber`, `jobId`, `warrantyTypeId`, `warrantyStatusId`
- **Warranty Claim**: `internalClaimNumber`, `warrantyClaimTypeId`, `warrantyClaimStatusId`, `jobId`
- **Warranty Equipment**: `warrantyId` (del path), `area`
- **Warranty Type**: `label`, `labelPlural`
- **Warranty Status**: `label`, `order`
- **Warranty Claim Type**: `label`, `labelPlural`
- **Warranty Claim Status**: `label`, `order`

### Catálogos — Protección contra borrado

Los catálogos no se pueden eliminar si tienen registros dependientes. Ejemplo: no se puede borrar un `warranty_type` si existen garantías que lo usen.

---

## Soft deletes

Las entidades **Warranties** y **Warranty Claims** usan soft delete (columna `deleted_at`):

- `DELETE` marca el registro con `deleted_at = NOW()` en lugar de eliminarlo.
- `GET` y `LIST` filtran automáticamente registros con `deleted_at IS NULL`.
- `UPDATE` sobre un registro eliminado retorna `410 Gone`.

Las entidades **Warranty Equipment** y **catálogos** usan **hard delete** (eliminación física).

---

## Estructura de archivos

```
pkg/
├── domain/
│   ├── warranty/                  # Entidad Warranty
│   │   ├── entity.go             # Struct + validación
│   │   ├── repository.go         # Interfaz de repositorio
│   │   ├── usecase.go            # Interfaz de servicio + struct UseCase
│   │   ├── errors.go             # Errores de dominio
│   │   ├── create.go             # UC: crear
│   │   ├── get_by_id.go          # UC: obtener por ID
│   │   ├── list.go               # UC: listar con filtros
│   │   ├── update.go             # UC: actualizar
│   │   └── delete.go             # UC: soft delete
│   ├── warranty_equipment/        # Entidad Warranty Equipment
│   │   ├── entity.go
│   │   ├── repository.go
│   │   ├── usecase.go
│   │   ├── errors.go
│   │   ├── create.go
│   │   ├── get_by_id.go
│   │   ├── list.go
│   │   ├── update.go
│   │   ├── delete.go
│   │   └── clone.go              # UC: clonar desde job_equipment
│   ├── warranty_claim/            # Entidad Warranty Claim
│   │   ├── entity.go
│   │   ├── repository.go
│   │   ├── usecase.go
│   │   ├── errors.go
│   │   ├── create.go
│   │   ├── get_by_id.go
│   │   ├── list.go
│   │   ├── update.go
│   │   └── delete.go
│   ├── warranty_type/             # Catálogo: tipos de garantía
│   ├── warranty_status/           # Catálogo: estados de garantía
│   ├── warranty_claim_type/       # Catálogo: tipos de reclamo
│   └── warranty_claim_status/     # Catálogo: estados de reclamo
├── repository/mysql/
│   ├── warranty/                  # MySQL repo + adapters
│   ├── warranty_equipment/        # MySQL repo + adapters + helpers + clone
│   ├── warranty_claim/            # MySQL repo + adapters + helpers
│   ├── warranty_type/
│   ├── warranty_status/
│   ├── warranty_claim_type/
│   └── warranty_claim_status/
└── rest/handler/
    ├── warranty/handler.go        # REST + Swagger
    ├── warranty_equipment/handler.go
    ├── warranty_claim/handler.go
    ├── warranty_type/handler.go
    ├── warranty_status/handler.go
    ├── warranty_claim_type/handler.go
    └── warranty_claim_status/handler.go

cmd/api/
├── container.go                   # DI: repos, adapters, use cases, handlers
└── main.go                        # Swagger tags

docs/
├── postman_warranties_collection.json  # Colección Postman completa
└── WARRANTIES_MODULE.md                # Este archivo
```

---

## Postman

La colección Postman se encuentra en `docs/postman_warranties_collection.json`. Incluye carpetas organizadas por entidad:

1. **Warranties** — CRUD completo con filtros
2. **Warranty Equipment** — CRUD como sub-recurso de warranty
3. **Warranty Claims** — CRUD completo con filtros
4. **Warranty Types** — CRUD de catálogo
5. **Warranty Statuses** — CRUD de catálogo
6. **Warranty Claim Types** — CRUD de catálogo
7. **Warranty Claim Statuses** — CRUD de catálogo

### Variables requeridas

| Variable | Valor por defecto | Descripción |
|----------|-------------------|-------------|
| `baseUrl` | `http://localhost:8080/api/v1` | URL base de la API |
| `accessToken` | — | Token JWT obtenido del login (`POST /auth/login`) |

### Swagger UI

Disponible en `http://localhost:8090/swagger/index.html` una vez que el servidor está corriendo. Todos los endpoints tienen anotaciones Swagger completas.
