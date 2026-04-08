package property

import (
	"context"
	"errors"
	"log/slog"
)

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	// Validar que la propiedad existe
	property, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get property for deletion",
			slog.String("error", err.Error()),
			slog.Int64("property_id", id))
		return err
	}

	// Validar que no esté ya eliminada
	if property.DeletedAt != nil {
		slog.WarnContext(ctx, "Property already deleted",
			slog.Int64("property_id", id))
		return errors.New("property already deleted")
	}

	// Validar que no tenga jobs asociados
	hasJobs, err := uc.repo.HasJobs(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to check property jobs",
			slog.String("error", err.Error()),
			slog.Int64("property_id", id))
		return err
	}

	if hasJobs {
		slog.WarnContext(ctx, "Cannot delete property with associated jobs",
			slog.Int64("property_id", id))
		return errors.New("cannot delete property with associated jobs")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "Failed to delete property",
			slog.String("error", err.Error()),
			slog.Int64("property_id", id))
		return err
	}

	slog.InfoContext(ctx, "Property deleted successfully",
		slog.Int64("property_id", id))

	return nil
}
