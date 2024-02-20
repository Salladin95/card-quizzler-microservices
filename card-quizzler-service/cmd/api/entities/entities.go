package entities

import (
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
	"time"
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
		ID:     id,
		UserID: dto.UserID,
		Title:  dto.Title,
	}, nil
}

type CreateModuleDto struct {
	Title  string `json:"title" validate:"required"`
	UserID string `json:"userID" validate:"required"`
	// TODO: REMOVE OMITEMPTY
	Terms []CreateTermDto `json:"terms" validate:"omitempty"`
}

func mockTerm() CreateTermDto {
	return CreateTermDto{
		ID:          uuid.New().String(),
		Title:       fmt.Sprintf("sth-%v", time.Now()),
		Description: fmt.Sprintf("sth-%v", time.Now()),
	}
}

func mockTerms() []CreateTermDto {
	i := 1
	var terms []CreateTermDto
	for i < 10 {
		terms = append(terms, mockTerm())
		i++
	}
	return terms
}
func (dto *CreateModuleDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateModuleDto) ToModels() (models.Module, []models.Term, error) {
	err := dto.Verify()
	var module models.Module
	var terms []models.Term
	//var module models.Module
	if err != nil {
		return module, terms, err
	}

	// Generate UUID for the module
	id := uuid.New()

	//TODO: REMOVE
	dto.Terms = mockTerms()

	module = models.Module{
		ID:     id,
		UserID: dto.UserID,
		Title:  dto.Title,
	}

	for _, v := range dto.Terms {
		termModel, err := v.ToModel(id)
		if err != nil {
			return models.Module{}, nil, err
		}
		terms = append(terms, termModel)
	}

	return module, terms, nil
}

type CreateTermDto struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (dto *CreateTermDto) Verify() error {
	return lib.Verify(dto)
}

func (dto *CreateTermDto) ToModel(moduleID uuid.UUID) (models.Term, error) {
	err := dto.Verify()
	var model models.Term
	if err != nil {
		return model, err
	}

	id, err := uuid.Parse(dto.ID)

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
