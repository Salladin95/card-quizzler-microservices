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

func (r *repo) UpdateModule(id uuid.UUID, dto entities.UpdateModuleDto) (models.Module, error) {
	module, err := r.GetModuleByID(id)
	if err != nil {
		return module, err
	}

	parsedDto, err := dto.ToModels(id)
	if err != nil {
		return module, err
	}

	// Update module title if provided in the DTO
	if dto.Title != "" {
		module.Title = dto.Title
	}

	if len(dto.UpdatedTerms) > 0 {
		// update terms
		module.Terms = parsedDto.UpdatedTerms
		//if err := r.db.Save(&parsedDto.UpdatedTerms).Error; err != nil {
		//	return module, err
		//}
	}

	// Append new terms to the module
	for _, term := range parsedDto.NewTerms {
		module.Terms = append(module.Terms, term)
	}

	// Save the updated module
	if err := r.db.Save(&module).Error; err != nil {
		return module, err
	}

	// DeleteTerms
	if err := r.db.Delete(&parsedDto.RemovedTerms).Error; err != nil {
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
		Select(clause.Associations).
		Delete(&module).
		Error

	if err != nil {
		return err
	}

	// delete all associated terms
	return r.db.Delete(&module.Terms).Error
}
