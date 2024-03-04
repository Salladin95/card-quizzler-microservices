package repositories

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repo struct {
	db     *gorm.DB
	broker rmqtools.MessageBroker
}

type Repository interface {
	GetFoldersByUID(ctx context.Context, uid string) ([]models.Folder, error)
	GetFolderByID(ctx context.Context, id uuid.UUID) (models.Folder, error)
	CreateFolder(ctx context.Context, dto entities.CreateFolderDto) (models.Folder, error)
	UpdateFolder(ctx context.Context, folderID uuid.UUID, dto entities.UpdateFolderDto) (models.Folder, error)
	DeleteFolder(ctx context.Context, id uuid.UUID) error

	DeleteModuleFromFolder(ctx context.Context, folderID uuid.UUID, moduleID uuid.UUID) error
	AddFolderToUser(uid string, folderID uuid.UUID) error
	AddModuleToFolder(ctx context.Context, folderID uuid.UUID, moduleID uuid.UUID) error
	GetModulesByUID(ctx context.Context, uid string) ([]models.Module, error)
	GetDifficultModulesByUID(ctx context.Context, uid string) ([]models.Module, error)
	GetModuleByID(ctx context.Context, id uuid.UUID) (models.Module, error)
	CreateModule(ctx context.Context, dto entities.CreateModuleDto) (models.Module, error)
	UpdateModule(ctx context.Context, id uuid.UUID, dto entities.UpdateModuleDto) (models.Module, error)
	DeleteModule(ctx context.Context, id uuid.UUID) error
	AddModuleToUser(uid string, moduleID uuid.UUID) error
	CreateUser(uid string) error
	UpdateTerms(terms []models.Term) error
}

func NewRepo(db *gorm.DB, broker rmqtools.MessageBroker) Repository {
	return &repo{db: db, broker: broker}
}

// withTransaction executes the provided function within a transaction.
// It begins a transaction, calls the provided function with the transaction,
// and commits the transaction if the function completes successfully.
// If an error occurs during any step of the transaction, it rolls back the transaction
// and returns an error.
func (r *repo) withTransaction(fn func(tx *gorm.DB) error) error {
	// Begin a new transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		// If an error occurs while starting the transaction, return an operation failure error
		return goErrorHandler.OperationFailure("start transaction", tx.Error)
	}
	defer func() {
		// Rollback the transaction if a panic occurs
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Call the provided function with the transaction
	if err := fn(tx); err != nil {
		// If an error occurs during the transaction, rollback the transaction and return the error
		tx.Rollback()
		return err
	}
	// Commit the transaction if no errors occurred
	if err := tx.Commit().Error; err != nil {
		// If an error occurs while committing the transaction, log the error
		goErrorHandler.OperationFailure("commit transaction", err)
	}
	return nil
}
