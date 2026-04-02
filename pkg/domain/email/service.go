package email

import (
	"context"
	"fmt"
	"time"

	domainPayroll "github.com/your-org/jvairv2/pkg/domain/payroll"
	infraEmail "github.com/your-org/jvairv2/pkg/infrastructure/email"
)

// Service define la interfaz del servicio de email
type Service interface {
	SendDispatchEmail(ctx context.Context, jobID int64, recipients []string) error
	SendDispatchSupervisorEmail(ctx context.Context, jobID int64, subject string, body string, recipients []string) error
	SendInvoiceEmail(ctx context.Context, invoiceID int64, recipients []string) error
	SendQuoteEmail(ctx context.Context, quoteID int64, recipients []string) error
	SendTaskNotificationEmail(ctx context.Context, taskID int64, recipients []string) error
	SendPayStubEmail(ctx context.Context, userID int64, recipients []string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, resetLink string) error
}

// EmailService implementa el servicio de email
type EmailService struct {
	emailService *infraEmail.Service
	jobRepo      JobRepository
	propertyRepo PropertyRepository
	customerRepo CustomerRepository
	userRepo     UserRepository
	residentRepo ResidentRepository
	invoiceRepo  InvoiceRepository
	quoteRepo    QuoteRepository
	taskRepo     TaskRepository
	payrollRepo  PayrollRepository
}

// JobRepository define los métodos necesarios del repositorio de jobs
type JobRepository interface {
	GetByID(ctx context.Context, id int64) (*JobData, error)
}

// PropertyRepository define los métodos necesarios del repositorio de properties
type PropertyRepository interface {
	GetByID(ctx context.Context, id int64) (*PropertyData, error)
}

// CustomerRepository define los métodos necesarios del repositorio de customers
type CustomerRepository interface {
	GetByID(ctx context.Context, id int64) (*CustomerData, error)
}

// UserRepository define los métodos necesarios del repositorio de users
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*UserData, error)
}

// ResidentRepository define los métodos necesarios del repositorio de residents
type ResidentRepository interface {
	ListByJobID(ctx context.Context, jobID int64) ([]*ResidentData, error)
}

// InvoiceRepository define los métodos necesarios del repositorio de invoices
type InvoiceRepository interface {
	GetByID(ctx context.Context, id int64) (*InvoiceData, error)
}

// QuoteRepository define los métodos necesarios del repositorio de quotes
type QuoteRepository interface {
	GetByID(ctx context.Context, id int64) (*QuoteData, error)
}

// TaskRepository define los métodos necesarios del repositorio de tasks
type TaskRepository interface {
	GetByID(ctx context.Context, id int64) (*TaskData, error)
}

// PayrollRepository define los métodos necesarios del repositorio de payroll
type PayrollRepository interface {
	GetPaystubData(ctx context.Context, userID int64) (*domainPayroll.PaystubData, error)
}

// JobData representa los datos de un job
type JobData struct {
	ID            int64
	WorkOrder     *string
	DispatchDate  *string
	DispatchNotes *string
	PropertyID    int64
	UserID        *int64
}

// PropertyData representa los datos de una propiedad
type PropertyData struct {
	ID         int64
	Street     string
	City       string
	State      string
	Zip        string
	CustomerID int64
}

// CustomerData representa los datos de un cliente
type CustomerData struct {
	ID   int64
	Name string
}

// UserData representa los datos de un usuario
type UserData struct {
	ID   int64
	Name string
}

// ResidentData representa los datos de un residente
type ResidentData struct {
	Name        string
	MobilePhone *string
	HomePhone   *string
}

// InvoiceData representa los datos de una factura
type InvoiceData struct {
	ID            int64
	InvoiceNumber string
	Amount        float64
	Description   *string
	Notes         *string
	JobID         int64
}

// QuoteData representa los datos de una cotización
type QuoteData struct {
	ID          int64
	QuoteNumber string
	Amount      float64
	Description *string
	Notes       *string
	Status      string
	JobID       int64
}

// TaskData representa los datos de una tarea
type TaskData struct {
	ID          int64
	Description string
	DueDate     *string
	Status      string
	Notes       *string
	JobID       int64
	UserID      int64
}

