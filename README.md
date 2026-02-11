# SecureTalon

Policy-driven agent runtime: runs tool intents (file, HTTP, Docker) only when policy allows, with full audit and safe replay.

## Backend

Requires `ADMIN_TOKEN`. Optional: `ADDR` (default `:8080`), `DATA_DIR`, `TOKEN_SECRET`.

```bash
ADMIN_TOKEN=your-secret-token go run ./cmd/securetalon
```

To use port 8090 (e.g. for the UI default):

```bash
ADDR=:8090 ADMIN_TOKEN=your-secret-token go run ./cmd/securetalon
```

## UI

```bash
cd ui
npm install
npm run dev
```

Open http://localhost:5173. In the UI, **Connect** with:

- **API base:** `http://localhost:8090` (or your backend URL)
- **Token:** same as `ADMIN_TOKEN`

See [ui/README.md](ui/README.md) for API mapping and build.

## Happy path demo

1. Start backend: `ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon`
2. Start UI: `cd ui && npm i && npm run dev` → connect with `http://localhost:8090` and token `demo`
3. **Create session** (Sessions → New session).
4. **Set policy:** Policies → select the session → add override: tool `file.read`, allow, roots e.g. `work` (chips), Save.
5. Create a file the agent can read (e.g. `mkdir work 2>nul & echo hello > work\input.txt` on Windows, or `mkdir -p work && echo hello > work/input.txt` on Unix).
6. **Send message** with parseable intent: in Session detail, send e.g. “Read work/input.txt” or use a message that triggers a file.read intent; or call the API with body `{"content":"...", "intents":[{"tool":"file.read","params":{"path":"work/input.txt"}}]}`.
7. **View run:** open the run from the session; confirm steps (policy_eval ALLOW/DENY, tool_exec ok/error) and status (queued → running → completed).
8. **Audit:** Audit → set session filter → Refresh → Validate chain → “Chain OK”.
9. **Replay:** Replay → enter run ID → Load Safe Replay → step through timeline (no re-execution).

See [scripts/demo-happy-path.ps1](scripts/demo-happy-path.ps1) for an automated API-only demo (PowerShell).
