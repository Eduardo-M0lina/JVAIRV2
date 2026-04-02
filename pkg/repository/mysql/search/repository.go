package search

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/search"
)

// Repository implementa search.Repository para MySQL
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia del repositorio
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// SearchJobs busca jobs por work_order, property, customer
func (r *Repository) SearchJobs(ctx context.Context, filters search.SearchFilters) ([]search.JobSearchResult, error) {
	query := `
		SELECT
			j.id,
			j.work_order,
			p.street AS property_street,
			p.city AS property_city,
			c.name AS customer_name,
			js.label AS status_name,
			u.name AS technician_name,
			j.closed
		FROM jobs j
		LEFT JOIN properties p ON j.property_id = p.id
		LEFT JOIN customers c ON p.customer_id = c.id
		LEFT JOIN job_statuses js ON j.job_status_id = js.id
		LEFT JOIN users u ON j.user_id = u.id
		WHERE j.deleted_at IS NULL
		AND (
			j.work_order LIKE ?
			OR p.property_code LIKE ?
			OR p.street LIKE ?
			OR p.city LIKE ?
			OR p.state LIKE ?
			OR p.zip LIKE ?
			OR c.name LIKE ?
			OR c.email LIKE ?
			OR c.phone LIKE ?
			OR c.mobile LIKE ?
		)
	`
	args := make([]interface{}, 10)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	// Si el usuario tiene restricción job_view_user_only
	if filters.UserOnly && filters.UserID != nil {
		query += " AND j.user_id = ?"
		args = append(args, *filters.UserID)
	}

	query += " ORDER BY j.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.JobSearchResult
	for rows.Next() {
		var job search.JobSearchResult
		err := rows.Scan(
			&job.ID,
			&job.WorkOrder,
			&job.PropertyStreet,
			&job.PropertyCity,
			&job.CustomerName,
			&job.StatusName,
			&job.TechnicianName,
			&job.Closed,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, job)
	}

	if results == nil {
		results = []search.JobSearchResult{}
	}
	return results, rows.Err()
}

