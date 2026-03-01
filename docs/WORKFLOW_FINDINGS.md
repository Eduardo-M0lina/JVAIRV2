# Hallazgos: Workflows en JVAIR (Proyecto Base Laravel)

## Resumen Ejecutivo

Los **Workflows** son una pieza central del sistema JVAIR. Representan **secuencias ordenadas de Job Statuses** que definen el flujo de vida de un Job. Cada Customer tiene un workflow por defecto asignado, y cuando se crea un Job, hereda automáticamente el workflow de su Customer.

**El CRUD de workflows ya está migrado a JVAIRV2 (Go)**, pero el frontend necesita entender exactamente **dónde y cómo se usan** los workflows en las pantallas del sistema original.

---

## 1. Estructura de Base de Datos

### Tabla `workflows`

| Columna | Tipo | Notas |
|---------|------|-------|
| `id` | bigint unsigned (PK) | Auto-increment (hasta 21 en producción) |
| `name` | varchar(255) | Nombre del workflow |
| `notes` | text (nullable) | Notas |
| `is_active` | tinyint(1) | Activo/Desactivado (default 1) |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |

> **No tiene soft delete** — los workflows se eliminan permanentemente.

### Tabla pivot `job_status_workflow`

| Columna | Tipo | Notas |
|---------|------|-------|
| `job_status_id` | bigint unsigned (FK → job_statuses) | Status asociado |
| `workflow_id` | bigint unsigned (FK → workflows) | Workflow al que pertenece |
| `order` | int (default 0) | Orden del status dentro del workflow |

> **No tiene PK ni timestamps** — tabla pivot pura.

### Relación con otras tablas

- **`customers.workflow_id`** → FK a `workflows.id` — Cada customer tiene un workflow por defecto
- **`jobs.workflow_id`** → FK a `workflows.id` — Cada job tiene un workflow asignado

---

## 2. Datos en Producción

### Workflows existentes (17 registros):

| ID | Nombre | Activo | Notas |
|----|--------|--------|-------|
| 1 | Z | ❌ No | Desactivado (renombrado a "Z") |
| 2 | Z | ❌ No | Desactivado |
| **3** | **2025 Change Order** | ✅ Sí | Activo |
| **5** | **2025 Contract Installation** | ✅ Sí | Activo |
| 6 | Z | ❌ No | Desactivado |
| 7 | Z | ❌ No | Desactivado |
| **8** | **2025 Work Orders** | ✅ Sí | Activo |
| 9 | Z | ❌ No | Desactivado |
| 10 | 2023 Contract Service | ❌ No | Desactivado |
| 12 | Z | ❌ No | Desactivado |
| 13 | Z | ❌ No | Desactivado |
| 14 | Z | ❌ No | Desactivado |
| 15 | Z | ❌ No | Desactivado |
| 17 | 2023 Builder Installation | ❌ No | Desactivado |
| **18** | **2025 Permits** | ✅ Sí | Activo |
| 19 | Z | ❌ No | Desactivado |
| 20 | Z | ❌ No | Desactivado |

> **Patrón importante:** Los workflows obsoletos se renombran a "Z" y se desactivan en lugar de eliminarlos (porque tienen customers/jobs asociados que impiden el borrado). Solo **4 workflows están activos** en producción.

### Pivot `job_status_workflow`

La tabla pivot tiene **muchos registros** con los statuses asociados a cada workflow. Cada workflow activo tiene entre 7-25+ statuses ordenados.

---

## 3. ¿Dónde se usan los Workflows en el Frontend? (Pantallas)

### 3.1. 📋 Settings → Job Workflows (CRUD Admin)

**Ubicación en sidebar:** `Settings > Job Workflows`

**Acceso controlado por:** Permiso `settings_manage` (verificado via `WorkflowPolicy`)

**Pantallas:**
- **Lista:** `/settings/workflows` → Muestra todos los workflows con nombre, customers asociados, fecha, status activo/desactivado, y opción de duplicar
- **Crear:** `/settings/workflows/create` → Formulario con nombre, checkbox activo, y **drag & drop** para seleccionar/ordenar statuses
- **Editar:** `/settings/workflows/{id}/edit` → Mismo que crear, con datos pre-llenados + opción de eliminar
- **Duplicar:** `/settings/workflows/duplicate/{id}` → Crea copia con nombre "Copy of X (id)"

