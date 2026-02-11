# API Mapping (UI ↔ Backend)

This UI consumes the backend REST APIs (see backend pack `docs/04-API-SPEC.md`).

## Auth
- Store:
  - `apiBase` (default `http://localhost:8080`)
  - `token` (admin token)
- Persist to `localStorage` only (MVP). Provide “forget” button.

## Sessions
- Create: `POST /v1/sessions`
- List: `GET /v1/sessions`
- Detail: `GET /v1/sessions/{id}`

UI features:
- Search/filter by label
- Create session modal

## Messages
- Post message: `POST /v1/sessions/{id}/messages` -> returns `run_id`
- List messages: `GET /v1/sessions/{id}/messages`

UI features:
- Chat-like thread
- “Send” starts run; show spinner tied to run status

## Runs
- Get run: `GET /v1/runs/{run_id}`

UI features:
- Step timeline with icons:
  - policy_eval, tool_exec, skill_run
- Inline view of deny reasons

## Policies
- Get effective: `GET /v1/policy/effective?session_id=...`
- Update overrides: `PUT /v1/sessions/{id}/policy`

UI features:
- Guardrails:
  - show diff: effective vs overrides
  - forms enforce least privilege (domains list, file roots, max bytes, TTL)
- “Simulate” button (phase 2; stub now)

## Skills
- List: `GET /v1/skills`
- Register: `POST /v1/skills`

UI features:
- Status badges: signed/unsigned, allowlisted registry
- Register form:
  - name, version, image@digest, signature, key id, manifest

## Audit
- Query: `GET /v1/audit?...`

UI features:
- Filters:
  - type, time range, session_id, run_id
- Hash chain status:
  - green if valid; red if broken (validate client-side)

## Replay
- Endpoint may be:
  - `POST /v1/runs/{run_id}/replay` (if you add it)
  - Or UI uses audit to render replay timeline (MVP)
MVP approach:
- Use audit events to render replay timeline locally (safe replay visualization)

