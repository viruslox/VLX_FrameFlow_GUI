<script>
  export let interfaces = {};

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }
</script>

<div class="card">
  <h2>Network Status</h2>
  {#if Object.keys(interfaces).length > 0}
    <table>
      <thead>
        <tr>
          <th>Interface</th>
          <th>Rx Bytes</th>
          <th>Tx Bytes</th>
        </tr>
      </thead>
      <tbody>
        {#each Object.entries(interfaces) as [iface, stats]}
          <tr>
            <td>{iface}</td>
            <td>{formatBytes(stats.rx_bytes)}</td>
            <td>{formatBytes(stats.tx_bytes)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Loading...</p>
  {/if}
</div>

<style>
  .card {
    background: #f4f4f4;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
  }
  table {
    width: 100%;
    border-collapse: collapse;
  }
  th, td {
    text-align: left;
    padding: 0.5rem;
    border-bottom: 1px solid #ddd;
  }
  @media (prefers-color-scheme: dark) {
    .card { background: #333; }
    th, td { border-bottom: 1px solid #555; }
  }
</style>
