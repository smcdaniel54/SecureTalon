# File read example

Allow `file.read` under a directory and run a message with a file-read intent.

**Prereqs:** Backend on `http://localhost:8090` with `ADMIN_TOKEN=demo`. Create a file to read, e.g.:

```powershell
New-Item -ItemType Directory -Force work | Out-Null; 'hello from SecureTalon' | Set-Content work/input.txt
```

**Run:** From repo root, `.\scripts\demo-happy-path.ps1` (or run the steps below).

**Steps:**

1. Create session: `POST /v1/sessions` with `{"label":"File read demo","metadata":{}}`.
2. Set policy: `PUT /v1/sessions/{id}/policy` with body:
   ```json
   {"overrides":[{"tool":"file.read","allow":true,"constraints":{"roots":["work"],"max_bytes":1048576}}]}
   ```
3. Post message with intent: `POST /v1/sessions/{id}/messages` with:
   ```json
   {"role":"user","content":"Read work/input.txt","intents":[{"tool":"file.read","params":{"path":"work/input.txt"}}]}
   ```
4. Get run: `GET /v1/runs/{run_id}` — check steps for `policy_eval` (allow) and `tool_exec` (ok) with file content in details.
5. Audit: `GET /v1/audit?session_id=...` and `GET /v1/audit/validate?session_id=...`.
6. Replay: `POST /v1/runs/{run_id}/replay`.

See [scripts/demo-happy-path.ps1](../../scripts/demo-happy-path.ps1) for a full PowerShell script.
