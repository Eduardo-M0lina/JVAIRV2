package payroll

import "context"

// Repository define las operaciones de persistencia para payroll
type Repository interface {
	// ListPayrollUsers lista usuarios con sus rates agrupados por estado
	ListPayrollUsers(ctx context.Context, filters PayrollFilters) ([]*PayrollUser, int64, error)

	// GetUserRates obtiene todos los rates de un usuario con filtros opcionales
	GetUserRates(ctx context.Context, userID int64, status *string, limit, offset int) ([]*PayrollRate, int64, error)

	// MarkRatesAsPaid marca los rates especificados como pagados
	MarkRatesAsPaid(ctx context.Context, rateIDs []int64) error

	// MarkRatesAsHolding marca los rates especificados como retenidos
	MarkRatesAsHolding(ctx context.Context, rateIDs []int64) error

	// GetPaystubData obtiene los datos del recibo de pago de un usuario
	GetPaystubData(ctx context.Context, userID int64) (*PaystubData, error)

	// GetStatusIDByLabel obtiene el ID del status por su label
	GetStatusIDByLabel(ctx context.Context, label string) (int64, error)
}

// UserExistsChecker verifica si un usuario existe
type UserExistsChecker interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
}