**UI especial — Drag & Drop:**
- La vista de crear/editar tiene dos contenedores: "Available Statuses" y "Selected Statuses"
- Los job statuses activos se muestran como badges con colores (clase CSS del status)
- Se pueden arrastrar entre contenedores para seleccionar y **reordenar**
- El JavaScript (`resources/js/admin/workflow-creator.js`) maneja el drag & drop con jQuery nativo
- Al soltar un status en "Selected Statuses", genera un `<input type="hidden" name="statuses[]" value="{id}">` para enviar el orden al servidor

### 3.2. 👤 Customer Create/Edit

**Ubicación:** Pantalla de creación/edición de Customer

**Pantallas donde aparece:**
- **Crear Customer** (`/customers/create`): Campo `<select>` "Workflow" — **obligatorio** — lista solo workflows activos
- **Editar Customer** (`/customers/{id}/edit`): Sección "Workflow Information" con `<select>` "Default Workflow" — **obligatorio** — lista solo workflows activos

**Acceso:** Cualquier usuario con permiso de crear/editar customers

**Comportamiento:** Cada customer DEBE tener un workflow asignado. Este workflow se hereda automáticamente cuando se crea un Job para ese customer.

### 3.3. 🔧 Job Create

**Ubicación:** Pantalla de creación de Job (`/jobs/create`)

**Comportamiento clave — El workflow NO se selecciona manualmente al crear un Job:**
- La vista de creación **NO muestra un selector de workflow**
- En el backend (`JobController@store`), el workflow se asigna automáticamente:
  ```php
  $validated['job']['workflow_id'] = $customer->workflow->id;
  $validated['job']['job_status_id'] = $customer->workflow->statuses->first()->id;
  ```
- El workflow viene del customer asociado
- El status inicial del Job es el **primer status** del workflow del customer

### 3.4. 🔧 Job Edit

**Ubicación:** Pantalla de edición de Job (`/jobs/{id}/edit`)

**Campos que aparecen:**
1. **Workflow** (`<select>`) — lista workflows activos — **controlado por permiso `job_edit_work_flow`**
2. **Job Status** (`<select>`) — lista los statuses **del workflow asignado al job** — **controlado por permiso `job_edit_job_status`**

**Comportamiento de permisos:**
- Si el usuario **tiene** `job_edit_work_flow` → selector habilitado, puede cambiar el workflow del job
- Si **no tiene** → selector visible pero **disabled** (solo lectura)
- Si el usuario **tiene** `job_edit_job_status` → selector de status habilitado
- Si **no tiene** → selector visible pero **disabled**

**Lógica importante en update (`JobController@update`):**
```php
if (!$request->user()->can('job_edit_work_flow')) {
    unset($data['workflow_id']);  // Ignora el campo si no tiene permiso
}
```

**Relación workflow ↔ statuses en la vista:**
```php
$statuses = $job->workflow->statuses->pluck('ordered_label', 'id');
```
Los statuses mostrados en el select de "Job Status" son **solo los del workflow asignado**, con formato `"1: Status Name"`, `"2: Status Name"`, etc.

### 3.5. 📊 Job Statuses Table (Settings)

**Ubicación:** `Settings > Job Statuses`

En la tabla de Job Statuses, hay una columna **"Workflows"** que muestra en qué workflows está incluido cada status, con links para ir al workflow correspondiente.

---

## 4. Sistema de Permisos

### 4.1. Para CRUD de Workflows (Settings)

| Acción | Permiso requerido | Roles que lo tienen |
|--------|-------------------|---------------------|
| Ver lista | `settings_manage` | administrator, executive |
| Ver detalle | `settings_manage` | administrator, executive |
| Crear | `settings_manage` | administrator, executive |
| Actualizar | `settings_manage` | administrator, executive |
| Eliminar | `settings_manage` **+ workflow sin customers asociados** | administrator, executive |
| Duplicar | `settings_manage` (usa policy de create) | administrator, executive |

> **Restricción de eliminación:** No se puede eliminar un workflow que tenga customers asociados. Por eso en producción los obsoletos se renombran a "Z" y se desactivan.

### 4.2. Para uso de Workflows en Jobs

| Acción | Permiso requerido | Notas |
|--------|-------------------|-------|
| Cambiar workflow de un Job | `job_edit_work_flow` | Solo admin y executive |
| Cambiar status de un Job | `job_edit_job_status` | Solo admin y executive |

### 4.3. Visibilidad en Sidebar

El menú "Job Workflows" dentro de Settings **solo es visible** si:
1. El usuario tiene permiso `settings_manage`
2. Se verifica con `@can('viewAny', App\Models\Workflow::class)`

---

## 5. Flujo Completo del Workflow

```
[Customer] ←→ [Workflow] ←→ [Job Statuses (ordenados)]
                  ↓
               [Job] → status_actual = uno de los statuses del workflow
```

### Ciclo de vida:

