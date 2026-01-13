package settings

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUseCase_Get(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	now := time.Now()
	twilioSID := "test_sid"
	twilioToken := "test_token"
	twilioNumber := "+1234567890"

	expectedSettings := &Settings{
		ID:                            1,
		IsTwilioEnabled:               true,
		TwilioSID:                     &twilioSID,
		TwilioAuthToken:               &twilioToken,
		TwilioFromNumber:              &twilioNumber,
		IsEnforceRoutinePasswordReset: true,
		PasswordExpireDays:            90,
		PasswordHistoryCount:          10,
		PasswordMinimumLength:         8,
		PasswordAge:                   5,
		PasswordIncludeNumbers:        true,
		PasswordIncludeSymbols:        true,
		CreatedAt:                     &now,
		UpdatedAt:                     &now,
	}

	mockRepo.On("Get", ctx).Return(expectedSettings, nil)

	settings, err := useCase.Get(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedSettings, settings)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Get_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	expectedError := errors.New("database error")

	mockRepo.On("Get", ctx).Return(nil, expectedError)

	settings, err := useCase.Get(ctx)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, settings)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Update_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	settings := &Settings{
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

	mockRepo.On("Update", ctx, settings).Return(nil)

	err := useCase.Update(ctx, settings)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUseCase_Update_ValidationError_PasswordExpireDays(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	settings := &Settings{
		ID:                    1,
		PasswordExpireDays:    0,
		PasswordHistoryCount:  10,
		PasswordMinimumLength: 8,
		PasswordAge:           5,
	}

	err := useCase.Update(ctx, settings)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPasswordExpireDays, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUseCase_Update_ValidationError_PasswordHistoryCount(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	settings := &Settings{
		ID:                    1,
		PasswordExpireDays:    90,
		PasswordHistoryCount:  -1,
		PasswordMinimumLength: 8,
		PasswordAge:           5,
	}

	err := useCase.Update(ctx, settings)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPasswordHistoryCount, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUseCase_Update_ValidationError_PasswordMinimumLength(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	settings := &Settings{
		ID:                    1,
		PasswordExpireDays:    90,
		PasswordHistoryCount:  10,
		PasswordMinimumLength: 3,
		PasswordAge:           5,
	}

	err := useCase.Update(ctx, settings)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPasswordMinimumLength, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUseCase_Update_ValidationError_PasswordAge(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	settings := &Settings{
		ID:                    1,
		PasswordExpireDays:    90,
		PasswordHistoryCount:  10,
		PasswordMinimumLength: 8,
		PasswordAge:           -1,
	}

	err := useCase.Update(ctx, settings)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPasswordAge, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUseCase_Update_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	useCase := NewUseCase(mockRepo)

	ctx := context.Background()
	settings := &Settings{
		ID:                    1,
		PasswordExpireDays:    90,
		PasswordHistoryCount:  10,
		PasswordMinimumLength: 8,
		PasswordAge:           5,
	}
	expectedError := errors.New("database error")

	mockRepo.On("Update", ctx, settings).Return(expectedError)

	err := useCase.Update(ctx, settings)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
