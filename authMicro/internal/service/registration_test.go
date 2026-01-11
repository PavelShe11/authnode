package service

import (
	"context"
	"testing"
	"time"

	"github.com/PavelShe11/studbridge/authMicro/internal/config"
	"github.com/PavelShe11/studbridge/authMicro/internal/entity"
	"github.com/PavelShe11/studbridge/authMicro/test/fixtures"
	"github.com/PavelShe11/studbridge/authMicro/test/mocks"
	"github.com/PavelShe11/studbridge/authMicro/utlis/hash"
	commonEntity "github.com/PavelShe11/studbridge/common/entity"
	"github.com/PavelShe11/studbridge/common/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// noopLogger is a test logger that does nothing
type noopLogger struct{}

func (n *noopLogger) Debug(args ...interface{})                 {}
func (n *noopLogger) Debugf(format string, args ...interface{}) {}
func (n *noopLogger) Info(args ...interface{})                  {}
func (n *noopLogger) Infof(format string, args ...interface{})  {}
func (n *noopLogger) Warn(args ...interface{})                  {}
func (n *noopLogger) Warnf(format string, args ...interface{})  {}
func (n *noopLogger) Error(args ...interface{})                 {}
func (n *noopLogger) Errorf(format string, args ...interface{}) {}
func (n *noopLogger) Fatal(args ...interface{})                 {}
func (n *noopLogger) Fatalf(format string, args ...interface{}) {}

func newNoopLogger() logger.Logger {
	return &noopLogger{}
}

// setupService - helper for creating service with mocks
func setupService(t *testing.T) (
	*RegistrationService,
	*mocks.MockRegistrationSessionRepository,
	*mocks.MockAccountProvider,
) {
	mockRepo := new(mocks.MockRegistrationSessionRepository)
	mockProvider := new(mocks.MockAccountProvider)
	testLogger := newNoopLogger()

	cfg := config.CodeGenConfig{
		CodePattern:   "[0-9]{6}",
		CodeMaxLength: 6,
		CodeTTL:       2 * time.Minute,
	}

	service := NewRegistrationService(mockRepo, mockProvider, testLogger, cfg)

	return service, mockRepo, mockProvider
}

// TestRegister_NewUser_Success - new user, code generated
func TestRegister_NewUser_Success(t *testing.T) {
	t.Parallel()

	// ARRANGE
	service, mockRepo, mockProvider := setupService(t)

	userData := fixtures.NewValidUserData()
	userData["email"] = "newuser@test.com"

	// Setup mocks
	mockProvider.On("ValidateAccountData", mock.Anything, userData, "en").Return(nil)
	mockProvider.On("GetAccountByEmail", mock.Anything, "newuser@test.com").Return(nil, nil)
	mockRepo.On("FindByEmail", mock.Anything, "newuser@test.com").Return(nil, nil)
	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(s *entity.RegistrationSession) bool {
		return s.Email == "newuser@test.com" && s.Code != ""
	})).Return(nil)
	mockRepo.On("CleanExpired", mock.Anything).Return(nil)

	// ACT
	result, err := service.Register(context.Background(), userData, "en")

	// ASSERT
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.CodeExpires, time.Now().Unix())
	assert.Equal(t, "[0-9]{6}", result.CodePattern)

	mockProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestRegister_ExistingUser_EmptyCode - existing email, no code (anti-enumeration)
func TestRegister_ExistingUser_EmptyCode(t *testing.T) {
	t.Parallel()

	service, mockRepo, mockProvider := setupService(t)

	userData := map[string]any{
		"email": "existing@test.com",
	}

	existingAccount := &entity.Account{
		Email: "existing@test.com",
	}

	mockProvider.On("ValidateAccountData", mock.Anything, userData, "en").Return(nil)
	mockProvider.On("GetAccountByEmail", mock.Anything, "existing@test.com").Return(existingAccount, nil)
	mockRepo.On("FindByEmail", mock.Anything, "existing@test.com").Return(nil, nil)
	mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(s *entity.RegistrationSession) bool {
		return s.Email == "existing@test.com" && s.Code == ""
	})).Return(nil)
	mockRepo.On("CleanExpired", mock.Anything).Return(nil)

	result, err := service.Register(context.Background(), userData, "en")

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestRegister_ValidationError_ReturnsError - ValidateAccountData fails
func TestRegister_ValidationError_ReturnsError(t *testing.T) {
	t.Parallel()

	service, mockRepo, mockProvider := setupService(t)

	userData := map[string]any{
		"email": "invalid-email",
	}

	validationErr := commonEntity.NewValidationError()

	mockProvider.On("ValidateAccountData", mock.Anything, userData, "en").Return(validationErr)
	mockRepo.On("CleanExpired", mock.Anything).Return(nil)

	result, err := service.Register(context.Background(), userData, "en")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, validationErr, err)

	mockProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestConfirmRegistration_Success - valid code, account created
func TestConfirmRegistration_Success(t *testing.T) {
	t.Parallel()

	service, mockRepo, mockProvider := setupService(t)

	plainCode := "123456"
	hashedCode, _ := hash.HashCode(plainCode)

	userData := map[string]any{
		"email": "test@test.com",
		"code":  plainCode,
	}

	existingSession := &entity.RegistrationSession{
		Email:       "test@test.com",
		Code:        hashedCode,
		CodeExpires: time.Now().Add(1 * time.Minute),
		CreatedAt:   time.Now(),
	}

	mockProvider.On("ValidateAccountData", mock.Anything, userData, "en").Return(nil)
	mockRepo.On("FindByEmail", mock.Anything, "test@test.com").Return(existingSession, nil)
	mockProvider.On("CreateAccount", mock.Anything, userData, "en").Return(nil)
	mockRepo.On("DeleteByEmail", mock.Anything, "test@test.com").Return(nil)

	err := service.ConfirmRegistration(context.Background(), userData, "en")

	assert.NoError(t, err)

	mockProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
