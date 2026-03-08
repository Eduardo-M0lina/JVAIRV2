# Módulo 13: Análisis de Dependencias con Módulos Anteriores

## 🔍 Resumen Ejecutivo

Después de analizar el código Laravel original, **el Módulo 13 tiene dependencias críticas con otros módulos**, especialmente el **Módulo de Jobs**. Varios controladores crean automáticamente registros de `JobActivityLog` cuando ocurren eventos importantes.

---

## ⚠️ DEPENDENCIAS CRÍTICAS ENCONTRADAS

### 1. **JobController → JobActivityLog** (ALTA PRIORIDAD)

El controlador de Jobs en Laravel crea automáticamente activity logs en múltiples operaciones:

#### 📍 Ubicación Laravel: `/JVAIR/app/Http/Controllers/JobController.php`

#### Eventos que Registran Actividades:

| Método | Línea | Evento | Tipo de Log |
|--------|-------|--------|-------------|
| `store()` | 181 | Job creado | `job_created` |
| `update()` | 389 | Job actualizado | `job_update` |
| `close()` | 432 | Job cerrado | `job_closed` |
| `dispatchEmail()` | 500 | Email enviado | `job_email_dispatched` |
| `dispatchSupervisorEmail()` | 540 | Email supervisor | `job_email_dispatched` |
| `dispatchSms()` | 603 | SMS enviado | `job_sms_dispatched` |
| `attemptedCall()` | 647 | Llamada registrada | `job_attempted_call` |

#### Código Laravel Ejemplo:

```php
// Al crear un Job
JobActivityLog::create([
    'job_id' => $job->id,
    'type' => 'job_created',
    'log' => 'Job was created',
    'user_id' => Auth::user()->id,
]);

// Al actualizar un Job
JobActivityLog::create([
    'job_id' => $job->id,
    'type' => 'job_update',
    'log' => 'Job was updated',
    'user_id' => Auth::user()->id,
]);

// Al cerrar un Job
JobActivityLog::create([
    'job_id' => $job->id,
    'type' => 'job_closed',
    'log' => 'Job was closed',
    'user_id' => Auth::user()->id,
]);
```

#### ✅ ACCIÓN REQUERIDA EN GO:

**Archivo**: `/Users/eduardo/projects/jvair/JVAIRV2/pkg/rest/handler/job/handler.go`

**Modificaciones Necesarias**:

1. **Importar el servicio de JobActivityLog**:
```go
import (
    jobActivityLogDomain "github.com/your-org/jvairv2/pkg/domain/job_activity_log"
)
```

2. **Agregar al Handler struct**:
```go
type Handler struct {
    service            job.Service
    activityLogService jobActivityLogDomain.Service  // NUEVO
}
```

3. **Actualizar constructor**:
```go
func NewHandler(service job.Service, activityLogService jobActivityLogDomain.Service) *Handler {
    return &Handler{
        service:            service,
        activityLogService: activityLogService,
    }
}
```

4. **Registrar actividades en cada operación**:

```go
// En Create()
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // ... crear job ...

    // Registrar actividad
    activityLog := &jobActivityLogDomain.JobActivityLog{
        JobID:  job.ID,
        UserID: userID, // del contexto/auth
        Type:   "job_created",
        Log:    "Job was created",
    }
    h.activityLogService.Create(r.Context(), activityLog)

    // ... respuesta ...
}

// En Update()
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    // ... actualizar job ...

    activityLog := &jobActivityLogDomain.JobActivityLog{
        JobID:  jobID,
        UserID: userID,
        Type:   "job_update",
        Log:    "Job was updated",
    }
    h.activityLogService.Create(r.Context(), activityLog)

    // ... respuesta ...
}

// En Close() - si existe
func (h *Handler) Close(w http.ResponseWriter, r *http.Request) {
    // ... cerrar job ...

    activityLog := &jobActivityLogDomain.JobActivityLog{
        JobID:  jobID,
        UserID: userID,
        Type:   "job_closed",
        Log:    "Job was closed",
    }
    h.activityLogService.Create(r.Context(), activityLog)

    // ... respuesta ...
}
```

---

### 2. **JobTaskController → JobActivityLog + Email** (ALTA PRIORIDAD)

#### 📍 Ubicación Laravel: `/JVAIR/app/Http/Controllers/JobTaskController.php`

#### Funcionalidad: `sendNotification()` (Líneas 100-122)

**Qué hace**:
1. Envía email al usuario asignado a la tarea
2. Registra el email en `job_emails` table
3. **Crea un JobActivityLog** indicando que se envió la notificación

