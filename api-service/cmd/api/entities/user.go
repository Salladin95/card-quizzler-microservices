package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
	"time"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Birthday  string    `json:"birthday"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
