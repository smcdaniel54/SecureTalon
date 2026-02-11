package broker

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (b *Broker) doHTTPFetch(params map[string]interface{}, constraints map[string]interface{}) (map[string]interface{}, error) {
	urlStr, _ := params["url"].(string)
	method, _ := params["method"].(string)
	if urlStr == "" {
		return nil, fmt.Errorf("url required")
	}
	if method == "" {
		method = "GET"
	}
	// Domain allowlist
	domains := constraints["domains"]
	if domains != nil && !domainAllowed(urlStr, domains) {
		return nil, fmt.Errorf("url domain not in allowlist")
	}
	methods, _ := constraints["methods"].([]interface{})
	if len(methods) > 0 && !methodAllowed(method, methods) {
		return nil, fmt.Errorf("method %s not in allowlist", method)
	}
	maxBytes := 200000
	if m, ok := constraints["max_bytes"].(float64); ok && m > 0 {
		maxBytes = int(m)
	}
	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SecureTalon/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(maxBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxBytes {
		return nil, fmt.Errorf("response exceeds max_bytes %d", maxBytes)
	}
	return map[string]interface{}{
		"status_code": resp.StatusCode,
		"body":        string(body),
		"bytes":       len(body),
	}, nil
}

func domainAllowed(urlStr string, domains interface{}) bool {
	// Parse host from URL (simple: take between // and / or end)
	u := urlStr
	if i := strings.Index(u, "//"); i >= 0 {
		u = u[i+2:]
	}
	if i := strings.Index(u, "/"); i >= 0 {
		u = u[:i]
	}
	host := strings.ToLower(strings.Split(u, ":")[0])
	switch v := domains.(type) {
	case []string:
		for _, d := range v {
			if host == strings.ToLower(d) || strings.HasSuffix(host, "."+strings.ToLower(d)) {
				return true
			}
		}
	case []interface{}:
		for _, d := range v {
			if s, ok := d.(string); ok {
				if host == strings.ToLower(s) || strings.HasSuffix(host, "."+strings.ToLower(s)) {
					return true
				}
			}
		}
	}
	return false
}

func methodAllowed(method string, allowlist []interface{}) bool {
	for _, m := range allowlist {
		if s, ok := m.(string); ok && strings.EqualFold(s, method) {
			return true
		}
	}
	return false
}

