package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func ParseLevel(level string) Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return LevelInfo
	}
}

type Logger struct {
	level  Level
	output io.Writer
}

func NewLogger(level Level) *Logger {
	return &Logger{
		level:  level,
		output: os.Stdout,
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "[%s] %s: %s\n", timestamp, level.String(), message)

	if level == LevelFatal {
		os.Exit(1)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LevelFatal, format, args...)
}

// ------------------------------------------------------
// DEFAULT LOGGER
// ------------------------------------------------------

// Initialize with default log level, will be overridden by config if available
var defaultLogger = NewLogger(LevelInfo)

func InitFromEnv() {
	levelStr := os.Getenv("LOG_LEVEL")
	// fmt.Println("Log level - " + levelStr)
	if levelStr != "" {
		newLevel := ParseLevel(levelStr)
		defaultLogger.level = newLevel
	}

	outputPath := os.Getenv("LOG_OUTPUT")
	if outputPath != "" && outputPath != "stdout" {
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			defaultLogger.output = file
		}
	}
}

func SetLevel(level Level) {
	defaultLogger.level = level
}

func SetOutput(w io.Writer) {
	defaultLogger.output = w
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}
