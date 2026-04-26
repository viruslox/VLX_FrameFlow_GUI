<script>
  import { onMount, onDestroy } from "svelte";
  import AnsiToHtml from "ansi-to-html";
  import { telemetryStore } from "./ws.js";

  const ansiConvert = new AnsiToHtml({ escapeXML: true });

  let lastResponse = "";
  let errorMsg = "";

  let bondingStatusOutput = "";
  let mediamtxStatus = "unknown";
  let mediamtxInterval;
  let gpsStatus = "unknown";
  let gpsInterval;
  let clientStatus = "unknown";
  let clientInterval;
  let cameramanStatus = "off";
  let cameramanInterval;

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
          : (JSON.stringify(data.output || data, null, 2) || ""),
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

  async function checkMediaMTXStatus(silent = false) {
    try {
      if (!silent) {
        lastResponse = "Loading...";
        errorMsg = "";
      }
      const res = await fetch(`http://${window.location.hostname}:8080/api/mediamtx/status`, {
        method: "POST"
      });
      const data = await res.json();

      if (!res.ok) {
        mediamtxStatus = "error";
        if (!silent) {
           errorMsg = data.error || "Request failed";
           lastResponse = "";
        }
        return;
      }

      // Check output string for typical service statuses
      const output = typeof data.output === "string" ? data.output.toLowerCase() : "";
      if (output.includes("stopped") || output.includes("inactive") || output.includes("dead")) {
         mediamtxStatus = "stopped";
      } else if (output.includes("running") || output.includes("active") || output.includes("executed")) {
         mediamtxStatus = "running";
      } else {
         mediamtxStatus = "running"; // fallback if command succeeds but output is unknown format
      }

      if (!silent) {
         lastResponse = ansiConvert.toHtml(
          typeof data.output === "string"
            ? data.output
            : (JSON.stringify(data.output || data, null, 2) || "")
         );
      }
    } catch (err) {
      mediamtxStatus = "error";
      if (!silent) {
         errorMsg = err.message;
         lastResponse = "";
      }
    }
  }

  async function checkGPSStatus(silent = false) {
    try {
      if (!silent) {
        lastResponse = "Loading...";
        errorMsg = "";
      }
      const res = await fetch(`http://${window.location.hostname}:8080/api/gps/status`, {
        method: "POST"
      });
      const data = await res.json();

      if (!res.ok) {
        gpsStatus = "error";
        if (!silent) {
           errorMsg = data.error || "Request failed";
           lastResponse = "";
        }
        return;
      }

      // Check output string for typical service statuses
      const output = typeof data.output === "string" ? data.output.toLowerCase() : "";
      if (output.includes("stopped") || output.includes("inactive") || output.includes("dead")) {
         gpsStatus = "stopped";
      } else if (output.includes("running") || output.includes("active") || output.includes("executed")) {
         gpsStatus = "running";
      } else {
         gpsStatus = "running"; // fallback if command succeeds but output is unknown format
      }

      if (!silent) {
         lastResponse = ansiConvert.toHtml(
          typeof data.output === "string"
            ? data.output
            : (JSON.stringify(data.output || data, null, 2) || "")
         );
      }
    } catch (err) {
      gpsStatus = "error";
      if (!silent) {
         errorMsg = err.message;
         lastResponse = "";
      }
    }
  }

  async function handleCameramanService() {
    try {
      errorMsg = "";
      lastResponse = "Loading...";

      const resStatus = await fetch(`/api/cameraman/status`, {
        method: "POST"
      });
      const dataStatus = await resStatus.json();

      if (!resStatus.ok) {
        throw new Error(dataStatus.error || "Request failed");
      }

      const resList = await fetch(`/api/cameraman/list-dev`, {
        method: "POST"
      });
      const dataList = await resList.json();

      if (!resList.ok) {
        throw new Error(dataList.error || "Request failed");
      }

      const outputStatus = typeof dataStatus.output === "string" ? dataStatus.output : (JSON.stringify(dataStatus.output || dataStatus, null, 2) || "");
      const outputList = typeof dataList.output === "string" ? dataList.output : (JSON.stringify(dataList.output || dataList, null, 2) || "");

      lastResponse = ansiConvert.toHtml(outputStatus + "\n" + outputList);
    } catch (err) {
      errorMsg = err.message;
      lastResponse = "";
    }
  }

  async function checkCameramanStatus(silent = false) {
    try {
      if (!silent) {
        lastResponse = "Loading...";
        errorMsg = "";
      }
      const res = await fetch(`/api/cameraman/status`, {
        method: "POST"
      });
      const data = await res.json();

      if (!res.ok) {
        cameramanStatus = "strange";
        if (!silent) {
           errorMsg = data.error || "Request failed";
           lastResponse = "";
        }
        return;
      }

      const output = typeof data.output === "string" ? data.output : "";
      const lowerOutput = output.toLowerCase();

      if (lowerOutput.includes("no active cameraman services running.")) {
         cameramanStatus = "off";
      } else if (lowerOutput.includes("cameraman services:") || lowerOutput.includes("active: active (running)")) {
         cameramanStatus = "running";
      } else {
         cameramanStatus = "strange";
      }

      if (!silent) {
         lastResponse = ansiConvert.toHtml(
          typeof data.output === "string"
            ? data.output
            : (JSON.stringify(data.output || data, null, 2) || "")
         );
      }
    } catch (err) {
      cameramanStatus = "strange";
      if (!silent) {
         errorMsg = err.message;
         lastResponse = "";
      }
    }
  }

  async function checkClientStatus(silent = false) {
    try {
      if (!silent) {
        lastResponse = "Loading...";
        errorMsg = "";
      }
      const res = await fetch(`http://${window.location.hostname}:8080/api/frameflow/client/status`, {
        method: "POST"
      });
      const data = await res.json();

      if (!res.ok) {
        clientStatus = "error";
        if (!silent) {
           errorMsg = data.error || "Request failed";
           lastResponse = "";
        }
        return;
      }

      // Check output string for typical service statuses
      const output = typeof data.output === "string" ? data.output.toLowerCase() : "";
      if (output.includes("stopped") || output.includes("inactive") || output.includes("dead")) {
         clientStatus = "stopped";
      } else if (output.includes("running") || output.includes("active") || output.includes("executed")) {
         clientStatus = "running";
      } else {
         clientStatus = "running"; // fallback if command succeeds but output is unknown format
      }

      if (!silent) {
         lastResponse = ansiConvert.toHtml(
          typeof data.output === "string"
            ? data.output
            : (JSON.stringify(data.output || data, null, 2) || "")
         );
      }
    } catch (err) {
      clientStatus = "error";
      if (!silent) {
         errorMsg = err.message;
         lastResponse = "";
      }
    }
  }

  onMount(() => {
    // initial fetch for bonding status
    handleAction("/api/frameflow/bonding", "GET", null, true);
    checkMediaMTXStatus(true);
    mediamtxInterval = setInterval(() => checkMediaMTXStatus(true), 60000);
    checkGPSStatus(true);
    gpsInterval = setInterval(() => checkGPSStatus(true), 60000);
    checkClientStatus(true);
    clientInterval = setInterval(() => checkClientStatus(true), 60000);
    checkCameramanStatus(true);
    cameramanInterval = setInterval(() => checkCameramanStatus(true), 60000);
  });

  onDestroy(() => {
    if (mediamtxInterval) {
      clearInterval(mediamtxInterval);
    }
    if (gpsInterval) {
      clearInterval(gpsInterval);
    }
    if (clientInterval) {
      clearInterval(clientInterval);
    }
    if (cameramanInterval) {
      clearInterval(cameramanInterval);
    }
  });
