package graphql

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/merchant-account/models"
	"unified-commerce/services/merchant-account/repository"
	merchantService "unified-commerce/services/merchant-account/service"
	"unified-commerce/services/shared/logger"
)

// Handler handles GraphQL requests for the merchant-account service
type Handler struct {
	repo    *repository.Repository
	service *merchantService.MerchantService
	logger  *logger.Logger
}

// NewHandler creates a new GraphQL handler
func NewHandler(repo *repository.Repository, svc *merchantService.MerchantService, log *logger.Logger) *Handler {
	return &Handler{
		repo:    repo,
		service: svc,
		logger:  log,
	}
}

// GraphQLHandler handles GraphQL HTTP requests
func (h *Handler) GraphQLHandler(c *gin.Context) {
	var request GraphQLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Failed to parse GraphQL request")
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": []map[string]interface{}{
				{
					"message": "Invalid GraphQL request format",
					"extensions": map[string]interface{}{
						"code": "INVALID_REQUEST",
					},
				},
			},
		})
		return
	}

	// Get user context from headers (set by gateway)
	userID := c.GetHeader("user-id")
	userEmail := c.GetHeader("user-email")
	userRoles := c.GetHeader("user-roles")

	// Create GraphQL context
	ctx := context.WithValue(c.Request.Context(), "user_id", userID)
	ctx = context.WithValue(ctx, "user_email", userEmail)
	ctx = context.WithValue(ctx, "user_roles", userRoles)

	// Execute GraphQL query
	response := h.executeQuery(ctx, request)

	c.JSON(http.StatusOK, response)
}

// GraphQLRequest represents a GraphQL HTTP request
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

// GraphQLResponse represents a GraphQL HTTP response
type GraphQLResponse struct {
	Data   interface{}              `json:"data,omitempty"`
	Errors []map[string]interface{} `json:"errors,omitempty"`
}

// executeQuery processes the GraphQL query and returns a response
func (h *Handler) executeQuery(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	h.logger.WithFields(map[string]interface{}{
		"operation": request.OperationName,
		"query":     request.Query,
	}).Info("Executing GraphQL query")

	// Parse query type (simple parsing for demo - in production use proper GraphQL parser)
	if contains(request.Query, "merchant(") {
		return h.handleMerchantQuery(ctx, request)
	} else if contains(request.Query, "merchants") {
		return h.handleMerchantsQuery(ctx, request)
	} else if contains(request.Query, "myMerchants") {
		return h.handleMyMerchantsQuery(ctx, request)
	} else if contains(request.Query, "store(") {
		return h.handleStoreQuery(ctx, request)
	} else if contains(request.Query, "stores") {
		return h.handleStoresQuery(ctx, request)
	} else if contains(request.Query, "subscription(") {
		return h.handleSubscriptionQuery(ctx, request)
	} else if contains(request.Query, "createMerchant") {
		return h.handleCreateMerchant(ctx, request)
	} else if contains(request.Query, "updateMerchant") {
		return h.handleUpdateMerchant(ctx, request)
	} else if contains(request.Query, "createStore") {
		return h.handleCreateStore(ctx, request)
	} else if contains(request.Query, "_entities") {
		return h.handleFederationEntities(ctx, request)
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"__schema": map[string]interface{}{
				"queryType": map[string]interface{}{
					"name": "Query",
				},
			},
		},
	}
}

// handleMerchantQuery handles single merchant queries
func (h *Handler) handleMerchantQuery(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	// Extract merchant ID from variables or query
	merchantID, ok := request.Variables["id"].(string)
	if !ok {
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Missing merchant ID"},
			},
		}
	}

	merchant, err := h.repo.Merchant.GetByID(ctx, merchantID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to fetch merchant")
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Failed to fetch merchant"},
			},
		}
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"merchant": h.merchantToGraphQL(merchant),
		},
	}
}

// handleMerchantsQuery handles merchant list queries
func (h *Handler) handleMerchantsQuery(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	// For demo purposes, return sample data
	merchants, _, err := h.repo.Merchant.List(ctx, 0, 100, map[string]interface{}{})
	if err != nil {
		h.logger.WithError(err).Error("Failed to fetch merchants")
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Failed to fetch merchants"},
			},
		}
	}

	var graphqlMerchants []interface{}
	for _, merchant := range merchants {
		graphqlMerchants = append(graphqlMerchants, h.merchantToGraphQL(&merchant))
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"merchants": map[string]interface{}{
				"edges": graphqlMerchants,
				"pageInfo": map[string]interface{}{
					"hasNextPage":     false,
					"hasPreviousPage": false,
					"startCursor":     "",
					"endCursor":       "",
				},
				"totalCount": len(graphqlMerchants),
			},
		},
	}
}

