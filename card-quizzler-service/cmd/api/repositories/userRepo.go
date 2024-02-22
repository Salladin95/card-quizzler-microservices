package repositories

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"time"
)

func (r *repo) CreateUser(uid string) error {
	return r.db.Create(models.User{ID: uid}).Error
}

func (r *repo) GetModulesByUID(uid string) ([]models.User, error) {
	var user []models.User
	return user, r.db.
		Preload("Modules.Terms").
		Where("id = ?", uid).
		Find(&user).
		Error
}

func (r *repo) GetFoldersByUID(uid string) ([]models.User, error) {
	var user []models.User
	return user, r.db.
		Preload("Folders.Modules.Terms").
		Where("id = ?", uid).
		Find(&user).
		Error
}

func (r *repo) AddModuleToUser(uid string, moduleID uuid.UUID) error {
	// Get the module by its ID
	module, err := r.GetModuleByID(moduleID)
	if err != nil {
		return fmt.Errorf("failed to get module by ID: %w", err)
	}

	// Create a copy of the module
	newModule := copyModule(module)
	newModule.UserID = uid // Set the user ID for the new module

	return r.db.Create(&newModule).Error
}

func (r *repo) AddFolderToUser(uid string, folderID uuid.UUID) error {
	// Get the folder by its ID
	folder, err := r.GetFolderByID(folderID)
	if err != nil {
		return fmt.Errorf("failed to get folder by ID: %w", err)
	}

	// Create a copy of the folder
	newFolder := copyFolder(folder)
	newFolder.UserID = uid // Set the user ID for the new folder

	return r.db.Create(&newFolder).Error
}

// Function to create a copy of a module
func copyModule(src models.Module) models.Module {
	return models.Module{
		ID:        uuid.New(),
		Title:     src.Title,
		UserID:    src.UserID,
		Terms:     copyTerms(src.Terms),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Function to create a copy of a folder
func copyFolder(src models.Folder) models.Folder {
	return models.Folder{
		ID:        uuid.New(),
		Title:     src.Title,
		UserID:    src.UserID,
		Modules:   copyModules(src.Modules),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Function to create a copy of an array of modules
func copyModules(src []models.Module) []models.Module {
	var copies []models.Module
	for _, module := range src {
		copies = append(copies, copyModule(module))
	}
	return copies
}

// Function to create a copy of an array of terms
func copyTerms(src []models.Term) []models.Term {
	var copies []models.Term
	for _, term := range src {
		copies = append(copies, models.Term{
			ID:          uuid.New(),
			Title:       term.Title,
			Description: term.Description,
			Modules:     term.Modules,
		})
	}
	return copies
}
