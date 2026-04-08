package property

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Update(ctx context.Context, property *Property) error {
	// Validar que la propiedad existe
	existing, err := uc.repo.GetByID(ctx, property.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get property for update",
			slog.String("error", err.Error()),
			slog.Int64("property_id", property.ID))
		return err
	}

	// Validar que no esté eliminada
	if existing.DeletedAt != nil {
		slog.WarnContext(ctx, "Cannot update deleted property",
			slog.Int64("property_id", property.ID))
		return errors.New("cannot update deleted property")
	}

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

	if err := uc.repo.Update(ctx, property); err != nil {
		slog.ErrorContext(ctx, "Failed to update property",
			slog.String("error", err.Error()),
			slog.Int64("property_id", property.ID))
		return err
	}

	slog.InfoContext(ctx, "Property updated successfully",
		slog.Int64("property_id", property.ID),
		slog.String("street", property.Street))

	return nil
}
