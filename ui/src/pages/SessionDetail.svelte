<script lang="ts">
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore, toastStore } from '../app/store'
  import { getSession, listMessages, postMessage } from '../lib/api'
  import type { Session, Message, ToolIntent } from '../lib/types'
  import PageHeader from '../components/PageHeader.svelte'

  export let params: { id?: string } = {}
  $: sessionId = params.id ?? ''

  let session: Session | null = null
  let messages: Message[] = []
  let loading = true
  let messageContent = ''
  let sending = false
  let lastRunId = ''

  function load() {
    const auth = $authStore
    if (!auth || !sessionId) return
    loading = true
    Promise.all([
      getSession(auth, sessionId),
      listMessages(auth, sessionId),
    ])
      .then(([s, m]) => { session = s; messages = m.messages ?? [] })
      .catch((e: Error & { details?: unknown }) => {
        toastStore.push('error', e.message, e.details)
      })
      .finally(() => { loading = false })
  }

  async function sendMessage() {
    const auth = $authStore
    if (!auth || !sessionId || !messageContent.trim()) return
    sending = true
    try {
      const r = await postMessage(auth, sessionId, {
        role: 'user',
        content: messageContent.trim(),
      })
      lastRunId = (r as { run_id: string }).run_id
      messageContent = ''
      toastStore.push('success', 'Message sent')
      load()
    } catch (e: unknown) {
      const err = e as Error & { details?: unknown }
      toastStore.push('error', err.message ?? 'Failed', err.details)
    } finally {
      sending = false
    }
  }

  $: if (sessionId) load()
</script>

<div class="page">
  {#if !$authStore}
    <p class="muted">Not connected. <a href="#/login">Login</a></p>
  {:else if loading && !session}
    <p class="muted">Loading…</p>
  {:else if session}
    <PageHeader title={session.label || session.id} subtitle="Session messages and runs." />
    <p><a href="#/sessions" on:click|preventDefault={() => push('/sessions')}>← Sessions</a></p>
    <div class="messages">
      {#each messages as msg}
        <div class="msg" class:assistant={msg.role === 'assistant'} class:user={msg.role === 'user'}>
          <span class="role">{msg.role}</span>
          <span class="content">{msg.content}</span>
          {#if msg.run_id}
            <a href="#/runs/{msg.run_id}" on:click|preventDefault={() => push(`/runs/${msg.run_id}`)}>Run {msg.run_id}</a>
          {/if}
        </div>
      {/each}
    </div>
    <form on:submit|preventDefault={sendMessage} class="input-row">
      <input type="text" bind:value={messageContent} placeholder="Message..." disabled={sending} />
      <button type="submit" class="primary" disabled={sending}>{sending ? 'Sending…' : 'Send'}</button>
    </form>
    {#if lastRunId}
      <p><a href="#/runs/{lastRunId}" on:click|preventDefault={() => push(`/runs/${lastRunId}`)}>View run {lastRunId}</a></p>
    {/if}
  {/if}
</div>

<style>
  .page { max-width: 640px; }
  .messages { margin: var(--space-4) 0; min-height: 200px; }
  .msg { padding: var(--space-2); margin: var(--space-1) 0; border-radius: var(--radius); }
  .msg.user { background: var(--accent-subtle); }
  .msg.assistant { background: var(--bg-elevated); border: 1px solid var(--border); }
  .role { font-size: 0.75rem; color: var(--text-muted); margin-right: var(--space-2); }
  .input-row { display: flex; gap: var(--space-2); margin-top: var(--space-4); }
  .input-row input { flex: 1; }
  .muted { color: var(--text-muted); }
</style>
