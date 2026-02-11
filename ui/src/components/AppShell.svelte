<script lang="ts">
  import { onMount } from 'svelte'
  import { push } from 'svelte-spa-router'
  import { authStore } from '../app/store'
  import { getStoredAuth } from '../lib/api'
  import Toast from './Toast.svelte'

  export let title = 'SecureTalon Admin'

  onMount(() => {
    if (!$authStore) {
      const stored = getStoredAuth()
      if (stored) authStore.setAuth(stored)
    }
  })

  $: path = (typeof window !== 'undefined' && window.location.hash)
    ? (window.location.hash.slice(2).split('?')[0] || '/')
    : ''

  function isActive(href: string): boolean {
    const p = href === '#/' ? '/' : href.replace('#', '')
    return path === p || (p !== '/' && path.startsWith(p + '/'))
  }
</script>

<svelte:head>
  <title>{title}</title>
</svelte:head>

<div class="shell">
  {#if $authStore}
    <aside class="nav">
      <div class="nav-brand">
        <img src="/securetalon-logo.png" alt="SecureTalon" class="nav-logo" />
        <span>SecureTalon</span>
      </div>
      <nav class="nav-links">
        <a href="#/" class:active={isActive('#/')} on:click|preventDefault={() => push('/')}>Dashboard</a>
        <a href="#/sessions" class:active={path === '/sessions' || path.startsWith('/sessions/')} on:click|preventDefault={() => push('/sessions')}>Sessions</a>
        <a href="#/policies" class:active={path === '/policies'} on:click|preventDefault={() => push('/policies')}>Policies</a>
        <a href="#/skills" class:active={path === '/skills'} on:click|preventDefault={() => push('/skills')}>Skills</a>
        <a href="#/audit" class:active={path === '/audit'} on:click|preventDefault={() => push('/audit')}>Audit</a>
        <a href="#/replay" class:active={path === '/replay'} on:click|preventDefault={() => push('/replay')}>Replay</a>
      </nav>
      <div class="nav-footer">
        <button class="disconnect" on:click={() => { authStore.clear(); push('/login') }}>Disconnect</button>
      </div>
    </aside>
    <main class="main">
      <slot />
    </main>
  {:else}
    <main class="main full">
      <slot />
    </main>
  {/if}
  <Toast />
</div>

<style>
  .shell { display: flex; min-height: 100vh; }
  .nav {
    width: 240px;
    background: var(--bg-elevated);
    border-right: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    padding: var(--space-4) 0;
    border-left: 3px solid var(--accent);
  }
  .nav-brand {
    font-weight: 600;
    padding: 0 var(--space-4) var(--space-4);
    font-size: 0.9375rem;
    display: flex;
    align-items: center;
    gap: var(--space-2);
    color: var(--accent);
    letter-spacing: -0.02em;
  }
  .nav-logo { height: 30px; width: auto; display: block; }
  .nav-links {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    padding: 0 var(--space-2);
  }
  .nav-links a {
    padding: var(--space-2) var(--space-3);
    color: var(--text-muted);
    text-decoration: none;
    font-size: 0.875rem;
    border-radius: var(--radius-sm);
    transition: background var(--duration-fast), color var(--duration-fast);
  }
  .nav-links a:hover {
    color: var(--text);
    background: var(--hover-overlay);
    text-decoration: none;
  }
  .nav-links a.active {
    color: var(--accent);
    background: var(--accent-subtle);
    font-weight: 500;
    box-shadow: inset 3px 0 0 var(--accent);
  }
  .nav-links a.active:hover {
    color: var(--accent-hover);
    background: var(--accent-subtle);
  }
  .nav-footer {
    padding: var(--space-3) var(--space-4) 0;
    border-top: 1px solid var(--border);
    margin-top: var(--space-2);
    padding-top: var(--space-3);
  }
  .disconnect {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    background: transparent;
    border: 1px solid var(--border);
    color: var(--text-muted);
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 0.8125rem;
    transition: border-color var(--duration-fast), color var(--duration-fast);
  }
  .disconnect:hover {
    color: var(--text);
    border-color: var(--text-subtle);
  }
  .main {
    flex: 1;
    padding: var(--space-6) var(--space-8);
    overflow: auto;
    background: var(--bg);
  }
  .main.full {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
