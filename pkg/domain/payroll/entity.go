package payroll

import (
	"time"
)

// PayrollUser representa un usuario con su información de payroll
type PayrollUser struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	RoleName *string `json:"roleName,omitempty"`

	// Rates agrupados por estado
	UnpaidRates  []*PayrollRate `json:"unpaidRates"`
	HoldingRates []*PayrollRate `json:"holdingRates"`
	PaidRates    []*PayrollRate `json:"paidRates"`

	// Totales calculados
	TotalUnpaid  float64 `json:"totalUnpaid"`
	TotalHolding float64 `json:"totalHolding"`
	TotalPaid    float64 `json:"totalPaid"`
}

// PayrollRate representa un job_rate con información enriquecida para payroll
type PayrollRate struct {
	ID              int64      `json:"id"`
	JobID           int64      `json:"jobId"`
	UserID          int64      `json:"userId"`
	JobRateStatusID int64      `json:"jobRateStatusId"`
	SalePrice       float64    `json:"salePrice"`
	RatePercent     float64    `json:"ratePercent"`
	RateFlat        float64    `json:"rateFlat"`
	TechParts       float64    `json:"techParts"`
	CompanyParts    float64    `json:"companyParts"`
	PartsReplaced   *string    `json:"partsReplaced,omitempty"`
	Deduction       float64    `json:"deduction"`
	Payment         float64    `json:"payment"`
	Paid            bool       `json:"paid"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`

	// Información del job (enriquecida)
	WorkOrder    *string `json:"workOrder,omitempty"`
	PropertyAddr *string `json:"propertyAddress,omitempty"`

	// Información del status
	StatusLabel *string `json:"statusLabel,omitempty"`
	StatusClass *string `json:"statusClass,omitempty"`
}

// PayrollListResponse representa la respuesta paginada de payroll
type PayrollListResponse struct {
	Users      []*PayrollUser `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
}

// PaystubData representa los datos del recibo de pago de un usuario
type PaystubData struct {
	User         *PayrollUser   `json:"user"`
	Rates        []*PayrollRate `json:"rates"`
	TotalPayment float64        `json:"totalPayment"`
	GeneratedAt  time.Time      `json:"generatedAt"`
}

// MarkRatesRequest representa la solicitud para marcar rates como pagados/retenidos
type MarkRatesRequest struct {
	RateIDs []int64 `json:"rateIds"`
}

// PayrollFilters representa los filtros para listar payroll
type PayrollFilters struct {
	UserID   *int64  `json:"userId,omitempty"`
	Status   *string `json:"status,omitempty"` // "unpaid", "holding", "paid", "all"
	Search   *string `json:"search,omitempty"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}
