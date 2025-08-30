package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/product-catalog/graphql"
	"unified-commerce/services/product-catalog/handlers"
	"unified-commerce/services/product-catalog/repository"
	"unified-commerce/services/product-catalog/service"
	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("product-catalog")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig("product-catalog")
	loggerConfig.Level = cfg.LogLevel
	log := logger.NewLogger(loggerConfig)

	// Connect to MongoDB
	mongoConfig := database.NewMongoConfigFromEnv(
		cfg.MongoURL,
		cfg.MongoDatabase,
		cfg.MongoUser,
		cfg.MongoPassword,
	)
	mongoDB, err := database.NewMongoConnection(mongoConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	defer mongoDB.Close(context.Background())

	// Initialize repositories
	repo := repository.NewRepository(mongoDB)

	// Initialize services
	productService := service.NewProductService(repo, log)

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

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		mongoStatus := "healthy"
		if err := mongoDB.Health(context.Background()); err != nil {
			mongoStatus = "unhealthy"
		}

		health := map[string]interface{}{
			"service": "product-catalog",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"mongodb": mongoStatus,
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	productHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(productService, log)
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
		log.WithFields(map[string]interface{}{"port": cfg.ServicePort}).Info("Starting Product Catalog Service")
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
