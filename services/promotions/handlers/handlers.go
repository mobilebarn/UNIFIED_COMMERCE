package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"unified-commerce/services/promotions/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// PromotionsHandler handles HTTP requests for promotions operations
type PromotionsHandler struct {
	service   *service.PromotionsService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewPromotionsHandler creates a new promotions handler
func NewPromotionsHandler(service *service.PromotionsService, logger *logger.Logger) *PromotionsHandler {
	return &PromotionsHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all promotions routes
func (h *PromotionsHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// Protected routes (authentication required)
		authConfig := middleware.DefaultAuthConfig()
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(authConfig))
		{
			// Promotion management
			promotions := protected.Group("/promotions")
			{
				promotions.POST("", h.CreatePromotion)
				promotions.GET("/:id", h.GetPromotion)
				promotions.GET("/merchant/:merchantId", h.GetPromotionsByMerchant)
				promotions.PUT("/:id", h.UpdatePromotion)
				promotions.DELETE("/:id", h.DeletePromotion)
			}

			// Discount code management
			discountCodes := protected.Group("/discount-codes")
			{
				discountCodes.POST("", h.CreateDiscountCode)
				discountCodes.GET("/:id", h.GetDiscountCode)
				discountCodes.GET("/code/:code", h.GetDiscountCodeByCode)
				discountCodes.GET("/promotion/:promotionId", h.GetDiscountCodesByPromotion)
				discountCodes.PUT("/:id", h.UpdateDiscountCode)
				discountCodes.DELETE("/:id", h.DeleteDiscountCode)
				discountCodes.POST("/validate", h.ValidateDiscountCode)
			}

			// Gift card management
			giftCards := protected.Group("/gift-cards")
			{
				giftCards.POST("", h.CreateGiftCard)
				giftCards.GET("/:id", h.GetGiftCard)
				giftCards.GET("/code/:code", h.GetGiftCardByCode)
				giftCards.GET("/merchant/:merchantId", h.GetGiftCardsByMerchant)
				giftCards.PUT("/:id", h.UpdateGiftCard)
				giftCards.POST("/:id/redeem", h.RedeemGiftCard)
				giftCards.POST("/:id/adjust", h.AdjustGiftCardBalance)
			}

			// Loyalty program management
			loyalty := protected.Group("/loyalty")
			{
				loyalty.POST("/programs", h.CreateLoyaltyProgram)
				loyalty.GET("/programs/:id", h.GetLoyaltyProgram)
				loyalty.GET("/programs/merchant/:merchantId", h.GetLoyaltyProgramsByMerchant)
				loyalty.PUT("/programs/:id", h.UpdateLoyaltyProgram)

				// Loyalty member management
				loyalty.POST("/members", h.CreateLoyaltyMember)
				loyalty.GET("/members/:id", h.GetLoyaltyMember)
				loyalty.GET("/members/customer/:customerId/program/:programId", h.GetLoyaltyMemberByCustomer)
				loyalty.PUT("/members/:id", h.UpdateLoyaltyMember)
				loyalty.POST("/members/:id/earn", h.EarnLoyaltyPoints)
				loyalty.POST("/members/:id/redeem", h.RedeemLoyaltyPoints)
				loyalty.GET("/members/:id/activities", h.GetLoyaltyActivities)

				// Loyalty tier management
				loyalty.POST("/tiers", h.CreateLoyaltyTier)
				loyalty.GET("/tiers/program/:programId", h.GetLoyaltyTiersByProgram)
				loyalty.PUT("/tiers/:id", h.UpdateLoyaltyTier)
			}

			// Administrative operations
			admin := protected.Group("/admin")
			{
				admin.GET("/reports/promotions", h.GetPromotionReport)
				admin.GET("/reports/loyalty", h.GetLoyaltyReport)
			}
		}

		// Public routes for discount code validation
		public := v1.Group("/public")
		{
			public.POST("/discount-codes/validate", h.ValidateDiscountCode)
		}
	}
}

