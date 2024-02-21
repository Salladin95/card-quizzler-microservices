package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID      string   `gorm:"primary_key;unique;" json:"id"`
	Folders []Folder `gorm:"many2many:user_folders;" json:"folders,omitempty"`
	Modules []Module `gorm:"many2many:user_modules;" json:"modules,omitempty"`
}

type Term struct {
	ID          uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Modules     []Module  `gorm:"many2many:module_terms;" json:"modules"`
}

type Module struct {
	ID        uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title     string    `json:"title"`
	Users     []User    `gorm:"many2many:user_modules;" json:"users"`
	Folders   []Folder  `gorm:"many2many:folder_modules;" json:"folders,omitempty"`
	Terms     []Term    `gorm:"many2many:module_terms;" json:"terms"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Folder struct {
	ID        uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title     string    `gorm:"unique;" json:"title"`
	Users     []User    `gorm:"many2many:user_folders;" json:"users"`
	Modules   []Module  `gorm:"many2many:folder_modules;" json:"modules"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
