package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"unified-commerce/services/cart/models"
	"unified-commerce/services/cart/repository"
	"unified-commerce/services/shared/logger"
)

// Service errors
var (
	ErrCartNotFound             = errors.New("cart not found")
	ErrLineItemNotFound         = errors.New("line item not found")
	ErrCheckoutNotFound         = errors.New("checkout not found")
	ErrInvalidQuantity          = errors.New("invalid quantity")
	ErrCartExpired              = errors.New("cart has expired")
	ErrCartCompleted            = errors.New("cart is already completed")
	ErrInsufficientInventory    = errors.New("insufficient inventory")
	ErrInvalidDiscountCode      = errors.New("invalid discount code")
	ErrDiscountAlreadyApplied   = errors.New("discount already applied")
	ErrInvalidCheckoutStep      = errors.New("invalid checkout step")
	ErrCheckoutAlreadyCompleted = errors.New("checkout already completed")
)

// CartService handles business logic for cart and checkout management
type CartService struct {
	repo   *repository.CartRepository
	logger *logger.Logger
}

// NewCartService creates a new cart service
func NewCartService(repo *repository.CartRepository, logger *logger.Logger) *CartService {
	return &CartService{
		repo:   repo,
		logger: logger,
	}
}

// Cart Management

// CreateCartRequest represents a request to create a cart
type CreateCartRequest struct {
	SessionID  string     `json:"session_id"`
	CustomerID *uuid.UUID `json:"customer_id"`
	MerchantID uuid.UUID  `json:"merchant_id" validate:"required"`
	Currency   string     `json:"currency"`
}

// CreateCart creates a new shopping cart
func (s *CartService) CreateCart(ctx context.Context, req *CreateCartRequest) (*models.Cart, error) {
	// Set default values
	if req.Currency == "" {
		req.Currency = "USD"
	}

	// Set expiration (30 days from now)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	cart := &models.Cart{
		SessionID:  req.SessionID,
		CustomerID: req.CustomerID,
		MerchantID: req.MerchantID,
		Status:     models.CartStatusActive,
		Currency:   req.Currency,
		ExpiresAt:  &expiresAt,
	}

	if err := s.repo.CreateCart(ctx, cart); err != nil {
		s.logger.WithError(err).Error("Failed to create cart")
		return nil, err
	}

	s.logger.WithField("cart_id", cart.ID).Info("Cart created successfully")
	return cart, nil
}

// GetCart retrieves a cart by ID
func (s *CartService) GetCart(ctx context.Context, id uuid.UUID) (*models.Cart, error) {
	cart, err := s.repo.GetCart(ctx, id)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, ErrCartNotFound
	}

	// Check if cart is expired
	if cart.ExpiresAt != nil && cart.ExpiresAt.Before(time.Now()) {
		return nil, ErrCartExpired
	}

	return cart, nil
}

// GetCartBySessionID retrieves a cart by session ID
func (s *CartService) GetCartBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	cart, err := s.repo.GetCartBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, ErrCartNotFound
	}

	// Check if cart is expired
	if cart.ExpiresAt != nil && cart.ExpiresAt.Before(time.Now()) {
		return nil, ErrCartExpired
	}

	return cart, nil
}

// GetCartByCustomerID retrieves active cart for a customer
func (s *CartService) GetCartByCustomerID(ctx context.Context, customerID uuid.UUID, merchantID uuid.UUID) (*models.Cart, error) {
	cart, err := s.repo.GetCartByCustomerID(ctx, customerID, merchantID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, ErrCartNotFound
	}

	// Check if cart is expired
	if cart.ExpiresAt != nil && cart.ExpiresAt.Before(time.Now()) {
		return nil, ErrCartExpired
	}

	return cart, nil
}

// UpdateCartRequest represents a request to update cart information
type UpdateCartRequest struct {
	CustomerEmail     string                 `json:"customer_email"`
	CustomerPhone     string                 `json:"customer_phone"`
	CustomerFirstName string                 `json:"customer_first_name"`
	CustomerLastName  string                 `json:"customer_last_name"`
	BillingAddress    *models.Address        `json:"billing_address"`
	ShippingAddress   *models.Address        `json:"shipping_address"`
	Notes             string                 `json:"notes"`
	Attributes        map[string]interface{} `json:"attributes"`
}

