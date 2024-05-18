package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/google/uuid"
)

type CreateModuleDto struct {
	Title  string          `json:"title" validate:"required"`
	UserID string          `json:"userID" validate:"required"`
	Terms  []CreateTermDto `json:"terms" validate:"required"`
	SecureAccess
}

func (dto *CreateModuleDto) Verify() error {
	if err := lib.Verify(dto); err == nil {
		return dto.SecureAccess.Verify()
	}
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

	var psd string

	if dto.Access == models.AccessPassword {
		psd, err = lib.HashPassword(dto.Password)
		if err != nil {
			return module, err
		}
	}

	module = models.Module{
		ID:         id,
		Title:      dto.Title,
		UserID:     dto.UserID,
		AuthorID:   dto.UserID,
		Access:     dto.Access,
		Password:   psd,
		Folders:    []models.Folder{},
		OriginalID: id,
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
	Title        string            `json:"title" validate:"omitempty"`
	NewTerms     []CreateTermDto   `json:"newTerms" validate:"omitempty"`
	UpdatedTerms []models.Term     `json:"updatedTerms" validate:"omitempty"`
	Access       models.AccessType `json:"access" validate:"omitempty"`
	Password     string            `json:"password" validate:"omitempty,min=4"`
}

func (dto *UpdateModuleDto) Verify() error {
	if dto.Access != "" {
		if err := ValidateSecureAccess(
			&SecureAccess{
				Access:   dto.Access,
				Password: dto.Password,
			},
		); err != nil {
			return err
		}
		return lib.Verify(dto)
	}
	return lib.Verify(dto)
}

// JoinTerms parses newTerms and joins them with updatedTerms
func (dto *UpdateModuleDto) JoinTerms(moduleID uuid.UUID) ([]models.Term, error) {
	var newTerms []models.Term
	if err := dto.Verify(); err != nil {
		return newTerms, err
	}

	for _, v := range dto.NewTerms {
		model, err := v.ToModel(moduleID)
		if err != nil {
			return newTerms, err
		}
		newTerms = append(newTerms, model)
	}

	newTerms = append(newTerms, dto.UpdatedTerms...)

	return newTerms, nil
}

type UpdateTermDto struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	ModuleID    uuid.UUID `json:"moduleID" validate:"required"`
	Title       string    `json:"title" validate:"omitempty"`
	Description string    `json:"description" validate:"omitempty"`
}

func (dto *UpdateTermDto) Verify() error {
	return lib.Verify(dto)
}
