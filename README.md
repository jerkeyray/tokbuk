```
██████ ▄████▄ ██ ▄█▀ █████▄ ██  ██ ██ ▄█▀ 
██   ██  ██ ████   ██▄▄██ ██  ██ ████   
██   ▀████▀ ██ ▀█▄ ██▄▄█▀ ▀████▀ ██ ▀█▄ 

```
                                                   
## tokbuk

**tokbuk** is a small, focused token bucket rate limiter written in Go.

It implements:

- **Token bucket** rate limiting with a fixed capacity and refill rate
- **Per-client bucket management**, so each client (IP, user ID, API key, etc.) gets its own limiter
- **HTTP middleware** for `net/http` that returns **HTTP 429 (Too Many Requests)** when a client exceeds its limit

There are no external dependencies beyond the Go standard library.

---

## How to run it

### Run the demo server

From the repository root:

```bash
go run ./cmd/demo
```

You should see log output like:

```text
middleware demo running at http://localhost:8080/test
```

In another terminal, send some requests:

```bash
# A few single requests (should be allowed)
curl -i http://localhost:8080/test

# Send many requests quickly to hit the limit
for i in $(seq 1 30); do
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/test
done
```

You should see a mix of:

- `200` responses with body `allowed`
- `429` responses with body `rate limit exceeded`

The limit in the demo is configured as:

- **Capacity**: 10 tokens
- **Refill rate**: 5 tokens/second per client IP

---

### Run tests

From the repository root:

```bash
go test ./...
```

This runs tests for both the core limiter and the HTTP middleware.

---

### Run benchmarks

The limiter package includes a simple benchmark for concurrent calls to `Allow`.

```bash
go test -bench=. ./internal/limiter
```

To see allocation data as well:

```bash
go test -bench=. -benchmem ./internal/limiter
```