// UpdateCart updates cart information
func (s *CartService) UpdateCart(ctx context.Context, id uuid.UUID, req *UpdateCartRequest) (*models.Cart, error) {
	cart, err := s.GetCart(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.CustomerEmail != "" {
		cart.CustomerEmail = req.CustomerEmail
	}
	if req.CustomerPhone != "" {
		cart.CustomerPhone = req.CustomerPhone
	}
	if req.CustomerFirstName != "" {
		cart.CustomerFirstName = req.CustomerFirstName
	}
	if req.CustomerLastName != "" {
		cart.CustomerLastName = req.CustomerLastName
	}
	if req.BillingAddress != nil {
		cart.BillingAddress = *req.BillingAddress
	}
	if req.ShippingAddress != nil {
		cart.ShippingAddress = *req.ShippingAddress
	}
	if req.Notes != "" {
		cart.Notes = req.Notes
	}
	if req.Attributes != nil {
		cart.Attributes = req.Attributes
	}

	if err := s.repo.UpdateCart(ctx, cart); err != nil {
		s.logger.WithError(err).Error("Failed to update cart")
		return nil, err
	}

	s.logger.WithField("cart_id", id).Info("Cart updated successfully")
	return cart, nil
}

// Line Item Management

// AddLineItemRequest represents a request to add a line item to a cart
type AddLineItemRequest struct {
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
	ProductImage     string            `json:"product_image"`
	Taxable          bool              `json:"taxable"`
	RequiresShipping bool              `json:"requires_shipping"`
	Properties       map[string]string `json:"properties"`
}

// AddLineItem adds a line item to a cart
func (s *CartService) AddLineItem(ctx context.Context, cartID uuid.UUID, req *AddLineItemRequest) (*models.CartLineItem, error) {
	cart, err := s.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	// Check if cart is completed
	if cart.Status == models.CartStatusCompleted {
		return nil, ErrCartCompleted
	}

	lineItem := &models.CartLineItem{
		CartID:           cartID,
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		Name:             req.Name,
		SKU:              req.SKU,
		Quantity:         req.Quantity,
		Price:            req.Price,
		CompareAtPrice:   req.CompareAtPrice,
		ProductTitle:     req.ProductTitle,
		VariantTitle:     req.VariantTitle,
		Vendor:           req.Vendor,
		ProductImage:     req.ProductImage,
		Taxable:          req.Taxable,
		RequiresShipping: req.RequiresShipping,
		Properties:       req.Properties,
	}

	if err := s.repo.AddLineItem(ctx, lineItem); err != nil {
		s.logger.WithError(err).Error("Failed to add line item to cart")
		return nil, err
	}

	// Get updated cart
	updatedCart, err := s.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}

	s.logger.WithField("cart_id", cartID).WithField("line_item_id", lineItem.ID).Info("Line item added to cart successfully")
	return lineItem, nil
}

// UpdateLineItemRequest represents a request to update a line item
type UpdateLineItemRequest struct {
	Quantity   int               `json:"quantity" validate:"min=0"`
	Properties map[string]string `json:"properties"`
}

// UpdateLineItem updates a line item in a cart
func (s *CartService) UpdateLineItem(ctx context.Context, lineItemID uuid.UUID, req *UpdateLineItemRequest) (*models.CartLineItem, error) {
	lineItem, err := s.repo.GetLineItem(ctx, lineItemID)
	if err != nil {
		return nil, err
	}
	if lineItem == nil {
		return nil, ErrLineItemNotFound
	}

	// Get cart to check status
	cart, err := s.GetCart(ctx, lineItem.CartID)
	if err != nil {
		return nil, err
	}

	// Check if cart is completed
	if cart.Status == models.CartStatusCompleted {
		return nil, ErrCartCompleted
	}

	// Update fields
	if req.Quantity > 0 {
		lineItem.Quantity = req.Quantity
	}
	if req.Properties != nil {
		lineItem.Properties = req.Properties
	}

	if err := s.repo.UpdateLineItem(ctx, lineItem); err != nil {
		s.logger.WithError(err).Error("Failed to update line item")
		return nil, err
	}

	s.logger.WithField("line_item_id", lineItemID).Info("Line item updated successfully")
	return lineItem, nil
}

