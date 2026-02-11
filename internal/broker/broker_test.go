package broker

import (
	"testing"

	"securetalon/internal/core"
	"securetalon/internal/policy"
)

func TestBrokerRejectsNoToken(t *testing.T) {
	v := policy.NewVerifier("secret")
	b := NewBroker(v)
	_, err := b.Execute(core.ToolIntent{Tool: "file.read", Params: map[string]interface{}{"path": "/work/foo"}}, nil)
	if err == nil {
		t.Fatal("expected error when token is nil")
	}
}

func TestBrokerConstraintEnforcement(t *testing.T) {
	secret := "secret"
	issuer := policy.NewIssuer(secret)
	verifier := policy.NewVerifier(secret)
	// Token allows only /work/allowed/
	tok, _ := issuer.Issue("sess_1", "agent", "file.read", map[string]interface{}{
		"roots": []interface{}{"/work/allowed"},
		"max_bytes": 1000,
	}, 60)

	b := NewBroker(verifier)
	// /etc/passwd must be denied (not under allowed root)
	_, err := b.Execute(core.ToolIntent{
		Tool:   "file.read",
		Params: map[string]interface{}{"path": "/etc/passwd"},
	}, tok)
	if err == nil {
		t.Fatal("expected constraint violation for /etc/passwd")
	}
}
