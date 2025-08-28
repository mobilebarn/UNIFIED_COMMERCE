package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// BaseService provides common functionality for all microservices
type BaseService struct {
	Name         string
	Config       *config.Config
	Logger       *logger.Logger
	PostgresDB   *database.PostgresDB
	MongoDB      *database.MongoDB
	RedisClient  *database.RedisClient
	Router       *gin.Engine
	Server       *http.Server
	HealthChecks []HealthCheck
}

// HealthCheck represents a health check function
type HealthCheck struct {
	Name  string
	Check func(ctx context.Context) error
}

// ServiceOptions holds configuration options for creating a service
type ServiceOptions struct {
	Name         string
	UsePostgres  bool
	UseMongoDB   bool
	UseRedis     bool
	CustomRoutes func(*gin.Engine)
	HealthChecks []HealthCheck
}

// NewBaseService creates a new base service with common setup
func NewBaseService(opts ServiceOptions) (*BaseService, error) {
	// Load configuration
	cfg, err := config.LoadConfig(opts.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig(opts.Name)
	loggerConfig.Level = cfg.LogLevel
	serviceLogger := logger.NewLogger(loggerConfig)

	service := &BaseService{
		Name:   opts.Name,
		Config: cfg,
		Logger: serviceLogger,
	}

	// Set up database connections based on service requirements
	if opts.UsePostgres {
		postgresConfig := database.NewPostgresConfigFromEnv(
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseUser,
			cfg.DatabasePassword,
			cfg.DatabaseName,
		)

		service.PostgresDB, err = database.NewPostgresConnection(postgresConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}

		service.HealthChecks = append(service.HealthChecks, HealthCheck{
			Name: "postgres",
			Check: func(ctx context.Context) error {
				return service.PostgresDB.Health(ctx)
			},
		})
	}

	if opts.UseMongoDB {
		mongoConfig := database.NewMongoConfigFromEnv(
			cfg.MongoURL,
			cfg.MongoDatabase,
			cfg.MongoUser,
			cfg.MongoPassword,
		)

		service.MongoDB, err = database.NewMongoConnection(mongoConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
		}

		service.HealthChecks = append(service.HealthChecks, HealthCheck{
			Name: "mongodb",
			Check: func(ctx context.Context) error {
				return service.MongoDB.Health(ctx)
			},
		})
	}

	if opts.UseRedis {
		redisConfig := database.NewRedisConfigFromEnv(
			cfg.RedisURL,
			cfg.RedisPassword,
			cfg.RedisDB,
		)

		service.RedisClient, err = database.NewRedisConnection(redisConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Redis: %w", err)
		}

		service.HealthChecks = append(service.HealthChecks, HealthCheck{
			Name: "redis",
			Check: func(ctx context.Context) error {
				return service.RedisClient.Health(ctx)
			},
		})
	}

	// Add custom health checks
	service.HealthChecks = append(service.HealthChecks, opts.HealthChecks...)

	// Set up HTTP router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	service.Router = gin.New()

	// Add global middleware
	service.Router.Use(middleware.RequestID())
	service.Router.Use(middleware.Logger())
	service.Router.Use(middleware.Recovery())
	service.Router.Use(middleware.CORS())

	// Add common routes
	service.setupCommonRoutes()

	// Add custom routes if provided
	if opts.CustomRoutes != nil {
		opts.CustomRoutes(service.Router)
	}

	// Create HTTP server
	service.Server = &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      service.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return service, nil
}

// setupCommonRoutes sets up routes that are common to all services
func (s *BaseService) setupCommonRoutes() {
	// Health check endpoint
	s.Router.GET("/health", s.healthHandler)

	// Readiness check endpoint (for Kubernetes)
	s.Router.GET("/ready", s.readinessHandler)

	// Metrics endpoint (placeholder for Prometheus metrics)
	s.Router.GET("/metrics", s.metricsHandler)
}

// healthHandler handles health check requests
func (s *BaseService) healthHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	health := map[string]interface{}{
		"service": s.Name,
		"status":  "healthy",
		"time":    time.Now().UTC(),
		"checks":  make(map[string]string),
	}

	allHealthy := true
	for _, check := range s.HealthChecks {
		if err := check.Check(ctx); err != nil {
			health["checks"].(map[string]string)[check.Name] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			health["checks"].(map[string]string)[check.Name] = "healthy"
		}
	}

	if !allHealthy {
		health["status"] = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	c.JSON(http.StatusOK, health)
}

// readinessHandler handles readiness check requests
func (s *BaseService) readinessHandler(c *gin.Context) {
	// For now, same as health check
	// In production, this might check if service is ready to receive traffic
	s.healthHandler(c)
}

// metricsHandler handles metrics requests (placeholder)
func (s *BaseService) metricsHandler(c *gin.Context) {
	// TODO: Implement Prometheus metrics
	c.JSON(http.StatusOK, gin.H{
		"service": s.Name,
		"metrics": "TODO: Implement Prometheus metrics",
	})
}

// Start starts the service with graceful shutdown
func (s *BaseService) Start() error {
	// Start server in a goroutine
	go func() {
		s.Logger.WithFields(map[string]interface{}{
			"port":    s.Config.ServicePort,
			"service": s.Name,
		}).Info("Starting HTTP server")

		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.Logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		s.Logger.WithError(err).Error("Server forced to shutdown")
		return err
	}

	s.Logger.Info("Server exited")
	return nil
}

// Stop gracefully stops the service
func (s *BaseService) Stop() error {
	var errors []error

	// Stop HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		errors = append(errors, fmt.Errorf("failed to shutdown HTTP server: %w", err))
	}

	// Close database connections
	if s.PostgresDB != nil {
		if err := s.PostgresDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close PostgreSQL connection: %w", err))
		}
	}

	if s.MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.MongoDB.Close(ctx); err != nil {
			errors = append(errors, fmt.Errorf("failed to close MongoDB connection: %w", err))
		}
	}

	if s.RedisClient != nil {
		if err := s.RedisClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close Redis connection: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	return nil
}
