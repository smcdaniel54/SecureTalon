# Docker skill example

Register a skill (Docker image by digest) and allow `docker.run` for that image; then send a message that runs the container.

**Prereqs:**

- Backend on `http://localhost:8090` with `ADMIN_TOKEN=demo`.
- Docker available to the backend (the broker runs containers).
- An image you can run by digest (e.g. `hello-world` or a small utility). Get digest: `docker inspect --format='{{index .RepoDigests 0}}' hello-world` (or use `image@sha256:...` from a registry).

**Steps:**

1. **Register skill (optional but recommended):**  
   `POST /v1/skills` with body:
   ```json
   {
     "name": "hello-world",
     "version": "1.0",
     "image": "hello-world@sha256:...",
     "manifest": {}
   }
   ```
   Use the actual digest from your environment.

2. **Create session:** `POST /v1/sessions` → `{"label":"Docker example","metadata":{}}`.

3. **Set policy:** `PUT /v1/sessions/{id}/policy` with body:
   ```json
   {
     "overrides": [{
       "tool": "docker.run",
       "allow": true,
       "constraints": {
         "images": ["hello-world@sha256:..."],
         "network_allowed": false,
         "mounts_allowed": false
       }
     }]
   }
   ```
   Replace `hello-world@sha256:...` with your image digest.

4. **Post message with intent:** `POST /v1/sessions/{id}/messages` with:
   ```json
   {
     "role": "user",
     "content": "Run hello-world",
     "intents": [{
       "tool": "docker.run",
       "params": {
         "image": "hello-world@sha256:...",
         "cmd": []
       }
     }]
   }
   ```

5. **Get run:** `GET /v1/runs/{run_id}` — steps: `policy_eval` (allow), `tool_exec` (ok) with container output in details.

**Run the script (after setting image digest):**

Edit `run.ps1` and set `$ImageDigest` to your image (e.g. from `docker inspect` or registry). Then from repo root:

```powershell
.\examples\docker-skill\run.ps1
```

**Note:** If the broker cannot run Docker (e.g. no Docker daemon), `tool_exec` may return an error; policy and audit/replay still demonstrate the flow.
