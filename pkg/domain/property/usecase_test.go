package property

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/customer"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		prop := &Property{
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
		}

		mockCustomer := &customer.Customer{
			ID:        1,
			Name:      "Test Customer",
			CreatedAt: &now,
		}

		mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(mockCustomer, nil)
		mockRepo.On("Create", ctx, prop).Return(nil)

		err := uc.Create(ctx, prop)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("invalid customer", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		prop := &Property{
			CustomerID: 999,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
		}

		mockCustomerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("customer not found"))

		err := uc.Create(ctx, prop)

		assert.Error(t, err)
		assert.Equal(t, "invalid customer_id", err.Error())
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("deleted customer", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		prop := &Property{
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
		}

		mockCustomer := &customer.Customer{
			ID:        1,
			Name:      "Test Customer",
			DeletedAt: &now,
		}

		mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(mockCustomer, nil)

		err := uc.Create(ctx, prop)

		assert.Error(t, err)
		assert.Equal(t, "customer is deleted", err.Error())
		mockCustomerRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		expectedProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			CreatedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expectedProp, nil)

		prop, err := uc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedProp, prop)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property not found"))

		prop, err := uc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, prop)
		mockRepo.AssertExpectations(t)
	})
}

func TestList(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success without filters", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		expectedProps := []*Property{
			{
				ID:         1,
				CustomerID: 1,
				Street:     "123 Main St",
				City:       "Atlanta",
				State:      "GA",
				Zip:        "30301",
				CreatedAt:  &now,
			},
		}

		filters := make(map[string]interface{})
		mockRepo.On("List", ctx, filters, 1, 10).Return(expectedProps, 1, nil)

		props, total, err := uc.List(ctx, filters, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedProps, props)
		assert.Equal(t, 1, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success with customer filter", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		mockCustomer := &customer.Customer{
			ID:        1,
			Name:      "Test Customer",
			CreatedAt: &now,
		}

		expectedProps := []*Property{
			{
				ID:         1,
				CustomerID: 1,
				Street:     "123 Main St",
				City:       "Atlanta",
				State:      "GA",
				Zip:        "30301",
				CreatedAt:  &now,
			},
		}

		filters := map[string]interface{}{
			"customer_id": int64(1),
		}

		mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(mockCustomer, nil)
		mockRepo.On("List", ctx, filters, 1, 10).Return(expectedProps, 1, nil)

		props, total, err := uc.List(ctx, filters, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedProps, props)
		assert.Equal(t, 1, total)
		mockRepo.AssertExpectations(t)
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("invalid customer filter", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		filters := map[string]interface{}{
			"customer_id": int64(999),
		}

		mockCustomerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("customer not found"))

		props, total, err := uc.List(ctx, filters, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, props)
		assert.Equal(t, 0, total)
		assert.Equal(t, "invalid customer_id", err.Error())
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("deleted customer filter", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		mockCustomer := &customer.Customer{
			ID:        1,
			Name:      "Test Customer",
			DeletedAt: &now,
		}

		filters := map[string]interface{}{
			"customer_id": int64(1),
		}

		mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(mockCustomer, nil)

		props, total, err := uc.List(ctx, filters, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, props)
		assert.Equal(t, 0, total)
		assert.Equal(t, "customer is deleted", err.Error())
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("invalid page defaults to 1", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		expectedProps := []*Property{}
		filters := make(map[string]interface{})
		mockRepo.On("List", ctx, filters, 1, 10).Return(expectedProps, 0, nil)

		props, total, err := uc.List(ctx, filters, 0, 10)

		assert.NoError(t, err)
		assert.Equal(t, expectedProps, props)
		assert.Equal(t, 0, total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid pageSize defaults to 10", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		expectedProps := []*Property{}
		filters := make(map[string]interface{})
		mockRepo.On("List", ctx, filters, 1, 10).Return(expectedProps, 0, nil)

		props, total, err := uc.List(ctx, filters, 1, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedProps, props)
		assert.Equal(t, 0, total)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			CreatedAt:  &now,
		}

		updatedProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "456 Oak Ave",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30302",
		}

		mockCustomer := &customer.Customer{
			ID:        1,
			Name:      "Test Customer",
			CreatedAt: &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)
		mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(mockCustomer, nil)
		mockRepo.On("Update", ctx, updatedProp).Return(nil)

		err := uc.Update(ctx, updatedProp)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("property not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		prop := &Property{
			ID:         999,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
		}

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property not found"))

		err := uc.Update(ctx, prop)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("cannot update deleted property", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			DeletedAt:  &now,
		}

		updatedProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "456 Oak Ave",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30302",
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)

		err := uc.Update(ctx, updatedProp)

		assert.Error(t, err)
		assert.Equal(t, "cannot update deleted property", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid customer", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			CreatedAt:  &now,
		}

		updatedProp := &Property{
			ID:         1,
			CustomerID: 999,
			Street:     "456 Oak Ave",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30302",
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)
		mockCustomerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("customer not found"))

		err := uc.Update(ctx, updatedProp)

		assert.Error(t, err)
		assert.Equal(t, "invalid customer_id", err.Error())
		mockRepo.AssertExpectations(t)
		mockCustomerRepo.AssertExpectations(t)
	})

	t.Run("deleted customer", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			CreatedAt:  &now,
		}

		updatedProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "456 Oak Ave",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30302",
		}

		mockCustomer := &customer.Customer{
			ID:        1,
			Name:      "Test Customer",
			DeletedAt: &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)
		mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(mockCustomer, nil)

		err := uc.Update(ctx, updatedProp)

		assert.Error(t, err)
		assert.Equal(t, "customer is deleted", err.Error())
		mockRepo.AssertExpectations(t)
		mockCustomerRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			CreatedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)
		mockRepo.On("HasJobs", ctx, int64(1)).Return(false, nil)
		mockRepo.On("Delete", ctx, int64(1)).Return(nil)

		err := uc.Delete(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("property not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("property not found"))

		err := uc.Delete(ctx, 999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("property already deleted", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			DeletedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)

		err := uc.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "property already deleted", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("cannot delete property with jobs", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			CreatedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)
		mockRepo.On("HasJobs", ctx, int64(1)).Return(true, nil)

		err := uc.Delete(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, "cannot delete property with associated jobs", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("error checking jobs", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockCustomerRepo := new(customer.MockRepository)
		uc := NewUseCase(mockRepo, mockCustomerRepo)

		existingProp := &Property{
			ID:         1,
			CustomerID: 1,
			Street:     "123 Main St",
			City:       "Atlanta",
			State:      "GA",
			Zip:        "30301",
			CreatedAt:  &now,
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(existingProp, nil)
		mockRepo.On("HasJobs", ctx, int64(1)).Return(false, errors.New("database error"))

		err := uc.Delete(ctx, 1)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
