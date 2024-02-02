package user

import (
	"cloud.google.com/go/firestore"
	user "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/model"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// firestoreToUser converts Firestore document data to a User instance.
func firestoreToUser(doc *firestore.DocumentSnapshot) (*user.User, error) {
	// Convert Firestore document data to a user.User instance
	var userData user.User
	if err := doc.DataTo(&userData); err != nil {
		// Handle the error if conversion fails
		return nil, goErrorHandler.OperationFailure("convert firestore document to User", err)
	}

	// Convert Firestore document ID to UUID
	id, err := uuid.Parse(doc.Ref.ID)
	if err != nil {
		// Handle the error if parsing UUID fails
		return nil, goErrorHandler.OperationFailure("parse UUID", err)
	}

	// Set the ID field of the user.User instance
	userData.ID = id

	// Return the user.User instance
	return &userData, nil
}

func (repo *repository) HashPassword(p string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", goErrorHandler.OperationFailure("hash password", err)
	}
	return string(hashedPassword), err
}

func (repo *repository) CompareHashAndPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
