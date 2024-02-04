package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `bson:"id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Password  string    `bson:"password" json:"password"`
	Email     string    `bson:"email" json:"email"`
	Birthday  string    `bson:"birthday" json:"birthday"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
