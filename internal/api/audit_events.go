// Package api: helpers to append audit events for session/message/run lifecycle.
package api

import (
	"securetalon/internal/core"
)

func (h *Handlers) emitSessionCreated(sess *core.Session) {
	if h.AuditStore == nil {
		return
	}
	ev := &core.AuditEvent{
		SessionID: sess.ID,
		Type:      "session.created",
		Data: map[string]interface{}{
			"label":    sess.Label,
			"status":   sess.Status,
			"metadata": sess.Metadata,
		},
	}
	_ = h.AuditStore.Append(ev)
}

func (h *Handlers) emitMessageAppended(sessionID string, msg *core.Message, runID string) {
	if h.AuditStore == nil {
		return
	}
	data := map[string]interface{}{
		"message_id": msg.ID,
		"role":       msg.Role,
		"content":    msg.Content,
	}
	if runID != "" {
		data["run_id"] = runID
	}
	ev := &core.AuditEvent{
		SessionID: sessionID,
		RunID:     runID,
		Type:      "message.appended",
		Data:      data,
	}
	_ = h.AuditStore.Append(ev)
}

func (h *Handlers) emitRunStarted(run *core.Run) {
	if h.AuditStore == nil {
		return
	}
	ev := &core.AuditEvent{
		SessionID: run.SessionID,
		RunID:     run.ID,
		Type:      "run.started",
		Data: map[string]interface{}{
			"status": run.Status,
		},
	}
	_ = h.AuditStore.Append(ev)
}

func (h *Handlers) emitRunFinished(runID string, sessionID string, status string) {
	if h.AuditStore == nil {
		return
	}
	ev := &core.AuditEvent{
		SessionID: sessionID,
		RunID:     runID,
		Type:      "run.finished",
		Data: map[string]interface{}{
			"status": status,
		},
	}
	_ = h.AuditStore.Append(ev)
}
