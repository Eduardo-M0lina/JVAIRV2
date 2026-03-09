package main

import (
	http "net/http"

	configs "github.com/your-org/jvairv2/configs"
	commonAuth "github.com/your-org/jvairv2/pkg/common/auth"
	"github.com/your-org/jvairv2/pkg/common/storage"
	ability "github.com/your-org/jvairv2/pkg/domain/ability"
	domainAlert "github.com/your-org/jvairv2/pkg/domain/alert"
	assignedRole "github.com/your-org/jvairv2/pkg/domain/assigned_role"
	domainAuth "github.com/your-org/jvairv2/pkg/domain/auth"
	customer "github.com/your-org/jvairv2/pkg/domain/customer"
	domainFile "github.com/your-org/jvairv2/pkg/domain/file"
	domainInvoice "github.com/your-org/jvairv2/pkg/domain/invoice"
	domainInvoicePayment "github.com/your-org/jvairv2/pkg/domain/invoice_payment"
	domainJob "github.com/your-org/jvairv2/pkg/domain/job"
	domainJobActivityLog "github.com/your-org/jvairv2/pkg/domain/job_activity_log"
	jobCategory "github.com/your-org/jvairv2/pkg/domain/job_category"
	domainJobEquip "github.com/your-org/jvairv2/pkg/domain/job_equipment"
	jobPriority "github.com/your-org/jvairv2/pkg/domain/job_priority"
	domainJobRate "github.com/your-org/jvairv2/pkg/domain/job_rate"
	domainJobRateStatus "github.com/your-org/jvairv2/pkg/domain/job_rate_status"
	domainJobResident "github.com/your-org/jvairv2/pkg/domain/job_resident"
	jobStatus "github.com/your-org/jvairv2/pkg/domain/job_status"
	domainJobTask "github.com/your-org/jvairv2/pkg/domain/job_task"
	domainJobVisit "github.com/your-org/jvairv2/pkg/domain/job_visit"
	permission "github.com/your-org/jvairv2/pkg/domain/permission"
	property "github.com/your-org/jvairv2/pkg/domain/property"
	domainPropEquip "github.com/your-org/jvairv2/pkg/domain/property_equipment"
	domainQuote "github.com/your-org/jvairv2/pkg/domain/quote"
	quoteStatus "github.com/your-org/jvairv2/pkg/domain/quote_status"
	role "github.com/your-org/jvairv2/pkg/domain/role"
	settings "github.com/your-org/jvairv2/pkg/domain/settings"
	domainSupervisor "github.com/your-org/jvairv2/pkg/domain/supervisor"
	taskStatus "github.com/your-org/jvairv2/pkg/domain/task_status"
	techJobStatus "github.com/your-org/jvairv2/pkg/domain/technician_job_status"
	user "github.com/your-org/jvairv2/pkg/domain/user"
	domainWarranty "github.com/your-org/jvairv2/pkg/domain/warranty"
	domainWarrantyClaim "github.com/your-org/jvairv2/pkg/domain/warranty_claim"
	warrantyClaimStatus "github.com/your-org/jvairv2/pkg/domain/warranty_claim_status"
	warrantyClaimType "github.com/your-org/jvairv2/pkg/domain/warranty_claim_type"
	domainWarrantyEquip "github.com/your-org/jvairv2/pkg/domain/warranty_equipment"
	warrantyStatus "github.com/your-org/jvairv2/pkg/domain/warranty_status"
	warrantyType "github.com/your-org/jvairv2/pkg/domain/warranty_type"
	workflow "github.com/your-org/jvairv2/pkg/domain/workflow"
	mysql "github.com/your-org/jvairv2/pkg/repository/mysql"
	mysqlAbility "github.com/your-org/jvairv2/pkg/repository/mysql/ability"
	mysqlAlert "github.com/your-org/jvairv2/pkg/repository/mysql/alert"
	mysqlAssignedRole "github.com/your-org/jvairv2/pkg/repository/mysql/assigned_role"
	mysqlCustomer "github.com/your-org/jvairv2/pkg/repository/mysql/customer"
	mysqlFile "github.com/your-org/jvairv2/pkg/repository/mysql/file"
	mysqlInvoice "github.com/your-org/jvairv2/pkg/repository/mysql/invoice"
	mysqlInvoicePayment "github.com/your-org/jvairv2/pkg/repository/mysql/invoice_payment"
	mysqlJob "github.com/your-org/jvairv2/pkg/repository/mysql/job"
	mysqlJobActivityLog "github.com/your-org/jvairv2/pkg/repository/mysql/job_activity_log"
	mysqlJobCategory "github.com/your-org/jvairv2/pkg/repository/mysql/job_category"
	mysqlJobEquip "github.com/your-org/jvairv2/pkg/repository/mysql/job_equipment"
	mysqlJobPriority "github.com/your-org/jvairv2/pkg/repository/mysql/job_priority"
	mysqlJobRate "github.com/your-org/jvairv2/pkg/repository/mysql/job_rate"
	mysqlJobRateStatus "github.com/your-org/jvairv2/pkg/repository/mysql/job_rate_status"
	mysqlJobResident "github.com/your-org/jvairv2/pkg/repository/mysql/job_resident"
	mysqlJobStatus "github.com/your-org/jvairv2/pkg/repository/mysql/job_status"
	mysqlJobTask "github.com/your-org/jvairv2/pkg/repository/mysql/job_task"
	mysqlJobVisit "github.com/your-org/jvairv2/pkg/repository/mysql/job_visit"
	mysqlPermission "github.com/your-org/jvairv2/pkg/repository/mysql/permission"
	mysqlProperty "github.com/your-org/jvairv2/pkg/repository/mysql/property"
	mysqlPropEquip "github.com/your-org/jvairv2/pkg/repository/mysql/property_equipment"
	mysqlQuote "github.com/your-org/jvairv2/pkg/repository/mysql/quote"
	mysqlQuoteStatus "github.com/your-org/jvairv2/pkg/repository/mysql/quote_status"
	mysqlRole "github.com/your-org/jvairv2/pkg/repository/mysql/role"
	mysqlSettings "github.com/your-org/jvairv2/pkg/repository/mysql/settings"
	mysqlSupervisor "github.com/your-org/jvairv2/pkg/repository/mysql/supervisor"
	mysqlTaskStatus "github.com/your-org/jvairv2/pkg/repository/mysql/task_status"
	mysqlTechJobStatus "github.com/your-org/jvairv2/pkg/repository/mysql/technician_job_status"
	mysqlUser "github.com/your-org/jvairv2/pkg/repository/mysql/user"
	mysqlWarranty "github.com/your-org/jvairv2/pkg/repository/mysql/warranty"
	mysqlWarrantyClaim "github.com/your-org/jvairv2/pkg/repository/mysql/warranty_claim"
	mysqlWarrantyClaimStatus "github.com/your-org/jvairv2/pkg/repository/mysql/warranty_claim_status"
	mysqlWarrantyClaimType "github.com/your-org/jvairv2/pkg/repository/mysql/warranty_claim_type"
	mysqlWarrantyEquip "github.com/your-org/jvairv2/pkg/repository/mysql/warranty_equipment"
	mysqlWarrantyStatus "github.com/your-org/jvairv2/pkg/repository/mysql/warranty_status"
	mysqlWarrantyType "github.com/your-org/jvairv2/pkg/repository/mysql/warranty_type"
	mysqlWorkflow "github.com/your-org/jvairv2/pkg/repository/mysql/workflow"
	handler "github.com/your-org/jvairv2/pkg/rest/handler"
	abilityHandler "github.com/your-org/jvairv2/pkg/rest/handler/ability"
	alertHandler "github.com/your-org/jvairv2/pkg/rest/handler/alert"
	assignedRoleHandler "github.com/your-org/jvairv2/pkg/rest/handler/assigned_role"
	authHandler "github.com/your-org/jvairv2/pkg/rest/handler/auth"
	customerHandler "github.com/your-org/jvairv2/pkg/rest/handler/customer"
	invoiceHandler "github.com/your-org/jvairv2/pkg/rest/handler/invoice"
	invoicePaymentHandler "github.com/your-org/jvairv2/pkg/rest/handler/invoice_payment"
	jobHandler "github.com/your-org/jvairv2/pkg/rest/handler/job"
	jobActivityLogHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_activity_log"
	jobCategoryHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_category"
	jobEquipHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_equipment"
	jobPriorityHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_priority"
	jobRateHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_rate"
	jobRateStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_rate_status"
	jobResidentHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_resident"
	jobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_status"
	jobTaskHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_task"
	jobVisitHandler "github.com/your-org/jvairv2/pkg/rest/handler/job_visit"
	permissionHandler "github.com/your-org/jvairv2/pkg/rest/handler/permission"
	propertyHandler "github.com/your-org/jvairv2/pkg/rest/handler/property"
	propEquipHandler "github.com/your-org/jvairv2/pkg/rest/handler/property_equipment"
	quoteHandler "github.com/your-org/jvairv2/pkg/rest/handler/quote"
	quoteStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/quote_status"
	roleHandler "github.com/your-org/jvairv2/pkg/rest/handler/role"
	settingsHandler "github.com/your-org/jvairv2/pkg/rest/handler/settings"
	supervisorHandler "github.com/your-org/jvairv2/pkg/rest/handler/supervisor"
	taskStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/task_status"
	techJobStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/technician_job_status"
	userHandler "github.com/your-org/jvairv2/pkg/rest/handler/user"
	warrantyHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty"
	warrantyClaimHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_claim"
	warrantyClaimStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_claim_status"
	warrantyClaimTypeHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_claim_type"
	warrantyEquipHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_equipment"
	warrantyStatusHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_status"
	warrantyTypeHandler "github.com/your-org/jvairv2/pkg/rest/handler/warranty_type"
	workflowHandler "github.com/your-org/jvairv2/pkg/rest/handler/workflow"
	middleware "github.com/your-org/jvairv2/pkg/rest/middleware"
	router "github.com/your-org/jvairv2/pkg/rest/router"
)

