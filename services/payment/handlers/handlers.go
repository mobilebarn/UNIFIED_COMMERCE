package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"unified-commerce/services/payment/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// PaymentHandler handles HTTP requests for payment operations
type PaymentHandler struct {
	service   *service.PaymentService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(service *service.PaymentService, logger *logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all payment routes
func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// Protected routes (authentication required)
		authConfig := middleware.DefaultAuthConfig()
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(authConfig))
		{
			// Payment management
			payments := protected.Group("/payments")
			{
				payments.POST("", h.CreatePayment)
				payments.GET("/:id", h.GetPayment)
				payments.GET("/order/:orderId", h.GetPaymentByOrderID)
				payments.GET("/merchant/:merchantId", h.GetPaymentsByMerchant)
				payments.POST("/:id/process", h.ProcessPayment)
				payments.POST("/:id/cancel", h.CancelPayment)
			}

			// Payment method management
			methods := protected.Group("/payment-methods")
			{
				methods.POST("", h.CreatePaymentMethod)
				methods.GET("/:id", h.GetPaymentMethod)
				methods.GET("/customer/:customerId", h.GetPaymentMethodsByCustomer)
				methods.PUT("/:id", h.UpdatePaymentMethod)
				methods.DELETE("/:id", h.DeletePaymentMethod)
			}

			// Gateway management
			gateways := protected.Group("/gateways")
			{
				gateways.POST("", h.CreateGateway)
				gateways.GET("/:id", h.GetGateway)
				gateways.GET("/name/:name", h.GetGatewayByName)
				gateways.GET("", h.GetAllGateways)
				gateways.PUT("/:id", h.UpdateGateway)
			}

			// Refund management
			refunds := protected.Group("/refunds")
			{
				refunds.POST("", h.CreateRefund)
				refunds.GET("/:id", h.GetRefund)
				refunds.GET("/payment/:paymentId", h.GetRefundsByPayment)
				refunds.POST("/:id/process", h.ProcessRefund)
			}

			// Event management
			events := protected.Group("/payment-events")
			{
				events.GET("/payment/:paymentId", h.GetPaymentEvents)
			}

			// Settlement management
			settlements := protected.Group("/settlements")
			{
				settlements.POST("", h.CreateSettlement)
				settlements.GET("/:id", h.GetSettlement)
				settlements.GET("/gateway/:gatewayId", h.GetSettlementsByGateway)
				settlements.PUT("/:id", h.UpdateSettlement)
			}

			// Administrative operations
			admin := protected.Group("/admin")
			{
				admin.GET("/reports/revenue", h.GetRevenueReport)
				admin.GET("/reports/settlements", h.GetSettlementReport)
			}
		}
	}
}

// Payment Management Handlers

// CreatePayment handles creating a new payment
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req service.CreatePaymentRequest
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

	payment, err := h.service.CreatePayment(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create payment")
		httputil.InternalServerError(c, "Failed to create payment")
		return
	}

	httputil.Created(c, payment, "Payment created successfully")
}

// GetPayment handles retrieving a payment by ID
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment ID")
		return
	}

	payment, err := h.service.GetPayment(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			httputil.NotFound(c, "Payment not found")
		default:
			h.logger.WithError(err).Error("Failed to get payment")
			httputil.InternalServerError(c, "Failed to get payment")
		}
		return
	}

	httputil.Success(c, payment, "Payment retrieved successfully")
}

// GetPaymentByOrderID handles retrieving a payment by order ID
func (h *PaymentHandler) GetPaymentByOrderID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("orderId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid order ID")
		return
	}

	payment, err := h.service.GetPaymentByOrderID(c.Request.Context(), orderID)
	if err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			httputil.NotFound(c, "Payment not found")
		default:
			h.logger.WithError(err).Error("Failed to get payment by order ID")
			httputil.InternalServerError(c, "Failed to get payment")
		}
		return
	}

	httputil.Success(c, payment, "Payment retrieved successfully")
}

