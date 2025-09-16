package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CustomerBehavior represents customer interaction data
type CustomerBehavior struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	CustomerID  string                 `bson:"customer_id" json:"customer_id"`
	SessionID   string                 `bson:"session_id" json:"session_id"`
	Action      string                 `bson:"action" json:"action"` // view, click, add_to_cart, purchase, etc.
	EntityID    string                 `bson:"entity_id" json:"entity_id"`
	EntityType  string                 `bson:"entity_type" json:"entity_type"` // product, category, brand, etc.
	Timestamp   time.Time              `bson:"timestamp" json:"timestamp"`
	UserAgent   string                 `bson:"user_agent" json:"user_agent"`
	IPAddress   string                 `bson:"ip_address" json:"ip_address"`
	Referrer    string                 `bson:"referrer" json:"referrer"`
	UTMSource   string                 `bson:"utm_source" json:"utm_source"`
	UTMMedium   string                 `bson:"utm_medium" json:"utm_medium"`
	UTMCampaign string                 `bson:"utm_campaign" json:"utm_campaign"`
	CustomData  map[string]interface{} `bson:"custom_data" json:"custom_data"`
}

// ProductRecommendation represents a product recommendation for a customer
type ProductRecommendation struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID         string             `bson:"customer_id" json:"customer_id"`
	ProductID          string             `bson:"product_id" json:"product_id"`
	Score              float64            `bson:"score" json:"score"`
	RecommendationType string             `bson:"recommendation_type" json:"recommendation_type"` // collaborative, content, popularity, etc.
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt          time.Time          `bson:"expires_at" json:"expires_at"`
}

// CustomerSegment represents a group of customers with similar characteristics
type CustomerSegment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CustomerIDs []string           `bson:"customer_ids" json:"customer_ids"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// BusinessMetric represents a key business metric
type BusinessMetric struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Value     float64            `bson:"value" json:"value"`
	Dimension string             `bson:"dimension" json:"dimension"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Tags      map[string]string  `bson:"tags" json:"tags"`
}

// AnalyticsReport represents a generated analytics report
type AnalyticsReport struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Type        string             `bson:"type" json:"type"`
	Data        interface{}        `bson:"data" json:"data"`
	GeneratedAt time.Time          `bson:"generated_at" json:"generated_at"`
	PeriodStart time.Time          `bson:"period_start" json:"period_start"`
	PeriodEnd   time.Time          `bson:"period_end" json:"period_end"`
}
