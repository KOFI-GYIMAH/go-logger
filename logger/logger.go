package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Logger struct {
	name      string
	level     LogLevel
	out       io.Writer
	color     bool
	formatter func(level LogLevel, message string, fields LogFields) string
	mutex     sync.Mutex
}

type LogFields struct {
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	Status     string `json:"status,omitempty"`
	Size       string `json:"size,omitempty"`
	Latency    string `json:"latency,omitempty"`
	CallerInfo string `json:"caller_info,omitempty"`
}

var log = NewLogger("app")

func NewLogger(name string) *Logger {
	return &Logger{
		name:      name,
		level:     LevelInfo,
		out:       os.Stdout,
		color:     true,
		formatter: defaultFormatter,
	}
}

// * Global logger helper functions
func Debug(msg string, fields LogFields) { log.Log(LevelDebug, msg, fields) }
func Info(msg string, fields LogFields)  { log.Log(LevelInfo, msg, fields) }
func Warn(msg string, fields LogFields)  { log.Log(LevelWarn, msg, fields) }
func Error(msg string, fields LogFields) { log.Log(LevelError, msg, fields) }
func Fatal(msg string, fields LogFields) { log.Log(LevelFatal, msg, fields); os.Exit(1) }

// * Instance logger methods
func (l *Logger) Debug(msg string, fields LogFields) {
	l.Log(LevelDebug, msg, fields)
}
func (l *Logger) Info(msg string, fields LogFields) {
	l.Log(LevelInfo, msg, fields)
}
func (l *Logger) Warn(msg string, fields LogFields) {
	l.Log(LevelWarn, msg, fields)
}
func (l *Logger) Error(msg string, fields LogFields) {
	l.Log(LevelError, msg, fields)
}
func (l *Logger) Fatal(msg string, fields LogFields) {
	l.Log(LevelFatal, msg, fields)
	os.Exit(1)
}

func (l *Logger) SetLevel(level LogLevel) { l.level = level }
func (l *Logger) SetOutput(out io.Writer) { l.out = out }
func (l *Logger) SetColor(enabled bool)   { l.color = enabled }
func (l *Logger) SetFormatter(fn func(LogLevel, string, LogFields) string) {
	l.formatter = fn
}

func (l *Logger) Log(level LogLevel, msg string, fields LogFields) {
	if level < l.level {
		return
	}

	if fields.CallerInfo == "" {
		_, file, line, _ := runtime.Caller(2)
		fileParts := strings.Split(file, "/")
		if len(fileParts) > 3 {
			fileParts = fileParts[len(fileParts)-3:]
		}
		fields.CallerInfo = fmt.Sprintf("%s:%d", strings.Join(fileParts, "/"), line)
	}

	msg = strings.ReplaceAll(msg, "\n", " ")

	output := l.formatter(level, msg, fields)
	if !l.color {
		output = stripColors(output)
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	fmt.Fprintln(l.out, output)
}

func (l *Logger) LogCtx(ctx context.Context, level LogLevel, msg string, fields LogFields) {
	l.Log(level, msg, fields)
}

func defaultFormatter(level LogLevel, message string, fields LogFields) string {
	timestamp := time.Now().Format(time.RFC3339)
	var levelStr string
	var levelColor *color.Color

	switch level {
	case LevelDebug:
		levelStr = "DEBUG"
		levelColor = color.New(color.FgHiCyan)
	case LevelInfo:
		levelStr = "INFO"
		levelColor = color.New(color.FgHiGreen)
	case LevelWarn:
		levelStr = "WARN"
		levelColor = color.New(color.FgHiYellow)
	case LevelError:
		levelStr = "ERROR"
		levelColor = color.New(color.FgHiRed)
	case LevelFatal:
		levelStr = "FATAL"
		levelColor = color.New(color.FgHiMagenta)
	}

	components := []string{
		levelColor.Sprintf("[%s]", levelStr),
		color.New(color.FgHiWhite).Sprint(timestamp),
	}

	if fields.Method != "" && fields.Path != "" {
		components = append(components, color.New(color.FgHiBlue).Sprintf("%s %s", fields.Method, fields.Path))
	}
	if fields.Status != "" {
		components = append(components, color.New(color.FgHiGreen).Sprint(fields.Status))
	}
	if fields.Size != "" {
		components = append(components, color.New(color.FgHiYellow).Sprint(fields.Size))
	}
	if fields.Latency != "" {
		components = append(components, color.New(color.FgHiMagenta).Sprint(fields.Latency))
	}

	components = append(components,
		color.New(color.FgHiWhite).Sprint("-"),
		color.New(color.FgCyan).Sprint(message))

	components = append(components,
		color.New(color.FgHiWhite).Sprint("-"),
		color.New(color.FgHiBlue).Sprint(fields.CallerInfo))

	return strings.Join(components, " ")
}

func JSONFormatter(level LogLevel, message string, fields LogFields) string {
	entry := map[string]any{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level.String(),
		"message":   message,
		"fields":    fields,
	}
	data, _ := json.Marshal(entry)
	return string(data)
}

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

func stripColors(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[mGK]`)
	return re.ReplaceAllString(s, "")
}

func NewLogFields(method, path, status, size, latency string) LogFields {
	return LogFields{
		Method:  method,
		Path:    path,
		Status:  status,
		Size:    size,
		Latency: latency,
	}
}
