package user

import (
	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type UserRepository struct {
	dbClient *firestore.Client
}

type IUserRepository interface {
	GetUsers() (interface{}, error)
	GetUserById(id uuid.UUID) (interface{}, error)
	GetUserByEmail(email string) (interface{}, error)
	CreateUser(createUserDto interface{}) (interface{}, error)
	UpdateUser(id uuid.UUID, updateUserDto interface{}) (interface{}, error)
	DeleteUser(id uuid.UUID) (interface{}, error)
}

func NewUserRepository(dbClient *firestore.Client) IUserRepository {
	return &UserRepository{dbClient: dbClient}
}
