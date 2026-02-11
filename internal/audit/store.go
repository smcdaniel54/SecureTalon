// Package audit provides append-only event log with hash chain for tamper evidence.
package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"securetalon/internal/core"
)

// Store is an append-only JSONL audit log with hash chaining.
type Store struct {
	mu      sync.Mutex
	dir     string
	prevHash string
}

// NewStore creates an audit store under dir (e.g. data/audit).
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	s := &Store{dir: dir}
	// Load last hash from existing file if any
	_ = s.loadLastHash()
	return s, nil
}

func (s *Store) loadLastHash() error {
	// MVP: single file audit.jsonl
	fpath := filepath.Join(s.dir, "audit.jsonl")
	f, err := os.Open(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			s.prevHash = ""
			return nil
		}
		return err
	}
	defer f.Close()
	var lastHash string
	dec := json.NewDecoder(f)
	for {
		var ev core.AuditEvent
		if err := dec.Decode(&ev); err != nil {
			break
		}
		lastHash = ev.Hash
	}
	s.prevHash = lastHash
	return nil
}

// Append writes one event with hash = sha256(prev_hash + canonical_json(event)).
func (s *Store) Append(ev *core.AuditEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ev.EventID == "" {
		ev.EventID = core.NewEventID()
	}
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}
	ev.PrevHash = s.prevHash
	ev.Hash = chainHash(s.prevHash, ev)
	s.prevHash = ev.Hash
	return s.appendLine(ev)
}

func (s *Store) appendLine(ev *core.AuditEvent) error {
	fpath := filepath.Join(s.dir, "audit.jsonl")
	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(ev)
}

// canonicalEvent returns a struct that serializes to the same format we hash.
type canonicalEvent struct {
	EventID   string                 `json:"event_id"`
	Timestamp time.Time              `json:"ts"`
	SessionID string                 `json:"session_id"`
	RunID     string                 `json:"run_id,omitempty"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	PrevHash  string                 `json:"prev_hash"`
}

func canonicalEventFrom(ev *core.AuditEvent) canonicalEvent {
	return canonicalEvent{
		EventID:   ev.EventID,
		Timestamp: ev.Timestamp,
		SessionID: ev.SessionID,
		RunID:     ev.RunID,
		Type:      ev.Type,
		Data:      ev.Data,
		PrevHash:  ev.PrevHash,
	}
}

// chainHash computes sha256(prevHash + canonical_json(ev)) without the hash field.
func chainHash(prevHash string, ev *core.AuditEvent) string {
	c := canonicalEventFrom(ev)
	c.PrevHash = prevHash
	raw, _ := json.Marshal(c)
	return sha256Hex(prevHash + string(raw))
}

// Query returns events matching filters. MVP: read full file and filter in memory.
func (s *Store) Query(sessionID, runID, since, until, eventType string, limit int) ([]core.AuditEvent, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fpath := filepath.Join(s.dir, "audit.jsonl")
	f, err := os.Open(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", nil
		}
		return nil, "", err
	}
	defer f.Close()
	var events []core.AuditEvent
	dec := json.NewDecoder(f)
	for dec.More() {
		var ev core.AuditEvent
		if err := dec.Decode(&ev); err != nil {
			break
		}
		if sessionID != "" && ev.SessionID != sessionID {
			continue
		}
		if runID != "" && ev.RunID != runID {
			continue
		}
		if eventType != "" && ev.Type != eventType {
			continue
		}
		if since != "" && ev.Timestamp.Format(time.RFC3339) < since {
			continue
		}
		if until != "" && ev.Timestamp.Format(time.RFC3339) > until {
			continue
		}
		events = append(events, ev)
		if limit > 0 && len(events) >= limit {
			break
		}
	}
	return events, "", nil
}

// ValidateChain verifies hash chain for a list of events. Returns first error index or -1.
func ValidateChain(events []core.AuditEvent) (int, error) {
	var prev string
	for i := range events {
		ev := &events[i]
		expected := chainHash(prev, ev)
		if ev.Hash != expected {
			return i, nil // caller can use index
		}
		prev = ev.Hash
	}
	return -1, nil
}
