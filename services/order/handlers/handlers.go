package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"unified-commerce/services/order/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// OrderHandler handles HTTP requests for order operations
type OrderHandler struct {
	service   *service.OrderService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(service *service.OrderService, logger *logger.Logger) *OrderHandler {
	return &OrderHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all order routes
func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// Protected routes (authentication required)
		authConfig := middleware.DefaultAuthConfig()
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(authConfig))
		{
			// Order management
			orders := protected.Group("/orders")
			{
				orders.POST("", h.CreateOrder)
				orders.GET("", h.GetOrders)
				orders.GET("/:id", h.GetOrder)
				orders.PUT("/:id", h.UpdateOrder)
				orders.POST("/:id/confirm", h.ConfirmOrder)
				orders.POST("/:id/cancel", h.CancelOrder)

				// Line item management
				orders.POST("/:id/line-items", h.AddLineItem)
				orders.PUT("/:id/line-items/:lineItemId", h.UpdateLineItem)
				orders.DELETE("/:id/line-items/:lineItemId", h.RemoveLineItem)

				// Order lookup
				orders.GET("/by-number/:orderNumber", h.GetOrderByNumber)
				orders.GET("/customer/:customerId", h.GetOrdersByCustomer)
			}

			// Fulfillment management
			fulfillments := protected.Group("/fulfillments")
			{
				fulfillments.POST("", h.CreateFulfillment)
				fulfillments.PUT("/:id", h.UpdateFulfillment)
				fulfillments.POST("/:id/ship", h.MarkFulfillmentShipped)
			}

			// Transaction management
			transactions := protected.Group("/transactions")
			{
				transactions.POST("", h.CreateTransaction)
				transactions.GET("/order/:orderId", h.GetTransactionsByOrder)
			}

			// Return management
			returns := protected.Group("/returns")
			{
				returns.POST("", h.CreateReturn)
				returns.GET("/:id", h.GetReturn)
				returns.PUT("/:id", h.UpdateReturn)
				returns.POST("/:id/process", h.ProcessReturn)
				returns.GET("/order/:orderId", h.GetReturnsByOrder)
			}

			// Analytics and reporting
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/orders/stats", h.GetOrderStats)
				analytics.GET("/orders/:id/events", h.GetOrderEvents)
			}
		}

		// Public routes for order tracking
		public := v1.Group("/public")
		{
			public.GET("/orders/:orderNumber/track", h.TrackOrder)
		}
	}
}

// Order Management Handlers

// CreateOrder handles creating a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create order")
		httputil.InternalServerError(c, "Failed to create order")
		return
	}

	httputil.Created(c, order, "Order created successfully")
}

// GetOrders handles retrieving orders with filters
func (h *OrderHandler) GetOrders(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	if merchantID == "" {
		httputil.BadRequest(c, "merchant_id is required")
		return
	}

	merchantUUID, err := uuid.Parse(merchantID)
	if err != nil {
		httputil.BadRequest(c, "Invalid merchant ID")
		return
	}

	// Parse filters
	filters := h.parseOrderFilters(c)
	pagination := httputil.GetPaginationParams(c)

	orders, total, err := h.service.GetOrdersByMerchant(c.Request.Context(), merchantUUID, filters, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get orders")
		httputil.InternalServerError(c, "Failed to get orders")
		return
	}

	httputil.SuccessWithMeta(c, orders, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   total,
	}, "Orders retrieved successfully")
}

// GetOrder handles retrieving a specific order
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid order ID")
		return
	}

	order, err := h.service.GetOrder(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to get order")
			httputil.InternalServerError(c, "Failed to get order")
		}
		return
	}

	httputil.Success(c, order, "Order retrieved successfully")
}

