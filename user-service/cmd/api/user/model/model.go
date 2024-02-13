package user

import (
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User Represents user's mongo document model
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Password  string             `bson:"password" json:"password"`
	Email     string             `bson:"email" json:"email"`
	Birthday  string             `bson:"birthday" json:"birthday"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// ToResponse converts a User object to a response object, omitting the password field
func (u *User) ToResponse() *user.Response {
	return &user.Response{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Birthday:  u.Birthday,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUserResponses(users []*User) []*user.Response {
	userResponses := make([]*user.Response, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}
	return userResponses
}