// RemoveLineItem removes a line item from a cart
func (s *CartService) RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	lineItem, err := s.repo.GetLineItem(ctx, lineItemID)
	if err != nil {
		return err
	}
	if lineItem == nil {
		return ErrLineItemNotFound
	}

	// Get cart to check status
	cart, err := s.GetCart(ctx, lineItem.CartID)
	if err != nil {
		return err
	}

	// Check if cart is completed
	if cart.Status == models.CartStatusCompleted {
		return ErrCartCompleted
	}

	if err := s.repo.RemoveLineItem(ctx, lineItemID); err != nil {
		s.logger.WithError(err).Error("Failed to remove line item")
		return err
	}

	s.logger.WithField("line_item_id", lineItemID).Info("Line item removed successfully")
	return nil
}

// Checkout Management

// CreateCheckoutRequest represents a request to create a checkout
type CreateCheckoutRequest struct {
	CartID uuid.UUID `json:"cart_id" validate:"required"`
	Email  string    `json:"email" validate:"required,email"`
	Phone  string    `json:"phone"`
}

// CreateCheckout creates a new checkout session
func (s *CartService) CreateCheckout(ctx context.Context, req *CreateCheckoutRequest) (*models.Checkout, error) {
	cart, err := s.GetCart(ctx, req.CartID)
	if err != nil {
		return nil, err
	}

	// Check if cart is completed
	if cart.Status == models.CartStatusCompleted {
		return nil, ErrCartCompleted
	}

	// Generate checkout token
	checkoutToken := s.generateCheckoutToken()

	checkout := &models.Checkout{
		CartID:           req.CartID,
		CheckoutToken:    checkoutToken,
		Status:           models.CheckoutStatusPending,
		Email:            req.Email,
		Phone:            req.Phone,
		RequiresShipping: s.requiresShipping(cart),
	}

	if err := s.repo.CreateCheckout(ctx, checkout); err != nil {
		s.logger.WithError(err).Error("Failed to create checkout")
		return nil, err
	}

	// Create checkout event
	event := &models.CheckoutEvent{
		CheckoutID:  checkout.ID,
		EventType:   models.CheckoutEventCreated,
		Description: "Checkout created",
	}
	if err := s.repo.CreateCheckoutEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create checkout event")
	}

	s.logger.WithField("checkout_id", checkout.ID).WithField("checkout_token", checkoutToken).Info("Checkout created successfully")
	return checkout, nil
}

// GetCheckout retrieves a checkout by ID
func (s *CartService) GetCheckout(ctx context.Context, id uuid.UUID) (*models.Checkout, error) {
	checkout, err := s.repo.GetCheckout(ctx, id)
	if err != nil {
		return nil, err
	}
	if checkout == nil {
		return nil, ErrCheckoutNotFound
	}

	return checkout, nil
}

// GetCheckoutByToken retrieves a checkout by token
func (s *CartService) GetCheckoutByToken(ctx context.Context, token string) (*models.Checkout, error) {
	checkout, err := s.repo.GetCheckoutByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if checkout == nil {
		return nil, ErrCheckoutNotFound
	}

	return checkout, nil
}

// UpdateCheckoutCustomerInfo updates customer information in checkout
func (s *CartService) UpdateCheckoutCustomerInfo(ctx context.Context, checkoutID uuid.UUID, req *CreateCheckoutRequest) (*models.Checkout, error) {
	checkout, err := s.GetCheckout(ctx, checkoutID)
	if err != nil {
		return nil, err
	}

	// Check if checkout is completed
	if checkout.Status == models.CheckoutStatusCompleted {
		return nil, ErrCheckoutAlreadyCompleted
	}

	// Update fields
	checkout.Email = req.Email
	checkout.Phone = req.Phone

	if err := s.repo.UpdateCheckout(ctx, checkout); err != nil {
		s.logger.WithError(err).Error("Failed to update checkout customer info")
		return nil, err
	}

	// Create event
	event := &models.CheckoutEvent{
		CheckoutID:  checkout.ID,
		EventType:   models.CheckoutEventCustomerInfoAdded,
		Description: "Customer information added",
	}
	if err := s.repo.CreateCheckoutEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create checkout event")
	}

	s.logger.WithField("checkout_id", checkoutID).Info("Checkout customer info updated successfully")
	return checkout, nil
}

