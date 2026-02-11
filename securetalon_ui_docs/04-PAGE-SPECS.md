# Page Specs (UI MVP)

## 1) Login / Connect
- Inputs:
  - API Base URL
  - Admin Token
- Buttons:
  - Connect (test `GET /v1/sessions?limit=1`)
  - Forget credentials

## 2) Dashboard
- Cards:
  - Recent Sessions
  - Running/Failed Runs
  - Policy Summary (deny-by-default active)
  - Skills Summary (signed count)
- Quick actions:
  - Create Session
  - Register Skill

## 3) Sessions
- Table: id, label, status, created_at
- Actions: Open, Copy ID
- Create modal

## 4) Session Detail
- Tabs:
  - Messages
  - Runs
  - Policy (shortcut)
  - Audit (filtered)
- Messages:
  - chat thread + input box
  - after send: show run_id link
- Runs list:
  - status badge + open

## 5) Run Detail
- Header: status, started/ended, session link
- Timeline:
  - steps with icons + expandable JSON
- Policy decisions:
  - show allow/deny reasons
  - “Fix safely” suggestions
- Tool results:
  - show summaries (bytes read, domain fetched, docker image digest)

## 6) Policies
- Show effective rules + overrides
- Editor:
  - add override for tool
  - constraints UI:
    - file roots picker (text)
    - domain allowlist chips
    - max bytes number
    - TTL seconds
- Danger zone:
  - shell.exec toggle (disabled by default), heavy warning

## 7) Skills
- List skills
- Register form
- Show manifest and requested tools
- Signature badge

## 8) Audit Explorer
- Filters
- Table view + timeline view toggle
- Hash chain validator:
  - show “valid chain” or “broken at event X”

## 9) Replay Viewer
- Choose run
- Render timeline using audit events
- Step through (prev/next)
- No re-execution; purely visualization (MVP)

