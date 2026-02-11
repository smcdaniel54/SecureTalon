<script lang="ts">
  export let open = false
  export let title = ''
</script>

{#if open}
  <div class="overlay" role="dialog" aria-modal="true" aria-labelledby="modal-title">
    <div class="modal">
      {#if title}<h2 id="modal-title" class="modal-title">{title}</h2>{/if}
      <slot />
    </div>
  </div>
{/if}

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  .modal {
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 1.25rem;
    max-width: 90vw;
    max-height: 90vh;
    overflow: auto;
    box-shadow: var(--shadow);
    color: var(--text);
  }
  .modal-title { margin: 0 0 1rem 0; font-size: 1.1rem; color: var(--text); }
</style>
