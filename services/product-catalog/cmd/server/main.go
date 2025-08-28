package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"unified-commerce/services/product-catalog/handlers"
	"unified-commerce/services/product-catalog/repository"
	"unified-commerce/services/product-catalog/service"
	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.Environment)

	// Connect to MongoDB
	mongoClient, err := connectMongoDB(cfg.MongoDB.URI)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	defer mongoClient.Disconnect(context.Background())

	// Initialize database
	db := mongoClient.Database(cfg.MongoDB.Database)

	// Initialize repositories
	productRepo := repository.NewProductRepository(db, log)
	categoryRepo := repository.NewCategoryRepository(db, log)
	collectionRepo := repository.NewCollectionRepository(db, log)

	// Initialize services
	productService := service.NewProductService(productRepo, categoryRepo, collectionRepo, log)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService, log)

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
	router.Use(middleware.RateLimit())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		health := map[string]interface{}{
			"service": "product-catalog",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"mongodb": checkMongoDB(db),
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	productHandler.RegisterRoutes(router)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("port", cfg.Server.Port).Info("Starting Product Catalog Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Product Catalog Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Product Catalog Service stopped")
}

// connectMongoDB establishes connection to MongoDB
func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// checkMongoDB checks MongoDB connection health
func checkMongoDB(db *mongo.Database) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Client().Ping(ctx, nil)
	if err != nil {
		return "unhealthy"
	}
	return "healthy"
}
