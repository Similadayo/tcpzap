# tcpzap

A high-performance, message-oriented TCP library in Go.

## Features

- **Message-Based**: Length-prefixed framing for easy message handling.
- **Low-Latency**: Custom AIMD congestion control optimized for short connections.
- **Resilient**: Configurable retries for transient failures.
- **Observable**: Latency metrics via callback.
- **Concurrent**: Thread-safe design with Go’s concurrency primitives.
- **Tested**: Comprehensive unit and integration tests.

## Installation

```bash
go get github.com/yourname/tcpzap
```

## Quick Start

Here’s a simple echo client and server:

### Example

```go
package main

import (
    "context"
    "log"
    "github.com/similadayo/tcpzap/internal/metrics"
 "github.com/similadayo/tcpzap/pkg/tcpzap"
)

func main() {
    cfg := tcpzap.DefaultConfig()
    cfg.MetricsFunc = func(m metrics.Metrics) {
        log.Printf("Latency: %v, Target: %s, Success: %v", m.Latency, m.Target, m.Success)
    }

    // Server
    ctx := context.Background()
    srv, err := tcpzap.NewServer(":8080", cfg)
    if err != nil {
        log.Fatal(err)
    }
    go func() { log.Fatal(srv.Serve(ctx, echoHandler{})) }()

    // Client
    cli, err := tcpzap.NewClient("localhost:8080", cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer cli.Close()

    resp, err := cli.Send(ctx, []byte("Hello"))
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Got: %q", resp)
}

type echoHandler struct{}

func (h echoHandler) Handle(_ context.Context, msg []byte) ([]byte, error) {
    return append([]byte("Echo: "), msg...), nil
}
```

### Run it

```bash
go run main.go
```

### Expected Output

```text
2025/03/15 01:07:58 Latency: 611.4µs, Target: localhost:8080, Success: true
2025/03/15 01:07:58 Got: "Echo: Hello"
```

## Configuration

Customize `tcpzap` with `Config`:

```go
cfg := tcpzap.Config{
    Timeout:     10 * time.Second,         // Connection/send timeout
    Retries:     5,                        // Max retry attempts
    RetryDelay:  200 * time.Millisecond,   // Delay between retries
    MetricsFunc: func(m metrics.Metrics) { // Optional metrics callback
        log.Printf("Latency: %v", m.Latency)
    },
}
```

## API Reference

- `tcpzap.NewClient(addr string, cfg Config) (*Client, error)`: Connects to a server.
- `Client.Send(ctx context.Context, msg []byte) ([]byte, error)`: Sends a message, returns the response.
- `Client.Close() error`: Closes the connection.
- `tcpzap.NewServer(addr string, cfg Config) (*Server, error)`: Starts a server.
- `Server.Serve(ctx context.Context, h Handler) error`: Listens for connections.
- `Server.Close() error`: Shuts down the server.
- **Handler Interface**: Implement `Handle(ctx context.Context, msg []byte) ([]byte, error)`.

See GoDoc for details (if published).

## Example: Chat Server

For a full example, check `cmd/echo/main.go` in the repo.

## Contributing

Found a bug or have an idea? Open an issue or PR on GitHub. Run tests with:

```bash
go test ./tests
```

## License

MIT
