package repositories

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
)

func (r *repo) CreateUser(uid string) error {
	return r.db.Create(models.User{ID: uid}).Error
}

func (r *repo) AddModuleToUser(uid string, moduleID uuid.UUID) error {
	return r.addUserAssociation(uid, "Modules", &models.Module{ID: moduleID})
}

func (r *repo) AddFolderToUser(uid string, folderID uuid.UUID) error {
	return r.addUserAssociation(uid, "Folders", &models.Folder{ID: folderID})
}

func (r *repo) addUserAssociation(uid string, associationName string, item interface{}) error {
	return r.db.Model(&models.User{ID: uid}).Association(associationName).Append(item)
}
