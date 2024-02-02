package user

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	userEntities "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
)

const CollectionPath = "auth-user"

func (repo *repository) getCollection() *firestore.CollectionRef {
	return repo.dbClient.Collection(CollectionPath)
}

// GetUsers retrieves all active users (not deleted)
func (repo *repository) GetUsers(ctx context.Context) ([]*auth.UserRecord, error) {
	iter, err := repo.authClient.GetUsers(ctx, []auth.UserIdentifier{})
	if err != nil {
		return nil, goErrorHandler.OperationFailure("fetching documents", err)
	}
	return iter.Users, nil
}

// GetById retrieves a user by ID, checking if the user is deleted
func (repo *repository) GetById(ctx context.Context, uid string) (*auth.UserRecord, error) {
	u, err := repo.authClient.GetUser(ctx, uid)
	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, fmt.Errorf("error getting user %s: %v\n", uid, err))
	}
	return u, nil
}

// GetByEmail retrieves a user by email, checking if the user is deleted
func (repo *repository) GetByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	u, err := repo.authClient.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, fmt.Errorf("error getting user %s: %v\n", email, err))
	}
	return u, nil
}

// CreateUser creates a new user in Firestore
func (repo *repository) CreateUser(ctx context.Context, dto userEntities.CreateUserDto) (*auth.UserRecord, error) {
	// Check if user with the same email already exists
	_, err := repo.GetByEmail(ctx, dto.Email)
	if err == nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrBadRequest, fmt.Errorf("user with email - %s already exists", dto.Email))
	}

	uid := uuid.New()
	params := (&auth.UserToCreate{}).
		UID(uid.String()).
		Email(dto.Email).
		EmailVerified(false).
		Password(dto.Password).
		DisplayName(dto.Name).
		Disabled(false)

	u, err := repo.authClient.CreateUser(ctx, params)

	if err != nil {
		return nil, goErrorHandler.OperationFailure("create user", err)
	}

	return u, nil
}

// UpdateUser updates a user in Firestore
func (repo *repository) UpdateUser(ctx context.Context, uid string, dto userEntities.UpdateDto) (*auth.UserRecord, error) {
	params := (&auth.UserToUpdate{}).Email(dto.Email).Password(dto.Password).DisplayName(dto.Name)

	u, err := repo.authClient.UpdateUser(ctx, uid, params)

	if err != nil {
		return nil, goErrorHandler.OperationFailure("create user", err)
	}

	return u, nil
}

// DeleteUser soft deletes a user in Firestore by setting "deletedAt"
func (repo *repository) DeleteUser(ctx context.Context, uid string) error {
	err := repo.authClient.DeleteUser(ctx, uid)
	if err != nil {
		return goErrorHandler.NewError(goErrorHandler.ErrNotFound, fmt.Errorf("error deleting user %s: %v\n", uid, err))
	}
	return nil
}
