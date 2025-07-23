package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var (
	currentLevel LogLevel = InfoLevel
	silent       bool     = false
	showTime     bool     = false
	showTag      bool     = false
)

type LogLevel int

func (l LogLevel) ToStringValue() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

func (l LogLevel) ToLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return DebugLevel
	}
}

func GetAllLogLevelsString() string {
	return DebugLevel.ToStringValue() + ", " +
		InfoLevel.ToStringValue() + ", " +
		WarnLevel.ToStringValue() + ", " +
		ErrorLevel.ToStringValue() + ", " +
		FatalLevel.ToStringValue()
}

// SetLevel sets the global log level.
func SetLevel(level LogLevel) {
	currentLevel = level
}

// EnableSilentMode suppresses all logs.
func EnableSilentMode() {
	silent = true
}

// EnableTimestamp enables timestamps in logs.
func EnableTimestamp() {
	showTime = true
}

func EnableTag() {
	showTag = true
}

// Debug prints debug logs if enabled.
func Debug(format string, a ...interface{}) {
	log(DebugLevel, "DEBUG", Cyan(format), a...)
}

// Info prints info level logs.
func Info(format string, a ...interface{}) {
	log(InfoLevel, "INFO", Green(format), a...)
}

// Warn prints warning logs.
func Warn(format string, a ...interface{}) {
	log(WarnLevel, "WARN", Yellow(format), a...)
}

// Error prints error logs.
func Error(format string, a ...interface{}) {
	log(ErrorLevel, "ERROR", Red(format), a...)
}

// Fatal prints fatal level logs.
func Fatal(format string, a ...interface{}) {
	log(FatalLevel, "FATAL", Red(format), a...)
}

func log(level LogLevel, tag, format string, a ...interface{}) {
	if silent || level < currentLevel {
		return
	}
	var tagInfo string
	var timeInfo string
	if showTag {
		tagInfo = fmt.Sprintf("[%s]", tag)
	}
	if showTime {
		timeInfo = fmt.Sprintf("[%s]", time.Now().Format("15:04:05"))
	}
	prefix := fmt.Sprintf("%s %s", timeInfo, tagInfo)
	fmt.Fprintf(os.Stderr, "%s %s\n", prefix, fmt.Sprintf(format, a...))
}
