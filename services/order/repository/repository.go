package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"unified-commerce/services/order/models"
	"unified-commerce/services/shared/logger"
)

// OrderRepository handles database operations for order management
type OrderRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB, logger *logger.Logger) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger,
	}
}

// Order Operations

// CreateOrder creates a new order with line items
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the order
		if err := tx.Create(order).Error; err != nil {
			r.logger.WithError(err).Error("Failed to create order")
			return err
		}

		// Create order event
		event := &models.OrderEvent{
			OrderID:     order.ID,
			EventType:   models.OrderEventCreated,
			Description: fmt.Sprintf("Order %s created", order.OrderNumber),
		}
		if err := tx.Create(event).Error; err != nil {
			r.logger.WithError(err).Error("Failed to create order event")
			return err
		}

		return nil
	})
}

// GetOrder retrieves an order by ID with all related data
func (r *OrderRepository) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Preload("Fulfillments").
		Preload("Fulfillments.LineItems").
		Preload("Transactions").
		Preload("Returns").
		Preload("Returns.LineItems").
		First(&order, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get order")
		return nil, err
	}
	return &order, nil
}

// GetOrderByNumber retrieves an order by order number
func (r *OrderRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Preload("Fulfillments").
		Preload("Transactions").
		Preload("Returns").
		First(&order, "order_number = ?", orderNumber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get order by number")
		return nil, err
	}
	return &order, nil
}

// GetOrdersByMerchant retrieves orders for a merchant with pagination
func (r *OrderRepository) GetOrdersByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Order{}).Where("merchant_id = ?", merchantID)

	// Apply filters
	query = r.applyOrderFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count orders")
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Where("merchant_id = ?", merchantID).
		Scopes(r.applyOrderFiltersScope(filters)).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get orders")
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrdersByCustomer retrieves orders for a customer
func (r *OrderRepository) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("customer_id = ?", customerID).Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count customer orders")
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get customer orders")
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateOrder updates an order
func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	if err := r.db.WithContext(ctx).Save(order).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update order")
		return err
	}
	return nil
}

// UpdateOrderStatus updates order status and creates an event
func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus, description string, userID *uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order status
		if err := tx.Model(&models.Order{}).
			Where("id = ?", orderID).
			Update("status", status).Error; err != nil {
			return err
		}

		// Create order event
		event := &models.OrderEvent{
			OrderID:     orderID,
			EventType:   r.getEventTypeFromStatus(status),
			Description: description,
			UserID:      userID,
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		return nil
	})
}

// CancelOrder cancels an order and restocks inventory
func (r *OrderRepository) CancelOrder(ctx context.Context, orderID uuid.UUID, reason string, userID *uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update order
		if err := tx.Model(&models.Order{}).
			Where("id = ?", orderID).
			Updates(map[string]interface{}{
				"status":       models.OrderStatusCancelled,
				"cancelled_at": &now,
			}).Error; err != nil {
			return err
		}

		// Create cancellation event
		event := &models.OrderEvent{
			OrderID:     orderID,
			EventType:   models.OrderEventCancelled,
			Description: fmt.Sprintf("Order cancelled: %s", reason),
			UserID:      userID,
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		return nil
	})
}

// Line Item Operations

// AddLineItem adds a new line item to an order
func (r *OrderRepository) AddLineItem(ctx context.Context, lineItem *models.OrderLineItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create line item
		if err := tx.Create(lineItem).Error; err != nil {
			return err
		}

		// Update order totals
		if err := r.recalculateOrderTotals(tx, lineItem.OrderID); err != nil {
			return err
		}

		return nil
	})
}

// UpdateLineItem updates a line item
func (r *OrderRepository) UpdateLineItem(ctx context.Context, lineItem *models.OrderLineItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update line item
		if err := tx.Save(lineItem).Error; err != nil {
			return err
		}

		// Update order totals
		if err := r.recalculateOrderTotals(tx, lineItem.OrderID); err != nil {
			return err
		}

		return nil
	})
}

// RemoveLineItem removes a line item from an order
func (r *OrderRepository) RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get line item to get order ID
		var lineItem models.OrderLineItem
		if err := tx.First(&lineItem, "id = ?", lineItemID).Error; err != nil {
			return err
		}

		// Delete line item
		if err := tx.Delete(&lineItem).Error; err != nil {
			return err
		}

		// Update order totals
		if err := r.recalculateOrderTotals(tx, lineItem.OrderID); err != nil {
			return err
		}

		return nil
	})
}

// Fulfillment Operations

