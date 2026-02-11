<script lang="ts">
  import { toastStore } from '../app/store'
</script>

{#each $toastStore as toast}
  <div class="toast" class:error={toast.type === 'error'} class:success={toast.type === 'success'} role="alert">
    <span class="msg">{toast.message}</span>
    {#if toast.details && Object.keys(toast.details).length > 0}
      <details class="details">
        <summary>Details</summary>
        <pre>{JSON.stringify(toast.details, null, 2)}</pre>
      </details>
    {/if}
    <button type="button" class="dismiss" on:click={() => toastStore.remove(toast.id)} aria-label="Dismiss">×</button>
  </div>
{/each}

<style>
  .toast {
    position: fixed;
    bottom: var(--space-4);
    right: var(--space-4);
    max-width: 400px;
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius);
    box-shadow: var(--shadow-lg);
    display: flex;
    align-items: flex-start;
    gap: var(--space-2);
    z-index: 9999;
    background: var(--bg-overlay);
    border: 1px solid var(--border);
    color: var(--text);
  }
  .toast.error { border-color: var(--error); background: var(--error-subtle); }
  .toast.success { border-color: var(--success); background: var(--success-subtle); }
  .msg { flex: 1; font-size: 0.875rem; }
  .details { font-size: 0.8rem; margin-top: var(--space-1); }
  .details pre { margin: var(--space-1) 0 0; padding: var(--space-2); background: var(--bg-input); border-radius: var(--radius-sm); overflow: auto; max-height: 120px; }
  .dismiss { background: none; border: none; font-size: 1.2rem; cursor: pointer; padding: 0 var(--space-1); line-height: 1; opacity: 0.7; color: var(--text); }
  .dismiss:hover { opacity: 1; }
</style>
