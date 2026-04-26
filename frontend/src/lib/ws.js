import { writable } from "svelte/store";

export const telemetryStore = writable({
  networkInterfaces: {},
  systemUsage: { cpu: 0, ram: 0, swap: 0 },
  wifiMode: "Not found",
});

export const connectionStatus = writable("disconnected");

export function connectWebSocket() {
  const wsUrl = `ws://${window.location.hostname}:9090/ws`;

  const connect = () => {
    connectionStatus.set("connecting");
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      connectionStatus.set("connected");
      console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "telemetry") {
          telemetryStore.update((store) => {
            return {
              ...store,
              networkInterfaces: data.network_interfaces || {},
              systemUsage: data.system_usage || { cpu: 0, ram: 0, swap: 0 },
              wifiMode: data.wifi_mode || "Not found",
            };
          });
        }
      } catch (err) {
        console.error("Failed to parse WebSocket message:", err);
      }
    };

    ws.onclose = () => {
      connectionStatus.set("disconnected");
      console.log("WebSocket disconnected. Reconnecting in 3s...");
      setTimeout(connect, 3000);
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
      ws.close();
    };
  };

  connect();
}
