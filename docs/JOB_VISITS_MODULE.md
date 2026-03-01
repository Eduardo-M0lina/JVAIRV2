# Módulo de Visitas de Trabajo (Job Visits) + Archivos

## Descripción General

Este módulo gestiona las **visitas de trabajo** realizadas por técnicos a los sitios de los clientes, así como los **archivos adjuntos** (fotos, documentos, etc.) asociados a cada visita. Los archivos se almacenan en **AWS S3** y se registran en la base de datos con una relación polimórfica.

---

## Estructura de Datos

### Job Visit (`job_visits`)

| Campo        | Tipo        | Descripción                                    |
|-------------|-------------|------------------------------------------------|
| `id`        | bigint (PK) | Identificador único                            |
| `job_id`    | bigint (FK) | Referencia al trabajo (`jobs`)                 |
| `user_id`   | bigint (FK) | Referencia al usuario/técnico (`users`)        |
| `viewable_by`| varchar    | JSON con IDs de usuarios que pueden ver la visita |
| `date`      | timestamp   | Fecha de la visita                             |
| `report`    | text        | Reporte o descripción de la visita             |
| `created_at`| timestamp   | Fecha de creación                              |
| `updated_at`| timestamp   | Fecha de última actualización                  |
| `deleted_at`| timestamp   | Soft delete                                    |

### File (`files`)

| Campo          | Tipo        | Descripción                                    |
|---------------|-------------|------------------------------------------------|
| `id`          | bigint (PK) | Identificador único                            |
| `type`        | varchar     | MIME type del archivo (ej: `image/jpeg`)       |
| `path`        | varchar     | Clave (key) del archivo en S3                  |
| `url`         | varchar     | URL pública del archivo en S3                  |
| `fileable_id` | int         | ID de la entidad asociada (polimórfico)        |
| `fileable_type`| varchar    | Tipo de entidad asociada (`App\Models\JobVisit`)|
| `created_at`  | timestamp   | Fecha de creación                              |
| `updated_at`  | timestamp   | Fecha de última actualización                  |

> **Nota:** La tabla `files` usa una relación polimórfica. El campo `fileable_type` contiene el nombre completo del modelo Laravel (`App\Models\JobVisit`) para mantener compatibilidad con el sistema existente.

---

## Endpoints

Base URL: `/api/v1`

### Visitas de Trabajo

#### `GET /jobs/{jobId}/visits`
Lista paginada de visitas de un trabajo específico. Incluye los archivos adjuntos de cada visita.

