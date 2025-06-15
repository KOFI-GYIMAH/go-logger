package main

import (
	. "github.com/KOFI-GYIMAH/go-logger/logger"
)

func main() {
	fields := NewLogFields("POST", "/api/checkout", "200 OK", "51KB", "120ms")

	// Using global logger
	Info("Checkout completed for user ID 3463", fields)

	// Using instance logger
	l := NewLogger("instance")
	l.SetColor(false) // disable color for testing
	l.SetFormatter(JSONFormatter)
	l.Error("Internal server error", fields)
}
