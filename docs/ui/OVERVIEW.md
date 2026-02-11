# SecureTalon UI Pack – Svelte + Vite (Design + Implementation Guide)

This pack defines a **first-class UI** for SecureTalon built with **Svelte + Vite**.
It is designed to be **backend-driven** and to consume the REST APIs defined in the backend design pack.

## Goals
- Make SecureTalon **drop-dead simple** for devs + power users.
- Provide **clarity + confidence** in security decisions:
  - “What ran?”
  - “Why allowed/denied?”
  - “What capabilities were granted?”
  - “What changed?”
- Provide a **tight, modern console**:
  - Sessions, Runs, Policies, Skills, Audit Timeline, Replay.

## Non-goals (UI MVP)
- Multi-tenant org management
- Advanced charts (phase 2)
- Marketplace purchasing UX (phase 2)
- Complex role provisioning (phase 2; MVP uses admin token only)

## UI MVP screens
1. **Login / Connect** (enter API base + admin token; stored locally)
2. **Dashboard** (recent sessions, running runs, alerts summary)
3. **Sessions** (list + create + open session)
4. **Session Detail** (messages + runs + quick actions)
5. **Run Detail** (step timeline + tool calls + results)
6. **Policy Editor** (effective policy + session overrides; guardrails)
7. **Skills** (list + register skill; signature status)
8. **Audit Explorer** (filter/search events; hash-chain validation status)
9. **Replay Viewer** (safe replay playback)

## UI principles (security-first UX)
- **Safe defaults**: “Deny by default” is visible and explicit.
- **Explainability**: every deny shows “why” + “how to fix safely.”
- **Progressive disclosure**: simple view first; advanced drawers for details.
- **Template-driven layouts**: consistent pages, forms, tables, modals.

