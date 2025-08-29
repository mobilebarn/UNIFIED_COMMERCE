package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsCollector holds all Prometheus metrics for a service
type MetricsCollector struct {
	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge

	// Database metrics
	DatabaseConnectionsActive prometheus.Gauge
	DatabaseConnectionsTotal  *prometheus.CounterVec
	DatabaseQueryDuration     *prometheus.HistogramVec

	// Application metrics
	AppRequestsTotal     *prometheus.CounterVec
	AppErrors            *prometheus.CounterVec
	AppBusinessMetrics   *prometheus.CounterVec
}

// NewMetricsCollector creates a new metrics collector for a service
func NewMetricsCollector(serviceName string) *MetricsCollector {
	metrics := &MetricsCollector{
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		HTTPRequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		),
		DatabaseConnectionsActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_connections_active",
				Help: "Number of active database connections",
			},
		),
		DatabaseConnectionsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_connections_total",
				Help: "Total number of database connections created",
			},
			[]string{"database_type", "status"},
		),
		DatabaseQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "database_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0},
			},
			[]string{"database_type", "operation"},
		),
		AppRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_requests_total",
				Help: "Total number of application requests by operation",
			},
			[]string{"operation", "status"},
		),
		AppErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_errors_total",
				Help: "Total number of application errors",
			},
			[]string{"error_type", "operation"},
		),
		AppBusinessMetrics: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_business_operations_total",
				Help: "Total number of business operations",
			},
			[]string{"operation", "entity_type"},
		),
	}

	// Register all metrics
	prometheus.MustRegister(
		metrics.HTTPRequestsTotal,
		metrics.HTTPRequestDuration,
		metrics.HTTPRequestsInFlight,
		metrics.DatabaseConnectionsActive,
		metrics.DatabaseConnectionsTotal,
		metrics.DatabaseQueryDuration,
		metrics.AppRequestsTotal,
		metrics.AppErrors,
		metrics.AppBusinessMetrics,
	)

	return metrics
}

// PrometheusMiddleware creates a Gin middleware for Prometheus metrics collection
func (mc *MetricsCollector) PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Increment in-flight requests
		mc.HTTPRequestsInFlight.Inc()
		defer mc.HTTPRequestsInFlight.Dec()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := string(rune(c.Writer.Status()))
		method := c.Request.Method
		endpoint := c.FullPath()

		// If endpoint is empty (404), use the raw path
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		mc.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		mc.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}

// RecordDatabaseOperation records database operation metrics
func (mc *MetricsCollector) RecordDatabaseOperation(dbType, operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	
	mc.DatabaseConnectionsTotal.WithLabelValues(dbType, status).Inc()
	mc.DatabaseQueryDuration.WithLabelValues(dbType, operation).Observe(duration.Seconds())
}

// RecordAppOperation records application operation metrics
func (mc *MetricsCollector) RecordAppOperation(operation, status string) {
	mc.AppRequestsTotal.WithLabelValues(operation, status).Inc()
}

// RecordAppError records application error metrics
func (mc *MetricsCollector) RecordAppError(errorType, operation string) {
	mc.AppErrors.WithLabelValues(errorType, operation).Inc()
}

// RecordBusinessOperation records business operation metrics
func (mc *MetricsCollector) RecordBusinessOperation(operation, entityType string) {
	mc.AppBusinessMetrics.WithLabelValues(operation, entityType).Inc()
}

// PrometheusHandler returns the Prometheus metrics handler
func PrometheusHandler() gin.HandlerFunc {
	handler := promhttp.Handler()
	return gin.WrapH(handler)
}
