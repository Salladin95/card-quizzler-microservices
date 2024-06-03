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
	SecureAccess
}

func (dto *CreateFolderDto) Verify() error {
	if err := lib.Verify(dto); err == nil {
		return dto.SecureAccess.Verify()
	}
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

	var psd string

	if dto.Access == models.AccessPassword {
		psd, err = lib.HashPassword(dto.Password)
		if err != nil {
			return createFolderParams, err
		}
	}

	return models.Folder{
		ID:         id,
		UserID:     dto.UserID,
		Title:      dto.Title,
		Modules:    []models.Module{},
		AuthorID:   dto.UserID,
		Access:     dto.Access,
		Password:   psd,
		OriginalID: id,
	}, nil
}

type UpdateFolderDto struct {
	Title    string            `json:"title" validate:"omitempty"`
	Access   models.AccessType `json:"access" validate:"omitempty"`
	Password string            `json:"password" validate:"omitempty,min=4"`
}

func (dto *UpdateFolderDto) Verify() error {
	return lib.Verify(dto)
}
