package main

import (
	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/json"
	"github.com/gxlog/gxlog/formatter/text"
	"github.com/gxlog/gxlog/iface"
	"github.com/gxlog/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
	testJSONFormatter()
	testTextFormatter()
}

func testJSONFormatter() {
	defer log.Timing("test Timing")()

	log.SetSlotFormatter(logger.Slot0, json.New(json.NewConfig()))

	log.WithMark(true).Info("json with mark")
	prelog := log.WithPrefix("prefix-123")
	prelog.WithContext("ah", "ha", "int", 123, "bool", false).Trace("json with context")
}

func testTextFormatter() {
	defer log.Timingf("%s", "test Timingf")()

	log.SetSlotFormatter(logger.Slot0, text.New(text.NewConfig()))

	log.Trace("test Trace")
	log.Tracef("%s", "test Tracef")
	log.Debug("test Debug")
	log.Debugf("%s", "test Debugf")
	log.WithMark(true).WithPrefix("prefix-123").WithContext("ah", "ha").Debug("test Debug with aux")
	log.Info("test Info")
	log.Infof("%s", "test Infof")
	log.Warn("test Warn")
	log.Warnf("%s", "test Warnf")
	log.Error("test Error")
	log.Errorf("%s", "test Errorf")
	log.LogError(iface.Error, "an error")
	// Fatal and Fatalf will output the stack of current goroutine by default.
	log.Fatal("test Fatal")
	log.Fatalf("%s", "test Fatalf")
}

// {"time":"2021-04-03T23:39:08+08:00","level":"I","file":"default.go","line":23,"pkg":"main","func":"testJSONFormatter","msg":"json with mark","marked":true}
// {"time":"2021-04-03T23:39:08+08:00","level":"T","file":"default.go","line":25,"pkg":"main","func":"testJSONFormatter","msg":"json with context","prefix":"prefix-123","contexts":[{"ah":"ha"},{"int":"123"},{"bool":"false"}]}
// {"time":"2021-04-03T23:39:08+08:00","level":"T","file":"default.go","line":26,"pkg":"main","func":"testJSONFormatter","msg":"test Timing (cost: 19.0473ms)"}
//23:39:08.020 T default.go:33 main.testTextFormatter [] test Trace
//23:39:08.020 T default.go:34 main.testTextFormatter [] test Tracef
//23:39:08.020 D default.go:35 main.testTextFormatter [] test Debug
//23:39:08.020 D default.go:36 main.testTextFormatter [] test Debugf
//23:39:08.020 D default.go:37 main.testTextFormatter prefix-123[(ah: ha)] test Debug with aux
//23:39:08.020 I default.go:38 main.testTextFormatter [] test Info
//23:39:08.020 I default.go:39 main.testTextFormatter [] test Infof
//23:39:08.020 W default.go:40 main.testTextFormatter [] test Warn
//23:39:08.020 W default.go:41 main.testTextFormatter [] test Warnf
//23:39:08.020 E default.go:42 main.testTextFormatter [] test Error
//23:39:08.020 E default.go:43 main.testTextFormatter [] test Errorf
//23:39:08.020 E default.go:44 main.testTextFormatter [] an error
//23:39:08.020 F default.go:46 main.testTextFormatter [] test Fatal
// goroutine 1 [running]:
// runtime/debug.Stack(0xc000136000, 0x1, 0x6)
//	C:/Go/src/runtime/debug/stack.go:24 +0xa5
// github.com/gxlog/gxlog/logger.(*Logger).Log(0xc000136000, 0x1, 0x6, 0xc00010fd68, 0x1, 0x1)
//	F:/Go/gxlog/logger/logger.go:149 +0xd2
// github.com/gxlog/gxlog/logger.(*Logger).Fatal(...)
//	F:/Go/gxlog/logger/logger.go:106
// main.testTextFormatter()
//	F:/Go/gxlog/_examples/default/default.go:46 +0x8af
// main.main()
//	F:/Go/gxlog/_examples/default/default.go:15 +0x2c
//23:39:08.020 F default.go:47 main.testTextFormatter [] test Fatalf
// goroutine 1 [running]:
// runtime/debug.Stack(0x0, 0xae056e, 0x2)
//	C:/Go/src/runtime/debug/stack.go:24 +0xa5
// github.com/gxlog/gxlog/logger.(*Logger).Logf(0xc000136000, 0x1, 0x6, 0xc00000a9a0, 0x5, 0xc00010fd58, 0x1, 0x1)
//	F:/Go/gxlog/logger/logger.go:168 +0x165
// github.com/gxlog/gxlog/logger.(*Logger).Fatalf(...)
//	F:/Go/gxlog/logger/logger.go:111
// main.testTextFormatter()
//	F:/Go/gxlog/_examples/default/default.go:47 +0x92e
// main.main()
//	F:/Go/gxlog/_examples/default/default.go:15 +0x2c
//23:39:08.020 T default.go:48 main.testTextFormatter [] test Timingf (cost: 148.6µs)
//
// Process finished with exit code 0
