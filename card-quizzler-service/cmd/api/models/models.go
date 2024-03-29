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

type Module struct {
	ID          uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title       string    `json:"title"`
	UserID      string    `json:"userID"`
	Folders     []Folder  `gorm:"many2many:module_folders;" json:"folders,omitempty"`
	Terms       []Term    `gorm:"foreignKey:ModuleID;" json:"terms"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsOpen      bool      `json:"isOpen" gorm:"default:false;column:is_open"`
	CopiesCount int       `json:"copiesCount" gorm:"default:0;column:copies_count"`
}

type Folder struct {
	ID          uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title       string    `gorm:"unique;" json:"title"`
	UserID      string    `json:"userID"`
	Modules     []Module  `gorm:"many2many:module_folders;" json:"modules"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsOpen      bool      `json:"isOpen" gorm:"default:false;column:is_open"`
	CopiesCount int       `json:"copiesCount" gorm:"default:0;column:copies_count"`
}
