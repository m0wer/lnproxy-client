package client

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	
	// Create a new logger that writes to our buffer with debug level
	logger := NewLogger(LevelDebug, &buf)
	
	// Test debug logs
	logger.Debug("This is a %s message", "debug")
	output := buf.String()
	if !strings.Contains(output, "DEBUG") || !strings.Contains(output, "This is a debug message") {
		t.Errorf("Debug log failed, got: %s", output)
	}
	
	// Reset the buffer
	buf.Reset()
	
	// Test info logs
	logger.Info("This is an %s message", "info")
	output = buf.String()
	if !strings.Contains(output, "INFO") || !strings.Contains(output, "This is an info message") {
		t.Errorf("Info log failed, got: %s", output)
	}
	
	// Reset the buffer
	buf.Reset()
	
	// Test log levels - set to INFO, DEBUG should not appear
	logger.SetLevel(LevelInfo)
	logger.Debug("This should not appear")
	logger.Info("This should appear")
	
	output = buf.String()
	if strings.Contains(output, "This should not appear") {
		t.Error("Debug message shouldn't appear when level is set to INFO")
	}
	if !strings.Contains(output, "This should appear") {
		t.Error("Info message should appear when level is set to INFO")
	}
	
	// Test with component
	buf.Reset()
	componentLogger := logger.WithComponent("TestComponent")
	componentLogger.Info("Component log test")
	
	output = buf.String()
	if !strings.Contains(output, "[TestComponent]") {
		t.Errorf("Component not in log output: %s", output)
	}
	
	// Test with prefix
	buf.Reset()
	prefixLogger := logger.WithPrefix("PREFIX")
	prefixLogger.Info("Prefix log test")
	
	output = buf.String()
	if !strings.Contains(output, "PREFIX") {
		t.Errorf("Prefix not in log output: %s", output)
	}
}

func TestGlobalLogger(t *testing.T) {
	// Save the original output to restore it later
	originalLogger := DefaultLogger()
	originalLevel := originalLogger.GetLevel()
	
	// Create a buffer to capture log output
	var buf bytes.Buffer
	SetGlobalOutput(&buf)
	
	// Test global debug logs
	Debug("Global %s test", "debug")
	output := buf.String()
	if !strings.Contains(output, "DEBUG") || !strings.Contains(output, "Global debug test") {
		t.Errorf("Global debug log failed, got: %s", output)
	}
	
	// Reset the buffer
	buf.Reset()
	
	// Test level filtering with global logger
	SetGlobalLevel(LevelInfo)
	Debug("This should not appear")
	Info("This should appear")
	
	output = buf.String()
	if strings.Contains(output, "This should not appear") {
		t.Error("Debug message shouldn't appear when global level is set to INFO")
	}
	if !strings.Contains(output, "This should appear") {
		t.Error("Info message should appear when global level is set to INFO")
	}
	
	// Restore original logger settings
	SetGlobalLevel(originalLevel)
}
