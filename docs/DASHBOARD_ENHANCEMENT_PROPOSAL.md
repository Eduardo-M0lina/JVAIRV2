# Propuesta de Mejora: Dashboard Enriquecido

> **Fecha**: Marzo 2026
> **Estado**: Propuesta
> **Endpoint base**: `GET /api/v1/dashboard`

---

## 1. Estado Actual (Migrado desde Laravel)

El dashboard migrado replica fielmente la funcionalidad original de Laravel:

| Vista | Datos | Fuente |
|-------|-------|--------|
| **Admin** | Jobs pendientes de dispatch (`job_status_id = 1`) | `jobs` |
| **Admin** | Jobs urgentes (`job_status_id = 9`, `job_priority_id = 4`) | `jobs` |
| **Técnico** | Jobs dispatched asignados (`job_status_id = 2`, `user_id = current`) | `jobs` |
| **Técnico** | Jobs urgentes asignados (`job_priority_id = 4`, `user_id = current`) | `jobs` |

### Estructura actual de respuesta

```json
{
  "stats": {
    "jobsAwaitingDispatch": 12,
    "jobsDispatched": 8,
    "jobsUrgent": 3,
    "jobsOpen": 45,
    "jobsClosedThisMonth": 22
  },
  "jobsAwaitingDispatch": [...],
  "jobsUrgent": [...]
}
```

---

## 2. Propuesta de Enriquecimiento

Aprovechando los 41 módulos CRUD ya migrados en Go, se propone agregar **nodos adicionales** a la respuesta del dashboard para convertirlo en un hub central de información. Cada sección adicional es un nodo nuevo en el JSON de respuesta, sin romper la estructura existente.

---

### 2.1 Nodo `invoiceSummary` — Resumen de Facturación

**Tablas**: `invoices`, `invoice_payments`

Provee al admin una vista rápida del estado financiero.

```json
"invoiceSummary": {
  "totalPending": 15,
  "totalPaid": 120,
  "totalOverdue": 4,
  "amountPending": 25430.50,
  "amountPaidThisMonth": 18200.00,
  "recentInvoices": [
    {
      "id": 1,
      "invoiceNumber": "INV-001",
      "jobId": 45,
      "customerName": "Acme Corp",
      "total": 1500.00,
      "status": "sent",
      "createdAt": "2026-03-15T10:00:00Z"
    }
  ]
}
```

**Query base**:
```sql
-- Pending invoices (sin payment completo)
SELECT COUNT(*), SUM(i.total) FROM invoices i
WHERE i.deleted_at IS NULL
  AND i.id NOT IN (SELECT invoice_id FROM invoice_payments WHERE status = 'paid');

-- Paid this month
SELECT COUNT(*), SUM(ip.amount) FROM invoice_payments ip
WHERE ip.created_at >= FIRST_DAY_OF_MONTH AND ip.status = 'paid';
```

**Visibilidad**: Solo admin/executive.

---

### 2.2 Nodo `quoteSummary` — Resumen de Cotizaciones

**Tablas**: `quotes`, `quote_statuses`

Vista rápida de cotizaciones pendientes de aprobación.

```json
"quoteSummary": {
  "totalPending": 8,
  "totalApproved": 45,
  "totalRejected": 12,
  "recentQuotes": [
    {
      "id": 1,
      "quoteNumber": "Q-001",
      "jobId": 30,
      "customerName": "XYZ Inc",
      "total": 3200.00,
      "statusName": "Pending",
      "createdAt": "2026-03-20T14:00:00Z"
    }
  ]
}
```

**Visibilidad**: Solo admin/executive.

---

### 2.3 Nodo `taskSummary` — Tareas Pendientes

**Tablas**: `job_tasks`, `task_statuses`

Muestra tareas asignadas al usuario actual (relevante para ambos roles).

```json
"taskSummary": {
  "totalPending": 5,
  "totalOverdue": 2,
  "tasks": [
    {
      "id": 1,
      "jobId": 45,
      "workOrder": "WO-123",
      "title": "Revisar equipo",
      "dueDate": "2026-03-25T00:00:00Z",
      "statusName": "Pending",
      "isOverdue": true
    }
  ]
}
```

**Query base**:
```sql
SELECT jt.*, js.label as status_name, j.work_order
FROM job_tasks jt
JOIN task_statuses ts ON ts.id = jt.task_status_id
JOIN jobs j ON j.id = jt.job_id
WHERE jt.user_id = ? AND jt.task_status_id NOT IN (/* completed statuses */)
ORDER BY jt.due_date ASC
LIMIT 10;
```

**Visibilidad**: Ambos roles (filtrado por `user_id` para técnicos, todos para admin).

---

### 2.4 Nodo `warrantySummary` — Resumen de Garantías

**Tablas**: `warranties`, `warranty_claims`, `warranty_claim_statuses`

Alerta sobre garantías próximas a vencer y reclamos abiertos.

```json
"warrantySummary": {
  "activeWarranties": 34,
  "expiringThisMonth": 5,
  "openClaims": 3,
  "recentClaims": [
    {
      "id": 1,
      "claimNumber": "WC-001",
      "jobId": 20,
      "customerName": "ABC LLC",
      "statusName": "Open",
      "createdAt": "2026-03-18T09:00:00Z"
    }
  ]
}
```

