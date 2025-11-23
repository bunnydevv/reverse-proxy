# Reverse Proxy Service

A production-ready reverse proxy service written in Go with advanced features including load balancing, health checks, and flexible configuration.

## Features

- **Multiple Load Balancing Algorithms**
  - Round Robin
  - Least Connections
  - Weighted Distribution

- **Health Checks**
  - Automatic backend health monitoring
  - Configurable intervals and timeouts
  - Automatic backend recovery

- **Flexible Configuration**
  - YAML-based configuration
  - Hot-reloadable settings
  - Multiple backend support

- **Production Ready**
  - Graceful shutdown
  - Connection tracking
  - Error handling
  - Logging

## Installation

```bash
go mod download
go build -o reverse-proxy
```

## Configuration

Edit `config.yaml` to configure your reverse proxy:

```yaml
server:
  address: ":8080"
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 120s

backends:
  - url: "http://localhost:8081"
    weight: 2
  - url: "http://localhost:8082"
    weight: 1

load_balancer:
  algorithm: "round-robin"  # Options: round-robin, least-connections, weighted

health_check:
  enabled: true
  interval: 10s
  timeout: 5s
  path: "/health"
```

## Usage

### Start the reverse proxy

```bash
./reverse-proxy -config config.yaml
```

### Command-line Options

- `-config`: Path to configuration file (default: `config.yaml`)

## Load Balancing Algorithms

### Round Robin
Distributes requests evenly across all healthy backends in sequence.

```yaml
load_balancer:
  algorithm: "round-robin"
```

### Least Connections
Routes requests to the backend with the fewest active connections.

```yaml
load_balancer:
  algorithm: "least-connections"
```

### Weighted
Distributes requests based on backend weights (higher weight = more requests).

```yaml
load_balancer:
  algorithm: "weighted"
backends:
  - url: "http://backend1:8081"
    weight: 3  # Receives 3x more requests
  - url: "http://backend2:8082"
    weight: 1
```

## Health Checks

The reverse proxy automatically monitors backend health:

- Sends HTTP GET requests to configured health check path
- Marks backends as unhealthy if they fail to respond
- Automatically recovers backends when they become healthy again
- Configurable check intervals and timeouts

## Architecture

```
Client Request
     ↓
Reverse Proxy (Port 8080)
     ↓
Load Balancer (Algorithm Selection)
     ↓
Backend Selection (Health Check)
     ↓
Backend Servers (8081, 8082, 8083, ...)
```

## Testing

### Create test backend servers

```bash
# Backend 1
python3 -m http.server 8081

# Backend 2
python3 -m http.server 8082

# Backend 3
python3 -m http.server 8083
```

### Send requests

```bash
curl http://localhost:8080
```

## Development

### Project Structure

```
.
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration management
├── proxy/
│   ├── proxy.go           # Main proxy logic
│   ├── load_balancer.go   # Load balancing algorithms
│   └── health_check.go    # Health check implementation
├── config.yaml            # Configuration file
└── README.md
```

### Building

```bash
go build -o reverse-proxy
```

### Running Tests

```bash
go test ./...
```

## License

MIT License
