# Docker skill demo

**What this does:** Registers a skill (Docker image by digest), sets a policy that allows `docker.run` only for that image with network and mounts disabled, sends a message with a docker.run intent, and shows the run. Requires Docker available to the backend.

## Prereqs

- Backend with Docker: `ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon`
- An image by digest. Example for `hello-world`:
  ```powershell
  docker pull hello-world
  docker inspect hello-world --format '{{index .RepoDigests 0}}'
  ```
  Use that value (e.g. `hello-world@sha256:...`) in `policy.json` and `message.json` below.

## Ready-to-run policy

Replace `YOUR_IMAGE_DIGEST` with your image (e.g. `hello-world@sha256:...`). Then `PUT /v1/sessions/{id}/policy` with this body. See `policy.json` in this folder.

```json
{
  "overrides": [
    {
      "tool": "docker.run",
      "allow": true,
      "constraints": {
        "images": ["YOUR_IMAGE_DIGEST"],
        "network_allowed": false,
        "mounts_allowed": false
      }
    }
  ]
}
```

## Ready-to-run intent

Message body for `POST /v1/sessions/{id}/messages` (replace `YOUR_IMAGE_DIGEST`):

```json
{
  "role": "user",
  "content": "Run container",
  "intents": [
    {
      "tool": "docker.run",
      "params": {
        "image": "YOUR_IMAGE_DIGEST",
        "cmd": []
      }
    }
  ]
}
```

## Run

1. Edit `policy.json` and `message.json` in this folder: set `YOUR_IMAGE_DIGEST` to your image (e.g. from `docker inspect`).
2. From repo root (PowerShell):

```powershell
.\examples\docker-skill-demo\run.ps1
```

Or set the digest in the script or env and run. Then open the UI → Sessions → run to see container output in step details.
