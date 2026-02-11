package api

import (
	"encoding/json"
	"net/http"
)

// ErrorBody matches the API spec: { "error": { "code", "message", "details" } }
type ErrorBody struct {
	Error struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details"`
	} `json:"error"`
}

// WriteError sends a standard JSON error response.
func WriteError(w http.ResponseWriter, code int, errCode, message string, details map[string]interface{}) {
	if details == nil {
		details = make(map[string]interface{})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorBody{
		Error: struct {
			Code    string                 `json:"code"`
			Message string                 `json:"message"`
			Details map[string]interface{} `json:"details"`
		}{errCode, message, details},
	})
}
