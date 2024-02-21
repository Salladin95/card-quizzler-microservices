package repositories

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (r *repo) CreateModule(dto entities.CreateModuleDto) (models.Module, error) {
	// Convert DTO to model
	module, err := dto.ToModel()
	if err != nil {
		return module, fmt.Errorf("failed to convert DTO to module: %w", err)
	}

	terms, err := parseCreateTermsPayload(dto.Terms, module)
	if err != nil {
		return module, fmt.Errorf("failed to parse terms from DTO: %w", err)
	}

	if err := r.db.
		Create(&module).
		Model(&module).
		Association("Terms").
		Append(&terms); err != nil {
		return module, fmt.Errorf("failed to create module: %v", err)
	}

	return module, nil
}

func (r *repo) UpdateModule(id uuid.UUID, dto entities.UpdateModuleDto) (models.Module, error) {
	module, err := r.GetModuleByID(id)

	if err != nil {
		return module, err
	}

	modules, err := dto.ToModels()

	if err != nil {
		return module, err
	}

	err = r.db.Model(&module).
		Where("id = ?", id).
		Association("Terms").
		Append(&modules.NewTerms)

	if err != nil {
		return module, err
	}

	err = r.db.Save(&modules.UpdatedTerms).Error
	if err != nil {
		return module, err
	}
	return r.GetModuleByID(id)
}

func (r *repo) GetModulesByUID(uid string) ([]models.Module, error) {
	var userModules []models.Module
	err := r.db.
		Preload("Terms").
		Preload("Users").
		Preload("Folders").
		Joins("JOIN user_modules ON modules.id = user_modules.module_id").
		Where("user_modules.user_id = ?", uid).
		Find(&userModules).
		Error
	return userModules, err
}

func (r *repo) GetModuleByID(id uuid.UUID) (models.Module, error) {
	var module models.Module
	err := r.db.
		Preload("Terms.Modules").
		Where("id = ?", id).
		First(&module).
		Error
	return module, err
}

func (r *repo) AddModuleToFolder(folderID uuid.UUID, moduleID uuid.UUID) error {
	var module models.Module
	if err := r.db.First(&module, moduleID).Error; err != nil {
		return err
	}

	// Create the association between the module and the module
	return r.db.
		Model(&module).
		Association("Folders").
		Append(&models.Folder{ID: folderID})
}

func (r *repo) AddTermToModule(termID uuid.UUID, moduleID uuid.UUID) error {
	var term models.Term
	if err := r.db.First(&term, termID).Error; err != nil {
		return err
	}

	// Create the association between the term and the module
	return r.db.
		Model(&term).
		Association("Modules").
		Append(&models.Module{ID: moduleID})
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