// handleMyMerchantsQuery handles user's merchant queries
func (h *Handler) handleMyMerchantsQuery(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	userID := ctx.Value("user_id").(string)
	if userID == "" {
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Authentication required"},
			},
		}
	}

	// Get merchants for the user
	merchants, err := h.repo.Merchant.GetByUserID(ctx, userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to fetch user merchants")
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Failed to fetch merchants"},
			},
		}
	}

	var graphqlMerchants []interface{}
	for _, merchant := range merchants {
		graphqlMerchants = append(graphqlMerchants, h.merchantToGraphQL(&merchant))
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"myMerchants": graphqlMerchants,
		},
	}
}

// handleStoreQuery handles single store queries
func (h *Handler) handleStoreQuery(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	// For demo - stores are managed through merchant
	return GraphQLResponse{
		Data: map[string]interface{}{
			"store": map[string]interface{}{
				"id":         "demo-store-1",
				"merchantId": "demo-merchant-1",
				"name":       "Demo Store",
				"storeType":  "ONLINE",
				"status":     "ACTIVE",
				"isActive":   true,
			},
		},
	}
}

// handleStoresQuery handles store list queries
func (h *Handler) handleStoresQuery(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	// For demo - return empty stores list
	return GraphQLResponse{
		Data: map[string]interface{}{
			"stores": map[string]interface{}{
				"edges": []interface{}{},
				"pageInfo": map[string]interface{}{
					"hasNextPage":     false,
					"hasPreviousPage": false,
				},
				"totalCount": 0,
			},
		},
	}
}

// handleSubscriptionQuery handles subscription queries
func (h *Handler) handleSubscriptionQuery(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	// For demo - return sample subscription
	return GraphQLResponse{
		Data: map[string]interface{}{
			"subscription": map[string]interface{}{
				"id":         "demo-subscription-1",
				"merchantId": "demo-merchant-1",
				"planId":     "demo-plan-1",
				"status":     "ACTIVE",
				"currency":   "USD",
				"amount":     29.99,
			},
		},
	}
}

// handleCreateMerchant handles merchant creation mutations
func (h *Handler) handleCreateMerchant(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	input, ok := request.Variables["input"].(map[string]interface{})
	if !ok {
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Missing input data"},
			},
		}
	}

	// Convert input to merchant model
	merchant := &models.Merchant{
		BusinessName:   input["businessName"].(string),
		LegalName:      getStringFromMap(input, "legalName"),
		BusinessType:   input["businessType"].(string),
		Industry:       getStringFromMap(input, "industry"),
		TaxID:          getStringFromMap(input, "taxId"),
		WebsiteURL:     getStringFromMap(input, "websiteUrl"),
		Description:    getStringFromMap(input, "description"),
		PrimaryEmail:   input["primaryEmail"].(string),
		PrimaryPhone:   getStringFromMap(input, "primaryPhone"),
		Status:         "pending",
		OnboardingStep: 1,
	}

	createdMerchant, err := h.service.CreateMerchant(ctx, merchant)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create merchant")
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Failed to create merchant"},
			},
		}
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"createMerchant": map[string]interface{}{
				"success":  true,
				"message":  "Merchant created successfully",
				"merchant": h.merchantToGraphQL(createdMerchant),
			},
		},
	}
}

// handleUpdateMerchant handles merchant update mutations
func (h *Handler) handleUpdateMerchant(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	merchantID, ok := request.Variables["id"].(string)
	if !ok {
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Missing merchant ID"},
			},
		}
	}

	input, ok := request.Variables["input"].(map[string]interface{})
	if !ok {
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Missing input data"},
			},
		}
	}

	updatedMerchant, err := h.service.UpdateMerchant(ctx, merchantID, input)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update merchant")
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Failed to update merchant"},
			},
		}
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"updateMerchant": map[string]interface{}{
				"success":  true,
				"message":  "Merchant updated successfully",
				"merchant": h.merchantToGraphQL(updatedMerchant),
			},
		},
	}
}

// handleCreateStore handles store creation mutations
func (h *Handler) handleCreateStore(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	input, ok := request.Variables["input"].(map[string]interface{})
	if !ok {
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Missing input data"},
			},
		}
	}

	store := &models.Store{
		MerchantID:  input["merchantId"].(string),
		Name:        input["name"].(string),
		Description: getStringFromMap(input, "description"),
		StoreType:   input["storeType"].(string),
		Domain:      getStringFromMap(input, "domain"),
		Status:      "draft",
		IsActive:    false,
	}

	createdStore, err := h.service.CreateStore(ctx, store)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create store")
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Failed to create store"},
			},
		}
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"createStore": map[string]interface{}{
				"success": true,
				"message": "Store created successfully",
				"store":   h.storeToGraphQL(createdStore),
			},
		},
	}
}