// Promotion Management Handlers

// CreatePromotion handles creating a new promotion
func (h *PromotionsHandler) CreatePromotion(c *gin.Context) {
	var req service.CreatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	promotion, err := h.service.CreatePromotion(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create promotion")
		httputil.InternalServerError(c, "Failed to create promotion")
		return
	}

	httputil.Created(c, promotion, "Promotion created successfully")
}

// GetPromotion handles retrieving a promotion by ID
func (h *PromotionsHandler) GetPromotion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid promotion ID")
		return
	}

	promotion, err := h.service.GetPromotion(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrPromotionNotFound:
			httputil.NotFound(c, "Promotion not found")
		default:
			h.logger.WithError(err).Error("Failed to get promotion")
			httputil.InternalServerError(c, "Failed to get promotion")
		}
		return
	}

	httputil.Success(c, promotion, "Promotion retrieved successfully")
}

// GetPromotionsByMerchant handles retrieving promotions for a merchant
func (h *PromotionsHandler) GetPromotionsByMerchant(c *gin.Context) {
	merchantID, err := uuid.Parse(c.Param("merchantId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid merchant ID")
		return
	}

	pagination := httputil.GetPaginationParams(c)
	filters := make(map[string]interface{})

	// Parse optional filters from query parameters
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if promoType := c.Query("type"); promoType != "" {
		filters["type"] = promoType
	}

	promotions, total, err := h.service.GetPromotionsByMerchant(c.Request.Context(), merchantID, filters, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get promotions by merchant")
		httputil.InternalServerError(c, "Failed to get promotions")
		return
	}

	response := map[string]interface{}{
		"data":        promotions,
		"total":       total,
		"page":        pagination.Page,
		"per_page":    pagination.PerPage,
		"total_pages": (total + int64(pagination.PerPage) - 1) / int64(pagination.PerPage),
	}

	httputil.Success(c, response, "Promotions retrieved successfully")
}

// UpdatePromotion handles updating a promotion
func (h *PromotionsHandler) UpdatePromotion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid promotion ID")
		return
	}

	var req service.UpdatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	promotion, err := h.service.UpdatePromotion(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrPromotionNotFound:
			httputil.NotFound(c, "Promotion not found")
		default:
			h.logger.WithError(err).Error("Failed to update promotion")
			httputil.InternalServerError(c, "Failed to update promotion")
		}
		return
	}

	httputil.Success(c, promotion, "Promotion updated successfully")
}

// DeletePromotion handles deleting a promotion
func (h *PromotionsHandler) DeletePromotion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid promotion ID")
		return
	}

	if err := h.service.DeletePromotion(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrPromotionNotFound:
			httputil.NotFound(c, "Promotion not found")
		default:
			h.logger.WithError(err).Error("Failed to delete promotion")
			httputil.InternalServerError(c, "Failed to delete promotion")
		}
		return
	}

	httputil.Success(c, nil, "Promotion deleted successfully")
}

// Discount Code Management Handlers

// CreateDiscountCode handles creating a new discount code
func (h *PromotionsHandler) CreateDiscountCode(c *gin.Context) {
	var req service.CreateDiscountCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	discountCode, err := h.service.CreateDiscountCode(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrPromotionNotFound:
			httputil.NotFound(c, "Promotion not found")
		default:
			h.logger.WithError(err).Error("Failed to create discount code")
			httputil.InternalServerError(c, "Failed to create discount code")
		}
		return
	}

	httputil.Created(c, discountCode, "Discount code created successfully")
}

// GetDiscountCode handles retrieving a discount code by ID
func (h *PromotionsHandler) GetDiscountCode(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid discount code ID")
		return
	}

	discountCode, err := h.service.GetDiscountCode(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrDiscountCodeNotFound:
			httputil.NotFound(c, "Discount code not found")
		default:
			h.logger.WithError(err).Error("Failed to get discount code")
			httputil.InternalServerError(c, "Failed to get discount code")
		}
		return
	}

	httputil.Success(c, discountCode, "Discount code retrieved successfully")
}

