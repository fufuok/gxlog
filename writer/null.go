package writer

import (
	"github.com/fufuok/gxlog/iface"
)

var nullWriter = Func(func([]byte, *iface.Record) {})

// Null returns the null Writer.
func Null() iface.Writer {
	return nullWriter
}
