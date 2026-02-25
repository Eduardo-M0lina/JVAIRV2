package property_equipment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/property"
)

func strPtr(s string) *string {
	return &s
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		eq := &PropertyEquipment{
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
		}

		mockProp := &property.Property{
			ID:        1,
			Street:    "123 Main St",
			CreatedAt: &now,
		}

		mockPropertyRepo.On("GetByID", ctx, int64(1)).Return(mockProp, nil)
		mockRepo.On("Create", ctx, eq).Return(nil)

		err := uc.Create(ctx, eq)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("invalid property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		eq := &PropertyEquipment{
			PropertyID: 999,
			Area:       strPtr("Main Floor"),
		}

		mockPropertyRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property not found"))

		err := uc.Create(ctx, eq)

		assert.Error(t, err)
		assert.Equal(t, "invalid property_id", err.Error())
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("deleted property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		eq := &PropertyEquipment{
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
		}

		mockProp := &property.Property{
			ID:        1,
			Street:    "123 Main St",
			DeletedAt: &now,
		}

		mockPropertyRepo.On("GetByID", ctx, int64(1)).Return(mockProp, nil)

		err := uc.Create(ctx, eq)

		assert.Error(t, err)
		assert.Equal(t, "property is deleted", err.Error())
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		eq := &PropertyEquipment{
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
		}

		mockProp := &property.Property{
			ID:        1,
			Street:    "123 Main St",
			CreatedAt: &now,
		}

		mockPropertyRepo.On("GetByID", ctx, int64(1)).Return(mockProp, nil)
		mockRepo.On("Create", ctx, eq).Return(errors.New("database error"))

		err := uc.Create(ctx, eq)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
		mockPropertyRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		expected := &PropertyEquipment{
			ID:         1,
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
			CreatedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expected, nil)

		eq, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, eq)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property equipment not found"))

		eq, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, eq)
		mockRepo.AssertExpectations(t)
	})
}

func TestList(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		mockProp := &property.Property{
			ID:        1,
			Street:    "123 Main St",
			CreatedAt: &now,
		}

		expected := []*PropertyEquipment{
			{
				ID:         1,
				PropertyID: 1,
				Area:       strPtr("Main Floor"),
				CreatedAt:  &now,
			},
		}

		mockPropertyRepo.On("GetByID", ctx, int64(1)).Return(mockProp, nil)
		mockRepo.On("List", ctx, int64(1)).Return(expected, nil)

		equipment, err := uc.List(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, equipment)
		mockRepo.AssertExpectations(t)
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("invalid property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		mockPropertyRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property not found"))

		equipment, err := uc.List(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, equipment)
		assert.Equal(t, "invalid property_id", err.Error())
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("deleted property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		mockProp := &property.Property{
			ID:        1,
			Street:    "123 Main St",
			DeletedAt: &now,
		}

		mockPropertyRepo.On("GetByID", ctx, int64(1)).Return(mockProp, nil)

		equipment, err := uc.List(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, equipment)
		assert.Equal(t, "property is deleted", err.Error())
		mockPropertyRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		existing := &PropertyEquipment{
			ID:         1,
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
			CreatedAt:  &now,
		}

		updated := &PropertyEquipment{
			ID:           1,
			PropertyID:   1,
			Area:         strPtr("Second Floor"),
			OutdoorBrand: strPtr("Carrier"),
		}

		mockProp := &property.Property{
			ID:        1,
			Street:    "123 Main St",
			CreatedAt: &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockPropertyRepo.On("GetByID", ctx, int64(1)).Return(mockProp, nil)
		mockRepo.On("Update", ctx, updated).Return(nil)

		err := uc.Update(ctx, updated)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockPropertyRepo.AssertExpectations(t)
	})

	t.Run("equipment not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		updated := &PropertyEquipment{
			ID:         999,
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
		}

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property equipment not found"))

		err := uc.Update(ctx, updated)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("equipment does not belong to property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		existing := &PropertyEquipment{
			ID:         1,
			PropertyID: 2,
			Area:       strPtr("Main Floor"),
			CreatedAt:  &now,
		}

		updated := &PropertyEquipment{
			ID:         1,
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
		}

		mockProp := &property.Property{
			ID:        1,
			Street:    "123 Main St",
			CreatedAt: &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockPropertyRepo.On("GetByID", ctx, int64(1)).Return(mockProp, nil)

		err := uc.Update(ctx, updated)

		assert.Error(t, err)
		assert.Equal(t, "equipment does not belong to this property", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		existing := &PropertyEquipment{
			ID:         1,
			PropertyID: 999,
			Area:       strPtr("Main Floor"),
			CreatedAt:  &now,
		}

		updated := &PropertyEquipment{
			ID:         1,
			PropertyID: 999,
			Area:       strPtr("Main Floor"),
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockPropertyRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property not found"))

		err := uc.Update(ctx, updated)

		assert.Error(t, err)
		assert.Equal(t, "invalid property_id", err.Error())
		mockRepo.AssertExpectations(t)
		mockPropertyRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		existing := &PropertyEquipment{
			ID:         1,
			PropertyID: 1,
			Area:       strPtr("Main Floor"),
			CreatedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockRepo.On("Delete", ctx, int64(1)).Return(nil)

		err := uc.Delete(ctx, 1, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("equipment not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property equipment not found"))

		err := uc.Delete(ctx, 999, 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("equipment does not belong to property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockPropertyRepo := new(property.MockRepository)
		uc := NewUseCase(mockRepo, mockPropertyRepo)

		existing := &PropertyEquipment{
			ID:         1,
			PropertyID: 2,
			Area:       strPtr("Main Floor"),
			CreatedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)

		err := uc.Delete(ctx, 1, 1)

		assert.Error(t, err)
		assert.Equal(t, "equipment does not belong to this property", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
