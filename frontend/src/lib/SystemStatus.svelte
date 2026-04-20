<script>
  export let load = [];
</script>

<div class="card">
  <h2>System Load</h2>
  {#if load && load.length === 3}
    <p>1 min: {load[0]}</p>
    <p>5 min: {load[1]}</p>
    <p>15 min: {load[2]}</p>
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
  @media (prefers-color-scheme: dark) {
    .card { background: #333; }
  }
</style>
