# Load-Balancer

Small HTTP reverse-proxy / load-balancer for local testing and learning.

## Overview

This project implements a basic reverse proxy that forwards requests to one or more backends per logical resource. Features include:

- Correct module path in `go.mod` for proper imports
- Reliable YAML config unmarshalling
- Standardized config key: `server.listen_port` (string)
- Thread-safe round-robin backend selection per resource
- Diagnostic logging (bind address printed on startup)

## File Structure

- `cmd/main.go` — Program entrypoint (`server.Run()`)
- `internal/config/config.go` — Reads `data/config.yaml` (uses viper)
- `data/config.yaml` — Example configuration (server + resources)
- `internal/server/server.go` — Router wiring and server startup
- `internal/server/proxy_handlers.go` — Request handler + proxy logic

## Quick Start (PowerShell)

1. **Build the project:**

    ```powershell
    cd C:\Users\HP\Desktop\LoadBalancer
    go build -o loadbalancer.exe ./cmd
    ```

2. **Run detached and capture logs:**

    ```powershell
    # start detached and redirect stdout/stderr to files
    Start-Process -FilePath .\loadbalancer.exe -WorkingDirectory (Get-Location) -RedirectStandardOutput .\server.out -RedirectStandardError .\server.err -NoNewWindow
    Start-Sleep -Milliseconds 500
    Invoke-WebRequest -Uri 'http://localhost:8089/ping' -Method Head -UseBasicParsing
    Get-Content .\server.out -Tail 50
    Get-Content .\server.err -Tail 50
    ```

3. **Or run in foreground:**

    ```powershell
    go run ./cmd
    ```

## Configuration

Edit `data/config.yaml`. Key fields:

- `server.host` (string): Host to bind (e.g. `localhost`)
- `server.listen_port` (string): Port to bind (e.g. `"8089"`)
- `resources`: List of resources. Each resource:
    - `name`: Logical name
    - `endpoint`: Path prefix to register (e.g. `/server1`)
    - `destination_urls`: Array of backend URLs (e.g. `["http://localhost:9001", "http://localhost:9002"]`)

**Example resource entry:**

```yaml
resources:
  - name: "Server1"
    endpoint: "/server1"
    destination_urls:
      - "http://localhost:9001"
      - "http://localhost:9002"
```

## Round-Robin Behavior

Each resource uses a `LoadBalancer` instance (see `internal/server/server.go`) for thread-safe round-robin backend selection. The handler in `internal/server/proxy_handlers.go` picks the next backend for each request and forwards using Go's `httputil.NewSingleHostReverseProxy`.

## Diagnostics & Troubleshooting

- Check `server.out` / `server.err` for logs when running detached.
- On startup, server prints: `listening on localhost:8089`.
- Use PowerShell's `Invoke-WebRequest` for testing:

    ```powershell
    Invoke-WebRequest -Uri 'http://localhost:8089/server1' -Method Head -UseBasicParsing
    ```

## Recent Changes

- `go.mod`: Module path set to `github.com/Indroneel007/Load-Balancer`
- `internal/config/config.go`: Fixed viper.Unmarshal usage
- `data/config.yaml`: Uses `listen_port` and `destination_urls` list
- `internal/server/server.go`: Added `LoadBalancer` type and bind address print
- `internal/server/proxy_handlers.go`: Handler updated for round-robin backend selection

## Next Steps

- Add config validation (fail-fast if `listen_port` missing or resources empty)
- Support both `destination_url` (single) and `destination_urls` (array) in YAML
- Add health-checking and backend removal for unhealthy backends
- Add structured logging and configurable log files/levels

If you need help or something doesn't match your local files, let me know which file to inspect.

---
_Last updated: Nov 1, 2025_
