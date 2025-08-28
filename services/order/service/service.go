package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"unified-commerce/services/order/models"
	"unified-commerce/services/order/repository"
	"unified-commerce/services/shared/logger"
)

// Service errors
var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrLineItemNotFound      = errors.New("line item not found")
	ErrFulfillmentNotFound   = errors.New("fulfillment not found")
	ErrTransactionNotFound   = errors.New("transaction not found")
	ErrReturnNotFound        = errors.New("return not found")
	ErrInvalidStatus         = errors.New("invalid status transition")
	ErrInsufficientPayment   = errors.New("insufficient payment")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrOrderNotCancellable   = errors.New("order cannot be cancelled")
	ErrInvalidQuantity       = errors.New("invalid quantity")
)

// OrderService handles business logic for order management
type OrderService struct {
	repo   *repository.OrderRepository
	logger *logger.Logger
}

// NewOrderService creates a new order service
func NewOrderService(repo *repository.OrderRepository, logger *logger.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger,
	}
}

// Order Service Methods

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	MerchantID      uuid.UUID               `json:"merchant_id" validate:"required"`
	CustomerID      *uuid.UUID              `json:"customer_id"`
	LocationID      *uuid.UUID              `json:"location_id"`
	Customer        models.CustomerInfo     `json:"customer" validate:"required"`
	BillingAddress  models.Address          `json:"billing_address"`
	ShippingAddress models.Address          `json:"shipping_address"`
	LineItems       []CreateLineItemRequest `json:"line_items" validate:"required,min=1"`
	ShippingMethod  string                  `json:"shipping_method"`
	ShippingRate    float64                 `json:"shipping_rate"`
	Currency        string                  `json:"currency"`
	Source          models.OrderSource      `json:"source"`
	Channel         string                  `json:"channel"`
	Tags            []string                `json:"tags"`
	Notes           string                  `json:"notes"`
}

// CreateLineItemRequest represents a line item in an order creation request
type CreateLineItemRequest struct {
	ProductID        uuid.UUID         `json:"product_id" validate:"required"`
	ProductVariantID *uuid.UUID        `json:"product_variant_id"`
	SKU              string            `json:"sku" validate:"required"`
	Name             string            `json:"name" validate:"required"`
	Quantity         int               `json:"quantity" validate:"required,min=1"`
	Price            float64           `json:"price" validate:"required,min=0"`
	CompareAtPrice   float64           `json:"compare_at_price"`
	ProductTitle     string            `json:"product_title"`
	VariantTitle     string            `json:"variant_title"`
	Vendor           string            `json:"vendor"`
	Taxable          bool              `json:"taxable"`
	RequiresShipping bool              `json:"requires_shipping"`
	Properties       map[string]string `json:"properties"`
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*models.Order, error) {
	// Generate order number
	orderNumber := s.generateOrderNumber()

	// Create order
	order := &models.Order{
		OrderNumber:       orderNumber,
		MerchantID:        req.MerchantID,
		CustomerID:        req.CustomerID,
		LocationID:        req.LocationID,
		Status:            models.OrderStatusPending,
		FulfillmentStatus: models.FulfillmentStatusUnfulfilled,
		PaymentStatus:     models.PaymentStatusPending,
		Customer:          req.Customer,
		BillingAddress:    req.BillingAddress,
		ShippingAddress:   req.ShippingAddress,
		ShippingMethod:    req.ShippingMethod,
		ShippingRate:      req.ShippingRate,
		Source:            req.Source,
		Channel:           req.Channel,
		Currency:          req.Currency,
		Tags:              req.Tags,
		Notes:             req.Notes,
	}

	if order.Currency == "" {
		order.Currency = "USD"
	}
	if order.Source == "" {
		order.Source = models.OrderSourceOnline
	}

	// Create line items
	var subtotal float64
	for _, item := range req.LineItems {
		linePrice := float64(item.Quantity) * item.Price
		subtotal += linePrice

		lineItem := models.OrderLineItem{
			OrderID:           order.ID,
			ProductID:         item.ProductID,
			ProductVariantID:  item.ProductVariantID,
			Name:              item.Name,
			SKU:               item.SKU,
			ProductTitle:      item.ProductTitle,
			VariantTitle:      item.VariantTitle,
			Vendor:            item.Vendor,
			Quantity:          item.Quantity,
			Price:             item.Price,
			CompareAtPrice:    item.CompareAtPrice,
			LinePrice:         linePrice,
			Taxable:           item.Taxable,
			RequiresShipping:  item.RequiresShipping,
			Properties:        item.Properties,
			FulfillmentStatus: models.FulfillmentStatusUnfulfilled,
		}
		order.LineItems = append(order.LineItems, lineItem)
	}

	// Calculate totals (simplified - real implementation would include tax calculation)
	order.SubtotalPrice = subtotal
	order.TotalShipping = req.ShippingRate
	order.TotalPrice = subtotal + req.ShippingRate

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.logger.WithError(err).Error("Failed to create order")
		return nil, err
	}

	s.logger.WithField("order_id", order.ID).WithField("order_number", order.OrderNumber).Info("Order created successfully")
	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

