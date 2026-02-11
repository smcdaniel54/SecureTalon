package api

import (
	"encoding/json"
	"net/http"

	"securetalon/internal/agent"
	"securetalon/internal/audit"
	"securetalon/internal/core"
	"securetalon/internal/policy"
	"securetalon/internal/replay"
)

// Handlers holds dependencies for HTTP handlers.
type Handlers struct {
	Store       *core.Store
	Policy      *policy.Engine
	AuditStore  *audit.Store
	Agent       *agent.Agent
}

// CreateSession handles POST /v1/sessions
func (h *Handlers) CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "POST required", nil)
		return
	}
	var body struct {
		Label    string            `json:"label"`
		Metadata map[string]string `json:"metadata"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}
	sess := h.Store.CreateSession(body.Label, body.Metadata)
	h.emitSessionCreated(sess)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sess)
}

// ListSessions handles GET /v1/sessions
func (h *Handlers) ListSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET required", nil)
		return
	}
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := parseInt(l); n > 0 && n <= 200 {
			limit = n
		}
	}
	cursor := r.URL.Query().Get("cursor")
	_ = cursor
	list, _ := h.Store.ListSessions(limit, cursor)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"sessions": list})
}

// GetSession handles GET /v1/sessions/{id}
func (h *Handlers) GetSession(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET required", nil)
		return
	}
	sess := h.Store.GetSession(sessionID)
	if sess == nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Session not found", map[string]interface{}{"session_id": sessionID})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sess)
}

// PostMessage handles POST /v1/sessions/{id}/messages (starts a run; returns run_id)
func (h *Handlers) PostMessage(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "POST required", nil)
		return
	}
	sess := h.Store.GetSession(sessionID)
	if sess == nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Session not found", map[string]interface{}{"session_id": sessionID})
		return
	}
	var body struct {
		Role     string              `json:"role"`
		Content  string              `json:"content"`
		Metadata map[string]string   `json:"metadata"`
		Intents  []core.ToolIntent   `json:"intents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}
	if body.Role == "" {
		body.Role = "user"
	}
	msg, ok := h.Store.AppendMessage(sessionID, body.Role, body.Content, body.Metadata)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "INTERNAL", "Failed to append message", nil)
		return
	}
	run := h.Store.CreateRun(sessionID)
	h.Store.SetMessageRunID(sessionID, run.ID)
	h.emitMessageAppended(sessionID, msg, run.ID)
	h.emitRunStarted(run)
	if h.Agent != nil {
		go h.Agent.Run(sessionID, run.ID, body.Intents)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"run_id": run.ID,
		"status": run.Status,
	})
}

// ListMessages handles GET /v1/sessions/{id}/messages
func (h *Handlers) ListMessages(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET required", nil)
		return
	}
	limit := 200
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := parseInt(l); n > 0 && n <= 500 {
			limit = n
		}
	}
	msgs, ok := h.Store.GetMessages(sessionID, limit)
	if !ok {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Session not found", map[string]interface{}{"session_id": sessionID})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"messages": msgs})
}

// GetRun handles GET /v1/runs/{run_id}
func (h *Handlers) GetRun(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET required", nil)
		return
	}
	run := h.Store.GetRun(runID)
	if run == nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Run not found", map[string]interface{}{"run_id": runID})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(run)
}

// PutSessionPolicy handles PUT /v1/sessions/{id}/policy
func (h *Handlers) PutSessionPolicy(w http.ResponseWriter, r *http.Request, sessionID string) {
	if h.Store.GetSession(sessionID) == nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Session not found", map[string]interface{}{"session_id": sessionID})
		return
	}
	var body struct {
		Overrides []policy.RuleOverride `json:"overrides"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}
	if h.Policy != nil {
		h.Policy.SetSessionPolicy(sessionID, &policy.SessionPolicy{Overrides: body.Overrides})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

// GetEffectivePolicy handles GET /v1/policy/effective?session_id=...
func (h *Handlers) GetEffectivePolicy(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID != "" && h.Store.GetSession(sessionID) == nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Session not found", map[string]interface{}{"session_id": sessionID})
		return
	}
	overrides := []interface{}{}
	if h.Policy != nil && sessionID != "" {
		if sp := h.Policy.SessionOverrides[sessionID]; sp != nil {
			for _, r := range sp.Overrides {
				overrides = append(overrides, r)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"default":   "deny",
		"overrides": overrides,
	})
}

// ListSkills handles GET /v1/skills
func (h *Handlers) ListSkills(w http.ResponseWriter, r *http.Request) {
	// Stub: empty list until skills package is wired
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"skills": []interface{}{}})
}

// RegisterSkill handles POST /v1/skills
func (h *Handlers) RegisterSkill(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}
	// Stub: accept and return 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
}

// QueryAudit handles GET /v1/audit?session_id=...&since=...&until=...&type=...&limit=500
func (h *Handlers) QueryAudit(w http.ResponseWriter, r *http.Request) {
	limit := 500
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := parseInt(l); n > 0 && n <= 1000 {
			limit = n
		}
	}
	sessionID := r.URL.Query().Get("session_id")
	since := r.URL.Query().Get("since")
	until := r.URL.Query().Get("until")
	eventType := r.URL.Query().Get("type")
	var events []core.AuditEvent
	if h.AuditStore != nil {
		events, _, _ = h.AuditStore.Query(sessionID, "", since, until, eventType, limit)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events":      events,
		"next_cursor": nil,
	})
}

// ValidateAuditChain handles GET /v1/audit/validate?session_id=...&limit=...
func (h *Handlers) ValidateAuditChain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/audit/validate" || r.Method != http.MethodGet {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Not found", nil)
		return
	}
	limit := 500
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, _ := parseInt(l); n > 0 && n <= 1000 {
			limit = n
		}
	}
	sessionID := r.URL.Query().Get("session_id")
	var events []core.AuditEvent
	if h.AuditStore != nil {
		events, _, _ = h.AuditStore.Query(sessionID, "", "", "", "", limit)
	}
	invalidIdx, _ := audit.ValidateChain(events)
	valid := invalidIdx < 0
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":        valid,
		"invalid_index": invalidIdx,
		"event_count":  len(events),
	})
}

// GetReplay handles GET /v1/replay?run_id=... (safe replay: no tool re-execution).
func (h *Handlers) GetReplay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET required", nil)
		return
	}
	runID := r.URL.Query().Get("run_id")
	if runID == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "run_id required", nil)
		return
	}
	if h.AuditStore == nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL", "Audit store not available", nil)
		return
	}
	events, valid, err := replay.Replay(h.AuditStore, runID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"run_id": runID,
		"valid":  valid,
		"events": events,
	})
}

// PostRunReplay handles POST /v1/runs/{run_id}/replay?safe=true. Returns timeline for UI (safe mode, no re-exec).
func (h *Handlers) PostRunReplay(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "POST required", nil)
		return
	}
	if h.AuditStore == nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL", "Audit store not available", nil)
		return
	}
	events, valid, err := replay.Replay(h.AuditStore, runID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL", err.Error(), nil)
		return
	}
	timeline := make([]map[string]interface{}, 0, len(events))
	for _, ev := range events {
		timeline = append(timeline, map[string]interface{}{
			"ts":        ev.Timestamp,
			"type":      ev.Type,
			"data":      ev.Data,
			"hash":      ev.Hash,
			"prev_hash": ev.PrevHash,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"run_id": runID,
		"mode":   "safe",
		"valid":  valid,
		"events": timeline,
	})
}
