package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost represents bcrypt cost factor (2^12 iterations)
	// Balances security vs performance for verification codes (~250-400ms)
	DefaultCost = 12
)

// HashCode generates a bcrypt hash of the verification code
// Uses constant cost factor for consistent timing
func HashCode(code string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(code), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyCode compares a plaintext code against a bcrypt hash
// Uses constant-time comparison internally (bcrypt.CompareHashAndPassword)
// to prevent timing attacks
func VerifyCode(hashedCode, plainCode string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(plainCode))
	return err == nil
}

// MustHashCode is a helper for testing that panics on error
// Should NOT be used in production code
func MustHashCode(code string) string {
	hash, err := HashCode(code)
	if err != nil {
		panic(err)
	}
	return hash
}