// GetOrderByNumber retrieves an order by order number
func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	order, err := s.repo.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

// GetOrdersByMerchant retrieves orders for a merchant
func (s *OrderService) GetOrdersByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, page, limit int) ([]*models.Order, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetOrdersByMerchant(ctx, merchantID, filters, limit, offset)
}

// GetOrdersByCustomer retrieves orders for a customer
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID uuid.UUID, page, limit int) ([]*models.Order, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetOrdersByCustomer(ctx, customerID, limit, offset)
}

// UpdateOrderRequest represents a request to update an order
type UpdateOrderRequest struct {
	Tags          []string `json:"tags"`
	Notes         string   `json:"notes"`
	InternalNotes string   `json:"internal_notes"`
}

// UpdateOrder updates an order
func (s *OrderService) UpdateOrder(ctx context.Context, id uuid.UUID, req *UpdateOrderRequest) (*models.Order, error) {
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// Update fields
	order.Tags = req.Tags
	order.Notes = req.Notes
	order.InternalNotes = req.InternalNotes

	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		s.logger.WithError(err).Error("Failed to update order")
		return nil, err
	}

	s.logger.WithField("order_id", order.ID).Info("Order updated successfully")
	return order, nil
}

// ConfirmOrder confirms a pending order
func (s *OrderService) ConfirmOrder(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error {
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return ErrOrderNotFound
	}

	if order.Status != models.OrderStatusPending {
		return ErrInvalidStatus
	}

	now := time.Now()
	if err := s.repo.UpdateOrderStatus(ctx, id, models.OrderStatusConfirmed, "Order confirmed", userID); err != nil {
		s.logger.WithError(err).Error("Failed to confirm order")
		return err
	}

	// Update processed timestamp
	order.ProcessedAt = &now
	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		s.logger.WithError(err).Error("Failed to update order processed time")
		return err
	}

	s.logger.WithField("order_id", id).Info("Order confirmed successfully")
	return nil
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(ctx context.Context, id uuid.UUID, reason string, userID *uuid.UUID) error {
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return ErrOrderNotFound
	}

	// Check if order can be cancelled
	if order.Status == models.OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if order.Status == models.OrderStatusDelivered || order.Status == models.OrderStatusReturned {
		return ErrOrderNotCancellable
	}

	if err := s.repo.CancelOrder(ctx, id, reason, userID); err != nil {
		s.logger.WithError(err).Error("Failed to cancel order")
		return err
	}

	s.logger.WithField("order_id", id).WithField("reason", reason).Info("Order cancelled successfully")
	return nil
}

// Line Item Management

// AddLineItemRequest represents a request to add a line item
type AddLineItemRequest struct {
	ProductID        uuid.UUID         `json:"product_id" validate:"required"`
	ProductVariantID *uuid.UUID        `json:"product_variant_id"`
	SKU              string            `json:"sku" validate:"required"`
	Name             string            `json:"name" validate:"required"`
	Quantity         int               `json:"quantity" validate:"required,min=1"`
	Price            float64           `json:"price" validate:"required,min=0"`
	Properties       map[string]string `json:"properties"`
}

// AddLineItem adds a line item to an order
func (s *OrderService) AddLineItem(ctx context.Context, orderID uuid.UUID, req *AddLineItemRequest) (*models.OrderLineItem, error) {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	linePrice := float64(req.Quantity) * req.Price
	lineItem := &models.OrderLineItem{
		OrderID:           orderID,
		ProductID:         req.ProductID,
		ProductVariantID:  req.ProductVariantID,
		Name:              req.Name,
		SKU:               req.SKU,
		Quantity:          req.Quantity,
		Price:             req.Price,
		LinePrice:         linePrice,
		Properties:        req.Properties,
		FulfillmentStatus: models.FulfillmentStatusUnfulfilled,
	}

	if err := s.repo.AddLineItem(ctx, lineItem); err != nil {
		s.logger.WithError(err).Error("Failed to add line item")
		return nil, err
	}

	s.logger.WithField("order_id", orderID).WithField("line_item_id", lineItem.ID).Info("Line item added successfully")
	return lineItem, nil
}

// UpdateLineItemRequest represents a request to update a line item
type UpdateLineItemRequest struct {
	Quantity int     `json:"quantity" validate:"min=1"`
	Price    float64 `json:"price" validate:"min=0"`
}

