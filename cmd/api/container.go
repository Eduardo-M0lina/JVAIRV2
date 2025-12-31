package main

import (
	"net/http"

	"github.com/your-org/jvairv2/configs"
	"github.com/your-org/jvairv2/pkg/repository/mysql"
	"github.com/your-org/jvairv2/pkg/rest/handler"
	"github.com/your-org/jvairv2/pkg/rest/router"
)

// Container contiene todas las dependencias de la aplicación
type Container struct {
	Config        *configs.Config
	DBConnection  *mysql.Connection
	HealthHandler *handler.HealthHandler
	Router        http.Handler
}

// NewContainer crea un nuevo contenedor con todas las dependencias inicializadas
func NewContainer(configPath string) (*Container, error) {
	// Cargar configuración
	config, err := configs.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Inicializar conexión a la base de datos
	dbConn, err := mysql.NewConnection(&config.DB)
	if err != nil {
		return nil, err
	}

	// Inicializar handlers
	healthHandler := handler.NewHealthHandler(dbConn)

	// Inicializar router
	r := router.New(healthHandler)

	return &Container{
		Config:        config,
		DBConnection:  dbConn,
		HealthHandler: healthHandler,
		Router:        r,
	}, nil
}

// Close cierra todas las conexiones
func (c *Container) Close() error {
	if c.DBConnection != nil {
		return c.DBConnection.Close()
	}
	return nil
}
