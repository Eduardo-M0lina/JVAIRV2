package settings

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockRepository es un mock de la interfaz Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Get(ctx context.Context) (*Settings, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Settings), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, settings *Settings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}
