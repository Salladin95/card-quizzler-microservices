package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
)

func (dto *CreateUserDto) Verify() error {
	return lib.Verify(dto)
}

type CreateFolderDto struct {
	Title  string `json:"title" validate:"required"`
	UserID string `json:"userID" validate:"required"`
	IsOpen bool   `json:"isOpen" validate:"omitempty"`
}

func (dto *CreateFolderDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateFolderDto) ToModel() (models.Folder, error) {
	err := dto.Verify()
	var createFolderParams models.Folder
	if err != nil {
		return createFolderParams, err
	}

	// Generate UUID for the module
	id := uuid.New()

	return models.Folder{
		ID:      id,
		UserID:  dto.UserID,
		Title:   dto.Title,
		Modules: []models.Module{},
		IsOpen:  dto.IsOpen,
	}, nil
}

type UpdateFolderDto struct {
	Title  string `json:"title" validate:"omitempty"`
	IsOpen bool   `json:"isOpen" validate:"omitempty"`
}

func (dto *UpdateFolderDto) Verify() error {
	return lib.Verify(dto)
}
