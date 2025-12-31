# JVAIRV2 - Migración de JVAIR a Go

Este proyecto es una migración del backend de JVAIR (PHP/Laravel) a Go, siguiendo los principios de Clean Architecture.

## Estructura del Proyecto

El proyecto sigue una arquitectura limpia con las siguientes capas:

```
JVAIRV2/
├── cmd/api/                      # Punto de entrada de la aplicación
│   ├── main.go                    # Servidor HTTP y graceful shutdown
│   └── container.go               # Inicialización y orquestación de dependencias
├── configs/                       # Configuración de la aplicación
│   ├── app.env                    # Variables de entorno
│   └── config.go                  # Cargador de configuración
├── pkg/                           # Núcleo del proyecto (código reutilizable)
│   ├── domain/                    # Capa de dominio (entidades e interfaces)
│   │   └── ...                    # Subdirectorios por contexto de negocio
│   ├── usecase/                   # Casos de uso (lógica de aplicación)
│   │   └── ...                    # Implementaciones de casos de uso
│   ├── repository/                # Implementaciones de repositorios
│   │   └── mysql/                 # Implementación MySQL
│   │       └── connection.go      # Conexión a la base de datos
│   └── rest/                      # API REST
│       ├── handler/               # Handlers HTTP
│       │   └── health_handler.go  # Endpoint de health check
│       └── router/                # Definición de rutas
│           └── router.go          # Configuración del router
└── README.md                      # Documentación del proyecto
```

### Responsabilidades por Capa

#### `cmd/api`
- **main.go**: Punto de entrada de la aplicación, configuración del servidor HTTP y manejo de graceful shutdown
- **container.go**: Inicialización e inyección de dependencias (configuración, base de datos, handlers, router)

#### `configs`
- Configuración de la aplicación
- Carga de variables de entorno

#### `pkg/domain`
- Entidades de dominio
- Interfaces de repositorio (puertos)
- Reglas de negocio puras
- No depende de ninguna otra capa

#### `pkg/usecase`
- Implementación de casos de uso
- Orquestación de reglas de negocio
- Depende solo de `domain`

#### `pkg/repository`
- Implementaciones concretas de repositorios
- Conexión a la base de datos
- Implementa interfaces definidas en `domain`

#### `pkg/rest`
- API REST
- Handlers HTTP
- Definición de rutas
- Transformación de DTOs

## Requisitos

- Go 1.21 o superior
- MySQL 5.7 o superior

## Configuración

1. Asegúrate de tener Go instalado:

```bash
go version
```

2. Configura la base de datos en `configs/app.env`:

```
DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=tu_password
DB_NAME=jvair
```

## Ejecución

1. Instala las dependencias:

```bash
cd /Users/eduardo/projects/jvair/JVAIRV2
go mod tidy
```

2. Ejecuta la aplicación:

```bash
go run cmd/api/main.go
```

3. Verifica que la aplicación esté funcionando correctamente:

```bash
curl http://localhost:8080/health
```

Deberías recibir una respuesta JSON con el estado de la aplicación.

## Desarrollo

Para agregar nuevos módulos, sigue estos pasos:

1. Define las entidades en `pkg/domain/{module}/entity.go`
2. Define las interfaces del repositorio en `pkg/domain/{module}/repository.go`
3. Implementa el repositorio en `pkg/repository/mysql/{module}_repository.go`
4. Implementa los casos de uso en `pkg/usecase/{module}/service.go`
5. Implementa los handlers HTTP en `pkg/rest/handler/{module}_handler.go`
6. Registra las rutas en `internal/bootstrap/router.go`

## Próximos Pasos

- Implementar autenticación y autorización
- Migrar módulo de Customers
- Migrar módulo de Properties
- Migrar módulo de Jobs
- Implementar tests unitarios y de integración
