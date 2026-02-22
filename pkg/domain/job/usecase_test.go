package job

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestUseCase() (*UseCase, *MockRepository, *MockJobCategoryChecker, *MockJobPriorityChecker, *MockJobStatusChecker, *MockWorkflowChecker, *MockPropertyChecker, *MockUserChecker, *MockTechnicianJobStatusChecker) {
	repo := new(MockRepository)
	catChecker := new(MockJobCategoryChecker)
	prioChecker := new(MockJobPriorityChecker)
	statusChecker := new(MockJobStatusChecker)
	wfChecker := new(MockWorkflowChecker)
	propChecker := new(MockPropertyChecker)
	userChecker := new(MockUserChecker)
	techChecker := new(MockTechnicianJobStatusChecker)

	uc := NewUseCase(repo, catChecker, prioChecker, statusChecker, wfChecker, propChecker, userChecker, techChecker)
	return uc, repo, catChecker, prioChecker, statusChecker, wfChecker, propChecker, userChecker, techChecker
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		uc, repo, catChecker, prioChecker, _, wfChecker, propChecker, _, _ := newTestUseCase()

		j := &Job{
			JobCategoryID: 1,
			JobPriorityID: 2,
			PropertyID:    3,
		}

		catChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		prioChecker.On("GetByID", ctx, int64(2)).Return(true, nil)
		propChecker.On("GetByID", ctx, int64(3)).Return(true, nil)
		propChecker.On("GetWorkflowID", ctx, int64(3)).Return(int64(10), nil)
		wfChecker.On("GetInitialStatusID", ctx, int64(10)).Return(int64(20), nil)
		repo.On("Create", ctx, j).Return(nil)

		err := uc.Create(ctx, j)

		assert.NoError(t, err)
		assert.Equal(t, int64(10), j.WorkflowID)
		assert.Equal(t, int64(20), j.JobStatusID)
		repo.AssertExpectations(t)
	})

	t.Run("missing category", func(t *testing.T) {
		uc, _, _, _, _, _, _, _, _ := newTestUseCase()

		j := &Job{
			JobPriorityID: 2,
			PropertyID:    3,
		}

		err := uc.Create(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, "job_category_id is required", err.Error())
	})

	t.Run("missing priority", func(t *testing.T) {
		uc, _, _, _, _, _, _, _, _ := newTestUseCase()

		j := &Job{
			JobCategoryID: 1,
			PropertyID:    3,
		}

		err := uc.Create(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, "job_priority_id is required", err.Error())
	})

	t.Run("missing property", func(t *testing.T) {
		uc, _, _, _, _, _, _, _, _ := newTestUseCase()

		j := &Job{
			JobCategoryID: 1,
			JobPriorityID: 2,
		}

		err := uc.Create(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, "property_id is required", err.Error())
	})

	t.Run("invalid category", func(t *testing.T) {
		uc, _, catChecker, _, _, _, _, _, _ := newTestUseCase()

		j := &Job{
			JobCategoryID: 999,
			JobPriorityID: 2,
			PropertyID:    3,
		}

		catChecker.On("GetByID", ctx, int64(999)).Return(nil, ErrInvalidJobCategory)

		err := uc.Create(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidJobCategory, err)
	})

	t.Run("invalid property", func(t *testing.T) {
		uc, _, catChecker, prioChecker, _, _, propChecker, _, _ := newTestUseCase()

		j := &Job{
			JobCategoryID: 1,
			JobPriorityID: 2,
			PropertyID:    999,
		}

		catChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		prioChecker.On("GetByID", ctx, int64(2)).Return(true, nil)
		propChecker.On("GetByID", ctx, int64(999)).Return(nil, ErrInvalidProperty)

		err := uc.Create(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidProperty, err)
	})

	t.Run("with user validation", func(t *testing.T) {
		uc, repo, catChecker, prioChecker, _, wfChecker, propChecker, userChecker, _ := newTestUseCase()

		userID := int64(5)
		j := &Job{
			JobCategoryID: 1,
			JobPriorityID: 2,
			PropertyID:    3,
			UserID:        &userID,
		}

		catChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		prioChecker.On("GetByID", ctx, int64(2)).Return(true, nil)
		propChecker.On("GetByID", ctx, int64(3)).Return(true, nil)
		propChecker.On("GetWorkflowID", ctx, int64(3)).Return(int64(10), nil)
		wfChecker.On("GetInitialStatusID", ctx, int64(10)).Return(int64(20), nil)
		userChecker.On("GetByID", ctx, int64(5)).Return(true, nil)
		repo.On("Create", ctx, j).Return(nil)

		err := uc.Create(ctx, j)

		assert.NoError(t, err)
		userChecker.AssertExpectations(t)
	})

	t.Run("invalid user", func(t *testing.T) {
		uc, _, catChecker, prioChecker, _, wfChecker, propChecker, userChecker, _ := newTestUseCase()

		userID := int64(999)
		j := &Job{
			JobCategoryID: 1,
			JobPriorityID: 2,
			PropertyID:    3,
			UserID:        &userID,
		}

		catChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		prioChecker.On("GetByID", ctx, int64(2)).Return(true, nil)
		propChecker.On("GetByID", ctx, int64(3)).Return(true, nil)
		propChecker.On("GetWorkflowID", ctx, int64(3)).Return(int64(10), nil)
		wfChecker.On("GetInitialStatusID", ctx, int64(10)).Return(int64(20), nil)
		userChecker.On("GetByID", ctx, int64(999)).Return(nil, ErrInvalidUser)

		err := uc.Create(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidUser, err)
	})
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		expected := &Job{
			ID:            1,
			JobCategoryID: 1,
			JobPriorityID: 2,
			JobStatusID:   3,
			WorkflowID:    4,
			PropertyID:    5,
			DateReceived:  now,
			CreatedAt:     &now,
		}

		repo.On("GetByID", ctx, int64(1)).Return(expected, nil)

		j, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, j)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		repo.On("GetByID", ctx, int64(999)).Return(nil, ErrJobNotFound)

		j, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, j)
		assert.Equal(t, ErrJobNotFound, err)
	})

	t.Run("deleted job returns not found", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		deleted := &Job{
			ID:           1,
			DateReceived: now,
			DeletedAt:    &now,
		}

		repo.On("GetByID", ctx, int64(1)).Return(deleted, nil)

		j, err := uc.GetByID(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, j)
		assert.Equal(t, ErrJobNotFound, err)
	})
}

