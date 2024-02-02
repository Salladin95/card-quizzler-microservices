package user

import (
	"cloud.google.com/go/firestore"
	"context"
	userEntities "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/entities"
	userModel "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/model"
)

type repository struct {
	dbClient *firestore.Client
}

type Repository interface {
	GetUsers(ctx context.Context) ([]*userModel.User, error)
	GetById(ctx context.Context, id string) (*userModel.User, error)
	GetByEmail(ctx context.Context, email string) (*userModel.User, error)
	CreateUser(ctx context.Context, createUserDto userEntities.SignUpDto) error
	UpdateUser(ctx context.Context, id string, updateUserDto userEntities.UpdateDto) error
	DeleteUser(ctx context.Context, id string) error
	CompareHashAndPassword(hash string, password string) error
}

func NewUserRepository(dbClient *firestore.Client) Repository {
	return &repository{dbClient: dbClient}
}
