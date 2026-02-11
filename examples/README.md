# Examples

Runnable examples for SecureTalon. Prerequisites: backend running (`ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon`).

## Ready-to-run demos (recommended)

Each demo includes **sample policy** and **sample intent** (in README and in `policy.json` / `message.json`) so you can copy-paste or run the script and “get it” immediately.

| Demo | What it does |
|------|----------------|
| [file-read-demo](file-read-demo/) | Allow `file.read` under `work/`; send intent; view run. Includes `policy.json` and `run.ps1` (creates `work/input.txt` if missing). |
| [http-fetch-demo](http-fetch-demo/) | Allow `http.fetch` to jsonplaceholder.typicode.com (GET); send intent; see response in run. Includes `policy.json`, `message.json`, `run.ps1`. |
| [docker-skill-demo](docker-skill-demo/) | Allow `docker.run` for an image by digest (no network/mounts); send intent; run container. Set `IMAGE_DIGEST` or edit `policy.json`/`message.json`. |

**Run from repo root (PowerShell):**  
`.\examples\file-read-demo\run.ps1` · `.\examples\http-fetch-demo\run.ps1` · `.\examples\docker-skill-demo\run.ps1`

## Other examples

| Example | Description |
|---------|-------------|
| [file-read](file-read/) | Same flow as file-read-demo; points to [scripts/demo-happy-path.ps1](../scripts/demo-happy-path.ps1). |
| [http-fetch](http-fetch/) | Inline script version of http-fetch-demo. |
| [docker-skill](docker-skill/) | Inline script version of docker-skill-demo. |
| [deny-then-fix](deny-then-fix/) | Trigger a deny, add policy override, retry. |

All examples use the REST API; you can replicate the same flows in the UI (Sessions, Policies, Audit, Replay).
