package repositories

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
	"time"
)

// GetModulesByTitle retrieves folders by title
func (r *repo) GetModulesByTitle(payload GetByTitlePayload) ([]models.Module, error) {
	var modules []models.Module
	if err := r.db.
		Where("title LIKE ?", "%"+payload.Title+"%").
		Order(payload.SortBy).
		Scopes(newPaginate(int(payload.Limit), int(payload.Page)).paginatedResult).
		Find(&modules).
		Error; err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	data := fetchedData{
		UserID: payload.Uid,
		Key:    fmt.Sprintf("%d:%d:%s:%s", payload.Limit, payload.Page, payload.SortBy, payload.Title),
		Data:   modules,
	}

	r.pushToQueue(payload.Ctx, constants.FetchedByTitleModulesKey, data)
	return modules, nil
}

// GetOpenModules retrieves modules where isOpen=true
func (r *repo) GetOpenModules(payload GetByUIDPayload) ([]models.Module, error) {
	var modules []models.Module
	if err := r.db.
		Preload("Terms").
		Where("access = ?", "open").
		Order(payload.SortBy).
		Scopes(newPaginate(int(payload.Limit), int(payload.Page)).paginatedResult).
		Find(&modules).
		Error; err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	data := fetchedData{
		UserID: payload.Uid,
		Key:    fmt.Sprintf("%d:%d:%s", payload.Limit, payload.Page, payload.SortBy),
		Data:   modules,
	}

	r.pushToQueue(payload.Ctx, constants.FetchedModulesKey, data)
	return modules, nil
}

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

	lib.LogInfo(
		"Before mutations",
		slog.String("service", "ModuleRepo"),
		slog.String("method", "CreateModule"),
		slog.Any("dto", dto),
		slog.Any("module", module),
	)

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

	r.pushToQueue(ctx, constants.CreatedModuleKey, module)

	lib.LogInfo(
		"After mutations",
		slog.String("service", "ModuleRepo"),
		slog.String("method", "CreateModule"),
		slog.Any("dto", dto),
		slog.Any("module", module),
	)

	// Return the created module
	return module, nil
}

