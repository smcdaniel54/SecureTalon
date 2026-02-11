package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// NewID returns a short prefixed ID (e.g. sess_abc123).
func NewID(prefix string) string {
	b := make([]byte, 8)
	rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}

// NewSessionID returns sess_...
func NewSessionID() string { return NewID("sess") }

// NewMessageID returns msg_...
func NewMessageID() string { return NewID("msg") }

// NewRunID returns run_...
func NewRunID() string { return NewID("run") }

// NewEventID returns evt_...
func NewEventID() string { return NewID("evt") }

// NewCapID returns cap_...
func NewCapID() string { return NewID("cap") }

// NewStepID returns s1, s2, ...
func NewStepID(n int) string { return fmt.Sprintf("s%d", n) }
