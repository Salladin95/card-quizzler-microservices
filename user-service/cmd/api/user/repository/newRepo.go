package user

import (
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepository(db *mongo.Client, cfg config.MongoCfg) Repository {
	return &repository{db: db, dbCfg: cfg}
}
