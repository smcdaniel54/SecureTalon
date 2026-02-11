# File read demo

**What this does:** Creates a session, sets a policy that allows `file.read` only under the `work/` directory, sends a message with a file-read intent, and shows the run. You get a working example in one command.

## Prereqs

- Backend: `ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon`
- Create a file to read (or let the script create it):

```powershell
New-Item -ItemType Directory -Force work | Out-Null; 'hello from SecureTalon' | Set-Content work/input.txt
```

## Ready-to-run policy

This is the exact policy applied by the script (`PUT /v1/sessions/{id}/policy`). Copy-paste to use in the UI or your own client.

```json
{
  "overrides": [
    {
      "tool": "file.read",
      "allow": true,
      "constraints": {
        "roots": ["work"],
        "max_bytes": 1048576
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
  "content": "Read work/input.txt",
  "intents": [
    {
      "tool": "file.read",
      "params": { "path": "work/input.txt" }
    }
  ]
}
```

## Run

From repo root (PowerShell):

```powershell
.\examples\file-read-demo\run.ps1
```

Then open the UI → Sessions → select the session → open the run to see `policy_eval` (allow) and `tool_exec` (ok) with file content in step details.