```php
public function sendNotification(TaskNotificationRequest $request, Job $job, JobTask $task): RedirectResponse
{
    // 1. Enviar email
    Mail::to($task->user->email)->send(new TaskNotificationEmail($task));

    // 2. Registrar email
    JobEmail::create([
        'job_id' => $task->job_id,
        'recipients' => $task->user->email,
        'type' => 'task_notification',
    ]);

    // 3. Registrar actividad
    JobActivityLog::create([
        'job_id' => $task->job_id,
        'type' => 'task_notification_sent',
        'log' => "An Email Notification was sent for task '{$task->task}' to user {$task->user->name}.",
        'user_id' => Auth::user()->id,
    ]);

    return redirect()->route('jobs.edit', $task->job)->with('success', 'Task Email Notification has been sent.');
}
```

#### ✅ ACCIÓN REQUERIDA EN GO:

**Archivo**: `/Users/eduardo/projects/jvair/JVAIRV2/pkg/rest/handler/job_task/handler.go`

**Agregar nuevo endpoint**:

```go
// SendNotification godoc
// @Summary Send task notification email
// @Description Send email notification to user assigned to task and log activity
// @Tags JobTasks
// @Accept json
// @Produce json
// @Param jobId path int true "Job ID"
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/jobs/{jobId}/tasks/{id}/send-notification [post]
// @Security BearerAuth
func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    taskID, _ := strconv.ParseInt(vars["id"], 10, 64)

    // 1. Obtener tarea
    task, err := h.service.GetByID(r.Context(), taskID)
    if err != nil {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }

    // 2. Obtener usuario (necesitas servicio de usuarios)
    user, err := h.userService.GetByID(r.Context(), task.UserID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // 3. Enviar email (necesitas implementar servicio de email)
    err = h.emailService.SendTaskNotification(user.Email, task)
    if err != nil {
        http.Error(w, "Failed to send email", http.StatusInternalServerError)
        return
    }

    // 4. Registrar en job_emails (si existe ese módulo)
    // ... código para registrar email ...

    // 5. Registrar actividad
    activityLog := &jobActivityLogDomain.JobActivityLog{
        JobID:  task.JobID,
        UserID: currentUserID, // del contexto/auth
        Type:   "task_notification_sent",
        Log:    fmt.Sprintf("An Email Notification was sent for task '%s' to user %s.", task.Task, user.Name),
    }
    h.activityLogService.Create(r.Context(), activityLog)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Task notification sent successfully",
    })
}
```

**Actualizar Handler struct**:
```go
type Handler struct {
    service            job_task.Service
    activityLogService jobActivityLogDomain.Service  // NUEVO
    userService        userDomain.Service            // NUEVO
    emailService       email.Service                 // NUEVO (pendiente implementar)
}
```

---

### 3. **JobActivityController → Alert** (PRIORIDAD MEDIA)

#### 📍 Ubicación Laravel: `/JVAIR/app/Http/Controllers/JobActivityController.php` (Línea 36-42)

Cuando se crea una actividad, opcionalmente puede crear una **Alert**:

```php
$note = JobActivityLog::create($data);

if ($create_alert) {
    Alert::create([
        'job_id' => $note->job_id,
        'user_id' => $note->user_id,
        'message' => $note->log,
    ]);
}
```

#### ✅ ACCIÓN REQUERIDA EN GO:

**Verificar si el módulo de Alerts ya fue migrado**. Si existe:

1. Agregar campo opcional `createAlert` en el request de JobActivityLog
2. Si es `true`, crear también un Alert

---

### 4. **Otros Módulos con Email** (PRIORIDAD BAJA)

Estos módulos envían emails pero **NO registran en JobActivityLog**:

| Módulo | Controlador | Método | Email |
|--------|-------------|--------|-------|
| Invoices | InvoiceController | `email()` | Envía factura por email |
| Quotes | QuoteController | `email()` | Envía cotización por email |
| Payroll | PayrollController | - | Envía recibo de pago |

**Nota**: Estos NO requieren integración con Módulo 13.

---

## 📋 CHECKLIST DE IMPLEMENTACIÓN

### Fase 1: JobController (CRÍTICO)
- [ ] Verificar si JobController ya está migrado en Go
- [ ] Agregar `activityLogService` al handler de Jobs
- [ ] Implementar registro de actividad en `Create()`
- [ ] Implementar registro de actividad en `Update()`
- [ ] Implementar registro de actividad en `Close()` (si existe)
- [ ] Implementar registro de actividad en `DispatchEmail()` (si existe)
- [ ] Implementar registro de actividad en `DispatchSMS()` (si existe)
- [ ] Implementar registro de actividad en `AttemptedCall()` (si existe)

### Fase 2: JobTaskController (CRÍTICO)
- [ ] Implementar servicio de Email en Go
- [ ] Crear templates de email para notificaciones
- [ ] Agregar endpoint `SendNotification` en JobTask handler
- [ ] Integrar con servicio de JobActivityLog
- [ ] Agregar ruta en router: `POST /jobs/{jobId}/tasks/{id}/send-notification`
- [ ] Actualizar Swagger con nuevo endpoint
- [ ] Actualizar colección Postman

