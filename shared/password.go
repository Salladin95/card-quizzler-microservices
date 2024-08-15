package lib

import (
	"github.com/Salladin95/goErrorHandler"
	"golang.org/x/crypto/bcrypt"
)

// CompareHashAndPassword compares a hashed password with a plaintext password.
// It takes a hashed password and a plaintext password as input and returns an error.
// If the passwords match, the error is nil; otherwise, an error is returned.
func CompareHashAndPassword(hashedPassword string, plaintextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		return goErrorHandler.OperationFailure("compare hash and password", err)
	}
	return nil
}

// HashPassword generates a hashed password using bcrypt with a default cost.
// It takes a plaintext password as input and returns the hashed password as a string.
// An error is returned if the hashing operation fails.
func HashPassword(plaintextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", goErrorHandler.OperationFailure("hash password", err)
	}
	return string(hashedPassword), nil
}
