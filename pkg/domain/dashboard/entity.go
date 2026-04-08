package dashboard

// DashboardStats contiene las estadísticas generales del dashboard
type DashboardStats struct {
	JobsAwaitingDispatch int `json:"jobsAwaitingDispatch"`
	JobsDispatched       int `json:"jobsDispatched"`
	JobsUrgent           int `json:"jobsUrgent"`
	JobsOpen             int `json:"jobsOpen"`
	JobsClosedThisMonth  int `json:"jobsClosedThisMonth"`
}

// DashboardJob representa un job resumido para el dashboard
type DashboardJob struct {
	ID            int64   `json:"id"`
	WorkOrder     *string `json:"workOrder,omitempty"`
	DateReceived  string  `json:"dateReceived"`
	JobCategoryID int64   `json:"jobCategoryId"`
	JobPriorityID int64   `json:"jobPriorityId"`
	JobStatusID   int64   `json:"jobStatusId"`
	PropertyID    int64   `json:"propertyId"`
	UserID        *int64  `json:"userId,omitempty"`
	DispatchDate  *string `json:"dispatchDate,omitempty"`
	DueDate       *string `json:"dueDate,omitempty"`
	QuickNotes    *string `json:"quickNotes,omitempty"`
	Closed        bool    `json:"closed"`
	// Campos adicionales del job
	CallAttempted  bool     `json:"callAttempted"`
	CageRequired   bool     `json:"cageRequired"`
	WeekNumber     *int     `json:"weekNumber,omitempty"`
	RouteNumber    *int     `json:"routeNumber,omitempty"`
	JobSalesPrice  *float64 `json:"jobSalesPrice,omitempty"`
	ScheduledTime  *string  `json:"scheduledTime,omitempty"`
	CompletionDate *string  `json:"completionDate,omitempty"`
	// Campos enriquecidos con JOINs
	PropertyStreet *string `json:"propertyStreet,omitempty"`
	PropertyCity   *string `json:"propertyCity,omitempty"`
	PropertyZip    *string `json:"propertyZip,omitempty"`
	CustomerName   *string `json:"customerName,omitempty"`
	TechnicianName *string `json:"technicianName,omitempty"`
	CategoryName   *string `json:"categoryName,omitempty"`
	PriorityName   *string `json:"priorityName,omitempty"`
	PriorityClass  *string `json:"priorityClass,omitempty"`
	StatusName     *string `json:"statusName,omitempty"`
	StatusClass    *string `json:"statusClass,omitempty"`
}

// AdminDashboard respuesta del dashboard para administradores
type AdminDashboard struct {
	Stats                DashboardStats `json:"stats"`
	JobsAwaitingDispatch []DashboardJob `json:"jobsAwaitingDispatch"`
	JobsUrgent           []DashboardJob `json:"jobsUrgent"`
}

// TechnicianDashboard respuesta del dashboard para técnicos
type TechnicianDashboard struct {
	Stats          DashboardStats `json:"stats"`
	JobsDispatched []DashboardJob `json:"jobsDispatched"`
	JobsUrgent     []DashboardJob `json:"jobsUrgent"`
}
