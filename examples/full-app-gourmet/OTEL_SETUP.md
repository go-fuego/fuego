# OpenTelemetry Observability Setup

This application includes comprehensive OpenTelemetry (OTEL) instrumentation that exports both metrics and distributed traces to an OTEL collector.

## Configuration

The application is configured using environment variables:

- `OTEL_EXPORTER_OTLP_ENDPOINT`: The endpoint of the OTEL collector (default: `otel-collector:4318`)
- `OTEL_SERVICE_NAME`: The name of this service in telemetry data (default: `gourmet-app`)
- `OTEL_SERVICE_VERSION`: The version of this service (default: `1.0.0`)

### Local Development

For local development, copy `.env.example` to `.env.local` and update the values:

```bash
cp .env.example .env.local
```

If running the OTEL collector locally (not in Docker), change the endpoint:

```
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
```

### Docker Deployment

When running in Docker, you can set environment variables in your Docker Compose file:

```yaml
services:
  gourmet-app:
    image: gourmet-app:latest
    environment:
      - OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318
      - OTEL_SERVICE_NAME=gourmet-app
      - OTEL_SERVICE_VERSION=1.0.0
    networks:
      - app-network

  otel-collector:
    image: otel/opentelemetry-collector:latest
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    networks:
      - app-network
    ports:
      - "4318:4318"  # OTLP HTTP receiver

networks:
  app-network:
    driver: bridge
```

Or pass them at runtime:

```bash
docker run -e OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318 \
           -e OTEL_SERVICE_NAME=gourmet-app \
           --network app-network \
           gourmet-app:latest
```

## Telemetry Data Collected

The application automatically collects both metrics and traces for all HTTP requests.

### Metrics

#### HTTP Server Metrics

- `http.server.request.count`: Total number of HTTP requests
  - Labels: `http.method`, `http.route`, `http.status_code`

- `http.server.request.duration`: HTTP request duration in milliseconds
  - Labels: `http.method`, `http.route`, `http.status_code`

- `http.server.active_requests`: Number of currently active HTTP requests

#### Custom Metrics

You can record custom metrics using the helper functions in the `otel` package:

```go
import (
    "github.com/go-fuego/fuego/examples/full-app-gourmet/otel"
    "go.opentelemetry.io/otel/attribute"
)

// Record a counter metric
otel.RecordCustomMetric(ctx, "recipes.created", 1,
    attribute.String("category", "dessert"))

// Record a duration metric
otel.RecordCustomDuration(ctx, "db.query.duration", duration,
    attribute.String("query", "get_recipes"))
```

### Distributed Traces

The application creates distributed trace spans for every HTTP request with the following attributes:

- `http.method`: HTTP method (GET, POST, etc.)
- `http.url`: Full request URL
- `http.target`: Request path
- `http.scheme`: URL scheme (http/https)
- `http.host`: Host header
- `http.user_agent`: User agent string
- `http.remote_addr`: Client IP address
- `http.status_code`: HTTP response status code

Traces automatically include:
- Span status (OK for 2xx/3xx, Error for 4xx/5xx)
- Trace context propagation (W3C Trace Context)
- Parent-child span relationships

#### Custom Spans

You can create custom spans within your handlers:

```go
import (
    "github.com/go-fuego/fuego/examples/full-app-gourmet/otel"
    "go.opentelemetry.io/otel/attribute"
)

func myHandler(c fuego.ContextNoBody) error {
    ctx := c.Context()

    // Start a custom span
    ctx, span := otel.StartSpan(ctx, "database.query",
        attribute.String("query", "SELECT * FROM recipes"))
    defer span.End()

    // Your code here...

    return nil
}
```

## OTEL Collector Configuration

Here's an example OTEL collector configuration that exports metrics to Prometheus and traces to Jaeger:

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
      grpc:
        endpoint: 0.0.0.0:4317

processors:
  batch:
    timeout: 10s
    send_batch_size: 1024

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
    namespace: gourmet

  otlp/jaeger:
    endpoint: "jaeger:4317"
    tls:
      insecure: true

  logging:
    loglevel: info

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [prometheus, logging]

    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp/jaeger, logging]
```

See `otel-collector-config.example.yaml` for a complete configuration with additional options.

## Verifying Telemetry

### Metrics

1. Check the application logs for successful OTEL initialization:
   ```
   INFO Initializing OpenTelemetry metrics endpoint=otel-collector:4318 service_name=gourmet-app version=1.0.0
   INFO OpenTelemetry metrics initialized successfully
   INFO Initializing OpenTelemetry traces endpoint=otel-collector:4318 service_name=gourmet-app version=1.0.0
   INFO OpenTelemetry traces initialized successfully
   ```

2. If using Prometheus exporter, metrics will be available at:
   ```
   http://localhost:8889/metrics
   ```

3. Access Prometheus UI to query metrics:
   ```
   http://localhost:9090
   ```

### Traces

1. Access Jaeger UI to view traces:
   ```
   http://localhost:16686
   ```

2. Select "gourmet-app" from the service dropdown

3. Click "Find Traces" to see all collected traces

4. Click on individual traces to see detailed span information

### Grafana

If using the full stack with Grafana:

1. Access Grafana:
   ```
   http://localhost:3000
   ```
   (username: admin, password: admin)

2. Add Prometheus as a data source (http://prometheus:9090)

3. Add Jaeger as a data source (http://jaeger:16686)

4. Create dashboards to visualize metrics and traces together

## Troubleshooting

- **Connection refused**: Ensure the OTEL collector is running and accessible on the specified endpoint
- **No metrics visible**: Check that the export interval (10 seconds) has passed and requests have been made to the app
- **No traces visible**: Ensure Jaeger is configured as an exporter in the OTEL collector and the collector is forwarding traces
- **Broken trace spans**: Verify trace context is being propagated by checking that handlers use the context from the request
- **Authentication errors**: This setup uses insecure connections for Docker internal networks. For production, configure TLS appropriately
- **High cardinality metrics**: Be careful with labels on custom metrics to avoid cardinality explosion

## Architecture

```
┌──────────────────────┐
│   Gourmet App        │
│   (Port 8083)        │
│                      │
│ ┌──────────────────┐ │
│ │ OTEL SDK         │ │
│ │ - Metrics        │ │
│ │ - Traces         │ │
│ └──────────────────┘ │
└──────────┬───────────┘
           │ HTTP/OTLP (Port 4318)
           ▼
┌──────────────────────┐
│  OTEL Collector      │
│  (Port 4318)         │
│                      │
│  Processors:         │
│  - Batch             │
│  - Memory Limiter    │
└──────────┬───────────┘
           │
           ├─────────► Prometheus (Port 8889)
           │            │
           │            └──► Grafana (Port 3000)
           │
           └─────────► Jaeger (Port 4317)
                       └──► Jaeger UI (Port 16686)
```

## Quick Start

1. Copy example files:
   ```bash
   cp docker-compose.example.yaml docker-compose.yaml
   cp otel-collector-config.example.yaml otel-collector-config.yaml
   cp prometheus.example.yml prometheus.yml
   ```

2. Start the stack:
   ```bash
   docker-compose up -d
   ```

3. Access the services:
   - Application: http://localhost:8083
   - Jaeger UI: http://localhost:16686
   - Prometheus: http://localhost:9090
   - Grafana: http://localhost:3000

4. Generate some traffic to the application

5. View traces in Jaeger and metrics in Prometheus/Grafana
