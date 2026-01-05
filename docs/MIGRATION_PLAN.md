# Plan de Migración Basado en Dependencias

Este documento describe el plan de migración para el proyecto JVAIR, organizado por niveles de dependencia entre módulos.

## Nivel 1: Módulos Base (Sin dependencias externas)

### 1. Administración de Usuarios y Roles
- **Tablas**: `users`, `roles`, `abilities`, `assigned_roles`, `permissions`
- **Descripción**: Este módulo ya está parcialmente implementado con la autenticación, pero falta la gestión completa de usuarios, roles y permisos.
- **Justificación**: Muchas tablas dependen de `users` y `roles`, por lo que es fundamental completar este módulo primero.

### 2. Configuraciones (Settings)
- **Tablas**: `settings`
- **Descripción**: Configuraciones generales del sistema.
- **Justificación**: No tiene dependencias externas y puede ser necesario para otros módulos.

### 3. Flujos de Trabajo (Workflows)
- **Tablas**: `workflows`
- **Descripción**: Define los flujos de trabajo para diferentes procesos.
- **Justificación**: Los clientes y estados de trabajos dependen de los flujos de trabajo.

## Nivel 2: Módulos con Dependencias Simples

### 4. Clientes (Customers)
- **Tablas**: `customers`
- **Dependencias**: `workflows`
- **Descripción**: Gestión de clientes.
- **Justificación**: Depende solo de workflows, pero es requerido por propiedades y trabajos.

### 5. Propiedades (Properties)
- **Tablas**: `properties`
- **Dependencias**: `customers`
- **Descripción**: Gestión de propiedades de clientes.
- **Justificación**: Depende de clientes y es requerido por trabajos.

### 6. Categorías y Estados de Trabajos
- **Tablas**: `job_categories`, `job_statuses`, `job_priorities`, `technician_job_statuses`, `task_statuses`
- **Dependencias**: `workflows` (para job_status_workflow)
- **Descripción**: Catálogos relacionados con trabajos.
- **Justificación**: Son requeridos por el módulo de trabajos.

## Nivel 3: Módulos con Dependencias Múltiples

### 7. Trabajos (Jobs)
- **Tablas**: `jobs`, `job_status_workflow`
- **Dependencias**: `properties`, `job_categories`, `job_statuses`, `job_priorities`, `technician_job_statuses`, `users`
- **Descripción**: Gestión de trabajos/servicios.
- **Justificación**: Depende de múltiples módulos previos y es requerido por varios módulos subsecuentes.

## Nivel 4: Módulos que Dependen de Trabajos

### 8. Cotizaciones (Quotes)
- **Tablas**: `quotes`, `quote_statuses`
- **Dependencias**: `jobs`
- **Descripción**: Gestión de cotizaciones para trabajos.
- **Justificación**: Depende de trabajos.

### 9. Facturas (Invoices)
- **Tablas**: `invoices`, `invoice_payments`
- **Dependencias**: `jobs`
- **Descripción**: Gestión de facturas para trabajos.
- **Justificación**: Depende de trabajos.

### 10. Equipos (Equipment)
- **Tablas**: `property_equipment`, `job_equipment`
- **Dependencias**: `properties`, `jobs`
- **Descripción**: Gestión de equipos asociados a propiedades y trabajos.
- **Justificación**: Depende de propiedades y trabajos.

### 11. Garantías (Warranties)
- **Tablas**: `warranties`, `warranty_claims`, `warranty_claim_types`, `warranty_claim_statuses`, `warranty_statuses`, `warranty_types`, `warranty_equipment`
- **Dependencias**: `jobs`, `equipment`
- **Descripción**: Gestión de garantías para trabajos y equipos.
- **Justificación**: Depende de trabajos y equipos.

## Nivel 5: Módulos Complementarios

### 12. Actividades y Comunicaciones
- **Tablas**: `job_activity_logs`, `job_visits`, `job_residents`, `job_tasks`, `job_rates`, `job_rate_statuses`
- **Dependencias**: `jobs`, `users`, `task_statuses`
- **Descripción**: Registro de actividades y comunicaciones relacionadas con trabajos.
- **Justificación**: Complementa la funcionalidad de trabajos.

### 13. Comunicaciones (Email/SMS)
- **Tablas**: `email_templates`, `sms_templates`, `job_emails`, `job_sms`
- **Dependencias**: `jobs`
- **Descripción**: Plantillas y registros de comunicaciones.
- **Justificación**: Complementa varios módulos con funcionalidad de comunicación.

### 14. Alertas (Alerts)
- **Tablas**: `alerts`
- **Dependencias**: `users` y varias entidades (mediante entity_id y entity_type)
- **Descripción**: Sistema de alertas y notificaciones.
- **Justificación**: Depende de múltiples módulos y proporciona funcionalidad transversal.

## Enfoque de Implementación

Para cada módulo, seguiremos este enfoque:

1. **Entidades de Dominio**: Definir las estructuras de datos y reglas de negocio.
2. **Repositorios**: Implementar la capa de acceso a datos.
3. **Casos de Uso**: Implementar la lógica de negocio.
4. **Endpoints REST**: Implementar la API REST.
5. **Tests**: Crear pruebas unitarias para todas las capas.

## Prioridad de Implementación

1. Completar el módulo de Administración de Usuarios
2. Implementar Configuraciones (Settings)
3. Implementar Flujos de Trabajo (Workflows)
4. Continuar con los siguientes niveles en orden
