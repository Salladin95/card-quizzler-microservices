package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
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

type ResponseWithToken struct {
	AccessToken string `json:"accessToken"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Birthday  string    `json:"birthday"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

type LogMessage struct {
	FromService string `json:"fromService" validate:"required"`
	Message     string `json:"message" validate:"required"`
	Level       string `json:"level" validate:"required"`
	Name        string `json:"name" validate:"omitempty"`
	Method      string `json:"method" validate:"omitempty"`
}

func (log *LogMessage) GenerateLog(message string, level string, method string, name string) LogMessage {
	return LogMessage{
		Level:       level,
		Method:      method,
		FromService: "api-service",
		Message:     message,
		Name:        name,
	}
}

type RequestEmailVerificationDto struct {
	Email string `json:"email" validate:"required"`
}

// Verify validates the structure and content of the SignUpDto.
func (dto *RequestEmailVerificationDto) Verify() error {
	return lib.Verify(dto)
}

type UpdateEmailDto struct {
	Email string `json:"email" validate:"required"`
	Code  int    `json:"code" validate:"required"`
}

func (dto *UpdateEmailDto) ToPayload(uid string) *userService.UpdateEmailPayload {
	return &userService.UpdateEmailPayload{
		Email: dto.Email,
		Code:  int64(dto.Code),
		Id:    uid,
	}
}

// Verify validates the structure and content of the SignUpDto.
func (dto *UpdateEmailDto) Verify() error {
	return lib.Verify(dto)
}

// UpdatePasswordDto represents the data transfer object for user's password
type UpdatePasswordDto struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=6"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}

func (dto *UpdatePasswordDto) ToPayload(uid string) *userService.UpdatePasswordPayload {
	return &userService.UpdatePasswordPayload{
		CurrentPassword: dto.CurrentPassword,
		NewPassword:     dto.NewPassword,
		Id:              uid,
	}
}

// Verify validates the structure and content of the UpdateUserDto.
func (dto *UpdatePasswordDto) Verify() error {
	return lib.Verify(dto)
}

// ResetPasswordDto represents the data transfer object for user's password
type ResetPasswordDto struct {
	NewPassword string `json:"newPassword" validate:"required,min=6"`
	Email       string `json:"email" validate:"required"`
	Code        int64  `json:"code" validate:"required"`
}

func (dto *ResetPasswordDto) ToPayload() *userService.ResetPasswordPayload {
	return &userService.ResetPasswordPayload{
		NewPassword: dto.NewPassword,
		Code:        dto.Code,
		Email:       dto.Email,
	}
}

// Verify validates the structure and content of the UpdateUserDto.
func (dto *ResetPasswordDto) Verify() error {
	return lib.Verify(dto)
}