// UpdateModule updates a module with the given ID using the provided DTO.
// It replaces module's terms with the terms provided in newTerms.
// If newTerms are provided, the module will contain only these terms, the same applies to updatedTerms.
func (r *repo) UpdateModule(payload UpdateModulePayload) (models.Module, error) {
	var module models.Module
	// joinedTerms contains terms which are need to be created also terms which are needs to be updated in the DB
	joinedTerms, err := payload.Dto.JoinTerms(payload.ModuleID)
	if err != nil {
		return module, err
	}

	// Define the function to be executed within the transaction
	if err := r.withTransaction(func(tx *gorm.DB) error {
		// Fetch module within the transaction
		if err := tx.
			Preload("Terms").
			Preload("Folders").
			First(&module, payload.ModuleID).Error; err != nil {
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		lib.LogInfo(
			"Before mutations",
			slog.String("service", "ModuleRepo"),
			slog.String("method", "UpdateModule"),
			slog.Any("payload", payload),
			slog.Any("module", module),
		)

		module.UpdatedAt = time.Now()

		// Update module's title if provided in the DTO
		if payload.Dto.Title != "" {
			module.Title = payload.Dto.Title
		}

		switch payload.Dto.Access {
		case models.AccessOnlyMe, models.AccessOpen:
			module.Password = ""
			module.Access = payload.Dto.Access
			break
		case models.AccessPassword:
			module.Access = payload.Dto.Access
			break
		}

		if module.Access == models.AccessPassword && payload.Dto.Password != "" {
			psd, err := lib.HashPassword(payload.Dto.Password)
			if err != nil {
				return err
			}
			module.Password = psd
		}

		var termsToDelete []models.Term
		if len(joinedTerms) > 0 {
			// Determine terms to delete
			termsToDelete = getTermsToDelete(module, joinedTerms)

			// replace module's terms
			module.Terms = joinedTerms
		}

		if err := tx.Save(&module).Error; err != nil {
			return goErrorHandler.OperationFailure("update module", err)
		}

		if len(joinedTerms) > 0 {
			// Save the updated terms
			if err := tx.Save(&module.Terms).Error; err != nil {
				return goErrorHandler.OperationFailure("update terms", err)
			}
		}

		// Delete terms
		if len(termsToDelete) > 0 {
			if err := tx.Delete(&termsToDelete).Error; err != nil {
				return goErrorHandler.OperationFailure("delete terms", err)
			}
		}

		lib.LogInfo(
			"After mutations",
			slog.String("service", "ModuleRepo"),
			slog.String("method", "UpdateModule"),
			slog.Any("payload", payload),
			slog.Any("module", module),
		)

		return nil
	}); err != nil {
		return module, err
	}

	r.pushToQueue(payload.Ctx, constants.MutatedModuleKey, module)
	// Execute the update function within a transaction
	return module, nil
}

// UpdateTerms updates given terms
func (r *repo) UpdateTerms(ctx context.Context, terms []models.Term) error {
	// Define the function to be executed within the transaction
	return r.withTransaction(func(tx *gorm.DB) error {
		// Save the updated terms
		if err := tx.Save(&terms).Error; err != nil {
			return goErrorHandler.OperationFailure("update terms", err)
		}

		modulesIDSMap := make(map[uuid.UUID]bool)
		for _, v := range terms {
			modulesIDSMap[v.ModuleID] = true
		}

		ids := make([]uuid.UUID, 0, len(modulesIDSMap))
		for id := range modulesIDSMap {
			ids = append(ids, id)
		}

		var modules []models.Module
		if err := tx.Find(&modules, ids).Error; err != nil {
			return err
		}

		for _, module := range modules {
			r.pushToQueue(ctx, constants.MutatedModuleKey, module)
		}

		return nil
	})
}

// UpdateTerm updates given term
func (r *repo) UpdateTerm(ctx context.Context, updateTermDTO entities.UpdateTermDto) error {
	// Define the function to be executed within the transaction
	return r.withTransaction(func(tx *gorm.DB) error {
		var term models.Term
		if err := tx.Find(&term, updateTermDTO.Id).Error; err != nil {
			return goErrorHandler.OperationFailure("get term", err)
		}

		if updateTermDTO.Title != "" {
			term.Title = updateTermDTO.Title
		}

		if updateTermDTO.Description != "" {
			term.Description = updateTermDTO.Description
		}

		// Save the updated term
		if err := tx.Save(&term).Error; err != nil {
			return goErrorHandler.OperationFailure("update term", err)
		}

		var module models.Module
		if err := tx.Find(&module, term.ModuleID).Error; err != nil {
			return goErrorHandler.OperationFailure("get module", err)
		}
		r.pushToQueue(ctx, constants.MutatedModuleKey, module)
		return nil
	})
}

// GetTerms fetches given terms
func (r *repo) GetTerms(termIDS []uuid.UUID) ([]models.Term, error) {
	var terms []models.Term
	// Define the function to be executed within the transaction
	return terms, r.withTransaction(func(tx *gorm.DB) error {
		// Save the updated terms
		if err := tx.Find(&terms, &termIDS).Error; err != nil {
			return goErrorHandler.OperationFailure("fetch terms", err)
		}
		return nil
	})
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

	r.pushToQueue(ctx, constants.FetchedModuleKey, module)
	// If no error occurred, return the retrieved module
	return module, nil
}

// AddModuleToFolder adds a module to a folder within a transaction.
func (r *repo) AddModuleToFolder(payload FolderModuleAssociation) error {
	// Retrieve the module from the database by its ID, preloading its associated terms
	var module models.Module
	// Execute the provided function within a transaction
	if err := r.withTransaction(func(tx *gorm.DB) error {
		if err := r.db.
			Preload("Terms").
			Preload("Folders").
			First(&module, payload.ModuleID).
			Error; err != nil {
			// If the module is not found, return a not found error
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		// Create the association between the module and the folder within the transaction
		if err := tx.
			Model(&module).
			Association("Folders").
			Append(&models.Folder{ID: payload.FolderID}); err != nil {
			// If an error occurs while creating the association, return a not found error
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		// If no errors occurred, return nil
		return nil
	}); err != nil {
		return err
	}
	r.pushToQueue(payload.Ctx, constants.MutatedModuleKey, module)
	r.pushToQueue(payload.Ctx, constants.MutatedFolderKey, models.Folder{ID: payload.FolderID, UserID: module.UserID})
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
	r.pushToQueue(ctx, constants.DeletedModuleKey, module)
	return nil
}
