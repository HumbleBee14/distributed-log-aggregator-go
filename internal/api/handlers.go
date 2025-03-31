package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/HumbleBee14/distributed-log-aggregator/internal/models"
	"github.com/HumbleBee14/distributed-log-aggregator/internal/storage"
	"github.com/go-playground/validator/v10"
)

// LogHandler handles the log-related API endpoints
type LogHandler struct {
	storage  *storage.RedisStorage
	validate *validator.Validate
}

// NewLogHandler creates a new log handler
func NewLogHandler(storage *storage.RedisStorage) *LogHandler {
	return &LogHandler{
		storage:  storage,
		validate: validator.New(),
	}
}

// HandleStoreLog handles log ingestion (POST /logs)
func (h *LogHandler) HandleStoreLog(w http.ResponseWriter, r *http.Request) {
	// Only accept POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var logEntry models.LogEntry
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the log entry
	if err := h.validate.Struct(logEntry); err != nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Store the log entry
	id, err := h.storage.StoreLog(logEntry)
	if err != nil {
		http.Error(w, "Failed to store log entry", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := models.LogCreateResponse{
		Status:  "success",
		Message: "Log entry created",
		ID:      id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// HandleQueryLogs handles log querying (GET /logs)
func (h *LogHandler) HandleQueryLogs(w http.ResponseWriter, r *http.Request) {
	// Only accept GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	serviceName := r.URL.Query().Get("service")
	if serviceName == "" {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	startTimeStr := r.URL.Query().Get("start")
	endTimeStr := r.URL.Query().Get("end")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		t, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start time format", http.StatusBadRequest)
			return
		}
		startTime = &t
	}

	if endTimeStr != "" {
		t, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end time format", http.StatusBadRequest)
			return
		}
		endTime = &t
	}

	// Query logs
	logs, err := h.storage.QueryLogs(serviceName, startTime, endTime)
	if err != nil {
		http.Error(w, "Failed to query logs", http.StatusInternalServerError)
		return
	}

	// Return the logs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// HandleHealthCheck handles the health check endpoint
func (h *LogHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "UP",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
