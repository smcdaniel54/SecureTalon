# Cursor AI – Master Implementation Prompt (Backend)

Copy/paste this whole prompt into Cursor AI (with the repo opened) to implement SecureTalon backend MVP.

---

## Objective
Implement **SecureTalon**: a **security-first agent platform** in Golang with these non-negotiables:
- Deny-by-default policy
- Capability tokens (signed, short TTL) required for every tool execution
- Tool Broker is the only path to OS/files/network/docker
- Skills run in Docker with hardened defaults
- Append-only audit log with hash chaining + safe replay

## Repo Setup
Create a new Go module:
- module name: `securetalon`
- Go version: latest stable
- Create folders:
  - `/cmd/securetalon`
  - `/internal/{api,core,policy,broker,audit,replay,skills,config,auth}`

## MVP Scope
### APIs
Implement endpoints from [API-SPEC.md](API-SPEC.md):
- sessions: create/list/get
- messages: post/list
- runs: get status
- policy: get effective, set session overrides
- skills: list/register
- audit: query

Auth: Bearer token `ADMIN_TOKEN` read from config/env.

### Core Types
Implement core structs and JSON schemas:
- Session, Message, Run, Step
- ToolIntent, Decision, CapabilityToken (token as JWT-like signed blob or custom HMAC)
- AuditEvent (with prev_hash/hash)

### Policy Engine
- Deny-by-default
- Support allowlist rules:
  - allow docker.run for signed, allowed images
  - allow file.read/write under allowed roots only
  - allow http.fetch to allowed domains/methods only
- On ALLOW: issue capability token with constraints and TTL (default 60s)

### Tool Broker
- Verify token signature + expiry + session binding
- Enforce constraints at execution time (broker is ultimate gate)
- Implement tools:
  - file.read (max bytes, allowed roots)
  - file.write (allowed roots, max bytes)
  - http.fetch (allowed domains, methods, max bytes) using Go net/http
  - docker.run (call `docker` CLI or Docker API; MVP can use CLI)
- For docker.run enforce flags listed in [DOCKER-RUNNER.md](DOCKER-RUNNER.md)

### Audit Store
- Append-only JSONL in `./data/audit/`
- Include hash chain per [AUDIT-AND-REPLAY.md](AUDIT-AND-REPLAY.md)
- Query by session/time/type

### Replay Engine
- Safe replay only (MVP):
  - Reconstruct run from audit
  - Return recorded outputs without executing tools

## Implementation Requirements
- Favor clarity and security over cleverness.
- Never log secrets; implement redaction.
- Validate all inputs.
- Return errors in the standard format in API spec.
- Add unit tests for:
  - token issuance/verification
  - deny-by-default policy
  - constraint enforcement (e.g., cannot read /etc/passwd)
  - audit hash chain validation

## Deliverables
- Working `securetalon` binary
- `make run` starts server on `:8080`
- `make test` passes
- Minimal sample skill registration in docs or example script

## Acceptance
System passes all checks in [ACCEPTANCE-TESTS.md](ACCEPTANCE-TESTS.md).

---

