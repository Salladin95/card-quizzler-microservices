package entities

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
)

type JsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
		FromService: "card-quizzler-service",
		Message:     message,
		Name:        name,
	}
}

type CreateUserDto struct {
	ID string `json:"id" validate:"required"`
}

type ResultTerm struct {
	models.Term
	ID     string `json:"id" validate:"required"`
	Answer bool   `json:"answer" validate:"required"`
}
