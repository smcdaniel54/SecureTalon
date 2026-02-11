# SecureTalon Admin Console (UI)

Svelte + Vite SPA for the SecureTalon backend.

## Run

```bash
cd ui
npm install
npm run dev
```

Open http://localhost:5173 (or the URL Vite prints).

**Connect** in the UI: set API base to `http://localhost:8090` (or your backend URL) and the admin token. Credentials are stored in `localStorage` only (MVP). Use **Disconnect** to clear.

## Build

```bash
npm run build
```

Output is in `dist/`.

## API base

Use **http://localhost:8090** when the backend is run with `ADDR=:8090` (see repo root README). Policy and Skills UIs use guardrails and confirmations for risky or broad actions.

## API mapping

| UI action           | Backend API |
|---------------------|-------------|
| Connect             | `GET /v1/sessions?limit=1` (test auth) |
| Create session      | `POST /v1/sessions` |
| List sessions       | `GET /v1/sessions` |
| Session detail      | `GET /v1/sessions/{id}`, `GET /v1/sessions/{id}/messages` |
| Post message        | `POST /v1/sessions/{id}/messages` → returns `run_id` |
| Run detail          | `GET /v1/runs/{run_id}` |
| Policies            | `GET /v1/policy/effective?session_id=...`, `PUT /v1/sessions/{id}/policy` |
| Skills              | `GET /v1/skills`, `POST /v1/skills` |
| Replay (safe)       | `POST /v1/runs/{run_id}/replay` → timeline only, no re-exec |
| Audit               | `GET /v1/audit?session_id=...&run_id=...&type=...&since=...&until=...&limit=500` |
| Validate chain      | `GET /v1/audit/validate?session_id=...&limit=500` |
