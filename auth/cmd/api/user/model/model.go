package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID  `json:"id" validate:"uuid"`
	Name      string     `json:"name" validate:"required"`
	Password  string     `json:"password" validate:"required"`
	Email     string     `json:"email" validate:"required,email"`
	Birthday  string     `json:"birthday" validate:"required"`
	CreatedAt time.Time  `json:"createdAt" validate:"required"`
	UpdatedAt time.Time  `json:"updatedAt" validate:"required"`
	DeletedAt *time.Time `json:"deletedAt"`
}
