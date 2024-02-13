package appEntities

// JsonResponse represents a simple JSON response message structure.
type JsonResponse struct {
	Message string `json:"message"` // Message field for JSON response messages
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
		FromService: "user-service",
		Message:     message,
		Name:        name,
	}
}
