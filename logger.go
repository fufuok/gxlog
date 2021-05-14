// Package gxlog provides the default Logger and the default Formatter.
package gxlog

import (
	"os"

	"github.com/fufuok/gxlog/formatter/text"
	"github.com/fufuok/gxlog/logger"
	"github.com/fufuok/gxlog/writer"
)

var (
	defaultLogger    *logger.Logger
	defaultFormatter *text.Formatter
)

func init() {
	defaultLogger = logger.New(logger.Config{})
	defaultFormatter = text.New(text.Config{})
	defaultLogger.Link(logger.Slot0, defaultFormatter, writer.Wrap(os.Stderr, nil))
}

// Logger returns the default Logger which has the default Formatter (a text
// formatter) and a writer wrapper of os.Stderr linked in Slot0.
// The rest slots are free.
func Logger() *logger.Logger {
	return defaultLogger
}

// Formatter returns the default Formatter. It is a text formatter.
func Formatter() *text.Formatter {
	return defaultFormatter
}
