package shared

import (
	"log"
	"os"
)

type Logger interface {
	Init() Logger
	Error(format string, v ...any)
	Info(format string, v ...any)
	Fatal(format string, v ...any)
	Warn(format string, v ...any)
}

type LoggerImpl struct {
	errorLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	fatalLogger *log.Logger
}

var loggerInstance Logger

func NewLogger() Logger {
	errorLogger := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	warnLogger := log.New(os.Stdout, "WARN\t", log.Ldate|log.Ltime)
	fatalLogger := log.New(os.Stdout, "FATAL\t", log.Ldate|log.Ltime)
	return &LoggerImpl{
		errorLogger: errorLogger,
		infoLogger:  infoLogger,
		warnLogger:  warnLogger,
		fatalLogger: fatalLogger,
	}
}

func (l *LoggerImpl) Error(format string, v ...any) {
	l.errorLogger.Printf(format, v)
}

func (l *LoggerImpl) Info(format string, v ...any) {
	l.infoLogger.Printf(format, v)
}

func (l *LoggerImpl) Warn(format string, v ...any) {
	l.infoLogger.Printf(format, v)
}

func (l *LoggerImpl) Fatal(format string, v ...any) {
	l.fatalLogger.Fatalf(format, v)
}

func (l *LoggerImpl) Init() Logger {
	instance := NewLogger()
	loggerInstance = instance
	return instance
}

func Log() Logger {
	return loggerInstance
}