// GetPaymentsByMerchant handles retrieving payments for a merchant
func (h *PaymentHandler) GetPaymentsByMerchant(c *gin.Context) {
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

	if currency := c.Query("currency"); currency != "" {
		filters["currency"] = currency
	}

	payments, total, err := h.service.GetPaymentsByMerchant(c.Request.Context(), merchantID, filters, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get payments by merchant")
		httputil.InternalServerError(c, "Failed to get payments")
		return
	}

	response := map[string]interface{}{
		"data":        payments,
		"total":       total,
		"page":        pagination.Page,
		"per_page":    pagination.PerPage,
		"total_pages": (total + int64(pagination.PerPage) - 1) / int64(pagination.PerPage),
	}

	httputil.Success(c, response, "Payments retrieved successfully")
}

// ProcessPayment handles processing a payment
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment ID")
		return
	}

	if err := h.service.ProcessPayment(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			httputil.NotFound(c, "Payment not found")
		case service.ErrPaymentAlreadyProcessed:
			httputil.BadRequest(c, "Payment already processed")
		case service.ErrInvalidPaymentStatus:
			httputil.BadRequest(c, "Invalid payment status")
		default:
			h.logger.WithError(err).Error("Failed to process payment")
			httputil.InternalServerError(c, "Failed to process payment")
		}
		return
	}

	httputil.Success(c, nil, "Payment processed successfully")
}

// CancelPayment handles cancelling a payment
func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment ID")
		return
	}

	if err := h.service.CancelPayment(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			httputil.NotFound(c, "Payment not found")
		case service.ErrInvalidPaymentStatus:
			httputil.BadRequest(c, "Payment cannot be cancelled")
		default:
			h.logger.WithError(err).Error("Failed to cancel payment")
			httputil.InternalServerError(c, "Failed to cancel payment")
		}
		return
	}

	httputil.Success(c, nil, "Payment cancelled successfully")
}

// Payment Method Management Handlers

// CreatePaymentMethod handles creating a new payment method
func (h *PaymentHandler) CreatePaymentMethod(c *gin.Context) {
	var req service.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	paymentMethod, err := h.service.CreatePaymentMethod(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create payment method")
		httputil.InternalServerError(c, "Failed to create payment method")
		return
	}

	httputil.Created(c, paymentMethod, "Payment method created successfully")
}

// GetPaymentMethod handles retrieving a payment method by ID
func (h *PaymentHandler) GetPaymentMethod(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment method ID")
		return
	}

	paymentMethod, err := h.service.GetPaymentMethod(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrPaymentMethodNotFound:
			httputil.NotFound(c, "Payment method not found")
		default:
			h.logger.WithError(err).Error("Failed to get payment method")
			httputil.InternalServerError(c, "Failed to get payment method")
		}
		return
	}

	httputil.Success(c, paymentMethod, "Payment method retrieved successfully")
}

// GetPaymentMethodsByCustomer handles retrieving payment methods for a customer
func (h *PaymentHandler) GetPaymentMethodsByCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid customer ID")
		return
	}

	paymentMethods, err := h.service.GetPaymentMethodsByCustomer(c.Request.Context(), customerID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get payment methods by customer")
		httputil.InternalServerError(c, "Failed to get payment methods")
		return
	}

	httputil.Success(c, paymentMethods, "Payment methods retrieved successfully")
}

// UpdatePaymentMethod handles updating a payment method
func (h *PaymentHandler) UpdatePaymentMethod(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment method ID")
		return
	}

	var req service.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	paymentMethod, err := h.service.UpdatePaymentMethod(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrPaymentMethodNotFound:
			httputil.NotFound(c, "Payment method not found")
		default:
			h.logger.WithError(err).Error("Failed to update payment method")
			httputil.InternalServerError(c, "Failed to update payment method")
		}
		return
	}

	httputil.Success(c, paymentMethod, "Payment method updated successfully")
}

