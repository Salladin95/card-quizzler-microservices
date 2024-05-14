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

type SecureAccess struct {
	Access   string `json:"access" validate:"required"`
	Password string `json:"password" validate:"omitempty,min=4"`
}

type Module struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Folders     []Folder  `json:"folders,omitempty"`
	Terms       []Term    `json:"terms"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CopiesCount int       `json:"copiesCount" gorm:"default:0;column:copies_count"`
	Access      string    `json:"access" gorm:"default:open;"`
	Password    *string   `json:"password"`
	AuthorID    string    `json:"authorID"` // user that has created this module in the first place
	UserID      string    `json:"userID"`   // user that owns this module
}

type Folder struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Modules     []Module  `json:"modules"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CopiesCount int       `json:"copiesCount" gorm:"default:0;column:copies_count"`
	Access      string    `json:"access" gorm:"default:open;"`
	Password    *string   `json:"password"`
	AuthorID    string    `json:"authorID"` // user that has created this folder in the first place
	UserID      string    `json:"userID"`   // user that owns this folder
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
	SecureAccess
}

type UpdateModuleDto struct {
	Title        string          `json:"title" validate:"omitempty"`
	NewTerms     []CreateTermDto `json:"newTerms" validate:"omitempty"`
	UpdatedTerms []Term          `json:"updatedTerms" validate:"omitempty"`
	SecureAccess
}

type CreateFolderDto struct {
	Title string `json:"title" validate:"required"`
	SecureAccess
}

type UpdateFolderDto struct {
	Title string `json:"title" validate:"omitempty"`
	SecureAccess
}

type UpdateTermDto struct {
	ModuleID    string `json:"moduleID" validate:"required"`
	Title       string `json:"title" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
}
