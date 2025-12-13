package otel

import (
	"context"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	httpMetricsOnce sync.Once
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	activeRequests  metric.Int64UpDownCounter
)

// initHTTPMetrics initializes the HTTP metrics instruments
func initHTTPMetrics() {
	httpMetricsOnce.Do(func() {
		meter := otel.Meter("gourmet-http")

		var err error
		requestCounter, err = meter.Int64Counter(
			"http.server.request.count",
			metric.WithDescription("Total number of HTTP requests"),
			metric.WithUnit("1"),
		)
		if err != nil {
			panic(err)
		}

		requestDuration, err = meter.Float64Histogram(
			"http.server.request.duration",
			metric.WithDescription("HTTP request duration"),
			metric.WithUnit("ms"),
		)
		if err != nil {
			panic(err)
		}

		activeRequests, err = meter.Int64UpDownCounter(
			"http.server.active_requests",
			metric.WithDescription("Number of active HTTP requests"),
			metric.WithUnit("1"),
		)
		if err != nil {
			panic(err)
		}
	})
}

// HTTPMetricsMiddleware returns a middleware that records HTTP metrics
func HTTPMetricsMiddleware(next http.Handler) http.Handler {
	initHTTPMetrics()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		// Track active requests
		activeRequests.Add(ctx, 1)
		defer activeRequests.Add(ctx, -1)

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start).Milliseconds()
		attrs := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.route", r.URL.Path),
			attribute.Int("http.status_code", wrapped.statusCode),
		}

		requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		requestDuration.Record(ctx, float64(duration), metric.WithAttributes(attrs...))
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

// RecordCustomMetric is a helper function to record custom metrics
func RecordCustomMetric(ctx context.Context, name string, value int64, attrs ...attribute.KeyValue) {
	meter := otel.Meter("gourmet-custom")
	counter, err := meter.Int64Counter(name)
	if err != nil {
		return
	}
	counter.Add(ctx, value, metric.WithAttributes(attrs...))
}

// RecordCustomDuration is a helper function to record custom duration metrics
func RecordCustomDuration(ctx context.Context, name string, duration time.Duration, attrs ...attribute.KeyValue) {
	meter := otel.Meter("gourmet-custom")
	histogram, err := meter.Float64Histogram(name, metric.WithUnit("ms"))
	if err != nil {
		return
	}
	histogram.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
}

// HTTPTraceMiddleware returns a middleware that creates trace spans for HTTP requests
func HTTPTraceMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("gourmet-http")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract trace context from incoming request
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// Start a new span
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.target", r.URL.Path),
				attribute.String("http.scheme", r.URL.Scheme),
				attribute.String("http.host", r.Host),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.remote_addr", r.RemoteAddr),
			),
		)
		defer span.End()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler with traced context
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// Record span attributes based on response
		span.SetAttributes(
			attribute.Int("http.status_code", wrapped.statusCode),
		)

		// Set span status based on HTTP status code
		if wrapped.statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}

// HTTPObservabilityMiddleware combines both metrics and tracing middleware
func HTTPObservabilityMiddleware(next http.Handler) http.Handler {
	initHTTPMetrics()
	tracer := otel.Tracer("gourmet-http")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Extract trace context from incoming request
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// Start a new span
		spanName := r.Method + " " + r.URL.Path
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.target", r.URL.Path),
				attribute.String("http.scheme", r.URL.Scheme),
				attribute.String("http.host", r.Host),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.remote_addr", r.RemoteAddr),
			),
		)
		defer span.End()

		// Track active requests
		activeRequests.Add(ctx, 1)
		defer activeRequests.Add(ctx, -1)

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler with traced context
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// Record metrics and span attributes
		duration := time.Since(start).Milliseconds()
		attrs := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.route", r.URL.Path),
			attribute.Int("http.status_code", wrapped.statusCode),
		}

		// Record metrics
		requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
		requestDuration.Record(ctx, float64(duration), metric.WithAttributes(attrs...))

		// Update span
		span.SetAttributes(attribute.Int("http.status_code", wrapped.statusCode))
		if wrapped.statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}

// StartSpan is a helper function to start a custom span
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	tracer := otel.Tracer("gourmet-custom")
	return tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}
