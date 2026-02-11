<script lang="ts">
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore, toastStore } from '../app/store'
  import { listSessions, createSession as apiCreateSession } from '../lib/api'
  import type { Session } from '../lib/types'
  import PageHeader from '../components/PageHeader.svelte'

  let sessions: Session[] = []
  let loading = true
  let showCreate = false
  let newLabel = ''
  let creating = false

  function load() {
    const auth = $authStore
    if (!auth) return
    loading = true
    listSessions(auth, 100)
      .then((r) => { sessions = r })
      .catch((e: Error & { details?: unknown }) => {
        toastStore.push('error', e.message, e.details)
      })
      .finally(() => { loading = false })
  }

  async function doCreate() {
    const auth = $authStore
    if (!auth) return
    creating = true
    try {
      const s = await apiCreateSession(auth, newLabel || 'New session')
      sessions = [s, ...sessions]
      showCreate = false
      newLabel = ''
      toastStore.push('success', 'Session created')
      push(`/sessions/${s.id}`)
    } catch (e: unknown) {
      const err = e as Error & { details?: unknown }
      toastStore.push('error', err.message ?? 'Failed', err.details)
    } finally {
      creating = false
    }
  }

  onMount(load)
</script>

<div class="page">
  <PageHeader title="Sessions" subtitle="List and open agent sessions." />
  <button class="primary" on:click={() => showCreate = true}>New session</button>
  {#if showCreate}
    <div class="card-inline">
      <input type="text" bind:value={newLabel} placeholder="Label" />
      <button class="primary" on:click={doCreate} disabled={creating}>{creating ? 'Creating…' : 'Create'}</button>
      <button class="secondary" on:click={() => { showCreate = false; newLabel = '' }}>Cancel</button>
    </div>
  {/if}
  {#if loading}
    <p class="muted">Loading…</p>
  {:else}
    <table>
      <thead><tr><th>Label</th><th>Status</th><th>Created</th><th></th></tr></thead>
      <tbody>
        {#each sessions as s}
          <tr>
            <td>{s.label || s.id}</td>
            <td>{s.status}</td>
            <td>{new Date(s.created_at).toLocaleString()}</td>
            <td><a href="#/sessions/{s.id}" on:click|preventDefault={() => push(`/sessions/${s.id}`)}>Open</a></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .card-inline {
    margin: var(--space-4) 0;
    padding: var(--space-4);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    background: var(--bg-elevated);
    display: flex;
    gap: var(--space-2);
    align-items: center;
  }
  .card-inline input { flex: 1; }
  .muted { color: var(--text-muted); }
</style>
