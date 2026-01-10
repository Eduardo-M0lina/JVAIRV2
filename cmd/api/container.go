package main

import (
	"net/http"

	"github.com/your-org/jvairv2/configs"
	commonAuth "github.com/your-org/jvairv2/pkg/common/auth"
	domainAuth "github.com/your-org/jvairv2/pkg/domain/auth"
	"github.com/your-org/jvairv2/pkg/repository/mysql"
	mysqlUser "github.com/your-org/jvairv2/pkg/repository/mysql/user"
	"github.com/your-org/jvairv2/pkg/rest/handler"
	abilityHandler "github.com/your-org/jvairv2/pkg/rest/handler/ability"
	assignedRoleHandler "github.com/your-org/jvairv2/pkg/rest/handler/assigned_role"
	authHandler "github.com/your-org/jvairv2/pkg/rest/handler/auth"
	permissionHandler "github.com/your-org/jvairv2/pkg/rest/handler/permission"
	roleHandler "github.com/your-org/jvairv2/pkg/rest/handler/role"
	userHandler "github.com/your-org/jvairv2/pkg/rest/handler/user"
	"github.com/your-org/jvairv2/pkg/rest/middleware"
	"github.com/your-org/jvairv2/pkg/rest/router"
)

// Container contiene todas las dependencias de la aplicación
type Container struct {
	Config              *configs.Config
	DBConnection        *mysql.Connection
	HealthHandler       *handler.HealthHandler
	AuthHandler         *authHandler.Handler
	UserHandler         *userHandler.Handler
	RoleHandler         *roleHandler.Handler
	AbilityHandler      *abilityHandler.Handler
	AssignedRoleHandler *assignedRoleHandler.Handler
	PermissionHandler   *permissionHandler.Handler
	AuthMiddleware      *middleware.AuthMiddleware
	Router              http.Handler
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

	// Inicializar repositorios
	userRepo := mysqlUser.NewRepository(dbConn.GetDB())

	// Inicializar servicios
	tokenStore := commonAuth.NewMemoryTokenStore()
	authService := commonAuth.NewJWTService(
		config.JWT.AccessSecret,
		config.JWT.RefreshSecret,
		config.JWT.AccessExpiration,
		config.JWT.RefreshExpiration,
		tokenStore,
	)

	// Inicializar casos de uso
	authUC := domainAuth.NewUseCase(userRepo, authService)

	// Inicializar handlers
	healthHandler := handler.NewHealthHandler(dbConn)
	authHandler := authHandler.NewHandler(authUC)

	// TODO: Implementar los casos de uso para estos handlers
	// Por ahora, usamos nil para permitir la compilación
	userHandler := &userHandler.Handler{}
	roleHandler := &roleHandler.Handler{}
	abilityHandler := &abilityHandler.Handler{}
	assignedRoleHandler := &assignedRoleHandler.Handler{}
	permissionHandler := &permissionHandler.Handler{}

	// Inicializar middlewares
	authMiddleware := middleware.NewAuthMiddleware(authUC)

	// Inicializar router
	r := router.New(
		healthHandler,
		authHandler,
		userHandler,
		roleHandler,
		abilityHandler,
		assignedRoleHandler,
		permissionHandler,
		authMiddleware,
	)

	return &Container{
		Config:              config,
		DBConnection:        dbConn,
		HealthHandler:       healthHandler,
		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		RoleHandler:         roleHandler,
		AbilityHandler:      abilityHandler,
		AssignedRoleHandler: assignedRoleHandler,
		PermissionHandler:   permissionHandler,
		AuthMiddleware:      authMiddleware,
		Router:              r,
	}, nil
}

// Close cierra todas las conexiones
func (c *Container) Close() error {
	if c.DBConnection != nil {
		return c.DBConnection.Close()
	}
	return nil
}
