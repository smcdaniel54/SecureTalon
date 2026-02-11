# Acceptance Tests (MVP)

These are black-box tests to confirm the backend is “secure by default” and usable.

## 1) Deny-by-default works
- Start SecureTalon with no policy rules allowing shell.
- Attempt `shell.exec` via tool intent.
- Expect: `POLICY_DENIED` and audit event `policy.decision` with `DENY`.

## 2) Capability token required
- Attempt `file.read` directly on Tool Broker endpoint without token.
- Expect: `401/403` and no file access.

## 3) Token constraints enforced
- Create policy allowing `file.read` under `/work/allowed/**`
- Attempt read `/etc/passwd`
- Expect: denied by broker with constraint violation

## 4) Skills run sandboxed
- Register `hello-world` skill image
- Run it via `docker.run`
- Verify:
  - `--network=none` (default)
  - no mounts are present
  - process runs as non-root if supported

## 5) Skill cannot access host
- Skill tries to read `/etc/hostname`
- Expect failure (read-only container root + no host mount)
- Audit logs record attempt

## 6) Audit log is complete and hash-chained
- Run a session with at least 1 tool execution
- Export audit events
- Verify hash chain validates end-to-end

## 7) Safe replay works
- Replay the last run in safe mode
- Ensure no tools are re-executed (no docker invocation)
- Outputs match recorded hashes

## 8) Secrets not leaked
- Configure a secret `demo_token=abc123`
- Run a tool requiring that secret (e.g., http auth header)
- Verify:
  - audit logs redact secret values
  - API responses never include secret plaintext

