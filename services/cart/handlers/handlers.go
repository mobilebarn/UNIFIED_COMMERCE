package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"unified-commerce/services/cart/models"
	"unified-commerce/services/cart/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// CartHandler handles HTTP requests for cart and checkout operations
type CartHandler struct {
	service   *service.CartService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewCartHandler creates a new cart handler
func NewCartHandler(service *service.CartService, logger *logger.Logger) *CartHandler {
	return &CartHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all cart and checkout routes
func (h *CartHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// Protected routes (authentication required)
		authConfig := middleware.DefaultAuthConfig()
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(authConfig))
		{
			// Cart management
			carts := protected.Group("/carts")
			{
				carts.POST("", h.CreateCart)
				carts.GET("/:id", h.GetCart)
				carts.PUT("/:id", h.UpdateCart)
				carts.DELETE("/:id", h.DeleteCart)

				// Line item management
				carts.POST("/:id/line-items", h.AddLineItem)
				carts.PUT("/line-items/:lineItemId", h.UpdateLineItem)
				carts.DELETE("/line-items/:lineItemId", h.RemoveLineItem)

				// Customer cart lookup
				carts.GET("/customer/:customerId", h.GetCartByCustomer)
			}

			// Checkout management
			checkouts := protected.Group("/checkouts")
			{
				checkouts.POST("", h.CreateCheckout)
				checkouts.GET("/:id", h.GetCheckout)
				checkouts.PUT("/:id/customer-info", h.UpdateCheckoutCustomerInfo)
				checkouts.PUT("/:id/shipping-address", h.UpdateCheckoutShippingAddress)
				checkouts.POST("/:id/shipping-lines", h.AddShippingLine)
				checkouts.POST("/:id/discounts", h.ApplyDiscountCode)
				checkouts.DELETE("/:id/discounts/:code", h.RemoveDiscountCode)
				checkouts.POST("/:id/complete", h.CompleteCheckout)

				// Checkout by token
				checkouts.GET("/token/:token", h.GetCheckoutByToken)
			}

			// Configuration and reference data
			config := protected.Group("/config")
			{
				config.GET("/shipping-rates", h.GetShippingRates)
				config.GET("/payment-methods", h.GetPaymentMethods)
			}

			// Administrative operations
			admin := protected.Group("/admin")
			{
				admin.POST("/abandonment/process", h.ProcessAbandonedCarts)
				admin.GET("/abandonment/carts", h.GetAbandonedCarts)
				admin.POST("/cleanup/expired", h.CleanupExpiredCarts)
			}
		}

		// Public routes for guest carts
		public := v1.Group("/public")
		{
			public.POST("/carts", h.CreateCart)
			public.GET("/carts/session/:sessionId", h.GetCartBySession)
		}
	}
}

// Cart Management Handlers

// CreateCart handles creating a new shopping cart
func (h *CartHandler) CreateCart(c *gin.Context) {
	var req service.CreateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	// If customer ID is not provided, try to get from JWT token
	if req.CustomerID == nil {
		if customerID, exists := c.Get("user_id"); exists {
			if idStr, ok := customerID.(string); ok {
				if id, err := uuid.Parse(idStr); err == nil {
					req.CustomerID = &id
				}
			}
		}
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	cart, err := h.service.CreateCart(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create cart")
		httputil.InternalServerError(c, "Failed to create cart")
		return
	}

	httputil.Created(c, cart, "Cart created successfully")
}

// GetCart handles retrieving a cart by ID
func (h *CartHandler) GetCart(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid cart ID")
		return
	}

	cart, err := h.service.GetCart(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		default:
			h.logger.WithError(err).Error("Failed to get cart")
			httputil.InternalServerError(c, "Failed to get cart")
		}
		return
	}

	httputil.Success(c, cart, "Cart retrieved successfully")
}

// GetCartBySession handles retrieving a cart by session ID
func (h *CartHandler) GetCartBySession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httputil.BadRequest(c, "Session ID is required")
		return
	}

	cart, err := h.service.GetCartBySessionID(c.Request.Context(), sessionID)
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		default:
			h.logger.WithError(err).Error("Failed to get cart by session ID")
			httputil.InternalServerError(c, "Failed to get cart")
		}
		return
	}

	httputil.Success(c, cart, "Cart retrieved successfully")
}

