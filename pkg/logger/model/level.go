package model

import "go.uber.org/zap/zapcore"

type Level uint8

const (
	// TraceLevel causes a logger to emit messages logged at "TRACE" level or more
	// severe. It is typically only enabled during development, and usually results
	// in very verbose logging output.
	TraceLevel Level = iota + 1

	// DebugLevel causes a logger to emit messages logged at "DEBUG" level or more
	// severe. It is typically only enabled when debugging or during development,
	// and usually results in very verbose logging output.
	DebugLevel

	// InfoLevel causes a logger to emit messages logged at "INFO" level or more
	// severe. It is typically used for general operational entries about what's
	// going on inside an application.
	InfoLevel

	// WarnLevel causes a logger to emit messages logged at "WARN" level or more
	// severe. It is typically used for non-critical entries that deserve attention.
	WarnLevel

	// ErrorLevel causes a logger to emit messages logged at "ERROR" level or more
	// severe. It is typically used for errors that should definitely be noted, and
	// is commonly used for hooks to send errors to an error tracking service.
	ErrorLevel

	// FatalLevel causes a logger to emit messages logged at "FATAL" level or more
	// severe. Messages logged at this level cause a logger to log the message and
	// then call os.Exit(1). It will exit even if the logging level is set to PanicLevel.
	FatalLevel
)

// String implements fmt.Stringer for Level.
func (l Level) String() string {
	switch l {
	case FatalLevel:
		return "FATAL"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	case TraceLevel:
		return "TRACE"
	default:
		return "INFO"
	}
}

func (l Level) ZapLevel() zapcore.Level {
	switch l {
	case FatalLevel:
		return zapcore.FatalLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case TraceLevel:
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

// ParseLevel converts log level string to level constant
// if wrong string received it returns debug level.
func ParseLevel(logLevel string) Level {
	switch logLevel {
	case "trace":
		return TraceLevel
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
		return InfoLevel
	}
}
