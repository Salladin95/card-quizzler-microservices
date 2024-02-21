package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
)

type JsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreateFolderDto struct {
	Title  string `json:"title" validate:"required"`
	UserID string `json:"userID" validate:"required"`
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
		ID:    id,
		Users: []models.User{{ID: dto.UserID}},
		Title: dto.Title,
	}, nil
}

type UpdateFolderDto struct {
	Title string `json:"title" validate:"omitempty"`
}

func (dto *UpdateFolderDto) Verify() error {
	return lib.Verify(dto)
}

type CreateModuleDto struct {
	Title  string          `json:"title" validate:"required"`
	UserID string          `json:"userID" validate:"required"`
	Terms  []CreateTermDto `json:"terms" validate:"required"`
}

func (dto *CreateModuleDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateModuleDto) ToModel() (models.Module, error) {
	err := dto.Verify()
	var module models.Module
	if err != nil {
		return module, err
	}

	// Generate UUID for the module
	id := uuid.New()

	module = models.Module{
		ID:    id,
		Title: dto.Title,
	}

	return module, nil
}

type CreateTermDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (dto *CreateTermDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateTermDto) ToModel() (models.Term, error) {
	err := dto.Verify()
	var model models.Term
	if err != nil {
		return model, err
	}

	id := uuid.New()

	return models.Term{
		ID:          id,
		Description: dto.Description,
		Title:       dto.Title,
	}, nil
}

type UpdateTermDto struct {
	ID          string `json:"id" validate:"required"`
	Title       string `json:"title" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
}

func (dto *UpdateTermDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *UpdateTermDto) ToModel() (models.Term, error) {
	err := dto.Verify()
	var model models.Term
	if err != nil {
		return model, err
	}

	id, err := uuid.Parse(dto.ID)

	if err != nil {
		id = uuid.New()
	}

	if err != nil {
		return model, err
	}

	return models.Term{
		ID:          id,
		Description: dto.Description,
		Title:       dto.Title,
	}, nil
}

type UpdateModuleDto struct {
	Title        string          `json:"title" validate:"omitempty"`
	NewTerms     []CreateTermDto `json:"newTerms" validate:"omitempty"`
	UpdatedTerms []UpdateTermDto `json:"updatedTerms" validate:"omitempty"`
}

type parsedUpdateModuleDto struct {
	Title         string `json:"title" validate:"omitempty"`
	NewModels     []models.Term
	UpdatedModels []models.Term
}

func (dto *UpdateModuleDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *UpdateModuleDto) ToModels() (parsedUpdateModuleDto, error) {
	err := dto.Verify()
	var models parsedUpdateModuleDto
	if err != nil {
		return models, err
	}

	models.Title = dto.Title

	for _, v := range dto.NewTerms {
		model, err := v.ToModel()
		if err != nil {
			return models, err
		}
		models.NewModels = append(models.NewModels, model)
	}

	for _, v := range dto.UpdatedTerms {
		model, err := v.ToModel()
		if err != nil {
			return models, err
		}
		models.UpdatedModels = append(models.UpdatedModels, model)
	}

	return models, nil
}
