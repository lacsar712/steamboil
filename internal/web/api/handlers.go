package api

import (
	"encoding/json"
	"net/http"

	"github.com/lacsar712/steamboil/internal/model"
)

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", withTimeout(s.handleHealth))
	mux.HandleFunc("GET /snapshot", withTimeout(s.handleSnapshot))
	mux.HandleFunc("GET /telemetry", withTimeout(s.handleTelemetry))
	mux.HandleFunc("GET /health/plant", withTimeout(s.handlePlantHealth))
	mux.HandleFunc("GET /warmup", withTimeout(s.handleWarmup))
	mux.HandleFunc("POST /purge/start", withTimeout(s.handleStartPurge))
	mux.HandleFunc("POST /purge/complete", withTimeout(s.handleCompletePurge))
	mux.HandleFunc("POST /ignite", withTimeout(s.handleIgnite))
	mux.HandleFunc("POST /trip/reset", withTimeout(s.handleResetTrip))
	mux.HandleFunc("POST /settings", withTimeout(s.handleSettings))
	return cors(logging(mux))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]string{"status": "ok", "unit": s.app.UnitID()})
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	writeOK(w, s.app.Snapshot())
}

func (s *Server) handleTelemetry(w http.ResponseWriter, r *http.Request) {
	writeOK(w, s.app.Telemetry())
}

func (s *Server) handlePlantHealth(w http.ResponseWriter, r *http.Request) {
	writeOK(w, s.app.PlantHealth())
}

func (s *Server) handleWarmup(w http.ResponseWriter, r *http.Request) {
	ready, detail := s.app.WarmupStatus()
	writeOK(w, map[string]interface{}{
		"ready":  ready,
		"detail": detail,
		"purge":  s.app.PurgeRemaining(),
		"warmup": s.app.CombustionWarmupRemaining(),
	})
}

func (s *Server) handleStartPurge(w http.ResponseWriter, r *http.Request) {
	holder := r.URL.Query().Get("holder")
	if holder == "" {
		holder = "hmi-operator"
	}
	if err := s.app.StartPurge(r.Context(), holder); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeOK(w, s.app.Snapshot())
}

func (s *Server) handleCompletePurge(w http.ResponseWriter, r *http.Request) {
	holder := r.URL.Query().Get("holder")
	if holder == "" {
		holder = "hmi-operator"
	}
	if err := s.app.CompletePurge(r.Context(), holder); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeOK(w, s.app.Snapshot())
}

func (s *Server) handleIgnite(w http.ResponseWriter, r *http.Request) {
	holder := r.URL.Query().Get("holder")
	if holder == "" {
		holder = "hmi-operator"
	}
	if err := s.app.Ignite(r.Context(), holder); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeOK(w, s.app.Snapshot())
}

func (s *Server) handleResetTrip(w http.ResponseWriter, r *http.Request) {
	holder := r.URL.Query().Get("holder")
	if holder == "" {
		holder = "hmi-operator"
	}
	if err := s.app.ResetTrip(r.Context(), holder); err != nil {
		writeErr(w, http.StatusConflict, err)
		return
	}
	writeOK(w, s.app.Snapshot())
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	var settings model.PlantSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if err := s.app.UpdateSettings(r.Context(), settings); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeOK(w, s.app.Snapshot())
}
