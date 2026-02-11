# UI Architecture

## Stack
- Svelte 4/5 (depending on your preference; default SvelteKit-style patterns without requiring SvelteKit)
- Vite
- TypeScript
- Tailwind CSS (optional but recommended for speed; keep styles consistent)
- Fetch-based API client (no heavy SDK needed)

> If you want SSR later, you can migrate to SvelteKit. MVP is a SPA.

## Project layout
```
/ui
  /src
    /app
      App.svelte
      routes.ts
      store.ts
    /components
      layout/
      forms/
      tables/
      badges/
      timeline/
      modals/
    /pages
      Login.svelte
      Dashboard.svelte
      Sessions.svelte
      SessionDetail.svelte
      RunDetail.svelte
      Policies.svelte
      Skills.svelte
      Audit.svelte
      Replay.svelte
    /lib
      api.ts
      types.ts
      format.ts
      validate.ts
      securityHints.ts
  index.html
  vite.config.ts
  tsconfig.json
```

## State management (MVP)
- Use Svelte stores:
  - `authStore` (apiBase, token)
  - `sessionsStore` (list cache)
  - `activeSessionStore` (session details/messages/runs)
  - `uiStore` (toasts, modals, loading states)

## API client
Centralize in `src/lib/api.ts`:
- Inject `Authorization: Bearer <token>`
- Standard error handling per backend spec
- Typed request/response (in `types.ts`)

## Routing
Lightweight router:
- Either `svelte-spa-router`
- Or a simple hash router in `routes.ts`
MVP recommendation: `svelte-spa-router` for speed.

## Security UX helpers
Create helper layer:
- `securityHints.ts` maps backend `POLICY_DENIED` reasons to friendly explanations
- “Fix safely” suggestions (e.g., allowlist specific domain, limit file root)