// UpdateLineItem updates a line item
func (s *OrderService) UpdateLineItem(ctx context.Context, lineItemID uuid.UUID, req *UpdateLineItemRequest) (*models.OrderLineItem, error) {
	// Get current line item
	order, err := s.repo.GetOrder(ctx, uuid.Nil) // This would need proper implementation to get line item
	if err != nil {
		return nil, err
	}

	var lineItem *models.OrderLineItem
	for i := range order.LineItems {
		if order.LineItems[i].ID == lineItemID {
			lineItem = &order.LineItems[i]
			break
		}
	}

	if lineItem == nil {
		return nil, ErrLineItemNotFound
	}

	// Update fields
	if req.Quantity > 0 {
		lineItem.Quantity = req.Quantity
	}
	if req.Price >= 0 {
		lineItem.Price = req.Price
	}
	lineItem.LinePrice = float64(lineItem.Quantity) * lineItem.Price

	if err := s.repo.UpdateLineItem(ctx, lineItem); err != nil {
		s.logger.WithError(err).Error("Failed to update line item")
		return nil, err
	}

	s.logger.WithField("line_item_id", lineItemID).Info("Line item updated successfully")
	return lineItem, nil
}

// RemoveLineItem removes a line item from an order
func (s *OrderService) RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	if err := s.repo.RemoveLineItem(ctx, lineItemID); err != nil {
		s.logger.WithError(err).Error("Failed to remove line item")
		return err
	}

	s.logger.WithField("line_item_id", lineItemID).Info("Line item removed successfully")
	return nil
}

// Fulfillment Management

// CreateFulfillmentRequest represents a request to create a fulfillment
type CreateFulfillmentRequest struct {
	OrderID         uuid.UUID                   `json:"order_id" validate:"required"`
	LocationID      *uuid.UUID                  `json:"location_id"`
	LineItems       []CreateFulfillmentLineItem `json:"line_items" validate:"required,min=1"`
	TrackingNumber  string                      `json:"tracking_number"`
	TrackingCompany string                      `json:"tracking_company"`
	Service         string                      `json:"service"`
	NotifyCustomer  bool                        `json:"notify_customer"`
}

// CreateFulfillmentLineItem represents a line item in a fulfillment
type CreateFulfillmentLineItem struct {
	LineItemID uuid.UUID `json:"line_item_id" validate:"required"`
	Quantity   int       `json:"quantity" validate:"required,min=1"`
}

// CreateFulfillment creates a new fulfillment
func (s *OrderService) CreateFulfillment(ctx context.Context, req *CreateFulfillmentRequest) (*models.Fulfillment, error) {
	order, err := s.repo.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	fulfillment := &models.Fulfillment{
		OrderID:         req.OrderID,
		LocationID:      req.LocationID,
		Status:          models.FulfillmentStatusUnfulfilled,
		TrackingNumber:  req.TrackingNumber,
		TrackingCompany: req.TrackingCompany,
		Service:         req.Service,
		ShipmentStatus:  models.ShipmentStatusPending,
	}

	// Add line items
	for _, item := range req.LineItems {
		fulfillmentLineItem := models.FulfillmentLineItem{
			FulfillmentID: fulfillment.ID,
			LineItemID:    item.LineItemID,
			Quantity:      item.Quantity,
		}
		fulfillment.LineItems = append(fulfillment.LineItems, fulfillmentLineItem)
	}

	if err := s.repo.CreateFulfillment(ctx, fulfillment); err != nil {
		s.logger.WithError(err).Error("Failed to create fulfillment")
		return nil, err
	}

	s.logger.WithField("fulfillment_id", fulfillment.ID).WithField("order_id", req.OrderID).Info("Fulfillment created successfully")
	return fulfillment, nil
}

// MarkFulfillmentShipped marks a fulfillment as shipped
func (s *OrderService) MarkFulfillmentShipped(ctx context.Context, fulfillmentID uuid.UUID, trackingNumber, trackingURL string) error {
	if err := s.repo.MarkFulfillmentShipped(ctx, fulfillmentID, trackingNumber, trackingURL); err != nil {
		s.logger.WithError(err).Error("Failed to mark fulfillment as shipped")
		return err
	}

	s.logger.WithField("fulfillment_id", fulfillmentID).Info("Fulfillment marked as shipped")
	return nil
}

// Payment Management

// CreateTransactionRequest represents a request to create a transaction
type CreateTransactionRequest struct {
	OrderID              uuid.UUID              `json:"order_id" validate:"required"`
	Kind                 models.TransactionKind `json:"kind" validate:"required"`
	Gateway              string                 `json:"gateway" validate:"required"`
	Amount               float64                `json:"amount" validate:"required,min=0"`
	Currency             string                 `json:"currency"`
	GatewayTransactionID string                 `json:"gateway_transaction_id"`
	PaymentMethodID      string                 `json:"payment_method_id"`
	AuthorizationCode    string                 `json:"authorization_code"`
}

