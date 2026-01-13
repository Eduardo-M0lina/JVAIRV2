package workflow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_List(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	filters := Filters{Name: "test"}
	page := 1
	pageSize := 10

	expectedWorkflows := []Workflow{
		{ID: 1, Name: "Workflow 1", IsActive: true},
		{ID: 2, Name: "Workflow 2", IsActive: true},
	}
	expectedTotal := int64(2)

	mockRepo.On("List", ctx, filters, page, pageSize).Return(expectedWorkflows, expectedTotal, nil)

	workflows, total, err := useCase.List(ctx, filters, page, pageSize)

	assert.NoError(t, err)
	assert.Equal(t, expectedWorkflows, workflows)
	assert.Equal(t, expectedTotal, total)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflowID := int64(1)

	expectedWorkflow := &Workflow{
		ID:       workflowID,
		Name:     "Test Workflow",
		IsActive: true,
	}

	expectedStatuses := []WorkflowStatus{
		{JobStatusID: 1, WorkflowID: workflowID, Order: 0, StatusName: "Pending"},
		{JobStatusID: 2, WorkflowID: workflowID, Order: 1, StatusName: "In Progress"},
	}

	mockRepo.On("GetByID", ctx, workflowID).Return(expectedWorkflow, nil)
	mockRepo.On("GetWorkflowStatuses", ctx, workflowID).Return(expectedStatuses, nil)

	workflow, err := useCase.GetByID(ctx, workflowID)

	assert.NoError(t, err)
	assert.Equal(t, expectedWorkflow.ID, workflow.ID)
	assert.Equal(t, expectedWorkflow.Name, workflow.Name)
	assert.Equal(t, expectedStatuses, workflow.Statuses)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflowID := int64(999)

	mockRepo.On("GetByID", ctx, workflowID).Return(nil, ErrWorkflowNotFound)

	workflow, err := useCase.GetByID(ctx, workflowID)

	assert.Error(t, err)
	assert.Equal(t, ErrWorkflowNotFound, err)
	assert.Nil(t, workflow)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Create_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflow := &Workflow{
		Name:     "New Workflow",
		IsActive: true,
	}
	statusIDs := []int64{1, 2, 3}

	mockRepo.On("Create", ctx, workflow).Return(nil).Run(func(args mock.Arguments) {
		wf := args.Get(1).(*Workflow)
		wf.ID = 1
	})

	expectedStatuses := []WorkflowStatus{
		{JobStatusID: 1, WorkflowID: 1, Order: 0},
		{JobStatusID: 2, WorkflowID: 1, Order: 1},
		{JobStatusID: 3, WorkflowID: 1, Order: 2},
	}
	mockRepo.On("SetWorkflowStatuses", ctx, int64(1), expectedStatuses).Return(nil)

	err := useCase.Create(ctx, workflow, statusIDs)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), workflow.ID)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Create_EmptyName(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflow := &Workflow{
		Name:     "",
		IsActive: true,
	}

	err := useCase.Create(ctx, workflow, nil)

	assert.Error(t, err)
	assert.Equal(t, ErrWorkflowNameRequired, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUseCase_Update_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflow := &Workflow{
		ID:       1,
		Name:     "Updated Workflow",
		IsActive: true,
	}
	statusIDs := []int64{1, 2}

	existingWorkflow := &Workflow{
		ID:       1,
		Name:     "Old Workflow",
		IsActive: true,
	}

	mockRepo.On("GetByID", ctx, workflow.ID).Return(existingWorkflow, nil)
	mockRepo.On("Update", ctx, workflow).Return(nil)

	expectedStatuses := []WorkflowStatus{
		{JobStatusID: 1, WorkflowID: 1, Order: 0},
		{JobStatusID: 2, WorkflowID: 1, Order: 1},
	}
	mockRepo.On("SetWorkflowStatuses", ctx, workflow.ID, expectedStatuses).Return(nil)

	err := useCase.Update(ctx, workflow, statusIDs)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Update_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflow := &Workflow{
		ID:       999,
		Name:     "Updated Workflow",
		IsActive: true,
	}

	mockRepo.On("GetByID", ctx, workflow.ID).Return(nil, ErrWorkflowNotFound)

	err := useCase.Update(ctx, workflow, nil)

	assert.Error(t, err)
	assert.Equal(t, ErrWorkflowNotFound, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUseCase_Update_EmptyName(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflow := &Workflow{
		ID:       1,
		Name:     "",
		IsActive: true,
	}

	err := useCase.Update(ctx, workflow, nil)

	assert.Error(t, err)
	assert.Equal(t, ErrWorkflowNameRequired, err)
	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestUseCase_Delete_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflowID := int64(1)

	existingWorkflow := &Workflow{
		ID:       workflowID,
		Name:     "Workflow to Delete",
		IsActive: true,
	}

	mockRepo.On("GetByID", ctx, workflowID).Return(existingWorkflow, nil)
	mockRepo.On("Delete", ctx, workflowID).Return(nil)

	err := useCase.Delete(ctx, workflowID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflowID := int64(999)

	mockRepo.On("GetByID", ctx, workflowID).Return(nil, ErrWorkflowNotFound)

	err := useCase.Delete(ctx, workflowID)

	assert.Error(t, err)
	assert.Equal(t, ErrWorkflowNotFound, err)
	mockRepo.AssertNotCalled(t, "Delete")
}

func TestUseCase_Duplicate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflowID := int64(1)

	existingWorkflow := &Workflow{
		ID:       workflowID,
		Name:     "Original Workflow",
		IsActive: true,
	}

	existingStatuses := []WorkflowStatus{
		{JobStatusID: 1, WorkflowID: workflowID, Order: 0},
		{JobStatusID: 2, WorkflowID: workflowID, Order: 1},
	}

	duplicatedWorkflow := &Workflow{
		ID:       2,
		Name:     "Original Workflow",
		IsActive: true,
	}

	mockRepo.On("GetByID", ctx, workflowID).Return(existingWorkflow, nil).Once()
	mockRepo.On("GetWorkflowStatuses", ctx, workflowID).Return(existingStatuses, nil).Once()
	mockRepo.On("Duplicate", ctx, workflowID).Return(duplicatedWorkflow, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(wf *Workflow) bool {
		return wf.ID == 2 && wf.Name == "Copy of Original Workflow (2)"
	})).Return(nil)

	newStatuses := []WorkflowStatus{
		{JobStatusID: 1, WorkflowID: 2, Order: 0},
		{JobStatusID: 2, WorkflowID: 2, Order: 1},
	}
	mockRepo.On("SetWorkflowStatuses", ctx, int64(2), newStatuses).Return(nil)
	mockRepo.On("GetByID", ctx, int64(2)).Return(duplicatedWorkflow, nil).Once()
	mockRepo.On("GetWorkflowStatuses", ctx, int64(2)).Return(newStatuses, nil).Once()

	result, err := useCase.Duplicate(ctx, workflowID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Duplicate_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflowID := int64(999)

	mockRepo.On("GetByID", ctx, workflowID).Return(nil, ErrWorkflowNotFound)

	result, err := useCase.Duplicate(ctx, workflowID)

	assert.Error(t, err)
	assert.Equal(t, ErrWorkflowNotFound, err)
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "Duplicate")
}

func TestUseCase_Duplicate_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	workflowID := int64(1)

	existingWorkflow := &Workflow{
		ID:       workflowID,
		Name:     "Original Workflow",
		IsActive: true,
	}

	existingStatuses := []WorkflowStatus{
		{JobStatusID: 1, WorkflowID: workflowID, Order: 0},
	}

	expectedError := errors.New("database error")

	mockRepo.On("GetByID", ctx, workflowID).Return(existingWorkflow, nil)
	mockRepo.On("GetWorkflowStatuses", ctx, workflowID).Return(existingStatuses, nil)
	mockRepo.On("Duplicate", ctx, workflowID).Return(nil, expectedError)

	result, err := useCase.Duplicate(ctx, workflowID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
