package user

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/config"
	userEntities "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db    *mongo.Client
	dbCfg config.MongoCfg
}

type Repository interface {
	GetUsers(ctx context.Context) ([]*user.User, error)
	GetById(ctx context.Context, uid string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	CreateUser(ctx context.Context, createUserDto userEntities.SignUpDto) (*user.User, error)
	UpdateUser(ctx context.Context, uid string, updateUserDto userEntities.UpdateUserDto) (*user.User, error)
	DeleteUser(ctx context.Context, uid string) error
}

func NewUserRepository(db *mongo.Client, cfg config.MongoCfg) Repository {
	return &repository{db: db, dbCfg: cfg}
}
