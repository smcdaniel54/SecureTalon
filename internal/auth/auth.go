// Package auth provides Bearer token validation for /v1/* endpoints.
// MVP: single ADMIN_TOKEN; no RBAC.
package auth

import (
	"net/http"
	"strings"
)

// Middleware returns a handler that requires Authorization: Bearer <ADMIN_TOKEN>.
// Missing or invalid token returns 401 with JSON error.
func Middleware(adminToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if adminToken == "" {
				writeJSONError(w, http.StatusInternalServerError, "AUTH_MISCONFIGURED", "ADMIN_TOKEN not configured")
				return
			}
			auth := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if !strings.HasPrefix(auth, prefix) {
				writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid Authorization header")
				return
			}
			token := strings.TrimSpace(auth[len(prefix):])
			if token != adminToken {
				writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// writeJSONError sends a standard error body. Kept minimal to avoid circular deps on api package.
func writeJSONError(w http.ResponseWriter, code int, errCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// Standard error format from API spec
	w.Write([]byte(`{"error":{"code":"` + errCode + `","message":"` + escapeJSON(message) + `","details":{}}}`))
}

func escapeJSON(s string) string {
	var b []byte
	for _, r := range s {
		switch r {
		case '"', '\\':
			b = append(b, '\\', byte(r))
		case '\n':
			b = append(b, '\\', 'n')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		default:
			if r < 32 {
				continue
			}
			b = append(b, byte(r))
		}
	}
	return string(b)
}
