<script lang="ts">
  import { onMount } from 'svelte'
  import { authStore, toastStore } from '../app/store'
  import { queryAudit, validateAuditChain } from '../lib/api'
  import type { AuditEvent } from '../lib/types'
  import Timeline from '../components/Timeline.svelte'

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
  <h1>Audit</h1>
  {#if !$authStore}
    <p>Not connected. <a href="#/login">Login</a></p>
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
      <button on:click={load}>Refresh</button>
      <button on:click={doValidateChain} disabled={validating} class="validate">{validating ? 'Validating…' : 'Validate chain'}</button>
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
  .page { max-width: 960px; margin: 1rem auto; padding: 1rem; }
  .toolbar { display: flex; gap: 0.5rem; margin: 1rem 0; flex-wrap: wrap; }
  .toolbar input { padding: 0.4rem; width: 180px; }
  .toolbar select { padding: 0.4rem; min-width: 200px; }
  .toolbar .validate { margin-left: 0.5rem; }
  .chain-result { padding: 0.5rem 1rem; border-radius: 6px; margin-bottom: 1rem; font-weight: 500; }
  .chain-result.valid { background: #e8f5e9; border: 1px solid #81c784; }
  .chain-result.invalid { background: #ffebee; border: 1px solid #e57373; }
  .chain-badge { font-weight: 600; margin-right: 0.25rem; }
  .chain-badge.ok { color: #2e7d32; }
  .chain-badge.broken { color: #c62828; }
  .view-toggle { margin-left: 0.5rem; }
  .view-toggle button { margin-right: 0.25rem; }
  .view-toggle button.active { font-weight: 600; background: #e3f2fd; }
  table { width: 100%; border-collapse: collapse; font-size: 0.85rem; }
  th, td { text-align: left; padding: 0.35rem; border-bottom: 1px solid #eee; vertical-align: top; }
  .badge { font-size: 0.75rem; padding: 0.15rem 0.4rem; border-radius: 4px; background: #e0e0e0; }
  .lifecycle-intent .badge { background: #bbdefb; }
  .lifecycle-decision .badge { background: #c8e6c9; }
  .lifecycle-capability .badge { background: #fff9c4; }
  .lifecycle-tool .badge { background: #b2dfdb; }
  .lifecycle-run .badge { background: #d1c4e9; }
  .lifecycle-session .badge { background: #f0f4c3; }
  .lifecycle-message .badge { background: #e1bee7; }
  .type { font-size: 0.75rem; }
  .data { margin: 0; font-size: 0.7rem; max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .hash { color: #666; }
</style>
