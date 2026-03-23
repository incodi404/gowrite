# GoWrite

> A Golang-based HTTP Proxy Server — built with intent, not vibes.

[![Go](https://img.shields.io/badge/Go-100%25-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20windows%20%7C%20mac-lightgrey?style=flat-square)]()

---

## What is GoWrite?

**GoWrite** is a lightweight, high-performance HTTP Proxy Server written in Go. While the software world drowns in AI-generated boilerplate and drag-and-drop abstractions, GoWrite is a deliberate act of engineering — built from first principles, with full control over every byte that passes through it.

This project exists because **real engineers still write real software**.

---

## Features

- 🔁 **HTTP Request Forwarding** — Intercepts and forwards client HTTP requests to target servers
- 🔒 **HTTPS Tunneling** — Full support for the `CONNECT` method for encrypted traffic
- 📋 **Request/Response Logging** — Visibility into all traffic flowing through the proxy
- ⚙️ **Configurable** — Tune port, timeouts, and other server settings via config

---

## Project Structure

```
gowrite/
├── cmd/
│   └── proxy/
│       └── main.go          # App entry point
│
├── internal/
│   ├── proxy/
│   │   ├── handler.go       # Handles incoming client requests
│   │   ├── forward.go       # Forwards requests to target servers
│   │   ├── response.go      # Sends response back to client
│   │   └── tunnel.go        # CONNECT method (HTTPS tunneling)
│   │
│   ├── config/
│   │   └── config.go        # Config loading (port, timeouts, etc.)
│   │
│   ├── logger/
│   │   └── logger.go        # Request / response logging
│   │
│   └── middleware/
│       └── auth.go          # Auth, filtering, rate limiting (optional)
│
├── pkg/
│   └── utils/
│       └── net.go           # Network helpers (copy streams, headers)
│
├── cert/                    # TLS certificates
├── go.mod
└── README.md
```

---

## Getting Started

### Prerequisites

- Go `1.21+`
- Git

### Clone & Run

```bash
git clone https://github.com/incodi404/gowrite.git
cd gowrite
go mod tidy
go run ./cmd/proxy
```

Configure your HTTP client or browser to use `localhost:<PORT>` as the proxy.

---

## Build

### Build for Linux (from Windows PowerShell)

```powershell
# Step 1 — Set environment
$env:CGO_ENABLED=0; $env:GOOS="linux"; $env:GOARCH="amd64"

# Step 2 — Build
go build -trimpath -ldflags="-s -w" -o proxy ./cmd/proxy
```

### Build for current platform

```bash
go build -o gowrite ./cmd/proxy
./gowrite
```

---

## How It Works

```
Client
  │
  ▼
┌──────────────────────────┐
│       GoWrite Proxy       │
│                           │
│  handler.go               │  ← Accepts & parses client request
│  forward.go               │  ← Proxies HTTP requests upstream
│  tunnel.go                │  ← Tunnels HTTPS via CONNECT
│  response.go              │  ← Streams response back to client
│  logger.go                │  ← Logs method, URL, status, latency
└──────────────────────────┘
  │
  ▼
Target Server
```

For plain **HTTP**, GoWrite parses the request, forwards it to the target, and pipes the response back.

For **HTTPS**, GoWrite handles the `CONNECT` method — establishing a raw TCP tunnel between client and server without inspecting encrypted payloads.

---

## Philosophy

> *"The best way to understand a system is to build it."*

In a world of vibe-coded frameworks and prompt-engineered repos, GoWrite is different. Understanding how a proxy works — how connections are tunneled, how headers are forwarded, how goroutines handle concurrent I/O — is worth infinitely more than scaffolding one with a prompt.

This is software written by someone who **actually thought about it**.

---

## Roadmap

- [ ] HTTPS MITM with custom CA certificate support
- [ ] Request/Response header injection
- [ ] URL-based traffic filtering & blacklisting
- [ ] Web UI for real-time traffic inspection
- [ ] Docker support
- [ ] Upstream load balancing

---

## Author

**incodi404** — building real things in a vibe-coded world.
GitHub: [@incodi404](https://github.com/incodi404)

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<p align="center"><i>Written with intent. Not vibes.</i></p>
