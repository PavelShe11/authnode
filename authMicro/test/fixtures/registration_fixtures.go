package fixtures

// NewValidUserData returns valid user registration data
func NewValidUserData() map[string]any {
	return map[string]any{
		"email":     "test@example.com",
		"firstName": "John",
		"lastName":  "Doe",
		"password":  "SecurePass123",
	}
}
