<script lang="ts">
  import { push } from 'svelte-spa-router'
  import { authStore } from '../app/store'
  import Toast from './Toast.svelte'

  export let title = 'SecureTalon Admin'
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
        <a href="#/" on:click|preventDefault={() => push('/')}>Dashboard</a>
        <a href="#/sessions" on:click|preventDefault={() => push('/sessions')}>Sessions</a>
        <a href="#/policies" on:click|preventDefault={() => push('/policies')}>Policies</a>
        <a href="#/skills" on:click|preventDefault={() => push('/skills')}>Skills</a>
        <a href="#/audit" on:click|preventDefault={() => push('/audit')}>Audit</a>
        <a href="#/replay" on:click|preventDefault={() => push('/replay')}>Replay</a>
      </nav>
      <button class="disconnect" on:click={() => { authStore.clear(); push('/login') }}>Disconnect</button>
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
    width: 200px;
    background: #1a1a1a;
    color: #e0e0e0;
    display: flex;
    flex-direction: column;
    padding: 1rem 0;
  }
  .nav-brand {
    font-weight: 600;
    padding: 0 1rem 1rem;
    font-size: 1rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .nav-logo { height: 28px; width: auto; display: block; }
  .nav-links { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
  .nav-links a {
    padding: 0.5rem 1rem;
    color: #b0b0b0;
    text-decoration: none;
    font-size: 0.9rem;
  }
  .nav-links a:hover { color: #fff; background: rgba(255,255,255,0.05); }
  .disconnect {
    margin: 0.5rem 1rem;
    padding: 0.4rem;
    background: transparent;
    border: 1px solid #555;
    color: #b0b0b0;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.85rem;
  }
  .disconnect:hover { color: #fff; border-color: #888; }
  .main { flex: 1; padding: 1rem 1.5rem; overflow: auto; }
  .main.full { width: 100%; }
</style>
