package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"unified-commerce/services/inventory/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// InventoryHandler handles HTTP requests for inventory operations
type InventoryHandler struct {
	service   *service.InventoryService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(service *service.InventoryService, logger *logger.Logger) *InventoryHandler {
	return &InventoryHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all inventory routes
func (h *InventoryHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		// Protected routes (authentication required)
		authConfig := middleware.DefaultAuthConfig()
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(authConfig))
		{
			// Location management
			locations := protected.Group("/locations")
			{
				locations.POST("", h.CreateLocation)
				locations.GET("", h.GetLocations)
				locations.GET("/:id", h.GetLocation)
				locations.PUT("/:id", h.UpdateLocation)
				locations.DELETE("/:id", h.DeleteLocation)
				locations.GET("/:id/summary", h.GetLocationSummary)
			}

			// Inventory items management
			inventory := protected.Group("/inventory")
			{
				inventory.POST("", h.CreateInventoryItem)
				inventory.GET("", h.GetInventoryItems)
				inventory.GET("/:id", h.GetInventoryItem)
				inventory.PUT("/:id", h.UpdateInventoryItem)
				inventory.POST("/:id/adjust", h.AdjustInventory)
				inventory.GET("/low-stock", h.GetLowStockItems)
				inventory.GET("/by-product/:productId", h.GetInventoryByProduct)
				inventory.GET("/check-availability", h.CheckStockAvailability)
			}

			// Stock movements
			movements := protected.Group("/stock-movements")
			{
				movements.GET("", h.GetStockMovements)
			}

			// Stock reservations
			reservations := protected.Group("/reservations")
			{
				reservations.POST("", h.CreateReservation)
				reservations.GET("/:id", h.GetReservation)
				reservations.GET("", h.GetReservationsByReference)
				reservations.POST("/:id/fulfill", h.FulfillReservation)
				reservations.POST("/:id/cancel", h.CancelReservation)
			}

			// Stock transfers (future implementation)
			transfers := protected.Group("/transfers")
			{
				transfers.POST("", h.CreateTransfer)
				transfers.GET("", h.GetTransfers)
				transfers.GET("/:id", h.GetTransfer)
				transfers.PUT("/:id", h.UpdateTransfer)
				transfers.POST("/:id/ship", h.ShipTransfer)
				transfers.POST("/:id/receive", h.ReceiveTransfer)
			}
		}
	}
}

// Location Handlers

// CreateLocation handles creating a new location
func (h *InventoryHandler) CreateLocation(c *gin.Context) {
	var req service.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	location, err := h.service.CreateLocation(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create location")
		httputil.InternalServerError(c, "Failed to create location")
		return
	}

	httputil.Created(c, location, "Location created successfully")
}

// GetLocations handles retrieving all locations for a merchant
func (h *InventoryHandler) GetLocations(c *gin.Context) {
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

	pagination := httputil.GetPaginationParams(c)
	locations, total, err := h.service.GetLocationsByMerchant(c.Request.Context(), merchantUUID, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get locations")
		httputil.InternalServerError(c, "Failed to get locations")
		return
	}

	httputil.SuccessWithMeta(c, locations, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   total,
	}, "Locations retrieved successfully")
}

// GetLocation handles retrieving a specific location
func (h *InventoryHandler) GetLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid location ID")
		return
	}

	location, err := h.service.GetLocation(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrLocationNotFound {
			httputil.NotFound(c, "Location not found")
		} else {
			h.logger.WithError(err).Error("Failed to get location")
			httputil.InternalServerError(c, "Failed to get location")
		}
		return
	}

	httputil.Success(c, location, "Location retrieved successfully")
}

// UpdateLocation handles updating a location
func (h *InventoryHandler) UpdateLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid location ID")
		return
	}

	var req service.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	location, err := h.service.UpdateLocation(c.Request.Context(), id, &req)
	if err != nil {
		if err == service.ErrLocationNotFound {
			httputil.NotFound(c, "Location not found")
		} else {
			h.logger.WithError(err).Error("Failed to update location")
			httputil.InternalServerError(c, "Failed to update location")
		}
		return
	}

	httputil.Success(c, location, "Location updated successfully")
}

