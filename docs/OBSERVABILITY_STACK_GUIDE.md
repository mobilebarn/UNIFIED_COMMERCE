# Observability Stack Implementation Guide

## Overview
This document provides instructions for implementing and configuring the observability stack for Retail OS. The observability stack includes logging, metrics, and distributed tracing to provide comprehensive monitoring and debugging capabilities.

## Architecture
The observability stack consists of the following components:

### Logging
- **ELK Stack**: Elasticsearch, Logstash, and Kibana for log aggregation, processing, and visualization
- **Filebeat**: Log shipper for collecting logs from containers and files

### Metrics
- **Prometheus**: Metrics collection and storage
- **Grafana**: Metrics visualization and dashboarding
- **OpenTelemetry Collector**: Metrics ingestion and export

### Distributed Tracing
- **OpenTelemetry Collector**: Traces ingestion and processing
- **Jaeger**: Distributed tracing storage and visualization

## Implementation Details

### OpenTelemetry Integration
All microservices have been updated to include OpenTelemetry instrumentation:

1. **Automatic Instrumentation**: Gin middleware automatically creates spans for HTTP requests
2. **Manual Instrumentation**: Services can create custom spans for business operations
3. **Context Propagation**: Trace context is propagated between services

### Metrics Collection
Services automatically collect the following metrics:

1. **HTTP Metrics**:
   - Request count by method, endpoint, and status
   - Request duration histograms
   - In-flight requests gauge

2. **Database Metrics**:
   - Active connections gauge
   - Connection count by status
   - Query duration histograms

3. **Application Metrics**:
   - Business operation counters
   - Error counters by type
   - Custom business metrics

### Logging
Services use structured logging with the following features:

1. **Context Enrichment**: Logs include request IDs, trace IDs, and span IDs
2. **Log Levels**: Support for debug, info, warn, error, and fatal levels
3. **Structured Format**: JSON formatted logs for easy parsing

## Configuration

### Environment Variables
Services use the following environment variables for observability configuration:

- `OTEL_EXPORTER_OTLP_ENDPOINT`: OTLP collector endpoint (default: localhost:4317)
- `OTEL_SERVICE_NAME`: Service name for tracing
- `LOG_LEVEL`: Logging level (default: info)

### Docker Configuration
The docker-compose.yml file includes services for:

1. **OpenTelemetry Collector**: Receives telemetry data from services
2. **Jaeger**: Stores and visualizes distributed traces
3. **Prometheus**: Collects and stores metrics
4. **Grafana**: Visualizes metrics and provides dashboards

## Usage

### Distributed Tracing
To create custom spans in your service code:

```go
// Start a new span
ctx, span := s.Tracer.StartSpan(context.Background(), "custom-operation")
defer span.End()

// Add attributes to the span
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("item.count", itemCount),
)

// Record an error if needed
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
}
```

### Metrics
To record custom metrics in your service code:

```go
// Record a business operation
s.Metrics.RecordBusinessOperation("order_created", "order")

// Record an application error
s.Metrics.RecordAppError("database", "query_failed")

// Record an application operation
s.Metrics.RecordAppOperation("user_login", "success")
```

### Logging
To add structured logging with trace context:

```go
// Log with context
s.Logger.WithFields(map[string]interface{}{
    "user_id": userID,
    "action": "login_attempt",
}).Info("User login attempt")

// Log an error with context
s.Logger.WithError(err).WithFields(map[string]interface{}{
    "user_id": userID,
}).Error("Failed to authenticate user")
```

## Monitoring and Visualization

### Jaeger UI
Access the Jaeger UI at: http://localhost:16686

Features:
- Trace search and filtering
- Trace detail views
- Service dependency graphs
- Performance statistics

### Grafana Dashboards
Access Grafana at: http://localhost:3001

Default credentials:
- Username: admin
- Password: admin

Pre-configured dashboards include:
- Service overview
- HTTP metrics
- Database metrics
- Business metrics

### Kibana
Access Kibana at: http://localhost:5601

Features:
- Log search and filtering
- Log visualization
- Dashboard creation

## Health Checks

### Service Health
All services expose health check endpoints at `/health` which include:
- Service status
- Database connectivity
- Dependency status

### Infrastructure Health
Infrastructure components include health checks for:
- Container status
- Port availability
- Resource utilization

## Troubleshooting

### Common Issues

1. **Traces not appearing in Jaeger**:
   - Check OTLP collector logs
   - Verify OTLP endpoint configuration
   - Ensure services can reach the collector

2. **Metrics not appearing in Grafana**:
   - Check Prometheus target status
   - Verify service metrics endpoints
   - Check network connectivity

3. **Logs not appearing in Kibana**:
   - Check Logstash pipeline configuration
   - Verify Filebeat configuration
   - Check Elasticsearch indices

### Debugging Commands

```bash
# Check OTLP collector logs
docker logs retail-os-otel-collector

# Check Jaeger logs
docker logs retail-os-jaeger

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Check service metrics
curl http://localhost:8001/metrics
```

## Security Considerations

### Network Security
- Services communicate over internal Docker network
- External access is restricted to necessary ports only
- TLS can be enabled for production deployments

### Data Security
- Sensitive data is not included in traces or logs
- Log data is retained according to compliance requirements
- Access to monitoring tools is protected by authentication

## Performance Optimization

### Resource Allocation
Recommended resource limits for observability components:
- OpenTelemetry Collector: 512MB RAM, 0.5 CPU
- Jaeger: 1GB RAM, 1 CPU
- Prometheus: 2GB RAM, 1 CPU
- Grafana: 512MB RAM, 0.5 CPU

### Sampling
- Trace sampling can be configured to reduce overhead
- Metric collection intervals can be adjusted
- Log levels can be adjusted for production

## Production Considerations

### High Availability
- Deploy multiple instances of observability components
- Use persistent storage for data retention
- Implement backup and restore procedures

### Scaling
- Horizontal scaling of OpenTelemetry collectors
- Sharding of Prometheus for large-scale metrics
- Load balancing for Grafana instances

### Data Retention
- Configure retention policies for traces, metrics, and logs
- Implement data archiving for long-term storage
- Regular cleanup of old data

## Future Enhancements

### Advanced Features
- Alerting based on metrics and logs
- Anomaly detection for performance issues
- Machine learning for log analysis
- Service mesh integration for enhanced observability

### Integration Opportunities
- Integration with incident management tools
- Alerting via Slack, email, and SMS
- Integration with CI/CD for deployment tracking
- Business intelligence dashboards