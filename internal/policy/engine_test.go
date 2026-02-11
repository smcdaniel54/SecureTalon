package policy

import (
	"testing"

	"securetalon/internal/core"
)

func TestDenyByDefault(t *testing.T) {
	issuer := NewIssuer("test-secret")
	engine := NewEngine(issuer)

	// shell.exec always denied
	result := engine.Evaluate(core.ToolIntent{Tool: "shell.exec", Params: map[string]interface{}{}}, "sess_1")
	if result.Decision != core.DecisionDeny {
		t.Fatalf("expected DENY for shell.exec, got %s", result.Decision)
	}

	// file.read without allowlist denied
	result = engine.Evaluate(core.ToolIntent{Tool: "file.read", Params: map[string]interface{}{"path": "/work/foo"}}, "sess_1")
	if result.Decision != core.DecisionDeny {
		t.Fatalf("expected DENY for file.read without allowlist, got %s", result.Decision)
	}
}

func TestAllowWithOverride(t *testing.T) {
	issuer := NewIssuer("test-secret")
	engine := NewEngine(issuer)
	engine.SetSessionPolicy("sess_1", &SessionPolicy{
		Overrides: []RuleOverride{
			{Tool: "file.read", Allow: true, Constraints: map[string]interface{}{
				"roots":     []string{"/work/allowed"},
				"max_bytes": 1024,
			}},
		},
	})

	result := engine.Evaluate(core.ToolIntent{
		Tool:   "file.read",
		Params: map[string]interface{}{"path": "/work/allowed/foo"},
	}, "sess_1")
	if result.Decision != core.DecisionAllow {
		t.Fatalf("expected ALLOW, got %s", result.Decision)
	}
	if result.Token == nil {
		t.Fatal("expected capability token")
	}
	if result.Token.Tool != "file.read" {
		t.Fatalf("token tool: got %s", result.Token.Tool)
	}
}