// DeleteLocation handles deleting a location
func (h *InventoryHandler) DeleteLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid location ID")
		return
	}

	if err := h.service.DeleteLocation(c.Request.Context(), id); err != nil {
		if err == service.ErrLocationNotFound {
			httputil.NotFound(c, "Location not found")
		} else {
			h.logger.WithError(err).Error("Failed to delete location")
			httputil.InternalServerError(c, "Failed to delete location")
		}
		return
	}

	httputil.Success(c, nil, "Location deleted successfully")
}

// GetLocationSummary handles retrieving inventory summary for a location
func (h *InventoryHandler) GetLocationSummary(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid location ID")
		return
	}

	summary, err := h.service.GetInventorySummaryByLocation(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get location summary")
		httputil.InternalServerError(c, "Failed to get location summary")
		return
	}

	httputil.Success(c, summary, "Location summary retrieved successfully")
}

// Inventory Item Handlers

// CreateInventoryItem handles creating a new inventory item
func (h *InventoryHandler) CreateInventoryItem(c *gin.Context) {
	var req service.CreateInventoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	item, err := h.service.CreateInventoryItem(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrLocationNotFound:
			httputil.NotFound(c, "Location not found")
		case service.ErrLocationInactive:
			httputil.BadRequest(c, "Location is inactive")
		default:
			h.logger.WithError(err).Error("Failed to create inventory item")
			httputil.InternalServerError(c, "Failed to create inventory item")
		}
		return
	}

	httputil.Created(c, item, "Inventory item created successfully")
}

// GetInventoryItems handles retrieving inventory items
func (h *InventoryHandler) GetInventoryItems(c *gin.Context) {
	locationID := c.Query("location_id")
	if locationID == "" {
		httputil.BadRequest(c, "location_id is required")
		return
	}

	locationUUID, err := uuid.Parse(locationID)
	if err != nil {
		httputil.BadRequest(c, "Invalid location ID")
		return
	}

	pagination := httputil.GetPaginationParams(c)
	items, total, err := h.service.GetInventoryItemsByLocation(c.Request.Context(), locationUUID, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get inventory items")
		httputil.InternalServerError(c, "Failed to get inventory items")
		return
	}

	httputil.SuccessWithMeta(c, items, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   total,
	}, "Inventory items retrieved successfully")
}

// GetInventoryItem handles retrieving a specific inventory item
func (h *InventoryHandler) GetInventoryItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid inventory item ID")
		return
	}

	item, err := h.service.GetInventoryItem(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrInventoryItemNotFound {
			httputil.NotFound(c, "Inventory item not found")
		} else {
			h.logger.WithError(err).Error("Failed to get inventory item")
			httputil.InternalServerError(c, "Failed to get inventory item")
		}
		return
	}

	httputil.Success(c, item, "Inventory item retrieved successfully")
}

// UpdateInventoryItem handles updating an inventory item
func (h *InventoryHandler) UpdateInventoryItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid inventory item ID")
		return
	}

	var req service.UpdateInventoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	item, err := h.service.UpdateInventoryItem(c.Request.Context(), id, &req)
	if err != nil {
		if err == service.ErrInventoryItemNotFound {
			httputil.NotFound(c, "Inventory item not found")
		} else {
			h.logger.WithError(err).Error("Failed to update inventory item")
			httputil.InternalServerError(c, "Failed to update inventory item")
		}
		return
	}

	httputil.Success(c, item, "Inventory item updated successfully")
}

// AdjustInventory handles adjusting inventory quantity
func (h *InventoryHandler) AdjustInventory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid inventory item ID")
		return
	}

	var req service.AdjustInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	item, err := h.service.AdjustInventory(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case service.ErrInventoryItemNotFound:
			httputil.NotFound(c, "Inventory item not found")
		case service.ErrInvalidQuantity:
			httputil.BadRequest(c, "Invalid quantity")
		default:
			h.logger.WithError(err).Error("Failed to adjust inventory")
			httputil.InternalServerError(c, "Failed to adjust inventory")
		}
		return
	}

	httputil.Success(c, item, "Inventory adjusted successfully")
}

