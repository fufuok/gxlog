package text

import (
	"regexp"
	"strings"
	"sync"

	"github.com/gratonos/gxlog"
)

var gHeaderRegexp = regexp.MustCompile("{{([^:%]*?)(?::([^%]*?))?(%.*?)?}}")

type Formatter struct {
	header      string
	enableColor bool

	colorMgr        *colorMgr
	headerAppenders []*headerAppender
	suffix          []byte
	buf             []byte

	lock sync.Mutex
}

func New(config *Config) *Formatter {
	formatter := &Formatter{
		colorMgr: newColorMgr(),
	}
	formatter.SetHeader(config.Header)
	formatter.enableColor = config.EnableColor
	formatter.MapColors(config.ColorMap)
	return formatter
}

func (this *Formatter) GetHeader() (header string) {
	this.lock.Lock()
	header = this.header
	this.lock.Unlock()

	return header
}

func (this *Formatter) SetHeader(header string) {
	this.lock.Lock()

	this.header = header
	this.headerAppenders = this.headerAppenders[:0]
	var staticText string
	for header != "" {
		indexes := gHeaderRegexp.FindStringSubmatchIndex(header)
		if indexes == nil {
			break
		}
		begin, end := indexes[0], indexes[1]
		staticText += header[:begin]
		element, property, fmtspec := extractElement(indexes, header)
		if this.addAppender(element, property, fmtspec, staticText) {
			staticText = ""
		}
		header = header[end:]
	}
	this.suffix = []byte(staticText + header)

	this.lock.Unlock()
}

func (this *Formatter) EnableColor() {
	this.lock.Lock()
	this.enableColor = true
	this.lock.Unlock()
}

func (this *Formatter) DisableColor() {
	this.lock.Lock()
	this.enableColor = false
	this.lock.Unlock()
}

func (this *Formatter) GetColor(level gxlog.LogLevel) (color ColorID) {
	this.lock.Lock()
	color = this.colorMgr.GetColor(level)
	this.lock.Unlock()

	return color
}

func (this *Formatter) SetColor(level gxlog.LogLevel, color ColorID) {
	this.lock.Lock()
	this.colorMgr.SetColor(level, color)
	this.lock.Unlock()
}

func (this *Formatter) MapColors(colorMap map[gxlog.LogLevel]ColorID) {
	this.lock.Lock()
	this.colorMgr.MapColors(colorMap)
	this.lock.Unlock()
}

func (this *Formatter) GetMarkedColor() (color ColorID) {
	this.lock.Lock()
	color = this.colorMgr.GetMarkedColor()
	this.lock.Unlock()

	return color
}

func (this *Formatter) SetMarkedColor(color ColorID) {
	this.lock.Lock()
	this.colorMgr.SetMarkedColor(color)
	this.lock.Unlock()
}

func (this *Formatter) Format(record *gxlog.Record) []byte {
	this.lock.Lock()

	var left, right []byte
	if this.enableColor {
		if record.Marked {
			left, right = this.colorMgr.GetMarkedColorEars()
		} else {
			left, right = this.colorMgr.GetColorEars(record.Level)
		}
	}
	this.buf = this.buf[:0]
	this.buf = append(this.buf, left...)
	for _, appender := range this.headerAppenders {
		this.buf = appender.AppendHeader(this.buf, record)
	}
	this.buf = append(this.buf, this.suffix...)
	this.buf = append(this.buf, right...)

	this.lock.Unlock()

	return this.buf
}

func (this *Formatter) addAppender(element, property, fmtspec, staticText string) bool {
	appender := newHeaderAppender(element, property, fmtspec, staticText)
	if appender == nil {
		return false
	}
	this.headerAppenders = append(this.headerAppenders, appender)
	return true
}

func extractElement(indexes []int, header string) (element, property, fmtspec string) {
	element = strings.ToLower(getField(header, indexes[2], indexes[3]))
	property = getField(header, indexes[4], indexes[5])
	fmtspec = getField(header, indexes[6], indexes[7])
	if fmtspec == "%" {
		fmtspec = ""
	}
	return element, property, fmtspec
}

func getField(str string, begin, end int) string {
	if begin < end {
		return strings.TrimSpace(str[begin:end])
	}
	return ""
}
