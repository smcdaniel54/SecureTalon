# Audit Log and Replay

## Goals
- **Forensics**: prove what happened
- **Debuggability**: reproduce failures
- **Governance**: show safe-by-design decisions
- **Replay**: re-run workflows without re-executing risky actions

---

## Event model
Every meaningful step emits an audit event:
- `session.created`
- `message.appended`
- `run.started`
- `policy.intent.received`
- `policy.decision`
- `capability.issued` (token hash only)
- `tool.executed`
- `skill.started`
- `skill.finished`
- `run.finished`

Each event includes:
- `event_id`
- `timestamp`
- `session_id`
- `run_id` (optional)
- `type`
- `data` (structured JSON)
- `prev_hash`
- `hash`

---

## Tamper evidence
Hash chain:
- `hash = sha256(prev_hash + canonical_json(event))`
- Store `prev_hash` in each event

If any event is changed, all subsequent hashes fail validation.

---

## Storage
MVP options:
1) **JSONL append-only** file per day or per session
2) **SQLite** table with append-only enforcement and WAL

Recommendation for MVP:
- JSONL for simplicity, plus periodic compaction/export.

---

## Replay
Two replay modes:

### A) Safe replay (default)
- Does **not** execute real tools
- Returns recorded tool outputs from audit stream
- Useful for UI playback and debugging

### B) Re-execution replay (phase 2)
- Re-runs policy evaluation and tool execution
- Requires approvals for any non-read-only actions

Replay steps:
1. Load event stream for run
2. Validate hash chain
3. Reconstruct state machine
4. For each tool event:
   - safe replay returns recorded output
   - re-exec calls policy/broker again

---

## UI integration later
The UI can show:
- timeline of events
- policy allow/deny badges
- diffs between planned vs executed
- “why denied” explanations

