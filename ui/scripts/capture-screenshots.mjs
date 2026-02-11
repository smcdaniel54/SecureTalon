#!/usr/bin/env node
/**
 * Capture real README screenshots: dashboard, policy-editor, audit-chain-ok, replay-viewer.
 *
 * Prereqs (must be running first):
 *   1. Backend: ADDR=:8090 ADMIN_TOKEN=demo go run ./cmd/securetalon   (from repo root)
 *   2. UI:      npm run dev   (in ui/)
 *
 * Then run: npm run capture-screenshots   (from ui/)
 * Or:       UI_BASE_URL=http://localhost:5173 API_BASE_URL=http://localhost:8090 npm run capture-screenshots
 *
 * Env: UI_BASE_URL (default http://localhost:5173), API_BASE_URL (default http://localhost:8090),
 *      ADMIN_TOKEN (default demo), SCREENSHOTS_DIR (default ../../docs/screenshots)
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
    javaScriptEnabled: true,
  });
  context.setDefaultTimeout(30000);
  const page = await context.newPage();

  const appUrl = UI_BASE.replace(/\/$/, '') + '/#/';

  try {
    // Load login, fill form, click Connect, wait for dashboard (real auth flow)
    await page.goto(UI_BASE.replace(/\/$/, '') + '/#/login', { waitUntil: 'networkidle', timeout: 25000 });
    await page.waitForTimeout(3000);
    const apiInput = page.locator('#apiBase');
    await apiInput.waitFor({ state: 'visible', timeout: 10000 });
    await apiInput.fill(API_BASE);
    await page.locator('#token').fill(TOKEN);
    await page.getByRole('button', { name: /Connect/i }).click();
    const loggedIn = await page.locator('aside.nav').waitFor({ state: 'visible', timeout: 25000 }).then(() => true).catch(() => false);
    if (!loggedIn) {
      await page.screenshot({ path: join(OUT_DIR, '_capture-debug.png'), fullPage: false }).catch(() => {});
      console.warn('Login did not complete; check _capture-debug.png.');
    } else {
      await page.waitForTimeout(2000);
    }

    // 1. Dashboard
    await page.goto(appUrl, { waitUntil: 'domcontentloaded', timeout: 10000 });
    await page.waitForTimeout(2000);
    await page.locator('main.main').evaluate((el) => el?.scrollIntoView({ block: 'start' })).catch(() => {});
    await page.waitForTimeout(500);
    await page.screenshot({ path: join(OUT_DIR, 'dashboard.png'), fullPage: false });
    console.log('Saved dashboard.png');

    // 2. Policy editor
    await page.goto(UI_BASE.replace(/\/$/, '') + '/#/policies', { waitUntil: 'domcontentloaded', timeout: 10000 });
    await page.waitForTimeout(2000);
    await page.locator('main.main').evaluate((el) => el?.scrollIntoView({ block: 'start' })).catch(() => {});
    await page.waitForTimeout(500);
    await page.screenshot({ path: join(OUT_DIR, 'policy-editor.png'), fullPage: false });
    console.log('Saved policy-editor.png');

    // 3. Audit
    await page.goto(UI_BASE.replace(/\/$/, '') + '/#/audit', { waitUntil: 'domcontentloaded', timeout: 10000 });
    await page.waitForTimeout(2000);
    const validateBtn = page.getByRole('button', { name: /Validate chain/i });
    if (await validateBtn.isVisible()) {
      await validateBtn.click();
      await page.waitForTimeout(2000);
    }
    await page.locator('main.main').evaluate((el) => el?.scrollIntoView({ block: 'start' })).catch(() => {});
    await page.waitForTimeout(500);
    await page.screenshot({ path: join(OUT_DIR, 'audit-chain-ok.png'), fullPage: false });
    console.log('Saved audit-chain-ok.png');

    // 4. Replay
    await page.goto(UI_BASE.replace(/\/$/, '') + '/#/replay', { waitUntil: 'domcontentloaded', timeout: 10000 });
    await page.waitForTimeout(2000);
    await page.locator('main.main').evaluate((el) => el?.scrollIntoView({ block: 'start' })).catch(() => {});
    await page.waitForTimeout(500);
    await page.screenshot({ path: join(OUT_DIR, 'replay-viewer.png'), fullPage: false });
    console.log('Saved replay-viewer.png');

  } catch (err) {
    await page.screenshot({ path: join(OUT_DIR, '_capture-debug.png'), fullPage: false }).catch(() => {});
    console.error('Capture failed. Ensure backend and UI are running, then retry.');
    throw err;
  } finally {
    await browser.close();
  }

  console.log('Done. Screenshots in', OUT_DIR);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
