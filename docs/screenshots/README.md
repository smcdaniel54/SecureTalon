# Screenshots for README

Add these PNG files here so the main [README](../../README.md) can display them:

| File | Content |
|------|---------|
| `dashboard.png` | Dashboard: recent sessions, quick links. |
| `policy-editor.png` | Policy editor: session selector, effective policy, session overrides, tool constraints, guardrails. |
| `audit-chain-ok.png` | Audit page: filters, “Chain OK” badge, table or timeline. |
| `replay-viewer.png` | Replay page: run ID input, “Load Safe Replay”, timeline with step-through (prev/next, jump to type). |

## Automated capture

From the repo root, with the **backend** and **UI** both running:

```bash
cd ui
npm install
npx playwright install chromium
npm run capture-screenshots
```

Screenshots are written to `docs/screenshots/` (1200×800 viewport). Optional env: `UI_BASE_URL`, `API_BASE_URL`, `ADMIN_TOKEN`, `SCREENSHOTS_DIR`.

## Manual capture

Run the UI (`cd ui && npm run dev`), connect to the backend, then capture each page at ~1200px wide.
