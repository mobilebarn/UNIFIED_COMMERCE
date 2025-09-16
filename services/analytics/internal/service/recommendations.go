package service

import (
	"context"
	"math/rand"
	"time"

	"retail-os/services/analytics/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecommendationAlgorithm represents a recommendation algorithm
type RecommendationAlgorithm string

const (
	// PopularityBased recommends products based on overall popularity
	PopularityBased RecommendationAlgorithm = "popularity"

	// CollaborativeFiltering recommends products based on similar users
	CollaborativeFiltering RecommendationAlgorithm = "collaborative"

	// ContentBased recommends products based on product attributes
	ContentBased RecommendationAlgorithm = "content"

	// TrendingBased recommends products that are currently trending
	TrendingBased RecommendationAlgorithm = "trending"

	// RandomBased recommends random products (for testing)
	RandomBased RecommendationAlgorithm = "random"
)

// GenerateRecommendations generates product recommendations using the specified algorithm
func (s *AnalyticsService) GenerateRecommendations(ctx context.Context, customerID string, algorithm RecommendationAlgorithm, limit int) ([]*models.ProductRecommendation, error) {
	switch algorithm {
	case PopularityBased:
		return s.generatePopularityBasedRecommendations(ctx, customerID, limit)
	case CollaborativeFiltering:
		return s.generateCollaborativeFilteringRecommendations(ctx, customerID, limit)
	case ContentBased:
		return s.generateContentBasedRecommendations(ctx, customerID, limit)
	case TrendingBased:
		return s.generateTrendingBasedRecommendations(ctx, customerID, limit)
	case RandomBased:
		return s.generateRandomRecommendations(ctx, customerID, limit)
	default:
		// Default to popularity-based recommendations
		return s.generatePopularityBasedRecommendations(ctx, customerID, limit)
	}
}

// generatePopularityBasedRecommendations generates recommendations based on product popularity
func (s *AnalyticsService) generatePopularityBasedRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	// In a real implementation, this would query a popularity ranking system
	// For now, we'll generate mock recommendations with popularity scores

	recommendations := []*models.ProductRecommendation{}

	// Sample popular products (in a real system, these would come from analytics data)
	popularProducts := []struct {
		ID    string
		Score float64
	}{
		{"product-001", 0.95},
		{"product-002", 0.92},
		{"product-003", 0.89},
		{"product-004", 0.87},
		{"product-005", 0.85},
		{"product-006", 0.82},
		{"product-007", 0.80},
		{"product-008", 0.78},
		{"product-009", 0.75},
		{"product-010", 0.72},
	}

	// Limit the number of recommendations
	actualLimit := limit
	if actualLimit <= 0 || actualLimit > len(popularProducts) {
		actualLimit = len(popularProducts)
	}

	now := time.Now()
	expiry := now.Add(24 * time.Hour)

	for i := 0; i < actualLimit; i++ {
		rec := &models.ProductRecommendation{
			ID:                 primitive.NewObjectID(),
			CustomerID:         customerID,
			ProductID:          popularProducts[i].ID,
			Score:              popularProducts[i].Score,
			RecommendationType: string(PopularityBased),
			CreatedAt:          now,
			ExpiresAt:          expiry,
		}

		recommendations = append(recommendations, rec)

		// Save recommendation to database
		if err := s.repo.SaveProductRecommendation(ctx, rec); err != nil {
			s.log.WithError(err).Error("Failed to save product recommendation")
		}
	}

	return recommendations, nil
}

// generateCollaborativeFilteringRecommendations generates recommendations based on similar users
func (s *AnalyticsService) generateCollaborativeFilteringRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	// In a real implementation, this would use collaborative filtering algorithms
	// For now, we'll generate mock recommendations based on customer segments

	recommendations := []*models.ProductRecommendation{}

	// Sample products for collaborative filtering (in a real system, these would come from ML models)
	collaborativeProducts := []struct {
		ID    string
		Score float64
	}{
		{"product-101", 0.88},
		{"product-102", 0.85},
		{"product-103", 0.82},
		{"product-104", 0.79},
		{"product-105", 0.76},
		{"product-106", 0.73},
		{"product-107", 0.70},
		{"product-108", 0.68},
		{"product-109", 0.65},
		{"product-110", 0.62},
	}

	// Limit the number of recommendations
	actualLimit := limit
	if actualLimit <= 0 || actualLimit > len(collaborativeProducts) {
		actualLimit = len(collaborativeProducts)
	}

	now := time.Now()
	expiry := now.Add(24 * time.Hour)

	for i := 0; i < actualLimit; i++ {
		rec := &models.ProductRecommendation{
			ID:                 primitive.NewObjectID(),
			CustomerID:         customerID,
			ProductID:          collaborativeProducts[i].ID,
			Score:              collaborativeProducts[i].Score,
			RecommendationType: string(CollaborativeFiltering),
			CreatedAt:          now,
			ExpiresAt:          expiry,
		}

		recommendations = append(recommendations, rec)

		// Save recommendation to database
		if err := s.repo.SaveProductRecommendation(ctx, rec); err != nil {
			s.log.WithError(err).Error("Failed to save product recommendation")
		}
	}

	return recommendations, nil
}