**Visibilidad**: Solo admin/executive.

---

### 2.5 Nodo `recentActivity` — Actividad Reciente

**Tablas**: `job_activity_logs`

Timeline de actividades recientes del sistema o del usuario.

```json
"recentActivity": [
  {
    "id": 1,
    "jobId": 45,
    "workOrder": "WO-123",
    "type": "status_change",
    "log": "Status changed to Dispatched",
    "userName": "John Doe",
    "createdAt": "2026-03-27T16:30:00Z"
  }
]
```

**Query base**:
```sql
SELECT jal.*, j.work_order, u.name as user_name
FROM job_activity_logs jal
JOIN jobs j ON j.id = jal.job_id
JOIN users u ON u.id = jal.user_id
WHERE jal.created_at >= NOW() - INTERVAL 7 DAY
ORDER BY jal.created_at DESC
LIMIT 15;
```

**Visibilidad**: Admin ve todo; técnicos ven solo actividad de sus jobs.

---

### 2.6 Nodo `alertSummary` — Alertas No Leídas

**Tablas**: `alerts`

Resumen de alertas pendientes para el usuario autenticado.

```json
"alertSummary": {
  "unreadCount": 7,
  "alerts": [
    {
      "id": 1,
      "alertType": "job_overdue",
      "message": "Job WO-123 is overdue",
      "messageLevel": "warning",
      "createdAt": "2026-03-27T10:00:00Z"
    }
  ]
}
```

**Visibilidad**: Ambos roles (filtrado por `user_id`).

---

### 2.7 Nodo `jobsByCategory` — Distribución de Jobs por Categoría

**Tablas**: `jobs`, `job_categories`

Datos para un gráfico de distribución (pie/bar chart) de jobs abiertos por categoría.

```json
"jobsByCategory": [
  { "categoryId": 1, "categoryName": "Maintenance", "count": 15 },
  { "categoryId": 2, "categoryName": "Installation", "count": 8 },
  { "categoryId": 3, "categoryName": "Repair", "count": 22 }
]
```

**Query base**:
```sql
SELECT jc.id, jc.label, COUNT(j.id) as count
FROM jobs j
JOIN job_categories jc ON jc.id = j.job_category_id
WHERE j.deleted_at IS NULL AND j.closed = 0
GROUP BY jc.id, jc.label
ORDER BY count DESC;
```

**Visibilidad**: Solo admin/executive.

---

### 2.8 Nodo `jobsByStatus` — Distribución de Jobs por Estado

**Tablas**: `jobs`, `job_statuses`

Datos para visualizar el pipeline de trabajo.

```json
"jobsByStatus": [
  { "statusId": 1, "statusName": "Awaiting Dispatch", "count": 12 },
  { "statusId": 2, "statusName": "Dispatched", "count": 8 },
  { "statusId": 3, "statusName": "In Progress", "count": 15 }
]
```

**Visibilidad**: Admin ve todos; técnicos ven solo sus jobs.

---

### 2.9 Nodo `jobsDueThisWeek` — Jobs con Vencimiento Esta Semana

**Tablas**: `jobs`

Lista de jobs abiertos cuyo `due_date` cae en la semana actual.

```json
"jobsDueThisWeek": [
  {
    "id": 1,
    "workOrder": "WO-456",
    "dueDate": "2026-03-28T00:00:00Z",
    "customerName": "Delta Corp",
    "propertyStreet": "123 Main St",
    "statusName": "In Progress",
    "priorityName": "High"
  }
]
```

**Visibilidad**: Ambos roles (filtrado por `user_id` para técnicos).

---

### 2.10 Nodo `technicianWorkload` — Carga de Trabajo por Técnico (Solo Admin)

**Tablas**: `jobs`, `users`, `assigned_roles`

Vista para que el admin vea qué tan cargado está cada técnico.

```json
"technicianWorkload": [
  { "userId": 5, "name": "Mike Smith", "openJobs": 8, "urgentJobs": 1 },
  { "userId": 6, "name": "Sarah Jones", "openJobs": 5, "urgentJobs": 0 },
  { "userId": 7, "name": "Bob Wilson", "openJobs": 12, "urgentJobs": 3 }
]
```

**Query base**:
```sql
SELECT u.id, u.name,
  COUNT(CASE WHEN j.closed = 0 THEN 1 END) as open_jobs,
  COUNT(CASE WHEN j.closed = 0 AND j.job_priority_id = 4 THEN 1 END) as urgent_jobs
FROM users u
JOIN assigned_roles ar ON ar.entity_id = u.id AND ar.entity_type = 'App\\Models\\User'
JOIN roles r ON r.id = ar.role_id AND r.name IN ('technician', 'installer')
LEFT JOIN jobs j ON j.user_id = u.id AND j.deleted_at IS NULL
WHERE u.deleted_at IS NULL AND u.is_active = 1
GROUP BY u.id, u.name
ORDER BY open_jobs DESC;
```

