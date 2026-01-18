package customer

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get customer for deletion",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		return err
	}

	if existing.DeletedAt != nil {
		slog.WarnContext(ctx, "Customer already deleted",
			slog.Int64("customer_id", id))
		return errors.New("customer already deleted")
	}

	hasProperties, err := uc.repo.HasProperties(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check customer properties",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		return err
	}

	if hasProperties {
		slog.WarnContext(ctx, "Cannot delete customer with properties",
			slog.Int64("customer_id", id))
		return errors.New("cannot delete customer with associated properties")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", id))
		return err
	}

	slog.InfoContext(ctx, "Customer deleted successfully",
		slog.Int64("customer_id", id))

	return nil
}
