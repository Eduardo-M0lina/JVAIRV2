# Módulo 13: Activities and Communications - Funcionalidad y Propósito

## 📋 Descripción General

El **Módulo 13: Activities and Communications** es el sistema de **comunicación, seguimiento y compensación** para trabajos (jobs). Permite registrar actividades, gestionar visitas técnicas, asignar tareas, administrar contactos de residentes y calcular comisiones/pagos para técnicos.

## 🎯 Propósito del Módulo

Este módulo sirve para:

1. **Comunicación y Registro**: Mantener un historial completo de todas las actividades relacionadas con un trabajo
2. **Gestión de Visitas**: Documentar visitas técnicas con reportes y archivos adjuntos
3. **Asignación de Tareas**: Crear y dar seguimiento a tareas específicas del trabajo
4. **Información de Contacto**: Mantener datos de residentes/contactos en la propiedad
5. **Compensación**: Calcular y gestionar pagos/comisiones para técnicos

---

## 📦 Componentes del Módulo

### 1. Job Activity Logs (Registro de Actividades)

**Propósito**: Bitácora de todas las actividades, notas y eventos relacionados con un trabajo.

**Casos de Uso**:
- Registrar llamadas telefónicas con clientes
- Documentar cambios de estado del trabajo
- Notas internas del equipo
- Historial de comunicaciones
- Registro de eventos importantes

**Campos Principales**:
- `job_id`: Trabajo asociado
- `user_id`: Usuario que registra la actividad
- `type`: Tipo de actividad (llamada, nota, cambio de estado, etc.)
- `log`: Descripción detallada de la actividad
- `created_at`: Fecha y hora del registro

**Operaciones**:
- ✅ Crear nueva actividad
- ✅ Listar actividades de un trabajo (con paginación)
- ✅ Eliminar actividad (hard delete)
- ❌ No se permite actualizar (mantiene integridad del historial)

**Ejemplo de Uso**:
```json
POST /api/v1/jobs/123/activities
{
  "userId": 5,
  "type": "phone_call",
  "log": "Cliente confirmó disponibilidad para visita técnica el martes 10am"
}
```

---

### 2. Job Visits (Visitas Técnicas)

**Propósito**: Documentar visitas técnicas al sitio del trabajo con reportes y evidencia fotográfica.

**Casos de Uso**:
- Registrar inspecciones iniciales
- Documentar avance del trabajo
- Adjuntar fotos del antes/después
- Reportes de visitas de supervisión
- Control de acceso (quién puede ver cada visita)

**Campos Principales**:
- `job_id`: Trabajo asociado
- `user_id`: Técnico que realizó la visita
- `date`: Fecha de la visita
- `report`: Reporte detallado de la visita
- `viewable_by`: Array JSON de roles que pueden ver la visita (control de acceso)
- Relación con archivos adjuntos (fotos, documentos)

**Operaciones**:
- ✅ Crear visita con archivos
- ✅ Listar visitas (filtrado por permisos)
- ✅ Actualizar visita
- ✅ Eliminar visita (soft delete)
- ✅ Descargar archivos como ZIP

**Característica Especial**: **Scope Viewable**
- Filtra visitas según el rol del usuario
- Permite control granular de visibilidad

**Ejemplo de Uso**:
```json
POST /api/v1/jobs/123/visits
{
  "userId": 8,
  "date": "2024-03-15",
  "report": "Inspección inicial completada. Se identificaron 3 unidades HVAC que requieren reemplazo.",
  "viewableBy": ["admin", "manager", "technician"],
  "files": [/* archivos adjuntos */]
}
```

---

### 3. Job Residents (Contactos de Residentes)

**Propósito**: Mantener información de contacto de residentes o personas en la propiedad del trabajo.

**Casos de Uso**:
- Guardar datos de contacto del propietario
- Información de inquilinos
- Contactos de emergencia
- Personas autorizadas para acceso
- Coordinación de horarios de visita

**Campos Principales**:
- `job_id`: Trabajo asociado
- `name`: Nombre completo
- `mobile_phone`: Teléfono móvil
- `home_phone`: Teléfono fijo
- `email`: Correo electrónico

**Operaciones**:
- ✅ Crear residente
- ✅ Listar residentes de un trabajo
- ✅ Actualizar información
- ✅ Eliminar residente (soft delete)

