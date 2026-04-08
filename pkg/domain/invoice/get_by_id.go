package invoice

import (
	"context"
	"log/slog"
)

// GetByID obtiene una factura por su ID
func (uc *UseCase) GetByID(ctx context.Context, id int64) (*Invoice, error) {
	inv, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get invoice",
			slog.Int64("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}

	if inv.IsDeleted() {
		return nil, ErrInvoiceNotFound
	}

	return inv, nil
}
