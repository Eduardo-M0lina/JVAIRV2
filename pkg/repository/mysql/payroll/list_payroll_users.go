package payroll

import (
	"context"
	"fmt"
	"strings"

	"github.com/angumol/jvairv2/pkg/domain/payroll"
)

func (r *repository) ListPayrollUsers(ctx context.Context, filters payroll.PayrollFilters) ([]*payroll.PayrollUser, int64, error) {
	// Primero obtenemos los usuarios activos con paginación
	baseQuery := `
		SELECT DISTINCT u.id, u.name, u.email, ro.name as role_name
		FROM users u
		LEFT JOIN assigned_roles ar ON ar.entity_id = u.id AND ar.entity_type = 'App\\Models\\User'
		LEFT JOIN roles ro ON ro.id = ar.role_id
		INNER JOIN job_rates jr ON jr.user_id = u.id AND jr.deleted_at IS NULL
		WHERE u.is_active = 1 AND u.deleted_at IS NULL
	`

	var args []interface{}
	var conditions []string

	if filters.UserID != nil {
		conditions = append(conditions, "u.id = ?")
		args = append(args, *filters.UserID)
	}

	if filters.Search != nil && *filters.Search != "" {
		conditions = append(conditions, "(u.name LIKE ? OR u.email LIKE ?)")
		searchTerm := "%" + *filters.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Contar total
	countQuery := "SELECT COUNT(DISTINCT u.id) FROM users u " +
		"LEFT JOIN assigned_roles ar ON ar.entity_id = u.id AND ar.entity_type = 'App\\\\Models\\\\User' " +
		"LEFT JOIN roles ro ON ro.id = ar.role_id " +
		"INNER JOIN job_rates jr ON jr.user_id = u.id AND jr.deleted_at IS NULL " +
		"WHERE u.is_active = 1 AND u.deleted_at IS NULL"

	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting payroll users: %w", err)
	}

	// Agregar paginación
	offset := (filters.Page - 1) * filters.PageSize
	baseQuery += " ORDER BY u.name ASC LIMIT ? OFFSET ?"
	args = append(args, filters.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying payroll users: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var users []*payroll.PayrollUser
	for rows.Next() {
		user := &payroll.PayrollUser{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.RoleName)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	// Para cada usuario, obtener sus rates agrupados por status
	for _, user := range users {
		if err := r.loadUserRates(ctx, user); err != nil {
			return nil, 0, fmt.Errorf("error loading rates for user %d: %w", user.ID, err)
		}
	}

	return users, total, nil
}

func (r *repository) loadUserRates(ctx context.Context, user *payroll.PayrollUser) error {
	query := `
		SELECT jr.id, jr.job_id, jr.user_id, jr.job_rate_status_id, jr.sale_price,
			   jr.rate_percent, jr.rate_flat, jr.tech_parts, jr.company_parts,
			   jr.parts_replaced, jr.deduction, jr.payment, jr.paid, jr.notes,
			   jr.created_at, jr.updated_at,
			   j.work_order,
			   CONCAT(p.street, ', ', p.city, ', ', p.state, ' ', p.zip) as property_address,
			   jrs.label as status_label, jrs.class as status_class
		FROM job_rates jr
		LEFT JOIN jobs j ON j.id = jr.job_id
		LEFT JOIN properties p ON p.id = j.property_id
		LEFT JOIN job_rate_statuses jrs ON jrs.id = jr.job_rate_status_id
		WHERE jr.user_id = ? AND jr.deleted_at IS NULL
		ORDER BY jr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	user.UnpaidRates = []*payroll.PayrollRate{}
	user.HoldingRates = []*payroll.PayrollRate{}
	user.PaidRates = []*payroll.PayrollRate{}

	for rows.Next() {
		rate := &payroll.PayrollRate{}
		err := rows.Scan(
			&rate.ID, &rate.JobID, &rate.UserID, &rate.JobRateStatusID,
			&rate.SalePrice, &rate.RatePercent, &rate.RateFlat,
			&rate.TechParts, &rate.CompanyParts, &rate.PartsReplaced,
			&rate.Deduction, &rate.Payment, &rate.Paid, &rate.Notes,
			&rate.CreatedAt, &rate.UpdatedAt,
			&rate.WorkOrder, &rate.PropertyAddr,
			&rate.StatusLabel, &rate.StatusClass,
		)
		if err != nil {
			return err
		}

		// Agrupar por status
		if rate.StatusLabel != nil {
			switch *rate.StatusLabel {
			case "Unpaid":
				user.UnpaidRates = append(user.UnpaidRates, rate)
				user.TotalUnpaid += rate.Payment
			case "Holding":
				user.HoldingRates = append(user.HoldingRates, rate)
				user.TotalHolding += rate.Payment
			case "Paid":
				user.PaidRates = append(user.PaidRates, rate)
				user.TotalPaid += rate.Payment
			}
		}
	}

	return rows.Err()
}
