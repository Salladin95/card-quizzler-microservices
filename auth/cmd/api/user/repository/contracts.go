package user

import (
	"cloud.google.com/go/firestore"
	"context"
	fireBaseAuth "firebase.google.com/go/v4/auth"
	userEntities "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/entities"
)

type repository struct {
	dbClient   *firestore.Client
	authClient *fireBaseAuth.Client
}

type Repository interface {
	GetUsers(ctx context.Context) ([]*fireBaseAuth.UserRecord, error)
	GetById(ctx context.Context, uid string) (*fireBaseAuth.UserRecord, error)
	GetByEmail(ctx context.Context, email string) (*fireBaseAuth.UserRecord, error)
	CreateUser(ctx context.Context, createUserDto userEntities.CreateUserDto) (*fireBaseAuth.UserRecord, error)
	UpdateUser(ctx context.Context, uid string, updateUserDto userEntities.UpdateDto) (*fireBaseAuth.UserRecord, error)
	DeleteUser(ctx context.Context, uid string) error
	CompareHashAndPassword(hash string, password string) error
}
