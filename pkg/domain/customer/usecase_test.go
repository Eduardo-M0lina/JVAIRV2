package customer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/workflow"
)

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	customer := &Customer{
		Name:       "Test Customer",
		WorkflowID: 1,
	}

	activeWorkflow := &workflow.Workflow{
		ID:       1,
		Name:     "Test Workflow",
		IsActive: true,
	}

	mockWorkflowRepo.On("GetByID", ctx, int64(1)).Return(activeWorkflow, nil)
	mockRepo.On("Create", ctx, customer).Return(nil)

	err := uc.Create(ctx, customer)

	assert.NoError(t, err)
	mockWorkflowRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreate_WorkflowNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	customer := &Customer{
		Name:       "Test Customer",
		WorkflowID: 999,
	}

	mockWorkflowRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := uc.Create(ctx, customer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid workflow ID")
	mockWorkflowRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreate_WorkflowNotActive(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	customer := &Customer{
		Name:       "Test Customer",
		WorkflowID: 1,
	}

	inactiveWorkflow := &workflow.Workflow{
		ID:       1,
		Name:     "Inactive Workflow",
		IsActive: false,
	}

	mockWorkflowRepo.On("GetByID", ctx, int64(1)).Return(inactiveWorkflow, nil)

	err := uc.Create(ctx, customer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "workflow is not active")
	mockWorkflowRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	expectedCustomer := &Customer{
		ID:         1,
		Name:       "Test Customer",
		WorkflowID: 1,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(expectedCustomer, nil)

	result, err := uc.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomer, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()

	mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := uc.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestList_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	filters := map[string]interface{}{"search": "test"}
	expectedCustomers := []*Customer{
		{ID: 1, Name: "Customer 1"},
		{ID: 2, Name: "Customer 2"},
	}

	mockRepo.On("List", ctx, filters, 1, 10).Return(expectedCustomers, 2, nil)

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomers, result)
	assert.Equal(t, 2, total)
	mockRepo.AssertExpectations(t)
}

func TestList_InvalidPagination(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	filters := map[string]interface{}{}

	mockRepo.On("List", ctx, filters, 1, 10).Return([]*Customer{}, 0, nil)

	// Test with invalid page (should default to 1)
	result, total, err := uc.List(ctx, filters, 0, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, total)
	mockRepo.AssertExpectations(t)
}

func TestList_WithValidWorkflowFilter(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	filters := map[string]interface{}{"workflow_id": int64(1)}
	expectedCustomers := []*Customer{
		{ID: 1, Name: "Customer 1", WorkflowID: 1},
	}

	activeWorkflow := &workflow.Workflow{
		ID:       1,
		Name:     "Test Workflow",
		IsActive: true,
	}

	mockWorkflowRepo.On("GetByID", ctx, int64(1)).Return(activeWorkflow, nil)
	mockRepo.On("List", ctx, filters, 1, 10).Return(expectedCustomers, 1, nil)

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomers, result)
	assert.Equal(t, 1, total)
	mockWorkflowRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestList_WithInvalidWorkflowFilter(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	filters := map[string]interface{}{"workflow_id": int64(999)}

	mockWorkflowRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid workflow_id")
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	mockWorkflowRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "List")
}

func TestList_WithInactiveWorkflowFilter(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	filters := map[string]interface{}{"workflow_id": int64(1)}

	inactiveWorkflow := &workflow.Workflow{
		ID:       1,
		Name:     "Inactive Workflow",
		IsActive: false,
	}

	mockWorkflowRepo.On("GetByID", ctx, int64(1)).Return(inactiveWorkflow, nil)

	result, total, err := uc.List(ctx, filters, 1, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "workflow is not active")
	assert.Nil(t, result)
	assert.Equal(t, 0, total)
	mockWorkflowRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "List")
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	existingCustomer := &Customer{
		ID:         1,
		Name:       "Old Name",
		WorkflowID: 1,
		DeletedAt:  nil,
	}

	updatedCustomer := &Customer{
		ID:         1,
		Name:       "New Name",
		WorkflowID: 1,
	}

	activeWorkflow := &workflow.Workflow{
		ID:       1,
		Name:     "Test Workflow",
		IsActive: true,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	mockWorkflowRepo.On("GetByID", ctx, int64(1)).Return(activeWorkflow, nil)
	mockRepo.On("Update", ctx, updatedCustomer).Return(nil)

	err := uc.Update(ctx, updatedCustomer)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockWorkflowRepo.AssertExpectations(t)
}

func TestUpdate_CustomerDeleted(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	now := time.Now()
	deletedCustomer := &Customer{
		ID:         1,
		Name:       "Deleted Customer",
		WorkflowID: 1,
		DeletedAt:  &now,
	}

	updatedCustomer := &Customer{
		ID:         1,
		Name:       "New Name",
		WorkflowID: 1,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(deletedCustomer, nil)

	err := uc.Update(ctx, updatedCustomer)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot update deleted customer")
	mockRepo.AssertExpectations(t)
	mockWorkflowRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	existingCustomer := &Customer{
		ID:         1,
		Name:       "Test Customer",
		WorkflowID: 1,
		DeletedAt:  nil,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	mockRepo.On("HasProperties", ctx, int64(1)).Return(false, nil)
	mockRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := uc.Delete(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_HasProperties(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	existingCustomer := &Customer{
		ID:         1,
		Name:       "Test Customer",
		WorkflowID: 1,
		DeletedAt:  nil,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(existingCustomer, nil)
	mockRepo.On("HasProperties", ctx, int64(1)).Return(true, nil)

	err := uc.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete customer with associated properties")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDelete_AlreadyDeleted(t *testing.T) {
	mockRepo := new(MockRepository)
	mockWorkflowRepo := new(workflow.MockRepository)
	uc := NewUseCase(mockRepo, mockWorkflowRepo)

	ctx := context.Background()
	now := time.Now()
	deletedCustomer := &Customer{
		ID:         1,
		Name:       "Deleted Customer",
		WorkflowID: 1,
		DeletedAt:  &now,
	}

	mockRepo.On("GetByID", ctx, int64(1)).Return(deletedCustomer, nil)

	err := uc.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "customer already deleted")
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "HasProperties")
	mockRepo.AssertNotCalled(t, "Delete")
}