// NewEmailService crea una nueva instancia del servicio de email
func NewEmailService(
	emailService *infraEmail.Service,
	jobRepo JobRepository,
	propertyRepo PropertyRepository,
	customerRepo CustomerRepository,
	userRepo UserRepository,
	residentRepo ResidentRepository,
	invoiceRepo InvoiceRepository,
	quoteRepo QuoteRepository,
	taskRepo TaskRepository,
	payrollRepo PayrollRepository,
) Service {
	return &EmailService{
		emailService: emailService,
		jobRepo:      jobRepo,
		propertyRepo: propertyRepo,
		customerRepo: customerRepo,
		userRepo:     userRepo,
		residentRepo: residentRepo,
		invoiceRepo:  invoiceRepo,
		quoteRepo:    quoteRepo,
		taskRepo:     taskRepo,
		payrollRepo:  payrollRepo,
	}
}

// SendDispatchEmail envía un email de dispatch
func (s *EmailService) SendDispatchEmail(ctx context.Context, jobID int64, recipients []string) error {
	// Obtener job
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("error al obtener job: %w", err)
	}

	// Obtener property
	property, err := s.propertyRepo.GetByID(ctx, job.PropertyID)
	if err != nil {
		return fmt.Errorf("error al obtener property: %w", err)
	}

	// Obtener customer
	customer, err := s.customerRepo.GetByID(ctx, property.CustomerID)
	if err != nil {
		return fmt.Errorf("error al obtener customer: %w", err)
	}

	// Obtener technician
	technicianName := "Unassigned"
	if job.UserID != nil {
		user, err := s.userRepo.GetByID(ctx, *job.UserID)
		if err == nil {
			technicianName = user.Name
		}
	}

	// Obtener residents
	residents, _ := s.residentRepo.ListByJobID(ctx, jobID)

	// Preparar address
	propertyAddress := fmt.Sprintf("%s, %s, %s %s", property.Street, property.City, property.State, property.Zip)

	// Preparar datos para template
	workOrder := ""
	if job.WorkOrder != nil {
		workOrder = *job.WorkOrder
	}

	dispatchDate := "No Dispatch Date"
	if job.DispatchDate != nil {
		dispatchDate = *job.DispatchDate
	}

	dispatchNotes := ""
	if job.DispatchNotes != nil {
		dispatchNotes = *job.DispatchNotes
	}

	// Convertir residents
	type Resident struct {
		Name        string
		MobilePhone string
		HomePhone   string
	}
	var residentsList []Resident
	for _, r := range residents {
		resident := Resident{Name: r.Name}
		if r.MobilePhone != nil {
			resident.MobilePhone = *r.MobilePhone
		}
		if r.HomePhone != nil {
			resident.HomePhone = *r.HomePhone
		}
		residentsList = append(residentsList, resident)
	}

	templateData := map[string]interface{}{
		"Subject": fmt.Sprintf("NEW JOB: %s", propertyAddress),
		"Job": map[string]interface{}{
			"PropertyAddress": propertyAddress,
			"DispatchDate":    dispatchDate,
			"TechnicianName":  technicianName,
			"CustomerName":    customer.Name,
			"WorkOrder":       workOrder,
			"DispatchNotes":   dispatchNotes,
			"Residents":       residentsList,
		},
	}

	// Enviar email usando template estático dispatch.html
	params := infraEmail.SendEmailParams{
		To:      recipients,
		Subject: fmt.Sprintf("NEW JOB: %s", propertyAddress),
	}

	return s.emailService.SendTemplatedEmail(ctx, "dispatch.html", templateData, params)
}

// SendInvoiceEmail envía un email de factura
func (s *EmailService) SendInvoiceEmail(ctx context.Context, invoiceID int64, recipients []string) error {
	// Obtener invoice
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return fmt.Errorf("error al obtener invoice: %w", err)
	}

	// Formatear fecha (usar fecha actual)
	date := time.Now().Format("January 2, 2006")

	// Preparar datos del template
	templateData := map[string]interface{}{
		"Invoice": map[string]interface{}{
			"InvoiceNumber":       invoice.InvoiceNumber,
			"Date":                date,
			"Balance":             fmt.Sprintf("%.2f", invoice.Amount),
			"Description":         getStringValue(invoice.Description),
			"AllowOnlinePayments": false, // TODO: Implementar lógica de pagos online
			"HasBalance":          invoice.Amount > 0,
			"PaymentURL":          "", // TODO: Generar URL de pago
		},
	}

	// Enviar email usando template estático invoice.html
	params := infraEmail.SendEmailParams{
		To:      recipients,
		Subject: fmt.Sprintf("Invoice # %s from JVAIR", invoice.InvoiceNumber),
	}

	return s.emailService.SendTemplatedEmail(ctx, "invoice.html", templateData, params)
}

