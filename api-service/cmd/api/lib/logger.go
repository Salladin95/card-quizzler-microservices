package lib

import (
	"github.com/rs/zerolog"
	"os"
)

var Logger zerolog.Logger

func InitLogger() {
	// Create a Zero Log logger instance
	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Set logger output to JSON format
	Logger = Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
