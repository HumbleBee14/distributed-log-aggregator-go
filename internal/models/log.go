package models

import (
	"time"
)

// LogEntry represents a single log message from a service
type LogEntry struct {
	ID          string    `json:"id,omitempty"`
	ServiceName string    `json:"service_name" validate:"required"`
	Timestamp   time.Time `json:"timestamp" validate:"required"`
	Message     string    `json:"message" validate:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type LogResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type LogCreateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      string `json:"id"`
}