**Ejemplo de Uso**:
```json
POST /api/v1/jobs/123/residents
{
  "name": "María González",
  "mobilePhone": "+1-555-0123",
  "homePhone": "+1-555-0124",
  "email": "maria.gonzalez@email.com"
}
```

---

### 4. Job Tasks (Tareas del Trabajo)

**Propósito**: Asignar y dar seguimiento a tareas específicas relacionadas con el trabajo.

**Casos de Uso**:
- Asignar tareas a técnicos específicos
- Seguimiento de pendientes
- Gestión de fechas límite
- Control de estado de tareas
- Vista global de todas las tareas pendientes

**Campos Principales**:
- `job_id`: Trabajo asociado
- `user_id`: Usuario asignado
- `task`: Descripción de la tarea
- `task_status_id`: Estado de la tarea (pendiente, en progreso, completada)
- `due_date`: Fecha límite

**Operaciones**:
- ✅ Crear tarea
- ✅ Listar tareas de un trabajo
- ✅ Listar TODAS las tareas (vista global)
- ✅ Actualizar tarea
- ✅ Eliminar tarea (soft delete)

**Característica Especial**: **Vista Global**
- Endpoint `/api/v1/tasks` lista todas las tareas de todos los trabajos
- Útil para dashboards y gestión de carga de trabajo

**Ejemplo de Uso**:
```json
POST /api/v1/jobs/123/tasks
{
  "userId": 10,
  "task": "Instalar 2 unidades HVAC en segundo piso",
  "taskStatusId": 1,
  "dueDate": "2024-03-20"
}
```

**Funcionalidad de Notificaciones** (Laravel):
- El controlador original incluye método `sendNotification`
- Envía email al usuario asignado
- Registra actividad en el log del trabajo
- **Nota**: Esta funcionalidad de email requiere implementación adicional en Go

---

### 5. Job Rate Statuses (Estados de Comisión)

**Propósito**: Catálogo de estados para las comisiones/pagos de técnicos.

**Casos de Uso**:
- Definir estados: "Pendiente", "Aprobado", "Pagado", "Rechazado"
- Clasificación visual con clases CSS
- Ordenamiento personalizado

**Campos Principales**:
- `label`: Etiqueta del estado (ej: "Pending Approval")
- `class`: Clase CSS para estilo visual (ej: "badge-warning")
- `order`: Orden de visualización

**Operaciones**:
- ✅ Crear estado
- ✅ Listar todos los estados (ordenados)
- ✅ Obtener estado por ID
- ✅ Actualizar estado
- ✅ Eliminar estado (soft delete)

**Ejemplo de Uso**:
```json
POST /api/v1/job-rate-statuses
{
  "label": "Pending Approval",
  "class": "badge-warning",
  "order": 1
}
```

---

### 6. Job Rates (Comisiones/Pagos)

**Propósito**: Calcular y gestionar comisiones y pagos para técnicos basados en el trabajo realizado.

**Casos de Uso**:
- Calcular comisión del técnico
- Registrar pagos
- Gestionar deducciones
- Tracking de partes reemplazadas
- Control de estado de pago

**Campos Principales**:
- `job_id`: Trabajo asociado
- `user_id`: Técnico
- `job_rate_status_id`: Estado del pago
- `sale_price`: Precio de venta del trabajo
- `rate_percent`: Porcentaje de comisión
- `rate_flat`: Comisión fija adicional
- `tech_parts`: Costo de partes del técnico
- `company_parts`: Costo de partes de la empresa
- `deduction`: Deducciones
- `payment`: **Pago calculado automáticamente**
- `paid`: Indicador si ya fue pagado
- `parts_replaced`: Descripción de partes reemplazadas
- `notes`: Notas adicionales

**Fórmula de Cálculo de Pago**:
```
payment = ((sale_price - tech_parts - company_parts) × (rate_percent / 100)) + rate_flat - deduction
```

**Ejemplo**:
- Precio de venta: $1,000
- Partes técnico: $100
- Partes empresa: $150
- Comisión: 15%
- Comisión fija: $50
- Deducción: $20

```
payment = ((1000 - 100 - 150) × 0.15) + 50 - 20
payment = (750 × 0.15) + 50 - 20
payment = 112.50 + 50 - 20
payment = $142.50
```