func TestList(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		expected := []*Job{
			{ID: 1, DateReceived: now, CreatedAt: &now},
			{ID: 2, DateReceived: now, CreatedAt: &now},
		}

		filters := map[string]interface{}{}
		repo.On("List", ctx, filters, 1, 10).Return(expected, 2, nil)

		jobs, total, err := uc.List(ctx, filters, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expected, jobs)
		assert.Equal(t, 2, total)
		repo.AssertExpectations(t)
	})

	t.Run("defaults page and pageSize", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		filters := map[string]interface{}{}
		repo.On("List", ctx, filters, 1, 10).Return([]*Job{}, 0, nil)

		jobs, total, err := uc.List(ctx, filters, 0, 0)

		assert.NoError(t, err)
		assert.Empty(t, jobs)
		assert.Equal(t, 0, total)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		filters := map[string]interface{}{}
		repo.On("List", ctx, filters, 1, 10).Return(nil, 0, errors.New("db error"))

		jobs, total, err := uc.List(ctx, filters, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, jobs)
		assert.Equal(t, 0, total)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		existing := &Job{
			ID:            1,
			JobCategoryID: 1,
			JobPriorityID: 2,
			JobStatusID:   3,
			WorkflowID:    4,
			PropertyID:    5,
			DateReceived:  now,
			CreatedAt:     &now,
		}

		updated := &Job{
			ID:            1,
			JobCategoryID: 1,
			JobPriorityID: 2,
			JobStatusID:   3,
			WorkflowID:    4,
			PropertyID:    5,
			DateReceived:  now,
			QuickNotes:    strPtr("Updated notes"),
		}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		repo.On("Update", ctx, updated).Return(nil)

		err := uc.Update(ctx, updated)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		j := &Job{ID: 999}
		repo.On("GetByID", ctx, int64(999)).Return(nil, ErrJobNotFound)

		err := uc.Update(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, ErrJobNotFound, err)
	})

	t.Run("deleted job", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		existing := &Job{
			ID:           1,
			DateReceived: now,
			DeletedAt:    &now,
		}

		j := &Job{ID: 1}
		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)

		err := uc.Update(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, ErrJobNotFound, err)
	})

	t.Run("tech status auto-updates job status", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, techChecker := newTestUseCase()

		existing := &Job{
			ID:            1,
			JobCategoryID: 1,
			JobPriorityID: 2,
			JobStatusID:   3,
			WorkflowID:    4,
			PropertyID:    5,
			DateReceived:  now,
			CreatedAt:     &now,
		}

		techID := int64(10)
		linkedStatusID := int64(50)
		updated := &Job{
			ID:                    1,
			JobCategoryID:         1,
			JobPriorityID:         2,
			JobStatusID:           3,
			WorkflowID:            4,
			PropertyID:            5,
			DateReceived:          now,
			TechnicianJobStatusID: &techID,
		}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		techChecker.On("GetByID", ctx, int64(10)).Return(true, nil)
		techChecker.On("GetLinkedJobStatusID", ctx, int64(10)).Return(&linkedStatusID, nil)
		repo.On("Update", ctx, updated).Return(nil)

		err := uc.Update(ctx, updated)

		assert.NoError(t, err)
		assert.Equal(t, int64(50), updated.JobStatusID)
		repo.AssertExpectations(t)
		techChecker.AssertExpectations(t)
	})

	t.Run("missing id", func(t *testing.T) {
		uc, _, _, _, _, _, _, _, _ := newTestUseCase()

		j := &Job{}

		err := uc.Update(ctx, j)

		assert.Error(t, err)
		assert.Equal(t, "id is required", err.Error())
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		existing := &Job{
			ID:           1,
			DateReceived: now,
			CreatedAt:    &now,
		}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		repo.On("Delete", ctx, int64(1)).Return(nil)

		err := uc.Delete(ctx, 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		repo.On("GetByID", ctx, int64(999)).Return(nil, ErrJobNotFound)

		err := uc.Delete(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, ErrJobNotFound, err)
	})

	t.Run("already deleted", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		existing := &Job{
			ID:           1,
			DateReceived: now,
			DeletedAt:    &now,
		}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)

		err := uc.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, ErrJobNotFound, err)
	})
}