// GetDiscountCodeByCode handles retrieving a discount code by code
func (h *PromotionsHandler) GetDiscountCodeByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		httputil.BadRequest(c, "Discount code is required")
		return
	}

	discountCode, err := h.service.GetDiscountCodeByCode(c.Request.Context(), code)
	if err != nil {
		switch err {
		case service.ErrDiscountCodeNotFound:
			httputil.NotFound(c, "Discount code not found")
		default:
			h.logger.WithError(err).Error("Failed to get discount code by code")
			httputil.InternalServerError(c, "Failed to get discount code")
		}
		return
	}

	httputil.Success(c, discountCode, "Discount code retrieved successfully")
}

// GetDiscountCodesByPromotion handles retrieving discount codes for a promotion
func (h *PromotionsHandler) GetDiscountCodesByPromotion(c *gin.Context) {
	promotionID, err := uuid.Parse(c.Param("promotionId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid promotion ID")
		return
	}

	discountCodes, err := h.service.GetDiscountCodesByPromotion(c.Request.Context(), promotionID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get discount codes by promotion")
		httputil.InternalServerError(c, "Failed to get discount codes")
		return
	}

	httputil.Success(c, discountCodes, "Discount codes retrieved successfully")
}

// UpdateDiscountCode handles updating a discount code
func (h *PromotionsHandler) UpdateDiscountCode(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid discount code ID")
		return
	}

	var req service.UpdateDiscountCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	discountCode, err := h.service.UpdateDiscountCode(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrDiscountCodeNotFound:
			httputil.NotFound(c, "Discount code not found")
		default:
			h.logger.WithError(err).Error("Failed to update discount code")
			httputil.InternalServerError(c, "Failed to update discount code")
		}
		return
	}

	httputil.Success(c, discountCode, "Discount code updated successfully")
}

// DeleteDiscountCode handles deleting a discount code
func (h *PromotionsHandler) DeleteDiscountCode(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid discount code ID")
		return
	}

	if err := h.service.DeleteDiscountCode(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrDiscountCodeNotFound:
			httputil.NotFound(c, "Discount code not found")
		default:
			h.logger.WithError(err).Error("Failed to delete discount code")
			httputil.InternalServerError(c, "Failed to delete discount code")
		}
		return
	}

	httputil.Success(c, nil, "Discount code deleted successfully")
}