**Operaciones**:
- ✅ Crear rate (calcula payment automáticamente)
- ✅ Listar rates de un trabajo
- ✅ Actualizar rate (recalcula payment)
- ✅ Eliminar rate (soft delete)
- ✅ **Calcular pago** (endpoint standalone para simulación)

**Ejemplo de Uso**:
```json
POST /api/v1/jobs/123/rates
{
  "userId": 8,
  "jobRateStatusId": 1,
  "salePrice": 1000.00,
  "ratePercent": 15.0,
  "rateFlat": 50.00,
  "techParts": 100.00,
  "companyParts": 150.00,
  "deduction": 20.00,
  "partsReplaced": "2x Compressor units, 1x Thermostat"
}
```

**Endpoint de Cálculo**:
```json
POST /api/v1/calculate-rate-payment
{
  "salePrice": 1000.00,
  "ratePercent": 15.0,
  "rateFlat": 50.00,
  "techParts": 100.00,
  "companyParts": 150.00,
  "deduction": 20.00
}

Response:
{
  "payment": 142.50
}
```

---

## 🔄 Flujo de Trabajo Típico

### Escenario: Trabajo de Instalación HVAC

1. **Creación del Trabajo** (Módulo anterior)
   - Se crea el trabajo #123 para instalación HVAC

2. **Registro de Contacto**
   ```
   POST /jobs/123/residents
   → Guardar info del propietario
   ```

3. **Primera Actividad**
   ```
   POST /jobs/123/activities
   → "Cliente contactado, visita programada para mañana"
   ```

4. **Asignación de Tarea**
   ```
   POST /jobs/123/tasks
   → Asignar inspección inicial al técnico Juan
   ```

5. **Visita Técnica**
   ```
   POST /jobs/123/visits
   → Registrar visita con fotos y reporte
   ```

6. **Actualización de Actividad**
   ```
   POST /jobs/123/activities
   → "Inspección completada, se requieren 3 unidades"
   ```

7. **Nueva Tarea**
   ```
   POST /jobs/123/tasks
   → Asignar instalación al técnico Pedro
   ```

8. **Registro de Comisión**
   ```
   POST /jobs/123/rates
   → Calcular y registrar pago del técnico
   ```

9. **Visita Final**
   ```
   POST /jobs/123/visits
   → Documentar trabajo completado con fotos
   ```

10. **Actividad de Cierre**
    ```
    POST /jobs/123/activities
    → "Trabajo completado, cliente satisfecho"
    ```

---

## 📧 Sistema de Mensajería y Email (Pendiente de Implementación)

### En Laravel (Original)

El sistema Laravel incluye funcionalidad de **notificaciones por email**:

**Ubicación**: `JobTaskController::sendNotification()`

**Funcionalidad**:
1. Envía email al usuario asignado a una tarea
2. Usa el sistema de Mail de Laravel
3. Registra la notificación como actividad en el job
4. Incluye detalles de la tarea en el email

**Código Laravel**:
```php
public function sendNotification(SendNotificationRequest $request, JobTask $jobTask)
{
    // Enviar email
    Mail::to($jobTask->user->email)->send(new TaskNotification($jobTask));

    // Registrar actividad
    JobActivityLog::create([
        'job_id' => $jobTask->job_id,
        'user_id' => auth()->id(),
        'type' => 'notification',
        'log' => "Email notification sent to {$jobTask->user->name}"
    ]);

    return redirect()->back()->with('success', 'Notification sent');
}
```

### ⚠️ En Go (Requiere Implementación)

**Estado Actual**: ❌ No implementado

**Para Implementar**:

1. **Servicio de Email**:
   - Integrar SMTP o servicio como SendGrid/AWS SES
   - Crear templates de email
   - Configurar credenciales en `configs/app.env`

2. **Estructura Sugerida**:
   ```
   pkg/
     notification/
       email/
         service.go          # Servicio de email
         templates.go        # Templates HTML
         task_notification.go # Notificación de tarea
   ```

3. **Agregar al JobTask Handler**:
   ```go
   func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
       // 1. Obtener tarea
       // 2. Obtener usuario asignado
       // 3. Enviar email
       // 4. Registrar actividad
   }
   ```

4. **Endpoint**:
   ```
   POST /api/v1/jobs/{jobId}/tasks/{id}/send-notification
   ```

---

## 🔗 Dependencias con Otros Módulos