// GetLowStockItems handles retrieving low stock items
func (h *InventoryHandler) GetLowStockItems(c *gin.Context) {
	var locationID *uuid.UUID
	if locID := c.Query("location_id"); locID != "" {
		id, err := uuid.Parse(locID)
		if err != nil {
			httputil.BadRequest(c, "Invalid location ID")
			return
		}
		locationID = &id
	}

	pagination := httputil.GetPaginationParams(c)
	items, total, err := h.service.GetLowStockItems(c.Request.Context(), locationID, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get low stock items")
		httputil.InternalServerError(c, "Failed to get low stock items")
		return
	}

	httputil.SuccessWithMeta(c, items, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   total,
	}, "Low stock items retrieved successfully")
}

// GetInventoryByProduct handles retrieving inventory for a specific product
func (h *InventoryHandler) GetInventoryByProduct(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		httputil.BadRequest(c, "Invalid product ID")
		return
	}

	items, err := h.service.GetInventoryItemsByProduct(c.Request.Context(), productID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get inventory by product")
		httputil.InternalServerError(c, "Failed to get inventory by product")
		return
	}

	httputil.Success(c, items, "Product inventory retrieved successfully")
}

// CheckStockAvailability handles checking stock availability
func (h *InventoryHandler) CheckStockAvailability(c *gin.Context) {
	sku := c.Query("sku")
	locationID := c.Query("location_id")
	quantityStr := c.Query("quantity")

	if sku == "" || locationID == "" || quantityStr == "" {
		httputil.BadRequest(c, "sku, location_id, and quantity are required")
		return
	}

	locationUUID, err := uuid.Parse(locationID)
	if err != nil {
		httputil.BadRequest(c, "Invalid location ID")
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		httputil.BadRequest(c, "Invalid quantity")
		return
	}

	available, availableQuantity, err := h.service.CheckStockAvailability(c.Request.Context(), sku, locationUUID, quantity)
	if err != nil {
		h.logger.WithError(err).Error("Failed to check stock availability")
		httputil.InternalServerError(c, "Failed to check stock availability")
		return
	}

	result := map[string]interface{}{
		"sku":                sku,
		"location_id":        locationID,
		"requested_quantity": quantity,
		"available_quantity": availableQuantity,
		"is_available":       available,
	}

	httputil.Success(c, result, "Stock availability checked successfully")
}

// Stock Movement Handlers

// GetStockMovements handles retrieving stock movements
func (h *InventoryHandler) GetStockMovements(c *gin.Context) {
	filters := make(map[string]interface{})

	// Parse query parameters for filters
	if locationID := c.Query("location_id"); locationID != "" {
		if id, err := uuid.Parse(locationID); err == nil {
			filters["location_id"] = id
		}
	}
	if productID := c.Query("product_id"); productID != "" {
		if id, err := uuid.Parse(productID); err == nil {
			filters["product_id"] = id
		}
	}
	if sku := c.Query("sku"); sku != "" {
		filters["sku"] = sku
	}
	if movementType := c.Query("type"); movementType != "" {
		filters["type"] = movementType
	}
	if reason := c.Query("reason"); reason != "" {
		filters["reason"] = reason
	}
	if reference := c.Query("reference"); reference != "" {
		filters["reference"] = reference
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

	pagination := httputil.GetPaginationParams(c)
	movements, total, err := h.service.GetStockMovements(c.Request.Context(), filters, pagination.Page, pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stock movements")
		httputil.InternalServerError(c, "Failed to get stock movements")
		return
	}

	httputil.SuccessWithMeta(c, movements, &httputil.MetaInfo{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Total:   total,
	}, "Stock movements retrieved successfully")
}

// Stock Reservation Handlers

// CreateReservation handles creating a new stock reservation
func (h *InventoryHandler) CreateReservation(c *gin.Context) {
	var req service.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	reservation, err := h.service.CreateReservation(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInventoryItemNotFound:
			httputil.NotFound(c, "Inventory item not found")
		case service.ErrInsufficientStock:
			httputil.BadRequest(c, "Insufficient stock available")
		default:
			h.logger.WithError(err).Error("Failed to create reservation")
			httputil.InternalServerError(c, "Failed to create reservation")
		}
		return
	}

	httputil.Created(c, reservation, "Stock reservation created successfully")
}

