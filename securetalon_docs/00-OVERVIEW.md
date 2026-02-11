# SecureTalon – Security‑First Agent Platform (Design Pack)

This folder contains the **authoritative markdown design spec** for SecureTalon’s **backend-first** implementation.
It is written to be directly consumable by **Cursor AI** to scaffold and implement the system end‑to‑end.

## What SecureTalon is
SecureTalon is a **security-first**, local‑first agent automation platform inspired by OpenClaw’s momentum, but designed for **real deployment**:
- **Zero-trust defaults**
- **Capability-based permissions**
- **Tool Broker** as the only path to host OS/network
- **Sandboxed skills** (Docker) by default
- **Signed skills** (curated hub)
- **Immutable audit log + replay**

## Scope of this design pack
This pack defines:
- Core services and boundaries
- Data models (capabilities, sessions, audit events)
- REST API contracts
- Skill execution model (Docker runner)
- Security hardening requirements
- Implementation milestones and acceptance tests

> Note: A separate UI pack (Svelte/Vite) will be created after backend MVP is functional. This pack includes only the minimal endpoints/UI considerations needed to support the future UI.

## Quick start for Cursor AI
1. Read **01-ARCHITECTURE.md**
2. Implement **02-SECURITY-MODEL.md**
3. Build services in **03-SERVICES.md**
4. Implement APIs from **04-API-SPEC.md**
5. Enable Docker execution from **05-DOCKER-RUNNER.md**
6. Add audit/replay from **06-AUDIT-AND-REPLAY.md**
7. Run through tests in **07-ACCEPTANCE-TESTS.md**

## Non-goals (for MVP)
- Multi-tenant SaaS
- Complex plugin marketplace UX
- Fine-grained OS sandboxing beyond Docker (gVisor/Firecracker are phase 2)
- Advanced anomaly detection (phase 2; see roadmap)

