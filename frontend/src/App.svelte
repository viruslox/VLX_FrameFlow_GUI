<script>
  import { onMount } from "svelte";
  import {
    telemetryStore,
    connectionStatus,
    connectWebSocket,
  } from "./lib/ws.js";
  import SystemStatus from "./lib/SystemStatus.svelte";
  import NetworkStatus from "./lib/NetworkStatus.svelte";
  import ControlPanel from "./lib/ControlPanel.svelte";

  onMount(() => {
    connectWebSocket();
  });
</script>

<main>
  <h1>VLX FrameFlow Dashboard</h1>
  <p class="status">Connection Status: {$connectionStatus}</p>

  <div class="grid">
    <div class="col">
      <SystemStatus usage={$telemetryStore.systemUsage} />
    </div>
    <div class="col">
      <NetworkStatus interfaces={$telemetryStore.networkInterfaces} />
    </div>
  </div>

  <ControlPanel />
</main>

<style>
  main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    font-family: sans-serif;
  }
  h1 {
    text-align: center;
    color: #333;
  }
  .status {
    text-align: center;
    font-weight: bold;
    margin-bottom: 2rem;
  }
  .grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  @media (max-width: 768px) {
    .grid {
      grid-template-columns: 1fr;
    }
  }
  @media (prefers-color-scheme: dark) {
    h1 {
      color: #fff;
    }
  }
</style>
