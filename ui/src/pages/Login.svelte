<script lang="ts">
  import { authStore } from '../app/store'
  import { connect, getStoredAuth } from '../lib/api'

  let apiBase = getStoredAuth()?.apiBase ?? 'http://localhost:8080'
  let token = getStoredAuth()?.token ?? ''
  let error = ''
  let loading = false

  async function handleConnect() {
    error = ''
    if (!apiBase.trim() || !token.trim()) {
      error = 'API base and token are required'
      return
    }
    loading = true
    try {
      await connect({ apiBase: apiBase.trim(), token: token.trim() })
      authStore.setAuth({ apiBase: apiBase.trim(), token: token.trim() })
    } catch (e) {
      error = e instanceof Error ? e.message : 'Connection failed'
    } finally {
      loading = false
    }
  }

  function handleForget() {
    authStore.clear()
    apiBase = 'http://localhost:8080'
    token = ''
    error = ''
  }
</script>

<div class="login">
  <img src="/securetalon-logo.png" alt="SecureTalon" class="login-logo" />
  <h1>SecureTalon</h1>
  <p class="subtitle">Connect to your API</p>
  <form on:submit|preventDefault={handleConnect}>
    <label for="apiBase">API Base URL</label>
    <input id="apiBase" type="url" bind:value={apiBase} placeholder="http://localhost:8080" />
    <label for="token">Admin Token</label>
    <input id="token" type="password" bind:value={token} placeholder="Bearer token" />
    {#if error}
      <p class="error">{error}</p>
    {/if}
    <div class="actions">
      <button type="submit" disabled={loading}>{loading ? 'Connecting…' : 'Connect'}</button>
      <button type="button" class="secondary" on:click={handleForget}>Forget</button>
    </div>
  </form>
</div>

<style>
  .login {
    max-width: 400px;
    margin: 2rem auto;
    padding: 1.5rem;
    text-align: center;
  }
  .login-logo { width: 160px; height: auto; margin-bottom: 1rem; display: block; margin-left: auto; margin-right: auto; }
  h1 { font-size: 1.5rem; margin-bottom: 0.25rem; }
  .subtitle { color: #666; margin-bottom: 1.5rem; }
  label { display: block; margin-top: 0.75rem; font-weight: 500; }
  input {
    width: 100%;
    padding: 0.5rem;
    margin-top: 0.25rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    box-sizing: border-box;
  }
  .error { color: #c00; margin-top: 0.5rem; font-size: 0.9rem; }
  .actions { margin-top: 1.25rem; display: flex; gap: 0.5rem; }
  button {
    padding: 0.5rem 1rem;
    border: 1px solid #333;
    border-radius: 4px;
    background: #333;
    color: #fff;
    cursor: pointer;
  }
  button:disabled { opacity: 0.7; cursor: not-allowed; }
  button.secondary { background: #fff; color: #333; }
</style>
