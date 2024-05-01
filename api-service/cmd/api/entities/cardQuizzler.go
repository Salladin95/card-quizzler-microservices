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
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	UserID      string    `json:"userID"`
	Folders     []Folder  `json:"folders,omitempty"`
	Terms       []Term    `json:"terms"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsOpen      bool      `json:"isOpen" gorm:"default:false;column:is_open"`
	CopiesCount int       `json:"copiesCount" gorm:"default:0;column:copies_count"`
}

type Folder struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	UserID      string    `json:"userID"`
	Modules     []Module  `json:"modules"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsOpen      bool      `json:"isOpen" gorm:"default:false;column:is_open"`
	CopiesCount int       `json:"copiesCount" gorm:"default:0;column:copies_count"`
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
	Title  string          `json:"title" validate:"required"`
	Terms  []CreateTermDto `json:"terms" validate:"required"`
	IsOpen bool            `json:"isOpen" validate:"omitempty"`
}

type UpdateModuleDto struct {
	Title        string          `json:"title" validate:"omitempty"`
	NewTerms     []CreateTermDto `json:"newTerms" validate:"omitempty"`
	UpdatedTerms []Term          `json:"updatedTerms" validate:"omitempty"`
	IsOpen       bool            `json:"isOpen" validate:"omitempty"`
}

type CreateFolderDto struct {
	Title  string `json:"title" validate:"required"`
	IsOpen bool   `json:"isOpen" validate:"omitempty"`
}

type UpdateFolderDto struct {
	Title  string `json:"title" validate:"omitempty"`
	IsOpen bool   `json:"isOpen" validate:"omitempty"`
}

type UpdateTermDto struct {
	ModuleID    string `json:"moduleID" validate:"required"`
	Title       string `json:"title" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
}
