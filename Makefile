# Unified Commerce Platform - Development Makefile

.PHONY: help setup start-infra stop-infra start-services stop-services start-frontend test build clean

# Default target
help:
	@echo "Unified Commerce Platform - Development Commands"
	@echo "================================================"
	@echo "setup          - Initial project setup"
	@echo "start-infra    - Start infrastructure services (databases, etc.)"
	@echo "stop-infra     - Stop infrastructure services"
	@echo "start-services - Start all microservices"
	@echo "stop-services  - Stop all microservices"
	@echo "start-frontend - Start frontend applications"
	@echo "test           - Run all tests"
	@echo "build          - Build all services and applications"
	@echo "clean          - Clean up generated files and containers"

# Initial setup
setup:
	@echo "Setting up Unified Commerce Platform..."
	@go mod init unified-commerce
	@echo "Installing development dependencies..."
	@go install github.com/air-verse/air@latest
	@echo "Setup complete!"

# Infrastructure services
start-infra:
	@echo "Starting infrastructure services..."
	@docker-compose up -d postgres mongodb redis elasticsearch kafka zookeeper prometheus grafana
	@echo "Infrastructure services started!"
	@echo "Services available at:"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  MongoDB: localhost:27017"
	@echo "  Redis: localhost:6379"
	@echo "  Elasticsearch: localhost:9200"
	@echo "  Kafka: localhost:9092"
	@echo "  Prometheus: http://localhost:9090"
	@echo "  Grafana: http://localhost:3000 (admin/admin)"

stop-infra:
	@echo "Stopping infrastructure services..."
	@docker-compose down
	@echo "Infrastructure services stopped!"

# Microservices
start-services:
	@echo "Starting microservices..."
	@cd services/identity && air &
	@cd services/merchant-account && air &
	@cd services/product-catalog && air &
	@cd services/inventory && air &
	@cd services/order && air &
	@cd services/cart-checkout && air &
	@cd services/payments && air &
	@cd services/promotions && air &
	@echo "Microservices started in development mode!"

stop-services:
	@echo "Stopping microservices..."
	@pkill -f "air"
	@echo "Microservices stopped!"

# Frontend applications
start-frontend:
	@echo "Starting frontend applications..."
	@cd storefront && npm run dev &
	@cd admin-panel && npm run dev &
	@cd gateway && npm run dev &
	@echo "Frontend applications started!"
	@echo "  Storefront: http://localhost:3001"
	@echo "  Admin Panel: http://localhost:3002"
	@echo "  GraphQL Gateway: http://localhost:4000"

# Testing
test:
	@echo "Running tests..."
	@go test ./services/... -v
	@cd storefront && npm test
	@cd admin-panel && npm test
	@echo "All tests completed!"

# Build
build:
	@echo "Building all services..."
	@cd services/identity && go build -o ../../bin/identity ./cmd/server
	@cd services/merchant-account && go build -o ../../bin/merchant-account ./cmd/server
	@cd services/product-catalog && go build -o ../../bin/product-catalog ./cmd/server
	@cd services/inventory && go build -o ../../bin/inventory ./cmd/server
	@cd services/order && go build -o ../../bin/order ./cmd/server
	@cd services/cart-checkout && go build -o ../../bin/cart-checkout ./cmd/server
	@cd services/payments && go build -o ../../bin/payments ./cmd/server
	@cd services/promotions && go build -o ../../bin/promotions ./cmd/server
	@cd storefront && npm run build
	@cd admin-panel && npm run build
	@cd gateway && npm run build
	@echo "Build completed!"

# Cleanup
clean:
	@echo "Cleaning up..."
	@docker-compose down -v
	@docker system prune -f
	@rm -rf bin/
	@rm -rf */node_modules/
	@rm -rf */dist/
	@rm -rf */.next/
	@echo "Cleanup completed!"

# Development shortcuts
dev: start-infra start-services start-frontend
	@echo "Full development environment started!"

logs:
	@docker-compose logs -f

status:
	@docker-compose ps
	@echo ""
	@echo "Microservices status:"
	@ps aux | grep -E "(identity|merchant|product|inventory|order|cart|payment|promotion)" | grep -v grep || echo "No microservices running"