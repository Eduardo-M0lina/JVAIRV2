package workflow

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	domainWorkflow "github.com/your-org/jvairv2/pkg/domain/workflow"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *Repository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}

	repo := &Repository{db: db}
	return db, mock, repo
}

func TestList_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	now := time.Now()
	notes := "Test notes"

	// Mock count query
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM workflows").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Mock select query
	rows := sqlmock.NewRows([]string{
		"id", "name", "notes", "is_active", "created_at", "updated_at",
	}).
		AddRow(1, "Workflow 1", notes, true, now, now).
		AddRow(2, "Workflow 2", nil, false, now, now)

	mock.ExpectQuery("SELECT (.+) FROM workflows ORDER BY id DESC LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnRows(rows)

	workflows, total, err := repo.List(context.Background(), domainWorkflow.Filters{}, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, workflows, 2)
	assert.Equal(t, "Workflow 1", workflows[0].Name)
	assert.Equal(t, &notes, workflows[0].Notes)
	assert.True(t, workflows[0].IsActive)
	assert.Equal(t, "Workflow 2", workflows[1].Name)
	assert.Nil(t, workflows[1].Notes)
	assert.False(t, workflows[1].IsActive)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithFilters(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	isActive := true
	filters := domainWorkflow.Filters{
		Name:     "Test",
		IsActive: &isActive,
		Search:   "search",
	}

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM workflows WHERE (.+)").
		WithArgs("%Test%", isActive, "%search%", "%search%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	rows := sqlmock.NewRows([]string{
		"id", "name", "notes", "is_active", "created_at", "updated_at",
	}).AddRow(1, "Test Workflow", nil, true, time.Now(), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM workflows WHERE (.+) ORDER BY id DESC LIMIT \\? OFFSET \\?").
		WithArgs("%Test%", isActive, "%search%", "%search%", 10, 0).
		WillReturnRows(rows)

	workflows, total, err := repo.List(context.Background(), filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, workflows, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	now := time.Now()
	notes := "Test notes"

	rows := sqlmock.NewRows([]string{
		"id", "name", "notes", "is_active", "created_at", "updated_at",
	}).AddRow(1, "Test Workflow", notes, true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM workflows WHERE id = \\?").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	workflow, err := repo.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	assert.Equal(t, int64(1), workflow.ID)
	assert.Equal(t, "Test Workflow", workflow.Name)
	assert.Equal(t, &notes, workflow.Notes)
	assert.True(t, workflow.IsActive)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	mock.ExpectQuery("SELECT (.+) FROM workflows WHERE id = \\?").
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	workflow, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, domainWorkflow.ErrWorkflowNotFound, err)
	assert.Nil(t, workflow)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	notes := "New workflow"
	workflow := &domainWorkflow.Workflow{
		Name:     "New Workflow",
		Notes:    &notes,
		IsActive: true,
	}

	mock.ExpectExec("INSERT INTO workflows").
		WithArgs(workflow.Name, workflow.Notes, workflow.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), workflow)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), workflow.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DatabaseError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	workflow := &domainWorkflow.Workflow{
		Name:     "New Workflow",
		IsActive: true,
	}

	expectedError := sql.ErrConnDone

	mock.ExpectExec("INSERT INTO workflows").
		WithArgs(workflow.Name, workflow.Notes, workflow.IsActive).
		WillReturnError(expectedError)

	err := repo.Create(context.Background(), workflow)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	notes := "Updated workflow"
	workflow := &domainWorkflow.Workflow{
		ID:       1,
		Name:     "Updated Workflow",
		Notes:    &notes,
		IsActive: false,
	}

	mock.ExpectExec("UPDATE workflows SET (.+) WHERE id = \\?").
		WithArgs(workflow.Name, workflow.Notes, workflow.IsActive, workflow.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(context.Background(), workflow)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM job_status_workflow WHERE workflow_id = \\?").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectExec("DELETE FROM workflows WHERE id = \\?").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_TransactionError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	expectedError := sql.ErrConnDone

	mock.ExpectBegin().WillReturnError(expectedError)

	err := repo.Delete(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDuplicate_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	now := time.Now()
	notes := "Original workflow"

	// Mock GetByID
	rows := sqlmock.NewRows([]string{
		"id", "name", "notes", "is_active", "created_at", "updated_at",
	}).AddRow(1, "Original Workflow", notes, true, now, now)

	mock.ExpectQuery("SELECT (.+) FROM workflows WHERE id = \\?").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	// Mock Create
	mock.ExpectExec("INSERT INTO workflows").
		WithArgs("Original Workflow", &notes, true).
		WillReturnResult(sqlmock.NewResult(2, 1))

	duplicated, err := repo.Duplicate(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, duplicated)
	assert.Equal(t, int64(2), duplicated.ID)
	assert.Equal(t, "Original Workflow", duplicated.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWorkflowStatuses_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	rows := sqlmock.NewRows([]string{
		"job_status_id", "workflow_id", "order", "label",
	}).
		AddRow(1, 1, 0, "Pending").
		AddRow(2, 1, 1, "In Progress").
		AddRow(3, 1, 2, "Completed")

	mock.ExpectQuery("SELECT (.+) FROM job_status_workflow (.+) INNER JOIN job_statuses (.+) WHERE (.+) ORDER BY (.+)").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	statuses, err := repo.GetWorkflowStatuses(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, statuses, 3)
	assert.Equal(t, int64(1), statuses[0].JobStatusID)
	assert.Equal(t, 0, statuses[0].Order)
	assert.Equal(t, "Pending", statuses[0].StatusName)
	assert.Equal(t, int64(2), statuses[1].JobStatusID)
	assert.Equal(t, 1, statuses[1].Order)
	assert.Equal(t, "In Progress", statuses[1].StatusName)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWorkflowStatuses_Empty(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	rows := sqlmock.NewRows([]string{
		"job_status_id", "workflow_id", "order", "label",
	})

	mock.ExpectQuery("SELECT (.+) FROM job_status_workflow (.+) INNER JOIN job_statuses (.+) WHERE (.+) ORDER BY (.+)").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	statuses, err := repo.GetWorkflowStatuses(context.Background(), 1)

	assert.NoError(t, err)
	assert.Len(t, statuses, 0)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetWorkflowStatuses_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	statuses := []domainWorkflow.WorkflowStatus{
		{JobStatusID: 1, WorkflowID: 1, Order: 0},
		{JobStatusID: 2, WorkflowID: 1, Order: 1},
	}

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM job_status_workflow WHERE workflow_id = \\?").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectPrepare("INSERT INTO job_status_workflow")
	mock.ExpectExec("INSERT INTO job_status_workflow").
		WithArgs(int64(1), int64(1), 0).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO job_status_workflow").
		WithArgs(int64(2), int64(1), 1).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	err := repo.SetWorkflowStatuses(context.Background(), 1, statuses)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetWorkflowStatuses_EmptyStatuses(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM job_status_workflow WHERE workflow_id = \\?").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.SetWorkflowStatuses(context.Background(), 1, []domainWorkflow.WorkflowStatus{})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetWorkflowStatuses_TransactionError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	expectedError := sql.ErrConnDone

	mock.ExpectBegin().WillReturnError(expectedError)

	err := repo.SetWorkflowStatuses(context.Background(), 1, []domainWorkflow.WorkflowStatus{})

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository(t *testing.T) {
	db, _, _ := setupMockDB(t)
	defer func() { _ = db.Close() }()

	repo := NewRepository(db)

	assert.NotNil(t, repo)
	assert.IsType(t, &Repository{}, repo)
}
