package handlers

import (

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"unified-commerce/services/merchant-account/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// MerchantHandler handles HTTP requests for merchant operations
type MerchantHandler struct {
	service   *service.MerchantService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewMerchantHandler creates a new merchant handler
func NewMerchantHandler(service *service.MerchantService, logger *logger.Logger) *MerchantHandler {
	return &MerchantHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all merchant routes
func (h *MerchantHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	
	// Authentication required for all routes
	authConfig := middleware.DefaultAuthConfig()
	v1.Use(middleware.JWTAuth(authConfig))

	// Merchant management routes
	merchants := v1.Group("/merchants")
	{
		merchants.POST("/", h.CreateMerchant)
		merchants.GET("/", h.GetUserMerchants)
		merchants.GET("/:id", h.GetMerchant)
		merchants.PUT("/:id", h.UpdateMerchant)
		merchants.POST("/:id/activate", h.ActivateMerchant)
		merchants.POST("/:id/suspend", h.SuspendMerchant)
		
		// Member management
		merchants.GET("/:id/members", h.GetMerchantMembers)
		merchants.POST("/:id/members/invite", h.InviteMember)
		merchants.POST("/:id/members/:memberId/accept", h.AcceptInvitation)
		merchants.DELETE("/:id/members/:memberId", h.RemoveMember)
		
		// Address management
		merchants.GET("/:id/addresses", h.GetMerchantAddresses)
		merchants.POST("/:id/addresses", h.CreateMerchantAddress)
		merchants.PUT("/:id/addresses/:addressId", h.UpdateMerchantAddress)
		merchants.POST("/:id/addresses/:addressId/default", h.SetDefaultAddress)
		
		// Subscription management
		merchants.GET("/:id/subscription", h.GetSubscription)
		merchants.POST("/:id/subscription", h.CreateSubscription)
		merchants.POST("/:id/subscription/cancel", h.CancelSubscription)
	}

	// Plan routes (public information)
	plans := v1.Group("/plans")
	{
		plans.GET("/", h.ListPlans)
		plans.GET("/:id", h.GetPlan)
	}

	// Admin routes (require admin role)
	admin := v1.Group("/admin")
	admin.Use(middleware.RequireRole("admin", "super_admin"))
	{
		admin.GET("/merchants", h.ListAllMerchants)
		admin.GET("/merchants/:id", h.GetMerchantAdmin)
		admin.POST("/merchants/:id/verify", h.VerifyMerchant)
		admin.POST("/merchants/:id/suspend", h.SuspendMerchantAdmin)
	}
}

// CreateMerchant handles merchant account creation
func (h *MerchantHandler) CreateMerchant(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	var req service.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	// Set the owner user ID from authenticated user
	req.OwnerUserID = userID.(string)

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	merchant, err := h.service.CreateMerchant(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrMerchantExists:
			httputil.Conflict(c, "Merchant with this email already exists")
		default:
			h.logger.WithError(err).Error("Failed to create merchant")
			httputil.InternalServerError(c, "Failed to create merchant account")
		}
		return
	}

	httputil.Created(c, merchant, "Merchant account created successfully")
}

// GetMerchant retrieves a specific merchant
func (h *MerchantHandler) GetMerchant(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	merchantID := c.Param("id")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	merchant, err := h.service.GetMerchant(c.Request.Context(), merchantID, userID.(string))
	if err != nil {
		switch err {
		case service.ErrMerchantNotFound:
			httputil.NotFound(c, "Merchant not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			h.logger.WithError(err).Error("Failed to get merchant")
			httputil.InternalServerError(c, "Failed to retrieve merchant")
		}
		return
	}

	httputil.Success(c, merchant)
}

// GetUserMerchants retrieves all merchants for the authenticated user
func (h *MerchantHandler) GetUserMerchants(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	merchants, err := h.service.GetMerchantsByUser(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user merchants")
		httputil.InternalServerError(c, "Failed to retrieve merchants")
		return
	}

	httputil.Success(c, merchants)
}

// UpdateMerchant handles merchant information updates
func (h *MerchantHandler) UpdateMerchant(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	merchantID := c.Param("id")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	var req service.UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	merchant, err := h.service.UpdateMerchant(c.Request.Context(), merchantID, userID.(string), &req)
	if err != nil {
		switch err {
		case service.ErrMerchantNotFound:
			httputil.NotFound(c, "Merchant not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			h.logger.WithError(err).Error("Failed to update merchant")
			httputil.InternalServerError(c, "Failed to update merchant")
		}
		return
	}

	httputil.Success(c, merchant, "Merchant updated successfully")
}

// InviteMember handles member invitation
func (h *MerchantHandler) InviteMember(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	merchantID := c.Param("id")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	var req service.InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	member, err := h.service.InviteMember(c.Request.Context(), merchantID, userID.(string), &req)
	if err != nil {
		switch err {
		case service.ErrMerchantNotFound:
			httputil.NotFound(c, "Merchant not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		case service.ErrInvalidRole:
			httputil.BadRequest(c, "Invalid role specified")
		default:
			h.logger.WithError(err).Error("Failed to invite member")
			httputil.InternalServerError(c, "Failed to invite member")
		}
		return
	}

	httputil.Created(c, member, "Member invited successfully")
}

// ListPlans returns all available subscription plans
func (h *MerchantHandler) ListPlans(c *gin.Context) {
	plans, err := h.service.ListPlans(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to list plans")
		httputil.InternalServerError(c, "Failed to retrieve plans")
		return
	}

	httputil.Success(c, plans)
}

// CreateSubscription handles subscription creation
func (h *MerchantHandler) CreateSubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	merchantID := c.Param("id")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	var req struct {
		PlanID string `json:"plan_id" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	subscription, err := h.service.CreateSubscription(c.Request.Context(), merchantID, req.PlanID, userID.(string))
	if err != nil {
		switch err {
		case service.ErrMerchantNotFound:
			httputil.NotFound(c, "Merchant not found")
		case service.ErrPlanNotFound:
			httputil.NotFound(c, "Plan not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			h.logger.WithError(err).Error("Failed to create subscription")
			httputil.InternalServerError(c, "Failed to create subscription")
		}
		return
	}

	httputil.Created(c, subscription, "Subscription created successfully")
}

// GetSubscription retrieves merchant's subscription
func (h *MerchantHandler) GetSubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		httputil.Unauthorized(c, "Authentication required")
		return
	}

	merchantID := c.Param("id")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	subscription, err := h.service.GetSubscription(c.Request.Context(), merchantID, userID.(string))
	if err != nil {
		switch err {
		case service.ErrMerchantNotFound:
			httputil.NotFound(c, "Merchant not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			h.logger.WithError(err).Error("Failed to get subscription")
			httputil.InternalServerError(c, "Failed to retrieve subscription")
		}
		return
	}

	if subscription == nil {
		httputil.NotFound(c, "No active subscription found")
		return
	}

	httputil.Success(c, subscription)
}

// Placeholder implementations for remaining endpoints

func (h *MerchantHandler) ActivateMerchant(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Activate merchant endpoint - to be implemented"})
}

func (h *MerchantHandler) SuspendMerchant(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Suspend merchant endpoint - to be implemented"})
}

func (h *MerchantHandler) GetMerchantMembers(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get merchant members endpoint - to be implemented"})
}

func (h *MerchantHandler) AcceptInvitation(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Accept invitation endpoint - to be implemented"})
}

func (h *MerchantHandler) RemoveMember(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Remove member endpoint - to be implemented"})
}

func (h *MerchantHandler) GetMerchantAddresses(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get merchant addresses endpoint - to be implemented"})
}

func (h *MerchantHandler) CreateMerchantAddress(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Create merchant address endpoint - to be implemented"})
}

func (h *MerchantHandler) UpdateMerchantAddress(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Update merchant address endpoint - to be implemented"})
}

func (h *MerchantHandler) SetDefaultAddress(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Set default address endpoint - to be implemented"})
}

func (h *MerchantHandler) GetPlan(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get plan endpoint - to be implemented"})
}

func (h *MerchantHandler) CancelSubscription(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Cancel subscription endpoint - to be implemented"})
}

func (h *MerchantHandler) ListAllMerchants(c *gin.Context) {
	pagination := httputil.GetPaginationParams(c)
	
	httputil.SuccessWithMeta(c, []interface{}{}, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   0,
	}, "Admin merchant list endpoint - to be implemented")
}

func (h *MerchantHandler) GetMerchantAdmin(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Admin get merchant endpoint - to be implemented"})
}

func (h *MerchantHandler) VerifyMerchant(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Verify merchant endpoint - to be implemented"})
}

func (h *MerchantHandler) SuspendMerchantAdmin(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Admin suspend merchant endpoint - to be implemented"})
}