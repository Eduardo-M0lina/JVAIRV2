package job_equipment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		eq := &JobEquipment{
			JobID: 1,
			Type:  "current",
			Area:  strPtr("Main Floor"),
		}

		mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
		mockRepo.On("Create", ctx, eq).Return(nil)

		err := uc.Create(ctx, eq)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockJobChecker.AssertExpectations(t)
	})

	t.Run("invalid job", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		eq := &JobEquipment{
			JobID: 999,
			Type:  "current",
			Area:  strPtr("Main Floor"),
		}

		mockJobChecker.On("JobExists", int64(999)).Return(false, nil)

		err := uc.Create(ctx, eq)

		assert.Error(t, err)
		assert.Equal(t, "invalid job_id", err.Error())
		mockJobChecker.AssertExpectations(t)
	})

	t.Run("job checker error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		eq := &JobEquipment{
			JobID: 1,
			Type:  "current",
			Area:  strPtr("Main Floor"),
		}

		mockJobChecker.On("JobExists", int64(1)).Return(false, errors.New("database error"))

		err := uc.Create(ctx, eq)

		assert.Error(t, err)
		assert.Equal(t, "invalid job_id", err.Error())
		mockJobChecker.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		eq := &JobEquipment{
			JobID: 1,
			Type:  "current",
			Area:  strPtr("Main Floor"),
		}

		mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
		mockRepo.On("Create", ctx, eq).Return(errors.New("database error"))

		err := uc.Create(ctx, eq)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
		mockJobChecker.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		expected := &JobEquipment{
			ID:        1,
			JobID:     1,
			Type:      "current",
			Area:      strPtr("Main Floor"),
			CreatedAt: &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expected, nil)

		eq, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, eq)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("job equipment not found"))

		eq, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, eq)
		mockRepo.AssertExpectations(t)
	})
}

func TestList(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success without type filter", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		expected := []*JobEquipment{
			{
				ID:        1,
				JobID:     1,
				Type:      "current",
				Area:      strPtr("Main Floor"),
				CreatedAt: &now,
			},
		}

		mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
		mockRepo.On("List", ctx, int64(1), "").Return(expected, nil)

		equipment, err := uc.List(ctx, 1, "")

		assert.NoError(t, err)
		assert.Equal(t, expected, equipment)
		mockRepo.AssertExpectations(t)
		mockJobChecker.AssertExpectations(t)
	})

	t.Run("success with type filter", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		expected := []*JobEquipment{
			{
				ID:        1,
				JobID:     1,
				Type:      "current",
				Area:      strPtr("Main Floor"),
				CreatedAt: &now,
			},
		}

		mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
		mockRepo.On("List", ctx, int64(1), "current").Return(expected, nil)

		equipment, err := uc.List(ctx, 1, "current")

		assert.NoError(t, err)
		assert.Equal(t, expected, equipment)
		mockRepo.AssertExpectations(t)
		mockJobChecker.AssertExpectations(t)
	})

	t.Run("invalid job", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		mockJobChecker.On("JobExists", int64(999)).Return(false, nil)

		equipment, err := uc.List(ctx, 999, "")

		assert.Error(t, err)
		assert.Nil(t, equipment)
		assert.Equal(t, "invalid job_id", err.Error())
		mockJobChecker.AssertExpectations(t)
	})

	t.Run("invalid type filter", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		mockJobChecker.On("JobExists", int64(1)).Return(true, nil)

		equipment, err := uc.List(ctx, 1, "invalid")

		assert.Error(t, err)
		assert.Nil(t, equipment)
		assert.Equal(t, "type must be one of: current, new", err.Error())
		mockJobChecker.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		existing := &JobEquipment{
			ID:        1,
			JobID:     1,
			Type:      "current",
			Area:      strPtr("Main Floor"),
			CreatedAt: &now,
		}

		updated := &JobEquipment{
			ID:           1,
			JobID:        1,
			Type:         "new",
			Area:         strPtr("Second Floor"),
			OutdoorBrand: strPtr("Carrier"),
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockJobChecker.On("JobExists", int64(1)).Return(true, nil)
		mockRepo.On("Update", ctx, updated).Return(nil)

		err := uc.Update(ctx, updated)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockJobChecker.AssertExpectations(t)
	})

	t.Run("equipment not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		updated := &JobEquipment{
			ID:    999,
			JobID: 1,
			Type:  "current",
			Area:  strPtr("Main Floor"),
		}

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("job equipment not found"))

		err := uc.Update(ctx, updated)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("equipment does not belong to job", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		existing := &JobEquipment{
			ID:        1,
			JobID:     2,
			Type:      "current",
			Area:      strPtr("Main Floor"),
			CreatedAt: &now,
		}

		updated := &JobEquipment{
			ID:    1,
			JobID: 1,
			Type:  "current",
			Area:  strPtr("Main Floor"),
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)

		err := uc.Update(ctx, updated)

		assert.Error(t, err)
		assert.Equal(t, "equipment does not belong to this job", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid job", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		existing := &JobEquipment{
			ID:        1,
			JobID:     1,
			Type:      "current",
			Area:      strPtr("Main Floor"),
			CreatedAt: &now,
		}

		updated := &JobEquipment{
			ID:    1,
			JobID: 1,
			Type:  "current",
			Area:  strPtr("Main Floor"),
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockJobChecker.On("JobExists", int64(1)).Return(false, nil)

		err := uc.Update(ctx, updated)

		assert.Error(t, err)
		assert.Equal(t, "invalid job_id", err.Error())
		mockRepo.AssertExpectations(t)
		mockJobChecker.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		existing := &JobEquipment{
			ID:        1,
			JobID:     1,
			Type:      "current",
			Area:      strPtr("Main Floor"),
			CreatedAt: &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		mockRepo.On("Delete", ctx, int64(1)).Return(nil)

		err := uc.Delete(ctx, 1, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("equipment not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("job equipment not found"))

		err := uc.Delete(ctx, 999, 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("equipment does not belong to job", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJobChecker := new(MockJobChecker)
		uc := NewUseCase(mockRepo, mockJobChecker)

		existing := &JobEquipment{
			ID:        1,
			JobID:     2,
			Type:      "current",
			Area:      strPtr("Main Floor"),
			CreatedAt: &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)

		err := uc.Delete(ctx, 1, 1)

		assert.Error(t, err)
		assert.Equal(t, "equipment does not belong to this job", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
