// @title           JVAIR API
// @version         1.0
// @description     API para el sistema JVAIR
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.jvair.com/support
// @contact.email  support@jvair.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8090

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ingrese el token con el prefijo 'Bearer '

// @tag.name Auth
// @tag.description Operaciones de autenticación

// @tag.name Users
// @tag.description Operaciones de usuarios

// @tag.name Roles
// @tag.description Operaciones de roles

// @tag.name Abilities
// @tag.description Operaciones de abilities

// @tag.name AssignedRoles
// @tag.description Operaciones de asignación de roles

// @tag.name Permissions
// @tag.description Operaciones de permisos

// @tag.name Settings
// @tag.description Configuraciones del sistema

// @tag.name Workflows
// @tag.description Operaciones de workflows

// @tag.name Customers
// @tag.description Operaciones de clientes

// @tag.name Properties
// @tag.description Operaciones de propiedades

// @tag.name PropertyEquipment
// @tag.description Equipos HVAC de propiedades

// @tag.name Supervisors
// @tag.description Operaciones de supervisores

// @tag.name Jobs
// @tag.description Operaciones de trabajos

// @tag.name JobEquipment
// @tag.description Equipos HVAC de trabajos

// @tag.name JobCategories
// @tag.description Categorías de trabajos

// @tag.name JobStatuses
// @tag.description Estados de trabajos

// @tag.name JobPriorities
// @tag.description Prioridades de trabajos

// @tag.name TechnicianJobStatuses
// @tag.description Estados de técnico de trabajo

// @tag.name TaskStatuses
// @tag.description Estados de tareas

// @tag.name Quotes
// @tag.description Operaciones de cotizaciones

// @tag.name QuoteStatuses
// @tag.description Estados de cotizaciones

// @tag.name Invoices
// @tag.description Operaciones de facturas

// @tag.name Invoice Payments
// @tag.description Pagos de facturas

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/your-org/jvairv2/configs"
	"github.com/your-org/jvairv2/pkg/common/logger"
)

func main() {
	// Cargar configuración primero para obtener el ambiente
	configPath := filepath.Join(".", "configs")
	cfg, err := configs.LoadConfig(configPath)
	if err != nil {
		// Si falla la carga de config, usar desarrollo por defecto
		logger.Init("development")
		slog.Error("Error al cargar configuración, usando valores por defecto", "error", err)
	} else {
		// Inicializar logger con el ambiente configurado
		logger.Init(cfg.App.Environment)
		slog.Info("Iniciando aplicación JVAIR V2",
			"environment", cfg.App.Environment,
			"version", "2.0",
		)
	}

	// Inicializar contenedor de dependencias
	container, err := NewContainer(configPath)
	if err != nil {
		slog.Error("Error al inicializar el contenedor", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := container.Close(); err != nil {
			slog.Error("Error al cerrar el contenedor", "error", err)
		}
	}()

	slog.Info("Conexión a la base de datos establecida correctamente")

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", container.Config.Server.Port),
		Handler:      container.Router,
		ReadTimeout:  container.Config.Server.ReadTimeout,
		WriteTimeout: container.Config.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Canal para recibir errores del servidor
	serverErrors := make(chan error, 1)

	// Iniciar el servidor en una goroutine
	go func() {
		slog.Info("Servidor HTTP iniciado", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Canal para recibir señales del sistema operativo
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Bloquear hasta recibir una señal o un error
	select {
	case err := <-serverErrors:
		slog.Error("Error al iniciar el servidor", "error", err)
		os.Exit(1)

	case <-shutdown:
		slog.Info("Señal de apagado recibida, cerrando servidor...")

		// Crear un contexto con timeout para el apagado
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Apagar el servidor
		err := server.Shutdown(ctx)
		if err != nil {
			slog.Error("Error al apagar el servidor", "error", err)
			err = server.Close()
			if err != nil {
				slog.Error("Error al cerrar el servidor", "error", err)
			}
		}

		slog.Info("Servidor apagado correctamente")
	}
}
