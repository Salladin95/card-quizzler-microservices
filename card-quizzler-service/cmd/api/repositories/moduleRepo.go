package repositories

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (r *repo) CreateModule(dto entities.CreateModuleDto) (models.Module, error) {
	module, err := dto.ToModel()

	if err != nil {
		return module, err
	}

	terms, err := parseCreateTermsPayload(dto.Terms, module)

	if err != nil {
		return module, err
	}

	err = r.db.Create(&module).Error
	if err != nil {
		return module, err
	}
	err = r.db.Model(&module).Association("Users").Append(&models.User{ID: dto.UserID})
	if err != nil {
		return module, err
	}
	err = r.db.Model(&module).Association("Terms").Append(&terms)
	if err != nil {
		return module, err
	}
	return module, nil
}

func (r *repo) GetModulesByUID(uid string) ([]models.Module, error) {
	var userModules []models.Module
	res := r.db.
		Preload("Terms").
		Find(&userModules).
		Where("users @> ARRAY[?]::text[]", uid)
	return userModules, res.Error
}

func (r *repo) GetModuleByID(id uuid.UUID) (models.Module, error) {
	var module models.Module
	res := r.db.
		Preload("Terms.Modules").
		Where("id = ?", id).
		First(&module)
	return module, res.Error
}

func (r *repo) AddModuleToFolder(folderID uuid.UUID, moduleID uuid.UUID) error {
	var module models.Module
	if err := r.db.First(&module, moduleID).Error; err != nil {
		return err
	}

	// Create the association between the module and the module
	res := r.db.
		Model(&module).
		Association("Folders").
		Append(&models.Folder{ID: folderID})

	return res
}

func (r *repo) AddTermToModule(termID uuid.UUID, moduleID uuid.UUID) error {
	var term models.Term
	if err := r.db.First(&term, termID).Error; err != nil {
		return err
	}

	// Create the association between the term and the module
	res := r.db.
		Model(&term).
		Association("Modules").
		Append(&models.Module{ID: moduleID})

	return res
}

func (r *repo) DeleteModule(id uuid.UUID) error {
	module, err := r.GetModuleByID(id)
	if err != nil {
		return err
	}
	module.Terms = extractAssociatedTerms(module.Terms)

	fmt.Println(module.Title)

	err = r.db.
		Select(clause.Associations).
		Delete(&module).
		Error

	if err != nil {
		return err
	}

	// delete all associated terms
	return r.db.Delete(&module.Terms).Error
}
