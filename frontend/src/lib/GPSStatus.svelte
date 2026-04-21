<script>
  export let gps = { lat: 0, lon: 0, speed: 0 };
</script>

<div class="card">
  <h2>GPS Status</h2>
  <p>Latitude: {gps.lat}</p>
  <p>Longitude: {gps.lon}</p>
  <p>Speed: {gps.speed} m/s</p>
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