// ValidateDiscountCode handles validating a discount code
func (h *PromotionsHandler) ValidateDiscountCode(c *gin.Context) {
	var req struct {
		Code        string  `json:"code" validate:"required"`
		OrderAmount float64 `json:"order_amount" validate:"required,min=0"`
		CustomerID  string  `json:"customer_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	var customerID *uuid.UUID
	if req.CustomerID != "" {
		id, err := uuid.Parse(req.CustomerID)
		if err != nil {
			httputil.BadRequest(c, "Invalid customer ID")
			return
		}
		customerID = &id
	}

	discountCode, discountAmount, err := h.service.ValidateDiscountCode(c.Request.Context(), req.Code, req.OrderAmount, customerID)
	if err != nil {
		switch err {
		case service.ErrInvalidCode:
			httputil.BadRequest(c, "Invalid discount code")
		case service.ErrCodeExpired:
			httputil.BadRequest(c, "Discount code has expired")
		case service.ErrCodeInactive:
			httputil.BadRequest(c, "Discount code is inactive")
		case service.ErrUsageLimitExceeded:
			httputil.BadRequest(c, "Discount code usage limit exceeded")
		case service.ErrCustomerUsageLimit:
			httputil.BadRequest(c, "Customer usage limit exceeded for this discount code")
		default:
			h.logger.WithError(err).Error("Failed to validate discount code")
			httputil.InternalServerError(c, "Failed to validate discount code")
		}
		return
	}

	response := map[string]interface{}{
		"discount_code":    discountCode,
		"discount_amount":  discountAmount,
		"discounted_total": req.OrderAmount - discountAmount,
	}

	httputil.Success(c, response, "Discount code validated successfully")
}

// Gift Card Management Handlers

// CreateGiftCard handles creating a new gift card
func (h *PromotionsHandler) CreateGiftCard(c *gin.Context) {
	var req service.CreateGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	giftCard, err := h.service.CreateGiftCard(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create gift card")
		httputil.InternalServerError(c, "Failed to create gift card")
		return
	}

	httputil.Created(c, giftCard, "Gift card created successfully")
}

// GetGiftCard handles retrieving a gift card by ID
func (h *PromotionsHandler) GetGiftCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid gift card ID")
		return
	}

	giftCard, err := h.service.GetGiftCard(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrGiftCardNotFound:
			httputil.NotFound(c, "Gift card not found")
		default:
			h.logger.WithError(err).Error("Failed to get gift card")
			httputil.InternalServerError(c, "Failed to get gift card")
		}
		return
	}

	httputil.Success(c, giftCard, "Gift card retrieved successfully")
}

// GetGiftCardByCode handles retrieving a gift card by code
func (h *PromotionsHandler) GetGiftCardByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		httputil.BadRequest(c, "Gift card code is required")
		return
	}

	giftCard, err := h.service.GetGiftCardByCode(c.Request.Context(), code)
	if err != nil {
		switch err {
		case service.ErrGiftCardNotFound:
			httputil.NotFound(c, "Gift card not found")
		default:
			h.logger.WithError(err).Error("Failed to get gift card by code")
			httputil.InternalServerError(c, "Failed to get gift card")
		}
		return
	}

	httputil.Success(c, giftCard, "Gift card retrieved successfully")
}

// GetGiftCardsByMerchant handles retrieving gift cards for a merchant
func (h *PromotionsHandler) GetGiftCardsByMerchant(c *gin.Context) {
	merchantID, err := uuid.Parse(c.Param("merchantId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid merchant ID")
		return
	}

	pagination := httputil.GetPaginationParams(c)
	filters := make(map[string]interface{})

	// Parse optional filters from query parameters
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	giftCards, total, err := h.service.GetGiftCardsByMerchant(c.Request.Context(), merchantID, filters, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get gift cards by merchant")
		httputil.InternalServerError(c, "Failed to get gift cards")
		return
	}

	response := map[string]interface{}{
		"data":        giftCards,
		"total":       total,
		"page":        pagination.Page,
		"per_page":    pagination.PerPage,
		"total_pages": (total + int64(pagination.PerPage) - 1) / int64(pagination.PerPage),
	}

	httputil.Success(c, response, "Gift cards retrieved successfully")
}

// UpdateGiftCard handles updating a gift card
func (h *PromotionsHandler) UpdateGiftCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid gift card ID")
		return
	}

	var req service.UpdateGiftCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	giftCard, err := h.service.UpdateGiftCard(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrGiftCardNotFound:
			httputil.NotFound(c, "Gift card not found")
		default:
			h.logger.WithError(err).Error("Failed to update gift card")
			httputil.InternalServerError(c, "Failed to update gift card")
		}
		return
	}

	httputil.Success(c, giftCard, "Gift card updated successfully")
}

// RedeemGiftCard handles redeeming a gift card
func (h *PromotionsHandler) RedeemGiftCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid gift card ID")
		return
	}

	var req struct {
		OrderID string  `json:"order_id"`
		Amount  float64 `json:"amount" validate:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	var orderID *uuid.UUID
	if req.OrderID != "" {
		id, err := uuid.Parse(req.OrderID)
		if err != nil {
			httputil.BadRequest(c, "Invalid order ID")
			return
		}
		orderID = &id
	}

	transaction, err := h.service.RedeemGiftCard(c.Request.Context(), id, orderID, req.Amount)
	if err != nil {
		switch err {
		case service.ErrGiftCardNotFound:
			httputil.NotFound(c, "Gift card not found")
		case service.ErrGiftCardExpired:
			httputil.BadRequest(c, "Gift card has expired")
		case service.ErrInsufficientBalance:
			httputil.BadRequest(c, "Insufficient gift card balance")
		case service.ErrInvalidAmount:
			httputil.BadRequest(c, "Invalid amount")
		default:
			h.logger.WithError(err).Error("Failed to redeem gift card")
			httputil.InternalServerError(c, "Failed to redeem gift card")
		}
		return
	}

	httputil.Success(c, transaction, "Gift card redeemed successfully")
}

// AdjustGiftCardBalance handles adjusting a gift card balance
func (h *PromotionsHandler) AdjustGiftCardBalance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid gift card ID")
		return
	}

	var req struct {
		Amount      float64 `json:"amount" validate:"required"`
		Description string  `json:"description" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	transaction, err := h.service.AdjustGiftCardBalance(c.Request.Context(), id, req.Amount, req.Description)
	if err != nil {
		switch err {
		case service.ErrGiftCardNotFound:
			httputil.NotFound(c, "Gift card not found")
		case service.ErrInsufficientBalance:
			httputil.BadRequest(c, "Insufficient gift card balance")
		case service.ErrInvalidAmount:
			httputil.BadRequest(c, "Invalid amount")
		default:
			h.logger.WithError(err).Error("Failed to adjust gift card balance")
			httputil.InternalServerError(c, "Failed to adjust gift card balance")
		}
		return
	}

	httputil.Success(c, transaction, "Gift card balance adjusted successfully")
}

// Loyalty Program Management Handlers

// CreateLoyaltyProgram handles creating a new loyalty program
func (h *PromotionsHandler) CreateLoyaltyProgram(c *gin.Context) {
	var req service.CreateLoyaltyProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	program, err := h.service.CreateLoyaltyProgram(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create loyalty program")
		httputil.InternalServerError(c, "Failed to create loyalty program")
		return
	}

	httputil.Created(c, program, "Loyalty program created successfully")
}

// GetLoyaltyProgram handles retrieving a loyalty program by ID
func (h *PromotionsHandler) GetLoyaltyProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty program ID")
		return
	}

	program, err := h.service.GetLoyaltyProgram(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrLoyaltyProgramNotFound:
			httputil.NotFound(c, "Loyalty program not found")
		default:
			h.logger.WithError(err).Error("Failed to get loyalty program")
			httputil.InternalServerError(c, "Failed to get loyalty program")
		}
		return
	}

	httputil.Success(c, program, "Loyalty program retrieved successfully")
}

// GetLoyaltyProgramsByMerchant handles retrieving loyalty programs for a merchant
func (h *PromotionsHandler) GetLoyaltyProgramsByMerchant(c *gin.Context) {
	merchantID, err := uuid.Parse(c.Param("merchantId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid merchant ID")
		return
	}

	programs, err := h.service.GetLoyaltyProgramsByMerchant(c.Request.Context(), merchantID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get loyalty programs by merchant")
		httputil.InternalServerError(c, "Failed to get loyalty programs")
		return
	}

	httputil.Success(c, programs, "Loyalty programs retrieved successfully")
}

// UpdateLoyaltyProgram handles updating a loyalty program
func (h *PromotionsHandler) UpdateLoyaltyProgram(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty program ID")
		return
	}

	var req service.UpdateLoyaltyProgramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	program, err := h.service.UpdateLoyaltyProgram(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrLoyaltyProgramNotFound:
			httputil.NotFound(c, "Loyalty program not found")
		default:
			h.logger.WithError(err).Error("Failed to update loyalty program")
			httputil.InternalServerError(c, "Failed to update loyalty program")
		}
		return
	}

	httputil.Success(c, program, "Loyalty program updated successfully")
}

// Loyalty Member Management Handlers

// CreateLoyaltyMember handles creating a new loyalty member
func (h *PromotionsHandler) CreateLoyaltyMember(c *gin.Context) {
	var req service.CreateLoyaltyMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	member, err := h.service.CreateLoyaltyMember(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrLoyaltyProgramNotFound:
			httputil.NotFound(c, "Loyalty program not found")
		default:
			h.logger.WithError(err).Error("Failed to create loyalty member")
			httputil.InternalServerError(c, "Failed to create loyalty member")
		}
		return
	}

	httputil.Created(c, member, "Loyalty member created successfully")
}

// GetLoyaltyMember handles retrieving a loyalty member by ID
func (h *PromotionsHandler) GetLoyaltyMember(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty member ID")
		return
	}

	member, err := h.service.GetLoyaltyMember(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrLoyaltyMemberNotFound:
			httputil.NotFound(c, "Loyalty member not found")
		default:
			h.logger.WithError(err).Error("Failed to get loyalty member")
			httputil.InternalServerError(c, "Failed to get loyalty member")
		}
		return
	}

	httputil.Success(c, member, "Loyalty member retrieved successfully")
}

// GetLoyaltyMemberByCustomer handles retrieving a loyalty member by customer ID and program ID
func (h *PromotionsHandler) GetLoyaltyMemberByCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid customer ID")
		return
	}

	programID, err := uuid.Parse(c.Param("programId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid program ID")
		return
	}

	member, err := h.service.GetLoyaltyMemberByCustomer(c.Request.Context(), customerID, programID)
	if err != nil {
		switch err {
		case service.ErrLoyaltyMemberNotFound:
			httputil.NotFound(c, "Loyalty member not found")
		default:
			h.logger.WithError(err).Error("Failed to get loyalty member by customer")
			httputil.InternalServerError(c, "Failed to get loyalty member")
		}
		return
	}

	httputil.Success(c, member, "Loyalty member retrieved successfully")
}

// UpdateLoyaltyMember handles updating a loyalty member
func (h *PromotionsHandler) UpdateLoyaltyMember(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty member ID")
		return
	}

	var req service.UpdateLoyaltyMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	member, err := h.service.UpdateLoyaltyMember(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrLoyaltyMemberNotFound:
			httputil.NotFound(c, "Loyalty member not found")
		default:
			h.logger.WithError(err).Error("Failed to update loyalty member")
			httputil.InternalServerError(c, "Failed to update loyalty member")
		}
		return
	}

	httputil.Success(c, member, "Loyalty member updated successfully")
}

// EarnLoyaltyPoints handles earning points for a loyalty member
func (h *PromotionsHandler) EarnLoyaltyPoints(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty member ID")
		return
	}

	var req struct {
		Points      int    `json:"points" validate:"required,min=1"`
		Description string `json:"description" validate:"required"`
		ReferenceID string `json:"reference_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	var referenceID *uuid.UUID
	if req.ReferenceID != "" {
		id, err := uuid.Parse(req.ReferenceID)
		if err != nil {
			httputil.BadRequest(c, "Invalid reference ID")
			return
		}
		referenceID = &id
	}

	if err := h.service.EarnLoyaltyPoints(c.Request.Context(), id, req.Points, req.Description, referenceID); err != nil {
		switch err {
		case service.ErrLoyaltyMemberNotFound:
			httputil.NotFound(c, "Loyalty member not found")
		default:
			h.logger.WithError(err).Error("Failed to earn loyalty points")
			httputil.InternalServerError(c, "Failed to earn loyalty points")
		}
		return
	}

	httputil.Success(c, nil, "Loyalty points earned successfully")
}

