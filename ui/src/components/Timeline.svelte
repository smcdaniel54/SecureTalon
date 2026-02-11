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
    padding: 0.75rem 1rem;
    margin: 0.5rem 0;
    border: 1px solid var(--border);
    border-left: 4px solid var(--accent);
    background: var(--bg-elevated);
    border-radius: var(--radius-sm);
  }
  .ts { font-size: 0.8rem; color: var(--text-muted); margin-right: 0.5rem; }
  .type { font-weight: 500; margin-right: 0.5rem; color: var(--text); }
  .status { font-size: 0.85rem; }
  .status.ok { color: var(--success); }
  .status.deny { color: var(--error); }
  .status.err { color: var(--warning); }
  .data { font-size: 0.8rem; margin: 0.5rem 0 0; padding: 0.5rem; background: var(--bg-input); border: 1px solid var(--border); border-radius: var(--radius-sm); overflow: auto; color: var(--text); }
</style>
