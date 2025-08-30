package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unified-commerce/services/payment/graphql"
	"unified-commerce/services/payment/handlers"
	"unified-commerce/services/payment/models"
	"unified-commerce/services/payment/repository"
	"unified-commerce/services/payment/service"
	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("payment")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig(cfg.ServiceName)
	loggerConfig.Level = cfg.LogLevel
	log := logger.NewLogger(loggerConfig)

	// Connect to PostgreSQL
	postgresConfig := database.NewPostgresConfigFromEnv(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)
	db, err := database.NewPostgresConnection(postgresConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to PostgreSQL")
	}

	// Run database migrations
	if err := runMigrations(db.DB); err != nil {
		log.WithError(err).Fatal("Failed to run database migrations")
	}

	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(db.DB, log)

	// Initialize services
	paymentService := service.NewPaymentService(paymentRepo, log)

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService, log)

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		postgresStatus := "healthy"
		if err := db.Health(context.Background()); err != nil {
			postgresStatus = "unhealthy"
		}

		health := map[string]interface{}{
			"service": "payment",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"postgres": postgresStatus,
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	paymentHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(paymentService, log)
	playgroundHandler := graphql.NewGraphQLPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	router.GET("/graphql/playground", gin.WrapH(playgroundHandler))

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("port", cfg.ServicePort).Info("Starting Payment Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Payment Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Payment Service stopped")
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Payment{},
		&models.PaymentMethod{},
		&models.PaymentGateway{},
		&models.Refund{},
		&models.PaymentEvent{},
		&models.Settlement{},
	)
}
