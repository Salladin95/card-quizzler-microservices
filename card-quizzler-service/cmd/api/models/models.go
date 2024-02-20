package models

import (
	"github.com/google/uuid"
	"time"
)

type Term struct {
	ID          uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ModuleID    uuid.UUID `json:"moduleID"`
}

type Module struct {
	ID        uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title     string    `gorm:"unique;" json:"title"`
	UserID    string    `gorm:"column:userId;" json:"userId"`
	Folders   []Folder  `gorm:"many2many:folder_modules;" json:"folders,omitempty"`
	Terms     []Term    `json:"terms" gorm:"foreignKey:ModuleID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Folder struct {
	ID        uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title     string    `gorm:"unique;" json:"title"`
	UserID    string    `gorm:"column:userId;" json:"userId"`
	Modules   []Module  `gorm:"many2many:folder_modules;" json:"modules"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
