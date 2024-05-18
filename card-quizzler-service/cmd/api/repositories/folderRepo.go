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
)

// GetFoldersByTitle retrieves folders by title
func (r *repo) GetFoldersByTitle(payload GetByTitlePayload) ([]models.Folder, error) {
	var folders []models.Folder
	if err := r.db.
		Preload("Modules").
		Where("title LIKE ?", "%"+payload.Title+"%").
		Order(payload.SortBy).
		Scopes(newPaginate(int(payload.Limit), int(payload.Page)).paginatedResult).
		Find(&folders).
		Error; err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	data := fetchedData{
		UserID: payload.Uid,
		Key:    fmt.Sprintf("%d:%d:%s:%s", payload.Limit, payload.Page, payload.SortBy, payload.Title),
		Data:   folders,
	}

	r.pushToQueue(payload.Ctx, constants.FetchedByTitleFoldersKey, data)
	return folders, nil
}

// GetOpenFolders retrieves folders where isOpen=true
func (r *repo) GetOpenFolders(payload GetByUIDPayload) ([]models.Folder, error) {
	var folders []models.Folder
	if err := r.db.
		Preload("Modules.Terms").
		Where("access = ?", "open").
		Order(payload.SortBy).
		Scopes(newPaginate(int(payload.Limit), int(payload.Page)).paginatedResult).
		Find(&folders).
		Error; err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	data := fetchedData{
		UserID: payload.Uid,
		Key:    fmt.Sprintf("%d:%d:%s", payload.Limit, payload.Page, payload.SortBy),
		Data:   folders,
	}

	r.pushToQueue(payload.Ctx, constants.FetchedFoldersKey, data)
	return folders, nil
}

// CreateFolder creates a new folder in the database using the provided DTO.
func (r *repo) CreateFolder(ctx context.Context, dto entities.CreateFolderDto) (models.Folder, error) {
	// Convert the DTO to a model
	folder, err := dto.ToModel()
	if err != nil {
		// If an error occurs during conversion, return the folder and the error
		return folder, err
	}

	lib.LogInfo(
		"Before mutations",
		slog.String("service", "FolderRepo"),
		slog.String("method", "CreateFolder"),
		slog.Any("dto", dto),
		slog.Any("module", folder),
	)

	// Attempt to create the folder in the database
	if err := r.db.Create(&folder).Error; err != nil {
		// If an error occurs during creation, return the folder and an operation failure error
		return folder, goErrorHandler.OperationFailure("create folder", err)
	}
	r.pushToQueue(ctx, constants.CreatedFolderKey, folder)
	lib.LogInfo(
		"After mutations",
		slog.String("service", "FolderRepo"),
		slog.String("method", "CreateFolder"),
		slog.Any("dto", dto),
		slog.Any("module", folder),
	)
	// If creation is successful, return the created folder and a nil error
	return folder, nil
}

