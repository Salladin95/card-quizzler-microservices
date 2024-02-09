package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
	"github.com/Salladin95/goErrorHandler"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
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

type JsonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type SignInResponse struct {
	AccessToken string `json:"accessToken"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Birthday  string    `json:"birthday"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type JwtUser struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Id    uuid.UUID `json:"id"`
}

// TokenPair represents a pair of JWTs: access token and refresh token.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type JwtUserClaims struct {
	JwtUser
	jwt.RegisteredClaims
}

func GetJwtUserClaims(name string, email string, id uuid.UUID, exp time.Duration) JwtUserClaims {
	return JwtUserClaims{
		JwtUser{name,
			email,
			id,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}

// Verify validates the structure and content of the SignInDto.
func (signInDto *SignInDto) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignInDto structure.
	if err := validate.Struct(signInDto); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

func (signInDto *SignInDto) ToAuthPayload() *userService.SignInPayload {
	return &userService.SignInPayload{
		Email:    signInDto.Email,
		Password: signInDto.Password,
	}
}

// Verify validates the structure and content of the SignUpDto.
func (signUpDto *SignUpDto) Verify() error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignUpDto structure.
	if err := validate.Struct(signUpDto); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

func (signUpDto *SignUpDto) ToAuthPayload() *userService.SignUpPayload {
	return &userService.SignUpPayload{
		Email:    signUpDto.Email,
		Password: signUpDto.Password,
		Name:     signUpDto.Name,
		Birthday: signUpDto.Birthday,
	}
}

func UnmarshalUser(data []byte) (*UserResponse, error) {
	var uRes UserResponse
	err := lib.UnmarshalData(data, &uRes)
	if err != nil {
		return nil, err
	}
	return &uRes, nil
}

type LogMessage struct {
	Message string `json:"message"`
}