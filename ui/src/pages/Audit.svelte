<script lang="ts">
  import { onMount } from 'svelte'
  import { authStore, toastStore } from '../app/store'
  import { queryAudit, validateAuditChain } from '../lib/api'
  import type { AuditEvent } from '../lib/types'
  import Timeline from '../components/Timeline.svelte'
  import PageHeader from '../components/PageHeader.svelte'

  let events: AuditEvent[] = []
  let loading = true
  let sessionFilter = ''
  let runIdFilter = ''
  let typeFilter = ''
  let sinceFilter = ''
  let untilFilter = ''
  let viewMode: 'table' | 'timeline' = 'table'
  let chainValid: boolean | null = null
  let chainInvalidIndex = -1
  let chainEventCount = 0
  let validating = false

  const EVENT_LIFECYCLE: Record<string, string> = {
    'session.created': 'session',
    'message.appended': 'message',
    'run.started': 'run',
    'policy.intent.received': 'intent',
    'policy.decision': 'decision',
    'capability.issued': 'capability',
    'tool.executed': 'tool',
    'run.finished': 'run',
  }

  function lifecycleClass(ev: AuditEvent): string {
    const l = EVENT_LIFECYCLE[ev.type] ?? 'other'
    return `lifecycle-${l}`
  }

  $: timelineItems = events.map((ev) => ({
    ts: ev.ts,
    type: ev.type,
    label: ev.type,
    status: 'ok',
    data: ev.data,
  }))

  function load() {
    const auth = $authStore
    if (!auth) return
    loading = true
    chainValid = null
    queryAudit(auth, {
      limit: 500,
      session_id: sessionFilter || undefined,
      run_id: runIdFilter || undefined,
      type: typeFilter || undefined,
      since: sinceFilter || undefined,
      until: untilFilter || undefined,
    })
      .then((r) => { events = r })
      .catch((e: Error & { details?: unknown }) => {
        toastStore.push('error', e.message, e.details)
      })
      .finally(() => { loading = false })
  }

  async function doValidateChain() {
    const auth = $authStore
    if (!auth) return
    validating = true
    chainValid = null
    try {
      const r = await validateAuditChain(auth, {
        session_id: sessionFilter || undefined,
        limit: 500,
      })
      chainValid = r.valid
      chainInvalidIndex = r.invalid_index
      chainEventCount = r.event_count
    } catch (e: Error & { details?: unknown }) {
      toastStore.push('error', e.message, e.details)
      chainValid = null
    } finally {
      validating = false
    }
  }

  onMount(load)
</script>

<div class="page">
  <PageHeader title="Audit" subtitle="Query audit log and validate event chain integrity." />
  {#if !$authStore}
    <p class="muted">Not connected. <a href="#/login">Login</a></p>
  {:else}
    <div class="toolbar">
      <input type="text" bind:value={sessionFilter} placeholder="Session ID" />
      <input type="text" bind:value={runIdFilter} placeholder="Run ID (optional)" />
      <select bind:value={typeFilter}>
        <option value="">All types</option>
        <option value="session.created">session.created</option>
        <option value="message.appended">message.appended</option>
        <option value="run.started">run.started</option>
        <option value="policy.intent.received">policy.intent.received</option>
        <option value="policy.decision">policy.decision</option>
        <option value="capability.issued">capability.issued</option>
        <option value="tool.executed">tool.executed</option>
        <option value="run.finished">run.finished</option>
      </select>
      <input type="datetime-local" bind:value={sinceFilter} placeholder="Since" title="Since (ISO)" />
      <input type="datetime-local" bind:value={untilFilter} placeholder="Until" title="Until (ISO)" />
      <button class="secondary" on:click={load}>Refresh</button>
      <button class="primary" on:click={doValidateChain} disabled={validating}>{validating ? 'Validating…' : 'Validate chain'}</button>
      <span class="view-toggle">
        <button class:active={viewMode === 'table'} on:click={() => viewMode = 'table'}>Table</button>
        <button class:active={viewMode === 'timeline'} on:click={() => viewMode = 'timeline'}>Timeline</button>
      </span>
    </div>
    {#if chainValid !== null}
      <div class="chain-result" class:valid={chainValid} class:invalid={!chainValid}>
        {#if chainValid}
          <span class="chain-badge ok">Chain OK</span> ({chainEventCount} events)
        {:else}
          <span class="chain-badge broken">Broken at event #{chainInvalidIndex}</span> ({chainEventCount} events checked)
        {/if}
      </div>
    {/if}
    {#if loading}
      <p>Loading…</p>
    {:else}
      {#if viewMode === 'timeline'}
        <Timeline items={timelineItems} />
      {:else}
      <table>
        <thead><tr><th>Time</th><th>Lifecycle</th><th>Type</th><th>Session</th><th>Run</th><th>Data</th><th>Hash</th></tr></thead>
        <tbody>
          {#each events as ev}
            <tr class={lifecycleClass(ev)}>
              <td>{new Date(ev.ts).toLocaleString()}</td>
              <td><span class="badge">{EVENT_LIFECYCLE[ev.type] ?? '—'}</span></td>
              <td><code class="type">{ev.type}</code></td>
              <td><code>{ev.session_id}</code></td>
              <td><code>{ev.run_id || '—'}</code></td>
              <td><pre class="data">{JSON.stringify(ev.data)}</pre></td>
              <td><code class="hash">{ev.hash.slice(0, 12)}…</code></td>
            </tr>
          {/each}
        </tbody>
      </table>
      {#if events.length === 0}
        <p>No audit events.</p>
      {/if}
      {/if}
    {/if}
  {/if}
</div>

<style>
  .page { max-width: 960px; }
  .toolbar { display: flex; gap: var(--space-2); margin: var(--space-4) 0; flex-wrap: wrap; }
  .toolbar input { width: 180px; }
  .toolbar select { min-width: 200px; }
  .chain-result { padding: var(--space-2) var(--space-4); border-radius: var(--radius-sm); margin-bottom: var(--space-4); font-weight: 500; }
  .chain-result.valid { background: var(--success-subtle); border: 1px solid var(--success); }
  .chain-result.invalid { background: var(--error-subtle); border: 1px solid var(--error); }
  .chain-badge { font-weight: 600; margin-right: var(--space-1); }
  .chain-badge.ok { color: var(--success); }
  .chain-badge.broken { color: var(--error); }
  .view-toggle { margin-left: var(--space-2); }
  .view-toggle button { margin-right: var(--space-1); }
  .view-toggle button.active { font-weight: 600; background: var(--accent-subtle); color: var(--accent); }
  .page table { font-size: 0.85rem; }
  .page td { vertical-align: top; }
  .badge { font-size: 0.75rem; padding: 0.15rem 0.4rem; border-radius: var(--radius-sm); background: var(--hover-overlay); color: var(--text); }
  .lifecycle-intent .badge { background: var(--accent-subtle); color: var(--accent); }
  .lifecycle-decision .badge { background: var(--success-subtle); color: var(--success); }
  .lifecycle-capability .badge { background: var(--warning-subtle); color: var(--warning); }
  .lifecycle-tool .badge { background: var(--accent-soft); color: var(--accent-hover); }
  .lifecycle-run .badge { background: var(--accent-subtle); color: var(--accent); }
  .lifecycle-session .badge { background: var(--warning-subtle); color: var(--warning); }
  .lifecycle-message .badge { background: var(--accent-soft); color: var(--accent); }
  .type { font-size: 0.75rem; }
  .data { margin: 0; font-size: 0.7rem; max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .hash { color: var(--text-muted); }
  .muted { color: var(--text-muted); }
</style>
