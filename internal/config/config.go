// Package config provides configuration loading from env and file.
// Required: ADMIN_TOKEN, data dir for audit store.
package config

import (
	"os"
	"path/filepath"
)

// Config holds server and security settings.
// ADMIN_TOKEN is required for all /v1/* endpoints.
type Config struct {
	// AdminToken is the bearer token for API auth (env: ADMIN_TOKEN).
	AdminToken string `yaml:"admin_token" json:"admin_token"`
	// DataDir is the root for audit store and other persistent data (env: DATA_DIR).
	DataDir string `yaml:"data_dir" json:"data_dir"`
	// Addr is the HTTP listen address (env: ADDR).
	Addr string `yaml:"addr" json:"addr"`
	// TokenSecret is used to sign capability tokens (env: TOKEN_SECRET).
	TokenSecret string `yaml:"token_secret" json:"token_secret"`
	// DockerMemoryLimit for skill containers (e.g. "512m").
	DockerMemoryLimit string `yaml:"docker_memory_limit" json:"docker_memory_limit"`
	// DockerCPULimit for skill containers (e.g. "1.0").
	DockerCPULimit string `yaml:"docker_cpu_limit" json:"docker_cpu_limit"`
	// AllowedRegistries for docker.run (comma-separated or in file).
	AllowedRegistries []string `yaml:"allowed_registries" json:"allowed_registries"`
}

// DefaultConfig returns defaults; env overrides.
func DefaultConfig() *Config {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return &Config{
		AdminToken:        os.Getenv("ADMIN_TOKEN"),
		DataDir:           dataDir,
		Addr:              addr,
		TokenSecret:       os.Getenv("TOKEN_SECRET"),
		DockerMemoryLimit: getEnv("DOCKER_MEMORY_LIMIT", "512m"),
		DockerCPULimit:    getEnv("DOCKER_CPU_LIMIT", "1.0"),
		AllowedRegistries: []string{},
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// AuditDir returns the audit log directory under DataDir.
func (c *Config) AuditDir() string {
	return filepath.Join(c.DataDir, "audit")
}

// EnsureDataDirs creates data and audit dirs if missing.
func (c *Config) EnsureDataDirs() error {
	if err := os.MkdirAll(c.DataDir, 0700); err != nil {
		return err
	}
	return os.MkdirAll(c.AuditDir(), 0700)
}
