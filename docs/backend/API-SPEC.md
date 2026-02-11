# REST API Spec (MVP)

Base: `/v1`

## Auth (MVP)
- Header: `Authorization: Bearer <ADMIN_TOKEN>`
- Admin token stored in config file or env var.
- Reject missing/invalid tokens with `401`.

---

## Sessions

### Create session
`POST /v1/sessions`
```json
{ "label": "demo", "metadata": { "owner": "stan" } }
```
Response `201`:
```json
{ "id": "sess_...", "created_at": "...", "label": "demo", "status": "active" }
```

### List sessions
`GET /v1/sessions?limit=50&cursor=...`

### Get session
`GET /v1/sessions/{session_id}`

---

## Messages

### Post message
`POST /v1/sessions/{session_id}/messages`
```json
{
  "role": "user",
  "content": "Run skill X with input Y",
  "metadata": { "source": "api" }
}
```

Response `202` (run started):
```json
{ "run_id": "run_...", "status": "queued" }
```

### List messages
`GET /v1/sessions/{session_id}/messages?limit=200`

---

## Runs

### Get run status
`GET /v1/runs/{run_id}`
Response:
```json
{
  "id": "run_...",
  "session_id": "sess_...",
  "status": "running",
  "started_at": "...",
  "ended_at": null,
  "steps": [
    { "step_id": "s1", "type": "policy_eval", "status": "ok" },
    { "step_id": "s2", "type": "tool_exec", "tool": "docker.run", "status": "ok" }
  ]
}
```

---

## Policies (MVP: static + per-session overrides)

### Get effective policy
`GET /v1/policy/effective?session_id=...`

### Update session policy overrides
`PUT /v1/sessions/{session_id}/policy`
```json
{
  "overrides": [
    {
      "tool": "http.fetch",
      "allow": true,
      "constraints": {
        "domains": ["api.example.com"],
        "methods": ["GET"],
        "max_bytes": 200000
      }
    }
  ]
}
```

---

## Skills

### List skills
`GET /v1/skills`
Response:
```json
{
  "skills": [
    { "name": "hello-world", "version": "1.0.0", "image": "registry/hello-world@sha256:...", "signed": true }
  ]
}
```

### Register skill (admin)
`POST /v1/skills`
```json
{
  "name": "hello-world",
  "version": "1.0.0",
  "image": "registry/hello-world@sha256:...",
  "signature": "base64...",
  "public_key_id": "key_001",
  "manifest": { "entrypoint": "/skill/run", "requested_tools": ["file.read"] }
}
```

---

## Audit

### Query audit events
`GET /v1/audit?session_id=...&since=...&until=...&type=...&limit=500`

Response:
```json
{
  "events": [
    {
      "event_id": "evt_...",
      "ts": "...",
      "session_id": "sess_...",
      "type": "policy.decision",
      "data": { "decision": "ALLOW", "tool": "docker.run", "reason": "matched allowlist" },
      "hash": "sha256:...",
      "prev_hash": "sha256:..."
    }
  ],
  "next_cursor": null
}
```

---

## Error format (standard)
All non-2xx responses:
```json
{
  "error": {
    "code": "POLICY_DENIED",
    "message": "Tool intent denied by policy",
    "details": { "tool": "shell.exec", "reason": "shell disabled by default" }
  }
}
```

