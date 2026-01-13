package settings

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	domainSettings "github.com/your-org/jvairv2/pkg/domain/settings"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *Repository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear mock de base de datos: %v", err)
	}

	repo := &Repository{db: db}
	return db, mock, repo
}

func TestGet_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	now := time.Now()
	twilioSID := "test_sid"
	twilioToken := "test_token"
	twilioNumber := "+1234567890"

	rows := sqlmock.NewRows([]string{
		"id", "is_twilio_enabled", "twilio_sid", "twilio_auth_token", "twilio_from_number",
		"is_enforce_routine_password_reset", "password_expire_days", "password_history_count",
		"password_minimum_length", "password_age", "password_include_numbers",
		"password_include_symbols", "created_at", "updated_at",
	}).AddRow(
		1, true, twilioSID, twilioToken, twilioNumber,
		true, 90, 10, 8, 5, true, true, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM settings LIMIT 1").
		WillReturnRows(rows)

	settings, err := repo.Get(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, int64(1), settings.ID)
	assert.True(t, settings.IsTwilioEnabled)
	assert.Equal(t, twilioSID, *settings.TwilioSID)
	assert.Equal(t, twilioToken, *settings.TwilioAuthToken)
	assert.Equal(t, twilioNumber, *settings.TwilioFromNumber)
	assert.True(t, settings.IsEnforceRoutinePasswordReset)
	assert.Equal(t, 90, settings.PasswordExpireDays)
	assert.Equal(t, 10, settings.PasswordHistoryCount)
	assert.Equal(t, 8, settings.PasswordMinimumLength)
	assert.Equal(t, 5, settings.PasswordAge)
	assert.True(t, settings.PasswordIncludeNumbers)
	assert.True(t, settings.PasswordIncludeSymbols)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_WithNullTwilioFields(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "is_twilio_enabled", "twilio_sid", "twilio_auth_token", "twilio_from_number",
		"is_enforce_routine_password_reset", "password_expire_days", "password_history_count",
		"password_minimum_length", "password_age", "password_include_numbers",
		"password_include_symbols", "created_at", "updated_at",
	}).AddRow(
		1, false, nil, nil, nil,
		true, 90, 10, 8, 5, true, true, now, now,
	)

	mock.ExpectQuery("SELECT (.+) FROM settings LIMIT 1").
		WillReturnRows(rows)

	settings, err := repo.Get(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.False(t, settings.IsTwilioEnabled)
	assert.Nil(t, settings.TwilioSID)
	assert.Nil(t, settings.TwilioAuthToken)
	assert.Nil(t, settings.TwilioFromNumber)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	mock.ExpectQuery("SELECT (.+) FROM settings LIMIT 1").
		WillReturnError(sql.ErrNoRows)

	settings, err := repo.Get(context.Background())

	assert.Error(t, err)
	assert.Equal(t, domainSettings.ErrSettingsNotFound, err)
	assert.Nil(t, settings)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGet_DatabaseError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	expectedError := sql.ErrConnDone

	mock.ExpectQuery("SELECT (.+) FROM settings LIMIT 1").
		WillReturnError(expectedError)

	settings, err := repo.Get(context.Background())

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, settings)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	twilioSID := "updated_sid"
	settings := &domainSettings.Settings{
		ID:                            1,
		IsTwilioEnabled:               true,
		TwilioSID:                     &twilioSID,
		TwilioAuthToken:               nil,
		TwilioFromNumber:              nil,
		IsEnforceRoutinePasswordReset: true,
		PasswordExpireDays:            90,
		PasswordHistoryCount:          10,
		PasswordMinimumLength:         8,
		PasswordAge:                   5,
		PasswordIncludeNumbers:        true,
		PasswordIncludeSymbols:        true,
	}

	mock.ExpectExec("UPDATE settings SET (.+) WHERE id = ?").
		WithArgs(
			settings.IsTwilioEnabled,
			settings.TwilioSID,
			settings.TwilioAuthToken,
			settings.TwilioFromNumber,
			settings.IsEnforceRoutinePasswordReset,
			settings.PasswordExpireDays,
			settings.PasswordHistoryCount,
			settings.PasswordMinimumLength,
			settings.PasswordAge,
			settings.PasswordIncludeNumbers,
			settings.PasswordIncludeSymbols,
			settings.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(context.Background(), settings)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_DatabaseError(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	settings := &domainSettings.Settings{
		ID:                            1,
		IsTwilioEnabled:               false,
		IsEnforceRoutinePasswordReset: true,
		PasswordExpireDays:            90,
		PasswordHistoryCount:          10,
		PasswordMinimumLength:         8,
		PasswordAge:                   5,
		PasswordIncludeNumbers:        true,
		PasswordIncludeSymbols:        true,
	}

	expectedError := sql.ErrConnDone

	mock.ExpectExec("UPDATE settings SET (.+) WHERE id = ?").
		WithArgs(
			settings.IsTwilioEnabled,
			settings.TwilioSID,
			settings.TwilioAuthToken,
			settings.TwilioFromNumber,
			settings.IsEnforceRoutinePasswordReset,
			settings.PasswordExpireDays,
			settings.PasswordHistoryCount,
			settings.PasswordMinimumLength,
			settings.PasswordAge,
			settings.PasswordIncludeNumbers,
			settings.PasswordIncludeSymbols,
			settings.ID,
		).
		WillReturnError(expectedError)

	err := repo.Update(context.Background(), settings)

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
