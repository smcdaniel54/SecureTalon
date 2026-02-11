#!/usr/bin/env node
/**
 * Capture README screenshots: dashboard, policy-editor, audit-chain-ok, replay-viewer.
 *
 * Prereqs:
 *   - Backend running (e.g. ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon)
 *   - UI running (npm run dev in ui/)
 *   - Optional: create a session and run something so Dashboard/Audit have content
 *
 * Run from repo root:
 *   cd ui && npm run capture-screenshots
 *
 * Or from ui/:
 *   node scripts/capture-screenshots.mjs
 *
 * Env:
 *   UI_BASE_URL   default http://localhost:5173
 *   API_BASE_URL  default http://localhost:8090 (used for login)
 *   ADMIN_TOKEN   default demo
 *   SCREENSHOTS_DIR default ../../docs/screenshots (relative to ui/scripts/)
 */

import { chromium } from 'playwright';
import { mkdir } from 'fs/promises';
import { dirname, join, resolve } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const UI_BASE = process.env.UI_BASE_URL || 'http://localhost:5173';
const API_BASE = process.env.API_BASE_URL || 'http://localhost:8090';
const TOKEN = process.env.ADMIN_TOKEN || 'demo';
const OUT_DIR = resolve(__dirname, process.env.SCREENSHOTS_DIR || '../../docs/screenshots');

const VIEWPORT = { width: 1200, height: 800 };

async function main() {
  await mkdir(OUT_DIR, { recursive: true });

  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({
    viewport: VIEWPORT,
    ignoreHTTPSErrors: true,
  });
  const page = await context.newPage();

  try {
    // Login if needed
    await page.goto(UI_BASE + '/#/login', { waitUntil: 'networkidle' });
    const apiInput = page.locator('#apiBase');
    if (await apiInput.isVisible()) {
      await apiInput.fill(API_BASE);
      await page.locator('#token').fill(TOKEN);
      await page.getByRole('button', { name: /Connect/i }).click();
      await page.waitForTimeout(2000);
    }

    // 1. Dashboard
    await page.goto(UI_BASE + '/#/', { waitUntil: 'networkidle' });
    await page.waitForTimeout(800);
    await page.screenshot({ path: join(OUT_DIR, 'dashboard.png'), fullPage: false });
    console.log('Saved dashboard.png');

    // 2. Policy editor
    await page.goto(UI_BASE + '/#/policies', { waitUntil: 'networkidle' });
    await page.waitForTimeout(800);
    await page.screenshot({ path: join(OUT_DIR, 'policy-editor.png'), fullPage: false });
    console.log('Saved policy-editor.png');

    // 3. Audit (try to show "Chain OK" by clicking Validate if present)
    await page.goto(UI_BASE + '/#/audit', { waitUntil: 'networkidle' });
    await page.waitForTimeout(800);
    const validateBtn = page.getByRole('button', { name: /Validate chain/i });
    if (await validateBtn.isVisible()) {
      await validateBtn.click();
      await page.waitForTimeout(1500);
    }
    await page.screenshot({ path: join(OUT_DIR, 'audit-chain-ok.png'), fullPage: false });
    console.log('Saved audit-chain-ok.png');

    // 4. Replay viewer
    await page.goto(UI_BASE + '/#/replay', { waitUntil: 'networkidle' });
    await page.waitForTimeout(800);
    await page.screenshot({ path: join(OUT_DIR, 'replay-viewer.png'), fullPage: false });
    console.log('Saved replay-viewer.png');

  } finally {
    await browser.close();
  }

  console.log('Done. Screenshots in', OUT_DIR);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
