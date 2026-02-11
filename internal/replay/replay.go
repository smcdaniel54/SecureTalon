// Package replay provides safe replay: reconstruct run timeline from audit events without re-executing tools.
package replay

import (
	"securetalon/internal/audit"
	"securetalon/internal/core"
)

// Replay returns the timeline of events for a run (from audit) and validates the hash chain.
// Safe: no tool execution; outputs are from recorded audit events.
func Replay(store *audit.Store, runID string) (events []core.AuditEvent, valid bool, err error) {
	events, _, err = store.Query("", runID, "", "", "", 0)
	if err != nil {
		return nil, false, err
	}
	invalidIdx, err := audit.ValidateChain(events)
	if err != nil {
		return nil, false, err
	}
	return events, invalidIdx < 0, nil
}