// GetOrderByNumber handles retrieving an order by order number
func (h *OrderHandler) GetOrderByNumber(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	if orderNumber == "" {
		httputil.BadRequest(c, "Order number is required")
		return
	}

	order, err := h.service.GetOrderByNumber(c.Request.Context(), orderNumber)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to get order by number")
			httputil.InternalServerError(c, "Failed to get order")
		}
		return
	}

	httputil.Success(c, order, "Order retrieved successfully")
}

// GetOrdersByCustomer handles retrieving orders for a customer
func (h *OrderHandler) GetOrdersByCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid customer ID")
		return
	}

	pagination := httputil.GetPaginationParams(c)
	orders, total, err := h.service.GetOrdersByCustomer(c.Request.Context(), customerID, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get customer orders")
		httputil.InternalServerError(c, "Failed to get customer orders")
		return
	}

	httputil.SuccessWithMeta(c, orders, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   total,
	}, "Customer orders retrieved successfully")
}

// UpdateOrder handles updating an order
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid order ID")
		return
	}

	var req service.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	order, err := h.service.UpdateOrder(c.Request.Context(), id, &req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to update order")
			httputil.InternalServerError(c, "Failed to update order")
		}
		return
	}

	httputil.Success(c, order, "Order updated successfully")
}

// ConfirmOrder handles confirming a pending order
func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid order ID")
		return
	}

	userID, _ := c.Get("user_id")
	var userUUID *uuid.UUID
	if userID != nil {
		if uid, ok := userID.(string); ok {
			if parsed, err := uuid.Parse(uid); err == nil {
				userUUID = &parsed
			}
		}
	}

	if err := h.service.ConfirmOrder(c.Request.Context(), id, userUUID); err != nil {
		switch err {
		case service.ErrOrderNotFound:
			httputil.NotFound(c, "Order not found")
		case service.ErrInvalidStatus:
			httputil.BadRequest(c, "Invalid order status for confirmation")
		default:
			h.logger.WithError(err).Error("Failed to confirm order")
			httputil.InternalServerError(c, "Failed to confirm order")
		}
		return
	}

	httputil.Success(c, nil, "Order confirmed successfully")
}

// CancelOrder handles cancelling an order
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid order ID")
		return
	}

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	var userUUID *uuid.UUID
	if userID != nil {
		if uid, ok := userID.(string); ok {
			if parsed, err := uuid.Parse(uid); err == nil {
				userUUID = &parsed
			}
		}
	}

	if err := h.service.CancelOrder(c.Request.Context(), id, req.Reason, userUUID); err != nil {
		switch err {
		case service.ErrOrderNotFound:
			httputil.NotFound(c, "Order not found")
		case service.ErrOrderAlreadyCancelled:
			httputil.BadRequest(c, "Order is already cancelled")
		case service.ErrOrderNotCancellable:
			httputil.BadRequest(c, "Order cannot be cancelled")
		default:
			h.logger.WithError(err).Error("Failed to cancel order")
			httputil.InternalServerError(c, "Failed to cancel order")
		}
		return
	}

	httputil.Success(c, nil, "Order cancelled successfully")
}

// Line Item Management Handlers

// AddLineItem handles adding a line item to an order
func (h *OrderHandler) AddLineItem(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid order ID")
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

	lineItem, err := h.service.AddLineItem(c.Request.Context(), orderID, &req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to add line item")
			httputil.InternalServerError(c, "Failed to add line item")
		}
		return
	}

	httputil.Created(c, lineItem, "Line item added successfully")
}

// UpdateLineItem handles updating a line item
func (h *OrderHandler) UpdateLineItem(c *gin.Context) {
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
		if err == service.ErrLineItemNotFound {
			httputil.NotFound(c, "Line item not found")
		} else {
			h.logger.WithError(err).Error("Failed to update line item")
			httputil.InternalServerError(c, "Failed to update line item")
		}
		return
	}

	httputil.Success(c, lineItem, "Line item updated successfully")
}

