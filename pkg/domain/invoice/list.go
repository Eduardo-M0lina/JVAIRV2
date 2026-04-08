package invoice

import (
	"context"
	"log/slog"
)

// List obtiene una lista paginada de facturas con filtros opcionales
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*Invoice, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	invoices, total, err := uc.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list invoices",
			slog.String("error", err.Error()))
		return nil, 0, err
	}

	slog.InfoContext(ctx, "Invoices listed successfully",
		slog.Int("total", total),
		slog.Int("page", page),
		slog.Int("pageSize", pageSize))

	return invoices, total, nil
}