// generateContentBasedRecommendations generates recommendations based on product content
func (s *AnalyticsService) generateContentBasedRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	// In a real implementation, this would use content-based filtering algorithms
	// For now, we'll generate mock recommendations based on product attributes

	recommendations := []*models.ProductRecommendation{}

	// Sample products for content-based filtering (in a real system, these would come from product data)
	contentProducts := []struct {
		ID    string
		Score float64
	}{
		{"product-201", 0.91},
		{"product-202", 0.88},
		{"product-203", 0.85},
		{"product-204", 0.82},
		{"product-205", 0.79},
		{"product-206", 0.76},
		{"product-207", 0.73},
		{"product-208", 0.70},
		{"product-209", 0.67},
		{"product-210", 0.64},
	}

	// Limit the number of recommendations
	actualLimit := limit
	if actualLimit <= 0 || actualLimit > len(contentProducts) {
		actualLimit = len(contentProducts)
	}

	now := time.Now()
	expiry := now.Add(24 * time.Hour)

	for i := 0; i < actualLimit; i++ {
		rec := &models.ProductRecommendation{
			ID:                 primitive.NewObjectID(),
			CustomerID:         customerID,
			ProductID:          contentProducts[i].ID,
			Score:              contentProducts[i].Score,
			RecommendationType: string(ContentBased),
			CreatedAt:          now,
			ExpiresAt:          expiry,
		}

		recommendations = append(recommendations, rec)

		// Save recommendation to database
		if err := s.repo.SaveProductRecommendation(ctx, rec); err != nil {
			s.log.WithError(err).Error("Failed to save product recommendation")
		}
	}

	return recommendations, nil
}

// generateTrendingBasedRecommendations generates recommendations based on trending products
func (s *AnalyticsService) generateTrendingBasedRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	// In a real implementation, this would query trending product data
	// For now, we'll generate mock recommendations based on recent popularity

	recommendations := []*models.ProductRecommendation{}

	// Sample trending products (in a real system, these would come from recent analytics data)
	trendingProducts := []struct {
		ID    string
		Score float64
	}{
		{"product-301", 0.93},
		{"product-302", 0.90},
		{"product-303", 0.87},
		{"product-304", 0.84},
		{"product-305", 0.81},
		{"product-306", 0.78},
		{"product-307", 0.75},
		{"product-308", 0.72},
		{"product-309", 0.69},
		{"product-310", 0.66},
	}

	// Limit the number of recommendations
	actualLimit := limit
	if actualLimit <= 0 || actualLimit > len(trendingProducts) {
		actualLimit = len(trendingProducts)
	}

	now := time.Now()
	expiry := now.Add(12 * time.Hour) // Trending recommendations expire faster

	for i := 0; i < actualLimit; i++ {
		rec := &models.ProductRecommendation{
			ID:                 primitive.NewObjectID(),
			CustomerID:         customerID,
			ProductID:          trendingProducts[i].ID,
			Score:              trendingProducts[i].Score,
			RecommendationType: string(TrendingBased),
			CreatedAt:          now,
			ExpiresAt:          expiry,
		}

		recommendations = append(recommendations, rec)

		// Save recommendation to database
		if err := s.repo.SaveProductRecommendation(ctx, rec); err != nil {
			s.log.WithError(err).Error("Failed to save product recommendation")
		}
	}

	return recommendations, nil
}

// generateRandomRecommendations generates random recommendations (for testing)
func (s *AnalyticsService) generateRandomRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	recommendations := []*models.ProductRecommendation{}

	// Generate random product IDs
	rand.Seed(time.Now().UnixNano())

	now := time.Now()
	expiry := now.Add(24 * time.Hour)

	for i := 0; i < limit; i++ {
		rec := &models.ProductRecommendation{
			ID:                 primitive.NewObjectID(),
			CustomerID:         customerID,
			ProductID:          "product-" + generateRandomID(5),
			Score:              rand.Float64(),
			RecommendationType: string(RandomBased),
			CreatedAt:          now,
			ExpiresAt:          expiry,
		}

		recommendations = append(recommendations, rec)

		// Save recommendation to database
		if err := s.repo.SaveProductRecommendation(ctx, rec); err != nil {
			s.log.WithError(err).Error("Failed to save product recommendation")
		}
	}

	return recommendations, nil
}

// generateRandomID generates a random ID string of specified length
func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
