package repositories

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateModule creates a new module in the database using the provided DTO.
func (r *repo) CreateModule(dto entities.CreateModuleDto) (models.Module, error) {
	// Convert DTO to model
	module, err := dto.ToModel()
	if err != nil {
		// Return a bad request error if failed to convert DTO to module
		return module, goErrorHandler.NewError(
			goErrorHandler.ErrBadRequest,
			fmt.Errorf("failed to convert DTO to module: %w", err),
		)
	}

	// Parse terms payload from DTO
	terms, err := parseCreateTermsPayload(dto.Terms, module.ID)
	if err != nil {
		// Return a bad request error if failed to parse terms from DTO
		return module, goErrorHandler.NewError(
			goErrorHandler.ErrBadRequest,
			fmt.Errorf("failed to parse terms from DTO: %w", err),
		)
	}

	// Create module and associate terms
	if err := r.db.
		Create(&module).
		Model(&module).
		Association("Terms").
		Append(&terms); err != nil {
		// Return an operation failure error if failed to create module or associate terms
		return module, goErrorHandler.OperationFailure("create module", err)
	}

	// Return the created module
	return module, nil
}

// UpdateModule updates a module with the given ID using the provided DTO.
// It replaces module's terms with the terms provided in newTerms.
// If newTerms are provided, the module will contain only these terms, the same applies to updatedTerms.
func (r *repo) UpdateModule(id uuid.UUID, dto entities.UpdateModuleDto) (models.Module, error) {
	var module models.Module
	// Parse DTO to models
	parsedDto, err := dto.ToModels(id)
	if err != nil {
		return module, err
	}

	// Determine terms to delete
	termsToDelete := getTermsToDelete(module, parsedDto.NewTerms)

	// Define the function to be executed within the transaction
	updateFunc := func(tx *gorm.DB) error {
		// Fetch module within the transaction
		if err := tx.Preload("Terms").Where("id = ?", id).First(&module).Error; err != nil {
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		// Update module's terms if new terms are provided
		if len(parsedDto.NewTerms) > 0 {
			module.Terms = parsedDto.NewTerms
		}

		// Update module's title if provided in the DTO
		if dto.Title != "" {
			if err := tx.Save(&module).Error; err != nil {
				return goErrorHandler.OperationFailure("update module", err)
			}
		}

		// Save the updated terms
		if err := tx.Save(&module.Terms).Error; err != nil {
			return goErrorHandler.OperationFailure("update terms", err)
		}

		// Delete terms
		if len(termsToDelete) > 0 {
			if err := tx.Delete(termsToDelete).Error; err != nil {
				return goErrorHandler.OperationFailure("delete terms", err)
			}
		}

		return nil
	}

	// Execute the update function within a transaction
	return module, r.withTransaction(updateFunc)
}

// GetModuleByID retrieves a module from the database by its ID,
func (r *repo) GetModuleByID(id uuid.UUID) (models.Module, error) {
	// Declare a variable to hold the retrieved module
	var module models.Module

	// Retrieve the module with the given ID from the database
	if err := r.db.
		Preload("Terms").    // Preload associated terms
		Where("id = ?", id). // Filter by module ID
		First(&module).      // Execute query and store result in 'module'
		Error; err != nil {  // Check for errors
		// If an error occurred, return a not found error
		return module, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	// If no error occurred, return the retrieved module
	return module, nil
}

// AddModuleToFolder adds a module to a folder within a transaction.
func (r *repo) AddModuleToFolder(folderID uuid.UUID, moduleID uuid.UUID) error {
	// Execute the provided function within a transaction
	return r.withTransaction(func(tx *gorm.DB) error {
		// Retrieve the module from the database by its ID, preloading its associated terms
		var module models.Module
		if err := r.db.
			Preload("Terms").
			Where("id = ?", moduleID).
			First(&module).
			Error; err != nil {
			// If the module is not found, return a not found error
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		// Create the association between the module and the folder within the transaction
		if err := tx.
			Model(&module).
			Association("Folders").
			Append(&models.Folder{ID: folderID}); err != nil {
			// If an error occurs while creating the association, return a not found error
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}
		// If no errors occurred, return nil
		return nil
	})
}

// DeleteModule deletes a module with the given ID from the database within a transaction.
func (r *repo) DeleteModule(id uuid.UUID) error {
	// Execute the provided function within a transaction
	return r.withTransaction(func(tx *gorm.DB) error {
		// Declare a variable to hold the module to be deleted
		var module models.Module

		// Retrieve the module from the database by its ID, preloading its associated terms
		if err := tx.Preload("Terms").Where("id = ?", id).First(&module).Error; err != nil {
			// If the module is not found, return a not found error
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		// Delete the module and its associated terms from the database
		if err := tx.Select("Terms", clause.Associations).Delete(&module).Error; err != nil {
			// If an error occurs while deleting the module, return an operation failure error
			return goErrorHandler.OperationFailure("delete module", err)
		}

		// If no errors occurred, return nil
		return nil
	})
}
