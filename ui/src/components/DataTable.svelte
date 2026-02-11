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
  .table-wrap { overflow-x: auto; }
  .data-table { width: 100%; border-collapse: collapse; font-size: 0.9rem; }
  .data-table th, .data-table td { text-align: left; padding: 0.5rem 0.75rem; border-bottom: 1px solid #eee; }
  .data-table th { font-weight: 600; background: #f5f5f5; }
</style>
