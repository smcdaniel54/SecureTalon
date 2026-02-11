package core

import (
	"sync"
	"time"
)

// Store holds sessions, messages, and runs in memory (MVP).
// Thread-safe.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	messages map[string][]Message // session_id -> messages
	runs     map[string]*Run
}

// NewStore creates an empty in-memory store.
func NewStore() *Store {
	return &Store{
		sessions: make(map[string]*Session),
		messages: make(map[string][]Message),
		runs:     make(map[string]*Run),
	}
}

// CreateSession adds a new session.
func (s *Store) CreateSession(label string, metadata map[string]string) *Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess := &Session{
		ID:        NewSessionID(),
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Status:    "active",
		Metadata:  metadata,
	}
	s.sessions[sess.ID] = sess
	s.messages[sess.ID] = nil
	return sess
}

// GetSession returns a session by ID or nil.
func (s *Store) GetSession(id string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[id]
}

// ListSessions returns sessions, optionally with cursor/limit (MVP: simple limit).
func (s *Store) ListSessions(limit int, cursor string) ([]*Session, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Session
	for _, sess := range s.sessions {
		out = append(out, sess)
	}
	// Sort by created_at desc (simplified: no sort, just limit)
	if limit <= 0 {
		limit = 50
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out, ""
}

// AppendMessage adds a message to a session and returns it.
func (s *Store) AppendMessage(sessionID string, role, content string, metadata map[string]string) (*Message, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[sessionID]; !ok {
		return nil, false
	}
	msg := Message{
		ID:        NewMessageID(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now().UTC(),
		Metadata:  metadata,
	}
	s.messages[sessionID] = append(s.messages[sessionID], msg)
	return &msg, true
}

// SetMessageRunID sets run_id on the last message of the session (the one that triggered the run).
func (s *Store) SetMessageRunID(sessionID, runID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	msgs := s.messages[sessionID]
	if len(msgs) == 0 {
		return
	}
	msgs[len(msgs)-1].RunID = runID
	s.messages[sessionID] = msgs
}

// GetMessages returns messages for a session (limit applied).
func (s *Store) GetMessages(sessionID string, limit int) ([]Message, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs, ok := s.messages[sessionID]
	if !ok {
		return nil, false
	}
	if limit <= 0 {
		limit = 200
	}
	start := len(msgs) - limit
	if start < 0 {
		start = 0
	}
	return msgs[start:], true
}

// CreateRun creates a new run for a session (queued).
func (s *Store) CreateRun(sessionID string) *Run {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := &Run{
		ID:        NewRunID(),
		SessionID: sessionID,
		Status:    "queued",
		StartedAt: time.Now().UTC(),
		Steps:     nil,
	}
	s.runs[r.ID] = r
	return r
}

// GetRun returns a run by ID.
func (s *Store) GetRun(id string) *Run {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.runs[id]
}

// UpdateRunStatus sets status and optionally ended_at and steps.
func (s *Store) UpdateRunStatus(id string, status string, endedAt *time.Time, steps []Step) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r := s.runs[id]; r != nil {
		r.Status = status
		if endedAt != nil {
			r.EndedAt = endedAt
		}
		if steps != nil {
			r.Steps = steps
		}
	}
}

// AppendRunStep appends a step to a run.
func (s *Store) AppendRunStep(runID string, step Step) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r := s.runs[runID]; r != nil {
		r.Steps = append(r.Steps, step)
	}
}