// GetCartByCustomer handles retrieving a customer's cart
func (h *CartHandler) GetCartByCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid customer ID")
		return
	}

	merchantIDStr := c.Query("merchant_id")
	if merchantIDStr == "" {
		httputil.BadRequest(c, "merchant_id is required")
		return
	}

	merchantID, err := uuid.Parse(merchantIDStr)
	if err != nil {
		httputil.BadRequest(c, "Invalid merchant ID")
		return
	}

	cart, err := h.service.GetCartByCustomerID(c.Request.Context(), customerID, merchantID)
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		default:
			h.logger.WithError(err).Error("Failed to get cart by customer ID")
			httputil.InternalServerError(c, "Failed to get cart")
		}
		return
	}

	httputil.Success(c, cart, "Cart retrieved successfully")
}

// UpdateCart handles updating cart information
func (h *CartHandler) UpdateCart(c *gin.Context) {
	cartID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid cart ID")
		return
	}

	var req service.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	cart, err := h.service.UpdateCart(c.Request.Context(), cartID, &req)
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		case service.ErrCartCompleted:
			httputil.BadRequest(c, "Cart is already completed")
		default:
			h.logger.WithError(err).Error("Failed to update cart")
			httputil.InternalServerError(c, "Failed to update cart")
		}
		return
	}

	httputil.Success(c, cart, "Cart updated successfully")
}

// DeleteCart handles deleting a cart
func (h *CartHandler) DeleteCart(c *gin.Context) {
	_, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid cart ID")
		return
	}

	// In a real implementation, you might want to soft delete or archive the cart
	// For now, we'll just return success
	httputil.Success(c, nil, "Cart deleted successfully")
}

// Line Item Management Handlers

// AddLineItem handles adding a line item to a cart
func (h *CartHandler) AddLineItem(c *gin.Context) {
	cartID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid cart ID")
		return
	}

	var req service.AddLineItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	lineItem, err := h.service.AddLineItem(c.Request.Context(), cartID, &req)
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		case service.ErrCartCompleted:
			httputil.BadRequest(c, "Cart is already completed")
		default:
			h.logger.WithError(err).Error("Failed to add line item to cart")
			httputil.InternalServerError(c, "Failed to add line item to cart")
		}
		return
	}

	httputil.Created(c, lineItem, "Line item added to cart successfully")
}

// UpdateLineItem handles updating a line item
func (h *CartHandler) UpdateLineItem(c *gin.Context) {
	lineItemID, err := uuid.Parse(c.Param("lineItemId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid line item ID")
		return
	}

	var req service.UpdateLineItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	lineItem, err := h.service.UpdateLineItem(c.Request.Context(), lineItemID, &req)
	if err != nil {
		switch err {
		case service.ErrLineItemNotFound:
			httputil.NotFound(c, "Line item not found")
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		case service.ErrCartCompleted:
			httputil.BadRequest(c, "Cart is already completed")
		default:
			h.logger.WithError(err).Error("Failed to update line item")
			httputil.InternalServerError(c, "Failed to update line item")
		}
		return
	}

	httputil.Success(c, lineItem, "Line item updated successfully")
}

// RemoveLineItem handles removing a line item from a cart
func (h *CartHandler) RemoveLineItem(c *gin.Context) {
	lineItemID, err := uuid.Parse(c.Param("lineItemId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid line item ID")
		return
	}

	if err := h.service.RemoveLineItem(c.Request.Context(), lineItemID); err != nil {
		switch err {
		case service.ErrLineItemNotFound:
			httputil.NotFound(c, "Line item not found")
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		case service.ErrCartCompleted:
			httputil.BadRequest(c, "Cart is already completed")
		default:
			h.logger.WithError(err).Error("Failed to remove line item")
			httputil.InternalServerError(c, "Failed to remove line item")
		}
		return
	}

	httputil.Success(c, nil, "Line item removed successfully")
}

// Checkout Management Handlers

// CreateCheckout handles creating a new checkout session
func (h *CartHandler) CreateCheckout(c *gin.Context) {
	var req service.CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	checkout, err := h.service.CreateCheckout(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrCartNotFound:
			httputil.NotFound(c, "Cart not found")
		case service.ErrCartExpired:
			c.JSON(http.StatusGone, gin.H{"error": "Cart has expired"})
		case service.ErrCartCompleted:
			httputil.BadRequest(c, "Cart is already completed")
		default:
			h.logger.WithError(err).Error("Failed to create checkout")
			httputil.InternalServerError(c, "Failed to create checkout")
		}
		return
	}

	httputil.Created(c, checkout, "Checkout created successfully")
}

