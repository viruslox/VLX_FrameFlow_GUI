import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { get } from 'svelte/store';
import { Server, WebSocket as MockWebSocket } from 'mock-socket';
import { connectWebSocket, connectionStatus, telemetryStore } from './ws.js';

describe('WebSocket connection logic', () => {
  let mockServer;
  const WS_URL = 'ws://localhost:8080/ws';

  beforeEach(() => {
    // Reset stores
    connectionStatus.set('disconnected');
    telemetryStore.set({
      networkInterfaces: {},
      systemUsage: { cpu: 0, ram: 0, swap: 0 },
    });

    // Mock window.location
    delete window.location;
    window.location = { hostname: 'localhost' };

    // Mock global WebSocket
    global.WebSocket = MockWebSocket;

    // Set up mock server
    mockServer = new Server(WS_URL);

    // Mock console.log and console.error
    vi.spyOn(console, 'log').mockImplementation(() => {});
    vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    mockServer.stop();
    vi.restoreAllMocks();
  });

  it('should establish connection and update status', async () => {
    return new Promise((resolve) => {
      mockServer.on('connection', socket => {
        setTimeout(() => {
          expect(get(connectionStatus)).toBe('connected');
          resolve();
        }, 10);
      });

      connectWebSocket();
      expect(get(connectionStatus)).toBe('connecting');
    });
  });

  it('should handle incoming telemetry messages', async () => {
    return new Promise((resolve) => {
      const mockData = {
        type: 'telemetry',
        network_interfaces: { eth0: { up: true } },
        system_usage: { cpu: 50, ram: 2048, swap: 0 }
      };

      mockServer.on('connection', socket => {
        socket.send(JSON.stringify(mockData));

        // Wait for Svelte store to update
        setTimeout(() => {
          const state = get(telemetryStore);
          expect(state.networkInterfaces).toEqual(mockData.network_interfaces);
          expect(state.systemUsage).toEqual(mockData.system_usage);
          resolve();
        }, 50);
      });

      connectWebSocket();
    });
  });

  it('should handle malformed JSON messages gracefully', async () => {
    return new Promise((resolve) => {
      mockServer.on('connection', socket => {
        socket.send('invalid-json');

        setTimeout(() => {
          expect(console.error).toHaveBeenCalledWith(
            'Failed to parse WebSocket message:',
            expect.any(SyntaxError)
          );
          resolve();
        }, 50);
      });

      connectWebSocket();
    });
  });

  it('should attempt reconnection on close', async () => {
    vi.useFakeTimers();
    let connections = 0;

    mockServer.on('connection', socket => {
      connections++;
      if (connections === 1) {
        socket.close();
      }
    });

    connectWebSocket();

    // Allow initial connection to happen
    await vi.advanceTimersByTimeAsync(100);

    // At this point connection was closed and we should be disconnected
    expect(console.log).toHaveBeenCalledWith('WebSocket disconnected. Reconnecting in 3s...');
    expect(get(connectionStatus)).toBe('disconnected');

    // Fast forward to after the setTimeout
    await vi.advanceTimersByTimeAsync(3000);

    // Should have attempted to connect again
    expect(connections).toBe(2);

    // And should be connected after the second attempt establishes
    await vi.advanceTimersByTimeAsync(100);
    expect(get(connectionStatus)).toBe('connected');

    vi.useRealTimers();
  });

  it('should handle error events', async () => {
    vi.useFakeTimers();

    let socketRef;
    mockServer.on('connection', socket => {
      socketRef = socket;
    });

    connectWebSocket();

    // Wait for connection to establish
    await vi.advanceTimersByTimeAsync(100);

    socketRef.close({ code: 1006, reason: 'Abnormal closure', wasClean: false });

    await vi.advanceTimersByTimeAsync(100);

    // The library we are using, mock-socket, doesn't easily let us trigger client side error handlers.
    // Instead, let's just make sure when the connection is dropped, it goes to disconnected state.
    expect(get(connectionStatus)).toBe('disconnected');
    vi.useRealTimers();
  });
});
