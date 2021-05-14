package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fufuok/gxlog/formatter/json"
	"github.com/fufuok/gxlog/formatter/text"
	"github.com/fufuok/gxlog/iface"
	"github.com/fufuok/gxlog/writer"
	"github.com/fufuok/gxlog/writer/file"

	"github.com/fufuok/gxlog"
	"github.com/fufuok/gxlog/logger"
)

const (
	// 修改该参数!!!
	Debug = true

	// 文件日志
	LogDir      = "."
	ProjectName = "app-ff"

	// 每 1 秒最多输出 3 条日志
	LogLimitTime = time.Second
	LogLimitNum  = 3

	// 日志级别: 1Trace 2Debug 3Info 4Warn 5Error(默认) 6Fatal 7Off
	LogLevel     = 2
	LogFileLevel = 4
)

var (
	Log       = gxlog.Logger()
	LogLimit  *logger.Logger
	LogPrefix *logger.Logger
)

func init() {
	if err := InitLogger(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func InitLogger() error {
	if err := LogConfig(); err != nil {
		return err
	}

	// 带限制的日志记录器
	LogLimit = Log.WithTimeLimit(LogLimitTime, LogLimitNum)

	// 某些模块特定前缀的日志记录器
	LogPrefix = Log.WithPrefix("***FUFU***")

	return nil
}

// 日志配置
// 1. Debug 时, 只输出到控制台
// 2. 生产环境时, 输出到日志文件, 并发送 JSON 日志到 ES (不同的日志级别)
func LogConfig() error {
	Log.SetSlotLevel(logger.Slot0, iface.Level(LogLevel))
	Log.SetSlotFormatter(logger.Slot0, text.New(text.NewConfig()))

	if Debug {
		// 测试环境只输出到控制台
		return nil
	}

	// 错误文件日志
	wt, err := file.Open(file.Config{
		Path: LogDir,
		Base: ProjectName,
	})
	if err != nil {
		return fmt.Errorf("logfile path err: %s\nbye.", err)
	}
	Log.CopySlot(logger.Slot1, logger.Slot0)
	Log.SetSlotWriter(logger.Slot1, wt)
	Log.SetSlotLevel(logger.Slot1, iface.Level(LogFileLevel))

	// JSON 日志, 发送到 ES
	Log.SetSlotFormatter(logger.Slot0, json.New(json.NewConfig()))
	Log.SetSlotWriter(logger.Slot0, writer.Func(func(bs []byte, _ *iface.Record) {
		LogCache(bs)
	}))

	return nil
}

func LogCache(bs []byte) {
	// 日志暂存, 定时批量发送到 ES
	fmt.Println("___LOG_CACHE:", string(bs))
}

func main() {
	Log.Debug("test DEBUG")
	Log.Info("test INFO")
	Log.Warn("test WARN")
	Log.WithContext("k1", "v1.string", "k2", 1.05).Error("test ERROR with Field")

	LogPrefix.WithMark(true).Warnf("test %s", "PREFIX AND MARK")

	for i := 0; i < 10; i++ {
		LogLimit.Infof(">>>test limit: %d", i)
		time.Sleep(200 * time.Millisecond)
	}
}

// Debug = true
//23:16:17.005 D project_demo.go:100 main.main [] test DEBUG
//23:16:17.005 I project_demo.go:101 main.main [] test INFO
//23:16:17.005 W project_demo.go:102 main.main [] test WARN
//23:16:17.005 E project_demo.go:103 main.main [(k1: v1.string) (k2: 1.05)] test ERROR with Field
//23:16:17.005 W project_demo.go:105 main.main ***FUFU***[] test PREFIX AND MARK
//23:16:17.005 I project_demo.go:108 main.main [] >>>test limit: 0
//23:16:17.206 I project_demo.go:108 main.main [] >>>test limit: 1
//23:16:17.406 I project_demo.go:108 main.main [] >>>test limit: 2
//23:16:18.009 I project_demo.go:108 main.main [] >>>test limit: 5
//23:16:18.209 I project_demo.go:108 main.main [] >>>test limit: 6
//23:16:18.410 I project_demo.go:108 main.main [] >>>test limit: 7

// Debug = false
// ___LOG_CACHE: {"time":"2021-04-03T23:18:45+08:00","level":"D","file":"project_demo.go","line":101,"pkg":"main","func":"main","msg":"test DEBUG"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:45+08:00","level":"I","file":"project_demo.go","line":102,"pkg":"main","func":"main","msg":"test INFO"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:45+08:00","level":"W","file":"project_demo.go","line":103,"pkg":"main","func":"main","msg":"test WARN"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:45+08:00","level":"E","file":"project_demo.go","line":104,"pkg":"main","func":"main","msg":"test ERROR with Field","contexts":[{"k1":"v1.string"},{"k2":"1.05"}]}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:45+08:00","level":"W","file":"project_demo.go","line":106,"pkg":"main","func":"main","msg":"test PREFIX AND MARK","prefix":"***FUFU***","marked":true}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:45+08:00","level":"I","file":"project_demo.go","line":109,"pkg":"main","func":"main","msg":">>>test limit: 0"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:46+08:00","level":"I","file":"project_demo.go","line":109,"pkg":"main","func":"main","msg":">>>test limit: 1"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:46+08:00","level":"I","file":"project_demo.go","line":109,"pkg":"main","func":"main","msg":">>>test limit: 2"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:46+08:00","level":"I","file":"project_demo.go","line":109,"pkg":"main","func":"main","msg":">>>test limit: 5"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:47+08:00","level":"I","file":"project_demo.go","line":109,"pkg":"main","func":"main","msg":">>>test limit: 6"}
//
// ___LOG_CACHE: {"time":"2021-04-03T23:18:47+08:00","level":"I","file":"project_demo.go","line":109,"pkg":"main","func":"main","msg":">>>test limit: 7"}

// 20210403/app-ff.231845.297245.log
//23:18:45.297 W project_demo.go:103 main.main [] test WARN
//23:18:45.806 E project_demo.go:104 main.main [(k1: v1.string) (k2: 1.05)] test ERROR with Field
//23:18:45.806 W project_demo.go:106 main.main ***FUFU***[] test PREFIX AND MARK