// RedeemLoyaltyPoints handles redeeming points for a loyalty member
func (h *PromotionsHandler) RedeemLoyaltyPoints(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty member ID")
		return
	}

	var req struct {
		Points      int    `json:"points" validate:"required,min=1"`
		Description string `json:"description" validate:"required"`
		ReferenceID string `json:"reference_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	var referenceID *uuid.UUID
	if req.ReferenceID != "" {
		id, err := uuid.Parse(req.ReferenceID)
		if err != nil {
			httputil.BadRequest(c, "Invalid reference ID")
			return
		}
		referenceID = &id
	}

	if err := h.service.RedeemLoyaltyPoints(c.Request.Context(), id, req.Points, req.Description, referenceID); err != nil {
		switch err {
		case service.ErrLoyaltyMemberNotFound:
			httputil.NotFound(c, "Loyalty member not found")
		default:
			h.logger.WithError(err).Error("Failed to redeem loyalty points")
			httputil.InternalServerError(c, "Failed to redeem loyalty points")
		}
		return
	}

	httputil.Success(c, nil, "Loyalty points redeemed successfully")
}

// GetLoyaltyActivities handles retrieving activities for a loyalty member
func (h *PromotionsHandler) GetLoyaltyActivities(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty member ID")
		return
	}

	pagination := httputil.GetPaginationParams(c)
	activities, err := h.service.GetLoyaltyActivities(c.Request.Context(), id, pagination.Page, pagination.PerPage)
	if err != nil {
		switch err {
		case service.ErrLoyaltyMemberNotFound:
			httputil.NotFound(c, "Loyalty member not found")
		default:
			h.logger.WithError(err).Error("Failed to get loyalty activities")
			httputil.InternalServerError(c, "Failed to get loyalty activities")
		}
		return
	}

	httputil.Success(c, activities, "Loyalty activities retrieved successfully")
}

// Loyalty Tier Management Handlers

// CreateLoyaltyTier handles creating a new loyalty tier
func (h *PromotionsHandler) CreateLoyaltyTier(c *gin.Context) {
	var req service.CreateLoyaltyTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	tier, err := h.service.CreateLoyaltyTier(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create loyalty tier")
		httputil.InternalServerError(c, "Failed to create loyalty tier")
		return
	}

	httputil.Created(c, tier, "Loyalty tier created successfully")
}

// GetLoyaltyTiersByProgram handles retrieving loyalty tiers for a program
func (h *PromotionsHandler) GetLoyaltyTiersByProgram(c *gin.Context) {
	programID, err := uuid.Parse(c.Param("programId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid program ID")
		return
	}

	tiers, err := h.service.GetLoyaltyTiersByProgram(c.Request.Context(), programID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get loyalty tiers by program")
		httputil.InternalServerError(c, "Failed to get loyalty tiers")
		return
	}

	httputil.Success(c, tiers, "Loyalty tiers retrieved successfully")
}

// UpdateLoyaltyTier handles updating a loyalty tier
func (h *PromotionsHandler) UpdateLoyaltyTier(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid loyalty tier ID")
		return
	}

	var req service.UpdateLoyaltyTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	tier, err := h.service.UpdateLoyaltyTier(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update loyalty tier")
		httputil.InternalServerError(c, "Failed to update loyalty tier")
		return
	}

	httputil.Success(c, tier, "Loyalty tier updated successfully")
}

// Administrative Operations Handlers

// GetPromotionReport handles retrieving promotion report
func (h *PromotionsHandler) GetPromotionReport(c *gin.Context) {
	// In a real implementation, this would generate a promotion report
	httputil.Success(c, nil, "Promotion report retrieved successfully")
}

// GetLoyaltyReport handles retrieving loyalty report
func (h *PromotionsHandler) GetLoyaltyReport(c *gin.Context) {
	// In a real implementation, this would generate a loyalty report
	httputil.Success(c, nil, "Loyalty report retrieved successfully")
}
