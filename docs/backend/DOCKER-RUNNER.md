# Docker Skill Runner

## Philosophy
Skills are **untrusted**. They must be sandboxed and cannot directly access the host.

SecureTalon runs skills as Docker containers with hardened defaults.

---

## Required Docker settings (MVP)
When invoking `docker run`, use:
- `--read-only`
- `--cap-drop=ALL`
- `--security-opt no-new-privileges`
- `--pids-limit=128`
- `--memory=512m` (configurable)
- `--cpus=1.0` (configurable)
- `--network=none` (default)
- `--user=1000:1000` (non-root) when image supports it
- `--tmpfs /tmp:rw,noexec,nosuid,size=64m`
- `--workdir /work`
- No volume mounts by default.

### Optional allowances (only via capability token constraints)
- Limited network egress (specific domains) — enforced by:
  - (MVP) broker-proxied HTTP fetch only (preferred)
  - (Phase 2) egress firewall rules / sidecar proxy
- Specific read-only mounts (e.g., `/work/input`) with size limits

---

## Execution contract
Skill container communicates through STDIN/STDOUT JSON.

### Input (stdin)
```json
{
  "session_id": "sess_...",
  "skill": "hello-world",
  "args": { "name": "Stan" },
  "capability_context": {
    "allowed_tools": ["file.read"],
    "token_refs": ["cap_..."]
  }
}
```

### Output (stdout)
```json
{
  "status": "ok",
  "result": { "message": "Hello Stan" },
  "requests": [
    {
      "type": "tool_intent",
      "tool": "file.read",
      "params": { "path": "/work/input.txt" }
    }
  ]
}
```

**Important**: The skill can *request* tool actions, but cannot execute them.  
SecureTalon will evaluate these tool_intents through Policy Engine and execute via Tool Broker.

---

## Image allowlisting and signing
- Skills must be referenced by immutable digest: `image@sha256:...`
- Broker verifies signature and digest.
- Only images from allowed registries.

---

## Developer workflow
- `make skill-new NAME=hello-world`
  - scaffolds manifest
  - creates Dockerfile template
  - registers skill locally for dev

---

## Failure modes
- Timeout: terminate container and return structured error
- Non-zero exit: capture stderr, redact secrets, log hashes
- Output parse failure: treat as failure; deny any implicit tool actions

