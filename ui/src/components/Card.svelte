<script lang="ts">
  export let title: string | undefined
</script>

<section class="card">
  {#if title}<h2 class="card-title">{title}</h2>{/if}
  <slot />
</section>

<style>
  .card {
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    padding: 1rem;
    margin-bottom: 1rem;
    background: #fafafa;
  }
  .card-title { font-size: 1rem; margin: 0 0 0.75rem 0; font-weight: 600; }
</style>
