package repositories

import "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"

func (r *repo) CreateUser(uid string) error {
	return r.db.Create(models.User{ID: uid}).Error
}
