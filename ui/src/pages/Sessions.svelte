<script lang="ts">
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore, toastStore } from '../app/store'
  import { listSessions, createSession as apiCreateSession } from '../lib/api'
  import type { Session } from '../lib/types'

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
  <h1>Sessions</h1>
  <button on:click={() => showCreate = true}>New session</button>
  {#if showCreate}
    <div class="modal">
      <input type="text" bind:value={newLabel} placeholder="Label" />
      <button on:click={doCreate} disabled={creating}>{creating ? 'Creating…' : 'Create'}</button>
      <button class="secondary" on:click={() => { showCreate = false; newLabel = '' }}>Cancel</button>
    </div>
  {/if}
  {#if loading}
    <p>Loading…</p>
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
  .page { max-width: 800px; margin: 1rem auto; padding: 1rem; }
  .modal { margin: 1rem 0; padding: 1rem; border: 1px solid #ccc; border-radius: 8px; display: flex; gap: 0.5rem; align-items: center; }
  .modal input { flex: 1; padding: 0.5rem; }
  table { width: 100%; border-collapse: collapse; }
  th, td { text-align: left; padding: 0.5rem; border-bottom: 1px solid #eee; }
  button { padding: 0.4rem 0.8rem; cursor: pointer; }
  button.secondary { background: #f0f0f0; border: 1px solid #ccc; }
</style>
