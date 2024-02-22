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
	ID          uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	ModuleID    uuid.UUID `json:"moduleID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type Module struct {
	ID        uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title     string    `json:"title"`
	UserID    string    `json:"userID"`
	Folders   []Folder  `gorm:"many2many:module_folders;" json:"folders,omitempty"`
	Terms     []Term    `gorm:"foreignKey:ModuleID;" json:"terms"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Folder struct {
	ID        uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title     string    `gorm:"unique;" json:"title"`
	UserID    string    `json:"userID"`
	Modules   []Module  `gorm:"many2many:module_folders;" json:"modules"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
