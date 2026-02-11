<script lang="ts">
  import { onMount } from 'svelte'
  import { authStore } from '../app/store'
  import { toastStore } from '../app/store'
  import { listSessions, getEffectivePolicy, putSessionPolicy, getRun } from '../lib/api'
  import type { Session, EffectivePolicy, RuleOverride } from '../lib/types'
  import PageHeader from '../components/PageHeader.svelte'
  import Card from '../components/Card.svelte'
  import FormPanel from '../components/FormPanel.svelte'

  let sessions: Session[] = []
  let sessionId = ''
  let policy: EffectivePolicy | null = null
  let loading = true
  let err = ''
  let saving = false
  let dangerConfirm = ''
  let fixSuggestion: { tool: string; hint: string } | null = null

  const TOOLS = ['file.read', 'file.write', 'http.fetch', 'docker.run'] as const

  function getQuery(): URLSearchParams {
    if (typeof window === 'undefined') return new URLSearchParams()
    return new URLSearchParams(window.location.search || window.location.hash.split('?')[1] || '')
  }

  function loadSessions() {
    const auth = $authStore
    if (!auth) return
    listSessions(auth, 200)
      .then((r) => { sessions = r })
      .catch((e) => { toastStore.push('error', e.message, (e as { details?: Record<string, unknown> }).details) })
  }

  function loadPolicy() {
    const auth = $authStore
    if (!auth) return
    policy = null
    if (!sessionId) { loading = false; return }
    loading = true
    err = ''
    getEffectivePolicy(auth, sessionId)
      .then((r) => { policy = r })
      .catch((e) => { err = e.message; toastStore.push('error', e.message, (e as { details?: Record<string, unknown> }).details) })
      .finally(() => { loading = false })
  }

  $: if (sessionId) loadPolicy()

  onMount(() => {
    loadSessions()
    const q = getQuery()
    const sid = q.get('session_id')
    const rid = q.get('run_id')
    if (sid) sessionId = sid
    if (sid && rid && $authStore) {
      getRun($authStore, rid).then((run) => {
        const denied = run.steps?.find((s) => s.type === 'policy_eval' && s.status === 'denied')
        if (denied?.details?.reason) {
          const tool = denied.tool || 'unknown'
          let hint = denied.details.reason as string
          if (tool === 'http.fetch' && hint.includes('domain')) hint = 'Add the requested domain to the allowlist.'
          else if ((tool === 'file.read' || tool === 'file.write') && hint.includes('path')) hint = 'Add the path root to allowed_roots (narrow scope).'
          fixSuggestion = { tool, hint }
        }
      }).catch(() => {})
    }
  })

  function emptyAllowlist(o: RuleOverride): boolean {
    if (!o.allow) return false
    const c = o.constraints || {}
    if (o.tool === 'http.fetch') return !(Array.isArray(c.domains) && c.domains.length > 0)
    if (o.tool === 'file.read' || o.tool === 'file.write') return !(Array.isArray(c.roots) && c.roots.length > 0)
    if (o.tool === 'docker.run') return !(Array.isArray(c.images) && c.images.length > 0)
    return false
  }

  function hasRiskyOptions(o: RuleOverride): boolean {
    const c = o.constraints || {}
    return !!(c.network_allowed || c.mounts_allowed || o.tool === 'shell.exec')
  }

  async function save() {
    const auth = $authStore
    if (!auth || !sessionId || !policy) return
    const withEmpty = policy.overrides.filter(emptyAllowlist)
    if (withEmpty.length > 0) {
      toastStore.push('error', 'Cannot save: an allowed tool has an empty allowlist (add domains, roots, or images).')
      return
    }
    const risky = policy.overrides.filter(hasRiskyOptions)
    if (risky.length > 0 && dangerConfirm !== 'I UNDERSTAND') {
      toastStore.push('error', 'Danger zone: type I UNDERSTAND to confirm risky options.')
      return
    }
    if (!confirm('Update session policy? This affects what tools the agent can run.')) return
    saving = true
    try {
      await putSessionPolicy(auth, sessionId, policy.overrides)
      toastStore.push('success', 'Policy saved.')
      dangerConfirm = ''
    } catch (e) {
      toastStore.push('error', (e as Error).message, (e as { details?: Record<string, unknown> }).details)
    } finally {
      saving = false
    }
  }

  function addOverride() {
    if (!policy) return
    policy = {
      ...policy,
      overrides: [...(policy.overrides || []), { tool: 'file.read', allow: true, constraints: { roots: [], max_bytes: 1048576 } }],
    }
  }

  function removeOverride(index: number) {
    if (!policy) return
    policy = { ...policy, overrides: policy.overrides.filter((_, i) => i !== index) }
  }

  function chipsToArray(s: string): string[] {
    return s.split(/[\n,]/).map((x) => x.trim()).filter(Boolean)
  }
  function arrayToChips(a: unknown): string {
    return Array.isArray(a) ? a.join(', ') : ''
  }

  function setRoots(o: RuleOverride, val: string) {
    if (!o.constraints) o.constraints = {}
    o.constraints.roots = chipsToArray(val)
    if (typeof o.constraints.max_bytes !== 'number') o.constraints.max_bytes = 1048576
    policy = policy ? { ...policy, overrides: [...(policy.overrides || [])] } : null
  }
  function setDomains(o: RuleOverride, val: string) {
    if (!o.constraints) o.constraints = {}
    o.constraints.domains = chipsToArray(val)
    if (typeof o.constraints.max_bytes !== 'number') o.constraints.max_bytes = 200000
    policy = policy ? { ...policy, overrides: [...(policy.overrides || [])] } : null
  }
  function setImages(o: RuleOverride, val: string) {
    if (!o.constraints) o.constraints = {}
    o.constraints.images = chipsToArray(val)
    policy = policy ? { ...policy, overrides: [...(policy.overrides || [])] } : null
  }
  function setMethods(o: RuleOverride, val: string) {
    if (!o.constraints) o.constraints = {}
    o.constraints.methods = val.split(',').map((x) => x.trim()).filter(Boolean)
    policy = policy ? { ...policy, overrides: [...(policy.overrides || [])] } : null
  }
