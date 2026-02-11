<script lang="ts">
  import { authStore, toastStore } from '../app/store'
  import { safeReplay } from '../lib/api'
  import type { ReplayResponse } from '../lib/types'
  import PageHeader from '../components/PageHeader.svelte'

  let runId = ''
  let loading = false
  let result: ReplayResponse | null = null
  let currentIndex = 0

  $: events = result?.events ?? []
  $: currentEvent = events[currentIndex] ?? null
  $: eventTypes = [...new Set(events.map((e) => e.type))]
  $: timelineItems = events.map((ev) => ({
    ts: ev.ts,
    type: ev.type,
    label: ev.type,
    status: 'ok',
    data: ev.data,
  }))
  $: highlightedIndex = currentIndex

  async function loadReplay() {
    const auth = $authStore
    if (!auth || !runId.trim()) {
      toastStore.push('error', 'Enter a run ID')
      return
    }
    loading = true
    result = null
    currentIndex = 0
    try {
      const r = await safeReplay(auth, runId.trim())
      result = r
    } catch (e: unknown) {
      const err = e as Error & { details?: unknown }
      toastStore.push('error', err.message ?? 'Replay failed', err.details)
    } finally {
      loading = false
    }
  }

  function prev() {
    if (currentIndex > 0) currentIndex -= 1
  }
  function next() {
    if (currentIndex < events.length - 1) currentIndex += 1
  }
  function jumpToType(type: string) {
    const i = events.findIndex((e) => e.type === type)
    if (i >= 0) currentIndex = i
  }
</script>

<div class="page">
  <PageHeader title="Replay (safe)" subtitle="View run timeline from audit log. No tools are re-executed — visualization only." />
  {#if !$authStore}
    <p class="muted">Not connected. <a href="#/login">Login</a></p>
  {:else}
    <div class="toolbar">
      <input type="text" bind:value={runId} placeholder="run_..." />
      <button class="primary" on:click={loadReplay} disabled={loading}>{loading ? 'Loading…' : 'Load Safe Replay'}</button>
    </div>
    {#if result}
      <section class="result card-section">
        <p>
          <strong>Run:</strong> {result.run_id} · <strong>Mode:</strong> {result.mode}
          · Chain: {result.valid ? '✓ Valid' : '✗ Invalid'}
        </p>
        <h2>Timeline</h2>
        <div class="step-controls">
          <button class="secondary" on:click={prev} disabled={currentIndex <= 0}>← Prev</button>
          <span class="step-info">Step {currentIndex + 1} of {events.length}</span>
          <button class="secondary" on:click={next} disabled={currentIndex >= events.length - 1}>Next →</button>
          <span class="jump-label">Jump to type:</span>
          <select
            value={currentEvent?.type ?? ''}
            on:change={(e) => jumpToType((e.currentTarget as HTMLSelectElement).value)}
          >
            {#each eventTypes as t}
              <option value={t}>{t}</option>
            {/each}
          </select>
        </div>
        {#if currentEvent}
          <div class="current-event">
            <strong>{currentEvent.type}</strong> · {new Date(currentEvent.ts).toLocaleString()}
            {#if Object.keys(currentEvent.data).length > 0}
              <pre class="ev-data">{JSON.stringify(currentEvent.data, null, 2)}</pre>
            {/if}
          </div>
        {/if}
        <div class="timeline-wrap">
          {#each timelineItems as item, i}
            <div class="timeline-item-wrap" class:current={i === currentIndex}>
              <button
                class="timeline-item-btn"
                on:click={() => currentIndex = i}
              >
                <span class="ts">{new Date(item.ts).toLocaleString()}</span>
                <span class="type">{item.type}</span>
              </button>
              {#if i === currentIndex && item.data && Object.keys(item.data).length > 0}
                <pre class="data">{JSON.stringify(item.data, null, 2)}</pre>
              {/if}
            </div>
          {/each}
        </div>
        {#if events.length === 0}
          <p class="muted">No events for this run.</p>
        {/if}
      </section>
    {/if}
  {/if}
</div>

<style>
  .page { max-width: 720px; }
  .toolbar { display: flex; gap: var(--space-2); margin-bottom: var(--space-4); }
  .toolbar input { width: 220px; }
  .card-section { border: 1px solid var(--border); border-radius: var(--radius); padding: var(--space-4); margin-top: var(--space-4); background: var(--bg-elevated); }
  .step-controls { display: flex; align-items: center; gap: var(--space-2); margin: var(--space-4) 0; flex-wrap: wrap; }
  .step-info { font-size: 0.9rem; color: var(--text-muted); }
  .jump-label { margin-left: var(--space-2); font-size: 0.9rem; color: var(--text-muted); }
  .step-controls select { min-width: 180px; }
  .current-event { padding: var(--space-3); margin-bottom: var(--space-4); background: var(--bg); border: 1px solid var(--border); border-radius: var(--radius-sm); }
  .current-event pre { margin: var(--space-2) 0 0; font-size: 0.85rem; }
  .timeline-wrap { margin-top: var(--space-4); }
  .timeline-item-wrap { margin-bottom: var(--space-1); }
  .timeline-item-wrap.current { box-shadow: inset 3px 0 0 var(--accent); padding-left: var(--space-2); margin-left: calc(-1 * var(--space-2)); }
  .timeline-item-btn { display: block; width: 100%; text-align: left; padding: var(--space-2); background: var(--bg-elevated); border: 1px solid var(--border); border-radius: var(--radius-sm); cursor: pointer; font-size: 0.9rem; transition: background var(--duration-fast); }
  .timeline-item-btn:hover { background: var(--hover-overlay); }
  .timeline-item-btn .ts { color: var(--text-muted); margin-right: var(--space-2); font-size: 0.85rem; }
  .data { font-size: 0.8rem; margin: var(--space-1) 0 var(--space-2); padding: var(--space-2); background: var(--bg-input); border: 1px solid var(--border); border-radius: var(--radius-sm); overflow: auto; }
  .muted { color: var(--text-muted); }
</style>