// handleFederationEntities handles Apollo Federation _entities queries
func (h *Handler) handleFederationEntities(ctx context.Context, request GraphQLRequest) GraphQLResponse {
	representations, ok := request.Variables["representations"].([]interface{})
	if !ok {
		return GraphQLResponse{
			Errors: []map[string]interface{}{
				{"message": "Missing representations"},
			},
		}
	}

	var entities []interface{}
	for _, rep := range representations {
		entity := rep.(map[string]interface{})
		typename := entity["__typename"].(string)
		id := entity["id"].(string)

		switch typename {
		case "Merchant":
			merchant, err := h.repo.GetMerchantByID(ctx, id)
			if err == nil {
				entities = append(entities, h.merchantToGraphQL(merchant))
			} else {
				entities = append(entities, nil)
			}
		case "Store":
			store, err := h.repo.GetStoreByID(ctx, id)
			if err == nil {
				entities = append(entities, h.storeToGraphQL(store))
			} else {
				entities = append(entities, nil)
			}
		case "Subscription":
			subscription, err := h.repo.GetSubscriptionByID(ctx, id)
			if err == nil {
				entities = append(entities, h.subscriptionToGraphQL(subscription))
			} else {
				entities = append(entities, nil)
			}
		default:
			entities = append(entities, nil)
		}
	}

	return GraphQLResponse{
		Data: map[string]interface{}{
			"_entities": entities,
		},
	}
}

// Helper functions for converting models to GraphQL format

func (h *Handler) merchantToGraphQL(merchant *models.Merchant) map[string]interface{} {
	if merchant == nil {
		return nil
	}

	settings, _ := json.Marshal(merchant.Settings)
	var settingsMap map[string]interface{}
	json.Unmarshal(settings, &settingsMap)

	return map[string]interface{}{
		"id":             merchant.ID,
		"businessName":   merchant.BusinessName,
		"legalName":      merchant.LegalName,
		"businessType":   merchant.BusinessType,
		"industry":       merchant.Industry,
		"taxId":          merchant.TaxID,
		"websiteUrl":     merchant.WebsiteURL,
		"description":    merchant.Description,
		"logoUrl":        merchant.LogoURL,
		"primaryEmail":   merchant.PrimaryEmail,
		"primaryPhone":   merchant.PrimaryPhone,
		"status":         merchant.Status,
		"isVerified":     merchant.IsVerified,
		"verifiedAt":     merchant.VerifiedAt,
		"onboardingStep": merchant.OnboardingStep,
		"settings":       settingsMap,
		"createdAt":      merchant.CreatedAt,
		"updatedAt":      merchant.UpdatedAt,
	}
}

func (h *Handler) storeToGraphQL(store *models.Store) map[string]interface{} {
	if store == nil {
		return nil
	}

	settings, _ := json.Marshal(store.Settings)
	var settingsMap map[string]interface{}
	json.Unmarshal(settings, &settingsMap)

	return map[string]interface{}{
		"id":          store.ID,
		"merchantId":  store.MerchantID,
		"name":        store.Name,
		"description": store.Description,
		"storeType":   store.StoreType,
		"status":      store.Status,
		"domain":      store.Domain,
		"isActive":    store.IsActive,
		"settings":    settingsMap,
		"createdAt":   store.CreatedAt,
		"updatedAt":   store.UpdatedAt,
	}
}

func (h *Handler) subscriptionToGraphQL(subscription *models.Subscription) map[string]interface{} {
	if subscription == nil {
		return nil
	}

	features, _ := json.Marshal(subscription.Features)
	var featuresArray []string
	json.Unmarshal(features, &featuresArray)

	limits, _ := json.Marshal(subscription.Limits)
	var limitsMap map[string]interface{}
	json.Unmarshal(limits, &limitsMap)

	return map[string]interface{}{
		"id":                 subscription.ID,
		"merchantId":         subscription.MerchantID,
		"planId":             subscription.PlanID,
		"status":             subscription.Status,
		"billingCycle":       subscription.BillingCycle,
		"currentPeriodStart": subscription.CurrentPeriodStart,
		"currentPeriodEnd":   subscription.CurrentPeriodEnd,
		"cancelAt":           subscription.CancelAt,
		"canceledAt":         subscription.CanceledAt,
		"trialStart":         subscription.TrialStart,
		"trialEnd":           subscription.TrialEnd,
		"amount":             subscription.Amount,
		"currency":           subscription.Currency,
		"features":           featuresArray,
		"limits":             limitsMap,
		"createdAt":          subscription.CreatedAt,
		"updatedAt":          subscription.UpdatedAt,
	}
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)/2-len(substr)/2:len(s)/2+len(substr)/2] == substr
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
