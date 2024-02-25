package entities

type JsonResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
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