// SendQuoteEmail envía un email de cotización
func (s *EmailService) SendQuoteEmail(ctx context.Context, quoteID int64, recipients []string) error {
	// Obtener quote
	quote, err := s.quoteRepo.GetByID(ctx, quoteID)
	if err != nil {
		return fmt.Errorf("error al obtener quote: %w", err)
	}

	// Obtener job
	job, err := s.jobRepo.GetByID(ctx, quote.JobID)
	if err != nil {
		return fmt.Errorf("error al obtener job: %w", err)
	}

	// Formatear fecha (usar fecha actual)
	date := time.Now().Format("January 2, 2006")

	// Preparar datos del template
	templateData := map[string]interface{}{
		"Quote": map[string]interface{}{
			"Date":        date,
			"WorkOrder":   getStringValue(job.WorkOrder),
			"Amount":      fmt.Sprintf("%.2f", quote.Amount),
			"Description": getStringValue(quote.Description),
		},
	}

	// Enviar email usando template estático quote.html
	params := infraEmail.SendEmailParams{
		To:      recipients,
		Subject: "Quote from JVAIR",
	}

	return s.emailService.SendTemplatedEmail(ctx, "quote.html", templateData, params)
}

// SendTaskNotificationEmail envía un email de notificación de tarea
func (s *EmailService) SendTaskNotificationEmail(ctx context.Context, taskID int64, recipients []string) error {
	// Obtener task
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("error al obtener task: %w", err)
	}

	// Obtener job
	job, err := s.jobRepo.GetByID(ctx, task.JobID)
	if err != nil {
		return fmt.Errorf("error al obtener job: %w", err)
	}

	// Obtener property
	property, err := s.propertyRepo.GetByID(ctx, job.PropertyID)
	if err != nil {
		return fmt.Errorf("error al obtener property: %w", err)
	}

	// Obtener customer
	customer, err := s.customerRepo.GetByID(ctx, property.CustomerID)
	if err != nil {
		return fmt.Errorf("error al obtener customer: %w", err)
	}

	// Formatear dirección de la propiedad
	propertyAddress := fmt.Sprintf("%s, %s, %s %s", property.Street, property.City, property.State, property.Zip)

	// Formatear due date
	dueDate := "No Due Date"
	if task.DueDate != nil {
		dueDate = *task.DueDate
	}

	// Preparar datos del template
	templateData := map[string]interface{}{
		"Task": map[string]interface{}{
			"DueDate":         dueDate,
			"CustomerName":    customer.Name,
			"PropertyAddress": propertyAddress,
			"WorkOrder":       getStringValue(job.WorkOrder),
			"TaskDescription": task.Description,
		},
	}

	// Enviar email usando template estático task-notification.html
	params := infraEmail.SendEmailParams{
		To:      recipients,
		Subject: "NEW TASK ASSIGNMENT",
	}

	return s.emailService.SendTemplatedEmail(ctx, "task-notification.html", templateData, params)
}

// SendDispatchSupervisorEmail envía un email a supervisores con cuerpo personalizado e información del job
func (s *EmailService) SendDispatchSupervisorEmail(ctx context.Context, jobID int64, subject string, body string, recipients []string) error {
	// Obtener job
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("error al obtener job: %w", err)
	}

	// Obtener property
	property, err := s.propertyRepo.GetByID(ctx, job.PropertyID)
	if err != nil {
		return fmt.Errorf("error al obtener property: %w", err)
	}

	// Obtener customer
	customer, err := s.customerRepo.GetByID(ctx, property.CustomerID)
	if err != nil {
		return fmt.Errorf("error al obtener customer: %w", err)
	}

	// Obtener technician si está asignado
	var technicianName string
	if job.UserID != nil {
		user, err := s.userRepo.GetByID(ctx, *job.UserID)
		if err == nil {
			technicianName = user.Name
		}
	}

	// Formatear dirección de la propiedad
	propertyAddress := fmt.Sprintf("%s, %s, %s %s", property.Street, property.City, property.State, property.Zip)

	// Preparar datos del template
	templateData := map[string]interface{}{
		"Job": map[string]interface{}{
			"WorkOrder":     getStringValue(job.WorkOrder),
			"DispatchDate":  getStringValue(job.DispatchDate),
			"DispatchNotes": getStringValue(job.DispatchNotes),
		},
		"Property": map[string]interface{}{
			"Address": propertyAddress,
		},
		"Customer": map[string]interface{}{
			"Name": customer.Name,
		},
		"Technician": technicianName,
		"Body":       body,
	}

	// Enviar email usando template estático dispatch-supervisor.html
	params := infraEmail.SendEmailParams{
		To:      recipients,
		Subject: subject,
	}

	return s.emailService.SendTemplatedEmail(ctx, "dispatch-supervisor.html", templateData, params)
}