1. **Admin crea Workflow** → Selecciona y ordena Job Statuses via drag & drop
2. **Admin asigna Workflow a Customer** → En el formulario de crear/editar customer
3. **Se crea un Job** → Hereda automáticamente:
   - `workflow_id` del Customer
   - `job_status_id` = primer status del workflow
4. **Job avanza por los statuses** → Usuario con permiso `job_edit_job_status` cambia el status dentro del select (que solo muestra statuses del workflow asignado)
5. **Opcionalmente** → Usuario con `job_edit_work_flow` puede reasignar otro workflow al job

---

## 6. Estado en JVAIRV2 (Go) — Ya Migrado

### ✅ Lo que ya existe:

| Componente | Archivo | Estado |
|-----------|---------|--------|
| **Entity** | `pkg/domain/workflow/entity.go` | ✅ Workflow + WorkflowStatus + Filters |
| **Repository interface** | `pkg/domain/workflow/repository.go` | ✅ List, GetByID, Create, Update, Delete, Duplicate, GetWorkflowStatuses, SetWorkflowStatuses |
| **UseCase** | `pkg/domain/workflow/usecase.go` | ✅ Completo con validación, CRUD, Duplicate |
| **MySQL Repository** | `pkg/repository/mysql/workflow/` | ✅ Todos los métodos implementados |
| **Handler** | `pkg/rest/handler/workflow/handler.go` | ✅ List, Get, Create, Update, Delete, Duplicate |
| **Router** | `pkg/rest/router/workflows.go` | ✅ Todas las rutas API |
| **Tests** | `usecase_test.go`, `repository_test.go`, `handler_test.go` | ✅ Tests unitarios |

### Rutas API migradas:
| Método | Ruta | Acción |
|--------|------|--------|
| GET | `/api/v1/workflows` | Listar workflows |
| POST | `/api/v1/workflows` | Crear workflow |
| GET | `/api/v1/workflows/{id}` | Obtener workflow con statuses |
| PUT | `/api/v1/workflows/{id}` | Actualizar workflow y statuses |
| DELETE | `/api/v1/workflows/{id}` | Eliminar workflow |
| POST | `/api/v1/workflows/{id}/duplicate` | Duplicar workflow |

---

## 7. Lo que el Frontend (JVAIRV2) Necesita Implementar

### 7.1. Pantalla Settings → Job Workflows

| Feature | Descripción | Prioridad |
|---------|-------------|-----------|
| **Lista de workflows** | Tabla con nombre, customers asociados, status activo/desactivado, acciones | Alta |
| **Filtro por activo/desactivado** | Select filter en la tabla | Media |
| **Crear workflow** | Formulario: nombre, is_active, **selector drag & drop de statuses** | Alta |
| **Editar workflow** | Mismo formulario con datos, + botón eliminar | Alta |
| **Duplicar workflow** | Botón en la tabla que llama POST `/{id}/duplicate` | Media |
| **Drag & drop de statuses** | UI para seleccionar y **reordenar** job statuses dentro del workflow | **Alta — es la feature más compleja** |
| **Protección de eliminación** | No permitir eliminar si tiene customers asociados | Alta |

### 7.2. Pantalla Customer Create/Edit

| Feature | Descripción | Prioridad |
|---------|-------------|-----------|
| **Select de Workflow** | Dropdown con workflows activos — **requerido** | Alta |
| **Sección "Workflow Information"** | En el edit, se muestra como sección separada | Media |

### 7.3. Pantalla Job Create

| Feature | Descripción | Prioridad |
|---------|-------------|-----------|
| **Asignación automática** | Al crear job, el backend asigna workflow del customer + primer status | Alta |
| **NO mostrar selector de workflow** | En la creación de Job no se selecciona workflow manualmente | Informativo |

### 7.4. Pantalla Job Edit

| Feature | Descripción | Prioridad |
|---------|-------------|-----------|
| **Select de Workflow** | Dropdown de workflows activos — **habilitado/disabled según permiso `job_edit_work_flow`** | Alta |
| **Select de Job Status** | Dropdown de statuses **del workflow asignado al job** — con formato "N: Status Name" | Alta |
| **Reacción al cambio de workflow** | Si se cambia el workflow, los statuses disponibles deben actualizarse dinámicamente | Alta |

---

## 8. Consideraciones Importantes para el Frontend

### 8.1. El Drag & Drop es la Feature Más Compleja
En el proyecto Laravel se implementó con jQuery nativo (`workflow-creator.js`). Para el frontend nuevo se recomienda usar una librería moderna como `@dnd-kit/core` (React) o similar que soporte reordenamiento.