// RemoveLineItem handles removing a line item from an order
func (h *OrderHandler) RemoveLineItem(c *gin.Context) {
	lineItemID, err := uuid.Parse(c.Param("lineItemId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid line item ID")
		return
	}

	if err := h.service.RemoveLineItem(c.Request.Context(), lineItemID); err != nil {
		if err == service.ErrLineItemNotFound {
			httputil.NotFound(c, "Line item not found")
		} else {
			h.logger.WithError(err).Error("Failed to remove line item")
			httputil.InternalServerError(c, "Failed to remove line item")
		}
		return
	}

	httputil.Success(c, nil, "Line item removed successfully")
}

// Fulfillment Management Handlers

// CreateFulfillment handles creating a new fulfillment
func (h *OrderHandler) CreateFulfillment(c *gin.Context) {
	var req service.CreateFulfillmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	fulfillment, err := h.service.CreateFulfillment(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to create fulfillment")
			httputil.InternalServerError(c, "Failed to create fulfillment")
		}
		return
	}

	httputil.Created(c, fulfillment, "Fulfillment created successfully")
}

// UpdateFulfillment handles updating a fulfillment
func (h *OrderHandler) UpdateFulfillment(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Update fulfillment - to be implemented"})
}

// MarkFulfillmentShipped handles marking a fulfillment as shipped
func (h *OrderHandler) MarkFulfillmentShipped(c *gin.Context) {
	fulfillmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid fulfillment ID")
		return
	}

	var req struct {
		TrackingNumber string `json:"tracking_number" validate:"required"`
		TrackingURL    string `json:"tracking_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	if err := h.service.MarkFulfillmentShipped(c.Request.Context(), fulfillmentID, req.TrackingNumber, req.TrackingURL); err != nil {
		if err == service.ErrFulfillmentNotFound {
			httputil.NotFound(c, "Fulfillment not found")
		} else {
			h.logger.WithError(err).Error("Failed to mark fulfillment as shipped")
			httputil.InternalServerError(c, "Failed to mark fulfillment as shipped")
		}
		return
	}

	httputil.Success(c, nil, "Fulfillment marked as shipped successfully")
}

// Transaction Management Handlers

// CreateTransaction handles creating a payment transaction
func (h *OrderHandler) CreateTransaction(c *gin.Context) {
	var req service.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	transaction, err := h.service.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to create transaction")
			httputil.InternalServerError(c, "Failed to create transaction")
		}
		return
	}

	httputil.Created(c, transaction, "Transaction created successfully")
}

// GetTransactionsByOrder handles retrieving transactions for an order
func (h *OrderHandler) GetTransactionsByOrder(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get transactions by order - to be implemented"})
}

// Return Management Handlers

// CreateReturn handles creating a return
func (h *OrderHandler) CreateReturn(c *gin.Context) {
	var req service.CreateReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	returnItem, err := h.service.CreateReturn(c.Request.Context(), &req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to create return")
			httputil.InternalServerError(c, "Failed to create return")
		}
		return
	}

	httputil.Created(c, returnItem, "Return created successfully")
}

// GetReturn handles retrieving a return
func (h *OrderHandler) GetReturn(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get return - to be implemented"})
}

// UpdateReturn handles updating a return
func (h *OrderHandler) UpdateReturn(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Update return - to be implemented"})
}

// ProcessReturn handles processing a return
func (h *OrderHandler) ProcessReturn(c *gin.Context) {
	returnID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid return ID")
		return
	}

	var req struct {
		RefundAmount float64 `json:"refund_amount" validate:"required,min=0"`
		RestockItems bool    `json:"restock_items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	if err := h.service.ProcessReturn(c.Request.Context(), returnID, req.RefundAmount, req.RestockItems); err != nil {
		if err == service.ErrReturnNotFound {
			httputil.NotFound(c, "Return not found")
		} else {
			h.logger.WithError(err).Error("Failed to process return")
			httputil.InternalServerError(c, "Failed to process return")
		}
		return
	}

	httputil.Success(c, nil, "Return processed successfully")
}

