package repositories

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

type Repository interface {
	CreateFolder(dto entities.CreateFolderDto) (models.Folder, error)
	GetFoldersByUID(uid string) ([]models.Folder, error)
	GetFolderByID(id uuid.UUID) (models.Folder, error)
	DeleteFolder(id uuid.UUID) error
	CreateModule(dto entities.CreateModuleDto) (models.Module, error)
	GetModulesByUID(uid string) ([]models.Module, error)
	GetModuleByID(id uuid.UUID) (models.Module, error)
	AddModuleToFolder(folderID uuid.UUID, moduleID uuid.UUID) error
	AddTermToModule(termID uuid.UUID, moduleID uuid.UUID) error
	DeleteModule(id uuid.UUID) error
	CreateUser(uid string) error
	// TODO: ADD METHOD REMOVE MODULE FROM FOLDER
}

func NewRepo(db *gorm.DB) Repository {
	return &repo{db: db}
}
