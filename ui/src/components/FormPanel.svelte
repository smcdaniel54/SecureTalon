<script lang="ts">
  export let title = ''
</script>

<div class="form-panel">
  {#if title}<h3 class="form-panel-title">{title}</h3>{/if}
  <slot />
</div>

<style>
  .form-panel { margin-bottom: 1rem; }
  .form-panel-title { font-size: 0.95rem; margin: 0 0 0.5rem 0; font-weight: 600; }
  .form-panel label { display: block; margin-top: 0.5rem; font-size: 0.9rem; }
  .form-panel input, .form-panel select, .form-panel textarea {
    width: 100%;
    padding: 0.4rem 0.5rem;
    margin-top: 0.2rem;
    box-sizing: border-box;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
</style>
