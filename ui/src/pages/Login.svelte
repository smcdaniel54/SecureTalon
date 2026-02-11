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
  <div class="login-bg" aria-hidden="true"></div>
  <div class="login-card">
    <div class="login-brand">
      <img src="/securetalon-logo.png" alt="SecureTalon" class="login-logo" />
      <h1>SecureTalon</h1>
      <p class="tagline">Security-first agent platform</p>
    </div>
    <p class="subtitle">Connect to your API to continue</p>
    <form on:submit|preventDefault={handleConnect}>
      <label for="apiBase">API Base URL</label>
      <input id="apiBase" type="url" bind:value={apiBase} placeholder="http://localhost:8080" />
      <label for="token">Admin Token</label>
      <input id="token" type="password" bind:value={token} placeholder="e.g. demo" />
      {#if error}
        <p class="error">{error}</p>
      {/if}
      <div class="actions">
        <button type="submit" class="primary" disabled={loading}>{loading ? 'Connecting…' : 'Connect'}</button>
        <button type="button" class="secondary" on:click={handleForget}>Forget</button>
      </div>
    </form>
  </div>
</div>

<style>
  .login {
    width: 100%;
    min-height: 100vh;
    padding: var(--space-8);
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
  }
  .login-bg {
    position: absolute;
    inset: 0;
    background:
      radial-gradient(ellipse 80% 50% at 50% -20%, rgba(13, 148, 136, 0.12), transparent),
      radial-gradient(ellipse 60% 40% at 100% 100%, rgba(13, 148, 136, 0.06), transparent);
    pointer-events: none;
  }
  .login-card {
    position: relative;
    max-width: 420px;
    width: 100%;
    padding: var(--space-10);
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-lg);
    text-align: center;
    border-top: 3px solid var(--accent);
  }
  .login-brand {
    margin-bottom: var(--space-6);
  }
  .login-logo {
    width: 180px;
    height: auto;
    margin: 0 auto var(--space-4);
    display: block;
  }
  .login-card h1 {
    font-size: 1.625rem;
    margin-bottom: var(--space-1);
    letter-spacing: -0.03em;
  }
  .tagline {
    font-size: 0.875rem;
    color: var(--text-muted);
    margin: 0;
  }
  .subtitle {
    color: var(--text-muted);
    margin-bottom: var(--space-6);
    font-size: 0.9rem;
  }
  .login-card label {
    text-align: left;
    margin-top: var(--space-4);
  }
  .login-card label:first-of-type {
    margin-top: 0;
  }
  .login-card input {
    width: 100%;
    margin-top: var(--space-1);
    padding: var(--space-3);
  }
  .error {
    color: var(--error);
    margin-top: var(--space-3);
    font-size: 0.875rem;
  }
  .actions {
    margin-top: var(--space-6);
    display: flex;
    gap: var(--space-3);
    justify-content: center;
  }
  .actions .primary {
    padding: var(--space-2) var(--space-6);
  }
</style>
