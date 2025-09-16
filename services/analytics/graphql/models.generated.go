package graphql

import (
	"retail-os/services/analytics/internal/models"
)

// CustomerBehavior model
type CustomerBehavior struct {
	ID          string                  `json:"id"`
	CustomerID  string                  `json:"customerId"`
	SessionID   string                  `json:"sessionId"`
	Action      CustomerAction          `json:"action"`
	EntityID    string                  `json:"entityId"`
	EntityType  string                  `json:"entityType"`
	Timestamp   string                  `json:"timestamp"`
	UserAgent   *string                 `json:"userAgent,omitempty"`
	IPAddress   *string                 `json:"ipAddress,omitempty"`
	Referrer    *string                 `json:"referrer,omitempty"`
	UTMSource   *string                 `json:"utmSource,omitempty"`
	UTMMedium   *string                 `json:"utmMedium,omitempty"`
	UTMCampaign *string                 `json:"utmCampaign,omitempty"`
	CustomData  *map[string]interface{} `json:"customData,omitempty"`
}

// ProductRecommendation model
type ProductRecommendation struct {
	ID                 string             `json:"id"`
	CustomerID         string             `json:"customerId"`
	ProductID          string             `json:"productId"`
	Product            *Product           `json:"product,omitempty"`
	Score              float64            `json:"score"`
	RecommendationType RecommendationType `json:"recommendationType"`
	CreatedAt          string             `json:"createdAt"`
	ExpiresAt          *string            `json:"expiresAt,omitempty"`
}

// CustomerSegment model
type CustomerSegment struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	CustomerIDs []string `json:"customerIds"`
	Customers   []User   `json:"customers"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

// BusinessMetric model
type BusinessMetric struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Value     float64            `json:"value"`
	Dimension *string            `json:"dimension,omitempty"`
	Timestamp string             `json:"timestamp"`
	Tags      *map[string]string `json:"tags,omitempty"`
}

// AnalyticsReport model
type AnalyticsReport struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Data        interface{} `json:"data"`
	GeneratedAt string      `json:"generatedAt"`
	PeriodStart string      `json:"periodStart"`
	PeriodEnd   string      `json:"periodEnd"`
}

// Input types
type TrackCustomerBehaviorInput struct {
	CustomerID  string                  `json:"customerId"`
	SessionID   string                  `json:"sessionId"`
	Action      CustomerAction          `json:"action"`
	EntityID    string                  `json:"entityId"`
	EntityType  string                  `json:"entityType"`
	UserAgent   *string                 `json:"userAgent,omitempty"`
	IPAddress   *string                 `json:"ipAddress,omitempty"`
	Referrer    *string                 `json:"referrer,omitempty"`
	UTMSource   *string                 `json:"utmSource,omitempty"`
	UTMMedium   *string                 `json:"utmMedium,omitempty"`
	UTMCampaign *string                 `json:"utmCampaign,omitempty"`
	CustomData  *map[string]interface{} `json:"customData,omitempty"`
}

type GenerateRecommendationsInput struct {
	CustomerID string             `json:"customerId"`
	Algorithm  RecommendationType `json:"algorithm"`
	Limit      *int               `json:"limit,omitempty"`
}

type CreateCustomerSegmentInput struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	CustomerIDs []string `json:"customerIds"`
}

type TrackBusinessMetricInput struct {
	Name      string             `json:"name"`
	Value     float64            `json:"value"`
	Dimension *string            `json:"dimension,omitempty"`
	Tags      *map[string]string `json:"tags,omitempty"`
}

type GenerateAnalyticsReportInput struct {
	Type        AnalyticsReportType `json:"type"`
	PeriodStart string              `json:"periodStart"`
	PeriodEnd   string              `json:"periodEnd"`
}

// Convert internal models to GraphQL models
func customerBehaviorToGraphQL(cb *models.CustomerBehavior) *CustomerBehavior {
	return &CustomerBehavior{
		ID:          cb.ID.Hex(),
		CustomerID:  cb.CustomerID,
		SessionID:   cb.SessionID,
		Action:      CustomerAction(cb.Action),
		EntityID:    cb.EntityID,
		EntityType:  cb.EntityType,
		Timestamp:   cb.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		UserAgent:   &cb.UserAgent,
		IPAddress:   &cb.IPAddress,
		Referrer:    &cb.Referrer,
		UTMSource:   &cb.UTMSource,
		UTMMedium:   &cb.UTMMedium,
		UTMCampaign: &cb.UTMCampaign,
	}
}

func productRecommendationToGraphQL(pr *models.ProductRecommendation) *ProductRecommendation {
	return &ProductRecommendation{
		ID:                 pr.ID.Hex(),
		CustomerID:         pr.CustomerID,
		ProductID:          pr.ProductID,
		Score:              pr.Score,
		RecommendationType: RecommendationType(pr.RecommendationType),
		CreatedAt:          pr.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ExpiresAt:          func() *string { s := pr.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"); return &s }(),
	}
}
