package entities

import (
	"github.com/google/uuid"
	"time"
)

type Term struct {
	ID                   uuid.UUID `json:"id"`
	ModuleID             uuid.UUID `json:"moduleID"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	NegativeAnswerStreak int       `json:"negativeAnswerStreak"`
	PositiveAnswerStreak int       `json:"negativeAnswerStreak"`
	IsDifficult          bool      `json:"isDifficult"`
}

type Module struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	UserID    string    `json:"userID"`
	Folders   []Folder  `json:"folders,omitempty"`
	Terms     []Term    `json:"terms"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Folder struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	UserID    string    `json:"userID"`
	Modules   []Module  `json:"modules"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type resultTerm struct {
	Term
	Answer bool `json:"answer" validate:"required"`
}

type QuizResultDto struct {
	Terms    []resultTerm `json:"terms" validate:"required"`
	ModuleID string       `json:"moduleID" validate:"omitempty"`
}

type CreateTermDto struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CreateModuleDto struct {
	Title string          `json:"title" validate:"required"`
	Terms []CreateTermDto `json:"terms" validate:"required"`
}

type UpdateModuleDto struct {
	Title        string          `json:"title" validate:"omitempty"`
	NewTerms     []CreateTermDto `json:"newTerms" validate:"omitempty"`
	UpdatedTerms []Term          `json:"updatedTerms" validate:"omitempty"`
}

type CreateFolderDto struct {
	Title string `json:"title" validate:"required"`
}

type UpdateFolderDto struct {
	Title string `json:"title" validate:"omitempty"`
}
