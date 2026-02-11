<script lang="ts">
  export let title: string | undefined
</script>

<section class="card">
  {#if title}<h2 class="card-title">{title}</h2>{/if}
  <slot />
</section>

<style>
  .card {
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: var(--space-5);
    margin-bottom: var(--space-4);
    background: var(--bg-elevated);
    border-left: 3px solid var(--accent);
    transition: border-color var(--duration-fast), box-shadow var(--duration-fast);
  }
  .card:hover {
    border-color: var(--border-focus);
    box-shadow: var(--shadow-sm);
  }
  .card-title {
    font-size: 1rem;
    margin: 0 0 var(--space-3) 0;
    font-weight: 600;
    color: var(--text);
    letter-spacing: -0.01em;
  }
</style>
