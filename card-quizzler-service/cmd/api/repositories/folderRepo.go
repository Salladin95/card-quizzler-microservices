package repositories

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (r *repo) CreateFolder(dto entities.CreateFolderDto) (models.Folder, error) {
	folder, err := dto.ToModel()

	if err != nil {
		return folder, err
	}

	createdFolder := r.db.Create(&folder)
	if createdFolder.Error != nil {
		return folder, createdFolder.Error
	}

	return folder, nil
}

func (r *repo) GetFoldersByUID(uid string) ([]models.Folder, error) {
	var folders []models.Folder
	return folders, r.db.
		Preload("Modules.Terms").
		Preload("Users").
		Where("user_id = ?", uid).
		Find(&folders).
		Error
}

func (r *repo) GetFolderByID(id uuid.UUID) (models.Folder, error) {
	var folder models.Folder
	return folder, r.db.Preload("Modules.Terms").First(&folder).Where("id", id).Error
}

func (r *repo) DeleteFolder(id uuid.UUID) error {
	var folder models.Folder
	err := r.db.First(&folder).Where("id", id).Error
	if err != nil {
		return err
	}
	// Delete all of a folder's has one, has many, and many2many associations
	return r.db.Select(clause.Associations).Delete(&folder, id).Error
}

func (r *repo) DeleteModuleFromFolder(folderID uuid.UUID, moduleID uuid.UUID) error {
	return r.db.Model(&models.Folder{ID: folderID}).
		Association("Modules").
		Delete(models.Module{ID: moduleID})
}
