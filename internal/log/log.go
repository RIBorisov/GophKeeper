package log

import (
	"time"

	"github.com/rs/zerolog"
)

var gLogger zerolog.Logger

func InitLogger(l zerolog.Level) {
	zerolog.SetGlobalLevel(l)

	out := zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.DateTime
			w.NoColor = false
		},
	)

	gLogger = zerolog.New(out).With().Timestamp().Logger()
}

func GetLogger() *zerolog.Logger {
	if &gLogger == nil {
		InitLogger(zerolog.Level(0))
	}
	return &gLogger
}

func Debug(msg string, kv ...any) {
	gLogger.Debug().Fields(kv).Msg(msg)
}

func Info(msg string, kv ...any) {
	gLogger.Info().Fields(kv).Msg(msg)
}

func Error(msg string, kv ...any) {
	gLogger.Error().Fields(kv).Msg(msg)
}

func Warning(msg string, kv ...any) {
	gLogger.Warn().Fields(kv).Msg(msg)
}

func Fatal(msg string, kv ...any) {
	gLogger.Fatal().Fields(kv).Msg(msg)
}
