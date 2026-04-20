<script>
  export let logs = [];
</script>

<div class="card">
  <h2>Video / FFmpeg Logs</h2>
  <div class="log-container">
    {#if logs && logs.length > 0}
      {#each logs as log}
        <div class="log-line">{log}</div>
      {/each}
    {:else}
      <p>No logs available.</p>
    {/if}
  </div>
</div>

<style>
  .card {
    background: #f4f4f4;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
  }
  .log-container {
    background: #222;
    color: #0f0;
    font-family: monospace;
    padding: 1rem;
    height: 150px;
    overflow-y: auto;
    border-radius: 4px;
  }
  .log-line {
    white-space: pre-wrap;
    word-break: break-all;
  }
  @media (prefers-color-scheme: dark) {
    .card { background: #333; }
  }
</style>
