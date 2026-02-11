# Screenshots for README

These PNGs are shown in the main [README](../../README.md). Keep them up to date so visitors see the real UI.

| File | Content |
|------|---------|
| `dashboard.png` | Dashboard: recent sessions, quick links. |
| `policy-editor.png` | Policy editor: session selector, effective policy, session overrides, tool constraints, guardrails. |
| `audit-chain-ok.png` | Audit page: filters, “Chain OK” badge, table or timeline. |
| `replay-viewer.png` | Replay page: run ID input, “Load Safe Replay”, timeline with step-through (prev/next, jump to type). |

**Recommended size:** 1200×800 (or ~1200px wide) so the README table looks good.

---

## How to capture real screenshots

### Option A: Manual (most reliable)

1. **Start the backend** (from repo root):
   ```powershell
   $env:ADDR = ":8090"; $env:ADMIN_TOKEN = "demo"; go run ./cmd/securetalon
   ```

2. **Start the UI** (in another terminal):
   ```powershell
   cd ui
   npm run dev
   ```
   Note the URL (e.g. http://localhost:5173).

3. **Open the UI in your browser**, go to the login page, enter:
   - API Base URL: `http://localhost:8090`
   - Admin Token: `demo`
   Click **Connect**.

4. **Capture each screen** at **1200×800** (or resize browser to ~1200px wide and full height, then crop to 800px tall if needed):
   - **Dashboard** — open `/#/` (home). Screenshot → save as `docs/screenshots/dashboard.png`.
   - **Policy editor** — open `/#/policies`. Screenshot → save as `docs/screenshots/policy-editor.png`.
   - **Audit** — open `/#/audit`. Optionally click “Validate chain” so “Chain OK” shows. Screenshot → save as `docs/screenshots/audit-chain-ok.png`.
   - **Replay** — open `/#/replay`. Screenshot → save as `docs/screenshots/replay-viewer.png`.

5. **Replace** the existing PNGs in `docs/screenshots/` with your four files.

### Option B: Automated (Playwright)

With **backend** and **UI** already running:

```powershell
cd ui
npm install
npx playwright install chromium
$env:UI_BASE_URL = "http://localhost:5173"   # or your UI port
$env:API_BASE_URL = "http://localhost:8090"
npm run capture-screenshots
```

Screenshots are written to `docs/screenshots/`. If the script fails (e.g. login form not found in headless), use **Option A** instead.

---

## One-shot script (PowerShell)

From repo root you can try:

```powershell
.\scripts\capture-readme-screenshots.ps1
```

This starts the backend and UI, runs the Playwright capture, then stops the servers. If automation fails in your environment, use manual capture (Option A) above.
