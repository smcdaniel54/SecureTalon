# Cursor AI – Combined Master Prompt (Backend + UI)

This prompt assumes you have a repo open in Cursor where you want to build SecureTalon (Go backend) plus the Admin Console UI (Svelte + Vite).
It implements **backend first**, then UI, using the design packs.

---

## Ground rules
- Security-first, deny-by-default. The LLM is NOT a security boundary.
- Implement hard constraints: capability tokens + broker mediation + docker sandboxing + audit hash chain.
- Keep MVP tight: ship a working secure core quickly.

---

## PART A — Backend (Go) MVP

### 1) Create repo structure
Create a Go module `securetalon` with:
```
/cmd/securetalon
/internal/{api,core,policy,broker,audit,replay,skills,config,auth}
/docs
/data (gitignored)
```

### 2) Implement config + auth
- Config sources: env + config file (YAML or JSON)
- Required:
  - `ADMIN_TOKEN`
  - data dir for audit store
- Auth: `Authorization: Bearer <ADMIN_TOKEN>` for all `/v1/*` endpoints.

### 3) Implement core types
Add `internal/core/types.go` and related files with:
- Session, Message
- Run, Step
- ToolIntent, Decision
- CapabilityToken (signed blob)
- AuditEvent (hash-chained)

### 4) Implement REST API (per backend spec)
Implement endpoints:
- `POST /v1/sessions`
- `GET /v1/sessions`
- `GET /v1/sessions/{id}`
- `POST /v1/sessions/{id}/messages` (starts a run; returns `run_id`)
- `GET /v1/sessions/{id}/messages`
- `GET /v1/runs/{run_id}`
- `GET /v1/policy/effective?session_id=...`
- `PUT /v1/sessions/{id}/policy`
- `GET /v1/skills`
- `POST /v1/skills`
- `GET /v1/audit?...`

Error format:
```json
{ "error": { "code": "POLICY_DENIED", "message": "...", "details": { } } }
```

### 5) Policy Engine (deny-by-default)
Implement allowlist rules + per-session overrides:
- file.read/write:
  - allow roots only
  - max bytes enforced
- http.fetch:
  - domains allowlist
  - methods allowlist
  - max bytes enforced
- docker.run:
  - only images by digest `@sha256:...`
  - must be registered + signed + allowlisted registry (simple allowlist in config for MVP)

Decision outputs:
- ALLOW -> issue capability token (TTL default 60s)
- DENY -> reason + safe suggestion

### 6) Capability tokens
Implement signed tokens:
- Prefer HMAC-SHA256 using server secret (MVP), or Ed25519 if you want keypairs.
- Token must bind:
  - session_id
  - subject
  - tool
  - constraints
  - exp
Broker must verify all fields, expiry, and constraints.

### 7) Tool Broker (the only executor)
Implement:
- Verify token and constraints
- Tools:
  - file.read
  - file.write
  - http.fetch
  - docker.run (skills)

Shell execution: OFF by default; stub only.

### 8) Docker Skill Runner (hardened defaults)
For `docker.run`, invoke docker with:
- `--read-only`
- `--cap-drop=ALL`
- `--security-opt no-new-privileges`
- `--pids-limit=128`
- `--memory`, `--cpus` limits
- `--network=none` default
- `--tmpfs /tmp:rw,noexec,nosuid,size=64m`
- no mounts by default

Skills communicate via STDIN/STDOUT JSON. Any skill-requested tool intents must be re-evaluated by Policy Engine and executed via Broker.

### 9) Audit Store (append-only + hash chain)
- JSONL append-only storage under `./data/audit/`
- Hash chain:
  - `hash = sha256(prev_hash + canonical_json(event))`
- Implement query filtering by session/time/type.

### 10) Replay (safe replay only)
- Validate hash chain
- Reconstruct run timeline
- Return recorded outputs (no tool execution)

### 11) Tests + Make targets
- `make run` -> starts server `:8080`
- `make test` -> runs unit tests
Unit tests:
- deny-by-default
- token issue/verify
- constraint enforcement (block `/etc/passwd`)
- audit hash chain validation

### 12) MVP acceptance checklist
Backend passes the acceptance tests described in the backend docs pack.

---

## PART B — UI (Svelte + Vite) Admin Console

### 1) Create UI project
In `/ui`:
- `npm create vite@latest` (Svelte + TS)
- Add `svelte-spa-router`
- Optional: Tailwind (recommended)

### 2) UI structure
```
/ui/src
  /app (AppShell, routes, stores)
  /pages (Login, Dashboard, Sessions, SessionDetail, RunDetail, Policies, Skills, Audit, Replay)
  /components (templates: Card, DataTable, Timeline, FormPanel, Modal, Toast)
  /lib (api.ts, types.ts, format.ts, securityHints.ts)
```

### 3) Auth UX
- Login page collects:
  - API base (default http://localhost:8080)
  - admin token
- Store in `localStorage`
- “Connect” tests `GET /v1/sessions?limit=1`
- “Forget” clears localStorage

### 4) Pages to implement (MVP)
- Dashboard (recent sessions/runs summary)
- Sessions list + create
- Session detail: messages + runs; posting a message starts a run
- Run detail: timeline of steps, policy decisions, tool results
- Policies: effective policy + session overrides editor with guardrails
- Skills: list + register
- Audit: filter + timeline/table, validate hash chain in browser
- Replay: render safe replay timeline from audit events (no re-exec)

### 5) Guardrails (security-first UI)
- Prevent saving empty allowlists
- Warnings for risky actions:
  - enabling network for skills
  - enabling shell.exec
  - allowing host mounts
- “Fix safely” suggestions when a run is denied

### 6) UI deliverables
- `npm run dev` hot reload works
- `/ui/README.md` explains setup and API mapping

---

## Implementation order (strict)
1) Backend scaffolding + auth + sessions/messages
2) Policy + tokens + broker
3) Docker runner
4) Audit + query
5) Replay
6) UI connect + sessions/messages
7) Run detail + policy editor + audit explorer + replay viewer

---

## Definition of Done
- From UI:
  - Create session
  - Send message -> run starts
  - View run timeline and policy decisions
  - Register skill and run it sandboxed
  - Explore audit events and see hash-chain status
  - Replay visualization works without re-executing tools

---

## Notes for Cursor
- Prefer small, well-named packages and interfaces.
- Use explicit JSON structs, avoid reflection-heavy magic.
- Redact secrets everywhere.
- Add comments describing security invariants in code.

