<script lang="ts">
  export let title = ''
</script>

<div class="form-panel">
  {#if title}<h3 class="form-panel-title">{title}</h3>{/if}
  <slot />
</div>

<style>
  .form-panel { margin-bottom: 1rem; }
  .form-panel-title { font-size: 0.95rem; margin: 0 0 var(--space-2) 0; font-weight: 600; color: var(--text); }
  /* Styles for slotted form fields (parents put label/input/select/textarea inside the slot) */
  .form-panel :global(label) { display: block; margin-top: 0.5rem; font-size: 0.9rem; }
  .form-panel :global(input),
  .form-panel :global(select),
  .form-panel :global(textarea) {
    width: 100%;
    margin-top: 0.25rem;
    box-sizing: border-box;
  }
</style>
