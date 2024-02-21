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
	res := r.db.Preload("Modules.Terms").Find(&folders).Where("user_ids @> ARRAY[?]::text[]", uid)
	return folders, res.Error
}

func (r *repo) GetFolderByID(id uuid.UUID) (models.Folder, error) {
	var folder models.Folder
	res := r.db.Preload("Modules.Terms").First(&folder).Where("id", id)
	return folder, res.Error
}

func (r *repo) DeleteFolder(id uuid.UUID) error {
	folder, err := r.GetFolderByID(id)
	if err != nil {
		return err
	}
	// Delete all of a user's has one, has many, and many2many associations
	err = r.db.Select(clause.Associations).Delete(&folder).Error
	return err
}