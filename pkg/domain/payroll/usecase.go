package payroll

import (
	"context"
	"fmt"
	"time"
)

// Service define las operaciones de negocio para payroll
type Service interface {
	// ListPayroll lista el payroll de todos los usuarios activos
	ListPayroll(ctx context.Context, filters PayrollFilters) (*PayrollListResponse, error)

	// GetUserPayroll obtiene el payroll detallado de un usuario específico
	GetUserPayroll(ctx context.Context, userID int64, status *string, page, pageSize int) ([]*PayrollRate, int64, error)

	// MarkAsPaid marca los rates especificados como pagados
	MarkAsPaid(ctx context.Context, userID int64, rateIDs []int64) error

	// MarkAsHolding marca los rates especificados como retenidos
	MarkAsHolding(ctx context.Context, userID int64, rateIDs []int64) error

	// GetPaystub obtiene el recibo de pago de un usuario
	GetPaystub(ctx context.Context, userID int64) (*PaystubData, error)
}

type service struct {
	repo        Repository
	userChecker UserExistsChecker
}

// NewUseCase crea una nueva instancia del servicio de payroll
func NewUseCase(repo Repository, userChecker UserExistsChecker) Service {
	return &service{
		repo:        repo,
		userChecker: userChecker,
	}
}

func (s *service) ListPayroll(ctx context.Context, filters PayrollFilters) (*PayrollListResponse, error) {
	// Valores por defecto
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}

	users, total, err := s.repo.ListPayrollUsers(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing payroll users: %w", err)
	}

	totalPages := int(total) / filters.PageSize
	if int(total)%filters.PageSize > 0 {
		totalPages++
	}

	return &PayrollListResponse{
		Users:      users,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *service) GetUserPayroll(ctx context.Context, userID int64, status *string, page, pageSize int) ([]*PayrollRate, int64, error) {
	// Verificar que el usuario existe
	exists, err := s.userChecker.UserExists(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return nil, 0, fmt.Errorf("user with ID %d not found", userID)
	}

	// Valores por defecto
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize
	return s.repo.GetUserRates(ctx, userID, status, pageSize, offset)
}

func (s *service) MarkAsPaid(ctx context.Context, userID int64, rateIDs []int64) error {
	if len(rateIDs) == 0 {
		return fmt.Errorf("no rate IDs provided")
	}

	// Verificar que el usuario existe
	exists, err := s.userChecker.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return s.repo.MarkRatesAsPaid(ctx, rateIDs)
}

func (s *service) MarkAsHolding(ctx context.Context, userID int64, rateIDs []int64) error {
	if len(rateIDs) == 0 {
		return fmt.Errorf("no rate IDs provided")
	}

	// Verificar que el usuario existe
	exists, err := s.userChecker.UserExists(ctx, userID)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return s.repo.MarkRatesAsHolding(ctx, rateIDs)
}

func (s *service) GetPaystub(ctx context.Context, userID int64) (*PaystubData, error) {
	// Verificar que el usuario existe
	exists, err := s.userChecker.UserExists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	paystub, err := s.repo.GetPaystubData(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting paystub data: %w", err)
	}

	paystub.GeneratedAt = time.Now()
	return paystub, nil
}
