package log

import (
	"os"
)

var log = &MulLog{
	logs: []Logger{
		NewWindowsDebugLog(),
	},
}

func AddLogger(logger Logger) {
	log.logs = append(log.logs, logger)
}

func Debug(format string, a ...interface{}) {
	log.Debug(format, a...)
}

func Info(format string, a ...interface{}) {
	log.Info(format, a...)
}

func Error(format string, a ...interface{}) {
	log.Error(format, a...)
}

func Fatal(format string, a ...interface{}) {
	log.Fatal(format, a...)
}

type Logger interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

type MulLog struct {
	logs []Logger
}

func (m *MulLog) Fatal(format string, a ...interface{}) {
	log.Fatal(format, a...)
	os.Exit(0)
}

func (m *MulLog) Debug(format string, a ...interface{}) {
	for _, log := range m.logs {
		log.Debug(format, a...)
	}
}

func (m *MulLog) Info(format string, a ...interface{}) {
	for _, log := range m.logs {
		log.Info(format, a...)
	}
}

func (m *MulLog) Error(format string, a ...interface{}) {
	for _, log := range m.logs {
		log.Error(format, a...)
	}
}
