<script lang="ts">
  import { push } from 'svelte-spa-router'
</script>
<div class="page">
  <p class="muted">Not found. <a href="#/" on:click|preventDefault={() => push('/')}>Dashboard</a></p>
</div>
<style>
  .muted { color: var(--text-muted); }
</style>
