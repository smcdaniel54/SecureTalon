# Contributing

Thanks for your interest in SecureTalon. This guide covers how to run the project, where to find design docs, how to add tools and skills, and coding standards.

## Getting started

1. **Run locally** — See the [README](README.md) (backend + UI, happy path demo).
2. **Architecture & API** — See [docs/backend](docs/backend/) for design, security model, and API spec.
3. **UI** — See [docs/ui](docs/ui/) and [ui/README.md](ui/README.md) for the admin console.
4. **Changes** — Open an issue or PR; keep scope clear and avoid breaking the policy/audit invariants.
5. **Security** — See [SECURITY.md](SECURITY.md) for threat model and how to report vulnerabilities responsibly.

---

## How to add a new tool

Tools are executed only by the **broker** after the **policy engine** has issued a capability token. Adding a new tool requires changes in three places.

### 1. Policy engine (`internal/policy/engine.go`)

- If the tool should be **denied by default** with no allowlist (like `shell.exec`), add an early return in `Evaluate` that returns `DecisionDeny` with a reason.
- If the tool is **allowable via session overrides**, no engine change is needed for a generic tool name; overrides are matched by `r.Tool == intent.Tool` and the broker will enforce constraints.

### 2. Broker (`internal/broker/`)

- **Execute** — In `broker.go`, add a `case "your.tool":` in the `switch intent.Tool` and call your new `doYourTool(intent.Params, token.Constraints)`.
- **Constraint checks** — In `checkConstraints`, add a branch for your tool (e.g. require certain constraint keys, validate params against them). The broker must **never** run without satisfying constraints.
- **Implementation** — Add a new file (e.g. `your_tool.go`) with:
  - `func (b *Broker) doYourTool(params map[string]interface{}, constraints map[string]interface{}) (map[string]interface{}, error)`
  - Parse params, enforce constraints (allowlists, limits), perform the operation, return a result map or error.

### 3. Tests

- **Policy** — In `internal/policy/engine_test.go`, add tests for deny-by-default (if applicable) and allow-with-override for the new tool.
- **Broker** — In `internal/broker/broker_test.go`, add tests that verify: token required, wrong tool rejected, constraints enforced, and a happy-path execution (with suitable mocks or temp resources).

### 4. UI (optional)

- **Policies page** — If the tool needs tool-specific constraint fields (like file roots, HTTP domains, Docker images), add a branch in the policy editor in `ui/src/pages/Policies.svelte` so users can configure the new tool’s constraints and guardrails.

### 5. Docs

- Update [docs/backend/API-SPEC.md](docs/backend/API-SPEC.md) and [docs/OPENCLAW-MIGRATION.md](docs/OPENCLAW-MIGRATION.md) if you add a new tool or constraint shape.

---

## How to add a skill

**Skills** in SecureTalon are runnable units (today: Docker images by digest) that can be registered and then allowed via session policy.

### 1. Registering a skill (API)

- **Endpoint:** `POST /v1/skills`
- **Body:** `name`, `version`, `image` (must include `@sha256:...` digest), optional `signature`, `public_key_id`, `manifest`.
- The backend may persist skills for allowlisting and display; the UI uses this to list skills and show a “signed” badge. See `internal/api/handlers.go` (`RegisterSkill`, `ListSkills`).

### 2. Allowing a skill in policy

- Session policy overrides for **`docker.run`** use constraints: `images` (list of allowed image digests), `network_allowed`, `mounts_allowed`.
- To “add a skill” from a user’s perspective: register the skill (image by digest), then in the Policies UI (or `PUT /v1/sessions/{id}/policy`) add an override for `docker.run` with `constraints.images` including that digest.

### 3. Adding skill registration in the UI

- **Skills page** (`ui/src/pages/Skills.svelte`) — List comes from `GET /v1/skills`; register form calls `POST /v1/skills` with validation (e.g. image must contain `@sha256:`).
- **Types** — Add or extend types in `ui/src/lib/types.ts` and use the API helpers in `ui/src/lib/api.ts` (`listSkills`, `registerSkill`).

### 4. Docs and examples

- Update [docs/backend](docs/backend/) if the skill model changes. The [Docker skill example](examples/docker-skill/) shows register + policy + intent flow.

---

## Coding standards

### Go (backend)

- **Formatting** — Use `gofmt` or `goimports`. CI runs `go build` and `go test ./...`.
- **Tests** — Put tests in `*_test.go` next to the code. Prefer table-driven tests for multiple cases. Cover:
  - **Policy:** deny-by-default, allow with override, token present and tool/constraints correct.
  - **Tokens:** issue + verify, reject wrong secret, reject expired.
  - **Audit:** append events, hash chain linkage, `ValidateChain` detects tampering.
  - **Broker:** token required, constraint enforcement, tool execution (with temp dirs or mocks where needed).
- **Errors** — Return errors with context (`fmt.Errorf("...: %w", err)`). API handlers use `WriteError` for consistent JSON error shape.
- **Packages** — Keep `internal/` layout: `agent`, `api`, `audit`, `broker`, `config`, `core`, `policy`, `replay`. No circular imports.

### UI (Svelte + TypeScript)

- **API** — Use the shared client in `ui/src/lib/api.ts` for all backend calls. Do not ad-hoc `fetch`; use helpers (`listSessions`, `getEffectivePolicy`, `putSessionPolicy`, etc.) and surface errors via the toast store.
- **Types** — Keep types in `ui/src/lib/types.ts`; use them in components and API.
- **Errors** — On API failure, push to the toast store with message and optional details (expandable in the toast). Avoid leaving users with no feedback.
- **Components** — Reuse existing ones where possible: `AppShell`, `PageHeader`, `Card`, `DataTable`, `Timeline`, `Toast`, `Modal`, `FormPanel`. Prefer simple CSS or Tailwind; keep the UI fast and uncluttered.
- **Lint/build** — `npm run build` must pass. Fix any reported unused selectors or TypeScript errors.

### General

- **Commits** — Prefer small, logical commits. Message should be clear (e.g. “Add broker constraint check for X”, “Fix audit chain validation when events empty”).
- **PRs** — Describe what changed and why. Reference any issue. Ensure CI passes.

---

## CI

GitHub Actions run on push/PR:

- **Backend** — `go test ./...` (policy enforcement, token validation, audit chain, broker).
- **UI** — `cd ui && npm ci && npm run build`.

See [.github/workflows/ci.yml](.github/workflows/ci.yml). Fix failing tests or build before merging.
