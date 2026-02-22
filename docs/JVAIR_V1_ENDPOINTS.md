# JVAIR V1 - Listado Completo de Endpoints

> Proyecto Laravel (JVAIR). Documentación generada a partir del archivo `routes/web.php`.
> **Nota:** El archivo `routes/api.php` no tiene rutas definidas.

---

## Tabla de Contenidos

1. [Rutas Públicas (sin autenticación)](#1-rutas-públicas-sin-autenticación)
2. [Autenticación](#2-autenticación)
3. [Dashboard](#3-dashboard)
4. [Alertas](#4-alertas)
5. [Tareas](#5-tareas)
6. [Cuenta de Usuario](#6-cuenta-de-usuario)
7. [Búsqueda](#7-búsqueda)
8. [Clientes (Customers)](#8-clientes-customers)
9. [Supervisores](#9-supervisores)
10. [Propiedades (Properties)](#10-propiedades-properties)
11. [Equipos de Propiedad (Property Equipment)](#11-equipos-de-propiedad-property-equipment)
12. [Trabajos (Jobs)](#12-trabajos-jobs)
13. [Equipos de Trabajo (Job Equipment)](#13-equipos-de-trabajo-job-equipment)
14. [Residentes de Trabajo (Job Residents)](#14-residentes-de-trabajo-job-residents)
15. [Visitas de Trabajo (Job Visits)](#15-visitas-de-trabajo-job-visits)
16. [Tareas de Trabajo (Job Tasks)](#16-tareas-de-trabajo-job-tasks)
17. [Tarifas de Trabajo (Job Rates)](#17-tarifas-de-trabajo-job-rates)
18. [Actividades/Notas de Trabajo (Job Activities)](#18-actividadesnotas-de-trabajo-job-activities)
19. [Reclamaciones de Garantía (Warranty Claims)](#19-reclamaciones-de-garantía-warranty-claims)
20. [Garantías (Warranties)](#20-garantías-warranties)
21. [Equipos de Garantía (Warranty Equipment)](#21-equipos-de-garantía-warranty-equipment)
22. [Nómina (Payroll)](#22-nómina-payroll)
23. [Configuración General (Settings)](#23-configuración-general-settings)
24. [Estados de Trabajo (Job Statuses)](#24-estados-de-trabajo-job-statuses)
25. [Estados de Tarea (Task Statuses)](#25-estados-de-tarea-task-statuses)
26. [Estados de Trabajo de Técnico (Tech Job Statuses)](#26-estados-de-trabajo-de-técnico-tech-job-statuses)
27. [Flujos de Trabajo (Workflows)](#27-flujos-de-trabajo-workflows)
28. [Categorías de Trabajo (Job Categories)](#28-categorías-de-trabajo-job-categories)
29. [Prioridades de Trabajo (Job Priorities)](#29-prioridades-de-trabajo-job-priorities)
30. [Estados de Reclamación de Garantía (Warranty Claim Statuses)](#30-estados-de-reclamación-de-garantía-warranty-claim-statuses)
31. [Estados de Garantía (Warranty Statuses)](#31-estados-de-garantía-warranty-statuses)
32. [Tipos de Garantía (Warranty Types)](#32-tipos-de-garantía-warranty-types)
33. [Plantillas de Email (Email Templates)](#33-plantillas-de-email-email-templates)
34. [Plantillas de Texto (Text Templates)](#34-plantillas-de-texto-text-templates)
35. [Usuarios](#35-usuarios)
36. [Roles](#36-roles)
37. [Archivos (Files)](#37-archivos-files)
38. [Cálculo de Pagos](#38-cálculo-de-pagos)
39. [Endpoints Comentados (Deshabilitados)](#39-endpoints-comentados-deshabilitados)

---

## Resumen de Conteo

| Sección | Endpoints Activos |
|---|---|
| Rutas Públicas | 3 |
| Autenticación | 13 |
| Dashboard | 2 |
| Alertas | 3 |
| Tareas | 1 |
| Cuenta de Usuario | 3 |
| Búsqueda | 8 |
| Clientes | 8 |
| Supervisores | 6 |
| Propiedades | 7 |
| Equipos de Propiedad | 3 |
| Trabajos | 13 |
| Equipos de Trabajo | 3 |
| Residentes de Trabajo | 3 |
| Visitas de Trabajo | 4 |
| Tareas de Trabajo | 4 |
| Tarifas de Trabajo | 3 |
| Actividades de Trabajo | 2 |
| Reclamaciones de Garantía | 6 |
| Garantías | 6 |
| Equipos de Garantía | 3 |
| Nómina | 6 |
| Configuración General | 2 |
| Estados de Trabajo | 6 |
| Estados de Tarea | 6 |
| Estados de Trabajo Técnico | 6 |
| Workflows | 7 |
| Categorías de Trabajo | 6 |
| Prioridades de Trabajo | 6 |
| Estados Reclamación Garantía | 6 |
| Estados de Garantía | 6 |
| Tipos de Garantía | 6 |
| Plantillas de Email | 6 |
| Plantillas de Texto | 6 |
| Usuarios | 9 |
| Roles | 6 |
| Archivos | 3 |
| Cálculo de Pagos | 1 |
| **TOTAL** | **~196** |

---

## 1. Rutas Públicas (sin autenticación)

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/` | `home` | `SiteController@home` | Página principal del sitio (frontend) |
| `GET` | `/invoices/{invoice}/pay` | `invoices.pay` | `InvoiceController@pay` | Página pública de pago online de una factura |
| `POST` | `/stripe/webhook` | `stripe.webhook` | `StripeWebhookController@webhook` | Webhook de Stripe para procesar eventos de pago |

---

## 2. Autenticación

Prefijo: `/auth`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/auth/login` | `auth.login` | `LoginController@show` | Mostrar formulario de login |
| `POST` | `/auth/login` | `auth.login` | `LoginController@login` | Procesar inicio de sesión |
| `GET/POST` | `/auth/logout` | `auth.logout` | `LoginController@logout` | Cerrar sesión |
| `GET` | `/auth/password/reset` | `auth.password.request` | `ForgotPasswordController@show` | Mostrar formulario de solicitud de reset de contraseña |
| `POST` | `/auth/password/email` | `auth.password.email` | `ForgotPasswordController@request` | Enviar email de reset de contraseña |
| `GET` | `/auth/password/reset/{token}` | `auth.password.reset` | `ResetPasswordController@show` | Mostrar formulario de reset con token |
| `POST` | `/auth/password/reset` | `auth.password.update` | `ResetPasswordController@reset` | Procesar reset de contraseña |
| `GET` | `/auth/password/reset-password` | `auth.password.enforceReset` | `PasswordSecurityController@show` | Mostrar formulario de cambio forzado de contraseña |
| `POST` | `/auth/password/reset-password` | `auth.password.enforceUpdate` | `PasswordSecurityController@resetPassword` | Procesar cambio forzado de contraseña |
| `GET` | `/auth/password/confirm` | `auth.password.confirm` | `ConfirmPasswordController@show` | Mostrar formulario de confirmación de contraseña |
| `POST` | `/auth/password/confirm` | — | `ConfirmPasswordController@confirm` | Procesar confirmación de contraseña |
| `GET` | `/auth/email/verify` | `auth.verification.notice` | `VerificationController@show` | Mostrar aviso de verificación de email |
| `GET` | `/auth/email/verify/{id}/{hash}` | `auth.verification.verify` | `VerificationController@verify` | Verificar email del usuario |
| `POST` | `/auth/email/resend` | `auth.verification.resend` | `VerificationController@resend` | Reenviar email de verificación |

---

## 3. Dashboard

**Middleware:** `auth`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/dashboard` | `dashboard` | `DashboardController@dashboard` | Panel principal del sistema |
| `GET` | `/dashboard/urgent` | `dashboard.urgent` | `DashboardController@urgent` | Vista de trabajos/tareas urgentes |

---

## 4. Alertas

**Middleware:** `auth`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/alerts` | `alerts` | `AlertController@alerts` | Listar todas las alertas |
| `GET` | `/alerts/{alert}/open` | `alerts.open` | `AlertController@openAlert` | Abrir/ver una alerta específica |
| `GET` | `/alerts/mark-call-log/{job}` | `alerts.mark-call-log` | `AlertController@markCallLogAsRead` | Marcar registro de llamada como leído para un trabajo |

---

## 5. Tareas

**Middleware:** `auth`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/tasks` | `tasks.index` | `JobTaskController@index` | Listar todas las tareas (vista global) |

---

## 6. Cuenta de Usuario

**Middleware:** `auth` — Prefijo: `/account`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/account` | `account.edit` | `AccountController@edit` | Mostrar formulario de edición de cuenta propia |
| `PUT` | `/account` | `account.update` | `AccountController@update` | Actualizar datos de cuenta propia |
| `POST` | `/account/toggle-sidebar-state` | `account.toggle-sidebar-state` | `AccountController@toggleSidebarState` | Alternar estado del sidebar (abierto/cerrado) |

---

## 7. Búsqueda

**Middleware:** `auth`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/search/jobs` | `search.jobs` | `SearchController@jobs` | Buscar trabajos |
| `GET` | `/search/customers` | `search.customers` | `SearchController@customers` | Buscar clientes |
| `GET` | `/search/properties` | `search.properties` | `SearchController@properties` | Buscar propiedades |
| `GET` | `/search/invoices` | `search.invoices` | `SearchController@invoices` | Buscar facturas |
| `GET` | `/search/quotes` | `search.quotes` | `SearchController@quotes` | Buscar cotizaciones |
| `GET` | `/search/warranties` | `search.warranties` | `SearchController@warranties` | Buscar garantías |
| `GET` | `/search/claims` | `search.claims` | `SearchController@claims` | Buscar reclamaciones |
| `GET` | `/search/users` | `search.users` | `SearchController@users` | Buscar usuarios |

---

## 8. Clientes (Customers)

**Middleware:** `auth` — Prefijo: `/customers`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/customers` | `customers.index` | `CustomerController@index` | Listar todos los clientes |
| `GET` | `/customers/create` | `customers.create` | `CustomerController@create` | Mostrar formulario de creación de cliente |
| `POST` | `/customers/create` | `customers.store` | `CustomerController@store` | Crear un nuevo cliente |
| `GET` | `/customers/{customer}/edit` | `customers.edit` | `CustomerController@edit` | Mostrar formulario de edición de cliente |
| `GET` | `/customers/{customer}/properties` | `customers.properties` | `CustomerController@properties` | Ver propiedades de un cliente |
| `GET` | `/customers/{customer}/jobs` | `customers.jobs` | `CustomerController@jobs` | Ver trabajos de un cliente |
| `GET` | `/customers/{customer}/supervisors` | `customers.supervisors` | `CustomerController@supervisors` | Ver supervisores de un cliente |
| `PUT` | `/customers/{customer}` | `customers.update` | `CustomerController@update` | Actualizar datos de un cliente |
| `DELETE` | `/customers/{customer}` | `customers.destroy` | `CustomerController@destroy` | Eliminar un cliente |

---

## 9. Supervisores

**Middleware:** `auth` — Prefijo: `/supervisors`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/supervisors` | `supervisors.index` | `SupervisorController@index` | Listar todos los supervisores |
| `GET` | `/supervisors/create` | `supervisors.create` | `SupervisorController@create` | Mostrar formulario de creación |
| `POST` | `/supervisors/create` | `supervisors.store` | `SupervisorController@store` | Crear un nuevo supervisor |
| `GET` | `/supervisors/{supervisor}/edit` | `supervisors.edit` | `SupervisorController@edit` | Mostrar formulario de edición |
| `PUT` | `/supervisors/{supervisor}` | `supervisors.update` | `SupervisorController@update` | Actualizar datos de un supervisor |
| `DELETE` | `/supervisors/{supervisor}` | `supervisors.destroy` | `SupervisorController@destroy` | Eliminar un supervisor |

---

## 10. Propiedades (Properties)

**Middleware:** `auth` — Prefijo: `/properties`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/properties` | `properties.index` | `PropertyController@index` | Listar todas las propiedades |
| `GET` | `/properties/create` | `properties.create` | `PropertyController@create` | Mostrar formulario de creación |
| `POST` | `/properties/create` | `properties.store` | `PropertyController@store` | Crear una nueva propiedad |
| `GET` | `/properties/{property}/edit` | `properties.edit` | `PropertyController@edit` | Mostrar formulario de edición |
| `GET` | `/properties/{property}/jobs` | `properties.jobs` | `PropertyController@jobs` | Ver trabajos asociados a una propiedad |
| `PUT` | `/properties/{property}` | `properties.update` | `PropertyController@update` | Actualizar datos de una propiedad |
| `DELETE` | `/properties/{property}` | `properties.destroy` | `PropertyController@destroy` | Eliminar una propiedad |

---

## 11. Equipos de Propiedad (Property Equipment)

**Middleware:** `auth` — Prefijo: `/properties/{property}/equipment`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/properties/{property}/equipment/create` | `properties.equipment.store` | `PropertyEquipmentController@store` | Agregar equipo a una propiedad |
| `PUT` | `/properties/{property}/equipment/{equipment}` | `properties.equipment.update` | `PropertyEquipmentController@update` | Actualizar equipo de una propiedad |
| `DELETE` | `/properties/{property}/equipment/{equipment}` | `properties.equipment.destroy` | `PropertyEquipmentController@destroy` | Eliminar equipo de una propiedad |

---

## 12. Trabajos (Jobs)

**Middleware:** `auth` — Prefijo: `/jobs`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/jobs` | `jobs.index` | `JobController@index` | Listar todos los trabajos |
| `GET` | `/jobs/export` | `jobs.export` | `JobController@export` | Exportar listado de trabajos |
| `GET` | `/jobs/create` | `jobs.create` | `JobController@create` | Mostrar formulario de creación de trabajo |
| `POST` | `/jobs/create` | `jobs.store` | `JobController@store` | Crear un nuevo trabajo |
| `GET` | `/jobs/{job}/edit` | `jobs.edit` | `JobController@edit` | Mostrar formulario de edición de trabajo |
| `GET` | `/jobs/{job}/emails` | `jobs.emails` | `JobController@emails` | Ver historial de emails enviados del trabajo |
| `GET` | `/jobs/{job}/sms` | `jobs.sms` | `JobController@sms` | Ver historial de SMS enviados del trabajo |
| `GET` | `/jobs/{job}/logs` | `jobs.logs` | `JobController@logs` | Ver logs/historial de actividad del trabajo |
| `PUT` | `/jobs/{job}` | `jobs.update` | `JobController@update` | Actualizar datos de un trabajo |
| `DELETE` | `/jobs/{job}` | `jobs.destroy` | `JobController@destroy` | Eliminar un trabajo |
| `PUT` | `/jobs/{job}/close` | `jobs.close` | `JobController@close` | Cerrar un trabajo |
| `PUT` | `/jobs/{job}/dispatch` | `jobs.dispatch` | `JobController@sendDispatchEmail` | Enviar email de despacho al técnico |
| `PUT` | `/jobs/{job}/dispatch-supervisor` | `jobs.dispatch-supervisor` | `JobController@sendDispatchSupervisorEmail` | Enviar email de despacho al supervisor |
| `PUT` | `/jobs/{job}/dispatch-sms` | `jobs.dispatch-sms` | `JobController@sendDispatchSMS` | Enviar SMS de despacho |
| `PUT` | `/jobs/{job}/log-attempted-call` | `jobs.log-attempted-call` | `JobController@logAttemptedCall` | Registrar intento de llamada al cliente |

---

## 13. Equipos de Trabajo (Job Equipment)

**Middleware:** `auth` — Prefijo: `/jobs/{job}/equipment`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/jobs/{job}/equipment/create` | `jobs.equipment.store` | `JobEquipmentController@store` | Agregar equipo a un trabajo |
| `PUT` | `/jobs/{job}/equipment/{equipment}` | `jobs.equipment.update` | `JobEquipmentController@update` | Actualizar equipo de un trabajo |
| `DELETE` | `/jobs/{job}/equipment/{equipment}` | `jobs.equipment.destroy` | `JobEquipmentController@destroy` | Eliminar equipo de un trabajo |

---

## 14. Residentes de Trabajo (Job Residents)

**Middleware:** `auth` — Prefijo: `/jobs/{job}/residents`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/jobs/{job}/residents/create` | `jobs.residents.store` | `JobResidentController@store` | Agregar residente a un trabajo |
| `PUT` | `/jobs/{job}/residents/{resident}` | `jobs.residents.update` | `JobResidentController@update` | Actualizar datos de un residente |
| `DELETE` | `/jobs/{job}/residents/{resident}` | `jobs.residents.destroy` | `JobResidentController@destroy` | Eliminar residente de un trabajo |

---

## 15. Visitas de Trabajo (Job Visits)

**Middleware:** `auth` — Prefijo: `/jobs/{job}/visits`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/jobs/{job}/visits/create` | `jobs.visits.store` | `JobVisitController@store` | Crear una visita/reporte para un trabajo |
| `PUT` | `/jobs/{job}/visits/{visit}` | `jobs.visits.update` | `JobVisitController@update` | Actualizar una visita |
| `DELETE` | `/jobs/{job}/visits/{visit}` | `jobs.visits.destroy` | `JobVisitController@destroy` | Eliminar una visita |
| `GET` | `/jobs/{job}/visits/{visit}/download` | `jobs.visits.download` | `JobVisitController@download` | Descargar reporte de visita |

---

## 16. Tareas de Trabajo (Job Tasks)

**Middleware:** `auth` — Prefijo: `/jobs/{job}/tasks`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/jobs/{job}/tasks/create` | `jobs.tasks.store` | `JobTaskController@store` | Crear una tarea para un trabajo |
| `PUT` | `/jobs/{job}/tasks/{task}` | `jobs.tasks.update` | `JobTaskController@update` | Actualizar una tarea |
| `DELETE` | `/jobs/{job}/tasks/{task}` | `jobs.tasks.destroy` | `JobTaskController@destroy` | Eliminar una tarea |
| `PUT` | `/jobs/{job}/tasks/{task}/send-notification` | `jobs.tasks.send-notification` | `JobTaskController@sendNotification` | Enviar notificación sobre una tarea |

---

## 17. Tarifas de Trabajo (Job Rates)

**Middleware:** `auth` — Prefijo: `/jobs/{job}/rates`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/jobs/{job}/rates/create` | `jobs.rates.store` | `JobRateController@store` | Crear una tarifa para un trabajo |
| `PUT` | `/jobs/{job}/rates/{rate}` | `jobs.rates.update` | `JobRateController@update` | Actualizar una tarifa |
| `DELETE` | `/jobs/{job}/rates/{rate}` | `jobs.rates.destroy` | `JobRateController@destroy` | Eliminar una tarifa |

---

## 18. Actividades/Notas de Trabajo (Job Activities)

**Middleware:** `auth` — Prefijo: `/jobs/{job}/activities`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/jobs/{job}/activities/create` | `jobs.activities.store` | `JobActivityController@store` | Crear una nota/actividad en un trabajo |
| `DELETE` | `/jobs/{job}/activities/{activity}` | `jobs.activities.destroy` | `JobActivityController@destroy` | Eliminar una nota/actividad |

---

## 19. Reclamaciones de Garantía (Warranty Claims)

**Middleware:** `auth` — Prefijo: `/warranty-claims`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/warranty-claims` | `warranty-claims.index` | `WarrantyClaimController@index` | Listar todas las reclamaciones de garantía |
| `GET` | `/warranty-claims/create` | `warranty-claims.create` | `WarrantyClaimController@create` | Mostrar formulario de creación |
| `POST` | `/warranty-claims/create` | `warranty-claims.store` | `WarrantyClaimController@store` | Crear una nueva reclamación |
| `GET` | `/warranty-claims/{claim}/edit` | `warranty-claims.edit` | `WarrantyClaimController@edit` | Mostrar formulario de edición |
| `PUT` | `/warranty-claims/{claim}` | `warranty-claims.update` | `WarrantyClaimController@update` | Actualizar una reclamación |
| `DELETE` | `/warranty-claims/{claim}` | `warranty-claims.destroy` | `WarrantyClaimController@destroy` | Eliminar una reclamación |

---

## 20. Garantías (Warranties)

**Middleware:** `auth` — Prefijo: `/warranties`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/warranties` | `warranties.index` | `WarrantyController@index` | Listar todas las garantías |
| `GET` | `/warranties/create` | `warranties.create` | `WarrantyController@create` | Mostrar formulario de creación |
| `POST` | `/warranties/create` | `warranties.store` | `WarrantyController@store` | Crear una nueva garantía |
| `GET` | `/warranties/{warranty}/edit` | `warranties.edit` | `WarrantyController@edit` | Mostrar formulario de edición |
| `PUT` | `/warranties/{warranty}` | `warranties.update` | `WarrantyController@update` | Actualizar una garantía |
| `DELETE` | `/warranties/{warranty}` | `warranties.destroy` | `WarrantyController@destroy` | Eliminar una garantía |

---

## 21. Equipos de Garantía (Warranty Equipment)

**Middleware:** `auth` — Prefijo: `/warranties/{warranty}/equipment`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/warranties/{warranty}/equipment/create` | `warranties.equipment.store` | `WarrantyEquipmentController@store` | Agregar equipo a una garantía |
| `PUT` | `/warranties/{warranty}/equipment/{equipment}` | `warranties.equipment.update` | `WarrantyEquipmentController@update` | Actualizar equipo de una garantía |
| `DELETE` | `/warranties/{warranty}/equipment/{equipment}` | `warranties.equipment.destroy` | `WarrantyEquipmentController@destroy` | Eliminar equipo de una garantía |

---

## 22. Nómina (Payroll)

**Middleware:** `auth` — Prefijo: `/payroll`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/payroll` | `payroll.payroll` | `PayrollController@payroll` | Vista principal de nómina |
| `GET` | `/payroll/{user}/pay` | `payroll.pay` | `PayrollController@pay` | Ver detalle de pago de un usuario |
| `POST` | `/payroll/{user}/pay/mark-paid` | `payroll.pay.mark-paid` | `PayrollController@markPaid` | Marcar pago como realizado |
| `POST` | `/payroll/{user}/pay/hold` | `payroll.pay.mark-held` | `PayrollController@markHeld` | Retener/pausar un pago |
| `GET` | `/payroll/{user}/paystub` | `payroll.paystub` | `PayrollController@viewPaystub` | Ver recibo de pago de un usuario |
| `POST` | `/payroll/{user}/paystub/email` | `payroll.paystub.email` | `PayrollController@emailPaystub` | Enviar recibo de pago por email |

---

## 23. Configuración General (Settings)

**Middleware:** `auth` — Prefijo: `/settings/general`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/general/edit` | `settings.general.edit` | `SettingsController@edit` | Mostrar formulario de configuración general |
| `PUT` | `/settings/general/{settings}` | `settings.general.update` | `SettingsController@update` | Actualizar configuración general |

---

## 24. Estados de Trabajo (Job Statuses)

**Middleware:** `auth` — Prefijo: `/settings/job-statuses`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/job-statuses` | `settings.job-statuses.index` | `JobStatusController@index` | Listar estados de trabajo |
| `GET` | `/settings/job-statuses/create` | `settings.job-statuses.create` | `JobStatusController@create` | Formulario de creación |
| `POST` | `/settings/job-statuses/create` | `settings.job-statuses.store` | `JobStatusController@store` | Crear estado de trabajo |
| `GET` | `/settings/job-statuses/{status}/edit` | `settings.job-statuses.edit` | `JobStatusController@edit` | Formulario de edición |
| `PUT` | `/settings/job-statuses/{status}` | `settings.job-statuses.update` | `JobStatusController@update` | Actualizar estado de trabajo |
| `DELETE` | `/settings/job-statuses/{status}` | `settings.job-statuses.destroy` | `JobStatusController@destroy` | Eliminar estado de trabajo |

---

## 25. Estados de Tarea (Task Statuses)

**Middleware:** `auth` — Prefijo: `/settings/task-statuses`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/task-statuses` | `settings.task-statuses.index` | `TaskStatusController@index` | Listar estados de tarea |
| `GET` | `/settings/task-statuses/create` | `settings.task-statuses.create` | `TaskStatusController@create` | Formulario de creación |
| `POST` | `/settings/task-statuses/create` | `settings.task-statuses.store` | `TaskStatusController@store` | Crear estado de tarea |
| `GET` | `/settings/task-statuses/{status}/edit` | `settings.task-statuses.edit` | `TaskStatusController@edit` | Formulario de edición |
| `PUT` | `/settings/task-statuses/{status}` | `settings.task-statuses.update` | `TaskStatusController@update` | Actualizar estado de tarea |
| `DELETE` | `/settings/task-statuses/{status}` | `settings.task-statuses.destroy` | `TaskStatusController@destroy` | Eliminar estado de tarea |

---

## 26. Estados de Trabajo de Técnico (Tech Job Statuses)

**Middleware:** `auth` — Prefijo: `/settings/tech-job-statuses`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/tech-job-statuses` | `settings.tech-job-statuses.index` | `TechJobStatusController@index` | Listar estados de trabajo de técnico |
| `GET` | `/settings/tech-job-statuses/create` | `settings.tech-job-statuses.create` | `TechJobStatusController@create` | Formulario de creación |
| `POST` | `/settings/tech-job-statuses/create` | `settings.tech-job-statuses.store` | `TechJobStatusController@store` | Crear estado de trabajo de técnico |
| `GET` | `/settings/tech-job-statuses/{status}/edit` | `settings.tech-job-statuses.edit` | `TechJobStatusController@edit` | Formulario de edición |
| `PUT` | `/settings/tech-job-statuses/{status}` | `settings.tech-job-statuses.update` | `TechJobStatusController@update` | Actualizar estado de trabajo de técnico |
| `DELETE` | `/settings/tech-job-statuses/{status}` | `settings.tech-job-statuses.destroy` | `TechJobStatusController@destroy` | Eliminar estado de trabajo de técnico |

---

## 27. Flujos de Trabajo (Workflows)

**Middleware:** `auth` — Prefijo: `/settings/workflows`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/workflows` | `settings.workflows.index` | `WorkflowController@index` | Listar flujos de trabajo |
| `GET` | `/settings/workflows/duplicate/{workflow}` | `settings.workflows.duplicate` | `WorkflowController@duplicate` | Duplicar un flujo de trabajo existente |
| `GET` | `/settings/workflows/create` | `settings.workflows.create` | `WorkflowController@create` | Formulario de creación |
| `POST` | `/settings/workflows/create` | `settings.workflows.store` | `WorkflowController@store` | Crear flujo de trabajo |
| `GET` | `/settings/workflows/{workflow}/edit` | `settings.workflows.edit` | `WorkflowController@edit` | Formulario de edición |
| `PUT` | `/settings/workflows/{workflow}` | `settings.workflows.update` | `WorkflowController@update` | Actualizar flujo de trabajo |
| `DELETE` | `/settings/workflows/{workflow}` | `settings.workflows.destroy` | `WorkflowController@destroy` | Eliminar flujo de trabajo |

---

## 28. Categorías de Trabajo (Job Categories)

**Middleware:** `auth` — Prefijo: `/settings/job-categories`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/job-categories` | `settings.job-categories.index` | `JobCategoryController@index` | Listar categorías de trabajo |
| `GET` | `/settings/job-categories/create` | `settings.job-categories.create` | `JobCategoryController@create` | Formulario de creación |
| `POST` | `/settings/job-categories/create` | `settings.job-categories.store` | `JobCategoryController@store` | Crear categoría |
| `GET` | `/settings/job-categories/{category}/edit` | `settings.job-categories.edit` | `JobCategoryController@edit` | Formulario de edición |
| `PUT` | `/settings/job-categories/{category}` | `settings.job-categories.update` | `JobCategoryController@update` | Actualizar categoría |
| `DELETE` | `/settings/job-categories/{category}` | `settings.job-categories.destroy` | `JobCategoryController@destroy` | Eliminar categoría |

---

## 29. Prioridades de Trabajo (Job Priorities)

**Middleware:** `auth` — Prefijo: `/settings/job-priorities`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/job-priorities` | `settings.job-priorities.index` | `JobPriorityController@index` | Listar prioridades de trabajo |
| `GET` | `/settings/job-priorities/create` | `settings.job-priorities.create` | `JobPriorityController@create` | Formulario de creación |
| `POST` | `/settings/job-priorities/create` | `settings.job-priorities.store` | `JobPriorityController@store` | Crear prioridad |
| `GET` | `/settings/job-priorities/{priority}/edit` | `settings.job-priorities.edit` | `JobPriorityController@edit` | Formulario de edición |
| `PUT` | `/settings/job-priorities/{priority}` | `settings.job-priorities.update` | `JobPriorityController@update` | Actualizar prioridad |
| `DELETE` | `/settings/job-priorities/{priority}` | `settings.job-priorities.destroy` | `JobPriorityController@destroy` | Eliminar prioridad |

---

## 30. Estados de Reclamación de Garantía (Warranty Claim Statuses)

**Middleware:** `auth` — Prefijo: `/settings/warranty-claim-statuses`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/warranty-claim-statuses` | `settings.warranty-claim-statuses.index` | `WarrantyClaimStatusController@index` | Listar estados de reclamación |
| `GET` | `/settings/warranty-claim-statuses/create` | `settings.warranty-claim-statuses.create` | `WarrantyClaimStatusController@create` | Formulario de creación |
| `POST` | `/settings/warranty-claim-statuses/create` | `settings.warranty-claim-statuses.store` | `WarrantyClaimStatusController@store` | Crear estado |
| `GET` | `/settings/warranty-claim-statuses/{status}/edit` | `settings.warranty-claim-statuses.edit` | `WarrantyClaimStatusController@edit` | Formulario de edición |
| `PUT` | `/settings/warranty-claim-statuses/{status}` | `settings.warranty-claim-statuses.update` | `WarrantyClaimStatusController@update` | Actualizar estado |
| `DELETE` | `/settings/warranty-claim-statuses/{status}` | `settings.warranty-claim-statuses.destroy` | `WarrantyClaimStatusController@destroy` | Eliminar estado |

---

## 31. Estados de Garantía (Warranty Statuses)

**Middleware:** `auth` — Prefijo: `/settings/warranty-statuses`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/warranty-statuses` | `settings.warranty-statuses.index` | `WarrantyStatusController@index` | Listar estados de garantía |
| `GET` | `/settings/warranty-statuses/create` | `settings.warranty-statuses.create` | `WarrantyStatusController@create` | Formulario de creación |
| `POST` | `/settings/warranty-statuses/create` | `settings.warranty-statuses.store` | `WarrantyStatusController@store` | Crear estado |
| `GET` | `/settings/warranty-statuses/{status}/edit` | `settings.warranty-statuses.edit` | `WarrantyStatusController@edit` | Formulario de edición |
| `PUT` | `/settings/warranty-statuses/{status}` | `settings.warranty-statuses.update` | `WarrantyStatusController@update` | Actualizar estado |
| `DELETE` | `/settings/warranty-statuses/{status}` | `settings.warranty-statuses.destroy` | `WarrantyStatusController@destroy` | Eliminar estado |

---

## 32. Tipos de Garantía (Warranty Types)

**Middleware:** `auth` — Prefijo: `/settings/warranty-types`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/warranty-types` | `settings.warranty-types.index` | `WarrantyTypeController@index` | Listar tipos de garantía |
| `GET` | `/settings/warranty-types/create` | `settings.warranty-types.create` | `WarrantyTypeController@create` | Formulario de creación |
| `POST` | `/settings/warranty-types/create` | `settings.warranty-types.store` | `WarrantyTypeController@store` | Crear tipo |
| `GET` | `/settings/warranty-types/{type}/edit` | `settings.warranty-types.edit` | `WarrantyTypeController@edit` | Formulario de edición |
| `PUT` | `/settings/warranty-types/{type}` | `settings.warranty-types.update` | `WarrantyTypeController@update` | Actualizar tipo |
| `DELETE` | `/settings/warranty-types/{type}` | `settings.warranty-types.destroy` | `WarrantyTypeController@destroy` | Eliminar tipo |

---

## 33. Plantillas de Email (Email Templates)

**Middleware:** `auth` — Prefijo: `/settings/email-templates`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/email-templates` | `settings.email-templates.index` | `EmailTemplateController@index` | Listar plantillas de email |
| `GET` | `/settings/email-templates/create` | `settings.email-templates.create` | `EmailTemplateController@create` | Formulario de creación |
| `POST` | `/settings/email-templates/create` | `settings.email-templates.store` | `EmailTemplateController@store` | Crear plantilla |
| `GET` | `/settings/email-templates/{template}/edit` | `settings.email-templates.edit` | `EmailTemplateController@edit` | Formulario de edición |
| `PUT` | `/settings/email-templates/{template}` | `settings.email-templates.update` | `EmailTemplateController@update` | Actualizar plantilla |
| `DELETE` | `/settings/email-templates/{template}` | `settings.email-templates.destroy` | `EmailTemplateController@destroy` | Eliminar plantilla |

---

## 34. Plantillas de Texto (Text Templates)

**Middleware:** `auth` — Prefijo: `/settings/text-templates`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/settings/text-templates` | `settings.text-templates.index` | `TextTemplateController@index` | Listar plantillas de texto/SMS |
| `GET` | `/settings/text-templates/create` | `settings.text-templates.create` | `TextTemplateController@create` | Formulario de creación |
| `POST` | `/settings/text-templates/create` | `settings.text-templates.store` | `TextTemplateController@store` | Crear plantilla |
| `GET` | `/settings/text-templates/{template}/edit` | `settings.text-templates.edit` | `TextTemplateController@edit` | Formulario de edición |
| `PUT` | `/settings/text-templates/{template}` | `settings.text-templates.update` | `TextTemplateController@update` | Actualizar plantilla |
| `DELETE` | `/settings/text-templates/{template}` | `settings.text-templates.destroy` | `TextTemplateController@destroy` | Eliminar plantilla |

---

## 35. Usuarios

**Middleware:** `auth` — Prefijo: `/users`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/users` | `users.index` | `UserController@index` | Listar todos los usuarios |
| `GET` | `/users/create` | `users.create` | `UserController@create` | Formulario de creación de usuario |
| `POST` | `/users/create` | `users.store` | `UserController@store` | Crear un nuevo usuario |
| `GET` | `/users/{user}/jobs` | `users.jobs` | `UserController@jobs` | Ver trabajos asignados a un usuario |
| `GET` | `/users/{user}/reports` | `users.reports` | `UserController@reports` | Ver reportes de un usuario |
| `GET` | `/users/{user}/rates` | `users.rates` | `UserController@rates` | Ver tarifas de un usuario |
| `GET` | `/users/{user}/edit` | `users.edit` | `UserController@edit` | Formulario de edición de usuario |
| `PUT` | `/users/{user}` | `users.update` | `UserController@update` | Actualizar datos de un usuario |
| `DELETE` | `/users/{user}` | `users.destroy` | `UserController@destroy` | Eliminar un usuario |

---

## 36. Roles

**Middleware:** `auth` — Prefijo: `/roles`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `GET` | `/roles` | `roles.index` | `RoleController@index` | Listar todos los roles |
| `GET` | `/roles/create` | `roles.create` | `RoleController@create` | Formulario de creación de rol |
| `POST` | `/roles/create` | `roles.store` | `RoleController@store` | Crear un nuevo rol |
| `GET` | `/roles/{role}/edit` | `roles.edit` | `RoleController@edit` | Formulario de edición de rol |
| `PUT` | `/roles/{role}` | `roles.update` | `RoleController@update` | Actualizar un rol |
| `DELETE` | `/roles/{role}` | `roles.destroy` | `RoleController@destroy` | Eliminar un rol |

---

## 37. Archivos (Files)

**Middleware:** `auth` — Prefijo: `/files`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/files/create` | `files.store` | `FileController@store` | Subir un archivo |
| `PUT` | `/files/{file}` | `files.update` | `FileController@update` | Actualizar un archivo |
| `DELETE` | `/files/{file}` | `files.destroy` | `FileController@destroy` | Eliminar un archivo |

---

## 38. Cálculo de Pagos

**Middleware:** `auth`

| Método | URI | Nombre | Controlador | Descripción |
|---|---|---|---|---|
| `POST` | `/calculate-rate-payment` | `calculate-payment` | `JobRateController@calculatePayment` | Calcular el pago de una tarifa de trabajo |

---

## 39. Endpoints Comentados (Deshabilitados)

Los siguientes endpoints están **comentados** en el código y **no están activos**:

### Facturas (Invoices) — Deshabilitado

| Método | URI | Controlador | Descripción |
|---|---|---|---|
| `GET` | `/invoices` | `InvoiceController@index` | Listar facturas |
| `GET` | `/invoices/create` | `InvoiceController@create` | Formulario de creación |
| `POST` | `/invoices/create` | `InvoiceController@store` | Crear factura |
| `GET` | `/invoices/{invoice}/edit` | `InvoiceController@edit` | Formulario de edición |
| `POST` | `/invoices/{invoice}/email` | `InvoiceController@email` | Enviar factura por email |
| `PUT` | `/invoices/{invoice}` | `InvoiceController@update` | Actualizar factura |
| `DELETE` | `/invoices/{invoice}` | `InvoiceController@destroy` | Eliminar factura |
| `GET` | `/invoices/{invoice}/payments/create` | `InvoicePaymentController@create` | Formulario de creación de pago |
| `POST` | `/invoices/{invoice}/payments/create` | `InvoicePaymentController@store` | Registrar pago |
| `GET` | `/invoices/{invoice}/payments/{payment}/edit` | `InvoicePaymentController@edit` | Formulario de edición de pago |
| `PUT` | `/invoices/{invoice}/payments/{payment}` | `InvoicePaymentController@update` | Actualizar pago |
| `DELETE` | `/invoices/{invoice}/payments/{payment}` | `InvoicePaymentController@destroy` | Eliminar pago |

### Cotizaciones (Quotes) — Deshabilitado

| Método | URI | Controlador | Descripción |
|---|---|---|---|
| `GET` | `/quotes` | `QuoteController@index` | Listar cotizaciones |
| `GET` | `/quotes/create` | `QuoteController@create` | Formulario de creación |
| `POST` | `/quotes/create` | `QuoteController@store` | Crear cotización |
| `GET` | `/quotes/{quote}/edit` | `QuoteController@edit` | Formulario de edición |
| `POST` | `/quotes/{quote}/email` | `QuoteController@email` | Enviar cotización por email |
| `PUT` | `/quotes/{quote}` | `QuoteController@update` | Actualizar cotización |
| `DELETE` | `/quotes/{quote}` | `QuoteController@destroy` | Eliminar cotización |

---

## Notas Técnicas

- **Framework:** Laravel (PHP)
- **Archivo de rutas principal:** `routes/web.php`
- **Archivo API:** `routes/api.php` (vacío, sin rutas definidas)
- **Autenticación:** Middleware `auth` de Laravel para todas las rutas protegidas
- **Patrón de rutas:** La aplicación sigue el patrón CRUD estándar de Laravel (index, create, store, edit, update, destroy)
- **Pagos:** Integración con Stripe para pagos online de facturas
- **Comunicaciones:** Soporte para envío de emails y SMS (despacho de trabajos, notificaciones de tareas, recibos de pago)
