# Custom HTTP Client Example

This example demonstrates how to use custom HTTP clients with the goitunes library.

## Features

- Retry logic for failed requests
- Request/response logging
- Custom timeouts and configurations
- Easy testing with mock clients

## Usage

```bash
cd examples/custom_http_client
go run main.go
```

## Code Explanation

### Implementing a Custom HTTP Client

Any type that implements the `infrahttp.Client` interface can be used:

```go
type Client interface {
    Do(req *http.Request) (*http.Response, error)
}
```

### Example 1: Retry Client

```go
type RetryHTTPClient struct {
    client      *http.Client
    maxRetries  int
    retryDelay  time.Duration
}

func (c *RetryHTTPClient) Do(req *http.Request) (*http.Response, error) {
    // Implement retry logic
    for attempt := 0; attempt <= c.maxRetries; attempt++ {
        resp, err := c.client.Do(req)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }
        time.Sleep(c.retryDelay)
    }
    return nil, errors.New("max retries exceeded")
}
```

### Example 2: Logging Client

```go
type LoggingHTTPClient struct {
    client *http.Client
}

func (c *LoggingHTTPClient) Do(req *http.Request) (*http.Response, error) {
    log.Printf("[HTTP] %s %s", req.Method, req.URL)
    return c.client.Do(req)
}
```

### Using Custom Client

```go
customClient := NewRetryHTTPClient(3, 2*time.Second)

client, err := goitunes.New("us",
    goitunes.WithHTTPClient(customClient),
)
```

## Use Cases

### 1. Retry Logic
Handle transient network failures automatically:
- Connection timeouts
- Temporary server errors (5xx)
- Rate limiting (with exponential backoff)

### 2. Logging & Monitoring
Track all API requests:
- Request/response logging
- Performance metrics
- Error tracking
- Debug information

### 3. Testing
Use mock clients for testing:
```go
type MockClient struct {
    responses map[string]*http.Response
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
    return m.responses[req.URL.Path], nil
}
```

### 4. Request Modification
Add custom behavior:
- Additional headers
- Request signing
- Proxy configuration
- Rate limiting

### 5. Response Caching
Cache responses to reduce API calls:
```go
type CachingClient struct {
    client *http.Client
    cache  map[string]*http.Response
}
```

## Benefits

- **Flexibility**: Full control over HTTP behavior
- **Testability**: Easy to mock for unit tests
- **Observability**: Add logging and metrics
- **Reliability**: Implement retry and circuit breaker patterns
- **Performance**: Add caching and connection pooling

## Advanced Example: Circuit Breaker

```go
type CircuitBreakerClient struct {
    client        *http.Client
    failureCount  int
    failureLimit  int
    resetTimeout  time.Duration
    lastFailure   time.Time
    state         string // "closed", "open", "half-open"
}

func (c *CircuitBreakerClient) Do(req *http.Request) (*http.Response, error) {
    if c.state == "open" {
        if time.Since(c.lastFailure) > c.resetTimeout {
            c.state = "half-open"
        } else {
            return nil, errors.New("circuit breaker open")
        }
    }
    
    resp, err := c.client.Do(req)
    
    if err != nil || resp.StatusCode >= 500 {
        c.failureCount++
        c.lastFailure = time.Now()
        if c.failureCount >= c.failureLimit {
            c.state = "open"
        }
        return resp, err
    }
    
    c.failureCount = 0
    c.state = "closed"
    return resp, nil
}
```

This pattern prevents cascading failures when the API is experiencing issues.

