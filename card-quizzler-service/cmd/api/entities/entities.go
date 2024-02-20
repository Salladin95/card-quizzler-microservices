package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	dbService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/db/sqlc"
	"github.com/google/uuid"
)

type JsonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type CreateFolderDto struct {
	Title  string `json:"title" validate:"required"`
	UserID string `json:"userID" validate:"required"`
}

func (dto *CreateFolderDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateFolderDto) ToCreateFolderParams() (dbService.CreateFolderParams, error) {
	err := dto.Verify()
	var createFolderParams dbService.CreateFolderParams
	if err != nil {
		return createFolderParams, err
	}

	// Generate UUID for the module
	moduleID := uuid.New().String()

	return dbService.CreateFolderParams{
		ID:     moduleID,
		UserID: dto.UserID,
		Title:  dto.Title,
	}, nil
}

type CreateModuleDto struct {
	Title  string `json:"title" validate:"required"`
	UserID string `json:"userID" validate:"required"`
	//Terms    []CreateTermDto `json:"terms" validate:"required"`
}

func (dto *CreateModuleDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateModuleDto) ToCreateModuleParams() (dbService.CreateModuleParams, error) {
	err := dto.Verify()
	var createModuleParams dbService.CreateModuleParams
	if err != nil {
		return createModuleParams, err
	}

	// Generate UUID for the module
	moduleID := uuid.New().String()

	return dbService.CreateModuleParams{
		ID:     moduleID,
		UserID: dto.UserID,
		Title:  dto.Title,
	}, nil
}

type CreateTermDto struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (dto *CreateTermDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateTermDto) ToCreateTermDto(moduleID string) (dbService.CreateTermParams, error) {
	err := dto.Verify()
	var CreateTermParams dbService.CreateTermParams
	if err != nil {
		return CreateTermParams, err
	}

	return dbService.CreateTermParams{
		ID:          dto.ID,
		Description: dto.Description,
		Title:       dto.Title,
		ModuleID:    moduleID,
	}, nil
}
