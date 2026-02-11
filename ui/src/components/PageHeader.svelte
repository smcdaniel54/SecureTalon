<script lang="ts">
  export let title: string
  export let subtitle: string | undefined
</script>

<header class="page-header">
  <h1>{title}</h1>
  {#if subtitle}<p class="subtitle">{subtitle}</p>{/if}
</header>

<style>
  .page-header {
    margin-bottom: var(--space-6);
    padding-bottom: var(--space-4);
    border-bottom: 3px solid var(--accent);
  }
  .page-header h1 {
    font-size: 1.5rem;
    margin: 0 0 var(--space-1) 0;
    letter-spacing: -0.02em;
  }
  .subtitle {
    color: var(--text-muted);
    font-size: 0.9rem;
    margin: 0;
  }
</style>
