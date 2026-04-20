import { writable } from 'svelte/store';

export const telemetryStore = writable({
  networkInterfaces: {},
  systemLoad: [],
  gps: { lat: 0, lon: 0, speed: 0 },
  ffmpegLogs: []
});

export const connectionStatus = writable('disconnected');

export function connectWebSocket() {
  const wsUrl = `ws://${window.location.hostname}:8080/ws`;

  const connect = () => {
    connectionStatus.set('connecting');
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      connectionStatus.set('connected');
      console.log('WebSocket connected');
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'telemetry') {
          telemetryStore.update(store => {
            return {
              ...store,
              networkInterfaces: data.network_interfaces || {},
              systemLoad: data.system_load || [],
              gps: data.gps || { lat: 0, lon: 0, speed: 0 },
              ffmpegLogs: data.ffmpeg_logs || []
            };
          });
        }
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err);
      }
    };

    ws.onclose = () => {
      connectionStatus.set('disconnected');
      console.log('WebSocket disconnected. Reconnecting in 3s...');
      setTimeout(connect, 3000);
    };

    ws.onerror = (err) => {
      console.error('WebSocket error:', err);
      ws.close();
    };
  };

  connect();
}