// CreateTransaction creates a new payment transaction
func (s *OrderService) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*models.Transaction, error) {
	order, err := s.repo.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	transaction := &models.Transaction{
		OrderID:              req.OrderID,
		Kind:                 req.Kind,
		Gateway:              req.Gateway,
		Status:               models.TransactionStatusPending,
		Amount:               req.Amount,
		Currency:             req.Currency,
		GatewayTransactionID: req.GatewayTransactionID,
		PaymentMethodID:      req.PaymentMethodID,
		AuthorizationCode:    req.AuthorizationCode,
	}

	if transaction.Currency == "" {
		transaction.Currency = order.Currency
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		s.logger.WithError(err).Error("Failed to create transaction")
		return nil, err
	}

	s.logger.WithField("transaction_id", transaction.ID).WithField("order_id", req.OrderID).Info("Transaction created successfully")
	return transaction, nil
}

// Return Management

// CreateReturnRequest represents a request to create a return
type CreateReturnRequest struct {
	OrderID        uuid.UUID              `json:"order_id" validate:"required"`
	Reason         models.ReturnReason    `json:"reason" validate:"required"`
	CustomerNotes  string                 `json:"customer_notes"`
	LineItems      []CreateReturnLineItem `json:"line_items" validate:"required,min=1"`
	RefundShipping bool                   `json:"refund_shipping"`
	RestockItems   bool                   `json:"restock_items"`
	NotifyCustomer bool                   `json:"notify_customer"`
}

// CreateReturnLineItem represents a line item in a return
type CreateReturnLineItem struct {
	LineItemID uuid.UUID           `json:"line_item_id" validate:"required"`
	Quantity   int                 `json:"quantity" validate:"required,min=1"`
	Reason     models.ReturnReason `json:"reason"`
	Notes      string              `json:"notes"`
}

// CreateReturn creates a new return
func (s *OrderService) CreateReturn(ctx context.Context, req *CreateReturnRequest) (*models.Return, error) {
	order, err := s.repo.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// Generate return number
	returnNumber := s.generateReturnNumber()

	returnItem := &models.Return{
		OrderID:        req.OrderID,
		ReturnNumber:   returnNumber,
		Status:         models.ReturnStatusPending,
		Reason:         req.Reason,
		CustomerNotes:  req.CustomerNotes,
		RefundShipping: req.RefundShipping,
		RestockItems:   req.RestockItems,
		NotifyCustomer: req.NotifyCustomer,
	}

	// Add line items
	for _, item := range req.LineItems {
		returnLineItem := models.ReturnLineItem{
			ReturnID:   returnItem.ID,
			LineItemID: item.LineItemID,
			Quantity:   item.Quantity,
			Reason:     item.Reason,
			Notes:      item.Notes,
		}
		returnItem.LineItems = append(returnItem.LineItems, returnLineItem)
	}

	if err := s.repo.CreateReturn(ctx, returnItem); err != nil {
		s.logger.WithError(err).Error("Failed to create return")
		return nil, err
	}

	s.logger.WithField("return_id", returnItem.ID).WithField("return_number", returnNumber).Info("Return created successfully")
	return returnItem, nil
}

// ProcessReturn processes a return and issues refund
func (s *OrderService) ProcessReturn(ctx context.Context, returnID uuid.UUID, refundAmount float64, restockItems bool) error {
	if err := s.repo.ProcessReturn(ctx, returnID, refundAmount, restockItems); err != nil {
		s.logger.WithError(err).Error("Failed to process return")
		return err
	}

	s.logger.WithField("return_id", returnID).WithField("refund_amount", refundAmount).Info("Return processed successfully")
	return nil
}

// Analytics and Reporting

// GetOrderStats retrieves order statistics for a merchant
func (s *OrderService) GetOrderStats(ctx context.Context, merchantID uuid.UUID, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	return s.repo.GetOrderStats(ctx, merchantID, dateFrom, dateTo)
}

// GetOrderEvents retrieves events for an order
func (s *OrderService) GetOrderEvents(ctx context.Context, orderID uuid.UUID) ([]*models.OrderEvent, error) {
	return s.repo.GetOrderEvents(ctx, orderID)
}

// Utility Methods

// generateOrderNumber generates a unique order number
func (s *OrderService) generateOrderNumber() string {
	// Simple implementation - in production you'd want a more sophisticated approach
	timestamp := time.Now().Format("20060102")
	random := rand.Intn(10000)
	return fmt.Sprintf("ORD-%s-%04d", timestamp, random)
}

// generateReturnNumber generates a unique return number
func (s *OrderService) generateReturnNumber() string {
	timestamp := time.Now().Format("20060102")
	random := rand.Intn(10000)
	return fmt.Sprintf("RET-%s-%04d", timestamp, random)
}
