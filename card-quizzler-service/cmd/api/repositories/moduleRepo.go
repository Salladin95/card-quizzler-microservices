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

	terms, err := parseCreateTermsPayload(dto.Terms, module.ID)
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

// UpdateModule updates a module with the given ID using the provided DTO.
// It replaces module's terms with the terms provided in newTerms.
// If newTerms are provided, the module will contain only these terms, the same applies to updatedTerms.
func (r *repo) UpdateModule(id uuid.UUID, dto entities.UpdateModuleDto) (models.Module, error) {
	// Retrieve the module by ID
	module, err := r.GetModuleByID(id)
	if err != nil {
		return module, err
	}

	// Parse DTO to models
	parsedDto, err := dto.ToModels(id)
	if err != nil {
		return module, err
	}

	// Determine terms to delete
	termsToDelete := getTermsToDelete(module, parsedDto.NewTerms)

	// Begin a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return models.Module{}, tx.Error
	}

	// Defer rollback if the transaction encounters an error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update module's terms if new terms are provided
	if len(parsedDto.NewTerms) > 0 {
		module.Terms = parsedDto.NewTerms
	}

	// Update module's title if provided in the DTO
	if dto.Title != "" {
		// Save the updated module title
		if err := tx.Save(&module).Error; err != nil {
			tx.Rollback()
			return module, err
		}
	}

	// Save the updated terms
	if err := tx.Save(&module.Terms).Error; err != nil {
		tx.Rollback()
		return module, err
	}

	// Delete terms
	if len(termsToDelete) > 0 {
		if err := tx.Delete(termsToDelete).Error; err != nil {
			tx.Rollback()
			return module, err
		}
	}

	// Commit the transaction if all operations succeed
	if err := tx.Commit().Error; err != nil {
		return module, err
	}

	// Return the updated module
	return module, nil
}

func (r *repo) GetModuleByID(id uuid.UUID) (models.Module, error) {
	var module models.Module
	err := r.db.
		Preload("Terms").
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

func (r *repo) DeleteModule(id uuid.UUID) error {
	module, err := r.GetModuleByID(id)
	if err != nil {
		return err
	}

	err = r.db.
		Select("Terms", clause.Associations).
		Delete(&module).
		Error

	return err
}
