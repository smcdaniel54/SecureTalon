# Cursor AI – Master Implementation Prompt (Frontend UI)

Copy/paste this prompt into Cursor AI with a new `ui/` folder in the SecureTalon repo.

---

## Objective
Build a **SecureTalon Admin Console UI** using **Svelte + Vite + TypeScript** (SPA).
The UI must be **security-first**, template-driven, and consume the backend APIs.

## Stack
- Svelte
- Vite
- TypeScript
- Router: `svelte-spa-router`
- Styling: simple CSS or Tailwind (recommended). If Tailwind, configure it properly.

## Create project
In `/ui`:
- Use `npm create vite@latest` with Svelte + TS
- Add `svelte-spa-router`

## App requirements
### Auth
- Login page to set:
  - `apiBase`
  - `adminToken`
- Save to `localStorage`
- Test connection on “Connect” (call `GET /v1/sessions?limit=1`)
- Provide “Forget” button

### Pages
Implement pages per [PAGE-SPECS.md](PAGE-SPECS.md):
- Login
- Dashboard
- Sessions
- Session Detail
- Run Detail
- Policies
- Skills
- Audit
- Replay

### Components/Templates
Implement reusable components from [TEMPLATES.md](TEMPLATES.md):
- AppShell, PageHeader, Card, DataTable, Timeline, FormPanel, Toasts, Modal

### API Client
Create `src/lib/api.ts`:
- `request<T>(method, path, body?)`
- Inject bearer token
- Parse standard error format:
  - show toast with error.message and error.details
- Types in `src/lib/types.ts` matching backend API responses

### Security UX
Implement guardrails from [UX-GUARDRAILS.md](UX-GUARDRAILS.md):
- Prevent saving insecure policies without explicit confirmations
- Provide “Fix safely” suggestions on DENY reasons
- Danger zone section for risky toggles (shell.exec)

### Audit hash chain validation
- On Audit page, validate `hash`/`prev_hash` chain client-side and show status.

### Dev experience
- `npm run dev` hot reload works
- Provide `.env.example` with `VITE_API_BASE` default

## Deliverables
- Working UI with navigation and all pages
- Clean, consistent design (no clutter)
- Detailed README in `/ui/README.md` including:
  - run commands
  - configuration
  - screenshots placeholders
  - how pages map to APIs

---

