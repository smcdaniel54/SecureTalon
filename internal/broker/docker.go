package broker

import (
	"fmt"
	"os/exec"
	"strings"
)

// doDockerRun runs a skill image by digest with hardened defaults.
// Constraints may include: image (digest), memory, cpus, network (allow).
func (b *Broker) doDockerRun(params map[string]interface{}, constraints map[string]interface{}) (map[string]interface{}, error) {
	image, _ := params["image"].(string)
	if image == "" {
		return nil, fmt.Errorf("image required (use image@sha256:...)")
	}
	if !strings.Contains(image, "@sha256:") {
		return nil, fmt.Errorf("docker.run only allowed with digest: image@sha256:...")
	}
	// Optional: check image against allowlist in constraints
	allowedImages, _ := constraints["images"].([]interface{})
	if len(allowedImages) > 0 {
		ok := false
		for _, a := range allowedImages {
			if s, _ := a.(string); s == image {
				ok = true
				break
			}
		}
		if !ok {
			return nil, fmt.Errorf("image not in allowlist")
		}
	}
	// Hardened defaults per docs/backend/DOCKER-RUNNER.md
	args := []string{
		"run", "--rm",
		"--read-only",
		"--cap-drop=ALL",
		"--security-opt", "no-new-privileges",
		"--pids-limit=128",
		"--memory=512m",
		"--cpus=1.0",
		"--network=none",
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=64m",
		"--workdir", "/work",
		image,
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	if err != nil {
		return map[string]interface{}{
			"error":  err.Error(),
			"exit":   exitCode,
			"stderr": string(out),
		}, nil
	}
	return map[string]interface{}{
		"stdout": string(out),
		"exit":   exitCode,
	}, nil
}
