# SecureTalon

<p align="center">
  <img src="assets/securetalon-logo.png" alt="SecureTalon logo" width="280" />
</p>

[![CI](https://github.com/smcdaniel54/SecureTalon/actions/workflows/ci.yml/badge.svg)](https://github.com/smcdaniel54/SecureTalon/actions/workflows/ci.yml)

**SecureTalon** is a security-first agent platform inspired by OpenClaw. It executes AI tool intents only through least-privilege policies, broker-mediated tools, Docker-sandboxed skills, and tamper-evident audit logs—so agents remain **powerful**, **observable**, and **safe by design**.

## Why SecureTalon?

OpenClaw demonstrated the power of autonomous agents.  
**SecureTalon makes them safe to use in real environments.**

SecureTalon adds:

- **Deny-by-default policies** — No tool runs unless you explicitly allow it with constraints.
- **Broker-mediated tools** — Every file, HTTP, and Docker call goes through the broker; no back doors.
- **Docker-sandboxed skills** — Skills run in containers by default; network and mounts are opt-in.
- **Tamper-evident audit logs** — Hash-chained events; you can validate integrity and prove what ran.
- **Safe replay without re-execution** — Step through any run from the audit log; no tools run again.
- **Built-in admin console** — Sessions, policies, audit, and replay in one Svelte UI.

## How SecureTalon differs from OpenClaw

| Area | OpenClaw | SecureTalon |
|------|----------|-------------|
| Default posture | permissive | deny-by-default |
| Tool access | direct | broker-mediated |
| Execution isolation | optional | Docker-sandboxed by default |
| Policy enforcement | ad-hoc | capability tokens |
| Audit | basic logs | hash-chained tamper-evident |
| Replay | limited | safe replay without re-execution |
| UI console | minimal | built-in Svelte admin console |

### For OpenClaw users

Concept mapping so conversion feels easy:

| OpenClaw | SecureTalon |
|----------|-------------|
| Agent | **Run** — one execution per message; steps show policy_eval and tool_exec |
| Tool | **Brokered tool** — file.read, http.fetch, docker.run go through the broker with constraints |
| Plugin | **Skill** — Docker image by digest; register via API, allow via session policy |
| exec.security / allowlist | **Session policy overrides** — per-session allowlist with constraints (roots, domains, images) |
| exec.ask / approvals | **Deny-by-default** — add an override to allow; no interactive ask yet |
| Logs | **Audit** (hash-chained) + **Replay** (safe, no re-execution) |

Full guide: [OpenClaw migration](docs/OPENCLAW-MIGRATION.md) (what to change, example conversions).

## Screenshots

The admin console is a major differentiator. Below: dashboard, policy editor, audit chain validation, and replay viewer.

| [Dashboard](docs/screenshots/dashboard.png) | [Policy editor](docs/screenshots/policy-editor.png) |
|---------------------------------------------|-----------------------------------------------------|
| Sessions and quick links                     | Session overrides, tool constraints, guardrails     |

| [Audit — Chain OK](docs/screenshots/audit-chain-ok.png) | [Replay viewer](docs/screenshots/replay-viewer.png)   |
|--------------------------------------------------------|------------------------------------------------------|
| Hash chain valid, filters, timeline                     | Safe replay timeline, step-through, jump to type      |

*Add the actual PNGs under `docs/screenshots/` (dashboard.png, policy-editor.png, audit-chain-ok.png, replay-viewer.png) to show your UI in the repo.*

## Security model (summary)

- **Deny-by-default** — Only explicitly allowed tools and constraints are permitted per session.
- **Capability tokens** — Short-lived, signed tokens issued by the policy engine; the broker is the only verifier and executor.
- **Broker mediation** — All tool intents (file, HTTP, Docker) go through the broker; constraints (paths, domains, images) are enforced before execution.
- **Tamper-evident audit** — Events are appended to a hash chain; you can validate integrity and replay runs without re-executing tools.

Details: [docs/backend/SECURITY-MODEL.md](docs/backend/SECURITY-MODEL.md).

## Use cases

- **Internal agent ops** — Run file, HTTP, or containerized skills under strict allowlists with full audit and replay.
- **Compliance-sensitive workloads** — Prove what ran, when, and with what policy via hash-chained logs.
- **Safe skill adoption** — Register Docker skills by digest; enforce network/mount policies and inspect runs in the console.
- **Post-incident review** — Use safe replay and audit to step through a run without re-executing any tools.

## Examples

Ready-to-run demos in **[examples/](examples/)** with sample policies and intents you can copy or run:

| Demo | Description |
|------|-------------|
| [file-read-demo](examples/file-read-demo/) | `policy.json` + `run.ps1`; allows `file.read` under `work/`, sends intent, shows run. |
| [http-fetch-demo](examples/http-fetch-demo/) | `policy.json` + `message.json` + `run.ps1`; allows GET to a domain, sends intent, shows response. |
| [docker-skill-demo](examples/docker-skill-demo/) | `policy.json` + `message.json` + `run.ps1`; allows `docker.run` by image digest (set `IMAGE_DIGEST` or edit JSON). |
| [deny-then-fix](examples/deny-then-fix/) | Trigger a deny, add override, retry. |

Run from repo root: `.\examples\file-read-demo\run.ps1` (and similarly for the others). Backend must be running.

## Getting started

### Backend

Requires `ADMIN_TOKEN`. Optional: `ADDR` (default `:8080`), `DATA_DIR`, `TOKEN_SECRET`.

```bash
ADMIN_TOKEN=your-secret-token go run ./cmd/securetalon
```

To use port 8090 (e.g. for the UI default):

```bash
ADDR=:8090 ADMIN_TOKEN=your-secret-token go run ./cmd/securetalon
```

### UI

```bash
cd ui
npm install
npm run dev
```

Open http://localhost:5173. In the UI, **Connect** with:

- **API base:** `http://localhost:8090` (or your backend URL)
- **Token:** same as `ADMIN_TOKEN`

See [ui/README.md](ui/README.md) for API mapping and build.

### Happy path demo

1. Start backend: `ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon`
2. Start UI: `cd ui && npm i && npm run dev` → connect with `http://localhost:8090` and token `demo`
3. **Create session** (Sessions → New session).
4. **Set policy:** Policies → select the session → add override: tool `file.read`, allow, roots e.g. `work` (chips), Save.
5. Create a file the agent can read (e.g. `mkdir work 2>nul & echo hello > work\input.txt` on Windows, or `mkdir -p work && echo hello > work/input.txt` on Unix).
6. **Send message** with parseable intent: in Session detail, send e.g. “Read work/input.txt” or use a message that triggers a file.read intent; or call the API with body `{"content":"...", "intents":[{"tool":"file.read","params":{"path":"work/input.txt"}}]}`.
7. **View run:** open the run from the session; confirm steps (policy_eval ALLOW/DENY, tool_exec ok/error) and status (queued → running → completed).
8. **Audit:** Audit → set session filter → Refresh → Validate chain → “Chain OK”.
9. **Replay:** Replay → enter run ID → Load Safe Replay → step through timeline (no re-execution).

See [scripts/demo-happy-path.ps1](scripts/demo-happy-path.ps1) for an automated API-only demo (PowerShell).

## Roadmap

- **MVP (done)** — Deny-by-default policy, capability tokens, broker (file/http/Docker), hash-chained audit, safe replay, Svelte admin console.
- **Next** — Approval workflows, RBAC, skill hub and signing, gVisor option, policy packs. See [docs/backend/ROADMAP.md](docs/backend/ROADMAP.md).

## Contributing

We welcome issues and pull requests. See [CONTRIBUTING.md](CONTRIBUTING.md) for how to run the project, how to add tools or skills, coding standards, and CI (policy, token, and audit chain tests run on every push/PR).

## Security

We take security seriously. See **[SECURITY.md](SECURITY.md)** for our threat model, responsible disclosure process, and what SecureTalon does and does not guarantee. If you believe you’ve found a vulnerability, report it privately (e.g. GitHub Security Advisories or maintainer contact); do not open a public issue.

## Distribution

Downloadable zips for public release **must not include** the `.git` directory (fine in the live repo, not in artifacts). GitHub’s **Code → Download ZIP** already omits `.git`.

To build a clean distribution zip yourself (e.g. for releases), run from the repo root:

```powershell
.\scripts\make-dist-zip.ps1
```

This produces `SecureTalon-YYYYMMDD.zip` with source and docs, excluding `.git`, `node_modules`, `ui/dist`, `data`, and `.env`.

---

Further documentation: [docs/](docs/) (backend architecture, API, UI specs).

**License:** [MIT](LICENSE)
