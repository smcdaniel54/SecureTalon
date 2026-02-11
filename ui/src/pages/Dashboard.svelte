<script lang="ts">
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore, toastStore } from '../app/store'
  import { listSessions } from '../lib/api'
  import type { Session } from '../lib/types'

  let sessions: Session[] = []
  let loading = true

  function load() {
    const auth = $authStore
    if (!auth) return
    loading = true
    listSessions(auth, 10)
      .then((r) => { sessions = r })
      .catch((e: Error & { details?: unknown }) => {
        toastStore.push('error', e.message, e.details)
      })
      .finally(() => { loading = false })
  }

  onMount(() => { if ($authStore) load() })
</script>
{#if $authStore}
  <div class="dashboard">
    <header class="dashboard-header">
      <h1>Dashboard</h1>
      <p class="dashboard-desc">Overview of sessions and quick actions</p>
    </header>
    <section class="card">
      <h2>Recent Sessions</h2>
      {#if loading}
        <p class="muted">Loading…</p>
      {:else if sessions.length === 0}
        <p class="muted empty">No sessions yet. Create one from Sessions or run a demo.</p>
      {:else}
        <ul>
          {#each sessions as s}
            <li><a href="#/sessions/{s.id}" on:click|preventDefault={() => push(`/sessions/${s.id}`)}>{s.label || s.id}</a></li>
          {/each}
        </ul>
      {/if}
    </section>
    <div class="actions">
      <button class="primary" on:click={() => push('/sessions')}>Sessions</button>
      <button class="secondary" on:click={() => push('/audit')}>Audit</button>
    </div>
  </div>
{:else}
  <p class="not-connected">Not connected. <a href="#/login">Login</a></p>
{/if}

<style>
  .dashboard { max-width: 680px; }
  .dashboard-header {
    margin-bottom: var(--space-6);
    padding-bottom: var(--space-4);
    border-bottom: 3px solid var(--accent);
  }
  .dashboard-header h1 { margin-bottom: var(--space-1); }
  .dashboard-desc { color: var(--text-muted); font-size: 0.9rem; margin: 0; }
  .dashboard .card { margin-top: 0; }
  .dashboard .muted { color: var(--text-muted); font-size: 0.9rem; }
  .dashboard .muted.empty { margin: 0; }
  .dashboard ul { list-style: none; padding: 0; margin: 0; }
  .dashboard li {
    padding: var(--space-2) 0;
    border-bottom: 1px solid var(--border);
    transition: background var(--duration-fast);
  }
  .dashboard li:last-child { border-bottom: none; }
  .dashboard li:hover {
    background: var(--hover-overlay);
    margin: 0 calc(-1 * var(--space-2));
    padding-left: var(--space-2);
    padding-right: var(--space-2);
    border-radius: var(--radius-sm);
  }
  .dashboard .actions { margin-top: var(--space-6); display: flex; gap: var(--space-2); }
  .not-connected { color: var(--text-muted); }
</style>
