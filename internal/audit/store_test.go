package audit

import (
	"testing"
	"time"

	"securetalon/internal/core"
)

func TestAuditHashChain(t *testing.T) {
	dir := t.TempDir()
	store, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	ev1 := &core.AuditEvent{
		SessionID: "sess_1",
		Type:      "session.created",
		Data:      map[string]interface{}{"label": "test"},
	}
	if err := store.Append(ev1); err != nil {
		t.Fatal(err)
	}
	if ev1.PrevHash != "" || ev1.Hash == "" {
		t.Fatalf("first event: prev_hash=%q hash=%q", ev1.PrevHash, ev1.Hash)
	}

	ev2 := &core.AuditEvent{
		SessionID: "sess_1",
		RunID:     "run_1",
		Type:      "run.started",
		Data:      map[string]interface{}{},
	}
	if err := store.Append(ev2); err != nil {
		t.Fatal(err)
	}
	if ev2.PrevHash != ev1.Hash {
		t.Fatalf("chain broken: ev2.prev_hash=%q ev1.hash=%q", ev2.PrevHash, ev1.Hash)
	}

	events, _, _ := store.Query("sess_1", "", "", "", "", 10)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	invalidIdx, err := ValidateChain(events)
	if err != nil {
		t.Fatal(err)
	}
	if invalidIdx >= 0 {
		t.Fatalf("chain invalid at index %d", invalidIdx)
	}
}

func TestValidateChainDetectsTampering(t *testing.T) {
	// Build two events manually; corrupt the second hash and validate
	ev1 := core.AuditEvent{
		EventID:   "evt_1",
		Timestamp: time.Now().UTC(),
		SessionID: "sess_1",
		Type:      "policy.decision",
		Data:      map[string]interface{}{"decision": "ALLOW"},
		PrevHash:  "",
	}
	ev1.Hash = chainHash("", &ev1)
	ev2 := core.AuditEvent{
		EventID:   "evt_2",
		Timestamp: time.Now().UTC(),
		SessionID: "sess_1",
		Type:      "tool.executed",
		Data:      map[string]interface{}{},
		PrevHash:  ev1.Hash,
	}
	ev2.Hash = "tampered"
	events := []core.AuditEvent{ev1, ev2}
	invalidIdx, _ := ValidateChain(events)
	if invalidIdx != 1 {
		t.Fatalf("expected invalid at index 1, got %d", invalidIdx)
	}
}