// SendPayStubEmail envía un email de recibo de pago a un usuario
func (s *EmailService) SendPayStubEmail(ctx context.Context, userID int64, recipients []string) error {
	// Obtener datos del paystub desde el repositorio de payroll
	paystub, err := s.payrollRepo.GetPaystubData(ctx, userID)
	if err != nil {
		return fmt.Errorf("error al obtener datos del paystub: %w", err)
	}

	// Debug: Log de datos recibidos
	fmt.Printf("[DEBUG] SendPayStubEmail - UserID: %d, User: %s, Rates count: %d, TotalPayment: %.2f\n",
		userID, paystub.User.Name, len(paystub.Rates), paystub.TotalPayment)
	for i, rate := range paystub.Rates {
		fmt.Printf("[DEBUG] Rate[%d]: ID=%d, WorkOrder=%v, Payment=%.2f\n",
			i, rate.ID, rate.WorkOrder, rate.Payment)
	}

	// Formatear fecha actual
	date := time.Now().Format("January 2, 2006")

	// Convertir rates a formato para el template
	unpaidRates := make([]map[string]interface{}, 0, len(paystub.Rates))
	for _, rate := range paystub.Rates {
		// Formatear fecha de completado
		completionDate := ""
		if rate.CreatedAt != nil {
			completionDate = rate.CreatedAt.Format("Jan 2, 2006")
		}

		rateMap := map[string]interface{}{
			"CompletionDate": completionDate,
			"WorkOrder":      getStringValue(rate.WorkOrder),
			"PropertyStreet": getStringValue(rate.PropertyAddr),
			"Payment":        fmt.Sprintf("%.2f", rate.Payment),
		}
		unpaidRates = append(unpaidRates, rateMap)
	}

	// Calcular totales
	totals := map[string]interface{}{
		"CurrentPeriodJobs":  len(paystub.Rates),
		"CurrentPeriodTotal": fmt.Sprintf("%.2f", paystub.TotalPayment),
		"YTDJobs":            len(paystub.Rates),
		"YTDTotal":           fmt.Sprintf("%.2f", paystub.TotalPayment),
		"AllTimeJobs":        len(paystub.Rates),
		"AllTimeTotal":       fmt.Sprintf("%.2f", paystub.TotalPayment),
	}

	// Preparar datos del template
	templateData := map[string]interface{}{
		"User": map[string]interface{}{
			"Name": paystub.User.Name,
		},
		"Date":        date,
		"UnpaidRates": unpaidRates,
		"Totals":      totals,
	}

	// Construir subject dinámico
	subject := fmt.Sprintf("Pay Report for %s as of %s", paystub.User.Name, date)

	// Enviar email usando template estático paystub.html
	params := infraEmail.SendEmailParams{
		To:      recipients,
		Subject: subject,
	}

	return s.emailService.SendTemplatedEmail(ctx, "paystub.html", templateData, params)
}

// SendPasswordResetEmail envía un email de recuperación de contraseña
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, toEmail, resetLink string) error {
	// Preparar datos del template
	templateData := map[string]interface{}{
		"ResetLink": resetLink,
		"Year":      time.Now().Year(),
	}

	// Enviar email usando template estático password-reset.html
	params := infraEmail.SendEmailParams{
		To:      []string{toEmail},
		Subject: "Recuperación de Contraseña - JVAIR",
	}

	return s.emailService.SendTemplatedEmail(ctx, "password-reset.html", templateData, params)
}

// getStringValue retorna el valor de un puntero a string o un valor por defecto
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
