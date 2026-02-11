<script lang="ts">
  export let title: string
  export let subtitle: string | undefined
</script>

<header class="page-header">
  <h1>{title}</h1>
  {#if subtitle}<p class="subtitle">{subtitle}</p>{/if}
</header>

<style>
  .page-header { margin-bottom: 1rem; }
  .page-header h1 { font-size: 1.5rem; margin: 0 0 0.25rem 0; }
  .subtitle { color: #666; font-size: 0.9rem; margin: 0; }
</style>
