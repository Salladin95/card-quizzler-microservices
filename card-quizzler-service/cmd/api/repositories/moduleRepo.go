package repositories

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// CreateModule creates a new module in the database using the provided DTO.
func (r *repo) CreateModule(ctx context.Context, dto entities.CreateModuleDto) (models.Module, error) {
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

	r.broker.PushToQueue(ctx, constants.CreatedModuleKey, module)

	// Return the created module
	return module, nil
}

// UpdateModule updates a module with the given ID using the provided DTO.
// It replaces module's terms with the terms provided in newTerms.
// If newTerms are provided, the module will contain only these terms, the same applies to updatedTerms.
func (r *repo) UpdateModule(ctx context.Context, id uuid.UUID, dto entities.UpdateModuleDto) (models.Module, error) {
	var module models.Module
	// Parse DTO to models
	parsedDto, err := dto.ToModels(id)
	if err != nil {
		return module, err
	}

	// Determine terms to delete
	termsToDelete := getTermsToDelete(module, parsedDto.NewTerms)

	// Define the function to be executed within the transaction
	if err := r.withTransaction(func(tx *gorm.DB) error {
		// Fetch module within the transaction
		if err := tx.
			Preload("Terms").
			Preload("Folders").
			First(&module, id).Error; err != nil {
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		module.UpdatedAt = time.Now()

		// Update module's terms if new terms are provided
		if len(parsedDto.NewTerms) > 0 {
			module.Terms = parsedDto.NewTerms
		}

		// Update module's title if provided in the DTO
		if dto.Title != "" {
			module.Title = dto.Title
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
	}); err != nil {
		return module, err
	}

	r.broker.PushToQueue(ctx, constants.MutatedModuleKey, module)
	// Execute the update function within a transaction
	return module, nil
}

// UpdateTerms updates given terms
func (r *repo) UpdateTerms(terms []models.Term) error {
	// Define the function to be executed within the transaction
	updateFunc := func(tx *gorm.DB) error {
		// Save the updated terms
		if err := tx.Save(&terms).Error; err != nil {
			return goErrorHandler.OperationFailure("update terms", err)
		}
		return nil
	}

	// Execute the update function within a transaction
	return r.withTransaction(updateFunc)
}

// GetModuleByID retrieves a module from the database by its ID,
func (r *repo) GetModuleByID(ctx context.Context, id uuid.UUID) (models.Module, error) {
	// Declare a variable to hold the retrieved module
	var module models.Module

	// Retrieve the module with the given ID from the database
	if err := r.db.
		Preload("Terms").   // Preload associated terms
		First(&module, id). // Execute query and store result in 'module'
		Error; err != nil { // Check for errors
		// If an error occurred, return a not found error
		return module, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	r.broker.PushToQueue(ctx, constants.FetchedModuleKey, module)
	// If no error occurred, return the retrieved module
	return module, nil
}

// AddModuleToFolder adds a module to a folder within a transaction.
func (r *repo) AddModuleToFolder(ctx context.Context, folderID uuid.UUID, moduleID uuid.UUID) error {
	// Retrieve the module from the database by its ID, preloading its associated terms
	var module models.Module
	// Execute the provided function within a transaction
	if err := r.withTransaction(func(tx *gorm.DB) error {
		if err := r.db.
			Preload("Terms").
			Preload("Folders").
			First(&module, moduleID).
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
	}); err != nil {
		return err
	}
	r.broker.PushToQueue(ctx, constants.MutatedModuleKey, module)
	r.broker.PushToQueue(ctx, constants.MutatedFolderKey, models.Folder{ID: folderID, UserID: module.UserID})
	return nil
}

// DeleteModule deletes a module with the given ID from the database within a transaction.
func (r *repo) DeleteModule(ctx context.Context, id uuid.UUID) error {
	// Declare a variable to hold the module to be deleted
	var module models.Module
	// Execute the provided function within a transaction
	if err := r.withTransaction(func(tx *gorm.DB) error {

		// Retrieve the module from the database by its ID, preloading its associated terms
		if err := tx.
			Preload("Terms").
			Preload("Folders").
			First(&module, id).Error; err != nil {
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
	}); err != nil {
		return err
	}
	r.broker.PushToQueue(ctx, constants.DeletedModuleKey, module)
	return nil
}
