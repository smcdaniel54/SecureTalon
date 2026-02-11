<script lang="ts">
  import { authStore, toastStore } from '../app/store'
  import { safeReplay } from '../lib/api'
  import type { ReplayResponse } from '../lib/types'

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
  <h1>Replay (safe)</h1>
  {#if !$authStore}
    <p>Not connected. <a href="#/login">Login</a></p>
  {:else}
    <p class="hint">View run timeline from audit log. No tools are re-executed — visualization only.</p>
    <div class="toolbar">
      <input type="text" bind:value={runId} placeholder="run_..." />
      <button on:click={loadReplay} disabled={loading}>{loading ? 'Loading…' : 'Load Safe Replay'}</button>
    </div>
    {#if result}
      <section class="result">
        <p>
          <strong>Run:</strong> {result.run_id} · <strong>Mode:</strong> {result.mode}
          · Chain: {result.valid ? '✓ Valid' : '✗ Invalid'}
        </p>
        <h2>Timeline</h2>
        <div class="step-controls">
          <button on:click={prev} disabled={currentIndex <= 0}>← Prev</button>
          <span class="step-info">Step {currentIndex + 1} of {events.length}</span>
          <button on:click={next} disabled={currentIndex >= events.length - 1}>Next →</button>
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
          <div class="current-event card">
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
  .page { max-width: 720px; margin: 1rem auto; padding: 1rem; }
  .hint { color: #666; margin-bottom: 1rem; }
  .toolbar { display: flex; gap: 0.5rem; margin-bottom: 1rem; }
  .toolbar input { padding: 0.4rem; width: 220px; }
  .result { border: 1px solid #ddd; border-radius: 8px; padding: 1rem; margin-top: 1rem; }
  .step-controls { display: flex; align-items: center; gap: 0.5rem; margin: 1rem 0; flex-wrap: wrap; }
  .step-info { font-size: 0.9rem; color: #666; }
  .jump-label { margin-left: 0.5rem; font-size: 0.9rem; }
  .step-controls select { padding: 0.35rem; min-width: 180px; }
  .current-event { padding: 0.75rem; margin-bottom: 1rem; background: #f5f5f5; border-radius: 6px; }
  .current-event pre { margin: 0.5rem 0 0; font-size: 0.85rem; }
  .timeline-wrap { margin-top: 1rem; }
  .timeline-item-wrap { margin-bottom: 0.25rem; }
  .timeline-item-wrap.current { border-left: 3px solid #2196f3; padding-left: 0.5rem; margin-left: -0.5rem; }
  .timeline-item-btn { display: block; width: 100%; text-align: left; padding: 0.5rem; background: #fafafa; border: 1px solid #eee; border-radius: 4px; cursor: pointer; font-size: 0.9rem; }
  .timeline-item-btn:hover { background: #f0f0f0; }
  .timeline-item-btn .ts { color: #666; margin-right: 0.5rem; font-size: 0.85rem; }
  .data { font-size: 0.8rem; margin: 0.25rem 0 0.5rem; padding: 0.5rem; background: #fff; border: 1px solid #eee; overflow: auto; }
  .muted { color: #666; }
  .card { border: 1px solid #e0e0e0; }
</style>
