# HTTP fetch example

Allow `http.fetch` to a domain (e.g. a public API or your own service) and send a message that triggers a GET.

**Prereqs:** Backend on `http://localhost:8090` with `ADMIN_TOKEN=demo`.

**Option A — Public API (e.g. JSON placeholder):**

1. Create session: `POST /v1/sessions` → `{"label":"HTTP demo","metadata":{}}`.
2. Set policy: `PUT /v1/sessions/{id}/policy` with body:
   ```json
   {
     "overrides": [{
       "tool": "http.fetch",
       "allow": true,
       "constraints": {
         "domains": ["https://jsonplaceholder.typicode.com"],
         "methods": ["GET"],
         "max_bytes": 200000
       }
     }]
   }
   ```
3. Post message with intent: `POST /v1/sessions/{id}/messages` with:
   ```json
   {
     "role": "user",
     "content": "Fetch posts from API",
     "intents": [{
       "tool": "http.fetch",
       "params": {
         "url": "https://jsonplaceholder.typicode.com/posts/1",
         "method": "GET"
       }
     }]
   }
   ```
4. Get run: `GET /v1/runs/{run_id}` — steps should show `policy_eval` (allow) and `tool_exec` (ok) with response snippet in details.

**Option B — Run the script:**

From repo root (PowerShell):

```powershell
.\examples\http-fetch\run.ps1
```

This creates a session, sets the policy above, posts the message, and prints the run ID and step summary.
