// Package policy implements deny-by-default evaluation and capability token issuance.
// Security invariant: only the Policy Engine issues tokens; only the Tool Broker verifies them.
package policy

import (
	"securetalon/internal/core"
)

// SessionPolicy holds per-session overrides (allowlist rules).
type SessionPolicy struct {
	Overrides []RuleOverride `json:"overrides"`
}

// RuleOverride is one allowlist rule (e.g. file.read under path, http.fetch to domain).
type RuleOverride struct {
	Tool       string                 `json:"tool"`
	Allow      bool                   `json:"allow"`
	Constraints map[string]interface{} `json:"constraints"`
}

// Engine evaluates ToolIntent against static + session policy and returns ALLOW+token or DENY.
type Engine struct {
	DefaultTTL       int64
	SessionOverrides map[string]*SessionPolicy
	Issuer           *Issuer
}

// NewEngine returns a deny-by-default engine. Pass issuer so ALLOW results include a signed token.
func NewEngine(issuer *Issuer) *Engine {
	return &Engine{
		DefaultTTL:       60,
		SessionOverrides: make(map[string]*SessionPolicy),
		Issuer:           issuer,
	}
}

// SetSessionPolicy sets overrides for a session.
func (e *Engine) SetSessionPolicy(sessionID string, sp *SessionPolicy) {
	if e.SessionOverrides == nil {
		e.SessionOverrides = make(map[string]*SessionPolicy)
	}
	e.SessionOverrides[sessionID] = sp
}

// Evaluate returns ALLOW + token or DENY + reason. SessionContext can be nil for MVP.
func (e *Engine) Evaluate(intent core.ToolIntent, sessionID string) core.PolicyResult {
	// Deny shell.exec by default (no allowlist in MVP)
	if intent.Tool == "shell.exec" {
		return core.PolicyResult{
			Decision:     core.DecisionDeny,
			Reason:       "shell disabled by default",
			SuggestedFix: "Use file.read/file.write or docker.run instead",
		}
	}

	// Check session overrides for an explicit allow
	overrides := e.SessionOverrides[sessionID]
	if overrides != nil {
		for _, r := range overrides.Overrides {
			if r.Tool == intent.Tool && r.Allow && r.Constraints != nil {
				// Issue token with those constraints
				return e.allowWithConstraints(intent, sessionID, r.Constraints)
			}
		}
	}

	// Deny by default
	return core.PolicyResult{
		Decision:     core.DecisionDeny,
		Reason:       "Tool intent denied by policy (no matching allowlist)",
		SuggestedFix: "Add a session policy override allowing this tool with constraints",
	}
}

// allowWithConstraints issues a capability token with the given constraints.
func (e *Engine) allowWithConstraints(intent core.ToolIntent, sessionID string, constraints map[string]interface{}) core.PolicyResult {
	var token *core.CapabilityToken
	if e.Issuer != nil {
		subject := intent.Subject
		if subject == "" {
			subject = "agent"
		}
		var err error
		token, err = e.Issuer.Issue(sessionID, subject, intent.Tool, constraints, e.DefaultTTL)
		if err != nil {
			return core.PolicyResult{
				Decision: core.DecisionDeny,
				Reason:   "failed to issue capability token: " + err.Error(),
			}
		}
	}
	return core.PolicyResult{
		Decision: core.DecisionAllow,
		Reason:   "matched allowlist",
		Token:    token,
	}
}
