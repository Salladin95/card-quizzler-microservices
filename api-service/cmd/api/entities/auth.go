package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
)

// SignInDto represents the data transfer object for user sign-in requests.
type SignInDto struct {
	Email    string `json:"email" validate:"required,email"`    // Email field with validation rules
	Password string `json:"password" validate:"min=6,required"` // Password field with validation rules
}

// SignUpDto represents the data transfer object for user sign-up requests.
type SignUpDto struct {
	Name     string `json:"name" validate:"required,min=1"`     // Name field with validation rules
	Password string `json:"password" validate:"required,min=6"` // Password field with validation rules
	Email    string `json:"email" validate:"required,email"`    // Email field with validation rules
	Birthday string `json:"birthday" validate:"required,min=1"` // Birthday field with validation rules
}

// TokenPair represents a pair of JWTs: access token and refresh token.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Verify validates the structure and content of the SignInDto.
func (signInDto *SignInDto) Verify() error {
	return lib.Verify(signInDto)
}

func (signInDto *SignInDto) ToAuthPayload() *userService.SignInPayload {
	return &userService.SignInPayload{
		Email:    signInDto.Email,
		Password: signInDto.Password,
	}
}

// Verify validates the structure and content of the SignUpDto.
func (signUpDto *SignUpDto) Verify() error {
	return lib.Verify(signUpDto)
}

func (signUpDto *SignUpDto) ToAuthPayload() *userService.SignUpPayload {
	return &userService.SignUpPayload{
		Email:    signUpDto.Email,
		Password: signUpDto.Password,
		Name:     signUpDto.Name,
		Birthday: signUpDto.Birthday,
	}
}

type ResponseWithToken struct {
	AccessToken string `json:"accessToken"`
}
