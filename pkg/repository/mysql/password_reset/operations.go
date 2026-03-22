package password_reset

import (
	"context"
	"database/sql"

	domainPasswordReset "github.com/your-org/jvairv2/pkg/domain/password_reset"
)

// Create crea un nuevo token de reseteo
func (r *Repository) Create(ctx context.Context, pr *domainPasswordReset.PasswordReset) error {
	query := `
		INSERT INTO password_resets (email, token, created_at)
		VALUES (?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, pr.Email, pr.Token, pr.CreatedAt)
	return err
}

// GetByToken obtiene un token de reseteo por su valor
func (r *Repository) GetByToken(ctx context.Context, token string) (*domainPasswordReset.PasswordReset, error) {
	query := `
		SELECT email, token, created_at
		FROM password_resets
		WHERE token = ?
		LIMIT 1
	`

	var pr domainPasswordReset.PasswordReset
	var createdAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&pr.Email,
		&pr.Token,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if createdAt.Valid {
		pr.CreatedAt = &createdAt.Time
	}

	return &pr, nil
}

// GetByEmail obtiene el token más reciente por email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*domainPasswordReset.PasswordReset, error) {
	query := `
		SELECT email, token, created_at
		FROM password_resets
		WHERE email = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	var pr domainPasswordReset.PasswordReset
	var createdAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&pr.Email,
		&pr.Token,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if createdAt.Valid {
		pr.CreatedAt = &createdAt.Time
	}

	return &pr, nil
}

// Delete elimina un token de reseteo
func (r *Repository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM password_resets WHERE token = ?`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

// DeleteByEmail elimina todos los tokens de un email
func (r *Repository) DeleteByEmail(ctx context.Context, email string) error {
	query := `DELETE FROM password_resets WHERE email = ?`
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

// DeleteExpired elimina tokens expirados (más antiguos que X horas)
func (r *Repository) DeleteExpired(ctx context.Context, hours int) error {
	query := `DELETE FROM password_resets WHERE created_at < DATE_SUB(NOW(), INTERVAL ? HOUR)`
	_, err := r.db.ExecContext(ctx, query, hours)
	return err
}
