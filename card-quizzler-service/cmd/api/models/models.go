package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID      string   `gorm:"primary_key;unique;" json:"id"`
	Folders []Folder `gorm:"foreignKey:UserID" json:"folders,omitempty"`
	Modules []Module `gorm:"foreignKey:UserID" json:"modules,omitempty"`
}

type Term struct {
	ID                   uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	ModuleID             uuid.UUID `json:"moduleID"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	NegativeAnswerStreak int       `gorm:"column:negative_answer_streak" json:"negativeAnswerStreak"`
	PositiveAnswerStreak int       `gorm:"column:positive_answer_streak" json:"negativeAnswerStreak"`
	IsDifficult          bool      `json:"isDifficult"`
}

// AccessType represents the type of access for the Module.
type AccessType string

// Constants representing the available access types.
const (
	AccessOpen     AccessType = "open"
	AccessOnlyMe   AccessType = "only me"
	AccessPassword AccessType = "password"
)

type Module struct {
	ID          uuid.UUID  `gorm:"primary_key;unique;" json:"id"`
	OriginalID  uuid.UUID  `gorm:"column:original_id"json:"originalID"`
	Title       string     `json:"title"`
	UserID      string     `json:"userID" gorm:"column:user_id"` // user that owns this folder
	AuthorID    string     `json:"authorID"`                     // user that has created this module in the first place
	Folders     []Folder   `gorm:"many2many:module_folders;" json:"folders,omitempty"`
	Terms       []Term     `gorm:"foreignKey:ModuleID;" json:"terms"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Access      AccessType `json:"access" gorm:"default:open;"`
	CopiesCount int        `json:"copiesCount" gorm:"default:0;column:copies_count"`
	Password    string     `json:"password"`
}

type Folder struct {
	ID          uuid.UUID  `gorm:"primary_key;unique;" json:"id"`
	OriginalID  uuid.UUID  `gorm:"column:original_id"json:"originalID"`
	Title       string     `json:"title"`
	UserID      string     `json:"userID" gorm:"column:user_id"` // user that owns this folder
	AuthorID    string     `json:"authorID"`                     // user that has created this folder in the first place
	Modules     []Module   `gorm:"many2many:module_folders;" json:"modules"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Access      AccessType `json:"access" gorm:"default:open;"`
	CopiesCount int        `json:"copiesCount" gorm:"default:0;column:copies_count"`
	Password    string     `json:"password"`
}