</script>

<div class="page">
  <PageHeader title="Policies" subtitle="Session overrides for tool access. Deny by default." />
  {#if !$authStore}
    <p>Not connected. <a href="#/login">Login</a></p>
  {:else}
    <div class="deny-banner">Deny by default — no tool runs without an explicit session override.</div>
    <FormPanel title="Session">
      <select id="session" bind:value={sessionId}>
        <option value="">-- Select session --</option>
        {#each sessions as s}
          <option value={s.id}>{s.label || s.id}</option>
        {/each}
      </select>
    </FormPanel>
    {#if fixSuggestion}
      <div class="fix-suggestion">
        <strong>Fix safely:</strong> {fixSuggestion.hint} (from run deny)
      </div>
    {/if}
    {#if sessionId}
      {#if loading}
        <p>Loading policy…</p>
      {:else if err}
        <p class="error">{err}</p>
      {:else if policy}
        <Card title="Effective policy (read-only)">
          <p><strong>Default:</strong> {policy.default}</p>
        </Card>
        <Card title="Session overrides (editable)">
          {#if policy.overrides && policy.overrides.length > 0}
            <ul class="overrides">
              {#each policy.overrides as override, i}
                <li class="override">
                  <select bind:value={override.tool}>
                    {#each TOOLS as t}
                      <option value={t}>{t}</option>
                    {/each}
                  </select>
                  <label><input type="checkbox" bind:checked={override.allow} /> Allow</label>
                  {#if override.tool === 'file.read' || override.tool === 'file.write'}
                    <FormPanel title="Allowed roots (one per line or comma)">
                      <input type="text" value={arrayToChips(override.constraints?.roots)} on:input={(e) => setRoots(override, e.currentTarget.value)} placeholder="/work" />
                    </FormPanel>
                    <label>max_bytes <input type="number" bind:value={override.constraints.max_bytes} min="1" /></label>
                  {:else if override.tool === 'http.fetch'}
                    <FormPanel title="Domains allowlist">
                      <input type="text" value={arrayToChips(override.constraints?.domains)} on:input={(e) => setDomains(override, e.currentTarget.value)} placeholder="api.example.com" />
                    </FormPanel>
                    <label>Methods (comma) <input type="text" value={Array.isArray(override.constraints?.methods) ? (override.constraints.methods as string[]).join(', ') : 'GET'} on:input={(e) => setMethods(override, e.currentTarget.value)} placeholder="GET, POST" /></label>
                    <label>max_bytes <input type="number" bind:value={override.constraints.max_bytes} min="1" /></label>
                  {:else if override.tool === 'docker.run'}
                    <FormPanel title="Allowed images (digest only, one per line)">
                      <input type="text" value={arrayToChips(override.constraints?.images)} on:input={(e) => setImages(override, e.currentTarget.value)} placeholder="repo/name@sha256:..." />
                    </FormPanel>
                    <label><input type="checkbox" bind:checked={override.constraints.network_allowed} /> Network allowed (danger)</label>
                    <label><input type="checkbox" bind:checked={override.constraints.mounts_allowed} /> Mounts allowed (danger)</label>
                  {/if}
                  <button type="button" class="remove" on:click={() => removeOverride(i)}>Remove</button>
                </li>
              {/each}
            </ul>
          {:else}
            <p class="muted">No overrides. All tools are denied by default.</p>
          {/if}
          <button on:click={addOverride}>Add override</button>
        </Card>
        {#if policy.overrides?.some(hasRiskyOptions)}
          <div class="danger-zone">
            <strong>Danger zone:</strong> You enabled network, mounts, or shell. Type <code>I UNDERSTAND</code> to enable save.
            <input type="text" bind:value={dangerConfirm} placeholder="I UNDERSTAND" />
          </div>
        {/if}
        <button class="primary" on:click={save} disabled={saving}>{saving ? 'Saving…' : 'Save overrides'}</button>
      {/if}
    {/if}
  {/if}
</div>

<style>
  .page { max-width: 680px; }
  .deny-banner { background: var(--warning-subtle); border: 1px solid var(--warning); border-radius: var(--radius-sm); padding: var(--space-2) var(--space-4); margin-bottom: var(--space-4); font-weight: 500; }
  .fix-suggestion { background: var(--accent-subtle); border-radius: var(--radius-sm); padding: var(--space-2) var(--space-4); margin-bottom: var(--space-4); font-size: 0.9rem; }
  .overrides { list-style: none; padding: 0; }
  .override { padding: var(--space-3); margin: var(--space-2) 0; background: var(--bg); border: 1px solid var(--border); border-radius: var(--radius-sm); }
  .override select { margin-right: var(--space-2); }
  .remove { margin-left: var(--space-2); color: var(--error); background: none; border: none; cursor: pointer; font-size: 0.85rem; }
  .danger-zone { background: var(--error-subtle); border: 1px solid var(--error); border-radius: var(--radius-sm); padding: var(--space-4); margin: var(--space-4) 0; }
  .danger-zone input { margin-top: var(--space-2); width: 100%; max-width: 200px; }
  .muted { color: var(--text-muted); }
  .error { color: var(--error); }
  .primary { margin-top: var(--space-2); }
</style>
