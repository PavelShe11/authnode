package fixtures

import (
	"time"

	"github.com/PavelShe11/studbridge/authMicro/internal/entity"
)

// NewValidUserData returns valid user registration data
func NewValidUserData() map[string]any {
	return map[string]any{
		"email":     "test@example.com",
		"firstName": "John",
		"lastName":  "Doe",
		"password":  "SecurePass123",
	}
}

// NewValidSession returns a typical registration session
func NewValidSession(email string) *entity.RegistrationSession {
	return &entity.RegistrationSession{
		Id:          "test-session-id",
		Email:       email,
		Code:        "hashed_code_123456",
		CodeExpires: time.Now().Add(2 * time.Minute),
		CreatedAt:   time.Now(),
	}
}

// NewExpiredSession returns an expired registration session
func NewExpiredSession(email string) *entity.RegistrationSession {
	return &entity.RegistrationSession{
		Id:          "expired-session-id",
		Email:       email,
		Code:        "hashed_code_123456",
		CodeExpires: time.Now().Add(-1 * time.Minute),
		CreatedAt:   time.Now().Add(-10 * time.Minute),
	}
}

// NewAccount returns an existing account
func NewAccount(email string) *entity.Account {
	return &entity.Account{
		Email: email,
	}
}