### Fase 3: Sistema de Email (NUEVO MÓDULO)
- [ ] Crear paquete `pkg/notification/email/`
- [ ] Implementar servicio SMTP o integración con SendGrid/AWS SES
- [ ] Crear templates HTML para emails
- [ ] Configurar credenciales en `configs/app.env`
- [ ] Implementar `SendTaskNotification()`
- [ ] Implementar otros tipos de email según necesidad

### Fase 4: Integración con Alerts (OPCIONAL)
- [ ] Verificar si módulo de Alerts existe
- [ ] Agregar campo `createAlert` en JobActivityLog request
- [ ] Implementar lógica de creación de Alert

---

## 🔧 ARCHIVOS QUE NECESITAN MODIFICACIÓN

### Archivos Existentes a Modificar:

1. **`/pkg/rest/handler/job/handler.go`**
   - Agregar `activityLogService`
   - Registrar actividades en operaciones CRUD

2. **`/pkg/rest/handler/job_task/handler.go`**
   - Agregar `SendNotification()` método
   - Agregar dependencias de email y activityLog

3. **`/cmd/api/container.go`**
   - Inyectar `activityLogService` en JobHandler
   - Inyectar servicios adicionales en JobTaskHandler

4. **`/pkg/rest/router/router.go`**
   - Agregar ruta para `SendNotification`

### Archivos Nuevos a Crear:

1. **`/pkg/notification/email/service.go`**
   - Servicio de envío de emails

2. **`/pkg/notification/email/templates.go`**
   - Templates HTML para emails

3. **`/pkg/notification/email/task_notification.go`**
   - Email específico para notificación de tareas

---

## 📊 IMPACTO Y PRIORIDAD

| Dependencia | Impacto | Prioridad | Esfuerzo |
|-------------|---------|-----------|----------|
| JobController → ActivityLog | **ALTO** | 🔴 CRÍTICO | 2-3 horas |
| JobTask → Email + ActivityLog | **ALTO** | 🔴 CRÍTICO | 4-6 horas |
| Sistema de Email | **MEDIO** | 🟡 ALTO | 6-8 horas |
| Alerts Integration | **BAJO** | 🟢 MEDIO | 2-3 horas |

**Total Estimado**: 14-20 horas de desarrollo

---

## 🎯 RECOMENDACIONES

1. **Implementar primero la integración con JobController**
   - Es la más crítica y la más simple
   - Mantiene el historial de auditoría de jobs

2. **Luego implementar el sistema de Email**
   - Necesario para JobTask notifications
   - Reutilizable para otros módulos

3. **Finalmente agregar SendNotification a JobTasks**
   - Depende del sistema de email
   - Completa la funcionalidad del Módulo 13

4. **Considerar implementación gradual**
   - Fase 1: Solo logging de actividades (sin email)
   - Fase 2: Agregar sistema de email
   - Fase 3: Integración completa

---

## ⚠️ NOTAS IMPORTANTES

1. **El Módulo 13 NO es standalone**
   - Requiere integración activa con otros módulos
   - Especialmente con el módulo de Jobs

2. **Sistema de Email es crítico**
   - Actualmente NO implementado en Go
   - Necesario para funcionalidad completa de notificaciones

3. **Auditoría y Compliance**
   - Los activity logs son esenciales para auditoría
   - Deben registrarse de forma consistente en todas las operaciones

4. **Autenticación**
   - Todos los logs requieren `user_id` del usuario autenticado
   - Necesitas extraer el user ID del contexto/token JWT

---

## 📝 PRÓXIMOS PASOS INMEDIATOS

1. **Verificar estado del JobController en Go**
   ```bash
   # Buscar el handler de jobs
   find /Users/eduardo/projects/jvair/JVAIRV2 -name "*job*handler.go" -type f
   ```

2. **Revisar si existe módulo de Email**
   ```bash
   # Buscar implementación de email
   find /Users/eduardo/projects/jvair/JVAIRV2 -name "*email*" -o -name "*mail*"
   ```

3. **Verificar módulo de Alerts**
   ```bash
   # Buscar implementación de alerts
   find /Users/eduardo/projects/jvair/JVAIRV2 -name "*alert*"
   ```

4. **Crear issue/ticket para cada dependencia**
   - Separar en tareas independientes
   - Priorizar según impacto

---

## ✅ CONCLUSIÓN

El **Módulo 13 está implementado** pero **NO está completamente integrado** con los módulos que lo necesitan. Se requiere:

1. ✅ Módulo 13 implementado (domain, repository, handlers)
2. ❌ Integración con JobController (registrar actividades)
3. ❌ Sistema de Email (nuevo módulo completo)
4. ❌ Endpoint SendNotification en JobTasks
5. ❌ Integración con Alerts (si existe)

**Sin estas integraciones, el sistema pierde funcionalidad crítica de auditoría y notificaciones.**
