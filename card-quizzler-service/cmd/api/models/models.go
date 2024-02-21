package models

import (
	"github.com/google/uuid"
	"time"
)

// BaseModel contains common fields like ID and timestamps
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	ID      string   `gorm:"primary_key;unique;" json:"id"`
	Folders []Folder `gorm:"foreignKey:UserID" json:"folders,omitempty"`
	Modules []Module `gorm:"foreignKey:UserID" json:"modules,omitempty"`
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
	UserID    string    `json:"userID"`
	Folders   []Folder  `gorm:"many2many:folder_modules;" json:"folders,omitempty"`
	Terms     []Term    `gorm:"many2many:module_terms;" json:"terms"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Folder struct {
	ID        uuid.UUID `gorm:"primary_key;unique;" json:"id"`
	Title     string    `gorm:"unique;" json:"title"`
	UserID    string    `json:"userID"`
	Modules   []Module  `gorm:"many2many:folder_modules;" json:"modules"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