// GetReservation handles retrieving a specific reservation
func (h *InventoryHandler) GetReservation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid reservation ID")
		return
	}

	reservation, err := h.service.GetReservation(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrReservationNotFound {
			httputil.NotFound(c, "Reservation not found")
		} else {
			h.logger.WithError(err).Error("Failed to get reservation")
			httputil.InternalServerError(c, "Failed to get reservation")
		}
		return
	}

	httputil.Success(c, reservation, "Reservation retrieved successfully")
}

// GetReservationsByReference handles retrieving reservations by reference
func (h *InventoryHandler) GetReservationsByReference(c *gin.Context) {
	reference := c.Query("reference")
	if reference == "" {
		httputil.BadRequest(c, "reference is required")
		return
	}

	reservations, err := h.service.GetReservationsByReference(c.Request.Context(), reference)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get reservations by reference")
		httputil.InternalServerError(c, "Failed to get reservations")
		return
	}

	httputil.Success(c, reservations, "Reservations retrieved successfully")
}

// FulfillReservation handles fulfilling a stock reservation
func (h *InventoryHandler) FulfillReservation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid reservation ID")
		return
	}

	var req struct {
		ActualQuantity int        `json:"actual_quantity" validate:"required,min=1"`
		UserID         *uuid.UUID `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	if err := h.service.FulfillReservation(c.Request.Context(), id, req.ActualQuantity, req.UserID); err != nil {
		switch err {
		case service.ErrReservationNotFound:
			httputil.NotFound(c, "Reservation not found")
		case service.ErrInvalidQuantity:
			httputil.BadRequest(c, "Invalid quantity")
		default:
			h.logger.WithError(err).Error("Failed to fulfill reservation")
			httputil.InternalServerError(c, "Failed to fulfill reservation")
		}
		return
	}

	httputil.Success(c, nil, "Reservation fulfilled successfully")
}

// CancelReservation handles cancelling a stock reservation
func (h *InventoryHandler) CancelReservation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.BadRequest(c, "Invalid reservation ID")
		return
	}

	if err := h.service.CancelReservation(c.Request.Context(), id); err != nil {
		if err == service.ErrReservationNotFound {
			httputil.NotFound(c, "Reservation not found")
		} else {
			h.logger.WithError(err).Error("Failed to cancel reservation")
			httputil.InternalServerError(c, "Failed to cancel reservation")
		}
		return
	}

	httputil.Success(c, nil, "Reservation cancelled successfully")
}

// Stock Transfer Handlers (placeholder implementations)

// CreateTransfer handles creating a new stock transfer
func (h *InventoryHandler) CreateTransfer(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Stock transfer creation - to be implemented"})
}

// GetTransfers handles retrieving stock transfers
func (h *InventoryHandler) GetTransfers(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Stock transfers listing - to be implemented"})
}

// GetTransfer handles retrieving a specific stock transfer
func (h *InventoryHandler) GetTransfer(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Stock transfer retrieval - to be implemented"})
}

// UpdateTransfer handles updating a stock transfer
func (h *InventoryHandler) UpdateTransfer(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Stock transfer update - to be implemented"})
}

// ShipTransfer handles shipping a stock transfer
func (h *InventoryHandler) ShipTransfer(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Stock transfer shipping - to be implemented"})
}

// ReceiveTransfer handles receiving a stock transfer
func (h *InventoryHandler) ReceiveTransfer(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Stock transfer receiving - to be implemented"})
}