// GetReturnsByOrder handles retrieving returns for an order
func (h *OrderHandler) GetReturnsByOrder(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get returns by order - to be implemented"})
}

// Analytics and Reporting Handlers

// GetOrderStats handles retrieving order statistics
func (h *OrderHandler) GetOrderStats(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	if merchantID == "" {
		httputil.BadRequest(c, "merchant_id is required")
		return
	}

	merchantUUID, err := uuid.Parse(merchantID)
	if err != nil {
		httputil.BadRequest(c, "Invalid merchant ID")
		return
	}

	// Parse date range
	dateFrom := c.DefaultQuery("date_from", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	dateTo := c.DefaultQuery("date_to", time.Now().Format("2006-01-02"))

	from, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		httputil.BadRequest(c, "Invalid date_from format (YYYY-MM-DD)")
		return
	}

	to, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		httputil.BadRequest(c, "Invalid date_to format (YYYY-MM-DD)")
		return
	}

	stats, err := h.service.GetOrderStats(c.Request.Context(), merchantUUID, from, to)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get order statistics")
		httputil.InternalServerError(c, "Failed to get order statistics")
		return
	}

	httputil.Success(c, stats, "Order statistics retrieved successfully")
}

// GetOrderEvents handles retrieving events for an order
func (h *OrderHandler) GetOrderEvents(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid order ID")
		return
	}

	events, err := h.service.GetOrderEvents(c.Request.Context(), orderID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get order events")
		httputil.InternalServerError(c, "Failed to get order events")
		return
	}

	httputil.Success(c, events, "Order events retrieved successfully")
}

// Public Handlers

// TrackOrder handles public order tracking
func (h *OrderHandler) TrackOrder(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	if orderNumber == "" {
		httputil.BadRequest(c, "Order number is required")
		return
	}

	order, err := h.service.GetOrderByNumber(c.Request.Context(), orderNumber)
	if err != nil {
		if err == service.ErrOrderNotFound {
			httputil.NotFound(c, "Order not found")
		} else {
			h.logger.WithError(err).Error("Failed to track order")
			httputil.InternalServerError(c, "Failed to track order")
		}
		return
	}

	// Return limited order information for public tracking
	trackingInfo := map[string]interface{}{
		"order_number":       order.OrderNumber,
		"status":             order.Status,
		"fulfillment_status": order.FulfillmentStatus,
		"tracking_number":    order.TrackingNumber,
		"tracking_url":       order.TrackingURL,
		"carrier":            order.Carrier,
		"shipped_at":         order.ShippedAt,
		"delivered_at":       order.DeliveredAt,
		"created_at":         order.CreatedAt,
	}

	httputil.Success(c, trackingInfo, "Order tracking information retrieved successfully")
}

// Helper Methods

// parseOrderFilters parses query parameters into order filters
func (h *OrderHandler) parseOrderFilters(c *gin.Context) map[string]interface{} {
	filters := make(map[string]interface{})

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filters["payment_status"] = paymentStatus
	}
	if fulfillmentStatus := c.Query("fulfillment_status"); fulfillmentStatus != "" {
		filters["fulfillment_status"] = fulfillmentStatus
	}
	if source := c.Query("source"); source != "" {
		filters["source"] = source
	}
	if customerID := c.Query("customer_id"); customerID != "" {
		if id, err := uuid.Parse(customerID); err == nil {
			filters["customer_id"] = id
		}
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filters["date_from"] = t
		}
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			filters["date_to"] = t
		}
	}
	if minTotal := c.Query("min_total"); minTotal != "" {
		if amount, err := strconv.ParseFloat(minTotal, 64); err == nil {
			filters["min_total"] = amount
		}
	}
	if maxTotal := c.Query("max_total"); maxTotal != "" {
		if amount, err := strconv.ParseFloat(maxTotal, 64); err == nil {
			filters["max_total"] = amount
		}
	}

	return filters
}
