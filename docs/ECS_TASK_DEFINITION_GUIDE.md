# ECS Task Definition y Variables de Entorno para JVAIRV2

Este documento resume cómo debe manejarse la configuración de entorno de **JVAIRV2** cuando se despliegue en **AWS ECS**, tomando como referencia la configuración existente de `WeCoolAtlanta-Prod`.

---

## Idea principal

JVAIRV2 **no debe leer directamente** el archivo JSON del Task Definition de ECS.

Ese JSON existe únicamente como configuración de AWS para definir cómo se ejecuta el contenedor. Cuando ECS levanta una tarea, AWS toma las variables declaradas en el Task Definition y las inyecta como **variables de entorno del sistema operativo dentro del contenedor**.

La aplicación Go solo necesita leer variables de entorno normales.

Flujo real:

```text
Task Definition JSON en ECS
  -> containerDefinitions[].environment[]
  -> ECS inyecta variables al contenedor
  -> viper.AutomaticEnv() las lee en tiempo de ejecución
```

---

## Cómo funciona hoy JVAIRV2

JVAIRV2 quedó preparado para trabajar en ambos escenarios:

### Desarrollo local

Lee variables desde:

- `configs/app.env`

### Producción en ECS

Si el archivo `app.env` no existe dentro del contenedor:

- la aplicación **no falla**
- Viper usa `AutomaticEnv()` para leer las variables inyectadas por ECS

Esto permite usar la misma aplicación tanto en local como en producción sin necesidad de cambiar código.

---

## ¿Aplica el mismo JSON de ECS que usa Laravel?

**Sí, parcialmente.**

El mismo enfoque aplica para JVAIRV2, y varias variables pueden reutilizarse exactamente igual:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_DEFAULT_REGION`
- `AWS_BUCKET`
- `AWS_URL`
- `DB_HOST`
- `DB_PORT`
- `DB_PASSWORD`
- `DB_DATABASE`

Pero JVAIRV2 necesita además variables propias de la API en Go, por ejemplo:

- `SERVER_PORT`
- `SERVER_TIMEOUT`
- `SERVER_READ_TIMEOUT`
- `SERVER_WRITE_TIMEOUT`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`
- `JWT_ACCESS_EXPIRATION`
- `JWT_REFRESH_EXPIRATION`
- `DB_DRIVER`
- `DB_USER`
- `DB_MAX_OPEN_CONNS`
- `DB_MAX_IDLE_CONNS`
- `DB_CONN_MAX_LIFETIME`
- `APP_ENV`

Por eso, para JVAIRV2 normalmente se crea **un nuevo Task Definition**, aunque reutilice muchos valores del Task Definition actual de Laravel.

---

## Importante: las variables del contenedor no se comparten automáticamente

Aunque dos servicios estén dentro del mismo entorno de AWS, cada servicio o tarea en ECS tiene su propia configuración.

Eso significa:

- El Task Definition `WeCoolAtlanta-Prod:74` aplica al contenedor definido allí
- JVAIRV2 debería tener su propio Task Definition
- Ambos pueden usar los mismos valores para base de datos y S3
- Pero no comparten automáticamente las variables por estar en el mismo entorno

---

## Ruta en AWS Console para crear el Task Definition

Ir a:

```text
https://us-east-1.console.aws.amazon.com/ecs/v2/task-definitions?region=us-east-1
```

Luego:

1. Click en `Create new task definition`
2. Elegir `AWS Fargate`
3. Definir el `Task definition family`
4. Configurar el contenedor
5. Agregar las variables de entorno
6. Guardar la revisión

---

## Sugerencia de nombre para JVAIRV2

Ejemplo:

```text
WeCoolAtlanta-API-Prod
```

Y para el contenedor:

```text
wecoolatlanta-api
```

---

## Variables sugeridas para JVAIRV2

### Variables reutilizables desde Laravel / entorno actual

| Variable | Valor esperado |
|---|---|
| `APP_ENV` | `production` |
| `DB_HOST` | host de RDS de producción |
| `DB_PORT` | `3306` |
| `DB_PASSWORD` | password de producción |
| `DB_DATABASE` | `wecoolatlanta` |
| `AWS_ACCESS_KEY_ID` | access key actual |
| `AWS_SECRET_ACCESS_KEY` | secret key actual |
| `AWS_DEFAULT_REGION` | `us-east-1` |
| `AWS_BUCKET` | `wecoolatlanta-prod` |
| `AWS_URL` | `https://images.wecoolatlanta.com` |

### Variables propias de JVAIRV2

| Variable | Valor sugerido |
|---|---|
| `SERVER_PORT` | `3001` |
| `SERVER_TIMEOUT` | `30s` |
| `SERVER_READ_TIMEOUT` | `15s` |
| `SERVER_WRITE_TIMEOUT` | `15s` |
| `DB_DRIVER` | `mysql` |
| `DB_USER` | usuario de BD para la API |
| `DB_MAX_OPEN_CONNS` | `25` |
| `DB_MAX_IDLE_CONNS` | `25` |
| `DB_CONN_MAX_LIFETIME` | `5m` |
| `JWT_ACCESS_SECRET` | secret seguro de producción |
| `JWT_REFRESH_SECRET` | secret seguro de producción |
| `JWT_ACCESS_EXPIRATION` | `15m` |
| `JWT_REFRESH_EXPIRATION` | `24h` |

---

## Recomendación de buenas prácticas

### Opción actual compatible

Usar variables directamente en el Task Definition.

Ventaja:

- simple
- consistente con el sistema actual

Desventaja:

- las credenciales sensibles quedan visibles en la definición de la tarea para quien tenga acceso

### Mejor práctica recomendada a futuro

Mover secretos sensibles a:

- **AWS Secrets Manager**, o
- **AWS Systems Manager Parameter Store**

Y luego referenciarlos desde ECS como secretos, en lugar de dejarlos como texto plano en `environment`.

Especialmente recomendable para:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `DB_PASSWORD`
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`

---

## Sobre S3 en JVAIRV2

JVAIRV2 ya fue ajustado para:

- usar `AWS_URL` cuando exista
- generar URLs públicas del tipo:

```text
https://images.wecoolatlanta.com/uploads/archivo.ext
```

Si `AWS_URL` no está configurado, hará fallback a la URL directa de S3.

Además, si no se definen `AWS_ACCESS_KEY_ID` y `AWS_SECRET_ACCESS_KEY`, el SDK de AWS puede usar el **IAM Role** de la tarea ECS, lo cual es una práctica mejor que usar claves fijas.

---

## Qué se necesitará cuando llegue el momento del deploy

Para desplegar JVAIRV2 en ECS hará falta:

1. Un **Dockerfile** para la aplicación Go
2. Una imagen publicada en **Amazon ECR**
3. Un **Task Definition** nuevo para JVAIRV2
4. Un **Service** de ECS apuntando a esa task
5. Variables de entorno configuradas en ECS
6. Idealmente, secretos movidos a **Secrets Manager**

---

## Resumen corto

- JVAIRV2 **no debe leer el JSON de ECS**
- ECS inyecta las variables del Task Definition al contenedor
- La app Go las lee con `viper.AutomaticEnv()`
- Sí se puede usar el mismo enfoque que Laravel
- Lo correcto es crear un **nuevo Task Definition** para JVAIRV2
- Se pueden reutilizar varias variables actuales de DB y S3
- A futuro, lo ideal es mover credenciales sensibles a **Secrets Manager**
