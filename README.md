# VLX_FrameFlow_GUI

## Overview
VLX FrameFlow GUI is a management layer for the VLX_FrameFlow suite. It provides a graphical interface and a centralized control system to manage Debian-based Single Board Computers (SBCs) acting as high-reliability mobile routers and video encoders.

## Project Goals
- Provide a web-based dashboard for real-time system monitoring.
- Manage network bonding, routing policies, and video encoding parameters without SSH.
- Visualize real-time telemetry (GPS, signal strength, bitrate) using low-latency updates.
- Enable remote field operations via a dedicated mobile interface.

## Architecture
The system is built with a focus on performance and low resource consumption for SBC hardware:

1. Backend (Golang):
- A lightweight compiled binary serving as a system daemon.
- Wraps core VLX_FrameFlow Bash scripts using 'os/exec'.
- Exposes a REST API for configuration and status.
- Implements WebSockets for real-time telemetry streaming to clients.
- Includes background workers to monitor system load, network interfaces, and external hardware logs.

2. Web Frontend (Web-based SPA):
- A reactive Single Page Application built with Svelte and Vite.
- Communicates via WebSockets to provide live updates on network performance, system load, GPS status, and video (FFmpeg) logs.
- Incorporates a Control Panel making fetch requests to the Go backend REST API to configure and manage VLX_FrameFlow services directly from the dashboard.

3. Mobile Client (Android):
- Future native or cross-platform application.
- Interfaces with the same Go backend API/WebSocket endpoints for mobile field use.

## File Tree Structure
```bash
.
├── backend/                # Go source code
│   ├── cmd/
│   │   └── controller/     # Main entry point for the daemon
│   ├── internal/           # Private application packages
│   │   ├── api/            # REST and WebSocket handlers
│   │   ├── system/         # Bash script wrappers and monitoring
│   │   └── models/         # Data structures and shared types
│   └── go.mod              # Go module definition
├── frontend/               # Web application source code (Svelte)
│   ├── src/
│   ├── public/
│   └── package.json
├── scripts/                # Local mock implementations of core VLX_FrameFlow bash scripts for testing and development
├── docs/                   # Documentation and API specifications
├── build/                  # Deployment scripts and compiled binaries
└── README.md
```

## Building and Deployment

To build the complete application (compiling the frontend and the Go backend for both AMD64 and ARM64 architectures), execute the root script:

```bash
./build.sh
```

Compiled binaries will be output to the `bin/` directory.

A systemd service template for deployment is located at `build/vlx_frameflow_gui.service`, which assumes the application is deployed to `/opt/VLX_FrameFlow_GUI`.
