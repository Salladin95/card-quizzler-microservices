package lib

import "log/slog"

func LogError(msg string, args ...any) {
	slog.Error(msg, args)
}

func LogInfo(msg string, args ...any) {
	slog.Info(msg, args)
}
