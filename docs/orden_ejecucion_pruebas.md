# Orden de Ejecución para Pruebas en Postman

Este documento describe el orden recomendado para ejecutar las pruebas en Postman con la API de JVAIR, respetando las dependencias entre los diferentes módulos.

## Secuencia de Pruebas

### 1. Autenticación (Auth)
- **Prioridad: Alta** - Requisito para todos los demás endpoints
- Endpoints a probar:
  1. `POST /auth/login` - Obtener token de acceso
  2. `POST /auth/refresh` - Refrescar token cuando sea necesario
  3. `POST /auth/logout` - Cerrar sesión al finalizar

### 2. Roles
- **Prioridad: Alta** - Necesarios para usuarios y permisos
- Endpoints a probar:
  1. `GET /roles` - Listar roles existentes
  2. `POST /roles` - Crear nuevo rol
  3. `GET /roles/{id}` - Verificar rol creado
  4. `PUT /roles/{id}` - Modificar rol (opcional)

### 3. Abilities
- **Prioridad: Alta** - Necesarias para configurar permisos
- Endpoints a probar:
  1. `GET /abilities` - Listar abilities existentes
  2. `POST /abilities` - Crear nueva ability
  3. `GET /abilities/{id}` - Verificar ability creada

### 4. Usuarios (Users)
- **Prioridad: Media** - Dependen de roles
- Endpoints a probar:
  1. `GET /users` - Listar usuarios existentes
  2. `POST /users` - Crear nuevo usuario con rol asignado
  3. `GET /users/{id}` - Verificar usuario creado
  4. `PUT /users/{id}` - Modificar usuario (opcional)
  5. `GET /users/{id}/roles` - Verificar roles de usuario

### 5. Asignación de Roles (AssignedRoles)
- **Prioridad: Media** - Dependen de usuarios y roles
- Endpoints a probar:
  1. `POST /assigned-roles` - Asignar rol a entidad
  2. `GET /assigned-roles/entity/{entity_type}/{entity_id}` - Verificar asignaciones
  3. `GET /assigned-roles/check/{role_id}/{entity_type}/{entity_id}` - Verificar asignación específica

### 6. Permisos (Permissions)
- **Prioridad: Baja** - Dependen de abilities y entidades
- Endpoints a probar:
  1. `POST /permissions` - Crear permiso
  2. `GET /permissions/check/{ability_id}/{entity_type}/{entity_id}` - Verificar permiso
  3. `GET /permissions/entity/{entity_type}/{entity_id}` - Listar permisos de entidad

### 7. Facturas (Invoices)
- **Prioridad: Media** - Dependen de jobs
- Endpoints a probar:
  1. `GET /invoices` - Listar facturas
  2. `POST /invoices` - Crear factura para un job existente
  3. `GET /invoices/{id}` - Verificar factura creada (incluye balance)
  4. `PUT /invoices/{id}` - Modificar factura
  5. `GET /invoices?status=unpaid` - Filtrar facturas pendientes
  6. `GET /invoices?status=paid` - Filtrar facturas pagadas
  7. `GET /invoices?jobId=1` - Filtrar facturas por trabajo
  8. `GET /invoices?search=INV` - Buscar facturas
  9. `DELETE /invoices/{id}` - Eliminar factura (soft delete)

### 8. Pagos de Facturas (Invoice Payments)
- **Prioridad: Media** - Dependen de invoices
- Endpoints a probar:
  1. `GET /invoices/{invoiceId}/payments` - Listar pagos de factura
  2. `POST /invoices/{invoiceId}/payments` - Crear pago para factura
  3. `GET /invoices/{invoiceId}/payments/{id}` - Verificar pago creado
  4. `PUT /invoices/{invoiceId}/payments/{id}` - Modificar pago
  5. `GET /invoices/{id}` - Verificar balance actualizado tras pago
  6. `DELETE /invoices/{invoiceId}/payments/{id}` - Eliminar pago (soft delete)

## Flujo de Prueba Completo

1. **Preparación Inicial**
   - Iniciar sesión y obtener token
   - Verificar roles y abilities existentes

2. **Configuración Básica**
   - Crear roles necesarios
   - Crear abilities necesarias

3. **Gestión de Usuarios**
   - Crear usuarios de prueba
   - Asignar roles a usuarios

4. **Configuración de Permisos**
   - Asignar permisos específicos
   - Verificar permisos asignados

5. **Pruebas de Integración**
   - Probar flujos completos que involucren múltiples endpoints
   - Verificar restricciones de acceso según roles y permisos

## Notas Importantes

- Siempre mantén activo un token válido. Si recibes errores 401, refresca el token.
- Algunos endpoints pueden requerir permisos específicos. Asegúrate de tener los roles adecuados.
- Para pruebas exhaustivas, crea múltiples usuarios con diferentes roles y permisos.
- Guarda IDs importantes (usuarios, roles, abilities) en variables de entorno para facilitar las pruebas.
