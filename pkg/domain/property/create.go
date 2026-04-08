package property

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Create(ctx context.Context, property *Property) error {
	// Validar que el customer existe y no está eliminado
	customer, err := uc.customerRepo.GetByID(ctx, property.CustomerID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to validate customer",
			slog.String("error", err.Error()),
			slog.Int64("customer_id", property.CustomerID))
		return errors.New("invalid customer_id")
	}

	if customer.DeletedAt != nil {
		slog.WarnContext(ctx, "Customer is deleted",
			slog.Int64("customer_id", property.CustomerID))
		return errors.New("customer is deleted")
	}

	if err := uc.repo.Create(ctx, property); err != nil {
		slog.ErrorContext(ctx, "Failed to create property",
			slog.String("error", err.Error()),
			slog.String("street", property.Street))
		return err
	}

	slog.InfoContext(ctx, "Property created successfully",
		slog.Int64("property_id", property.ID),
		slog.String("street", property.Street))

	return nil
}