// Container contiene todas las dependencias de la aplicación
type Container struct {
	Config                     *configs.Config
	DBConnection               *mysql.Connection
	HealthHandler              *handler.HealthHandler
	AuthHandler                *authHandler.Handler
	UserHandler                *userHandler.Handler
	RoleHandler                *roleHandler.Handler
	AbilityHandler             *abilityHandler.Handler
	AssignedRoleHandler        *assignedRoleHandler.Handler
	PermissionHandler          *permissionHandler.Handler
	SettingsHandler            *settingsHandler.Handler
	WorkflowHandler            *workflowHandler.Handler
	CustomerHandler            *customerHandler.Handler
	PropertyHandler            *propertyHandler.Handler
	JobHandler                 *jobHandler.Handler
	JobCategoryHandler         *jobCategoryHandler.Handler
	JobStatusHandler           *jobStatusHandler.Handler
	JobPriorityHandler         *jobPriorityHandler.Handler
	TechJobStatusHandler       *techJobStatusHandler.Handler
	TaskStatusHandler          *taskStatusHandler.Handler
	QuoteHandler               *quoteHandler.Handler
	QuoteStatusHandler         *quoteStatusHandler.Handler
	SupervisorHandler          *supervisorHandler.Handler
	PropEquipHandler           *propEquipHandler.Handler
	JobEquipHandler            *jobEquipHandler.Handler
	AuthMiddleware             *middleware.AuthMiddleware
	Router                     http.Handler
	InvoiceHandler             *invoiceHandler.Handler
	InvoicePaymentHandler      *invoicePaymentHandler.Handler
	WarrantyTypeHandler        *warrantyTypeHandler.Handler
	WarrantyStatusHandler      *warrantyStatusHandler.Handler
	WarrantyClaimTypeHandler   *warrantyClaimTypeHandler.Handler
	WarrantyClaimStatusHandler *warrantyClaimStatusHandler.Handler
	WarrantyHandler            *warrantyHandler.Handler
	WarrantyEquipHandler       *warrantyEquipHandler.Handler
	WarrantyClaimHandler       *warrantyClaimHandler.Handler
	JobVisitHandler            *jobVisitHandler.Handler
	JobActivityLogHandler      *jobActivityLogHandler.Handler
	JobResidentHandler         *jobResidentHandler.Handler
	JobRateStatusHandler       *jobRateStatusHandler.Handler
	JobTaskHandler             *jobTaskHandler.Handler
	JobRateHandler             *jobRateHandler.Handler
	AlertHandler               *alertHandler.Handler
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
	jobRepo := mysqlJob.NewRepository(dbConn.GetDB())
	jobCategoryChecker := mysqlJob.NewJobCategoryCheckerAdapter(dbConn.GetDB())
	jobPriorityChecker := mysqlJob.NewJobPriorityCheckerAdapter(dbConn.GetDB())
	jobStatusChecker := mysqlJob.NewJobStatusCheckerAdapter(dbConn.GetDB())
	workflowChecker := mysqlJob.NewWorkflowCheckerAdapter(dbConn.GetDB())
	propertyChecker := mysqlJob.NewPropertyCheckerAdapter(dbConn.GetDB())
	userChecker := mysqlJob.NewUserCheckerAdapter(dbConn.GetDB())
	techJobStatusChecker := mysqlJob.NewTechnicianJobStatusCheckerAdapter(dbConn.GetDB())
	jobUC := domainJob.NewUseCase(jobRepo, jobCategoryChecker, jobPriorityChecker, jobStatusChecker, workflowChecker, propertyChecker, userChecker, techJobStatusChecker)
	quoteStatusRepo := mysqlQuoteStatus.NewRepository(dbConn.GetDB())
	quoteStatusUC := quoteStatus.NewUseCase(quoteStatusRepo)
	quoteRepo := mysqlQuote.NewRepository(dbConn.GetDB())
	quoteJobChecker := mysqlQuote.NewJobCheckerAdapter(dbConn.GetDB())
	quoteQSChecker := mysqlQuote.NewQuoteStatusCheckerAdapter(dbConn.GetDB())
	quoteUC := domainQuote.NewUseCase(quoteRepo, quoteJobChecker, quoteQSChecker)
	supervisorRepo := mysqlSupervisor.NewRepository(dbConn.GetDB())
	supervisorUC := domainSupervisor.NewUseCase(supervisorRepo, customerRepo)
	propEquipRepo := mysqlPropEquip.NewRepository(dbConn.GetDB())
	propEquipUC := domainPropEquip.NewUseCase(propEquipRepo, propertyRepo)
	jobEquipRepo := mysqlJobEquip.NewRepository(dbConn.GetDB())
	jobEquipJobChecker := mysqlJobEquip.NewJobCheckerAdapter(dbConn.GetDB())
	jobEquipUC := domainJobEquip.NewUseCase(jobEquipRepo, jobEquipJobChecker)
	invoiceRepo := mysqlInvoice.NewRepository(dbConn.GetDB())
	invoiceJobChecker := mysqlInvoice.NewJobCheckerAdapter(dbConn.GetDB())
	invoiceUC := domainInvoice.NewUseCase(invoiceRepo, invoiceJobChecker)
	invoicePaymentRepo := mysqlInvoicePayment.NewRepository(dbConn.GetDB())
	invoiceChecker := mysqlInvoice.NewInvoiceCheckerAdapter(dbConn.GetDB())
	invoicePaymentUC := domainInvoicePayment.NewUseCase(invoicePaymentRepo, invoiceChecker)

	// Job Visits + Files
	s3Client, err := storage.NewS3Client(&config.S3)
	if err != nil {
		// S3 es opcional - log warning y continuar sin soporte de archivos
		_ = err
	}
	jobVisitRepo := mysqlJobVisit.NewRepository(dbConn.GetDB())
	jobVisitJobChecker := mysqlJobVisit.NewJobExistsChecker(dbConn.GetDB())
	jobVisitUserChecker := mysqlJobVisit.NewUserExistsChecker(dbConn.GetDB())
	jobVisitUC := domainJobVisit.NewUseCase(jobVisitRepo, jobVisitJobChecker, jobVisitUserChecker)
	fileRepo := mysqlFile.NewRepository(dbConn.GetDB())
	var fileUC domainFile.Service
	if s3Client != nil {
		fileUC = domainFile.NewUseCase(fileRepo, s3Client)
	} else {
		fileUC = domainFile.NewUseCase(fileRepo, nil)
	}

	// Warranty catalogs
	warrantyTypeRepo := mysqlWarrantyType.NewRepository(dbConn.GetDB())
	warrantyTypeUC := warrantyType.NewUseCase(warrantyTypeRepo)
	warrantyStatusRepo := mysqlWarrantyStatus.NewRepository(dbConn.GetDB())
	warrantyStatusUC := warrantyStatus.NewUseCase(warrantyStatusRepo)
	warrantyClaimTypeRepo := mysqlWarrantyClaimType.NewRepository(dbConn.GetDB())
	warrantyClaimTypeUC := warrantyClaimType.NewUseCase(warrantyClaimTypeRepo)
	warrantyClaimStatusRepo := mysqlWarrantyClaimStatus.NewRepository(dbConn.GetDB())
	warrantyClaimStatusUC := warrantyClaimStatus.NewUseCase(warrantyClaimStatusRepo)

	// Warranties
	warrantyRepo := mysqlWarranty.NewRepository(dbConn.GetDB())
	warrantyJobChecker := mysqlWarranty.NewJobCheckerAdapter(dbConn.GetDB())
	warrantyTypeChecker := mysqlWarranty.NewWarrantyTypeCheckerAdapter(dbConn.GetDB())
	warrantyStatusChecker := mysqlWarranty.NewWarrantyStatusCheckerAdapter(dbConn.GetDB())
	warrantyUC := domainWarranty.NewUseCase(warrantyRepo, warrantyJobChecker, warrantyTypeChecker, warrantyStatusChecker)

	// Warranty Equipment
	warrantyEquipRepo := mysqlWarrantyEquip.NewRepository(dbConn.GetDB())
	warrantyEquipChecker := mysqlWarrantyEquip.NewWarrantyCheckerAdapter(dbConn.GetDB())
	warrantyEquipUC := domainWarrantyEquip.NewUseCase(warrantyEquipRepo, warrantyEquipChecker)

	// Warranty Claims
	warrantyClaimRepo := mysqlWarrantyClaim.NewRepository(dbConn.GetDB())
	warrantyClaimJobChecker := mysqlWarrantyClaim.NewJobCheckerAdapter(dbConn.GetDB())
	warrantyClaimTypeChecker := mysqlWarrantyClaim.NewClaimTypeCheckerAdapter(dbConn.GetDB())
	warrantyClaimStatusChecker := mysqlWarrantyClaim.NewClaimStatusCheckerAdapter(dbConn.GetDB())
	warrantyClaimUC := domainWarrantyClaim.NewUseCase(warrantyClaimRepo, warrantyClaimJobChecker, warrantyClaimTypeChecker, warrantyClaimStatusChecker)

	// Module 13: Activities and Communications
	// Job Activity Logs
	jobActivityLogRepo := mysqlJobActivityLog.NewRepository(dbConn.GetDB())
	jobActivityLogJobChecker := mysqlJobActivityLog.NewJobExistsChecker(dbConn.GetDB())
	jobActivityLogUserChecker := mysqlJobActivityLog.NewUserExistsChecker(dbConn.GetDB())
	jobActivityLogUC := domainJobActivityLog.NewUseCase(jobActivityLogRepo, jobActivityLogJobChecker, jobActivityLogUserChecker)

	// Job Residents
	jobResidentRepo := mysqlJobResident.NewRepository(dbConn.GetDB())
	jobResidentJobChecker := mysqlJobResident.NewJobExistsChecker(dbConn.GetDB())
	jobResidentUC := domainJobResident.NewUseCase(jobResidentRepo, jobResidentJobChecker)

	// Job Rate Statuses
	jobRateStatusRepo := mysqlJobRateStatus.NewRepository(dbConn.GetDB())
	jobRateStatusUC := domainJobRateStatus.NewUseCase(jobRateStatusRepo)

	// Job Tasks
	jobTaskRepo := mysqlJobTask.NewRepository(dbConn.GetDB())
	jobTaskJobChecker := mysqlJobTask.NewJobExistsChecker(dbConn.GetDB())
	jobTaskUserChecker := mysqlJobTask.NewUserExistsChecker(dbConn.GetDB())
	jobTaskStatusChecker := mysqlJobTask.NewTaskStatusExistsChecker(dbConn.GetDB())
	jobTaskUC := domainJobTask.NewUseCase(jobTaskRepo, jobTaskJobChecker, jobTaskUserChecker, jobTaskStatusChecker)

	// Job Rates
	jobRateRepo := mysqlJobRate.NewRepository(dbConn.GetDB())
	jobRateJobChecker := mysqlJobRate.NewJobExistsChecker(dbConn.GetDB())
	jobRateUserChecker := mysqlJobRate.NewUserExistsChecker(dbConn.GetDB())
	jobRateStatusChecker := mysqlJobRate.NewJobRateStatusExistsChecker(dbConn.GetDB())
	jobRateUC := domainJobRate.NewUseCase(jobRateRepo, jobRateJobChecker, jobRateUserChecker, jobRateStatusChecker)

	// Module 15: Alerts
	alertRepo := mysqlAlert.NewRepository(dbConn.GetDB())
	alertUserChecker := mysqlAlert.NewUserExistsChecker(dbConn.GetDB())
	alertUC := domainAlert.NewUseCase(alertRepo, alertUserChecker)

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
	customerHandler := customerHandler.NewHandler(customerUC, propertyUC, jobUC)
	propHandler := propertyHandler.NewHandler(propertyUC)
	jobCatHandler := jobCategoryHandler.NewHandler(jobCategoryUC)
	jobStatHandler := jobStatusHandler.NewHandler(jobStatusUC)
	jobPrioHandler := jobPriorityHandler.NewHandler(jobPriorityUC)
	techJobStatHandler := techJobStatusHandler.NewHandler(techJobStatusUC)
	taskStatHandler := taskStatusHandler.NewHandler(taskStatusUC)
	jobHdlr := jobHandler.NewHandler(jobUC, jobActivityLogUC)
	quoteHdlr := quoteHandler.NewHandler(quoteUC)
	quoteStatHandler := quoteStatusHandler.NewHandler(quoteStatusUC)
	supervisorHdlr := supervisorHandler.NewHandler(supervisorUC)
	propEquipHdlr := propEquipHandler.NewHandler(propEquipUC)
	jobEquipHdlr := jobEquipHandler.NewHandler(jobEquipUC)
	invHdlr := invoiceHandler.NewHandler(invoiceUC)
	invPayHdlr := invoicePaymentHandler.NewHandler(invoicePaymentUC)
	wtHdlr := warrantyTypeHandler.NewHandler(warrantyTypeUC)
	wsHdlr := warrantyStatusHandler.NewHandler(warrantyStatusUC)
	wctHdlr := warrantyClaimTypeHandler.NewHandler(warrantyClaimTypeUC)
	wcsHdlr := warrantyClaimStatusHandler.NewHandler(warrantyClaimStatusUC)
	wHdlr := warrantyHandler.NewHandler(warrantyUC)
	weHdlr := warrantyEquipHandler.NewHandler(warrantyEquipUC)
	wcHdlr := warrantyClaimHandler.NewHandler(warrantyClaimUC)
	jvHdlr := jobVisitHandler.NewHandler(jobVisitUC, fileUC)
	jalHdlr := jobActivityLogHandler.NewHandler(jobActivityLogUC)
	jrHdlr := jobResidentHandler.NewHandler(jobResidentUC)
	jrsHdlr := jobRateStatusHandler.NewHandler(jobRateStatusUC)
	jtHdlr := jobTaskHandler.NewHandler(jobTaskUC)
	jraHdlr := jobRateHandler.NewHandler(jobRateUC)
	alHdlr := alertHandler.NewHandler(alertUC)

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
		jobHdlr,
		jobCatHandler,
		jobStatHandler,
		jobPrioHandler,
		techJobStatHandler,
		taskStatHandler,
		quoteHdlr,
		quoteStatHandler,
		supervisorHdlr,
		propEquipHdlr,
		jobEquipHdlr,
		invHdlr,
		invPayHdlr,
		wtHdlr,
		wsHdlr,
		wctHdlr,
		wcsHdlr,
		wHdlr,
		weHdlr,
		wcHdlr,
		jvHdlr,
		jalHdlr,
		jrHdlr,
		jrsHdlr,
		jtHdlr,
		jraHdlr,
		alHdlr,
		authMiddleware,
		userUC,
	)

	return &Container{
		Config:                     config,
		DBConnection:               dbConn,
		HealthHandler:              healthHandler,
		AuthHandler:                authHandler,
		UserHandler:                userHandler,
		RoleHandler:                roleHandler,
		AbilityHandler:             abilityHandler,
		AssignedRoleHandler:        assignedRoleHandler,
		PermissionHandler:          permissionHandler,
		SettingsHandler:            settingsHandler,
		WorkflowHandler:            workflowHandler,
		CustomerHandler:            customerHandler,
		PropertyHandler:            propHandler,
		JobHandler:                 jobHdlr,
		JobCategoryHandler:         jobCatHandler,
		JobStatusHandler:           jobStatHandler,
		JobPriorityHandler:         jobPrioHandler,
		TechJobStatusHandler:       techJobStatHandler,
		TaskStatusHandler:          taskStatHandler,
		QuoteHandler:               quoteHdlr,
		QuoteStatusHandler:         quoteStatHandler,
		SupervisorHandler:          supervisorHdlr,
		PropEquipHandler:           propEquipHdlr,
		JobEquipHandler:            jobEquipHdlr,
		AuthMiddleware:             authMiddleware,
		Router:                     r,
		InvoiceHandler:             invHdlr,
		InvoicePaymentHandler:      invPayHdlr,
		WarrantyTypeHandler:        wtHdlr,
		WarrantyStatusHandler:      wsHdlr,
		WarrantyClaimTypeHandler:   wctHdlr,
		WarrantyClaimStatusHandler: wcsHdlr,
		WarrantyHandler:            wHdlr,
		WarrantyEquipHandler:       weHdlr,
		WarrantyClaimHandler:       wcHdlr,
		JobVisitHandler:            jvHdlr,
		JobActivityLogHandler:      jalHdlr,
		JobResidentHandler:         jrHdlr,
		JobRateStatusHandler:       jrsHdlr,
		JobTaskHandler:             jtHdlr,
		JobRateHandler:             jraHdlr,
		AlertHandler:               alHdlr,
	}, nil
}

// Close cierra todas las conexiones
func (c *Container) Close() error {
	if c.DBConnection != nil {
		return c.DBConnection.Close()
	}
	return nil
}
