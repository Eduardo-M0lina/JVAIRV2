package supervisor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/customer"
)

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	sup := &Supervisor{
		Name:       "John Doe",
		CustomerID: 1,
	}

	existingCustomer := &customer.Customer{
		ID:   1,
		Name: "ACME Corp",
	}

	mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	mockRepo.On("Create", ctx, sup).Return(nil)

	err := uc.Create(ctx, sup)

	assert.NoError(t, err)
	mockCustomerRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreate_CustomerNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	sup := &Supervisor{
		Name:       "John Doe",
		CustomerID: 999,
	}

	mockCustomerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := uc.Create(ctx, sup)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid customer ID")
	mockCustomerRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreate_RepoError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	sup := &Supervisor{
		Name:       "John Doe",
		CustomerID: 1,
	}

	existingCustomer := &customer.Customer{ID: 1, Name: "ACME Corp"}
	mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	mockRepo.On("Create", ctx, sup).Return(errors.New("db error"))

	err := uc.Create(ctx, sup)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockCustomerRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	expected := &Supervisor{
		ID:         1,
		Name:       "John Doe",
		CustomerID: 1,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(expected, nil)

	result, err := uc.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := uc.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestList_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	filters := map[string]interface{}{"search": "john"}
	expected := []*Supervisor{
		{ID: 1, Name: "John Doe", CustomerID: 1},
		{ID: 2, Name: "John Smith", CustomerID: 2},
	}

	mockRepo.On("List", ctx, filters, 1, 10).Return(expected, 2, nil)

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, 2, total)
	mockRepo.AssertExpectations(t)
}

func TestList_InvalidPagination(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	filters := map[string]interface{}{}

	mockRepo.On("List", ctx, filters, 1, 10).Return([]*Supervisor{}, 0, nil)

	result, total, err := uc.List(ctx, filters, 0, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, total)
	mockRepo.AssertExpectations(t)
}

func TestList_WithValidCustomerFilter(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	filters := map[string]interface{}{"customer_id": int64(1)}
	expected := []*Supervisor{
		{ID: 1, Name: "John Doe", CustomerID: 1},
	}

	existingCustomer := &customer.Customer{ID: 1, Name: "ACME Corp"}
	mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	mockRepo.On("List", ctx, filters, 1, 10).Return(expected, 1, nil)

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, total)
	mockCustomerRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestList_WithInvalidCustomerFilter(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	filters := map[string]interface{}{"customer_id": int64(999)}

	mockCustomerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid customer_id")
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	mockCustomerRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "List")
}

func TestList_RepoError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	filters := map[string]interface{}{}

	mockRepo.On("List", ctx, filters, 1, 10).Return(nil, 0, errors.New("db error"))

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	mockRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	existing := &Supervisor{
		ID:         1,
		Name:       "Old Name",
		CustomerID: 1,
		DeletedAt:  nil,
	}

	updated := &Supervisor{
		ID:         1,
		Name:       "New Name",
		CustomerID: 1,
	}

	existingCustomer := &customer.Customer{ID: 1, Name: "ACME Corp"}

	mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	mockCustomerRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	mockRepo.On("Update", ctx, updated).Return(nil)

	err := uc.Update(ctx, updated)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCustomerRepo.AssertExpectations(t)
}

func TestUpdate_SupervisorDeleted(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	now := time.Now()
	deleted := &Supervisor{
		ID:         1,
		Name:       "Deleted Supervisor",
		CustomerID: 1,
		DeletedAt:  &now,
	}

	updated := &Supervisor{
		ID:         1,
		Name:       "New Name",
		CustomerID: 1,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(deleted, nil)

	err := uc.Update(ctx, updated)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot update deleted supervisor")
	mockRepo.AssertExpectations(t)
	mockCustomerRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdate_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	updated := &Supervisor{
		ID:         999,
		Name:       "New Name",
		CustomerID: 1,
	}

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := uc.Update(ctx, updated)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
	mockCustomerRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdate_InvalidCustomer(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	existing := &Supervisor{
		ID:         1,
		Name:       "John Doe",
		CustomerID: 1,
		DeletedAt:  nil,
	}

	updated := &Supervisor{
		ID:         1,
		Name:       "John Doe",
		CustomerID: 999,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	mockCustomerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := uc.Update(ctx, updated)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid customer ID")
	mockRepo.AssertExpectations(t)
	mockCustomerRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	existing := &Supervisor{
		ID:         1,
		Name:       "John Doe",
		CustomerID: 1,
		DeletedAt:  nil,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := uc.Delete(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_AlreadyDeleted(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	now := time.Now()
	deleted := &Supervisor{
		ID:         1,
		Name:       "Deleted Supervisor",
		CustomerID: 1,
		DeletedAt:  &now,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(deleted, nil)

	err := uc.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "supervisor already deleted")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDelete_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := uc.Delete(ctx, 999)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDelete_RepoError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockCustomerRepo := new(customer.MockRepository)
	uc := NewUseCase(mockRepo, mockCustomerRepo)

	ctx := context.Background()
	existing := &Supervisor{
		ID:         1,
		Name:       "John Doe",
		CustomerID: 1,
		DeletedAt:  nil,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(errors.New("db error"))

	err := uc.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}
