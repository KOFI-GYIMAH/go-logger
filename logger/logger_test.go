package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestLogLevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger("test")
	logger.SetOutput(buf)
	logger.SetLevel(LevelError)

	logger.Info("this should not appear", LogFields{})
	logger.Error("this should appear", LogFields{})

	output := buf.String()
	if strings.Contains(output, "this should not appear") {
		t.Error("log level filtering failed: INFO message was logged")
	}
	if !strings.Contains(output, "this should appear") {
		t.Error("log level filtering failed: ERROR message was not logged")
	}
}

func TestDefaultFormatterOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger("test")
	logger.SetOutput(buf)
	logger.SetColor(false)
	logger.SetLevel(LevelDebug)

	logger.Debug("test debug message", LogFields{Method: "GET", Path: "/status"})

	out := buf.String()
	if !strings.Contains(out, "DEBUG") || !strings.Contains(out, "GET /status") {
		t.Errorf("default formatter did not output expected content: %s", out)
	}
}

func TestJSONFormatterOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger("json")
	logger.SetOutput(buf)
	logger.SetLevel(LevelInfo)
	logger.SetFormatter(JSONFormatter)

	logger.Info("json output test", LogFields{Status: "200 OK", Path: "/ping"})

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatalf("failed to parse JSON log: %v", err)
	}

	if logEntry["message"] != "json output test" {
		t.Error("JSON log message not found")
	}

	fields, ok := logEntry["fields"].(map[string]interface{})
	if !ok {
		t.Fatal("JSON fields not found or not a map")
	}
	if fields["status"] != "200 OK" || fields["path"] != "/ping" {
		t.Error("JSON log fields incorrect")
	}
}

func TestLoggerConcurrency(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger("concurrent")
	logger.SetOutput(buf)
	logger.SetColor(false)
	logger.SetLevel(LevelInfo)

	done := make(chan bool)

	for i := range 10 {
		go func(i int) {
			logger.Info("concurrent log", LogFields{Status: "200"})
			done <- true
		}(i)
	}

	for range 10 {
		<-done
	}

	out := buf.String()
	count := strings.Count(out, "concurrent log")
	if count != 10 {
		t.Errorf("expected 10 log lines, got %d", count)
	}
}

func TestLoggerWithContext(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger("ctx")
	logger.SetOutput(buf)
	logger.SetLevel(LevelInfo)

	ctx := context.WithValue(context.Background(), "userID", "abc123")
	logger.LogCtx(ctx, LevelInfo, "log with context", LogFields{Path: "/ctx"})

	out := buf.String()
	if !strings.Contains(out, "log with context") {
		t.Error("context-based logging failed")
	}
}
