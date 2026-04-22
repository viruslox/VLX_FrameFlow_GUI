<script>
  import { onMount } from "svelte";
  import AnsiToHtml from "ansi-to-html";

  const ansiConvert = new AnsiToHtml();

  let lastResponse = "";
  let errorMsg = "";

  let bondingStatusOutput = "";
  let bitrate = "2500";

  let camV = "0";
  let camA = "1";

  // Devices mapping
  const videoDevices = ["0", "1", "2", "3"];
  const audioDevices = ["0", "1", "2", "3"];

  $: deviceName = `V${camV}A${camA}`;

  async function handleAction(
    endpoint,
    method = "POST",
    body = null,
    isBondingStatus = false,
  ) {
    try {
      errorMsg = "";
      if (isBondingStatus) {
        bondingStatusOutput = "Loading...";
      } else {
        lastResponse = "Loading...";
      }

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

      // Convert ansi colors to html
      const formatted = ansiConvert.toHtml(
        typeof data.output === "string"
          ? data.output
          : JSON.stringify(data.output || data, null, 2),
      );

      if (isBondingStatus) {
        bondingStatusOutput = formatted;
      } else {
        lastResponse = formatted;
      }
    } catch (err) {
      if (isBondingStatus) {
        bondingStatusOutput = err.message;
      } else {
        errorMsg = err.message;
        lastResponse = "";
      }
    }
  }

  function handleBitrateUpdate() {
    handleAction("/api/mediamtx/config", "POST", { bitrate: bitrate });
  }

  onMount(() => {
    // initial fetch for bonding status
    handleAction("/api/frameflow/bonding", "GET", null, true);
  });
</script>

<div class="card">
  <h2 style="text-align: center;">Bonding</h2>
  <div class="controls-grid three-cols">
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

    <!-- Bonding Status -->
    <div class="control-group console-box">
      <div
        style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;"
      >
        <h3 style="margin: 0;">Status</h3>
        <button
          on:click={() =>
            handleAction("/api/frameflow/bonding", "GET", null, true)}
          >Refresh</button
        >
      </div>
      <div class="mini-response">
        {@html bondingStatusOutput || "No data"}
      </div>
    </div>

    <!-- Access Point (FrameFlow AP) -->
    <div class="control-group">
      <h3>Access Point</h3>
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
  </div>
</div>

<div class="card">
  <h2 style="text-align: center;">FrameFlow Services</h2>
  <div class="controls-grid three-cols">
    <!-- Cameraman -->
    <div class="control-group">
      <h3>Cameraman</h3>
      <div class="form-group" style="margin-bottom: 0.5rem;">
        <label for="camV">Video (Vx):</label>
        <select id="camV" bind:value={camV}>
          {#each videoDevices as v}
            <option value={v}>V{v}</option>
          {/each}
        </select>
        <label for="camA">Audio (Vy):</label>
        <select id="camA" bind:value={camA}>
          {#each audioDevices as a}
            <option value={a}>A{a}</option>
          {/each}
        </select>
      </div>
      <p style="margin: 0.2rem 0; font-size: 0.9em;">Device: {deviceName}</p>
      <div class="buttons">
        <button
          on:click={() =>
            handleAction("/api/cameraman/start", "POST", {
              device: deviceName,
            })}>Start</button
        >
        <button
          on:click={() =>
            handleAction("/api/cameraman/stop", "POST", { device: deviceName })}
          >Stop</button
        >
        <button
          on:click={() =>
            handleAction("/api/cameraman/status", "POST", {
              device: deviceName,
            })}>Status</button
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
  </div>
</div>

<div class="card response-area">
  <h4>Response:</h4>
  {#if errorMsg}
    <pre class="error">{errorMsg}</pre>
  {:else if lastResponse}
    <pre>{@html lastResponse}</pre>
  {:else}
    <p class="muted">No recent actions.</p>
  {/if}
</div>

<style>
  .card {
    background: #f4f4f4;
    padding: 1rem;
    border-radius: 8px;
  }
  .controls-grid {
    display: grid;
    gap: 1rem;
    margin-bottom: 0.5rem;
  }
  .three-cols {
    grid-template-columns: repeat(3, 1fr);
  }

  @media (max-width: 768px) {
    .three-cols {
      grid-template-columns: 1fr;
    }
  }

  .control-group {
    background: #fff;
    padding: 1rem;
    border-radius: 6px;
    border: 1px solid #ddd;
    display: flex;
    flex-direction: column;
  }
  .console-box {
    background: #222;
    color: #0f0;
    border-color: #444;
  }
  .console-box h3 {
    color: #fff;
  }

  .control-group h3 {
    margin-top: 0;
    font-size: 1.1rem;
  }
  .buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: auto;
  }
  .form-group {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .form-group input,
  .form-group select {
    padding: 0.3rem;
    border: 1px solid #ccc;
    border-radius: 4px;
    width: 80px;
  }
  .form-group select {
    width: 60px;
  }
  button {
    font-size: 0.9rem;
    padding: 0.4em 0.8em;
    cursor: pointer;
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

  .mini-response {
    background: transparent;
    color: #0f0;
    font-family: monospace;
    font-size: 0.85rem;
    white-space: pre-wrap;
    word-wrap: break-word;
    overflow-y: auto;
    max-height: 150px;
  }

  pre {
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
    font-family: monospace;
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
    .console-box {
      background: #222;
    }
    .form-group input,
    .form-group select {
      background: #555;
      color: #fff;
      border-color: #666;
    }
    button {
      background-color: #555;
      color: #fff;
      border: 1px solid #666;
    }
    button:hover {
      background-color: #666;
      border-color: #888;
    }
  }
</style>
