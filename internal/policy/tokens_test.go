package policy

import (
	"testing"
	"time"
)

func TestIssueAndVerify(t *testing.T) {
	secret := "test-secret"
	issuer := NewIssuer(secret)
	verifier := NewVerifier(secret)

	tok, err := issuer.Issue("sess_1", "agent", "file.read", map[string]interface{}{"roots": []string{"/work"}}, 60)
	if err != nil {
		t.Fatal(err)
	}
	if tok.CapID == "" || tok.Signature == "" {
		t.Fatal("token missing cap_id or signature")
	}
	if err := verifier.Verify(tok); err != nil {
		t.Fatalf("verify failed: %v", err)
	}
}

func TestVerifyRejectsWrongSecret(t *testing.T) {
	issuer := NewIssuer("secret-a")
	verifier := NewVerifier("secret-b")

	tok, _ := issuer.Issue("sess_1", "agent", "file.read", nil, 60)
	if verifier.Verify(tok) == nil {
		t.Fatal("expected verify to fail with wrong secret")
	}
}

func TestVerifyRejectsExpired(t *testing.T) {
	secret := "test"
	issuer := NewIssuer(secret)
	verifier := NewVerifier(secret)

	tok, _ := issuer.Issue("sess_1", "agent", "file.read", nil, -1) // expired
	// Force exp in the past
	tok.Exp = time.Now().UTC().Unix() - 10
	if verifier.Verify(tok) == nil {
		t.Fatal("expected verify to fail for expired token")
	}
}
