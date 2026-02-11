package api

import (
	"net/http"
	"regexp"
	"strings"
)

// Router serves /v1/* with path params. Auth middleware must wrap this.
var (
	reSessionID = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	reRunID     = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// NewRouter returns an http.Handler that routes /v1/* to Handlers.
func NewRouter(h *Handlers) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/sessions", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sessions" {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
			return
		}
		if r.Method == http.MethodPost {
			h.CreateSession(w, r)
			return
		}
		if r.Method == http.MethodGet {
			h.ListSessions(w, r)
			return
		}
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
	})
	mux.HandleFunc("/v1/sessions/", func(w http.ResponseWriter, r *http.Request) {
		trimmed := strings.TrimPrefix(r.URL.Path, "/v1/sessions/")
		parts := strings.SplitN(trimmed, "/", 2)
		sessionID := parts[0]
		if !reSessionID.MatchString(sessionID) {
			WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid session id", nil)
			return
		}
		if len(parts) == 1 {
			if r.Method == http.MethodGet {
				h.GetSession(w, r, sessionID)
				return
			}
			if r.Method == http.MethodPut && sessionID != "" {
				// PUT /v1/sessions/{id}/policy handled below via longer path
				WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
				return
			}
			WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
			return
		}
		rest := parts[1]
		if rest == "messages" {
			if r.Method == http.MethodPost {
				h.PostMessage(w, r, sessionID)
				return
			}
			if r.Method == http.MethodGet {
				h.ListMessages(w, r, sessionID)
				return
			}
			WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
			return
		}
		if rest == "policy" && r.Method == http.MethodPut {
			h.PutSessionPolicy(w, r, sessionID)
			return
		}
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
	})
	mux.HandleFunc("/v1/runs/", func(w http.ResponseWriter, r *http.Request) {
		trimmed := strings.TrimPrefix(r.URL.Path, "/v1/runs/")
		parts := strings.SplitN(trimmed, "/", 2)
		runID := parts[0]
		if runID == "" || !reRunID.MatchString(runID) {
			WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid run id", nil)
			return
		}
		if len(parts) == 2 && parts[1] == "replay" {
			if r.Method == http.MethodPost {
				h.PostRunReplay(w, r, runID)
				return
			}
			WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "POST required for replay", nil)
			return
		}
		if r.Method == http.MethodGet {
			h.GetRun(w, r, runID)
			return
		}
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
	})
	mux.HandleFunc("/v1/policy/effective", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/policy/effective" {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
			return
		}
		if r.Method == http.MethodGet {
			h.GetEffectivePolicy(w, r)
			return
		}
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
	})
	mux.HandleFunc("/v1/skills", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/skills" {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
			return
		}
		if r.Method == http.MethodGet {
			h.ListSkills(w, r)
			return
		}
		if r.Method == http.MethodPost {
			h.RegisterSkill(w, r)
			return
		}
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
	})
	mux.HandleFunc("/v1/audit", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/audit/validate" {
			h.ValidateAuditChain(w, r)
			return
		}
		if r.URL.Path != "/v1/audit" {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
			return
		}
		if r.Method == http.MethodGet {
			h.QueryAudit(w, r)
			return
		}
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
	})
	mux.HandleFunc("/v1/replay", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/replay" {
			WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
			return
		}
		if r.Method == http.MethodGet {
			h.GetReplay(w, r)
			return
		}
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
	})
	return mux
}
