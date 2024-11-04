package log

import (
	"time"

	"github.com/rs/zerolog"
)

var globalLogger zerolog.Logger

func InitLogger(l zerolog.Level) {
	zerolog.SetGlobalLevel(l)

	out := zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.DateTime
			w.NoColor = false
		},
	)

	globalLogger = zerolog.New(out).With().Timestamp().Logger()
}

func Debug(msg string, kv ...any) {
	globalLogger.Debug().Fields(kv).Msg(msg)
}

func Info(msg string, kv ...any) {
	globalLogger.Info().Fields(kv).Msg(msg)
}

func Error(msg string, kv ...any) {
	globalLogger.Error().Fields(kv).Msg(msg)
}

func Warning(msg string, kv ...any) {
	globalLogger.Warn().Fields(kv).Msg(msg)
}

func Fatal(msg string, kv ...any) {
	globalLogger.Fatal().Fields(kv).Msg(msg)
}
