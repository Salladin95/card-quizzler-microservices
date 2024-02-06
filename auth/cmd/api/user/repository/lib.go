package user

import (
	"github.com/Salladin95/goErrorHandler"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a hashed password using bcrypt with a default cost.
// It takes a plaintext password as input and returns the hashed password as a string.
// An error is returned if the hashing operation fails.
func (repo *repository) HashPassword(plaintextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", goErrorHandler.OperationFailure("hash password", err)
	}
	return string(hashedPassword), nil
}

// CompareHashAndPassword compares a hashed password with a plaintext password.
// It takes a hashed password and a plaintext password as input and returns an error.
// If the passwords match, the error is nil; otherwise, an error is returned.
func (repo *repository) CompareHashAndPassword(hashedPassword string, plaintextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	return err
}

// toObjectId converts a hexadecimal string representation of MongoDB ObjectID to primitive.ObjectID.
// It returns the converted ObjectID and an error.
// If the conversion fails, an error with context is returned.
func toObjectId(id string) (primitive.ObjectID, error) {
	// Convert the hexadecimal string to primitive.ObjectID
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return docID, goErrorHandler.OperationFailure("create objectID", err)
	}
	return docID, nil
}
