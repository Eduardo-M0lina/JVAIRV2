package customer

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/customer"
)

func setupTest(t *testing.T) (*Repository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}

	repo := &Repository{db: db}

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestCreate_Success(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()
	c := &customer.Customer{
		Name:       "Test Customer",
		WorkflowID: 1,
	}

	mock.ExpectExec("INSERT INTO customers").
		WithArgs(
			c.Name, c.Email, c.Phone, c.Mobile, c.Fax, c.PhoneOther, c.Website,
			c.ContactName, c.ContactEmail, c.ContactPhone,
			c.BillingAddressStreet, c.BillingAddressCity, c.BillingAddressState, c.BillingAddressZip,
			c.WorkflowID, c.Notes,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(ctx, c)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), c.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate_DatabaseError(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()
	c := &customer.Customer{
		Name:       "Test Customer",
		WorkflowID: 1,
	}

	mock.ExpectExec("INSERT INTO customers").
		WithArgs(
			c.Name, c.Email, c.Phone, c.Mobile, c.Fax, c.PhoneOther, c.Website,
			c.ContactName, c.ContactEmail, c.ContactPhone,
			c.BillingAddressStreet, c.BillingAddressCity, c.BillingAddressState, c.BillingAddressZip,
			c.WorkflowID, c.Notes,
		).
		WillReturnError(sql.ErrConnDone)

	err := repo.Create(ctx, c)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "mobile", "fax", "phone_other", "website",
		"contact_name", "contact_email", "contact_phone",
		"billing_address_street", "billing_address_city", "billing_address_state", "billing_address_zip",
		"workflow_id", "notes", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		1, "Test Customer", email, nil, nil, nil, nil, nil,
		nil, nil, nil,
		nil, nil, nil, nil,
		1, nil, now, now, nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM customers WHERE id = \\? AND deleted_at IS NULL").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Test Customer", result.Name)
	assert.Equal(t, &email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectQuery("SELECT (.+) FROM customers WHERE id = \\? AND deleted_at IS NULL").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "customer not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_Success(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()
	filters := map[string]interface{}{}
	now := time.Now()

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM customers WHERE deleted_at IS NULL").
		WillReturnRows(countRows)

	// Mock list query
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "mobile", "fax", "phone_other", "website",
		"contact_name", "contact_email", "contact_phone",
		"billing_address_street", "billing_address_city", "billing_address_state", "billing_address_zip",
		"workflow_id", "notes", "created_at", "updated_at", "deleted_at",
	}).
		AddRow(1, "Customer 1", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1, nil, now, now, nil).
		AddRow(2, "Customer 2", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1, nil, now, now, nil)

	mock.ExpectQuery("SELECT (.+) FROM customers WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \\? OFFSET \\?").
		WithArgs(10, 0).
		WillReturnRows(rows)

	result, total, err := repo.List(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
	assert.Equal(t, "Customer 1", result[0].Name)
	assert.Equal(t, "Customer 2", result[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestList_WithSearch(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()
	filters := map[string]interface{}{
		"search": "ACME",
	}
	now := time.Now()

	searchPattern := "%ACME%"

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM customers WHERE deleted_at IS NULL AND").
		WithArgs(searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern).
		WillReturnRows(countRows)

	// Mock list query
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "phone", "mobile", "fax", "phone_other", "website",
		"contact_name", "contact_email", "contact_phone",
		"billing_address_street", "billing_address_city", "billing_address_state", "billing_address_zip",
		"workflow_id", "notes", "created_at", "updated_at", "deleted_at",
	}).AddRow(1, "ACME Corp", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1, nil, now, now, nil)

	mock.ExpectQuery("SELECT (.+) FROM customers WHERE deleted_at IS NULL AND").
		WithArgs(searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, 10, 0).
		WillReturnRows(rows)

	result, total, err := repo.List(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, result, 1)
	assert.Equal(t, "ACME Corp", result[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()
	c := &customer.Customer{
		ID:         1,
		Name:       "Updated Customer",
		WorkflowID: 1,
	}

	mock.ExpectExec("UPDATE customers SET").
		WithArgs(
			c.Name, c.Email, c.Phone, c.Mobile, c.Fax, c.PhoneOther, c.Website,
			c.ContactName, c.ContactEmail, c.ContactPhone,
			c.BillingAddressStreet, c.BillingAddressCity, c.BillingAddressState, c.BillingAddressZip,
			c.WorkflowID, c.Notes, c.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(ctx, c)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectExec("UPDATE customers SET deleted_at = NOW\\(\\) WHERE id = \\? AND deleted_at IS NULL").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHasProperties_True(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM properties WHERE customer_id = \\? AND deleted_at IS NULL").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repo.HasProperties(ctx, 1)

	assert.NoError(t, err)
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHasProperties_False(t *testing.T) {
	repo, mock, cleanup := setupTest(t)
	defer cleanup()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM properties WHERE customer_id = \\? AND deleted_at IS NULL").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repo.HasProperties(ctx, 1)

	assert.NoError(t, err)
	assert.False(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}