// CreateFulfillment creates a new fulfillment
func (r *OrderRepository) CreateFulfillment(ctx context.Context, fulfillment *models.Fulfillment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create fulfillment
		if err := tx.Create(fulfillment).Error; err != nil {
			return err
		}

		// Update line item fulfillment status
		for _, fulfillmentLineItem := range fulfillment.LineItems {
			if err := tx.Model(&models.OrderLineItem{}).
				Where("id = ?", fulfillmentLineItem.LineItemID).
				Updates(map[string]interface{}{
					"fulfilled_quantity": gorm.Expr("fulfilled_quantity + ?", fulfillmentLineItem.Quantity),
					"fulfillment_status": r.calculateLineItemFulfillmentStatus(tx, fulfillmentLineItem.LineItemID),
				}).Error; err != nil {
				return err
			}
		}

		// Update order fulfillment status
		if err := r.updateOrderFulfillmentStatus(tx, fulfillment.OrderID); err != nil {
			return err
		}

		// Create fulfillment event
		event := &models.OrderEvent{
			OrderID:     fulfillment.OrderID,
			EventType:   models.OrderEventFulfilled,
			Description: fmt.Sprintf("Fulfillment created with tracking number: %s", fulfillment.TrackingNumber),
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateFulfillment updates a fulfillment
func (r *OrderRepository) UpdateFulfillment(ctx context.Context, fulfillment *models.Fulfillment) error {
	if err := r.db.WithContext(ctx).Save(fulfillment).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update fulfillment")
		return err
	}
	return nil
}

// MarkFulfillmentShipped marks a fulfillment as shipped
func (r *OrderRepository) MarkFulfillmentShipped(ctx context.Context, fulfillmentID uuid.UUID, trackingNumber, trackingURL string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update fulfillment
		if err := tx.Model(&models.Fulfillment{}).
			Where("id = ?", fulfillmentID).
			Updates(map[string]interface{}{
				"shipment_status": models.ShipmentStatusInTransit,
				"tracking_number": trackingNumber,
				"tracking_url":    trackingURL,
				"shipped_at":      &now,
			}).Error; err != nil {
			return err
		}

		// Get fulfillment to get order ID
		var fulfillment models.Fulfillment
		if err := tx.First(&fulfillment, "id = ?", fulfillmentID).Error; err != nil {
			return err
		}

		// Update order if all fulfillments are shipped
		if err := r.updateOrderShippingStatus(tx, fulfillment.OrderID); err != nil {
			return err
		}

		// Create shipped event
		event := &models.OrderEvent{
			OrderID:     fulfillment.OrderID,
			EventType:   models.OrderEventShipped,
			Description: fmt.Sprintf("Order shipped with tracking number: %s", trackingNumber),
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		return nil
	})
}

// Transaction Operations

// CreateTransaction creates a new payment transaction
func (r *OrderRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create transaction
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// Update order payment status
		if err := r.updateOrderPaymentStatus(tx, transaction.OrderID); err != nil {
			return err
		}

		// Create payment event
		eventType := models.OrderEventPaymentAuthorized
		if transaction.Kind == models.TransactionKindCapture || transaction.Kind == models.TransactionKindSale {
			eventType = models.OrderEventPaymentCaptured
		}

		event := &models.OrderEvent{
			OrderID:     transaction.OrderID,
			EventType:   eventType,
			Description: fmt.Sprintf("Payment %s: %s %.2f", transaction.Kind, transaction.Currency, transaction.Amount),
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		return nil
	})
}

// Return Operations

// CreateReturn creates a new return
func (r *OrderRepository) CreateReturn(ctx context.Context, returnItem *models.Return) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create return
		if err := tx.Create(returnItem).Error; err != nil {
			return err
		}

		// Create return event
		event := &models.OrderEvent{
			OrderID:     returnItem.OrderID,
			EventType:   models.OrderEventReturned,
			Description: fmt.Sprintf("Return %s requested", returnItem.ReturnNumber),
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		return nil
	})
}

// ProcessReturn processes a return and restocks inventory if needed
func (r *OrderRepository) ProcessReturn(ctx context.Context, returnID uuid.UUID, refundAmount float64, restockItems bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update return
		if err := tx.Model(&models.Return{}).
			Where("id = ?", returnID).
			Updates(map[string]interface{}{
				"status":        models.ReturnStatusProcessed,
				"refund_amount": refundAmount,
				"restock_items": restockItems,
				"processed_at":  &now,
			}).Error; err != nil {
			return err
		}

		// Get return to get order ID
		var returnItem models.Return
		if err := tx.Preload("LineItems").First(&returnItem, "id = ?", returnID).Error; err != nil {
			return err
		}

		// Update order status if all items are returned
		if err := r.updateOrderReturnStatus(tx, returnItem.OrderID); err != nil {
			return err
		}

		return nil
	})
}

// Event Operations

// CreateOrderEvent creates a new order event
func (r *OrderRepository) CreateOrderEvent(ctx context.Context, event *models.OrderEvent) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create order event")
		return err
	}
	return nil
}

// GetOrderEvents retrieves events for an order
func (r *OrderRepository) GetOrderEvents(ctx context.Context, orderID uuid.UUID) ([]*models.OrderEvent, error) {
	var events []*models.OrderEvent
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get order events")
		return nil, err
	}
	return events, nil
}

