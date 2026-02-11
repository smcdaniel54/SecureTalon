// Package broker is the only component that executes tool operations (file, http, docker).
// It verifies capability tokens and enforces constraints. No tool runs without a valid token.
package broker

import (
	"fmt"
	"securetalon/internal/core"
	"securetalon/internal/policy"
)

// Broker executes tool intents after verifying the capability token and constraints.
type Broker struct {
	Verifier *policy.Verifier
}

// NewBroker returns a broker that uses the given verifier.
func NewBroker(v *policy.Verifier) *Broker {
	return &Broker{Verifier: v}
}

// Execute verifies the token and runs the tool. Returns result or error.
// Constraint enforcement: e.g. file.read only under allowed roots, block /etc/passwd.
func (b *Broker) Execute(intent core.ToolIntent, token *core.CapabilityToken) (result map[string]interface{}, err error) {
	if token == nil {
		return nil, fmt.Errorf("capability token required")
	}
	if err := b.Verifier.Verify(token); err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if token.Tool != intent.Tool {
		return nil, fmt.Errorf("token tool mismatch")
	}
	// Enforce constraints against intent params
	if err := b.checkConstraints(intent, token); err != nil {
		return nil, err
	}
	switch intent.Tool {
	case "file.read":
		return b.doFileRead(intent.Params, token.Constraints)
	case "file.write":
		return b.doFileWrite(intent.Params, token.Constraints)
	case "http.fetch":
		return b.doHTTPFetch(intent.Params, token.Constraints)
	case "docker.run":
		return b.doDockerRun(intent.Params, token.Constraints)
	case "shell.exec":
		return nil, fmt.Errorf("shell.exec disabled by default")
	default:
		return nil, fmt.Errorf("unknown tool: %s", intent.Tool)
	}
}

// checkConstraints enforces token constraints on intent params (e.g. path under allowed root).
func (b *Broker) checkConstraints(intent core.ToolIntent, token *core.CapabilityToken) error {
	tool := intent.Tool
	params := intent.Params
	constraints := token.Constraints
	if constraints == nil {
		return fmt.Errorf("no constraints on token")
	}
	switch tool {
	case "file.read", "file.write":
		path, _ := params["path"].(string)
		if path == "" {
			return fmt.Errorf("path required")
		}
		allowedRoots := constraints["roots"]
		if allowedRoots == nil {
			return fmt.Errorf("constraint roots required for file access")
		}
		if !pathUnderAllowedRoots(path, allowedRoots) {
			return fmt.Errorf("path %s not under allowed roots (e.g. block /etc/passwd)", path)
		}
	case "http.fetch":
		// domain allowlist checked in doHTTPFetch
	case "docker.run":
		// image digest allowlist checked in doDockerRun
	}
	return nil
}

// pathUnderAllowedRoots returns true if path is under one of the allowed root prefixes.
func pathUnderAllowedRoots(path string, allowedRoots interface{}) bool {
	// allowedRoots can be []interface{} from JSON
	switch v := allowedRoots.(type) {
	case []string:
		for _, root := range v {
			if pathUnder(path, root) {
				return true
			}
		}
	case []interface{}:
		for _, r := range v {
			if s, ok := r.(string); ok && pathUnder(path, s) {
				return true
			}
		}
	}
	return false
}

func pathUnder(path, root string) bool {
	if root == "" {
		return false
	}
	if len(root) > len(path) {
		return false
	}
	if path[:len(root)] != root {
		return false
	}
	if len(path) == len(root) {
		return true
	}
	return path[len(root)] == '/' || path[len(root)] == '\\'
}
