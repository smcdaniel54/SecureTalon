<script lang="ts">
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore, toastStore } from '../app/store'
  import { getSession, listMessages, postMessage } from '../lib/api'
  import type { Session, Message, ToolIntent } from '../lib/types'

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
    <p>Not connected. <a href="#/login">Login</a></p>
  {:else if loading && !session}
    <p>Loading…</p>
  {:else if session}
    <h1>{session.label || session.id}</h1>
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
      <button type="submit" disabled={sending}>{sending ? 'Sending…' : 'Send'}</button>
    </form>
    {#if lastRunId}
      <p><a href="#/runs/{lastRunId}" on:click|preventDefault={() => push(`/runs/${lastRunId}`)}>View run {lastRunId}</a></p>
    {/if}
  {/if}
</div>

<style>
  .page { max-width: 640px; margin: 1rem auto; padding: 1rem; }
  .messages { margin: 1rem 0; min-height: 200px; }
  .msg { padding: 0.5rem; margin: 0.25rem 0; border-radius: 8px; }
  .msg.user { background: #e3f2fd; }
  .msg.assistant { background: #f5f5f5; }
  .role { font-size: 0.75rem; color: #666; margin-right: 0.5rem; }
  .input-row { display: flex; gap: 0.5rem; margin-top: 1rem; }
  .input-row input { flex: 1; padding: 0.5rem; }
</style>
