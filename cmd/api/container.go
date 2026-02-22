package main

import (
	http "net/http"

	configs "github.com/your-org/jvairv2/configs"
	commonAuth "github.com/your-org/jvairv2/pkg/common/auth"
	ability "github.com/your-org/jvairv2/pkg/domain/ability"
	assignedRole "github.com/your-org/jvairv2/pkg/domain/assigned_role"
	domainAuth "github.com/your-org/jvairv2/pkg/domain/auth"
	customer "github.com/your-org/jvairv2/pkg/domain/customer"
	jobCategory "github.com/your-org/jvairv2/pkg/domain/job_category"
	jobPriority "github.com/your-org/jvairv2/pkg/domain/job_priority"
	jobStatus "github.com/your-org/jvairv2/pkg/domain/job_status"
	permission "github.com/your-org/jvairv2/pkg/domain/permission"
	property "github.com/your-org/jvairv2/pkg/domain/property"
	role "github.com/your-org/jvairv2/pkg/domain/role"
	settings "github.com/your-org/jvairv2/pkg/domain/settings"
	taskStatus "github.com/your-org/jvairv2/pkg/domain/task_status"
	techJobStatus "github.com/your-org/jvairv2/pkg/domain/technician_job_status"
	user "github.com/your-org/jvairv2/pkg/domain/user"
	workflow "github.com/your-org/jvairv2/pkg/domain/workflow"
	mysql "github.com/your-org/jvairv2/pkg/repository/mysql"
	mysqlAbility "github.com/your-org/jvairv2/pkg/repository/mysql/ability"
	mysqlAssignedRole "github.com/your-org/jvairv2/pkg/repository/mysql/assigned_role"
	mysqlCustomer "github.com/your-org/jvairv2/pkg/repository/mysql/customer"
	mysqlJobCategory "github.com/your-org/jvairv2/pkg/repository/mysql/job_category"
	mysqlJobPriority "github.com/your-org/jvairv2/pkg/repository/mysql/job_priority"
	mysqlJobStatus "github.com/your-org/jvairv2/pkg/repository/mysql/job_status"
	mysqlPermission "github.com/your-org/jvairv2/pkg/repository/mysql/permission"
	mysqlProperty "github.com/your-org/jvairv2/pkg/repository/mysql/property"
	mysqlRole "github.com/your-org/jvairv2/pkg/repository/mysql/role"
	mysqlSettings "github.com/your-org/jvairv2/pkg/repository/mysql/settings"
	mysqlTaskStatus "github.com/your-org/jvairv2/pkg/repository/mysql/task_status"
	mysqlTechJobStatus "github.com/your-org/jvairv2/pkg/repository/mysql/technician_job_status"
	mysqlUser "github.com/your-org/jvairv2/pkg/repository/mysql/user"
	mysqlWorkflow "github.com/your-org/jvairv2/pkg/repository/mysql/workflow"
	handler "github.com/your-org/jvairv2/pkg/rest/handler"
	abilityHandler "github.com/your-org/jvairv2/pkg/rest/handler/ability"
	assignedRoleHandler "github.com/your-org/jvairv2/pkg/rest/handler/assigned_role"
	authHandler "github.com/your-org/jvairv2/pkg/rest/handler/auth"
	customerHandler "github.com/your-org/jvairv2/pkg/rest/handler/customer"
	jobCategoryHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_category"
	jobPriorityHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_priority"
	jobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_status"
	permissionHandler "github.com/your-org/jvairv2/pkg/rest/handler/permission"
	propertyHandler "github.com/your-org/jvairv2/pkg/rest/handler/property"
	roleHandler "github.com/your-org/jvairv2/pkg/rest/handler/role"
	settingsHandler "github.com/your-org/jvairv2/pkg/rest/handler/settings"
	taskStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/task_status"
	techJobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/technician_job_status"
	userHandler "github.com/your-org/jvairv2/pkg/rest/handler/user"
	workflowHandler "github.com/your-org/jvairv2/pkg/rest/handler/workflow"
	middleware "github.com/your-org/jvairv2/pkg/rest/middleware"
	router "github.com/your-org/jvairv2/pkg/rest/router"
)

