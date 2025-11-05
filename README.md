# Load-Balancer

Small HTTP reverse-proxy / load-balancer for local testing and learning.

## Overview

This project implements a basic reverse proxy that forwards requests to one or more backends per logical resource. Features include:

- Correct module path in `go.mod` for proper imports
- Reliable YAML config unmarshalling
- Standardized config key: `server.listen_port` (string)
- **Thread-safe round-robin and least response time backend selection per resource**
- **SSL redirection headers set for secure forwarding**
- **Mutex usage for safe concurrent server requests**
- Diagnostic logging (bind address printed on startup)

## File Structure

- `cmd/main.go` — Program entrypoint (`server.Run()`)
- `internal/config/config.go` — Reads `data/config.yaml` (uses viper)
- `data/config.yaml` — Example configuration (server + resources)
- `internal/server/server.go` — Router wiring, server startup, load balancing algorithms
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

## Load Balancing Algorithms

- **Round-Robin:** Each resource uses a thread-safe round-robin algorithm for backend selection.
- **Least Response Time:** Requests can be routed to the backend with the lowest recent response time for improved performance.

## SSL Redirection

- The proxy sets appropriate headers to support SSL redirection, ensuring secure forwarding of requests to backends.

## Concurrency

- Mutexes are used to safely handle concurrent requests and backend selection, preventing race conditions.

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
- `internal/server/server.go`: Added round-robin and least response time algorithms, mutex for concurrency, bind address print
- `internal/server/proxy_handlers.go`: Handler updated for backend selection and SSL headers

## Next Steps

- Add config validation (fail-fast if `listen_port` missing or resources empty)
- Support both `destination_url` (single) and `destination_urls` (array) in YAML
- Add health-checking and backend removal for unhealthy backends
- Add structured logging and configurable log files/levels

If you need help or something doesn't match your local files, let me know which file to inspect.

---
_Last updated: Nov 1, 2025_
