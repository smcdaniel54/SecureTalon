// Package core defines session, message, run, and capability types.
// Security invariant: only Policy Engine issues tokens; only Tool Broker verifies and executes.
package core

import "time"

// Session represents a conversation/agent session.
type Session struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Status    string            `json:"status"` // active, closed
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Message is a single message in a session.
type Message struct {
	ID        string            `json:"id"`
	Role      string            `json:"role"` // user, assistant, system
	Content   string            `json:"content"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	RunID     string            `json:"run_id,omitempty"` // set when message triggers a run
}

// Run represents an agent run (triggered by a user message).
type Run struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Status    string    `json:"status"` // queued, running, completed, failed
	StartedAt time.Time `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Steps     []Step    `json:"steps,omitempty"`
}

// Step is one step in a run (policy eval or tool exec).
type Step struct {
	StepID   string                 `json:"step_id"`
	Type     string                 `json:"type"` // policy_eval, tool_exec
	Status   string                 `json:"status"`
	Tool     string                 `json:"tool,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// ToolIntent is a request to execute a tool (from agent/skill).
type ToolIntent struct {
	Tool     string                 `json:"tool"`
	Params   map[string]interface{} `json:"params"`
	Subject  string                 `json:"subject,omitempty"`
}

// Decision is the result of policy evaluation.
type Decision string

const (
	DecisionAllow          Decision = "ALLOW"
	DecisionDeny           Decision = "DENY"
	DecisionRequireApproval Decision = "REQUIRE_APPROVAL"
)

// PolicyResult is returned by Policy Engine.
type PolicyResult struct {
	Decision     Decision        `json:"decision"`
	Reason       string          `json:"reason,omitempty"`
	SuggestedFix string          `json:"suggested_fix,omitempty"`
	Token        *CapabilityToken `json:"token,omitempty"`
}

// CapabilityToken is a signed, short-lived grant for one tool action.
// Binds: session_id, subject, tool, constraints, exp.
type CapabilityToken struct {
	CapID      string                 `json:"cap_id"`
	SessionID  string                 `json:"session_id"`
	Subject    string                 `json:"subject"`
	Tool       string                 `json:"tool"`
	Constraints map[string]interface{} `json:"constraints"`
	Iat        int64                  `json:"iat"`
	Exp        int64                  `json:"exp"`
	Nonce      string                 `json:"nonce"`
	Signature  string                 `json:"signature"`
}

// AuditEvent is one entry in the append-only audit log (hash-chained).
type AuditEvent struct {
	EventID   string                 `json:"event_id"`
	Timestamp time.Time              `json:"ts"`
	SessionID string                 `json:"session_id"`
	RunID     string                 `json:"run_id,omitempty"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	PrevHash  string                 `json:"prev_hash"`
	Hash      string                 `json:"hash"`
}