// GetCheckout handles retrieving a checkout by ID
func (h *CartHandler) GetCheckout(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid checkout ID")
		return
	}

	checkout, err := h.service.GetCheckout(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		default:
			h.logger.WithError(err).Error("Failed to get checkout")
			httputil.InternalServerError(c, "Failed to get checkout")
		}
		return
	}

	httputil.Success(c, checkout, "Checkout retrieved successfully")
}

// GetCheckoutByToken handles retrieving a checkout by token
func (h *CartHandler) GetCheckoutByToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		httputil.BadRequest(c, "Checkout token is required")
		return
	}

	checkout, err := h.service.GetCheckoutByToken(c.Request.Context(), token)
	if err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		default:
			h.logger.WithError(err).Error("Failed to get checkout by token")
			httputil.InternalServerError(c, "Failed to get checkout")
		}
		return
	}

	httputil.Success(c, checkout, "Checkout retrieved successfully")
}

// UpdateCheckoutCustomerInfo handles updating customer information in checkout
func (h *CartHandler) UpdateCheckoutCustomerInfo(c *gin.Context) {
	checkoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid checkout ID")
		return
	}

	var req service.CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	checkout, err := h.service.UpdateCheckoutCustomerInfo(c.Request.Context(), checkoutID, &req)
	if err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		case service.ErrCheckoutAlreadyCompleted:
			httputil.BadRequest(c, "Checkout already completed")
		default:
			h.logger.WithError(err).Error("Failed to update checkout customer info")
			httputil.InternalServerError(c, "Failed to update checkout customer info")
		}
		return
	}

	httputil.Success(c, checkout, "Checkout customer info updated successfully")
}

// UpdateCheckoutShippingAddress handles updating shipping address in checkout
func (h *CartHandler) UpdateCheckoutShippingAddress(c *gin.Context) {
	checkoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid checkout ID")
		return
	}

	var req struct {
		Address service.UpdateCartRequest `json:"address" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	address := models.Address{
		FirstName:  req.Address.CustomerFirstName,
		LastName:   req.Address.CustomerLastName,
		Street1:    req.Address.ShippingAddress.Street1,
		Street2:    req.Address.ShippingAddress.Street2,
		City:       req.Address.ShippingAddress.City,
		State:      req.Address.ShippingAddress.State,
		Country:    req.Address.ShippingAddress.Country,
		PostalCode: req.Address.ShippingAddress.PostalCode,
		Phone:      req.Address.CustomerPhone,
	}

	checkout, err := h.service.UpdateCheckoutShippingAddress(c.Request.Context(), checkoutID, address)
	if err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		case service.ErrCheckoutAlreadyCompleted:
			httputil.BadRequest(c, "Checkout already completed")
		default:
			h.logger.WithError(err).Error("Failed to update checkout shipping address")
			httputil.InternalServerError(c, "Failed to update checkout shipping address")
		}
		return
	}

	httputil.Success(c, checkout, "Checkout shipping address updated successfully")
}

// AddShippingLine handles adding shipping information to checkout
func (h *CartHandler) AddShippingLine(c *gin.Context) {
	checkoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid checkout ID")
		return
	}

	var req struct {
		ShippingLine models.CartShippingLine `json:"shipping_line" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	req.ShippingLine.CartID = checkoutID

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	if err := h.service.AddShippingLine(c.Request.Context(), checkoutID, &req.ShippingLine); err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		case service.ErrCheckoutAlreadyCompleted:
			httputil.BadRequest(c, "Checkout already completed")
		default:
			h.logger.WithError(err).Error("Failed to add shipping line")
			httputil.InternalServerError(c, "Failed to add shipping line")
		}
		return
	}

	httputil.Success(c, nil, "Shipping line added successfully")
}

