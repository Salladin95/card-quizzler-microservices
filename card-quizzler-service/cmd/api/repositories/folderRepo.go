package repositories

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateFolder creates a new folder in the database using the provided DTO.
func (r *repo) CreateFolder(ctx context.Context, dto entities.CreateFolderDto) (models.Folder, error) {
	// Convert the DTO to a model
	folder, err := dto.ToModel()
	if err != nil {
		// If an error occurs during conversion, return the folder and the error
		return folder, err
	}
	// Attempt to create the folder in the database
	if err := r.db.Create(&folder).Error; err != nil {
		// If an error occurs during creation, return the folder and an operation failure error
		return folder, goErrorHandler.OperationFailure("create folder", err)
	}
	// If creation is successful, return the created folder and a nil error
	return folder, nil
}

// UpdateFolder updates a folder in the database with the given ID using the provided DTO.
func (r *repo) UpdateFolder(ctx context.Context, folderID uuid.UUID, dto entities.UpdateFolderDto) (models.Folder, error) {
	// Declare a variable to hold the folder
	var folder models.Folder

	// Execute the provided function within a transaction
	return folder, r.withTransaction(func(tx *gorm.DB) error {
		// Retrieve the folder by ID from the database, preloading its associated modules and terms
		if err := tx.
			Preload("Modules.Terms").
			First(&folder).
			Where("id", folderID).Error; err != nil {
			// If the folder is not found, return a not found error
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}

		// Update the folder's title with the data from the DTO
		folder.Title = dto.Title

		// Save the updated folder in the database
		if err := tx.Save(&folder).Error; err != nil {

			// If an error occurs while updating the folder, return an operation failure error
			return goErrorHandler.OperationFailure("update folder", err)
		}

		// If no errors occurred, return nil
		return nil
	})
}

// GetFolderByID retrieves a folder from the database by its ID.
func (r *repo) GetFolderByID(ctx context.Context, id uuid.UUID) (models.Folder, error) {
	// Declare a variable to hold the folder
	var folder models.Folder

	// Retrieve the folder by its ID from the database, preloading its associated modules and terms
	if err := r.db.Preload("Modules.Terms").First(&folder).Where("id", id).Error; err != nil {
		// If the folder is not found, return a not found error
		return folder, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

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
			First(&folder).
			Where("id", id).Error; err != nil {
			// If the folder is not found, return the error
			return err
		}

		// Delete all of a folder's has one, has many, and many-to-many associations
		if err := r.db.Select(clause.Associations).Delete(&folder, id).Error; err != nil {
			// If an error occurs while deleting associations, return the error
			return err
		}

		// If no errors occurred, return nil
		return nil
	})
}

// DeleteModuleFromFolder deletes a module from a folder in the database.
// It removes the association between the specified module and folder.
func (r *repo) DeleteModuleFromFolder(ctx context.Context, folderID uuid.UUID, moduleID uuid.UUID) error {
	// Delete the association between the specified module and folder
	if err := r.db.Model(&models.Folder{ID: folderID}).
		Association("Modules").
		Delete(models.Module{ID: moduleID}); err != nil {
		// If an error occurs during deletion, return the error
		return goErrorHandler.OperationFailure("delete module from folder", err)
	}
	// If deletion is successful, return nil
	return nil
}