// DeletePaymentMethod handles deleting a payment method
func (h *PaymentHandler) DeletePaymentMethod(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment method ID")
		return
	}

	if err := h.service.DeletePaymentMethod(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrPaymentMethodNotFound:
			httputil.NotFound(c, "Payment method not found")
		default:
			h.logger.WithError(err).Error("Failed to delete payment method")
			httputil.InternalServerError(c, "Failed to delete payment method")
		}
		return
	}

	httputil.Success(c, nil, "Payment method deleted successfully")
}

// Gateway Management Handlers

// CreateGateway handles creating a new payment gateway
func (h *PaymentHandler) CreateGateway(c *gin.Context) {
	var req service.CreateGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	gateway, err := h.service.CreateGateway(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create payment gateway")
		httputil.InternalServerError(c, "Failed to create payment gateway")
		return
	}

	httputil.Created(c, gateway, "Payment gateway created successfully")
}

// GetGateway handles retrieving a payment gateway by ID
func (h *PaymentHandler) GetGateway(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid gateway ID")
		return
	}

	gateway, err := h.service.GetGateway(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrGatewayNotFound:
			httputil.NotFound(c, "Payment gateway not found")
		default:
			h.logger.WithError(err).Error("Failed to get payment gateway")
			httputil.InternalServerError(c, "Failed to get payment gateway")
		}
		return
	}

	httputil.Success(c, gateway, "Payment gateway retrieved successfully")
}

// GetGatewayByName handles retrieving a payment gateway by name
func (h *PaymentHandler) GetGatewayByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		httputil.BadRequest(c, "Gateway name is required")
		return
	}

	gateway, err := h.service.GetGatewayByName(c.Request.Context(), name)
	if err != nil {
		switch err {
		case service.ErrGatewayNotFound:
			httputil.NotFound(c, "Payment gateway not found")
		default:
			h.logger.WithError(err).Error("Failed to get payment gateway by name")
			httputil.InternalServerError(c, "Failed to get payment gateway")
		}
		return
	}

	httputil.Success(c, gateway, "Payment gateway retrieved successfully")
}

// GetAllGateways handles retrieving all payment gateways
func (h *PaymentHandler) GetAllGateways(c *gin.Context) {
	gateways, err := h.service.GetAllGateways(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get payment gateways")
		httputil.InternalServerError(c, "Failed to get payment gateways")
		return
	}

	httputil.Success(c, gateways, "Payment gateways retrieved successfully")
}

// UpdateGateway handles updating a payment gateway
func (h *PaymentHandler) UpdateGateway(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid gateway ID")
		return
	}

	var req service.UpdateGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	gateway, err := h.service.UpdateGateway(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrGatewayNotFound:
			httputil.NotFound(c, "Payment gateway not found")
		default:
			h.logger.WithError(err).Error("Failed to update payment gateway")
			httputil.InternalServerError(c, "Failed to update payment gateway")
		}
		return
	}

	httputil.Success(c, gateway, "Payment gateway updated successfully")
}

// Refund Management Handlers

// CreateRefund handles creating a new refund
func (h *PaymentHandler) CreateRefund(c *gin.Context) {
	var req service.CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	refund, err := h.service.CreateRefund(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrPaymentNotFound:
			httputil.NotFound(c, "Payment not found")
		case service.ErrInvalidPaymentStatus:
			httputil.BadRequest(c, "Payment cannot be refunded")
		case service.ErrRefundAmountExceeds:
			httputil.BadRequest(c, "Refund amount exceeds payment amount")
		default:
			h.logger.WithError(err).Error("Failed to create refund")
			httputil.InternalServerError(c, "Failed to create refund")
		}
		return
	}

	httputil.Created(c, refund, "Refund created successfully")
}

// GetRefund handles retrieving a refund by ID
func (h *PaymentHandler) GetRefund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid refund ID")
		return
	}

	refund, err := h.service.GetRefund(c.Request.Context(), id)
	if err != nil {
		switch err {
		case service.ErrRefundNotFound:
			httputil.NotFound(c, "Refund not found")
		default:
			h.logger.WithError(err).Error("Failed to get refund")
			httputil.InternalServerError(c, "Failed to get refund")
		}
		return
	}

	httputil.Success(c, refund, "Refund retrieved successfully")
}

