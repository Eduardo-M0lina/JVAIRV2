package property

import (
	"context"
	"fmt"
	"log/slog"
)

func (uc *UseCase) SearchByAddress(ctx context.Context, address string) ([]*Property, error) {
	if len(address) < 3 {
		return nil, fmt.Errorf("address must be at least 3 characters long")
	}

	properties, err := uc.repo.SearchByAddress(ctx, address)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to search properties by address",
			slog.String("error", err.Error()),
			slog.String("address", address))
		return nil, err
	}

	return properties, nil
}
