package repository

import (
	"context"
	"time"

	"retail-os/services/analytics/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRepository is a mock implementation of the Repository interface for testing
type MockRepository struct{}

// SaveCustomerBehavior saves customer behavior data
func (m *MockRepository) SaveCustomerBehavior(ctx context.Context, behavior *models.CustomerBehavior) error {
	// For testing, we don't actually save anything
	return nil
}

// GetCustomerBehaviors retrieves customer behavior data
func (m *MockRepository) GetCustomerBehaviors(ctx context.Context, customerID string, limit int) ([]*models.CustomerBehavior, error) {
	// For testing, return empty slice
	return []*models.CustomerBehavior{}, nil
}

// SaveProductRecommendation saves a product recommendation
func (m *MockRepository) SaveProductRecommendation(ctx context.Context, recommendation *models.ProductRecommendation) error {
	// For testing, we don't actually save anything
	return nil
}

// GetProductRecommendations retrieves product recommendations for a customer
func (m *MockRepository) GetProductRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	// For testing, return empty slice
	return []*models.ProductRecommendation{}, nil
}

// SaveCustomerSegment saves a customer segment
func (m *MockRepository) SaveCustomerSegment(ctx context.Context, segment *models.CustomerSegment) error {
	// For testing, we don't actually save anything
	return nil
}

// GetCustomerSegment retrieves a customer segment by ID
func (m *MockRepository) GetCustomerSegment(ctx context.Context, segmentID string) (*models.CustomerSegment, error) {
	// For testing, return a mock segment
	return &models.CustomerSegment{
		ID:          primitive.NewObjectID(),
		Name:        "Test Segment",
		Description: "Test segment for testing",
		CustomerIDs: []string{"test-customer-123"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// SaveBusinessMetric saves a business metric
func (m *MockRepository) SaveBusinessMetric(ctx context.Context, metric *models.BusinessMetric) error {
	// For testing, we don't actually save anything
	return nil
}

// GetBusinessMetrics retrieves business metrics
func (m *MockRepository) GetBusinessMetrics(ctx context.Context, name string, start, end time.Time) ([]*models.BusinessMetric, error) {
	// For testing, return empty slice
	return []*models.BusinessMetric{}, nil
}

// SaveAnalyticsReport saves an analytics report
func (m *MockRepository) SaveAnalyticsReport(ctx context.Context, report *models.AnalyticsReport) error {
	// For testing, we don't actually save anything
	return nil
}

// GetAnalyticsReport retrieves an analytics report by ID
func (m *MockRepository) GetAnalyticsReport(ctx context.Context, reportID string) (*models.AnalyticsReport, error) {
	// For testing, return nil
	return nil, nil
}