**Query Parameters:**
| Parámetro   | Tipo   | Default      | Descripción                     |
|------------|--------|-------------|----------------------------------|
| `page`     | int    | 1           | Número de página                 |
| `pageSize` | int    | 10          | Registros por página             |
| `search`   | string | -           | Buscar en el campo `report`      |
| `userId`   | int    | -           | Filtrar por usuario              |
| `sort`     | string | `created_at`| Campo de ordenamiento (`date`, `created_at`) |
| `direction`| string | `DESC`      | Dirección (`ASC`, `DESC`)        |

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": 1,
      "jobId": 100,
      "userId": 5,
      "viewableBy": "[\"1\",\"2\"]",
      "date": "01-15-2024",
      "report": "Inspección del sistema HVAC completada.",
      "files": [
        {
          "id": 10,
          "type": "image",
          "url": "https://bucket.s3.amazonaws.com/uploads/1234_photo.jpg",
          "createdAt": "2024-01-15T10:30:00Z"
        }
      ],
      "createdAt": "2024-01-15T10:00:00Z",
      "updatedAt": "2024-01-15T10:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 10,
  "totalItems": 1,
  "totalPages": 1
}
```

---

#### `GET /jobs/{jobId}/visits/{visitId}`
Obtiene una visita específica con sus archivos adjuntos.

**Response:** `200 OK`
```json
{
  "id": 1,
  "jobId": 100,
  "userId": 5,
  "date": "01-15-2024",
  "report": "Inspección completada.",
  "files": [],
  "createdAt": "2024-01-15T10:00:00Z",
  "updatedAt": "2024-01-15T10:00:00Z"
}
```

---

#### `POST /jobs/{jobId}/visits`
Crea una nueva visita de trabajo.

**Body (JSON):**
```json
{
  "userId": 5,
  "viewableBy": "[\"1\",\"2\"]",
  "date": "01-15-2024",
  "report": "Visita de inspección al sitio."
}
```

| Campo        | Requerido | Descripción                           |
|-------------|-----------|---------------------------------------|
| `userId`    | Sí        | ID del técnico que realizó la visita  |
| `viewableBy`| No        | JSON array con IDs de usuarios        |
| `date`      | No        | Fecha en formato `MM-DD-YYYY`. Si no se envía, se usa la fecha actual |
| `report`    | No        | Texto del reporte                     |

**Response:** `201 Created`

---

#### `PUT /jobs/{jobId}/visits/{visitId}`
Actualiza una visita existente.

**Body (JSON):**
```json
{
  "userId": 5,
  "viewableBy": "[\"1\",\"2\",\"3\"]",
  "date": "01-16-2024",
  "report": "Reporte actualizado después de la visita de seguimiento."
}
```

**Response:** `200 OK`

---

#### `DELETE /jobs/{jobId}/visits/{visitId}`
Elimina una visita (soft delete). **También elimina todos los archivos adjuntos** de S3 y de la base de datos.

**Response:** `204 No Content`

---

### Archivos de Visitas

#### `GET /jobs/{jobId}/visits/{visitId}/files`
Lista los archivos adjuntos de una visita.

**Response:** `200 OK`
```json
[
  {
    "id": 10,
    "type": "image",
    "url": "https://bucket.s3.amazonaws.com/uploads/1234_photo.jpg",
    "createdAt": "2024-01-15T10:30:00Z"
  }
]
```

---

#### `POST /jobs/{jobId}/visits/{visitId}/files`
Sube uno o más archivos a una visita. Usa `multipart/form-data`.

**Form Data:**
| Campo   | Tipo | Descripción                                      |
|---------|------|--------------------------------------------------|
| `files` | file | Archivos a subir (máx. 40MB c/u, máx. 20 archivos) |

**Response:** `201 Created`
```json
[
  {
    "id": 11,
    "type": "image",
    "url": "https://bucket.s3.amazonaws.com/uploads/1709312345_foto.jpg",
    "createdAt": "2024-01-15T11:00:00Z"
  }
]
```

---

#### `GET /jobs/{jobId}/visits/{visitId}/files/{fileId}/download`
Descarga un archivo de S3. Retorna el archivo binario con los headers:
- `Content-Type`: MIME type del archivo
- `Content-Disposition`: `attachment; filename="nombre_archivo.ext"`

---

#### `DELETE /jobs/{jobId}/visits/{visitId}/files/{fileId}`
Elimina un archivo de S3 y su registro de la base de datos (hard delete).

**Response:** `204 No Content`

---

## Configuración de S3

Configurar las siguientes variables de entorno en `configs/app.env`:

```env
AWS_DEFAULT_REGION=us-east-1
AWS_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_ENDPOINT=                    # Opcional, para S3-compatible (MinIO, DigitalOcean Spaces, etc.)
```

> **Nota:** Si las credenciales de S3 no son válidas o no se configuran, el módulo de archivos seguirá funcionando pero las operaciones de upload/download/delete fallarán con errores de S3.

---

## Arquitectura

```
pkg/
├── common/storage/
│   └── s3.go                    # Cliente S3 (Upload, Delete, GetObject)
├── domain/
│   ├── job_visit/
│   │   ├── entity.go            # Entidad JobVisit
│   │   ├── errors.go            # Errores del dominio
│   │   ├── repository.go        # Interfaz Repository + Checkers
│   │   ├── usecase.go           # Service interface + UseCase struct
│   │   ├── create.go            # Lógica de creación
│   │   ├── get_by_id.go         # Obtener por ID
│   │   ├── list.go              # Listado paginado
│   │   ├── update.go            # Actualización
│   │   └── delete.go            # Soft delete
│   └── file/
│       ├── entity.go            # Entidad File (polimórfica)
│       ├── errors.go            # Errores del dominio
│       ├── repository.go        # Interfaz Repository
│       ├── usecase.go           # Service + StorageService interfaces
│       ├── upload.go            # Upload a S3 + crear registro
│       ├── get_by_id.go         # Obtener por ID
│       ├── list.go              # Listar por fileable
│       ├── download.go          # Descargar de S3
│       └── delete.go            # Eliminar de S3 + BD
├── repository/mysql/
│   ├── job_visit/
│   │   ├── repository.go        # Constructor
│   │   ├── create.go            # INSERT
│   │   ├── get_by_id.go         # SELECT by ID
│   │   ├── list.go              # SELECT paginado con filtros
│   │   ├── update.go            # UPDATE
│   │   ├── delete.go            # Soft DELETE
│   │   └── adapters.go          # JobExistsChecker, UserExistsChecker
│   └── file/
│       ├── repository.go        # Constructor
│       ├── create.go            # INSERT
│       ├── get_by_id.go         # SELECT by ID
│       ├── list.go              # SELECT by fileable
│       └── delete.go            # Hard DELETE
└── rest/handler/job_visit/
    └── handler.go               # HTTP handlers (CRUD visitas + archivos)
```

---

## Colección Postman

Importar el archivo `docs/postman_job_visits_collection.json` en Postman para probar todos los endpoints.

**Variables de colección:**
- `baseUrl`: URL base del servidor (default: `http://localhost:8080`)
- `token`: Token JWT de autenticación
- `jobId`: ID del trabajo
- `visitId`: ID de la visita
- `fileId`: ID del archivo
