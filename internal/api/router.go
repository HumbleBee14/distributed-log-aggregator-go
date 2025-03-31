package api

import (
	"net/http"

	"github.com/HumbleBee14/distributed-log-aggregator/internal/storage"
)

// Router handles HTTP routing
type Router struct {
	handler *LogHandler
}

// NewRouter creates a new router
func NewRouter(storage *storage.RedisStorage) *Router {
	return &Router{
		handler: NewLogHandler(storage),
	}
}

// Setup sets up all routes
func (r *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/logs", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.handler.HandleStoreLog(w, req)
		case http.MethodGet:
			r.handler.HandleQueryLogs(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	mux.HandleFunc("/health", r.handler.HandleHealthCheck)

	handler := applyMiddleware(mux)

	return handler
}

func applyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// to handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// TODO: Logger

		next.ServeHTTP(w, r)
	})
}