// SearchCustomers busca customers por name, email, phone, etc.
func (r *Repository) SearchCustomers(ctx context.Context, filters search.SearchFilters) ([]search.CustomerSearchResult, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.email,
			c.phone,
			c.mobile,
			c.website
		FROM customers c
		WHERE c.deleted_at IS NULL
		AND (
			c.name LIKE ?
			OR c.email LIKE ?
			OR c.phone LIKE ?
			OR c.mobile LIKE ?
			OR c.fax LIKE ?
			OR c.phone_other LIKE ?
			OR c.website LIKE ?
			OR c.contact_name LIKE ?
			OR c.contact_email LIKE ?
			OR c.contact_phone LIKE ?
			OR c.billing_address_street LIKE ?
			OR c.billing_address_city LIKE ?
		)
	`
	args := make([]interface{}, 12)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	// Si el usuario tiene restricción job_view_user_only, filtrar por customers con jobs del usuario
	if filters.UserOnly && filters.UserID != nil {
		query += ` AND c.id IN (
			SELECT DISTINCT p.customer_id
			FROM properties p
			INNER JOIN jobs j ON j.property_id = p.id
			WHERE j.user_id = ? AND j.deleted_at IS NULL
		)`
		args = append(args, *filters.UserID)
	}

	query += " ORDER BY c.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.CustomerSearchResult
	for rows.Next() {
		var customer search.CustomerSearchResult
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Mobile,
			&customer.Website,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, customer)
	}

	if results == nil {
		results = []search.CustomerSearchResult{}
	}
	return results, rows.Err()
}

// SearchProperties busca properties por property_code, street, city, zip
func (r *Repository) SearchProperties(ctx context.Context, filters search.SearchFilters) ([]search.PropertySearchResult, error) {
	query := `
		SELECT
			p.id,
			p.property_code,
			p.street,
			p.city,
			p.state,
			p.zip,
			c.name AS customer_name
		FROM properties p
		LEFT JOIN customers c ON p.customer_id = c.id
		WHERE p.deleted_at IS NULL
		AND (
			p.property_code LIKE ?
			OR p.street LIKE ?
			OR p.city LIKE ?
			OR p.state LIKE ?
			OR p.zip LIKE ?
		)
	`
	args := make([]interface{}, 5)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	// Si el usuario tiene restricción job_view_user_only
	if filters.UserOnly && filters.UserID != nil {
		query += ` AND p.id IN (
			SELECT DISTINCT j.property_id
			FROM jobs j
			WHERE j.user_id = ? AND j.deleted_at IS NULL
		)`
		args = append(args, *filters.UserID)
	}

	query += " ORDER BY p.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.PropertySearchResult
	for rows.Next() {
		var prop search.PropertySearchResult
		err := rows.Scan(
			&prop.ID,
			&prop.PropertyCode,
			&prop.Street,
			&prop.City,
			&prop.State,
			&prop.Zip,
			&prop.CustomerName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, prop)
	}

	if results == nil {
		results = []search.PropertySearchResult{}
	}
	return results, rows.Err()
}

// SearchInvoices busca invoices por invoice_number, work_order, customer
func (r *Repository) SearchInvoices(ctx context.Context, filters search.SearchFilters) ([]search.InvoiceSearchResult, error) {
	query := `
		SELECT
			i.id,
			i.invoice_number,
			i.total,
			j.work_order,
			c.name AS customer_name
		FROM invoices i
		LEFT JOIN jobs j ON i.job_id = j.id
		LEFT JOIN properties p ON j.property_id = p.id
		LEFT JOIN customers c ON p.customer_id = c.id
		WHERE i.deleted_at IS NULL
		AND (
			i.invoice_number LIKE ?
			OR j.work_order LIKE ?
			OR c.name LIKE ?
			OR c.email LIKE ?
			OR c.phone LIKE ?
		)
	`
	args := make([]interface{}, 5)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	query += " ORDER BY i.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.InvoiceSearchResult
	for rows.Next() {
		var invoice search.InvoiceSearchResult
		err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNumber,
			&invoice.Total,
			&invoice.WorkOrder,
			&invoice.CustomerName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, invoice)
	}

	if results == nil {
		results = []search.InvoiceSearchResult{}
	}
	return results, rows.Err()
}

// SearchQuotes busca quotes por quote_number, work_order, customer
func (r *Repository) SearchQuotes(ctx context.Context, filters search.SearchFilters) ([]search.QuoteSearchResult, error) {
	query := `
		SELECT
			q.id,
			q.quote_number,
			q.amount,
			j.work_order,
			c.name AS customer_name,
			qs.label AS status_name
		FROM quotes q
		LEFT JOIN jobs j ON q.job_id = j.id
		LEFT JOIN properties p ON j.property_id = p.id
		LEFT JOIN customers c ON p.customer_id = c.id
		LEFT JOIN quote_statuses qs ON q.quote_status_id = qs.id
		WHERE q.deleted_at IS NULL
		AND (
			q.quote_number LIKE ?
			OR j.work_order LIKE ?
			OR c.name LIKE ?
			OR c.email LIKE ?
			OR c.phone LIKE ?
		)
	`
	args := make([]interface{}, 5)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	query += " ORDER BY q.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.QuoteSearchResult
	for rows.Next() {
		var quote search.QuoteSearchResult
		err := rows.Scan(
			&quote.ID,
			&quote.QuoteNumber,
			&quote.Total,
			&quote.WorkOrder,
			&quote.CustomerName,
			&quote.StatusName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, quote)
	}

	if results == nil {
		results = []search.QuoteSearchResult{}
	}
	return results, rows.Err()
}

// SearchWarranties busca warranties por warranty_number, agreement_number, work_order
func (r *Repository) SearchWarranties(ctx context.Context, filters search.SearchFilters) ([]search.WarrantySearchResult, error) {
	query := `
		SELECT
			w.id,
			w.warranty_number,
			w.agreement_number,
			j.work_order,
			c.name AS customer_name,
			wt.label AS type_name,
			ws.label AS status_name
		FROM warranties w
		LEFT JOIN jobs j ON w.job_id = j.id
		LEFT JOIN properties p ON j.property_id = p.id
		LEFT JOIN customers c ON p.customer_id = c.id
		LEFT JOIN warranty_types wt ON w.warranty_type_id = wt.id
		LEFT JOIN warranty_statuses ws ON w.warranty_status_id = ws.id
		WHERE w.deleted_at IS NULL
		AND (
			w.warranty_number LIKE ?
			OR w.agreement_number LIKE ?
			OR j.work_order LIKE ?
			OR c.name LIKE ?
			OR c.email LIKE ?
		)
	`
	args := make([]interface{}, 5)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	query += " ORDER BY w.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.WarrantySearchResult
	for rows.Next() {
		var warranty search.WarrantySearchResult
		err := rows.Scan(
			&warranty.ID,
			&warranty.WarrantyNumber,
			&warranty.AgreementNumber,
			&warranty.WorkOrder,
			&warranty.CustomerName,
			&warranty.TypeName,
			&warranty.StatusName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, warranty)
	}

	if results == nil {
		results = []search.WarrantySearchResult{}
	}
	return results, rows.Err()
}

// SearchWarrantyClaims busca warranty claims por claim_number, internal_claim_number, etc.
func (r *Repository) SearchWarrantyClaims(ctx context.Context, filters search.SearchFilters) ([]search.WarrantyClaimSearchResult, error) {
	query := `
		SELECT
			wc.id,
			wc.claim_number,
			wc.internal_claim_number,
			wc.warranty_part,
			wc.manufacturer,
			j.work_order,
			wcs.label AS status_name
		FROM warranty_claims wc
		LEFT JOIN jobs j ON wc.job_id = j.id
		LEFT JOIN warranty_claim_statuses wcs ON wc.warranty_claim_status_id = wcs.id
		WHERE wc.deleted_at IS NULL
		AND (
			wc.claim_number LIKE ?
			OR wc.internal_claim_number LIKE ?
			OR wc.warranty_part LIKE ?
			OR wc.manufacturer LIKE ?
			OR j.work_order LIKE ?
		)
	`
	args := make([]interface{}, 5)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	query += " ORDER BY wc.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.WarrantyClaimSearchResult
	for rows.Next() {
		var claim search.WarrantyClaimSearchResult
		err := rows.Scan(
			&claim.ID,
			&claim.ClaimNumber,
			&claim.InternalClaimNumber,
			&claim.WarrantyPart,
			&claim.Manufacturer,
			&claim.WorkOrder,
			&claim.StatusName,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, claim)
	}

	if results == nil {
		results = []search.WarrantyClaimSearchResult{}
	}
	return results, rows.Err()
}

// SearchUsers busca users por name, email
func (r *Repository) SearchUsers(ctx context.Context, filters search.SearchFilters) ([]search.UserSearchResult, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			r.title AS role_name,
			u.is_active
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.deleted_at IS NULL
		AND (
			u.name LIKE ?
			OR u.email LIKE ?
		)
	`
	args := make([]interface{}, 2)
	searchPattern := "%" + filters.Query + "%"
	for i := range args {
		args[i] = searchPattern
	}

	query += " ORDER BY u.created_at DESC LIMIT ?"
	args = append(args, filters.Limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []search.UserSearchResult
	for rows.Next() {
		var user search.UserSearchResult
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.RoleName,
			&user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, user)
	}

	if results == nil {
		results = []search.UserSearchResult{}
	}
	return results, rows.Err()
}
