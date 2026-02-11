# Security

This document describes SecureTalon’s threat model, how to report vulnerabilities, and what the project does and does not guarantee. It is intended to support enterprise evaluation and deployment decisions.

## Threat model

**What we protect against**

- **Unauthorized tool execution** — No tool (file, HTTP, Docker) runs unless the policy engine has issued a capability token for that session, tool, and constraints. The broker is the only component that executes tools and it verifies the token before running anything.
- **Policy bypass** — Policy is enforced in process; tool intents are denied by default. Only explicit session overrides (with constraints) result in token issuance. The LLM or client cannot grant itself broader access by persuasion or prompt injection at the policy layer.
- **Audit tampering** — Audit events are appended to a hash chain. Tampering or reordering can be detected by validating the chain (`GET /v1/audit/validate`). Replay is read-only and does not re-execute tools.
- **Overly broad tool scope** — Constraints (paths, domains, methods, images, network/mounts) are enforced by the broker against the intent parameters before execution. For example, a token that allows `file.read` under `work/` cannot be used to read outside that subtree.

**Assumptions**

- **Admin token** — Anyone with `ADMIN_TOKEN` (or equivalent) can create sessions, set policy, and trigger runs. Protect this credential; treat the API as trusted-admin only unless you add your own auth layer.
- **Deployment** — SecureTalon does not enforce TLS, network segmentation, or host hardening. You are responsible for deploying the backend and UI in a secure environment (HTTPS, firewall, etc.).
- **Docker** — When the broker runs containers, it relies on the host’s Docker daemon. Isolation is at the container level; use standard Docker security practices and consider runtimes (e.g. gVisor) for higher assurance where needed.

**Out of scope (for this threat model)**

- Security of the LLM itself (training data, prompt extraction, etc.).
- Security of third-party skills or images you allow in policy (we enforce allowlists and constraints, not the integrity of the code inside the image).
- Supply-chain security of the Go toolchain, Node, or other build-time dependencies (follow your organization’s practices).

## Responsible disclosure

We take security vulnerabilities seriously and ask that you report them in a way that allows us to fix issues before they are made public.

**How to report**

- **Do not** open a public GitHub issue for a security vulnerability.
- Report privately: use [GitHub Security Advisories](https://github.com/smcdaniel54/SecureTalon/security/advisories/new) (if enabled) or contact the maintainers through a private channel (e.g. contact information in the repo or your existing relationship).
- Include a clear description, steps to reproduce, and impact. If you have a suggested fix or mitigation, we welcome it.

**What we do**

- We will acknowledge receipt of your report in a timely manner.
- We will work to confirm the issue and develop a fix or mitigation.
- We will coordinate with you on disclosure timing and will credit you in the advisory or release notes if you wish (unless you prefer to remain anonymous).

We do not support bug bounties or compensation at this time; we appreciate responsible disclosure as a community contribution.

## What SecureTalon guarantees

Within the assumptions above, SecureTalon is designed to provide the following guarantees:

- **Deny-by-default execution** — No tool runs without a matching session policy override that allows that tool with constraints. Unlisted tools or out-of-scope parameters result in denial and an audit record.
- **Broker-mediated execution only** — All file, HTTP, and Docker execution goes through the broker. The broker verifies the capability token and constraints before performing the action. There is no “back door” for the agent to run tools without going through the policy engine and broker.
- **Short-lived, scoped tokens** — Capability tokens are time-limited and bound to session, subject, tool, and constraints. They are not reusable for other tools or broader scope.
- **Tamper-evident audit** — The audit log is hash-chained. You can verify that events have not been altered or reordered using the validation endpoint. Safe replay reconstructs the timeline from the log without re-executing any tools.
- **Constraint enforcement** — Paths, domains, methods, and image allowlists are enforced by the broker against the intent parameters. A token that allows only certain paths or domains cannot be used to exceed those limits at execution time.

These guarantees apply to the current design and implementation; specific versions may have known limitations documented in release notes or issues.

## What SecureTalon does not guarantee

To set accurate expectations:

- **Correctness of policy configuration** — SecureTalon enforces the policy you configure; it does not guarantee that your policy is correct or complete. Misconfiguration (e.g. overly permissive roots or domains) can still lead to undesired access. Use the UI guardrails and review overrides carefully.
- **LLM or application logic** — We do not guarantee that the model or your application will only request intended tools or parameters. We guarantee that whatever is requested is enforced by policy and broker; we do not control what the model “decides” to ask for.
- **Security of allowed third-party content** — Allowed HTTP domains, Docker images, or file paths may host or contain malicious content. We enforce allowlists and constraints; we do not scan or attest to the safety of that content.
- **Availability or integrity of the host** — If the host or Docker daemon is compromised, an attacker may be able to bypass or influence execution. SecureTalon is not a substitute for host hardening, network security, or secrets management.
- **Cryptographic strength of token signing** — Token signing uses the configured secret (e.g. `TOKEN_SECRET` or `ADMIN_TOKEN`). Key strength, rotation, and storage are your responsibility.

For a deeper treatment of the security model (capability tokens, constraints, prompt-injection resistance, skills), see [docs/backend/SECURITY-MODEL.md](docs/backend/SECURITY-MODEL.md).
