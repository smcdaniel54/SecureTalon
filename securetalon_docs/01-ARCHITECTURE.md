# Architecture

## High-level components
SecureTalon is built around a **hardened control plane** and **isolated execution plane**.

### Control Plane (Go binary)
- **Gateway API**: REST endpoints for sessions, skills, policies, audits.
- **Ingress adapters** (optional in MVP): Slack/Discord/Email connectors can post messages into sessions later.
- **Session Router**: groups messages and agent runs by session.
- **Planner / Agent Runtime**: prepares “intents” (tool requests) and executes skill flows.
- **Policy Engine**: evaluates every tool/skill request and issues **capability tokens**.
- **Tool Broker**: the only component allowed to execute OS commands, access files, or call network resources.
- **Audit Log**: append-only event stream of all actions.
- **Replay Engine**: deterministic replay of prior runs using the audit stream.

### Execution Plane
- **Skill Runner (Docker)**: executes a skill in an isolated container.
- Skills are **untrusted by default**. They must request actions through Tool Broker.

> IMPORTANT: The LLM / agent planner is *never* trusted to enforce security. Only the Policy Engine + Tool Broker enforce.

---

## Trust boundaries (zero-trust)
**Untrusted inputs**:
- Inbound messages from any channel
- Skill outputs
- Any retrieved “memory” or external documents
- Model outputs

**Trusted components**:
- Policy Engine (non-LLM)
- Tool Broker (non-LLM)
- Audit Log (append-only)
- Host OS (under our control)

---

## Request flow (happy path)
1. Client/UI sends request to `/v1/sessions/{id}/messages`
2. Session Router appends message to session and emits audit event
3. Agent Runtime produces a **Plan** and a set of **ToolIntents**
4. For each ToolIntent, Agent Runtime calls Policy Engine: `Evaluate(intent, sessionContext)`
5. Policy Engine returns:
   - `ALLOW` + a **capability token** scoped to that intent (tool, args pattern, TTL)
   - or `DENY` + reason + suggested least-privilege alternative
   - or `REQUIRE_APPROVAL` (phase 2; can stub now)
6. Agent Runtime submits ToolIntent + token to Tool Broker
7. Tool Broker validates token, executes action (or calls Docker Skill Runner), returns structured result
8. Results appended to session and audit log; Agent Runtime continues.

---

## Default mode expectations
- **Deny by default**: no tool executes without a token
- **Sandbox by default**: skills run in containers; host mounts prohibited by default
- **No secrets to LLM**: secrets are stored and injected only into Tool Broker/runner, never returned as plaintext
- **Auditable by default**: every step logs a structured event (inputs, decision, outputs, hashes)

---

## Suggested repo layout
```
/cmd/securetalon            # main
/internal/
  api/                      # http handlers, request/response types
  core/                     # session router, agent runtime orchestrator
  policy/                   # policy engine, rules, token issuer/verifier
  broker/                   # tool broker + runners (docker, local safe ops)
  audit/                    # append-only store + query
  replay/                   # deterministic replay
  skills/                   # skill manifests, registry client, signatures
  config/                   # config parsing & defaults
  auth/                     # local auth (admin token), RBAC later
/docs/                      # this pack
```