// Analytics Operations

// GetOrderStats retrieves order statistics for a merchant
func (r *OrderRepository) GetOrderStats(ctx context.Context, merchantID uuid.UUID, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total orders
	var totalOrders int64
	if err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("merchant_id = ? AND created_at BETWEEN ? AND ?", merchantID, dateFrom, dateTo).
		Count(&totalOrders).Error; err != nil {
		return nil, err
	}
	stats["total_orders"] = totalOrders

	// Total revenue
	var totalRevenue float64
	if err := r.db.WithContext(ctx).Model(&models.Order{}).
		Where("merchant_id = ? AND created_at BETWEEN ? AND ? AND status NOT IN ?",
			merchantID, dateFrom, dateTo, []string{"cancelled", "refunded"}).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, err
	}
	stats["total_revenue"] = totalRevenue

	// Average order value
	if totalOrders > 0 {
		stats["average_order_value"] = totalRevenue / float64(totalOrders)
	} else {
		stats["average_order_value"] = 0
	}

	// Orders by status
	statusStats := make(map[string]int64)
	var statusResults []struct {
		Status string
		Count  int64
	}
	if err := r.db.WithContext(ctx).Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Where("merchant_id = ? AND created_at BETWEEN ? AND ?", merchantID, dateFrom, dateTo).
		Group("status").
		Scan(&statusResults).Error; err != nil {
		return nil, err
	}
	for _, result := range statusResults {
		statusStats[result.Status] = result.Count
	}
	stats["orders_by_status"] = statusStats

	return stats, nil
}

// Helper functions

func (r *OrderRepository) applyOrderFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "status":
			query = query.Where("status = ?", value)
		case "payment_status":
			query = query.Where("payment_status = ?", value)
		case "fulfillment_status":
			query = query.Where("fulfillment_status = ?", value)
		case "source":
			query = query.Where("source = ?", value)
		case "customer_id":
			query = query.Where("customer_id = ?", value)
		case "date_from":
			query = query.Where("created_at >= ?", value)
		case "date_to":
			query = query.Where("created_at <= ?", value)
		case "min_total":
			query = query.Where("total_price >= ?", value)
		case "max_total":
			query = query.Where("total_price <= ?", value)
		}
	}
	return query
}

func (r *OrderRepository) applyOrderFiltersScope(filters map[string]interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return r.applyOrderFilters(db, filters)
	}
}

func (r *OrderRepository) recalculateOrderTotals(tx *gorm.DB, orderID uuid.UUID) error {
	var subtotal float64
	if err := tx.Model(&models.OrderLineItem{}).
		Where("order_id = ?", orderID).
		Select("COALESCE(SUM(line_price), 0)").
		Scan(&subtotal).Error; err != nil {
		return err
	}

	// Update order totals (simplified - in real implementation would include tax calculation)
	return tx.Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"subtotal_price": subtotal,
			"total_price":    subtotal, // Simplified
		}).Error
}

func (r *OrderRepository) calculateLineItemFulfillmentStatus(tx *gorm.DB, lineItemID uuid.UUID) models.FulfillmentStatus {
	var lineItem models.OrderLineItem
	if err := tx.First(&lineItem, "id = ?", lineItemID).Error; err != nil {
		return models.FulfillmentStatusUnfulfilled
	}

	if lineItem.FulfilledQuantity == 0 {
		return models.FulfillmentStatusUnfulfilled
	} else if lineItem.FulfilledQuantity < lineItem.Quantity {
		return models.FulfillmentStatusPartiallyFulfilled
	} else {
		return models.FulfillmentStatusFulfilled
	}
}

func (r *OrderRepository) updateOrderFulfillmentStatus(tx *gorm.DB, orderID uuid.UUID) error {
	// This would implement logic to check all line items and update order fulfillment status
	// Simplified implementation
	return nil
}

func (r *OrderRepository) updateOrderShippingStatus(tx *gorm.DB, orderID uuid.UUID) error {
	// This would implement logic to check if order should be marked as shipped
	// Simplified implementation
	return nil
}

func (r *OrderRepository) updateOrderPaymentStatus(tx *gorm.DB, orderID uuid.UUID) error {
	// This would implement logic to calculate payment status based on transactions
	// Simplified implementation
	return nil
}

func (r *OrderRepository) updateOrderReturnStatus(tx *gorm.DB, orderID uuid.UUID) error {
	// This would implement logic to check if order should be marked as returned
	// Simplified implementation
	return nil
}

func (r *OrderRepository) getEventTypeFromStatus(status models.OrderStatus) models.OrderEventType {
	switch status {
	case models.OrderStatusConfirmed:
		return models.OrderEventConfirmed
	case models.OrderStatusShipped:
		return models.OrderEventShipped
	case models.OrderStatusDelivered:
		return models.OrderEventDelivered
	case models.OrderStatusCancelled:
		return models.OrderEventCancelled
	default:
		return models.OrderEventCreated
	}
}