// UpdateFolder updates a folder in the database with the given ID using the provided DTO.
func (r *repo) UpdateFolder(payload UpdateFolderPayload) (models.Folder, error) {
	// Declare a variable to hold the folder
	var folder models.Folder
	// Execute the provided function within a transaction
	return folder, r.withTransaction(func(tx *gorm.DB) error {
		// Retrieve the folder by ID from the database, preloading its associated modules and terms
		if err := tx.
			Preload("Modules.Terms").
			First(&folder, payload.FolderID).
			Error; err != nil {
			// If the folder is not found, return a not found error
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		lib.LogInfo(
			"Before mutations",
			slog.String("service", "FolderRepo"),
			slog.String("method", "UpdateFolder"),
			slog.Any("payload", payload),
			slog.Any("folder", folder),
		)

		if payload.Dto.Title != "" {
			folder.Title = payload.Dto.Title
		}

		switch payload.Dto.Access {
		case models.AccessOnlyMe, models.AccessOpen:
			folder.Password = ""
			folder.Access = payload.Dto.Access
			break
		case models.AccessPassword:
			folder.Access = payload.Dto.Access
			break
		}

		if folder.Access == models.AccessPassword && payload.Dto.Password != "" {
			psd, err := lib.HashPassword(payload.Dto.Password)
			if err != nil {
				return err
			}
			folder.Password = psd
		}

		// Save the updated folder in the database
		if err := tx.Save(&folder).Error; err != nil {
			// If an error occurs while updating the folder, return an operation failure error
			return goErrorHandler.OperationFailure("update folder", err)
		}
		r.pushToQueue(payload.Ctx, constants.MutatedFolderKey, folder)

		lib.LogInfo(
			"After mutations",
			slog.String("service", "FolderRepo"),
			slog.String("method", "UpdateFolder"),
			slog.Any("payload", payload),
			slog.Any("folder", folder),
		)

		// If no errors occurred, return nil
		return nil
	})
}

// GetFolderByID retrieves a folder from the database by its ID.
func (r *repo) GetFolderByID(ctx context.Context, id uuid.UUID) (models.Folder, error) {
	// Declare a variable to hold the folder
	var folder models.Folder

	// Retrieve the folder by its ID from the database, preloading its associated modules and terms
	if err := r.db.
		Preload("Modules.Terms").
		First(&folder, id).
		Error; err != nil {
		// If the folder is not found, return a not found error
		return folder, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	r.pushToQueue(ctx, constants.FetchedFolderKey, folder)

	// If the folder is found, return the folder and a nil error
	return folder, nil
}

// DeleteFolder deletes a folder from the database by its ID.
// It begins a transaction, retrieves the folder by ID,
// deletes all of its associated modules and their terms,
// and finally deletes the folder itself from the database.
// If the folder is not found, it returns a not found error.
// If an error occurs during the database operation or transaction execution,
// it returns the underlying error.
func (r *repo) DeleteFolder(ctx context.Context, id uuid.UUID) error {
	// Execute the provided function within a transaction
	return r.withTransaction(func(tx *gorm.DB) error {
		// Declare a variable to hold the folder
		var folder models.Folder

		// Retrieve the folder by its ID from the database, preloading its associated modules and terms
		if err := tx.
			Preload("Modules.Terms").
			First(&folder, id).
			Error; err != nil {
			// If the folder is not found, return the error
			return err
		}

		// Delete all of a folder's has one, has many, and many-to-many associations
		if err := r.db.Select(clause.Associations).Delete(&folder, id).Error; err != nil {
			// If an error occurs while deleting associations, return the error
			return err
		}
		r.pushToQueue(ctx, constants.DeletedFolderKey, folder)

		// If no errors occurred, return nil
		return nil
	})
}

// DeleteModuleFromFolder deletes a module from a folder in the database.
// It removes the association between the specified module and folder.
func (r *repo) DeleteModuleFromFolder(payload FolderModuleAssociation) error {
	// Execute the provided function within a transaction
	return r.withTransaction(func(tx *gorm.DB) error {
		// Declare variables to hold the folder and module
		var folder models.Folder
		var module models.Module

		// Fetch the folder and module within the transaction
		if err := tx.First(&folder, payload.FolderID).Error; err != nil {
			return fmt.Errorf("failed to find folder: %w", err)
		}

		if err := tx.First(&module, payload.ModuleID).Error; err != nil {
			return fmt.Errorf("failed to find module: %w", err)
		}

		// Remove the association between the module and folder
		if err := tx.Model(&folder).Association("Modules").Delete(&module); err != nil {
			return fmt.Errorf("failed to delete module from folder: %w", err)
		}

		// Push to queue after successful deletion
		r.pushToQueue(payload.Ctx, constants.MutatedFolderKey, folder)
		r.pushToQueue(payload.Ctx, constants.MutatedModuleKey, module)

		return nil
	})
}
