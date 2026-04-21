<script>
  let lastResponse = "";
  let errorMsg = "";

  let bitrate = "2500";
  let bondingInterface = "bond0";

  async function handleAction(endpoint, method = "POST", body = null) {
    try {
      errorMsg = "";
      lastResponse = "Loading...";

      const options = {
        method,
        headers: {},
      };

      if (body) {
        options.headers["Content-Type"] = "application/json";
        options.body = JSON.stringify(body);
      }

      const res = await fetch(
        `http://${window.location.hostname}:8080${endpoint}`,
        options,
      );
      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || "Request failed");
      }

      lastResponse = JSON.stringify(data.output || data, null, 2);
    } catch (err) {
      errorMsg = err.message;
      lastResponse = "";
    }
  }

  function handleBondingUpdate() {
    // In a real app we'd probably POST this to a specific configuration endpoint
    handleAction("/api/frameflow/bonding", "POST", {
      interface: bondingInterface,
    });
  }

  function handleBitrateUpdate() {
    // In a real app we'd probably POST this to a specific configuration endpoint
    handleAction("/api/mediamtx/config", "POST", { bitrate: bitrate });
  }
</script>

<div class="card">
  <h2>Control Panel</h2>

  <div class="controls-grid">
    <!-- FrameFlow Client -->
    <div class="control-group">
      <h3>FrameFlow Client</h3>
      <div class="buttons">
        <button on:click={() => handleAction("/api/frameflow/client/start")}
          >Start</button
        >
        <button on:click={() => handleAction("/api/frameflow/client/stop")}
          >Stop</button
        >
        <button on:click={() => handleAction("/api/frameflow/client/status")}
          >Status</button
        >
        <button on:click={() => handleAction("/api/frameflow/client/reset")}
          >Reset</button
        >
      </div>
    </div>

    <!-- FrameFlow AP -->
    <div class="control-group">
      <h3>FrameFlow AP</h3>
      <div class="buttons">
        <button on:click={() => handleAction("/api/frameflow/ap/start")}
          >Start</button
        >
        <button on:click={() => handleAction("/api/frameflow/ap/stop")}
          >Stop</button
        >
        <button on:click={() => handleAction("/api/frameflow/ap/status")}
          >Status</button
        >
      </div>
    </div>

    <!-- FrameFlow Bonding -->
    <div class="control-group">
      <h3>Bonding</h3>
      <div class="form-group">
        <label for="bond-iface">Interface:</label>
        <input type="text" id="bond-iface" bind:value={bondingInterface} />
        <button on:click={handleBondingUpdate}>Update</button>
      </div>
      <div class="buttons" style="margin-top: 0.5rem;">
        <button on:click={() => handleAction("/api/frameflow/bonding", "GET")}
          >Get Status</button
        >
      </div>
    </div>

    <!-- MediaMTX -->
    <div class="control-group">
      <h3>MediaMTX</h3>
      <div class="form-group">
        <label for="bitrate">Bitrate (kbps):</label>
        <input type="number" id="bitrate" bind:value={bitrate} />
        <button on:click={handleBitrateUpdate}>Update</button>
      </div>
      <div class="buttons" style="margin-top: 0.5rem;">
        <button on:click={() => handleAction("/api/mediamtx/start")}
          >Start</button
        >
        <button on:click={() => handleAction("/api/mediamtx/stop")}>Stop</button
        >
        <button on:click={() => handleAction("/api/mediamtx/status")}
          >Status</button
        >
      </div>
    </div>

    <!-- GPS Tracker -->
    <div class="control-group">
      <h3>GPS Tracker</h3>
      <div class="buttons">
        <button on:click={() => handleAction("/api/gps/start")}>Start</button>
        <button on:click={() => handleAction("/api/gps/stop")}>Stop</button>
        <button on:click={() => handleAction("/api/gps/status")}>Status</button>
      </div>
    </div>

    <!-- Cameraman -->
    <div class="control-group">
      <h3>Cameraman (V0A1)</h3>
      <div class="buttons">
        <button
          on:click={() =>
            handleAction("/api/cameraman/start", "POST", { device: "V0A1" })}
          >Start</button
        >
        <button
          on:click={() =>
            handleAction("/api/cameraman/stop", "POST", { device: "V0A1" })}
          >Stop</button
        >
        <button
          on:click={() =>
            handleAction("/api/cameraman/status", "POST", { device: "V0A1" })}
          >Status</button
        >
      </div>
    </div>
  </div>

  <div class="response-area">
    <h4>Response:</h4>
    {#if errorMsg}
      <pre class="error">{errorMsg}</pre>
    {:else if lastResponse}
      <pre>{lastResponse}</pre>
    {:else}
      <p class="muted">No recent actions.</p>
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
  .controls-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1rem;
    margin-bottom: 1.5rem;
  }
  .control-group {
    background: #fff;
    padding: 1rem;
    border-radius: 6px;
    border: 1px solid #ddd;
  }
  .control-group h3 {
    margin-top: 0;
    font-size: 1.1rem;
  }
  .buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  .form-group {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .form-group input {
    padding: 0.3rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    width: 80px;
  }
  button {
    font-size: 0.9rem;
    padding: 0.4em 0.8em;
  }
  .response-area {
    background: #222;
    color: #0f0;
    padding: 1rem;
    border-radius: 6px;
    min-height: 100px;
    max-height: 300px;
    overflow-y: auto;
  }
  .response-area h4 {
    margin-top: 0;
    color: #fff;
  }
  pre {
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
  .error {
    color: #ff4444;
  }
  .muted {
    color: #888;
    font-style: italic;
  }

  @media (prefers-color-scheme: dark) {
    .card {
      background: #333;
    }
    .control-group {
      background: #444;
      border-color: #555;
    }
    .form-group input {
      background: #555;
      color: #fff;
      border-color: #666;
    }
    button {
      background-color: #555;
      color: #fff;
    }
    button:hover {
      background-color: #666;
      border-color: #888;
    }
  }
</style>
