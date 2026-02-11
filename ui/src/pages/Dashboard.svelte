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
    <h1>Dashboard</h1>
    <section class="card">
      <h2>Recent Sessions</h2>
      {#if loading}
        <p>Loading…</p>
      {:else if sessions.length === 0}
        <p>No sessions yet.</p>
      {:else}
        <ul>
          {#each sessions as s}
            <li><a href="#/sessions/{s.id}" on:click|preventDefault={() => push(`/sessions/${s.id}`)}>{s.label || s.id}</a></li>
          {/each}
        </ul>
      {/if}
    </section>
    <p><button on:click={() => push('/sessions')}>Sessions</button> <button on:click={() => push('/audit')}>Audit</button></p>
  </div>
{:else}
  <p>Not connected. <a href="#/login">Login</a></p>
{/if}

<style>
  .dashboard { max-width: 600px; margin: 1rem auto; padding: 1rem; }
  .card { border: 1px solid #ddd; border-radius: 8px; padding: 1rem; margin: 1rem 0; }
  ul { list-style: none; padding: 0; }
  li { padding: 0.25rem 0; }
  a { color: #06c; }
  button { margin-right: 0.5rem; padding: 0.4rem 0.8rem; cursor: pointer; }
</style>