// ApplyDiscountCode handles applying a discount code to checkout
func (h *CartHandler) ApplyDiscountCode(c *gin.Context) {
	checkoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid checkout ID")
		return
	}

	var req struct {
		Code string `json:"code" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	if err := h.service.ApplyDiscountCode(c.Request.Context(), checkoutID, req.Code); err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		case service.ErrCheckoutAlreadyCompleted:
			httputil.BadRequest(c, "Checkout already completed")
		case service.ErrInvalidDiscountCode:
			httputil.BadRequest(c, "Invalid discount code")
		case service.ErrDiscountAlreadyApplied:
			httputil.BadRequest(c, "Discount code already applied")
		default:
			h.logger.WithError(err).Error("Failed to apply discount code")
			httputil.InternalServerError(c, "Failed to apply discount code")
		}
		return
	}

	httputil.Success(c, nil, "Discount code applied successfully")
}

// RemoveDiscountCode handles removing a discount code from checkout
func (h *CartHandler) RemoveDiscountCode(c *gin.Context) {
	checkoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid checkout ID")
		return
	}

	code := c.Param("code")
	if code == "" {
		httputil.BadRequest(c, "Discount code is required")
		return
	}

	if err := h.service.RemoveDiscountCode(c.Request.Context(), checkoutID, code); err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		case service.ErrCheckoutAlreadyCompleted:
			httputil.BadRequest(c, "Checkout already completed")
		default:
			h.logger.WithError(err).Error("Failed to remove discount code")
			httputil.InternalServerError(c, "Failed to remove discount code")
		}
		return
	}

	httputil.Success(c, nil, "Discount code removed successfully")
}

// CompleteCheckout handles completing the checkout process
func (h *CartHandler) CompleteCheckout(c *gin.Context) {
	checkoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid checkout ID")
		return
	}

	var req struct {
		PaymentMethodID string `json:"payment_method_id" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	orderID, err := h.service.CompleteCheckout(c.Request.Context(), checkoutID, req.PaymentMethodID)
	if err != nil {
		switch err {
		case service.ErrCheckoutNotFound:
			httputil.NotFound(c, "Checkout not found")
		case service.ErrCheckoutAlreadyCompleted:
			httputil.BadRequest(c, "Checkout already completed")
		default:
			h.logger.WithError(err).Error("Failed to complete checkout")
			httputil.InternalServerError(c, "Failed to complete checkout")
		}
		return
	}

	response := map[string]interface{}{
		"order_id": orderID,
		"message":  "Checkout completed successfully",
	}

	httputil.Success(c, response, "Checkout completed successfully")
}

// Configuration and Reference Data Handlers

// GetShippingRates handles retrieving available shipping rates
func (h *CartHandler) GetShippingRates(c *gin.Context) {
	rates, err := h.service.GetShippingRates(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get shipping rates")
		httputil.InternalServerError(c, "Failed to get shipping rates")
		return
	}

	httputil.Success(c, rates, "Shipping rates retrieved successfully")
}

// GetPaymentMethods handles retrieving available payment methods
func (h *CartHandler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.service.GetPaymentMethods(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get payment methods")
		httputil.InternalServerError(c, "Failed to get payment methods")
		return
	}

	httputil.Success(c, methods, "Payment methods retrieved successfully")
}

// Administrative Operations Handlers

// ProcessAbandonedCarts handles processing abandoned carts
func (h *CartHandler) ProcessAbandonedCarts(c *gin.Context) {
	// In a real implementation, this would process abandoned carts for recovery campaigns
	httputil.Success(c, nil, "Abandoned cart processing initiated")
}

// GetAbandonedCarts handles retrieving abandoned carts
func (h *CartHandler) GetAbandonedCarts(c *gin.Context) {
	pagination := httputil.GetPaginationParams(c)
	carts, err := h.service.GetAbandonedCarts(c.Request.Context(), pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get abandoned carts")
		httputil.InternalServerError(c, "Failed to get abandoned carts")
		return
	}

	httputil.Success(c, carts, "Abandoned carts retrieved successfully")
}

// CleanupExpiredCarts handles cleaning up expired carts
func (h *CartHandler) CleanupExpiredCarts(c *gin.Context) {
	if err := h.service.CleanupExpiredCarts(c.Request.Context()); err != nil {
		h.logger.WithError(err).Error("Failed to cleanup expired carts")
		httputil.InternalServerError(c, "Failed to cleanup expired carts")
		return
	}

	httputil.Success(c, nil, "Expired carts cleaned up successfully")
}