### 8.2. La Relación Workflow-Status es N:M con Orden
- Un status puede estar en múltiples workflows
- El orden es **por workflow** (no global)
- Al guardar, se envía un array de `statusIDs` en el orden deseado

### 8.3. Los Statuses del Job vienen del Workflow
Cuando se muestra el select de "Job Status" en la edición de un Job, los statuses disponibles son **solo los del workflow asignado**. Si se cambia el workflow, hay que recargar los statuses.

### 8.4. El formato "ordered_label"
Los statuses se muestran con formato `"1: Received"`, `"2: In Progress"`, etc. donde el número es `order + 1` (human-readable).

### 8.5. No se pueden Eliminar Workflows con Customers
La policy verifica `$workflow->customers->isEmpty()` antes de permitir la eliminación. El frontend debe manejar este error gracefully.

### 8.6. Los Workflows Inactivos se renombran a "Z"
Patrón de producción: en vez de eliminar, se renombra a "Z" y se desactiva. El frontend podría mejorar esto con un soft-disable más limpio.

---

## 9. Archivos Relevantes del Proyecto Base (Laravel)

| Archivo | Descripción |
|---------|-------------|
| `app/Models/Workflow.php` | Modelo: name, notes, is_active + relaciones statuses, customers, jobs |
| `app/Models/JobStatus.php` | Modelo con relación workflows (BelongsToMany) y método orderedLabel() |
| `app/Models/Customer.php` | Modelo con `workflow_id` FK y relación BelongsTo workflow |
| `app/Models/Job.php` | Modelo con `workflow_id` FK y relación BelongsTo workflow |
| `app/Http/Controllers/WorkflowController.php` | CRUD completo + duplicate |
| `app/Http/Controllers/JobController.php` | Usa workflow en create (auto-assign) y edit (select + permission) |
| `app/Http/Controllers/CustomerController.php` | Pasa workflows activos a create/edit |
| `app/Policies/WorkflowPolicy.php` | Todas las acciones requieren `settings_manage`, delete requiere sin customers |
| `app/Http/Requests/Workflows/StoreRequest.php` | Validación: name required, statuses required + exists, is_active boolean |
| `resources/views/pages/workflows/` | Vistas: index, create, edit |
| `resources/views/components/workflows/table.blade.php` | Componente tabla con filtro activo/desactivado |
| `resources/views/components/workflows/tabs.blade.php` | Tabs de navegación |
| `resources/views/pages/jobs/partials/section-details.blade.php` | Sección de Job edit con select de workflow y status |
| `resources/views/pages/jobs/create.blade.php` | Vista de crear Job — NO tiene selector de workflow |
| `resources/views/pages/customers/edit.blade.php` | Sección "Workflow Information" con select |
| `resources/views/pages/customers/create.blade.php` | Select de workflow requerido |
| `resources/views/layouts/app.blade.php` | Sidebar: menú de workflows bajo Settings (permiso settings_manage) |
| `resources/js/admin/workflow-creator.js` | JS para drag & drop de statuses |
| `database/migrations/2020_09_20_182444_create_workflows_table.php` | Migración workflows |
| `database/migrations/2020_09_20_182447_create_job_status_workflow_table.php` | Migración tabla pivot |
| `database/seeders/WorkflowSeeder.php` | Seeder: crea workflow "Default" con todos los statuses activos |
| `database/migrations/2023_01_22_211309_add_default_roles.php` | Migración que asigna `job_edit_work_flow` a executive y administrator |
| `resources/views/pages/roles/partials/section-permissions.blade.php` | UI para gestionar permisos de roles incluyendo "Work Flow" |

---

## 10. Archivos Relevantes de JVAIRV2 (Go) — Ya Migrados

| Archivo | Descripción |
|---------|-------------|
| `pkg/domain/workflow/entity.go` | Entity: Workflow, WorkflowStatus, Filters |
| `pkg/domain/workflow/repository.go` | Interface del repositorio |
| `pkg/domain/workflow/usecase.go` | Lógica de negocio completa |
| `pkg/domain/workflow/mock.go` | Mock para tests |
| `pkg/repository/mysql/workflow/` | Implementación MySQL (list, get, create, update, delete, duplicate, statuses) |
| `pkg/rest/handler/workflow/handler.go` | Handler HTTP con request/response structs |
| `pkg/rest/router/workflows.go` | Rutas API |
| `pkg/domain/workflow/usecase_test.go` | Tests del usecase |
| `pkg/repository/mysql/workflow/repository_test.go` | Tests del repositorio |
| `pkg/rest/handler/workflow/handler_test.go` | Tests del handler |
