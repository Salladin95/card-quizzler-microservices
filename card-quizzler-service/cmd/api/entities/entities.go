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
		ID:      id,
		UserID:  dto.UserID,
		Title:   dto.Title,
		Modules: []models.Module{},
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
		ID:      id,
		Title:   dto.Title,
		UserID:  dto.UserID,
		Folders: []models.Folder{},
	}

	return module, nil
}

type CreateTermDto struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (dto *CreateTermDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateTermDto) ToModel(moduleID uuid.UUID) (models.Term, error) {
	var model models.Term
	err := dto.Verify()
	if err != nil {
		return model, err
	}

	id := uuid.New()

	if err != nil {
		return model, err
	}

	return models.Term{
		ID:          id,
		Description: dto.Description,
		Title:       dto.Title,
		ModuleID:    moduleID,
	}, nil
}

type UpdateModuleDto struct {
	Title        string          `json:"title" validate:"omitempty"`
	NewTerms     []CreateTermDto `json:"newTerms" validate:"omitempty"`
	UpdatedTerms []models.Term   `json:"updatedTerms" validate:"omitempty"`
}

type parsedUpdateModuleDto struct {
	Title    string `json:"title" validate:"omitempty"`
	NewTerms []models.Term
}

func (dto *UpdateModuleDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *UpdateModuleDto) ToModels(moduleID uuid.UUID) (parsedUpdateModuleDto, error) {
	var models parsedUpdateModuleDto
	if err := dto.Verify(); err != nil {
		return models, err
	}

	models.Title = dto.Title

	for _, v := range dto.NewTerms {
		model, err := v.ToModel(moduleID)
		if err != nil {
			return models, err
		}
		models.NewTerms = append(models.NewTerms, model)
	}

	models.NewTerms = append(models.NewTerms, dto.UpdatedTerms...)

	return models, nil
}