// Container contiene todas las dependencias de la aplicación
type Container struct {
	Config               *configs.Config
	DBConnection         *mysql.Connection
	HealthHandler        *handler.HealthHandler
	AuthHandler          *authHandler.Handler
	UserHandler          *userHandler.Handler
	RoleHandler          *roleHandler.Handler
	AbilityHandler       *abilityHandler.Handler
	AssignedRoleHandler  *assignedRoleHandler.Handler
	PermissionHandler    *permissionHandler.Handler
	SettingsHandler      *settingsHandler.Handler
	WorkflowHandler      *workflowHandler.Handler
	CustomerHandler      *customerHandler.Handler
	PropertyHandler      *propertyHandler.Handler
	JobCategoryHandler   *jobCategoryHandler.Handler
	JobStatusHandler     *jobStatusHandler.Handler
	JobPriorityHandler   *jobPriorityHandler.Handler
	TechJobStatusHandler *techJobStatusHandler.Handler
	TaskStatusHandler    *taskStatusHandler.Handler
	AuthMiddleware       *middleware.AuthMiddleware
	Router               http.Handler
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
	assignedRoleRepo := mysqlAssignedRole.NewRepository(dbConn.GetDB())
	roleRepo := mysqlRole.NewRepository(dbConn.GetDB())
	abilityRepo := mysqlAbility.NewRepository(dbConn.GetDB())
	permissionRepo := mysqlPermission.NewRepository(dbConn.GetDB())
	settingsRepo := mysqlSettings.NewRepository(dbConn.GetDB())
	workflowRepo := mysqlWorkflow.NewRepository(dbConn.GetDB())
	customerRepo := mysqlCustomer.NewRepository(dbConn.GetDB())

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
	userUC := user.NewUseCase(userRepo, assignedRoleRepo, roleRepo)
	roleUC := role.NewUseCase(roleRepo)
	abilityUC := ability.NewUseCase(abilityRepo)
	assignedRoleUC := assignedRole.NewUseCase(assignedRoleRepo, roleRepo)
	permissionUC := permission.NewUseCase(permissionRepo, abilityRepo)
	settingsUC := settings.NewUseCase(settingsRepo)
	workflowUC := workflow.NewUseCase(workflowRepo)
	customerUC := customer.NewUseCase(customerRepo, workflowRepo)
	propertyRepo := mysqlProperty.NewRepository(dbConn.DB)
	propertyUC := property.NewUseCase(propertyRepo, customerRepo)
	jobCategoryRepo := mysqlJobCategory.NewRepository(dbConn.GetDB())
	jobStatusRepo := mysqlJobStatus.NewRepository(dbConn.GetDB())
	jobPriorityRepo := mysqlJobPriority.NewRepository(dbConn.GetDB())
	techJobStatusRepo := mysqlTechJobStatus.NewRepository(dbConn.GetDB())
	taskStatusRepo := mysqlTaskStatus.NewRepository(dbConn.GetDB())
	jobCategoryUC := jobCategory.NewUseCase(jobCategoryRepo)
	jobStatusUC := jobStatus.NewUseCase(jobStatusRepo)
	jobPriorityUC := jobPriority.NewUseCase(jobPriorityRepo)
	techJobStatusUC := techJobStatus.NewUseCase(techJobStatusRepo, jobStatusRepo)
	taskStatusUC := taskStatus.NewUseCase(taskStatusRepo)

	// Inicializar handlers
	healthHandler := handler.NewHealthHandler(dbConn)
	authHandler := authHandler.NewHandler(authUC)

	// Inicializar handlers con sus casos de uso
	userHandler := userHandler.NewHandler(userUC)
	roleHandler := roleHandler.NewHandler(roleUC)
	abilityHandler := abilityHandler.NewHandler(abilityUC)
	assignedRoleHandler := assignedRoleHandler.NewHandler(assignedRoleUC)
	permissionHandler := permissionHandler.NewHandler(permissionUC)
	settingsHandler := settingsHandler.NewHandler(settingsUC)
	workflowHandler := workflowHandler.NewHandler(workflowUC)
	customerHandler := customerHandler.NewHandler(customerUC, propertyUC)
	propHandler := propertyHandler.NewHandler(propertyUC)
	jobCatHandler := jobCategoryHandler.NewHandler(jobCategoryUC)
	jobStatHandler := jobStatusHandler.NewHandler(jobStatusUC)
	jobPrioHandler := jobPriorityHandler.NewHandler(jobPriorityUC)
	techJobStatHandler := techJobStatusHandler.NewHandler(techJobStatusUC)
	taskStatHandler := taskStatusHandler.NewHandler(taskStatusUC)

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
		settingsHandler,
		workflowHandler,
		customerHandler,
		propHandler,
		jobCatHandler,
		jobStatHandler,
		jobPrioHandler,
		techJobStatHandler,
		taskStatHandler,
		authMiddleware,
		userUC,
	)

	return &Container{
		Config:               config,
		DBConnection:         dbConn,
		HealthHandler:        healthHandler,
		AuthHandler:          authHandler,
		UserHandler:          userHandler,
		RoleHandler:          roleHandler,
		AbilityHandler:       abilityHandler,
		AssignedRoleHandler:  assignedRoleHandler,
		PermissionHandler:    permissionHandler,
		SettingsHandler:      settingsHandler,
		WorkflowHandler:      workflowHandler,
		CustomerHandler:      customerHandler,
		PropertyHandler:      propHandler,
		JobCategoryHandler:   jobCatHandler,
		JobStatusHandler:     jobStatHandler,
		JobPriorityHandler:   jobPrioHandler,
		TechJobStatusHandler: techJobStatHandler,
		TaskStatusHandler:    taskStatHandler,
		AuthMiddleware:       authMiddleware,
		Router:               r,
	}, nil
}

// Close cierra todas las conexiones
func (c *Container) Close() error {
	if c.DBConnection != nil {
		return c.DBConnection.Close()
	}
	return nil
}
