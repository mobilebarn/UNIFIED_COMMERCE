package analytics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"retail-os/services/analytics/internal/models"
	"retail-os/services/analytics/internal/service"
	"retail-os/services/shared/logger"
)

// Handler handles HTTP requests for analytics
type Handler struct {
	service *service.AnalyticsService
	log     *logger.Logger
}

// NewHandler creates a new analytics handler
func NewHandler(service *service.AnalyticsService, log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// RegisterRoutes registers the analytics routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Customer behavior tracking
	router.POST("/api/v1/analytics/behavior", h.trackCustomerBehavior)
	router.GET("/api/v1/analytics/behavior/:customerID", h.getCustomerBehaviors)

	// Product recommendations
	router.POST("/api/v1/analytics/recommendations/generate/:customerID", h.generateProductRecommendations)
	router.GET("/api/v1/analytics/recommendations/:customerID", h.getProductRecommendations)

	// Customer segments
	router.POST("/api/v1/analytics/segments", h.createCustomerSegment)
	router.GET("/api/v1/analytics/segments/:id", h.getCustomerSegment)

	// Business metrics
	router.POST("/api/v1/analytics/metrics", h.trackBusinessMetric)
	router.GET("/api/v1/analytics/metrics/:name", h.getBusinessMetrics)
}

// trackCustomerBehavior records a customer behavior event
func (h *Handler) trackCustomerBehavior(c *gin.Context) {
	var behavior models.CustomerBehavior
	if err := c.ShouldBindJSON(&behavior); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.TrackCustomerBehavior(c.Request.Context(), &behavior); err != nil {
		h.log.WithError(err).Error("Failed to track customer behavior")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track customer behavior"})
		return
	}

	c.JSON(http.StatusCreated, behavior)
}

// getCustomerBehaviors retrieves customer behaviors
func (h *Handler) getCustomerBehaviors(c *gin.Context) {
	customerID := c.Param("customerID")

	// Get limit from query parameter, default to 100
	limit := 100
	if limitParam := c.Query("limit"); limitParam != "" {
		// Parse limit parameter
		// In a real implementation, you would parse this properly
	}

	behaviors, err := h.service.GetCustomerBehaviors(c.Request.Context(), customerID, limit)
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve customer behaviors")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer behaviors"})
		return
	}

	c.JSON(http.StatusOK, behaviors)
}

// generateProductRecommendations generates product recommendations for a customer
func (h *Handler) generateProductRecommendations(c *gin.Context) {
	customerID := c.Param("customerID")

	// Get algorithm from query parameter, default to popularity
	algorithm := c.Query("algorithm")
	if algorithm == "" {
		algorithm = "popularity"
	}

	// Get limit from query parameter, default to 10
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		// In a real implementation, you would parse this properly
	}

	recommendations, err := h.service.GenerateRecommendations(c.Request.Context(), customerID, service.RecommendationAlgorithm(algorithm), limit)
	if err != nil {
		h.log.WithError(err).Error("Failed to generate product recommendations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate product recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// getProductRecommendations retrieves product recommendations for a customer
func (h *Handler) getProductRecommendations(c *gin.Context) {
	customerID := c.Param("customerID")

	// Get limit from query parameter, default to 10
	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		// Parse limit parameter
		// In a real implementation, you would parse this properly
	}

	recommendations, err := h.service.GetProductRecommendations(c.Request.Context(), customerID, limit)
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve product recommendations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product recommendations"})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// createCustomerSegment creates a new customer segment
func (h *Handler) createCustomerSegment(c *gin.Context) {
	var segment models.CustomerSegment
	if err := c.ShouldBindJSON(&segment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCustomerSegment(c.Request.Context(), &segment); err != nil {
		h.log.WithError(err).Error("Failed to create customer segment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer segment"})
		return
	}

	c.JSON(http.StatusCreated, segment)
}

// getCustomerSegment retrieves a customer segment by ID
func (h *Handler) getCustomerSegment(c *gin.Context) {
	id := c.Param("id")

	segment, err := h.service.GetCustomerSegment(c.Request.Context(), id)
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve customer segment")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer segment"})
		return
	}

	c.JSON(http.StatusOK, segment)
}

// trackBusinessMetric records a business metric
func (h *Handler) trackBusinessMetric(c *gin.Context) {
	var metric models.BusinessMetric
	if err := c.ShouldBindJSON(&metric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.TrackBusinessMetric(c.Request.Context(), &metric); err != nil {
		h.log.WithError(err).Error("Failed to track business metric")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track business metric"})
		return
	}

	c.JSON(http.StatusCreated, metric)
}

// getBusinessMetrics retrieves business metrics for a time period
func (h *Handler) getBusinessMetrics(c *gin.Context) {
	name := c.Param("name")

	// Get start and end times from query parameters
	// In a real implementation, you would parse these properly
	start := time.Now().AddDate(0, 0, -30) // Default to last 30 days
	end := time.Now()

	metrics, err := h.service.GetBusinessMetrics(c.Request.Context(), name, start, end)
	if err != nil {
		h.log.WithError(err).Error("Failed to retrieve business metrics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve business metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
