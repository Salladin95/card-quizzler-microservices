package user

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	userEntities "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/entities"
	userModel "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/model"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"time"
)

const CollectionPath = "auth-user"

func (repo *repository) getCollection() *firestore.CollectionRef {
	return repo.dbClient.Collection(CollectionPath)
}

// GetUsers retrieves all active users (not deleted)
func (repo *repository) GetUsers(ctx context.Context) ([]*userModel.User, error) {
	users, err := repo.getCollection().Where("deletedAt", "==", nil).Documents(ctx).GetAll()
	if err != nil {
		return nil, goErrorHandler.OperationFailure("fetching documents", err)
	}

	var userData []*userModel.User
	for _, u := range users {
		user, err := firestoreToUser(u)
		if err != nil {
			return nil, goErrorHandler.OperationFailure("convert firestore userModel document to User", err)
		}
		userData = append(userData, user)
	}

	return userData, nil
}

// GetById retrieves a user by ID, checking if the user is deleted
func (repo *repository) GetById(ctx context.Context, id string) (*userModel.User, error) {
	doc, err := repo.getCollection().Doc(id).Get(ctx)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("fetching document", err)
	}

	user, err := firestoreToUser(doc)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("convert firestore userModel document to User", err)
	}

	// Check if the user is deleted
	if deletedAt, ok := doc.Data()["deletedAt"]; ok && deletedAt != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, fmt.Errorf("user not found"))
	}

	return user, nil
}

// GetByEmail retrieves a user by email, checking if the user is deleted
func (repo *repository) GetByEmail(ctx context.Context, email string) (*userModel.User, error) {
	iter := repo.getCollection().Where("email", "==", email).Where("deletedAt", "==", nil).Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, fmt.Errorf("document not found - %v", err))
		}
		return nil, goErrorHandler.OperationFailure("retrieve document", err)
	}

	user, err := firestoreToUser(doc)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("convert firestore userModel document to User", err)
	}

	return user, nil
}

// CreateUser creates a new user in Firestore
func (repo *repository) CreateUser(ctx context.Context, dto userEntities.SignUpDto) error {
	// Check if user with the same email already exists
	_, err := repo.GetByEmail(ctx, dto.Email)
	if err == nil {
		return goErrorHandler.NewError(goErrorHandler.ErrBadRequest, fmt.Errorf("user with email - %s already exists", dto.Email))
	}

	psd, err := repo.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	// Generate a unique ID for the user document
	userID := uuid.New().String()

	// Set user data in Firestore
	_, err = repo.getCollection().Doc(userID).Set(ctx, map[string]interface{}{
		"name":      dto.Name,
		"password":  psd,
		"birthday":  dto.Birthday,
		"email":     dto.Email,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
		"deletedAt": nil,
	})

	if err != nil {
		return goErrorHandler.OperationFailure("create user", err)
	}

	return nil
}

// UpdateUser updates a user in Firestore
func (repo *repository) UpdateUser(ctx context.Context, id string, dto userEntities.UpdateDto) error {
	// Check if the user exists
	_, err := repo.GetById(ctx, id)
	if err != nil {
		return err
	}

	// Update user data in Firestore
	updateData := make(map[string]interface{})
	updateData["updatedAt"] = time.Now()

	// Helper function to update a field if it's non-empty
	updateField := func(field string, value interface{}) {
		if value != "" {
			updateData[field] = value
		}
	}

	updateField("name", dto.Name)
	updateField("email", dto.Email)
	updateField("birthday", dto.Birthday)

	if dto.Password != "" {
		psd, err := repo.HashPassword(dto.Password)
		if err != nil {
			return err
		}
		updateField("password", psd)
	}

	_, err = repo.getCollection().Doc(id).Set(ctx, updateData, firestore.MergeAll)
	if err != nil {
		return goErrorHandler.OperationFailure("update user", err)
	}
	return nil
}

// DeleteUser soft deletes a user in Firestore by setting "deletedAt"
func (repo *repository) DeleteUser(ctx context.Context, id string) error {
	// Check if the user exists
	_, err := repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	// Soft delete by setting "deletedAt" field
	_, err = repo.getCollection().Doc(id).Set(ctx, map[string]interface{}{
		"updatedAt": time.Now(),
		"deletedAt": time.Now(),
	}, firestore.MergeAll)

	if err != nil {
		return goErrorHandler.OperationFailure("delete user", err)
	}
	return nil
}
