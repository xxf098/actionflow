package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"golang.org/x/term"
)

func New(cfg LogConfig) zerolog.Logger {
	logger := zerolog.
		New(os.Stderr).
		With().
		Timestamp().
		Logger()

	logger = logger.Output(&PlainOutput{Out: colorable.NewColorableStderr()})

	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		panic(err)
	}
	return logger.Level(lvl)
}

type LogConfig struct {
	Level  string
	Format string
}

func TeeCloud(w io.Writer) zerolog.LevelWriter {
	return zerolog.MultiLevelWriter(w)
}

func jsonLogs(format string) bool {
	switch format {
	case "json":
		return true
	case "plain":
		return false
	case "tty":
		return false
	case "tty2":
		return false
	case "auto":
		return !term.IsTerminal(int(os.Stdout.Fd()))
	default:
		fmt.Fprintf(os.Stderr, "invalid --log-format %q\n", format)
		os.Exit(1)
	}
	return false
}
