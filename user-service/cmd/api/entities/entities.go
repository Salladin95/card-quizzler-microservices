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
