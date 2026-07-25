package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/makarivmri12/bug/internal/engine"
)

// Server represents the HTTP API server
type Server struct {
	router    *chi.Mux
	logger    *zap.Logger
	scanEngine *engine.ScanEngine
}

// NewServer creates a new API server
func NewServer(logger *zap.Logger, scanEngine *engine.ScanEngine) *Server {
	s := &Server{
		router:     chi.NewRouter(),
		logger:     logger,
		scanEngine: scanEngine,
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"X-Total-Count"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Logging middleware
	s.router.Use(s.loggingMiddleware)

	// Health check
	s.router.Get("/health", s.handleHealth)

	// Scan endpoints
	s.router.Post("/api/v1/scans", s.handleCreateScan)
	s.router.Get("/api/v1/scans/{scanID}", s.handleGetScan)
	s.router.Get("/api/v1/scans/{scanID}/findings", s.handleGetFindings)

	// Target endpoints
	s.router.Post("/api/v1/targets", s.handleCreateTarget)
	s.router.Get("/api/v1/targets/{targetID}", s.handleGetTarget)

	// Finding endpoints
	s.router.Get("/api/v1/findings/{findingID}", s.handleGetFinding)
	s.router.Put("/api/v1/findings/{findingID}", s.handleUpdateFinding)
}

// loggingMiddleware logs HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("incoming request",
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.String("remote_addr", r.RemoteAddr),
		)
		next.ServeHTTP(w, r)
	})
}

// Health check handler
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// CreateScan handler
func (s *Server) handleCreateScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"scan_id":"scan-123","status":"running"}`))
}

// GetScan handler
func (s *Server) handleGetScan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"scan_id":"scan-123","status":"completed"}`))
}

// GetFindings handler
func (s *Server) handleGetFindings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"findings":[]}`))
}

// CreateTarget handler
func (s *Server) handleCreateTarget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"target_id":"target-123"}`))
}

// GetTarget handler
func (s *Server) handleGetTarget(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"target_id":"target-123","url":"https://example.com"}`))
}

// GetFinding handler
func (s *Server) handleGetFinding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"finding_id":"finding-123","severity":"HIGH"}`))
}

// UpdateFinding handler
func (s *Server) handleUpdateFinding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"finding_id":"finding-123","status":"confirmed"}`))
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
