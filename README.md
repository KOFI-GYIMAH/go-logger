# Go Structured Logger

[![Go Reference](https://pkg.go.dev/badge/github.com/KOFI-GYIMAH/go-logger.svg)](https://pkg.go.dev/github.com/KOFI-GYIMAH/go-logger)  
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)  
[![Go Report Card](https://goreportcard.com/badge/github.com/KOFI-GYIMAH/go-logger)](https://goreportcard.com/report/github.com/KOFI-GYIMAH/go-logger)

A structured, colorful, and customizable logger for Go applications. Built for flexibility and designed for use in both development and production environments, with support for structured HTTP request logging and optional JSON formatting.

---

## 🚀 Features

- 🌈 Colorful output (with optional auto-disable for non-TTY environments)
- 📝 Structured log format with support for HTTP fields
- ⚡ Fast and lightweight
- 📌 Automatic caller info (`file:line`)
- 🧩 Fully customizable output, writer, and formatter
- 🌐 JSON log formatter for production environments
- 🎚️ Log levels: `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`

---

## 📦 Installation

```bash
go get github.com/KOFI-GYIMAH/go-logger
```

---

## 🧑‍💻 Basic Usage

```go
package main

import (
	"github.com/KOFI-GYIMAH/go-logger"
)

func main() {
	logger.Info("checkout completed for userId: 3463", logger.LogFields{
		CallerInfo: "cmd/main.go:12",
	})
}
```

---

## 🌐 HTTP Request Logging Example

```go
package main

import (
	"net/http"
	"time"
	
	"github.com/KOFI-GYIMAH/go-logger"
)

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Business logic...

	logger.Info("checkout completed for userId: 3463", logger.LogFields{
		Method:     r.Method,
		Path:       r.URL.Path,
		Status:     "200 OK",
		Size:       "42KB",
		Latency:    time.Since(start).String(),
		CallerInfo: "handlers/main.go:26",
	})
}
```

---

## ✨ Output Formats

### 1. Text Formatter (Default)

```text
[INFO] 2025-06-14T12:30:25Z POST /api/checkout 200 OK 51KB 120ms - checkout completed for userId: 3463 - go-logger/logger.go:132
```

**Fields:**

- `[LEVEL]`: Log level  
- Timestamp (RFC3339 format)  
- HTTP Method + Path (if provided)  
- Status, Size, Latency (if provided)  
- Message  
- Caller location  

### 2. JSON Formatter (Optional)

```json
{
  "level": "error",
  "timestamp": "2025-06-14T12:30:25Z",
  "message": "internal server error",
  "fields": {
    "method": "POST",
    "path": "/api/checkout",
    "status": "500 Internal Server Error",
    "latency": "250ms",
    "caller_info": "handlers/api.go:47"
  }
}
```

---

## ⚙️ Configuration Options

```go
// Set log level (DEBUG, INFO, WARN, ERROR, FATAL)
logger.SetLevel(logger.LevelDebug)

// Disable color output (e.g. in JSON or file logs)
logger.SetColor(false)

// Change log output (default: os.Stdout)
logger.SetOutput(os.Stderr)

// Set a custom formatter function
logger.SetFormatter(func(level logger.LogLevel, message string, fields logger.LogFields) string {
	return fmt.Sprintf("[%s] %s - %s", level, fields.Path, message)
})
```

---

## 🔄 Switching Between Formatters

```go
// Use default colorful formatter
logger.SetFormatter(logger.DefaultFormatter)

// Use JSON formatter (ideal for production/log ingestion)
logger.SetFormatter(logger.JSONFormatter)
```

---

## 🔧 Auto-Switch Formatter by Environment

```go
import "os"

if os.Getenv("ENV") == "production" {
	logger.SetColor(false)
	logger.SetFormatter(logger.JSONFormatter)
} else {
	logger.SetFormatter(logger.DefaultFormatter)
}
```

---

## 📊 Log Levels

| Level | Description |
|-------|-------------|
| `DEBUG` | Fine-grained debug info (development only) |
| `INFO`  | Routine operations and events |
| `WARN`  | Unusual situations, not errors |
| `ERROR` | Problems that should be investigated |
| `FATAL` | Critical issue; terminates application |

---

## 💡 Best Practices

- Use structured fields for logging HTTP requests
- Disable colors in CI/CD and file-based logs
- Enable `DEBUG` level only in local/dev environments
- Use the JSON formatter with a log collector in production

---

## ⚡ Performance Considerations

- Zero-cost log skipping if below threshold level
- Caller info is lazily computed (only when not provided)
- Color stripping uses a compiled regex (efficient)

---

## ✅ Tests

Run unit tests using:

```bash
go test ./...
```

Includes test coverage for:
- Log level filtering
- Custom formatting
- Output redirection
- Field inclusion
- JSON vs text log output

---

## 🤝 Contributing

PRs are welcome! To contribute:

1. Fork the repo  
2. Create your feature branch (`git checkout -b feature/foo`)  
3. Add tests and documentation  
4. Push to the branch  
5. Open a PR ✨  

---

## 📄 License

MIT © [Kofi Gyimah](https://github.com/KOFI-GYIMAH)