func TestClose(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		uc, repo, _, _, statusChecker, _, _, _, _ := newTestUseCase()

		existing := &Job{
			ID:           1,
			DateReceived: now,
			Closed:       false,
			CreatedAt:    &now,
		}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		statusChecker.On("GetByID", ctx, int64(5)).Return(true, nil)
		repo.On("Close", ctx, int64(1), int64(5)).Return(nil)

		err := uc.Close(ctx, 1, 5)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		statusChecker.AssertExpectations(t)
	})

	t.Run("already closed", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		existing := &Job{
			ID:           1,
			DateReceived: now,
			Closed:       true,
			CreatedAt:    &now,
		}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)

		err := uc.Close(ctx, 1, 5)

		assert.Error(t, err)
		assert.Equal(t, ErrJobAlreadyClosed, err)
	})

	t.Run("not found", func(t *testing.T) {
		uc, repo, _, _, _, _, _, _, _ := newTestUseCase()

		repo.On("GetByID", ctx, int64(999)).Return(nil, ErrJobNotFound)

		err := uc.Close(ctx, 999, 5)

		assert.Error(t, err)
		assert.Equal(t, ErrJobNotFound, err)
	})

	t.Run("invalid status", func(t *testing.T) {
		uc, repo, _, _, statusChecker, _, _, _, _ := newTestUseCase()

		existing := &Job{
			ID:           1,
			DateReceived: now,
			Closed:       false,
			CreatedAt:    &now,
		}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		statusChecker.On("GetByID", ctx, int64(999)).Return(nil, ErrInvalidJobStatus)

		err := uc.Close(ctx, 1, 999)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidJobStatus, err)
	})
}

func strPtr(s string) *string {
	return &s
}