// UpdateCheckoutShippingAddress updates shipping address in checkout
func (s *CartService) UpdateCheckoutShippingAddress(ctx context.Context, checkoutID uuid.UUID, address models.Address) (*models.Checkout, error) {
	checkout, err := s.GetCheckout(ctx, checkoutID)
	if err != nil {
		return nil, err
	}

	// Check if checkout is completed
	if checkout.Status == models.CheckoutStatusCompleted {
		return nil, ErrCheckoutAlreadyCompleted
	}

	// Update cart shipping address
	cart, err := s.GetCart(ctx, checkout.CartID)
	if err != nil {
		return nil, err
	}

	cart.ShippingAddress = address
	if err := s.repo.UpdateCart(ctx, cart); err != nil {
		s.logger.WithError(err).Error("Failed to update cart shipping address")
		return nil, err
	}

	// Update checkout step
	checkout.CompletedStep = models.CheckoutStepShipping

	if err := s.repo.UpdateCheckout(ctx, checkout); err != nil {
		s.logger.WithError(err).Error("Failed to update checkout")
		return nil, err
	}

	// Create event
	event := &models.CheckoutEvent{
		CheckoutID:  checkout.ID,
		EventType:   models.CheckoutEventShippingAdded,
		Description: "Shipping address added",
	}
	if err := s.repo.CreateCheckoutEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create checkout event")
	}

	s.logger.WithField("checkout_id", checkoutID).Info("Checkout shipping address updated successfully")
	return checkout, nil
}

// AddShippingLine adds shipping information to checkout
func (s *CartService) AddShippingLine(ctx context.Context, checkoutID uuid.UUID, shippingLine *models.CartShippingLine) error {
	checkout, err := s.GetCheckout(ctx, checkoutID)
	if err != nil {
		return err
	}

	// Check if checkout is completed
	if checkout.Status == models.CheckoutStatusCompleted {
		return ErrCheckoutAlreadyCompleted
	}

	if err := s.repo.AddShippingLine(ctx, shippingLine); err != nil {
		s.logger.WithError(err).Error("Failed to add shipping line")
		return err
	}

	s.logger.WithField("checkout_id", checkoutID).Info("Shipping line added successfully")
	return nil
}

// ApplyDiscountCode applies a discount code to checkout
func (s *CartService) ApplyDiscountCode(ctx context.Context, checkoutID uuid.UUID, discountCode string) error {
	checkout, err := s.GetCheckout(ctx, checkoutID)
	if err != nil {
		return err
	}

	// Check if checkout is completed
	if checkout.Status == models.CheckoutStatusCompleted {
		return ErrCheckoutAlreadyCompleted
	}

	// Validate discount code (simplified - would integrate with promotions service in real implementation)
	if !s.isValidDiscountCode(discountCode) {
		return ErrInvalidDiscountCode
	}

	// Check if discount already applied
	cart, err := s.GetCart(ctx, checkout.CartID)
	if err != nil {
		return err
	}

	for _, discount := range cart.DiscountApplications {
		if discount.Code == discountCode {
			return ErrDiscountAlreadyApplied
		}
	}

	// Apply discount (simplified - would calculate actual discount in real implementation)
	discount := &models.CartDiscountApplication{
		CartID:           checkout.CartID,
		Type:             models.DiscountTypeCode,
		Code:             discountCode,
		Title:            fmt.Sprintf("Discount: %s", discountCode),
		Value:            10.00, // Simplified
		ValueType:        models.DiscountValueTypeFixed,
		AllocationMethod: models.AllocationMethodAcross,
		TargetSelection:  models.TargetSelectionAll,
		TargetType:       models.TargetTypeLineItem,
	}

	if err := s.repo.ApplyDiscount(ctx, discount); err != nil {
		s.logger.WithError(err).Error("Failed to apply discount")
		return err
	}

	// Create event
	event := &models.CheckoutEvent{
		CheckoutID:  checkout.ID,
		EventType:   models.CheckoutEventDiscountApplied,
		Description: fmt.Sprintf("Discount code applied: %s", discountCode),
	}
	if err := s.repo.CreateCheckoutEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create checkout event")
	}

	s.logger.WithField("checkout_id", checkoutID).WithField("discount_code", discountCode).Info("Discount code applied successfully")
	return nil
}

