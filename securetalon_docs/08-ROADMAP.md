# Roadmap

## MVP (Secure plumbing first)
- Go gateway + auth token
- Sessions/messages/runs
- Policy engine: deny-by-default + simple allowlist rules
- Capability tokens: sign/verify + TTL
- Tool broker: file.read/file.write (safe), http.fetch (allowlist), docker.run
- Docker hardening defaults
- Audit JSONL + hash chain
- Safe replay

## Phase 2 (Viral + enterprise-ready)
- Svelte/Vite UI dashboard (sessions, runs, policies, audit timeline)
- Approval workflows (REQUIRE_APPROVAL)
- Curated skill hub + signing keys management
- RBAC (admin/operator/viewer)
- gVisor option for higher isolation
- Domain-specific policy packs (Slack bot, CRM sync, lead gen)

## Phase 3 (Security intelligence)
- Sentinel agents (advisory): anomaly detection, risk scoring
- Policy recommendations (human-approved)
- Alerting (email/webhook)
- Compliance report exports

