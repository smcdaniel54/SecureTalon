package broker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const defaultMaxBytes = 1024 * 1024 // 1MB

func (b *Broker) doFileRead(params map[string]interface{}, constraints map[string]interface{}) (map[string]interface{}, error) {
	path, _ := params["path"].(string)
	if path == "" {
		return nil, fmt.Errorf("path required")
	}
	maxBytes := defaultMaxBytes
	if m, ok := constraints["max_bytes"].(float64); ok && m > 0 {
		maxBytes = int(m)
	}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(io.LimitReader(f, int64(maxBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(content) > maxBytes {
		return nil, fmt.Errorf("file exceeds max_bytes %d", maxBytes)
	}
	return map[string]interface{}{
		"path":    path,
		"content": string(content),
		"bytes":   len(content),
	}, nil
}

func (b *Broker) doFileWrite(params map[string]interface{}, constraints map[string]interface{}) (map[string]interface{}, error) {
	path, _ := params["path"].(string)
	content, _ := params["content"].(string)
	if path == "" {
		return nil, fmt.Errorf("path required")
	}
	maxBytes := defaultMaxBytes
	if m, ok := constraints["max_bytes"].(float64); ok && m > 0 {
		maxBytes = int(m)
	}
	if len(content) > maxBytes {
		return nil, fmt.Errorf("content exceeds max_bytes %d", maxBytes)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil && !os.IsExist(err) {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"path":  path,
		"bytes": len(content),
	}, nil
}
