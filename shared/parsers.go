package lib

import (
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"time"
)

// ParseDuration parses hours
func ParseDuration(config map[string]string, key string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(config[key] + "h")
	if err != nil {
		return defaultValue
	}
	return duration
}

func ParseUUID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return parsedID, goErrorHandler.ParseUUIDFailure()
	}
	return parsedID, nil
}