**Visibilidad**: Solo admin/executive.

---

## 3. Estructura Final Propuesta de Respuesta

### Admin Dashboard
```json
{
  "stats": { ... },
  "jobsAwaitingDispatch": [...],
  "jobsUrgent": [...],
  "invoiceSummary": { ... },
  "quoteSummary": { ... },
  "taskSummary": { ... },
  "warrantySummary": { ... },
  "recentActivity": [...],
  "alertSummary": { ... },
  "jobsByCategory": [...],
  "jobsByStatus": [...],
  "jobsDueThisWeek": [...],
  "technicianWorkload": [...]
}
```

### Technician Dashboard
```json
{
  "stats": { ... },
  "jobsDispatched": [...],
  "jobsUrgent": [...],
  "taskSummary": { ... },
  "recentActivity": [...],
  "alertSummary": { ... },
  "jobsByStatus": [...],
  "jobsDueThisWeek": [...]
}
```

---

## 4. Plan de Implementación

### Fase A — Quick Wins (bajo esfuerzo, alto valor)
| Nodo | Esfuerzo | Tablas existentes | Requiere nuevos repos |
|------|----------|-------------------|-----------------------|
| `alertSummary` | Bajo | `alerts` | No — repo alert ya existe |
| `taskSummary` | Bajo | `job_tasks`, `task_statuses` | No — repo job_task ya existe |
| `recentActivity` | Bajo | `job_activity_logs` | No — repo job_activity_log ya existe |
| `jobsByCategory` | Bajo | `jobs`, `job_categories` | No — query directa |
| `jobsByStatus` | Bajo | `jobs`, `job_statuses` | No — query directa |
| `jobsDueThisWeek` | Bajo | `jobs` | No — query directa |

### Fase B — Valor de Negocio (esfuerzo medio)
| Nodo | Esfuerzo | Tablas existentes | Requiere nuevos repos |
|------|----------|-------------------|-----------------------|
| `invoiceSummary` | Medio | `invoices`, `invoice_payments` | No — repos ya existen |
| `quoteSummary` | Medio | `quotes`, `quote_statuses` | No — repos ya existen |
| `warrantySummary` | Medio | `warranties`, `warranty_claims` | No — repos ya existen |
| `technicianWorkload` | Medio | `jobs`, `users`, `assigned_roles` | No — query directa con JOINs |

### Fase C — Opcional / Futuro
- Gráficos de tendencia temporal (jobs creados/cerrados por semana)
- Revenue tracking (sum de `job_sales_price` por mes)
- Performance metrics por técnico (tiempo promedio de cierre)

---

## 5. Consideraciones Técnicas

### Performance
- Todas las queries del dashboard deben ejecutarse en **paralelo** usando goroutines para minimizar latencia.
- Considerar **caching** con TTL de 30-60 segundos para stats que no necesitan ser real-time.
- Los conteos (`COUNT`) son queries ligeras sobre índices existentes.

### Implementación sugerida
```go
// Ejecutar queries en paralelo
var wg sync.WaitGroup
var stats *DashboardStats
var invoiceSummary *InvoiceSummary
// ... más variables

wg.Add(N)
go func() { defer wg.Done(); stats, _ = repo.GetStats(ctx, userID) }()
go func() { defer wg.Done(); invoiceSummary, _ = repo.GetInvoiceSummary(ctx) }()
// ...
wg.Wait()
```

### Compatibilidad
- Los nodos nuevos son **aditivos**: el frontend existente sigue funcionando sin cambios.
- Cada nodo nuevo puede ser habilitado/deshabilitado con query params opcionales:
  - `GET /api/v1/dashboard?include=invoices,quotes,tasks`
  - Sin params → retorna solo la vista básica migrada de Laravel.

---

## 6. Resumen de Endpoints Existentes Reutilizables

Todos estos módulos ya están migrados y sus repositorios pueden ser reutilizados directamente:

| Módulo | Endpoint | Datos para Dashboard |
|--------|----------|---------------------|
| Jobs | `GET /api/v1/jobs` | Conteos por estado, categoría, prioridad |
| Invoices | `GET /api/v1/invoices` | Resumen de facturación |
| Invoice Payments | `GET /api/v1/invoices/{id}/payments` | Montos pagados |
| Quotes | `GET /api/v1/quotes` | Cotizaciones pendientes |
| Job Tasks | `GET /api/v1/tasks` | Tareas asignadas |
| Warranties | `GET /api/v1/warranties` | Garantías activas |
| Warranty Claims | `GET /api/v1/warranty-claims` | Reclamos abiertos |
| Job Activity Logs | `GET /api/v1/jobs/{id}/activities` | Actividad reciente |
| Alerts | `GET /api/v1/alerts` | Alertas no leídas |
| Users | `GET /api/v1/users` | Carga de trabajo por técnico |
| Job Categories | `GET /api/v1/job-categories` | Distribución por categoría |
| Job Statuses | `GET /api/v1/job-statuses` | Pipeline de trabajo |
| Job Rates | `GET /api/v1/jobs/{id}/rates` | Info de payroll (futuro) |
