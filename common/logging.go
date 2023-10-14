package common

import (
	"fmt"
	"io"
	"log"
	"os"
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
	LogLevelDebug             = iota
	LogLevelInfo              = iota
	LogLevelWarning           = iota
	LogLevelError             = iota
	LogLevelCritical          = iota
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

	levelValue :=  strconv.ParseInt(text, 10, 0)
	switch levelValue {
	case 0:
		// Nothing, means the value wasn't text
	case LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError, LogLevelCritical:
		return parsedLevelValue, nil
	default:
		return LogLevelUnknown, fmt.Errorf("received bad log level number: %d", )
	}

	
	for level, name := range LogLevelNames {
		if level == levelValue || name == text {
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
		if  err == nil {
			return envLogLevel
		} 
		fmt.Fprintln(os.Stderr, "error parsing LOG_LEVEL: %s; %v", envLogLevelString, err)
	}

	fmt.Fprintln(os.Stderr, "using default log level: %s", defaultLogLevel.Name())
	return defaultLogLevel
}



const debugLogFlags = log.Lmsgprefix | log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
const nonDebugLogFlags = log.Lmsgprefix | log.Ldate | log.Ltime

type Logger struct {
	wrapped log.Logger
	level LogLevel
}

func ProvideLogger(writer LogWriter, prefix LogPrefix, debugMode DebugMode, logLevel LogLevel) Logger {
	flags := nonDebugLogFlags
	if debugMode {
		flags = debugLogFlags
	}

	return Logger {
		wrapped: log.New(writer, prefix, flags),
		level: logLevel,
	}
}
