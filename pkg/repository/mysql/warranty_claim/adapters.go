package warranty_claim

import (
	"context"
	"database/sql"

	domainWC "github.com/angumol/jvairv2/pkg/domain/warranty_claim"
)

// JobCheckerAdapter adapta la verificación de jobs para el use case de warranty claims
type JobCheckerAdapter struct {
	db *sql.DB
}

func NewJobCheckerAdapter(db *sql.DB) domainWC.JobChecker {
	return &JobCheckerAdapter{db: db}
}

func (a *JobCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM jobs WHERE id = ? AND deleted_at IS NULL)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainWC.ErrInvalidJob
	}
	return true, nil
}

// ClaimTypeCheckerAdapter adapta la verificación de tipos de reclamo
type ClaimTypeCheckerAdapter struct {
	db *sql.DB
}

func NewClaimTypeCheckerAdapter(db *sql.DB) domainWC.ClaimTypeChecker {
	return &ClaimTypeCheckerAdapter{db: db}
}

func (a *ClaimTypeCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM warranty_claim_types WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainWC.ErrInvalidClaimType
	}
	return true, nil
}

// ClaimStatusCheckerAdapter adapta la verificación de estados de reclamo
type ClaimStatusCheckerAdapter struct {
	db *sql.DB
}

func NewClaimStatusCheckerAdapter(db *sql.DB) domainWC.ClaimStatusChecker {
	return &ClaimStatusCheckerAdapter{db: db}
}

func (a *ClaimStatusCheckerAdapter) GetByID(ctx context.Context, id int64) (interface{}, error) {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM warranty_claim_statuses WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		return nil, domainWC.ErrInvalidClaimStatus
	}
	return true, nil
}
