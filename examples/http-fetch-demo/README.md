# HTTP fetch demo

**What this does:** Creates a session, sets a policy that allows `http.fetch` only to `https://jsonplaceholder.typicode.com` with GET, sends a message with an HTTP GET intent, and shows the run. No file system or Docker required—ideal for trying SecureTalon quickly.

## Prereqs

- Backend: `ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon`

## Ready-to-run policy

Exact policy applied by the script (`PUT /v1/sessions/{id}/policy`). See also `policy.json` in this folder.

```json
{
  "overrides": [
    {
      "tool": "http.fetch",
      "allow": true,
      "constraints": {
        "domains": ["https://jsonplaceholder.typicode.com"],
        "methods": ["GET"],
        "max_bytes": 200000
      }
    }
  ]
}
```

## Ready-to-run intent

Message body for `POST /v1/sessions/{id}/messages`:

```json
{
  "role": "user",
  "content": "Fetch post 1 from API",
  "intents": [
    {
      "tool": "http.fetch",
      "params": {
        "url": "https://jsonplaceholder.typicode.com/posts/1",
        "method": "GET"
      }
    }
  ]
}
```

## Run

From repo root (PowerShell):

```powershell
.\examples\http-fetch-demo\run.ps1
```

Then open the UI → Sessions → run → step details to see the HTTP response.
