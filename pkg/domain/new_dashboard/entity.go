package new_dashboard

import "time"

// TimeRange representa los rangos de tiempo disponibles para filtrar el dashboard
type TimeRange string

const (
	TimeRange7Days     TimeRange = "7days"
	TimeRange30Days    TimeRange = "30days"
	TimeRange90Days    TimeRange = "90days"
	TimeRangeThisMonth TimeRange = "thisMonth"
	TimeRangeLastMonth TimeRange = "lastMonth"
	TimeRangeThisYear  TimeRange = "thisYear"
)

// EnhancedStats contiene estadísticas expandidas del dashboard
type EnhancedStats struct {
	// Stats originales
	JobsAwaitingDispatch int `json:"jobsAwaitingDispatch"`
	JobsDispatched       int `json:"jobsDispatched"`
	JobsUrgent           int `json:"jobsUrgent"`
	JobsOpen             int `json:"jobsOpen"`
	JobsClosedThisMonth  int `json:"jobsClosedThisMonth"`
	// Conteos de nuevos nodos
	TotalAlerts          int `json:"totalAlerts"`
	TotalTasksPending    int `json:"totalTasksPending"`
	TotalTasksOverdue    int `json:"totalTasksOverdue"`
	TotalInvoicesPending int `json:"totalInvoicesPending,omitempty"`
	TotalQuotesPending   int `json:"totalQuotesPending,omitempty"`
	TotalWarrantyClaims  int `json:"totalWarrantyClaims,omitempty"`
}

// AlertSummary contiene el resumen de alertas
type AlertSummary struct {
	UnreadCount int      `json:"unreadCount"`
	Alerts      []*Alert `json:"alerts"`
}

// Alert representa una alerta individual
type Alert struct {
	ID           int64     `json:"id"`
	AlertType    string    `json:"alertType"`
	Message      string    `json:"message"`
	MessageLevel string    `json:"messageLevel"`
	CreatedAt    time.Time `json:"createdAt"`
}

// TaskSummary contiene el resumen de tareas
type TaskSummary struct {
	TotalPending int     `json:"totalPending"`
	TotalOverdue int     `json:"totalOverdue"`
	Tasks        []*Task `json:"tasks"`
}

// Task representa una tarea individual
type Task struct {
	ID         int64     `json:"id"`
	JobID      int64     `json:"jobId"`
	WorkOrder  *string   `json:"workOrder,omitempty"`
	Title      string    `json:"title"`
	DueDate    time.Time `json:"dueDate"`
	StatusName string    `json:"statusName"`
	IsOverdue  bool      `json:"isOverdue"`
}

// Activity representa una actividad reciente
type Activity struct {
	ID        int64     `json:"id"`
	JobID     int64     `json:"jobId"`
	WorkOrder *string   `json:"workOrder,omitempty"`
	Type      string    `json:"type"`
	Log       string    `json:"log"`
	UserName  string    `json:"userName"`
	CreatedAt time.Time `json:"createdAt"`
}

// InvoiceSummary contiene el resumen de facturación
type InvoiceSummary struct {
	TotalPending        int              `json:"totalPending"`
	TotalPaid           int              `json:"totalPaid"`
	TotalOverdue        int              `json:"totalOverdue"`
	AmountPending       float64          `json:"amountPending"`
	AmountPaidThisMonth float64          `json:"amountPaidThisMonth"`
	RecentInvoices      []*RecentInvoice `json:"recentInvoices"`
}

