# Load-Balancer

Small HTTP reverse-proxy / load-balancer used for local testing and learning.

This repo implements a basic reverse proxy that can forward requests to one
or more backends per logical resource. Recent work in this branch includes:

- fixing the module path in `go.mod` so imports resolve correctly
- fixing config unmarshalling so the YAML is read reliably
- standardizing the config key to `server.listen_port` (string)
- adding a simple thread-safe round-robin algorithm per resource
- small diagnostic logging to print the bind address on startup

Files you may care about
- `cmd/main.go` — program entrypoint (calls `server.Run()`)
- `internal/config/config.go` — reads `data/config.yaml` (uses viper)
- `data/config.yaml` — example configuration (server + resources)
- `internal/server/server.go` — router wiring and server startup
- `internal/server/proxy_handlers.go` — request handler + proxy logic

Quick start (PowerShell)

1. Build the project:

```powershell
cd C:\Users\HP\Desktop\LoadBalancer
go build -o loadbalancer.exe ./cmd
```

2. Run detached and capture logs (recommended for background runs):

```powershell
# start detached and redirect stdout/stderr to files
Start-Process -FilePath .\loadbalancer.exe -WorkingDirectory (Get-Location) -RedirectStandardOutput .\server.out -RedirectStandardError .\server.err -NoNewWindow
Start-Sleep -Milliseconds 500
Invoke-WebRequest -Uri 'http://localhost:8089/ping' -Method Head -UseBasicParsing
Get-Content .\server.out -Tail 50
Get-Content .\server.err -Tail 50
```

3. Or run in-foreground to see logs immediately:

```powershell
go run ./cmd
```

Configuration

The config is in `data/config.yaml`. Key points:

- `server.host` (string) — host to bind, e.g. `localhost`
- `server.listen_port` (string) — port to bind, e.g. `"8089"`
- `resources` — list of resources. Each resource should have:
	- `name` — a logical name
	- `endpoint` — the path prefix to register on the router (e.g. `/server1`)
	- `destination_urls` — an array of backend URLs (e.g. `["http://localhost:9001", "http://localhost:9002"]`)

Example resource entry:

```yaml
resources:
	- name: "Server1"
		endpoint: "/server1"
		destination_urls:
			- "http://localhost:9001"
			- "http://localhost:9002"
```

Round-robin behavior

Each resource now gets its own `LoadBalancer` instance (in `internal/server/server.go`) which
maintains a counter and exposes `Next(n int) int` to return the next backend index in a
thread-safe round-robin fashion. The handler in `internal/server/proxy_handlers.go` picks the
next backend for every incoming request and forwards the request using Go's
`httputil.NewSingleHostReverseProxy`.

Diagnostics & troubleshooting

- If the server fails to bind or not responding, check `server.out` / `server.err` when running detached.
- The server prints a line like `listening on localhost:8089` at startup (added for debugging).
- Use PowerShell's `Invoke-WebRequest` on Windows rather than `curl -I` for consistent output:

```powershell
Invoke-WebRequest -Uri 'http://localhost:8089/server1' -Method Head -UseBasicParsing
```

What changed in code (high level)

- `go.mod` — module path corrected to `github.com/Indroneel007/Load-Balancer`
- `internal/config/config.go` — fixed viper.Unmarshal usage; now unmarshals into a local struct and assigns
- `data/config.yaml` — switched to `listen_port` and `destination_urls` list format
- `internal/server/server.go` — added `LoadBalancer` type + `Next` method, wired each resource to a lb + destinations list, prints bind address
- `internal/server/proxy_handlers.go` — handler updated to accept a `LoadBalancer` and destinations slice and select backends round-robin

Next suggested improvements

- Add validation for config (fail-fast if `listen_port` is missing or resources are empty)
- Allow both `destination_url` (single) and `destination_urls` (array) in the YAML for more forgiving parsing (already partially supported)
- Add health-checking and active backend removal for unhealthy backends
- Add structured logging and configurable log files/levels

If you want, I can:
- add config validation and fail-fast behavior now
- remove the debug `listening on ...` print once you're satisfied
- add a minimal test that ensures round-robin cycles through backends

If anything here doesn't match what you see locally, tell me which file to inspect and I'll reconcile it.

---
Generated update: included fixes, run instructions, and round-robin support added on Nov 1, 2025.
