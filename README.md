# HTTP Proxy Server

A lightweight proxy server build in Golang to understand Golang and Networking.

## File Structure

http-proxy/
├── cmd/
│   └── proxy/
│       └── main.go        # App entry point
│
├── internal/
│   ├── proxy/
│   │   ├── handler.go    # Handles incoming client requests
│   │   ├── forward.go    # Forwards requests to target servers
│   │   ├── response.go   # Sends response back to client
│   │   └── tunnel.go     # CONNECT method (HTTPS tunneling)
│   │
│   ├── config/
│   │   └── config.go     # Config loading (port, timeouts, etc.)
│   │
│   ├── logger/
│   │   └── logger.go     # Request / response logging
│   │
│   └── middleware/
│       └── auth.go       # Auth, filtering, rate limiting (optional)
│
├── pkg/
│   └── utils/
│       └── net.go        # Network helpers (copy streams, headers)
│
├── go.mod
└── README.md

### Build file for Linux
Command-1: $env:CGO_ENABLED=0; $env:GOOS="linux"; $env:GOARCH="amd64";
Command-2: go build -trimpath -ldflags="-s -w" -o proxy ./cmd/proxy