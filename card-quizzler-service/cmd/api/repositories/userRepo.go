package repositories

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// CreateUser creates a new user in the database with the given UID.
func (r *repo) CreateUser(uid string) error {
	// Attempt to create a new user with the provided UID in the database
	if err := r.db.Create(models.User{ID: uid}).Error; err != nil {
		return goErrorHandler.OperationFailure("create user", err)
	}
	return nil
}

type fetchedData struct {
	Key    string      `json:"key"`
	UserID string      `json:"userID"`
	Data   interface{} `json:"data"`
}

// GetFoldersByUID retrieves folders associated with a user by their UID from the database
func (r *repo) GetFoldersByUID(payload UidSortPayload) ([]models.Folder, error) {
	var userFolders []models.Folder
	if err := r.db.
		Preload("Modules.Terms").
		Where("user_id = ?", payload.Uid).
		Order(payload.SortBy).
		Scopes(newPaginate(int(payload.Limit), int(payload.Page)).paginatedResult).
		Find(&userFolders).
		Error; err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	data := fetchedData{
		UserID: payload.Uid,
		Key:    fmt.Sprintf("%d:%d:%s", payload.Limit, payload.Page, payload.SortBy),
		Data:   userFolders,
	}

	r.broker.PushToQueue(payload.Ctx, constants.FetchedUserFoldersKey, data)
	return userFolders, nil
}

// GetModulesByUID retrieves modules associated with a user by their UID from the database.
func (r *repo) GetModulesByUID(payload UidSortPayload) ([]models.Module, error) {
	var userModules []models.Module
	if err := r.db.
		Preload("Terms").
		Where("user_id = ?", payload.Uid).
		Order(payload.SortBy).
		Scopes(newPaginate(int(payload.Limit), int(payload.Page)).paginatedResult).
		Find(&userModules).
		Error; err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	data := fetchedData{
		UserID: payload.Uid,
		Key:    fmt.Sprintf("%d:%d:%s", payload.Limit, payload.Page, payload.SortBy),
		Data:   userModules,
	}

	r.broker.PushToQueue(payload.Ctx, constants.FetchedUserModulesKey, data)
	return userModules, nil
}

// GetDifficultModulesByUID retrieves difficult modules associated with a user by their UID from the database.
func (r *repo) GetDifficultModulesByUID(ctx context.Context, uid string) ([]models.Module, error) {
	var difficultTerms []models.Term
	if err := r.db.
		Joins("JOIN modules ON terms.module_id = modules.id").
		Where("modules.user_id = ?", uid).
		Where("terms.is_difficult = ?", true).
		Find(&difficultTerms).
		Error; err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	const maxTerms = 30

	var newModules []models.Module
	numModules := (len(difficultTerms) + (maxTerms - 1)) / maxTerms

	for i := 0; i < numModules; i++ {
		start := i * maxTerms
		end := (i + 1) * maxTerms
		if end > len(difficultTerms) {
			end = len(difficultTerms)
		}
		moduleTerms := difficultTerms[start:end]
		newModule := models.Module{
			ID:     uuid.New(),
			Title:  fmt.Sprintf("Module - %d", len(newModules)+1),
			UserID: uid,
			Terms:  moduleTerms,
		}
		newModules = append(newModules, newModule)
	}

	r.broker.PushToQueue(ctx, constants.FetchedDifficultModulesKey, newModules)
	return newModules, nil
}

// AddModuleToUser associates a module with a user.
// It fetches the module from the database by ID and creates a copy of it.
// The copy is assigned to the user specified by the UID parameter.
// It performs these operations within a transaction to ensure atomicity.
func (r *repo) AddModuleToUser(uid string, moduleID uuid.UUID) error {
	return r.withTransaction(func(tx *gorm.DB) error {
		var module models.Module
		// Fetch the module from the database
		if err := tx.
			Preload("Terms").
			Where("id = ?", moduleID).
			First(&module).
			Error; err != nil {
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}
		// Create a copy of the module
		newModule := copyModule(module)
		newModule.UserID = uid // Set the user ID for the new module

		// Create the new module within the transaction
		if err := tx.Create(&newModule).Error; err != nil {
			return goErrorHandler.OperationFailure("add module to user", err)
		}
		return nil
	})
}

// AddFolderToUser associates a folder with a user.
// It fetches the folder from the database by ID and creates a copy of it.
// The copy is assigned to the user specified by the UID parameter.
// It performs these operations within a transaction to ensure atomicity.
func (r *repo) AddFolderToUser(uid string, folderID uuid.UUID) error {
	return r.withTransaction(func(tx *gorm.DB) error {
		var folder models.Folder
		// Fetch the folder from the database
		if err := tx.
			Preload("Modules.Terms").
			First(&folder).
			Where("id = ?", folderID).Error; err != nil {
			return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
		}
		// Create a copy of the folder
		newFolder := copyFolder(folder)
		newFolder.UserID = uid // Set the user ID for the new folder

		// Create the new folder within the transaction
		if err := tx.Create(&newFolder).Error; err != nil {
			return goErrorHandler.OperationFailure("add folder to user", err)
		}
		return nil
	})
}

// copyModule creates a copy of the provided module.
// It generates a new UUID for the module ID and sets the creation and update timestamps.
// It also creates copies of the terms associated with the module, if any.
func copyModule(src models.Module) models.Module {
	// Generate a new UUID for the module ID
	moduleID := uuid.New()

	// Create a new module with the copied attributes
	return models.Module{
		ID:        moduleID,
		Title:     src.Title,
		UserID:    src.UserID,
		Terms:     copyTerms(src.Terms, moduleID), // Copy associated terms
		CreatedAt: time.Now(),                     // Set creation timestamp
		UpdatedAt: time.Now(),                     // Set update timestamp
	}
}

// copyFolder creates a copy of the provided folder.
// It generates a new UUID for the folder ID and sets the creation and update timestamps.
// It also creates copies of the modules associated with the folder, if any.
func copyFolder(src models.Folder) models.Folder {
	// Generate a new UUID for the folder ID
	folderID := uuid.New()

	// Create a new folder with the copied attributes
	return models.Folder{
		ID:        folderID,
		Title:     src.Title,
		UserID:    src.UserID,
		Modules:   copyModules(src.Modules), // Copy associated modules
		CreatedAt: time.Now(),               // Set creation timestamp
		UpdatedAt: time.Now(),               // Set update timestamp
	}
}

// copyModules creates a copy of the provided array of modules.
func copyModules(src []models.Module) []models.Module {
	var copies []models.Module
	for _, module := range src {
		copies = append(copies, copyModule(module))
	}
	return copies
}

// copyTerms creates a copy of the provided array of terms associated with a module.
func copyTerms(src []models.Term, moduleID uuid.UUID) []models.Term {
	var copies []models.Term
	for _, term := range src {
		copies = append(copies, models.Term{
			ID:          uuid.New(),
			Title:       term.Title,
			Description: term.Description,
			ModuleID:    moduleID,
		})
	}
	return copies
}
