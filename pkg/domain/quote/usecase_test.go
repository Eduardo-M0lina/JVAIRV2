package quote

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func newTestUseCase() (*UseCase, *MockRepository, *MockJobChecker, *MockQuoteStatusChecker) {
	repo := new(MockRepository)
	jobChecker := new(MockJobChecker)
	quoteStatusChecker := new(MockQuoteStatusChecker)
	uc := NewUseCase(repo, jobChecker, quoteStatusChecker)
	return uc, repo, jobChecker, quoteStatusChecker
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, repo, jobChecker, qsChecker := newTestUseCase()

		q := &Quote{
			JobID:         1,
			QuoteNumber:   "Q-001",
			QuoteStatusID: 1,
			Amount:        500.00,
		}

		jobChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		qsChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		repo.On("Create", ctx, q).Return(nil)

		err := uc.Create(ctx, q)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		jobChecker.AssertExpectations(t)
		qsChecker.AssertExpectations(t)
	})

	t.Run("validation error - missing quote_number", func(t *testing.T) {
		uc, _, _, _ := newTestUseCase()

		q := &Quote{
			JobID:         1,
			QuoteStatusID: 1,
			Amount:        500.00,
		}

		err := uc.Create(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, "quote_number is required", err.Error())
	})

	t.Run("validation error - missing job_id", func(t *testing.T) {
		uc, _, _, _ := newTestUseCase()

		q := &Quote{
			QuoteNumber:   "Q-001",
			QuoteStatusID: 1,
			Amount:        500.00,
		}

		err := uc.Create(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, "job_id is required", err.Error())
	})

	t.Run("validation error - missing quote_status_id", func(t *testing.T) {
		uc, _, _, _ := newTestUseCase()

		q := &Quote{
			JobID:       1,
			QuoteNumber: "Q-001",
			Amount:      500.00,
		}

		err := uc.Create(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, "quote_status_id is required", err.Error())
	})

	t.Run("invalid job", func(t *testing.T) {
		uc, _, jobChecker, _ := newTestUseCase()

		q := &Quote{
			JobID:         999,
			QuoteNumber:   "Q-001",
			QuoteStatusID: 1,
			Amount:        500.00,
		}

		jobChecker.On("GetByID", ctx, int64(999)).Return(nil, ErrInvalidJob)

		err := uc.Create(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidJob, err)
	})

	t.Run("invalid quote status", func(t *testing.T) {
		uc, _, jobChecker, qsChecker := newTestUseCase()

		q := &Quote{
			JobID:         1,
			QuoteNumber:   "Q-001",
			QuoteStatusID: 999,
			Amount:        500.00,
		}

		jobChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		qsChecker.On("GetByID", ctx, int64(999)).Return(nil, ErrInvalidQuoteStatus)

		err := uc.Create(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidQuoteStatus, err)
	})

	t.Run("repository error", func(t *testing.T) {
		uc, repo, jobChecker, qsChecker := newTestUseCase()

		q := &Quote{
			JobID:         1,
			QuoteNumber:   "Q-001",
			QuoteStatusID: 1,
			Amount:        500.00,
		}

		jobChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		qsChecker.On("GetByID", ctx, int64(1)).Return(true, nil)
		repo.On("Create", ctx, q).Return(errors.New("db error"))

		err := uc.Create(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		expected := &Quote{ID: 1, QuoteNumber: "Q-001"}
		repo.On("GetByID", ctx, int64(1)).Return(expected, nil)

		result, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		repo.On("GetByID", ctx, int64(999)).Return(nil, ErrQuoteNotFound)

		result, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrQuoteNotFound, err)
	})
}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		expected := []*Quote{
			{ID: 1, QuoteNumber: "Q-001"},
			{ID: 2, QuoteNumber: "Q-002"},
		}
		filters := map[string]interface{}{}
		repo.On("List", ctx, filters, 1, 15).Return(expected, int64(2), nil)

		result, total, err := uc.List(ctx, filters, 1, 15)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, int64(2), total)
	})

	t.Run("default pagination", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		filters := map[string]interface{}{}
		repo.On("List", ctx, filters, 1, 15).Return([]*Quote{}, int64(0), nil)

		result, total, err := uc.List(ctx, filters, 0, 0)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, int64(0), total)
	})

	t.Run("repository error", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		filters := map[string]interface{}{}
		repo.On("List", ctx, filters, 1, 15).Return(nil, int64(0), errors.New("db error"))

		result, total, err := uc.List(ctx, filters, 1, 15)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, repo, jobChecker, qsChecker := newTestUseCase()

		existing := &Quote{ID: 1, JobID: 1, QuoteNumber: "Q-001", QuoteStatusID: 1, Amount: 500.00}
		q := &Quote{ID: 1, JobID: 1, QuoteNumber: "Q-001-updated", QuoteStatusID: 1, Amount: 600.00}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		repo.On("Update", ctx, q).Return(nil)
		_ = jobChecker
		_ = qsChecker

		err := uc.Update(ctx, q)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("success with changed job", func(t *testing.T) {
		uc, repo, jobChecker, _ := newTestUseCase()

		existing := &Quote{ID: 1, JobID: 1, QuoteNumber: "Q-001", QuoteStatusID: 1, Amount: 500.00}
		q := &Quote{ID: 1, JobID: 2, QuoteNumber: "Q-001", QuoteStatusID: 1, Amount: 500.00}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		jobChecker.On("GetByID", ctx, int64(2)).Return(true, nil)
		repo.On("Update", ctx, q).Return(nil)

		err := uc.Update(ctx, q)

		assert.NoError(t, err)
		jobChecker.AssertExpectations(t)
	})

	t.Run("success with changed status", func(t *testing.T) {
		uc, repo, _, qsChecker := newTestUseCase()

		existing := &Quote{ID: 1, JobID: 1, QuoteNumber: "Q-001", QuoteStatusID: 1, Amount: 500.00}
		q := &Quote{ID: 1, JobID: 1, QuoteNumber: "Q-001", QuoteStatusID: 2, Amount: 500.00}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		qsChecker.On("GetByID", ctx, int64(2)).Return(true, nil)
		repo.On("Update", ctx, q).Return(nil)

		err := uc.Update(ctx, q)

		assert.NoError(t, err)
		qsChecker.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		q := &Quote{ID: 999, JobID: 1, QuoteNumber: "Q-001", QuoteStatusID: 1, Amount: 500.00}
		repo.On("GetByID", ctx, int64(999)).Return(nil, ErrQuoteNotFound)

		err := uc.Update(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, ErrQuoteNotFound, err)
	})

	t.Run("invalid job on change", func(t *testing.T) {
		uc, repo, jobChecker, _ := newTestUseCase()

		existing := &Quote{ID: 1, JobID: 1, QuoteNumber: "Q-001", QuoteStatusID: 1, Amount: 500.00}
		q := &Quote{ID: 1, JobID: 999, QuoteNumber: "Q-001", QuoteStatusID: 1, Amount: 500.00}

		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		jobChecker.On("GetByID", ctx, int64(999)).Return(nil, ErrInvalidJob)

		err := uc.Update(ctx, q)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidJob, err)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		existing := &Quote{ID: 1, QuoteNumber: "Q-001"}
		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		repo.On("Delete", ctx, int64(1)).Return(nil)

		err := uc.Delete(ctx, 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		repo.On("GetByID", ctx, int64(999)).Return(nil, ErrQuoteNotFound)

		err := uc.Delete(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, ErrQuoteNotFound, err)
	})

	t.Run("repository error", func(t *testing.T) {
		uc, repo, _, _ := newTestUseCase()

		existing := &Quote{ID: 1, QuoteNumber: "Q-001"}
		repo.On("GetByID", ctx, int64(1)).Return(existing, nil)
		repo.On("Delete", ctx, int64(1)).Return(errors.New("db error"))

		err := uc.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
