# OpenClaw migration guide

If you’re coming from OpenClaw, this guide maps concepts to SecureTalon and shows how to convert configs and flows.

## Concept mapping

| OpenClaw | SecureTalon |
|----------|-------------|
| **exec.security** / policy knobs | **Session policy overrides** — per-session allowlists in the Policy Engine (UI or `PUT /v1/sessions/{id}/policy`) |
| **exec.ask** / approvals | **Deny-by-default** — intents are denied until you add an override; future: explicit approval workflows |
| **Per-agent allowlist** | **Per-session overrides** — each session has a list of `{ tool, allow, constraints }` (e.g. `file.read` with `roots`, `http.fetch` with `domains`) |
| **Tool profiles / tool groups** | **Overrides per tool** — you allow `file.read`, `http.fetch`, `docker.run` with tool-specific constraints (paths, domains, images) |
| **Tool execution** | **Broker-mediated execution** — agent sends intents; policy engine issues a capability token only if an override allows; broker verifies token and runs the tool |
| **Logs / observability** | **Audit** — `GET /v1/audit` with hash chain; **Replay** — `POST /v1/runs/{id}/replay` (no re-execution) |
| **Elevated mode** | **Explicit overrides** — e.g. allow `docker.run` with `network_allowed` or `mounts_allowed` only when you add that to the session policy (with guardrails in the UI) |

## What needs to change

1. **Policy** — Move from OpenClaw’s exec knobs and per-agent allowlists to **session-level overrides** in SecureTalon. Each tool you want to allow needs an override with **constraints** (e.g. domains for HTTP, roots for file, images for Docker).
2. **Tool calls** — Your agent (or client) sends **intents** (tool + params) in the message body; SecureTalon’s agent runs the loop and the **broker** executes only after a capability token is issued.
3. **Approvals** — Today SecureTalon uses “allow in policy = run”; there is no interactive ask/approval step yet. To mimic “ask”, use a session that has no override for that tool (deny), then add an override after review (see [Deny then fix safely](../examples/deny-then-fix/README.md)).
4. **Observability** — Use the **Audit** and **Replay** UI (or `/v1/audit`, `/v1/runs/{id}/replay`) instead of OpenClaw’s logs; validate the hash chain for tamper evidence.

## Example conversions

### 1. “Allow file read under /workspace”

**OpenClaw (conceptual):** allowlist tool for file read, scope to a path.

**SecureTalon:** Session policy override:

```json
{
  "overrides": [
    {
      "tool": "file.read",
      "allow": true,
      "constraints": {
        "roots": ["workspace"],
        "max_bytes": 1048576
      }
    }
  ]
}
```

`PUT /v1/sessions/{session_id}/policy` with the body above. Only paths under `workspace` (or `workspace/...`) are allowed.

### 2. “Allow HTTP GET to api.example.com”

**SecureTalon:** Session policy override:

```json
{
  "overrides": [
    {
      "tool": "http.fetch",
      "allow": true,
      "constraints": {
        "domains": ["https://api.example.com"],
        "methods": ["GET"],
        "max_bytes": 200000
      }
    }
  ]
}
```

### 3. “Allow running a specific Docker image (by digest)”

**SecureTalon:** Register the skill (optional for broker allowlist), then session policy:

```json
{
  "overrides": [
    {
      "tool": "docker.run",
      "allow": true,
      "constraints": {
        "images": ["myorg/worker@sha256:abc123..."],
        "network_allowed": false,
        "mounts_allowed": false
      }
    }
  ]
}
```

### 4. “Deny by default, then allow after review”

- Do **not** add an override for a tool → intent is **denied**; the run records a `policy_eval` step with status `denied` and a reason.
- In the UI, open **Policies** for that session (or use the “Fix safely” flow from a denied run link); add an override with the suggested constraint; Save.
- Retry the message/intent; it is now allowed if it matches the new override.

See the [examples](../examples/) directory for runnable HTTP, Docker, and deny-then-fix demos.
