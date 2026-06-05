package handlers

import (
	"database-scaling-lab/internal/repository"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests sorted by path.",
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
	}, []string{"path", "method"})
)

type Server struct {
	Repo *repository.PostgresRepository
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/users/", s.handleGetUser)
	mux.HandleFunc("/events", s.handleCreateEvent)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		httpDuration.WithLabelValues("/health", r.Method).Observe(time.Since(start).Seconds())
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := s.Repo.Pool.Ping(r.Context())
	if err != nil {
		http.Error(w, `{"status":"unhealthy"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"healthy"}`))
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		httpDuration.WithLabelValues("/users/:id", r.Method).Observe(time.Since(start).Seconds())
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Method split unsupported", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Malformed UUID parameter", http.StatusBadRequest)
		return
	}

	user, err := s.Repo.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, "User record lookup collision or failure", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		httpDuration.WithLabelValues("/events", r.Method).Observe(time.Since(start).Seconds())
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Method split unsupported", http.StatusMethodNotAllowed)
		return
	}

	// For rapid high-throughput generation testing, mock parameters are accepted via query
	uStr := r.URL.Query().Get("user_id")
	eType := r.URL.Query().Get("type")

	uID, err := uuid.Parse(uStr)
	if err != nil || eType == "" {
		http.Error(w, "Invalid parameters provided", http.StatusBadRequest)
		return
	}

	err = s.Repo.InsertEvent(r.Context(), uID, eType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
