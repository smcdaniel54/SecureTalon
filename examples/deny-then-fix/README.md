# Deny then fix safely

Demonstrate deny-by-default: send an intent that is **not** allowed by policy (e.g. `file.read` with no override, or `http.fetch` to a domain not in the allowlist). Observe the **denied** run; then add a policy override and retry (or use the UI “Fix safely” flow).

**Prereqs:** Backend on `http://localhost:8090` with `ADMIN_TOKEN=demo`. Create a file to read later, e.g. `work/input.txt`.

**Flow:**

1. **Create session** — `POST /v1/sessions` → `{"label":"Deny then fix","metadata":{}}`.
2. **Do not** add any `file.read` override (or add one that does **not** include the path you will request).
3. **Post message with intent** — e.g. `file.read` for `work/input.txt`. The run will record a `policy_eval` step with status **denied** and a reason (e.g. “no matching allowlist”).
4. **Inspect run** — `GET /v1/runs/{run_id}` or open in UI. Note the deny reason and suggested fix.
5. **Fix policy** — `PUT /v1/sessions/{id}/policy` with an override that allows `file.read` and `roots: ["work"]`.
6. **Retry** — Post the same message again (or a new message with the same intent). The new run should show `policy_eval` (allow) and `tool_exec` (ok).

**In the UI:**

- Open **Policies** → select the session. You’ll see “Effective policy” (read-only) and “Session overrides” (editable).
- Add an override for `file.read` with roots `work`, Save.
- In **Session detail**, send again: “Read work/input.txt” (with intent) or use the same API body.
- **Fix safely (deep link):** If you open Policies with `?session_id=...&run_id=...` (from a denied run), the UI can suggest an override based on the deny reason; add it and Save, then retry.

**Run the script:**

From repo root (PowerShell):

```powershell
.\examples\deny-then-fix\run.ps1
```

This creates a session, sends a file-read intent with **no** file override (deny), prints the run and deny reason, then adds the override and retries (allow + tool_exec).
