package service

import (
	"context"
	"time"

	"retail-os/services/analytics/internal/models"
	"retail-os/services/analytics/internal/repository"
	"retail-os/services/shared/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AnalyticsService provides analytics business logic
type AnalyticsService struct {
	repo repository.RepositoryInterface
	log  *logger.Logger
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo repository.RepositoryInterface, log *logger.Logger) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
		log:  log,
	}
}

// TrackCustomerBehavior records a customer behavior event
func (s *AnalyticsService) TrackCustomerBehavior(ctx context.Context, behavior *models.CustomerBehavior) error {
	// Set timestamp if not provided
	if behavior.Timestamp.IsZero() {
		behavior.Timestamp = time.Now()
	}

	// Set ID if not provided
	if behavior.ID.IsZero() {
		behavior.ID = primitive.NewObjectID()
	}

	s.log.WithFields(map[string]interface{}{
		"customer_id": behavior.CustomerID,
		"action":      behavior.Action,
		"entity_type": behavior.EntityType,
		"entity_id":   behavior.EntityID,
	}).Info("Tracking customer behavior")

	return s.repo.SaveCustomerBehavior(ctx, behavior)
}

// GetCustomerBehaviors retrieves recent customer behaviors
func (s *AnalyticsService) GetCustomerBehaviors(ctx context.Context, customerID string, limit int) ([]*models.CustomerBehavior, error) {
	s.log.WithFields(map[string]interface{}{
		"customer_id": customerID,
		"limit":       limit,
	}).Info("Retrieving customer behaviors")

	return s.repo.GetCustomerBehaviors(ctx, customerID, limit)
}

// GenerateProductRecommendations generates product recommendations for a customer
func (s *AnalyticsService) GenerateProductRecommendations(ctx context.Context, customerID string) ([]*models.ProductRecommendation, error) {
	s.log.WithField("customer_id", customerID).Info("Generating product recommendations")

	// Generate popularity-based recommendations as default
	return s.GenerateRecommendations(ctx, customerID, PopularityBased, 10)
}

// GetProductRecommendations retrieves product recommendations for a customer
func (s *AnalyticsService) GetProductRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	s.log.WithFields(map[string]interface{}{
		"customer_id": customerID,
		"limit":       limit,
	}).Info("Retrieving product recommendations")

	return s.repo.GetProductRecommendations(ctx, customerID, limit)
}

// CreateCustomerSegment creates a new customer segment
func (s *AnalyticsService) CreateCustomerSegment(ctx context.Context, segment *models.CustomerSegment) error {
	// Set timestamps
	now := time.Now()
	segment.CreatedAt = now
	segment.UpdatedAt = now

	// Set ID if not provided
	if segment.ID.IsZero() {
		segment.ID = primitive.NewObjectID()
	}

	s.log.WithFields(map[string]interface{}{
		"segment_name": segment.Name,
		"customer_ids": len(segment.CustomerIDs),
	}).Info("Creating customer segment")

	return s.repo.SaveCustomerSegment(ctx, segment)
}

// GetCustomerSegment retrieves a customer segment by ID
func (s *AnalyticsService) GetCustomerSegment(ctx context.Context, id string) (*models.CustomerSegment, error) {
	s.log.WithField("segment_id", id).Info("Retrieving customer segment")

	return s.repo.GetCustomerSegment(ctx, id)
}

// TrackBusinessMetric records a business metric
func (s *AnalyticsService) TrackBusinessMetric(ctx context.Context, metric *models.BusinessMetric) error {
	// Set timestamp if not provided
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}

	// Set ID if not provided
	if metric.ID.IsZero() {
		metric.ID = primitive.NewObjectID()
	}

	s.log.WithFields(map[string]interface{}{
		"metric_name": metric.Name,
		"value":       metric.Value,
		"dimension":   metric.Dimension,
	}).Info("Tracking business metric")

	return s.repo.SaveBusinessMetric(ctx, metric)
}

// GetBusinessMetrics retrieves business metrics for a time period
func (s *AnalyticsService) GetBusinessMetrics(ctx context.Context, name string, start, end time.Time) ([]*models.BusinessMetric, error) {
	s.log.WithFields(map[string]interface{}{
		"metric_name": name,
		"start":       start,
		"end":         end,
	}).Info("Retrieving business metrics")

	return s.repo.GetBusinessMetrics(ctx, name, start, end)
}
