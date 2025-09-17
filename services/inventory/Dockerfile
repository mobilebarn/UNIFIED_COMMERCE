# Use Go 1.21 for Railway compatibility
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o app ./cmd/server

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary
COPY --from=builder /app/app .

# Expose port
EXPOSE 8005

# Run the application
CMD ["./app"]