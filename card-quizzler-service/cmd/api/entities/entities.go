package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
)

type JsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CreateUserDto struct {
	ID string `json:"id" validate:"required"`
}

type ResultTerm struct {
	models.Term
	ID     string `json:"id" validate:"required"`
	Answer bool   `json:"answer" validate:"required"`
}