// RemoveDiscountCode removes a discount code from checkout
func (s *CartService) RemoveDiscountCode(ctx context.Context, checkoutID uuid.UUID, discountCode string) error {
	checkout, err := s.GetCheckout(ctx, checkoutID)
	if err != nil {
		return err
	}

	// Check if checkout is completed
	if checkout.Status == models.CheckoutStatusCompleted {
		return ErrCheckoutAlreadyCompleted
	}

	if err := s.repo.RemoveDiscount(ctx, checkout.CartID, discountCode); err != nil {
		s.logger.WithError(err).Error("Failed to remove discount")
		return err
	}

	// Create event
	event := &models.CheckoutEvent{
		CheckoutID:  checkout.ID,
		EventType:   models.CheckoutEventDiscountRemoved,
		Description: fmt.Sprintf("Discount code removed: %s", discountCode),
	}
	if err := s.repo.CreateCheckoutEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create checkout event")
	}

	s.logger.WithField("checkout_id", checkoutID).WithField("discount_code", discountCode).Info("Discount code removed successfully")
	return nil
}

// CompleteCheckout completes the checkout process and creates an order
func (s *CartService) CompleteCheckout(ctx context.Context, checkoutID uuid.UUID, paymentMethodID string) (string, error) {
	checkout, err := s.GetCheckout(ctx, checkoutID)
	if err != nil {
		return "", err
	}

	// Check if checkout is completed
	if checkout.Status == models.CheckoutStatusCompleted {
		return "", ErrCheckoutAlreadyCompleted
	}

	// Update checkout payment information
	checkout.PaymentMethodID = paymentMethodID
	checkout.CompletedStep = models.CheckoutStepComplete

	if err := s.repo.UpdateCheckout(ctx, checkout); err != nil {
		s.logger.WithError(err).Error("Failed to update checkout")
		return "", err
	}

	// In a real implementation, this would integrate with the Order Service to create an order
	// For now, we'll generate a mock order ID
	orderID := s.generateOrderID()

	// Mark checkout as completed
	if err := s.repo.CompleteCheckout(ctx, checkoutID, orderID); err != nil {
		s.logger.WithError(err).Error("Failed to complete checkout")
		return "", err
	}

	s.logger.WithField("checkout_id", checkoutID).WithField("order_id", orderID).Info("Checkout completed successfully")
	return orderID, nil
}

// Utility Methods

// generateCheckoutToken generates a unique checkout token
func (s *CartService) generateCheckoutToken() string {
	// Simple implementation - in production you'd want a more sophisticated approach
	timestamp := time.Now().Unix()
	random := rand.Intn(1000000)
	return fmt.Sprintf("chk_%d_%06d", timestamp, random)
}

// generateOrderID generates a mock order ID
func (s *CartService) generateOrderID() string {
	timestamp := time.Now().Format("20060102")
	random := rand.Intn(10000)
	return fmt.Sprintf("ORD-%s-%04d", timestamp, random)
}

// requiresShipping determines if a cart requires shipping
func (s *CartService) requiresShipping(cart *models.Cart) bool {
	for _, item := range cart.LineItems {
		if item.RequiresShipping {
			return true
		}
	}
	return false
}

// isValidDiscountCode validates a discount code (simplified)
func (s *CartService) isValidDiscountCode(code string) bool {
	// In a real implementation, this would check against the Promotions Service
	// For now, we'll accept any non-empty code
	return code != ""
}

// GetShippingRates retrieves available shipping rates
func (s *CartService) GetShippingRates(ctx context.Context) ([]*models.ShippingRate, error) {
	return s.repo.GetShippingRates(ctx)
}

// GetPaymentMethods retrieves available payment methods
func (s *CartService) GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	return s.repo.GetPaymentMethods(ctx)
}

// Abandonment Management

// MarkAbandonedCarts marks carts as abandoned
func (s *CartService) MarkAbandonedCarts(ctx context.Context, abandonmentThreshold time.Duration) error {
	if err := s.repo.MarkCartAbandoned(ctx, abandonmentThreshold); err != nil {
		s.logger.WithError(err).Error("Failed to mark abandoned carts")
		return err
	}

	s.logger.Info("Abandoned carts marked successfully")
	return nil
}

// GetAbandonedCarts retrieves abandoned carts for recovery
func (s *CartService) GetAbandonedCarts(ctx context.Context, page, limit int) ([]*models.Cart, error) {
	offset := (page - 1) * limit
	return s.repo.GetAbandonedCarts(ctx, limit, offset)
}

// Cleanup Operations

// CleanupExpiredCarts removes expired carts
func (s *CartService) CleanupExpiredCarts(ctx context.Context) error {
	if err := s.repo.CleanupExpiredCarts(ctx); err != nil {
		s.logger.WithError(err).Error("Failed to cleanup expired carts")
		return err
	}

	s.logger.Info("Expired carts cleaned up successfully")
	return nil
}
