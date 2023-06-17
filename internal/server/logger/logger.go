package logger

import (
	"github.com/rs/zerolog"
	"os"
)

var Log zerolog.Logger

func InitLogger() {
	Log = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
}
