package service

import (
	"context"
	"retail-os/services/analytics/internal/repository"
	"retail-os/services/shared/logger"
	"testing"
)

func TestAnalyticsService_TrackCustomerBehavior(t *testing.T) {
	// Create a mock repository
	mockRepo := &repository.MockRepository{}

	// Create a logger
	log := logger.NewLogger(logger.DefaultConfig("analytics-test"))

	// Create the service
	service := NewAnalyticsService(mockRepo, log)

	// Check that the service was created successfully
	if service == nil {
		t.Error("Expected service to be created, but got nil")
	}

	// Note: We would add more comprehensive tests here once we have a proper mock repository
}

func TestAnalyticsService_GenerateRecommendations(t *testing.T) {
	// Create a mock repository
	mockRepo := &repository.MockRepository{}

	// Create a logger
	log := logger.NewLogger(logger.DefaultConfig("analytics-test"))

	// Create the service
	service := NewAnalyticsService(mockRepo, log)

	// Check that the service was created successfully
	if service == nil {
		t.Error("Expected service to be created, but got nil")
	}

	// Test the recommendation algorithms
	customerID := "test-customer-123"
	ctx := context.Background()

	// Test popularity-based recommendations
	recommendations, err := service.generatePopularityBasedRecommendations(ctx, customerID, 5)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if len(recommendations) > 5 {
		t.Errorf("Expected at most 5 recommendations, but got %d", len(recommendations))
	}

	// Test trending-based recommendations
	recommendations, err = service.generateTrendingBasedRecommendations(ctx, customerID, 5)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if len(recommendations) > 5 {
		t.Errorf("Expected at most 5 recommendations, but got %d", len(recommendations))
	}

	// Test random recommendations
	recommendations, err = service.generateRandomRecommendations(ctx, customerID, 5)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if len(recommendations) > 5 {
		t.Errorf("Expected at most 5 recommendations, but got %d", len(recommendations))
	}
}