</script>

<div class="card">
  <h2 style="text-align: center;">Bonding</h2>
  <div class="controls-grid three-cols">
    <!-- FrameFlow Client -->
    <div class="control-group">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h3 style="margin: 0;">FrameFlow Client</h3>
        <div class="indicator {clientStatus}"></div>
      </div>
      <div class="buttons">
        <button on:click={async () => { await handleAction("/api/frameflow/client/start"); await checkClientStatus(true); }}
          >Start</button
        >
        <button on:click={async () => { await handleAction("/api/frameflow/client/stop"); await checkClientStatus(true); }}
          >Stop</button
        >
        <button on:click={() => checkClientStatus(false)}
          >Status</button
        >
        <button on:click={async () => { await handleAction("/api/frameflow/client/reset"); await checkClientStatus(true); }}
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
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h3 style="margin: 0;">Access Point</h3>
        <span style="font-size: 0.85em; color: #666;">{$telemetryStore.wifiMode}</span>
      </div>
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
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h3 style="margin: 0;">Cameraman</h3>
        <div class="indicator {cameramanStatus}"></div>
      </div>
      <div class="form-group" style="margin-bottom: 0.5rem; margin-top: 0.5rem;">
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
        <button
          on:click={handleCameramanService} style="margin-left: auto;">Service</button
        >
      </div>
    </div>

    <!-- MediaMTX -->
    <div class="control-group">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h3 style="margin: 0;">MediaMTX</h3>
        <div class="indicator {mediamtxStatus}"></div>
      </div>
      <div class="buttons">
        <button on:click={async () => { await handleAction("/api/mediamtx/start"); await checkMediaMTXStatus(true); }}
          >Start</button
        >
        <button on:click={async () => { await handleAction("/api/mediamtx/stop"); await checkMediaMTXStatus(true); }}>Stop</button
        >
        <button on:click={() => checkMediaMTXStatus(false)}
          >Status</button
        >
      </div>
    </div>

    <!-- GPS Tracker -->
    <div class="control-group">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h3 style="margin: 0;">GPS Tracker</h3>
        <div class="indicator {gpsStatus}"></div>
      </div>
      <div class="buttons">
        <button on:click={async () => { await handleAction("/api/gps/start"); await checkGPSStatus(true); }}>Start</button>
        <button on:click={async () => { await handleAction("/api/gps/stop"); await checkGPSStatus(true); }}>Stop</button>
        <button on:click={() => checkGPSStatus(false)}>Status</button>
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

  .indicator {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background-color: gray;
  }
  .indicator.running {
    background-color: #0f0;
    box-shadow: 0 0 5px #0f0;
  }
  .indicator.stopped, .indicator.strange {
    background-color: #ff0;
    box-shadow: 0 0 5px #ff0;
  }
  .indicator.error {
    background-color: #f00;
    box-shadow: 0 0 5px #f00;
  }
  .indicator.off {
    background-color: gray;
    box-shadow: none;
  }
</style>
