package repository

import (
	"context"
	"time"

	"retail-os/services/analytics/internal/models"
	"retail-os/services/shared/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RepositoryInterface defines the interface for the analytics repository
type RepositoryInterface interface {
	SaveCustomerBehavior(ctx context.Context, behavior *models.CustomerBehavior) error
	GetCustomerBehaviors(ctx context.Context, customerID string, limit int) ([]*models.CustomerBehavior, error)
	SaveProductRecommendation(ctx context.Context, recommendation *models.ProductRecommendation) error
	GetProductRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error)
	SaveCustomerSegment(ctx context.Context, segment *models.CustomerSegment) error
	GetCustomerSegment(ctx context.Context, id string) (*models.CustomerSegment, error)
	SaveBusinessMetric(ctx context.Context, metric *models.BusinessMetric) error
	GetBusinessMetrics(ctx context.Context, name string, start, end time.Time) ([]*models.BusinessMetric, error)
	SaveAnalyticsReport(ctx context.Context, report *models.AnalyticsReport) error
	GetAnalyticsReport(ctx context.Context, id string) (*models.AnalyticsReport, error)
}

// Repository handles data access for analytics
type Repository struct {
	mongoDB    *database.MongoDB
	postgresDB *database.PostgresDB
}

// NewRepository creates a new analytics repository
func NewRepository(mongoDB *database.MongoDB, postgresDB *database.PostgresDB) *Repository {
	return &Repository{
		mongoDB:    mongoDB,
		postgresDB: postgresDB,
	}
}

// SaveCustomerBehavior saves customer behavior data
func (r *Repository) SaveCustomerBehavior(ctx context.Context, behavior *models.CustomerBehavior) error {
	collection := r.mongoDB.Client.Database("analytics").Collection("customer_behaviors")
	_, err := collection.InsertOne(ctx, behavior)
	return err
}

// GetCustomerBehaviors retrieves customer behaviors for a specific customer
func (r *Repository) GetCustomerBehaviors(ctx context.Context, customerID string, limit int) ([]*models.CustomerBehavior, error) {
	collection := r.mongoDB.Client.Database("analytics").Collection("customer_behaviors")
	filter := bson.M{"customer_id": customerID}

	// Sort by timestamp descending
	opts := &options.FindOptions{}
	opts.SetSort(bson.D{{"timestamp", -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var behaviors []*models.CustomerBehavior
	if err = cursor.All(ctx, &behaviors); err != nil {
		return nil, err
	}

	return behaviors, nil
}

// SaveProductRecommendation saves a product recommendation
func (r *Repository) SaveProductRecommendation(ctx context.Context, recommendation *models.ProductRecommendation) error {
	collection := r.mongoDB.Client.Database("analytics").Collection("product_recommendations")
	_, err := collection.InsertOne(ctx, recommendation)
	return err
}

// GetProductRecommendations retrieves product recommendations for a customer
func (r *Repository) GetProductRecommendations(ctx context.Context, customerID string, limit int) ([]*models.ProductRecommendation, error) {
	collection := r.mongoDB.Client.Database("analytics").Collection("product_recommendations")

	// Filter by customer ID and ensure recommendations haven't expired
	now := time.Now()
	filter := bson.M{
		"customer_id": customerID,
		"expires_at":  bson.M{"$gte": now},
	}

	// Sort by score descending
	opts := &options.FindOptions{}
	opts.SetSort(bson.D{{"score", -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recommendations []*models.ProductRecommendation
	if err = cursor.All(ctx, &recommendations); err != nil {
		return nil, err
	}

	return recommendations, nil
}

// SaveCustomerSegment saves a customer segment
func (r *Repository) SaveCustomerSegment(ctx context.Context, segment *models.CustomerSegment) error {
	collection := r.mongoDB.Client.Database("analytics").Collection("customer_segments")
	_, err := collection.InsertOne(ctx, segment)
	return err
}

// GetCustomerSegment retrieves a customer segment by ID
func (r *Repository) GetCustomerSegment(ctx context.Context, id string) (*models.CustomerSegment, error) {
	collection := r.mongoDB.Client.Database("analytics").Collection("customer_segments")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var segment models.CustomerSegment
	err = collection.FindOne(ctx, filter).Decode(&segment)
	if err != nil {
		return nil, err
	}

	return &segment, nil
}

// SaveBusinessMetric saves a business metric
func (r *Repository) SaveBusinessMetric(ctx context.Context, metric *models.BusinessMetric) error {
	collection := r.mongoDB.Client.Database("analytics").Collection("business_metrics")
	_, err := collection.InsertOne(ctx, metric)
	return err
}

// GetBusinessMetrics retrieves business metrics for a time period
func (r *Repository) GetBusinessMetrics(ctx context.Context, name string, start, end time.Time) ([]*models.BusinessMetric, error) {
	collection := r.mongoDB.Client.Database("analytics").Collection("business_metrics")

	filter := bson.M{
		"name": name,
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	// Sort by timestamp ascending
	opts := &options.FindOptions{}
	opts.SetSort(bson.D{{"timestamp", 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metrics []*models.BusinessMetric
	if err = cursor.All(ctx, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

// SaveAnalyticsReport saves an analytics report
func (r *Repository) SaveAnalyticsReport(ctx context.Context, report *models.AnalyticsReport) error {
	collection := r.mongoDB.Client.Database("analytics").Collection("analytics_reports")
	_, err := collection.InsertOne(ctx, report)
	return err
}

// GetAnalyticsReport retrieves an analytics report by ID
func (r *Repository) GetAnalyticsReport(ctx context.Context, id string) (*models.AnalyticsReport, error) {
	collection := r.mongoDB.Client.Database("analytics").Collection("analytics_reports")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var report models.AnalyticsReport
	err = collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}