### Módulos que USAN este módulo:

1. **Módulo de Jobs** (Principal)
   - Todos los componentes dependen de `job_id`
   - Verificación de existencia de job en cada operación

2. **Módulo de Users**
   - Activities, Visits, Tasks, Rates requieren `user_id`
   - Notificaciones de email usan datos de usuario

3. **Módulo de Task Statuses**
   - JobTasks usa `task_status_id`
   - Catálogo de estados de tareas

4. **Módulo de Files** (para JobVisits)
   - Relación polimórfica con archivos adjuntos
   - Upload y descarga de documentos/fotos

### Módulos que DEPENDEN de este módulo:

1. **Dashboard/Reportes**
   - Estadísticas de actividades
   - Métricas de visitas completadas
   - Tracking de tareas pendientes
   - Reporte de comisiones

2. **Sistema de Notificaciones**
   - Alertas de tareas vencidas
   - Notificaciones de nuevas actividades
   - Recordatorios de visitas programadas

---

## 📊 Endpoints Disponibles

### Job Activities (3 endpoints)
- `POST /api/v1/jobs/{jobId}/activities` - Crear
- `GET /api/v1/jobs/{jobId}/activities` - Listar
- `DELETE /api/v1/jobs/{jobId}/activities/{id}` - Eliminar

### Job Residents (4 endpoints)
- `POST /api/v1/jobs/{jobId}/residents` - Crear
- `GET /api/v1/jobs/{jobId}/residents` - Listar
- `PUT /api/v1/jobs/{jobId}/residents/{id}` - Actualizar
- `DELETE /api/v1/jobs/{jobId}/residents/{id}` - Eliminar

### Job Rate Statuses (5 endpoints)
- `POST /api/v1/job-rate-statuses` - Crear
- `GET /api/v1/job-rate-statuses` - Listar
- `GET /api/v1/job-rate-statuses/{id}` - Obtener
- `PUT /api/v1/job-rate-statuses/{id}` - Actualizar
- `DELETE /api/v1/job-rate-statuses/{id}` - Eliminar

### Job Tasks (6 endpoints)
- `POST /api/v1/jobs/{jobId}/tasks` - Crear
- `GET /api/v1/jobs/{jobId}/tasks` - Listar por job
- `GET /api/v1/tasks` - Listar todas (global)
- `PUT /api/v1/jobs/{jobId}/tasks/{id}` - Actualizar
- `DELETE /api/v1/jobs/{jobId}/tasks/{id}` - Eliminar
- ⚠️ `POST /api/v1/jobs/{jobId}/tasks/{id}/send-notification` - Pendiente

### Job Rates (5 endpoints)
- `POST /api/v1/jobs/{jobId}/rates` - Crear
- `GET /api/v1/jobs/{jobId}/rates` - Listar
- `PUT /api/v1/jobs/{jobId}/rates/{id}` - Actualizar
- `DELETE /api/v1/jobs/{jobId}/rates/{id}` - Eliminar
- `POST /api/v1/calculate-rate-payment` - Calcular pago

### Job Visits (ya implementado)
- Múltiples endpoints con manejo de archivos

**Total**: 29 endpoints

---

## ✅ Checklist de Implementación

- [x] JobActivityLogs - CRUD básico
- [x] JobResidents - CRUD completo
- [x] JobRateStatuses - Catálogo
- [x] JobTasks - CRUD + vista global
- [x] JobRates - CRUD + cálculo de pago
- [x] JobVisits - Ya implementado
- [ ] Sistema de notificaciones por email
- [ ] Templates de email
- [ ] Endpoint SendNotification para tareas
- [ ] Integración con servicio SMTP
- [ ] Tests unitarios
- [ ] Tests de integración

---

## 🎯 Valor de Negocio

Este módulo es **crítico** para:

1. **Transparencia**: Historial completo de comunicaciones
2. **Accountability**: Registro de quién hizo qué y cuándo
3. **Coordinación**: Asignación clara de tareas y responsabilidades
4. **Documentación**: Evidencia fotográfica y reportes de visitas
5. **Compensación Justa**: Cálculo automático y transparente de pagos
6. **Seguimiento**: Vista global del estado de todos los trabajos

Sin este módulo, el sistema sería solo un CRUD de trabajos sin capacidad de gestión operativa real.
