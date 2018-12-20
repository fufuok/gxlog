// +build nacl plan9 windows

package syslog

import (
	"fmt"

	"github.com/gxlog/gxlog"
)

const cError = "not implemented on nacl, plan9 or windows"

type Writer struct{}

func Open(cfg *Config) (*Writer, error) {
	return nil, fmt.Errorf("writer/syslog.Open: %s", cError)
}

func (this *Writer) Close() error {
	return fmt.Errorf("writer/syslog.Close: %s", cError)
}

func (this *Writer) Write(bs []byte, record *gxlog.Record) {}

func (this *Writer) ReportOnErr() bool { return false }

func (this *Writer) SetReportOnErr(ok bool) {}

func (this *Writer) MapSeverity(severityMap map[gxlog.Level]Severity) {}
