package main

import (
	"fmt"
	"os"

	"github.com/fufuok/gxlog/formatter/json"
	"github.com/fufuok/gxlog/formatter/text"
	"github.com/fufuok/gxlog/iface"
	"github.com/fufuok/gxlog/logger"
	"github.com/fufuok/gxlog/writer"
	"github.com/fufuok/gxlog/writer/file"
)

func main() {
	log := logger.New(logger.Config{
		Disabled:   logger.DynamicContext | logger.LimitByTime,
		TrackLevel: iface.Off,
	})

	fileWriter, err := file.Open(file.Config{
		Path: "/tmp/gxlog",
		Base: "test",
		// ErrorHandler will be called when an error occurs.
		ErrorHandler: writer.Report,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fileWriter.Close()

	// Logs will be formatted to text format and output to os.Stderr, then
	// formatted to json format and output to log files in /tmp/gxlog.
	log.Link(logger.Slot0, text.New(text.Config{
		Coloring: true,
	}), writer.Wrap(os.Stderr, nil))
	log.Link(logger.Slot1, json.New(json.Config{}), fileWriter)

	log.Info("test")
}
