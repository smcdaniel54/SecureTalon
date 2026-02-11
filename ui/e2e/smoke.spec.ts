import { test } from 'playwright/test'
import { join } from 'path'

const API_BASE = process.env.API_BASE_URL || 'http://localhost:8080'
const TOKEN = process.env.ADMIN_TOKEN || 'demo'
const SCREENSHOTS_DIR = process.env.SCREENSHOTS_DIR || join(process.cwd(), '..', 'docs', 'screenshots')

async function apiPost(path: string, body: object) {
  const res = await fetch(`${API_BASE.replace(/\/$/, '')}${path}`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${TOKEN}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  })
  if (!res.ok) throw new Error(`API ${path}: ${res.status} ${await res.text()}`)
  return res.json() as Promise<Record<string, unknown>>
}

test.describe('smoke + screenshots', () => {
  let runId: string

  test.beforeAll(async () => {
    const session = await apiPost('/v1/sessions', { label: 'E2E smoke', metadata: {} }) as { id: string }
    const msg = await apiPost(`/v1/sessions/${session.id}/messages`, { content: 'Hello' }) as { run_id: string }
    runId = msg.run_id
  })

  test('full smoke: login, then capture each screen at most complete state', async ({ page }) => {
    await page.goto('/#/login')
    await page.getByLabel(/API Base URL/i).fill(API_BASE)
    await page.getByLabel(/Admin Token/i).fill(TOKEN)
    await page.getByRole('button', { name: /Connect/i }).click()
    await page.locator('aside.nav').waitFor({ state: 'visible', timeout: 15000 })

    // 1. Dashboard — with sessions list
    await page.goto('/#/')
    await page.getByText('Recent Sessions').waitFor({ state: 'visible', timeout: 10000 })
    await page.waitForTimeout(1500)
    await page.screenshot({ path: join(SCREENSHOTS_DIR, 'dashboard.png'), fullPage: false })

    // 2. Policies — page loaded
    await page.goto('/#/policies')
    await page.getByRole('heading', { name: /Policies/i }).waitFor({ state: 'visible', timeout: 10000 })
    await page.waitForTimeout(1500)
    await page.screenshot({ path: join(SCREENSHOTS_DIR, 'policy-editor.png'), fullPage: false })

    // 3. Audit — click Validate chain; capture when result appears or after short wait
    await page.goto('/#/audit')
    await page.getByRole('heading', { name: /Audit/i }).waitFor({ state: 'visible', timeout: 10000 })
    const validateBtn = page.getByRole('button', { name: /Validate chain/i })
    if (await validateBtn.isVisible()) {
      await validateBtn.click()
      await page.locator('.chain-result').waitFor({ state: 'visible', timeout: 20000 }).catch(() => {})
    }
    await page.waitForTimeout(1500)
    await page.screenshot({ path: join(SCREENSHOTS_DIR, 'audit-chain-ok.png'), fullPage: false })

    // 4. Replay — timeline loaded
    await page.goto('/#/replay')
    await page.getByPlaceholder(/run_/).fill(runId)
    await page.getByRole('button', { name: /Load Safe Replay/i }).click()
    await page.getByText('Timeline').waitFor({ state: 'visible', timeout: 15000 })
    await page.waitForTimeout(1500)
    await page.screenshot({ path: join(SCREENSHOTS_DIR, 'replay-viewer.png'), fullPage: false })
  })
})
