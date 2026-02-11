<script lang="ts">
  export let columns: { key: string; label: string }[] = []
  export let rows: Record<string, unknown>[] = []
</script>

<div class="table-wrap">
  <table class="data-table">
    <thead>
      <tr>
        {#each columns as col}
          <th>{col.label}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each rows as row}
        <tr>
          {#each columns as col}
            <td>{String(row[col.key] ?? '—')}</td>
          {/each}
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .table-wrap { overflow-x: auto; border-radius: var(--radius); border: 1px solid var(--border); }
  .data-table { width: 100%; border-collapse: collapse; font-size: 0.875rem; }
  .data-table th, .data-table td { text-align: left; padding: 0.5rem 0.75rem; border-bottom: 1px solid var(--border); }
  .data-table th { font-weight: 600; background: var(--bg-elevated); color: var(--text); }
  .data-table td { color: var(--text); }
  .data-table tbody tr { transition: background var(--duration-fast); }
  .data-table tbody tr:hover { background: var(--hover-overlay); }
</style>
