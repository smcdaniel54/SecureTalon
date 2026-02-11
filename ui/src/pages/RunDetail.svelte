<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore, toastStore } from '../app/store'
  import { getRun } from '../lib/api'
  import type { Run } from '../lib/types'
  import PageHeader from '../components/PageHeader.svelte'

  export let params: { id?: string } = {}
  $: runId = params.id ?? ''

  let run: Run | null = null
  let loading = true
  let pollTimer: ReturnType<typeof setInterval> | null = null

  function load(isPoll = false) {
    const auth = $authStore
    if (!auth || !runId) return
    if (!isPoll) loading = true
    getRun(auth, runId)
      .then((r) => {
        run = r as Run
        if ((r as Run).status === 'running' || (r as Run).status === 'queued') {
          if (!pollTimer) pollTimer = setInterval(() => load(true), 1500)
        } else if (pollTimer) {
          clearInterval(pollTimer)
          pollTimer = null
        }
      })
      .catch((e: Error & { details?: unknown }) => {
        if (!isPoll) toastStore.push('error', e.message, e.details)
      })
      .finally(() => { if (!isPoll) loading = false })
  }

  $: if (runId) load()

  onDestroy(() => {
    if (pollTimer) clearInterval(pollTimer)
  })
</script>

<div class="page">
  {#if !$authStore}
    <p class="muted">Not connected. <a href="#/login">Login</a></p>
  {:else if !runId}
    <p class="muted">Run ID required.</p>
  {:else if loading && !run}
    <p class="muted">Loading…</p>
  {:else if run}
    <PageHeader title="Run {run.id}" subtitle="Steps and status for this run." />
    <p>
      Status: <strong class="status" class:running={run.status === 'running' || run.status === 'queued'}>{run.status}</strong>
      · Started: {new Date(run.started_at).toLocaleString()}
      {#if run.ended_at} · Ended: {new Date(run.ended_at).toLocaleString()}{/if}
    </p>
    <p><a href="#/sessions/{run.session_id}" on:click|preventDefault={() => push(`/sessions/${run!.session_id}`)}>← Session</a></p>
    <h2>Steps</h2>
    {#if run.steps && run.steps.length > 0}
      <ul class="steps">
        {#each run.steps as step}
          <li class="step" class:denied={step.status === 'denied'} class:error={step.status === 'error'} class:ok={step.status === 'ok'}>
            <span class="step-id">{step.step_id}</span>
            <span class="step-type">{step.type}</span>
            <span class="step-status">{step.status}</span>
            {#if step.tool}<span class="step-tool">{step.tool}</span>{/if}
            {#if step.details && Object.keys(step.details).length}
              <pre class="details">{JSON.stringify(step.details, null, 2)}</pre>
            {/if}
          </li>
        {/each}
      </ul>
    {:else if run.status === 'running' || run.status === 'queued'}
      <p class="muted">Steps will appear as the agent runs…</p>
    {:else}
      <p class="muted">No steps recorded.</p>
    {/if}
  {/if}
</div>

<style>
  .page { max-width: 720px; }
  .steps { list-style: none; padding: 0; }
  .step { padding: var(--space-3); margin: var(--space-2) 0; border-left: 4px solid var(--border); background: var(--bg-elevated); border-radius: var(--radius-sm); }
  .step.denied { border-left-color: var(--error); }
  .step.error { border-left-color: var(--warning); }
  .step.ok { border-left-color: var(--success); }
  .status.running { color: var(--accent); }
  .muted { color: var(--text-muted); }
  .step-id { font-weight: 600; margin-right: var(--space-2); }
  .step-type { margin-right: var(--space-2); color: var(--text-muted); }
  .details { font-size: 0.85rem; overflow: auto; margin: var(--space-2) 0 0; padding: var(--space-2); background: var(--bg-input); border: 1px solid var(--border); border-radius: var(--radius-sm); }
</style>