// GetRefundsByPayment handles retrieving refunds for a payment
func (h *PaymentHandler) GetRefundsByPayment(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("paymentId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment ID")
		return
	}

	refunds, err := h.service.GetRefundsByPayment(c.Request.Context(), paymentID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get refunds by payment")
		httputil.InternalServerError(c, "Failed to get refunds")
		return
	}

	httputil.Success(c, refunds, "Refunds retrieved successfully")
}

// ProcessRefund handles processing a refund
func (h *PaymentHandler) ProcessRefund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid refund ID")
		return
	}

	if err := h.service.ProcessRefund(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrRefundNotFound:
			httputil.NotFound(c, "Refund not found")
		case service.ErrInvalidRefundStatus:
			httputil.BadRequest(c, "Refund cannot be processed")
		default:
			h.logger.WithError(err).Error("Failed to process refund")
			httputil.InternalServerError(c, "Failed to process refund")
		}
		return
	}

	httputil.Success(c, nil, "Refund processed successfully")
}

// Event Management Handlers

// GetPaymentEvents handles retrieving events for a payment
func (h *PaymentHandler) GetPaymentEvents(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("paymentId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid payment ID")
		return
	}

	events, err := h.service.GetPaymentEvents(c.Request.Context(), paymentID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get payment events")
		httputil.InternalServerError(c, "Failed to get payment events")
		return
	}

	httputil.Success(c, events, "Payment events retrieved successfully")
}

// Settlement Management Handlers

// CreateSettlement handles creating a new settlement
func (h *PaymentHandler) CreateSettlement(c *gin.Context) {
	var req service.CreateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	settlement, err := h.service.CreateSettlement(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create settlement")
		httputil.InternalServerError(c, "Failed to create settlement")
		return
	}

	httputil.Created(c, settlement, "Settlement created successfully")
}

// GetSettlement handles retrieving a settlement by ID
func (h *PaymentHandler) GetSettlement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid settlement ID")
		return
	}

	settlement, err := h.service.GetSettlement(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get settlement")
		httputil.InternalServerError(c, "Failed to get settlement")
		return
	}

	httputil.Success(c, settlement, "Settlement retrieved successfully")
}

// GetSettlementsByGateway handles retrieving settlements for a gateway
func (h *PaymentHandler) GetSettlementsByGateway(c *gin.Context) {
	gatewayID, err := uuid.Parse(c.Param("gatewayId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid gateway ID")
		return
	}

	pagination := httputil.GetPaginationParams(c)
	settlements, total, err := h.service.GetSettlementsByGateway(c.Request.Context(), gatewayID, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get settlements by gateway")
		httputil.InternalServerError(c, "Failed to get settlements")
		return
	}

	response := map[string]interface{}{
		"data":        settlements,
		"total":       total,
		"page":        pagination.Page,
		"per_page":    pagination.PerPage,
		"total_pages": (total + int64(pagination.PerPage) - 1) / int64(pagination.PerPage),
	}

	httputil.Success(c, response, "Settlements retrieved successfully")
}

// UpdateSettlement handles updating a settlement
func (h *PaymentHandler) UpdateSettlement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid settlement ID")
		return
	}

	var req service.UpdateSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	settlement, err := h.service.UpdateSettlement(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to update settlement")
		httputil.InternalServerError(c, "Failed to update settlement")
		return
	}

	httputil.Success(c, settlement, "Settlement updated successfully")
}

// Administrative Operations Handlers

// GetRevenueReport handles retrieving revenue report
func (h *PaymentHandler) GetRevenueReport(c *gin.Context) {
	// In a real implementation, this would generate a revenue report
	httputil.Success(c, nil, "Revenue report retrieved successfully")
}

// GetSettlementReport handles retrieving settlement report
func (h *PaymentHandler) GetSettlementReport(c *gin.Context) {
	// In a real implementation, this would generate a settlement report
	httputil.Success(c, nil, "Settlement report retrieved successfully")
}
