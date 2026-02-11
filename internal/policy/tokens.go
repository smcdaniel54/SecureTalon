// Package policy provides capability token signing and verification.
// HMAC-SHA256 using server secret for MVP.
package policy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"securetalon/internal/core"
	"time"
)

// Issuer signs capability tokens.
type Issuer struct {
	Secret []byte
}

// NewIssuer creates an issuer with the given secret (e.g. from config).
func NewIssuer(secret string) *Issuer {
	if secret == "" {
		secret = "default-secret-change-in-production"
	}
	return &Issuer{Secret: []byte(secret)}
}

// Issue creates a signed capability token for the given session, subject, tool, and constraints.
func (i *Issuer) Issue(sessionID, subject, tool string, constraints map[string]interface{}, ttlSeconds int64) (*core.CapabilityToken, error) {
	now := time.Now().UTC().Unix()
	if ttlSeconds <= 0 {
		ttlSeconds = 60
	}
	exp := now + ttlSeconds
	capID := core.NewCapID()
	nonce := fmt.Sprintf("%d-%s", now, capID)
	tok := &core.CapabilityToken{
		CapID:       capID,
		SessionID:   sessionID,
		Subject:     subject,
		Tool:        tool,
		Constraints: constraints,
		Iat:         now,
		Exp:         exp,
		Nonce:       nonce,
	}
	sig, err := signToken(tok, i.Secret)
	if err != nil {
		return nil, err
	}
	tok.Signature = sig
	return tok, nil
}

// Verifier checks token signature and expiry.
type Verifier struct {
	Secret []byte
}

// NewVerifier creates a verifier with the same secret as the issuer.
func NewVerifier(secret string) *Verifier {
	if secret == "" {
		secret = "default-secret-change-in-production"
	}
	return &Verifier{Secret: []byte(secret)}
}

// Verify returns nil if the token is valid and not expired.
func (v *Verifier) Verify(tok *core.CapabilityToken) error {
	if tok == nil {
		return fmt.Errorf("nil token")
	}
	now := time.Now().UTC().Unix()
	if tok.Exp < now {
		return fmt.Errorf("token expired")
	}
	if tok.Iat > now {
		return fmt.Errorf("token not yet valid")
	}
	expectedSig, err := signToken(tok, v.Secret)
	if err != nil {
		return err
	}
	if !hmac.Equal([]byte(tok.Signature), []byte(expectedSig)) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

// signToken computes HMAC-SHA256 of canonical token payload (without signature).
func signToken(tok *core.CapabilityToken, secret []byte) (string, error) {
	payload := struct {
		CapID       string                 `json:"cap_id"`
		SessionID   string                 `json:"session_id"`
		Subject     string                 `json:"subject"`
		Tool        string                 `json:"tool"`
		Constraints map[string]interface{} `json:"constraints"`
		Iat         int64                  `json:"iat"`
		Exp         int64                  `json:"exp"`
		Nonce       string                 `json:"nonce"`
	}{
		tok.CapID, tok.SessionID, tok.Subject, tok.Tool, tok.Constraints, tok.Iat, tok.Exp, tok.Nonce,
	}
	canon, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, secret)
	h.Write(canon)
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
