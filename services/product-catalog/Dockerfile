# Multi-stage build for Go application
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application (corrected path)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o product-catalog ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/product-catalog .

# Create uploads directory
RUN mkdir -p uploads && chown -R appuser:appgroup uploads

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app/

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8003

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8003/health || exit 1

# Command to run
CMD ["./product-catalog"]