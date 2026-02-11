<script lang="ts">
  import { push } from 'svelte-spa-router'
</script>
<p>Not found. <a href="#/" on:click|preventDefault={() => push('/')}>Dashboard</a></p>
