package common

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type WrappedLogger *log.Logger

// LogPrefix should be injected from elsewhere
// See LoggingWireSetFactory
type LogPrefix string

// LogWriter should be injected from elsewhere
// See LoggingWireSetFactory
type LogWriter io.Writer

// LogLevel is what it says on the tin
// Any value other than the debug/info/warning/error/critical will raise an error during configuration
type LogLevel int

const (
	LogLevelUnknown  LogLevel = 0
	LogLevelDebug    LogLevel = iota
	LogLevelInfo     LogLevel = iota
	LogLevelWarning  LogLevel = iota
	LogLevelError    LogLevel = iota
	LogLevelCritical LogLevel = iota
)

var LogLevelNames = map[LogLevel]string{
	LogLevelDebug:    "DEBUG",
	LogLevelInfo:     "INFO",
	LogLevelWarning:  "WARNING",
	LogLevelError:    "ERROR",
	LogLevelCritical: "CRITICAL",
}

func ParseLogLevel(text string) (LogLevel, error) {
	text = strings.TrimSpace(text)
	text = strings.ToUpper(text)

	levelValue, _ := strconv.ParseInt(text, 10, 0)
	switch LogLevel(levelValue) {
	case 0:
		// Nothing, means the value wasn't text
	case LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError, LogLevelCritical:
		return LogLevel(levelValue), nil
	default:
		return LogLevelUnknown, fmt.Errorf("received bad log level number: %d", levelValue)
	}

	for level, name := range LogLevelNames {
		if int64(level) == levelValue || name == text {
			return level, nil
		}
	}

	return 0, fmt.Errorf("no log level match: %v", text)
}

func (ll LogLevel) Name() string {
	if name, ok := LogLevelNames[ll]; ok {
		return name
	}
	return fmt.Sprintf("INVALID_LOG_LEVEL_%d", ll)
}

const defaultLogLevel = LogLevelWarning

func ProvideLogLevel(debugMode DebugMode, environ Environ) LogLevel {
	if debugMode {
		return LogLevelDebug
	}

	envLogLevelString, _ := environ.GetString("LOG_LEVEL")
	if envLogLevelString != "" {
		envLogLevel, err := ParseLogLevel(envLogLevelString)
		if err == nil {
			return envLogLevel
		}
		confWarningLogger.Printf("error parsing LOG_LEVEL: %s; %v\n", envLogLevelString, err)
	}

	confWarningLogger.Printf("using default log level: %s", defaultLogLevel.Name())
	return defaultLogLevel
}

const debugLogFlags = log.Lmsgprefix | log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
const nonDebugLogFlags = log.Lmsgprefix | log.Ldate | log.Ltime

type Logger struct {
	wrapped *log.Logger
	level   LogLevel
}

func ProvideLogger(writer LogWriter, prefix LogPrefix, debugMode DebugMode, logLevel LogLevel) Logger {
	flags := nonDebugLogFlags
	if debugMode {
		flags = debugLogFlags
	}

	return Logger{
		wrapped: log.New(writer, string(prefix), flags),
		level:   logLevel,
	}
}

func (logger Logger) Logf(callDepth int, level LogLevel, format string, args ...any) {
	if level < logger.level {
		return
	}

	formatWithLevel := fmt.Sprintf("[%s] %s", level.Name(), format)
	message := fmt.Sprintf(formatWithLevel, args...)
	// logger.wrapped.Print(message)
	logger.wrapped.Output(callDepth+2, message)

	if level >= LogLevelCritical {
		panic("received critical signal")
	}
}

func (logger Logger) Debugf(format string, args ...any) {
	logger.Logf(1, LogLevelDebug, format, args...)
}

func (logger Logger) Infof(format string, args ...any) {
	logger.Logf(1, LogLevelInfo, format, args...)
}

func (logger Logger) Warningf(format string, args ...any) {
	logger.Logf(1, LogLevelWarning, format, args...)
}

func (logger Logger) Errorf(format string, args ...any) {
	logger.Logf(1, LogLevelError, format, args...)
}

func (logger Logger) Criticalf(format string, args ...any) {
	logger.Logf(1, LogLevelCritical, format, args...)
}

var confWarningLogger *log.Logger = log.Default()

func init() {
	confWarningLogger = log.Default()
	confWarningLogger.SetPrefix("[conf warning]")
}
