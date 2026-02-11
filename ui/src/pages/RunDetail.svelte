<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore, toastStore } from '../app/store'
  import { getRun } from '../lib/api'
  import type { Run } from '../lib/types'

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
    <p>Not connected. <a href="#/login">Login</a></p>
  {:else if !runId}
    <p>Run ID required.</p>
  {:else if loading && !run}
    <p>Loading…</p>
  {:else if run}
    <h1>Run {run.id}</h1>
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
  .page { max-width: 720px; margin: 1rem auto; padding: 1rem; }
  .steps { list-style: none; padding: 0; }
  .step { padding: 0.75rem; margin: 0.5rem 0; border-left: 4px solid #ddd; background: #fafafa; border-radius: 4px; }
  .step.denied { border-left-color: #f44336; }
  .step.error { border-left-color: #ff9800; }
  .step.ok { border-left-color: #4caf50; }
  .status.running { color: #2196f3; }
  .muted { color: #666; }
  .step-id { font-weight: 600; margin-right: 0.5rem; }
  .step-type { margin-right: 0.5rem; color: #666; }
  .details { font-size: 0.85rem; overflow: auto; margin: 0.5rem 0 0; padding: 0.5rem; background: #fff; border: 1px solid #eee; }
</style>
