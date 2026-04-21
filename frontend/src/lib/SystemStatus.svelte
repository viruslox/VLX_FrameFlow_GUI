<script>
  export let usage = { cpu: 0, ram: 0, swap: 0 };
</script>

<div class="card">
  <h2>System Usage</h2>
  <p>CPU: {usage.cpu.toFixed(1)}%</p>
  <p>RAM: {usage.ram.toFixed(1)}%</p>
  <p>Swap: {usage.swap.toFixed(1)}%</p>
</div>

<style>
  .card {
    background: #f4f4f4;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1rem;
  }
  @media (prefers-color-scheme: dark) {
    .card {
      background: #333;
    }
  }
</style>
