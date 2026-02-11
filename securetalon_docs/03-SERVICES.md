# Services and Responsibilities

## 1) API Gateway
Responsibilities:
- HTTP server, JSON API
- Authentication (MVP: admin API key)
- Validation & request shaping
- Rate limiting (basic)
- Error handling (no secret leaks)

## 2) Session Router
Responsibilities:
- Maintain sessions, message threads
- Start agent runs
- Append messages/events
- Emit audit events

Data:
- `Session` (id, created_at, labels, status)
- `Message` (id, role, content, timestamp, metadata)

## 3) Agent Runtime (Planner/Executor)
Responsibilities:
- Convert messages into plan steps
- Create tool/skill intents
- Ask Policy Engine for decisions
- Submit to Tool Broker
- Collect results and continue

MVP assumption:
- Agent Runtime can be deterministic and rule-based (no LLM required yet),
  allowing secure plumbing to be built first.
- Later, plug in LLM planning.

## 4) Policy Engine
Responsibilities:
- Deny-by-default evaluation
- Capability token issuance
- Rule evaluation (static rules + per-session context)
- (Phase 2) approvals, RBAC, risk scoring

Inputs:
- `ToolIntent`
- `SessionContext`

Outputs:
- `Decision`: ALLOW/DENY/REQUIRE_APPROVAL
- `CapabilityToken` if allowed
- `Reason` and `SuggestedFix`

## 5) Tool Broker
Responsibilities:
- Verify tokens
- Execute approved tool operations:
  - file ops (safe read/write under constraints)
  - http fetch (domain allowlist)
  - docker run (skills)
  - shell exec (very restricted; off by default)
- Return structured results (no raw secrets)
- Emit audit events

## 6) Skill Runner (Docker)
Responsibilities:
- Run skill images
- Enforce container hardening options
- Capture stdout/stderr
- Enforce timeouts/resource limits
- Return result object

## 7) Audit Store
Responsibilities:
- Append-only event log
- Query by session/tool/time
- Hash chain for tamper evidence
- Export for replay

Storage options:
- MVP: local disk (SQLite + WAL or append-only JSONL)
- Later: pluggable backends

## 8) Replay Engine
Responsibilities:
- Re-run prior session deterministically by consuming audit stream
- Replace real tool runs with recorded outputs (safe replay)
- Optionally “replay with changes” (phase 2)

