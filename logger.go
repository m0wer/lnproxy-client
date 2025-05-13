package client

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel defines the severity of a log message
type LogLevel int

const (
	// Log levels from most to least severe
	LevelError LogLevel = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

// String returns the string representation of the LogLevel
func (l LogLevel) String() string {
	switch l {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging functionality
type Logger struct {
	level     LogLevel
	output    io.Writer
	mu        sync.Mutex
	prefix    string
	component string
}

var (
	// defaultLogger is the global logger instance
	defaultLogger     *Logger
	defaultLoggerOnce sync.Once
)

// DefaultLogger returns the default logger instance set to debug level
func DefaultLogger() *Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = &Logger{
			level:  LevelDebug, // Default to debug level
			output: os.Stderr,
		}
	})
	return defaultLogger
}

// NewLogger creates a new logger with the specified parameters
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
	}
}

// SetLevel changes the log level of the logger
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// SetOutput changes the output destination of the logger
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

// WithPrefix returns a new Logger with the specified prefix
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		level:     l.level,
		output:    l.output,
		prefix:    prefix,
		component: l.component,
	}
}

// WithComponent returns a new Logger with the specified component
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		level:     l.level,
		output:    l.output,
		prefix:    l.prefix,
		component: component,
	}
}

// log writes a log message with the given level and format
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level > l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	prefix := ""
	if l.prefix != "" {
		prefix = l.prefix + " "
	}

	component := ""
	if l.component != "" {
		component = "[" + l.component + "] "
	}

	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("%s [%s] %s%s%s\n", timestamp, level.String(), prefix, component, message)
	
	// We don't check for errors here as there's not much we can do if logging fails
	_, _ = io.WriteString(l.output, logLine)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Global helper functions that use the default logger

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	DefaultLogger().Debug(format, args...)
}

// Info logs an informational message using the default logger
func Info(format string, args ...interface{}) {
	DefaultLogger().Info(format, args...)
}

// Warn logs a warning message using the default logger
func Warn(format string, args ...interface{}) {
	DefaultLogger().Warn(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	DefaultLogger().Error(format, args...)
}

// SetGlobalLevel sets the log level for the default logger
func SetGlobalLevel(level LogLevel) {
	DefaultLogger().SetLevel(level)
}

// SetGlobalOutput sets the output for the default logger
func SetGlobalOutput(output io.Writer) {
	DefaultLogger().SetOutput(output)
}
