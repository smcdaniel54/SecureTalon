# Security Model

## Core principle
**The model is not a security boundary.**  
Security must be enforced by hard technical constraints:
- Capability tokens (cryptographic)
- Broker mediation
- Container isolation
- Signed skill artifacts
- Immutable audit log

---

## Capability-based security (CBS)

### Capability token
A capability is a short-lived, cryptographically signed token granting a single action (or narrow set of actions).

**Fields (minimum):**
- `cap_id` (uuid)
- `session_id`
- `subject` (agent/service identity)
- `tool` (e.g., `file.read`, `http.fetch`, `docker.run`, `shell.exec`)
- `constraints` (args allowlist, regex patterns, resource limits)
- `iat`, `exp`
- `nonce`
- `signature`

**Properties:**
- **Least privilege**: tokens grant minimal scope.
- **Time-boxed**: short TTL by default (e.g., 60 seconds).
- **Non-transferrable**: bound to session + subject.
- **Revocable**: broker checks revocation list (in-memory for MVP, persisted later).

### Token issuance
Policy Engine is the only issuer.  
Tool Broker is the only verifier/executor.

### Constraints examples
- `file.read` allowed only under `/workspace/projects/foo/**` with max 1MB.
- `http.fetch` allowed only to `https://api.example.com/*` with GET only.
- `docker.run` allowed only images in allowlist, no host mounts, network disabled unless explicitly granted.

---

## Prompt-injection resistance
All inbound text is **data**, not instructions.
Rules:
- Agent Runtime must label external content as `UNTRUSTED_TEXT`.
- Planner may summarize it, but **policy** cannot be changed due to it.
- Any tool request that originates from untrusted text must be treated as **high risk** (deny or require approval).

---

## Secrets handling
- Secrets are stored in OS keychain/vault or `.secrets/` with strict file perms (MVP can use local encrypted file).
- Secrets are referenced by handle, not raw value:
  - Tool request contains `secret_ref: "slack_bot_token"`
  - Broker resolves and injects into execution env.
- Secrets never:
  - appear in logs (redacted)
  - appear in LLM context
  - are returned in tool output

---

## Skill supply-chain security
- Skills must have a **manifest** and a signed artifact.
- Only allow execution if:
  - signature verifies against trusted keys
  - version is pinned
  - image hash matches allowlist
- No “curl | bash” style install in official workflow.

---

## Sandbox policy
Default execution environment for skills:
- Docker container
- read-only root filesystem
- no privileged mode
- drop all Linux capabilities
- no host mounts (unless explicitly allowed by token)
- network disabled (unless explicitly allowed by token)
- resource limits (CPU/mem/time)

Phase 2 options:
- gVisor for stronger syscall isolation
- Firecracker microVM for untrusted high-risk skills

---

## Audit & non-repudiation
Every security-relevant action must be logged:
- Policy evaluation decisions (allow/deny)
- Token issuance (hash only, not full token)
- Tool execution (inputs metadata, output hash, exit status)
- Skill runs (image digest, resource constraints)

Audit log is append-only; tamper evidence via hash chaining (see 06-AUDIT-AND-REPLAY.md).