// RecentInvoice representa una factura reciente
type RecentInvoice struct {
	ID            int64     `json:"id"`
	InvoiceNumber string    `json:"invoiceNumber"`
	JobID         *int64    `json:"jobId,omitempty"`
	CustomerName  string    `json:"customerName"`
	Total         float64   `json:"total"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
}

// QuoteSummary contiene el resumen de cotizaciones
type QuoteSummary struct {
	TotalPending  int            `json:"totalPending"`
	TotalApproved int            `json:"totalApproved"`
	TotalRejected int            `json:"totalRejected"`
	RecentQuotes  []*RecentQuote `json:"recentQuotes"`
}

// RecentQuote representa una cotización reciente
type RecentQuote struct {
	ID           int64     `json:"id"`
	QuoteNumber  string    `json:"quoteNumber"`
	JobID        *int64    `json:"jobId,omitempty"`
	CustomerName string    `json:"customerName"`
	Total        float64   `json:"total"`
	StatusName   string    `json:"statusName"`
	CreatedAt    time.Time `json:"createdAt"`
}

// WarrantySummary contiene el resumen de garantías
type WarrantySummary struct {
	ActiveWarranties  int            `json:"activeWarranties"`
	ExpiringThisMonth int            `json:"expiringThisMonth"`
	OpenClaims        int            `json:"openClaims"`
	RecentClaims      []*RecentClaim `json:"recentClaims"`
}

// RecentClaim representa un reclamo de garantía reciente
type RecentClaim struct {
	ID           int64     `json:"id"`
	ClaimNumber  string    `json:"claimNumber"`
	JobID        *int64    `json:"jobId,omitempty"`
	CustomerName string    `json:"customerName"`
	StatusName   string    `json:"statusName"`
	CreatedAt    time.Time `json:"createdAt"`
}

// CategoryCount representa el conteo de jobs por categoría
type CategoryCount struct {
	CategoryID   int64  `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Count        int    `json:"count"`
}

// StatusCount representa el conteo de jobs por estado
type StatusCount struct {
	StatusID   int64  `json:"statusId"`
	StatusName string `json:"statusName"`
	Count      int    `json:"count"`
}

// DueJob representa un job con vencimiento próximo
type DueJob struct {
	ID             int64     `json:"id"`
	WorkOrder      *string   `json:"workOrder,omitempty"`
	DueDate        time.Time `json:"dueDate"`
	CustomerName   string    `json:"customerName"`
	PropertyStreet string    `json:"propertyStreet"`
	StatusName     string    `json:"statusName"`
	PriorityName   string    `json:"priorityName"`
}

// TechnicianLoad representa la carga de trabajo de un técnico
type TechnicianLoad struct {
	UserID     int64  `json:"userId"`
	Name       string `json:"name"`
	OpenJobs   int    `json:"openJobs"`
	UrgentJobs int    `json:"urgentJobs"`
}

// AdminEnhancedDashboard es la respuesta completa para administradores
type AdminEnhancedDashboard struct {
	Stats              EnhancedStats     `json:"stats"`
	AlertSummary       *AlertSummary     `json:"alertSummary"`
	TaskSummary        *TaskSummary      `json:"taskSummary"`
	RecentActivity     []*Activity       `json:"recentActivity"`
	InvoiceSummary     *InvoiceSummary   `json:"invoiceSummary"`
	QuoteSummary       *QuoteSummary     `json:"quoteSummary"`
	WarrantySummary    *WarrantySummary  `json:"warrantySummary"`
	JobsByCategory     []*CategoryCount  `json:"jobsByCategory"`
	JobsByStatus       []*StatusCount    `json:"jobsByStatus"`
	JobsDueThisWeek    []*DueJob         `json:"jobsDueThisWeek"`
	TechnicianWorkload []*TechnicianLoad `json:"technicianWorkload"`
}

// TechnicianEnhancedDashboard es la respuesta para técnicos (filtrada por user_id)
type TechnicianEnhancedDashboard struct {
	Stats           EnhancedStats  `json:"stats"`
	AlertSummary    *AlertSummary  `json:"alertSummary"`
	TaskSummary     *TaskSummary   `json:"taskSummary"`
	RecentActivity  []*Activity    `json:"recentActivity"`
	JobsByStatus    []*StatusCount `json:"jobsByStatus"`
	JobsDueThisWeek []*DueJob      `json:"jobsDueThisWeek"`
}
