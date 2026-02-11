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
    bottom: 1rem;
    right: 1rem;
    max-width: 400px;
    padding: 0.75rem 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0,0,0,0.15);
    display: flex;
    align-items: flex-start;
    gap: 0.5rem;
    z-index: 9999;
    background: #fff;
    border: 1px solid #ddd;
  }
  .toast.error { border-color: #e57373; background: #ffebee; }
  .toast.success { border-color: #81c784; background: #e8f5e9; }
  .msg { flex: 1; font-size: 0.9rem; }
  .details { font-size: 0.8rem; margin-top: 0.25rem; }
  .details pre { margin: 0.25rem 0 0; padding: 0.5rem; background: rgba(0,0,0,0.05); border-radius: 4px; overflow: auto; max-height: 120px; }
  .dismiss { background: none; border: none; font-size: 1.2rem; cursor: pointer; padding: 0 0.25rem; line-height: 1; opacity: 0.7; }
  .dismiss:hover { opacity: 1; }
</style>
