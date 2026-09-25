package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const FormatJSON = "json"

func New(logLevel string, logFormat string) (logger zerolog.Logger) {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(logLevel); err == nil {
		level = newLevel
	}

	var out io.Writer = os.Stdout
	if logFormat != FormatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return zerolog.New(out).Level(level).With().Timestamp().Stack().Logger()
}
