<script lang="ts">
  export interface TimelineItem {
    id?: string
    ts: string
    type: string
    label?: string
    data?: Record<string, unknown>
    status?: string
  }
  export let items: TimelineItem[] = []
</script>

<ol class="timeline">
  {#each items as item}
    <li class="timeline-item">
      <span class="ts">{typeof item.ts === 'string' ? new Date(item.ts).toLocaleString() : item.ts}</span>
      <span class="type">{item.type}</span>
      {#if item.label}<span class="label">{item.label}</span>{/if}
      {#if item.status}<span class="status" class:ok={item.status === 'ok'} class:deny={item.status === 'denied'} class:err={item.status === 'error'}>{item.status}</span>{/if}
      {#if item.data && Object.keys(item.data).length > 0}
        <pre class="data">{JSON.stringify(item.data, null, 2)}</pre>
      {/if}
    </li>
  {/each}
</ol>

<style>
  .timeline { list-style: none; padding: 0; margin: 0; }
  .timeline-item {
    padding: 0.75rem;
    margin: 0.5rem 0;
    border-left: 4px solid #2196f3;
    background: #fafafa;
    border-radius: 4px;
  }
  .ts { font-size: 0.8rem; color: #666; margin-right: 0.5rem; }
  .type { font-weight: 500; margin-right: 0.5rem; }
  .status { font-size: 0.85rem; }
  .status.ok { color: #2e7d32; }
  .status.deny { color: #c62828; }
  .status.err { color: #e65100; }
  .data { font-size: 0.8rem; margin: 0.5rem 0 0; padding: 0.5rem; background: #fff; border: 1px solid #eee; overflow: auto; }
</style>
