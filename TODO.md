# Phase 1: Foundation & Scaffold
- [x] Setup repository structure (backend/ and frontend/ folders).
- [x] Initialize Go module and install basic dependencies (e.g., Gin or Fiber for HTTP, Gorilla for WebSockets).
- [x] Setup basic HTTP server with a health-check endpoint.
- [x] Initialize Frontend SPA boilerplate.

# Phase 2: Backend Core & Script Wrapping
- [x] Create Go interfaces for executing and capturing output from underlying Bash scripts (`os/exec`).
- [x] Implement mock Bash scripts for local development without an SBC.
- [x] Develop REST API endpoints for starting/stopping services (Routing, Video Encoding).
- [x] Implement the WebSocket hub for real-time data broadcasting.

# Phase 3: Telemetry & System Monitoring
- [x] Create Go workers to parse network interfaces and system load.
- [x] Create Go workers to read GPS serial data (if applicable) and FFmpeg logs.
- [x] Push telemetry data through WebSockets to connected clients.

# Phase 4: Frontend Development
- [x] Build the main Dashboard UI (Network status, Video status).
- [x] Implement WebSocket client to handle real-time UI updates.
- [x] Create forms for configuration (editing bonding settings, bitrates, etc.).

# Phase 5: Build & Packaging
- [x] Embed the compiled Frontend SPA directly into the Go binary using the `embed` package.
- [x] Write a build script (`build.sh`) to compile for `linux/arm64` and `linux/amd64`.
- [x] Create a systemd service file to run the Go binary on boot.

# Phase 6: Mobile (Future)
- [ ] Evaluate Flutter/Kotlin for Android app.
- [ ] Connect Android app to Go REST/WS APIs.
