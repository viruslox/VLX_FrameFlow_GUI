# VLX_FrameFlow_GUI Overview & Integration Analysis

## Project Overview

VLX_FrameFlow_GUI is a comprehensive management layer designed to interface with the VLX_FrameFlow suite. It offers a graphical dashboard and centralized control system for managing Debian-based Single Board Computers (SBCs) that act as high-reliability mobile routers and video encoders.

The application architecture consists of:
1. **Backend (Golang):** A lightweight daemon that exposes a REST API for configuration and utilizes WebSockets for real-time telemetry streaming (e.g., system load, network interfaces, GPS data, FFmpeg logs). It acts as a wrapper around the core VLX_FrameFlow Bash scripts.
2. **Frontend (Svelte):** A reactive Single Page Application (SPA) that provides a real-time dashboard. It connects to the Go backend via WebSockets for live updates and makes HTTP requests to the REST API for service configuration.
3. **Deployment:** The Svelte frontend is compiled and embedded into the Go binary. A single unified binary is produced for target architectures (`linux/amd64` and `linux/arm64`), simplifying deployment.

## Integration with VLX_FrameFlow

The GUI integrates with the underlying `VLX_FrameFlow` repository (https://github.com/viruslox/VLX_FrameFlow) through shell script execution:

*   **Executor Module:** The Go backend leverages the `os/exec` package within `backend/internal/system/executor.go` to invoke shell commands.
*   **Command Execution:** It executes commands using `bash -ic '<script> "$@"' -- <args>`. This approach allows it to inherit the environment of the interactive shell, which is crucial because the `VLX_FrameFlow` suite heavily relies on environment variables and potentially shell aliases (like `VLX_FrameFlow`, `VLX_mediamtx`, `VLX_gps_tracker`, `VLX_cameraman`) defined during its installation process.
*   **API Mapping:** The REST API endpoints defined in `backend/internal/api/handlers.go` directly map to specific script commands. For example:
    *   `/api/frameflow/client/:action` calls `VLX_FrameFlow client <action>`
    *   `/api/mediamtx/:action` calls `VLX_mediamtx <action>`
    *   `/api/cameraman/:action` calls `VLX_cameraman <device> <action>`
*   **Mocking for Development:** The repository includes a `scripts/` directory containing mock implementations of these core bash scripts (`VLX_FrameFlow`, `VLX_cameraman`, etc.). This allows the GUI to be developed and tested locally without requiring a full SBC environment or the actual VLX_FrameFlow suite to be installed.

## Codebase Maintenance Performed

During the review process, the following maintenance actions were performed:
*   **Backend:** Ran `go fmt ./...` and `go vet ./...` to ensure code formatting and identify potential issues. The Go code was already clean.
*   **Frontend:** Installed necessary ESLint and Prettier dependencies. Formatted the Svelte codebase using Prettier and configured ESLint to lint the files. Two missing `key` properties in `{#each}` blocks within Svelte components (`NetworkStatus.svelte` and `VideoStatus.svelte`) were identified and fixed to adhere to Svelte best practices.
*   **Cleanup:** A search was conducted for temporary files (`*.tmp`, `*.log`, `.DS_Store`), and none were found, indicating a clean repository state.
