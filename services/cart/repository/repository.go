package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"unified-commerce/services/cart/models"
	"unified-commerce/services/shared/logger"
)

// CartRepository handles database operations for cart and checkout management
type CartRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *gorm.DB, logger *logger.Logger) *CartRepository {
	return &CartRepository{
		db:     db,
		logger: logger,
	}
}

// Cart Operations

// CreateCart creates a new shopping cart
func (r *CartRepository) CreateCart(ctx context.Context, cart *models.Cart) error {
	if err := r.db.WithContext(ctx).Create(cart).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create cart")
		return err
	}
	return nil
}

// GetCart retrieves a cart by ID with all related data
func (r *CartRepository) GetCart(ctx context.Context, id uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Preload("LineItems.DiscountAllocations").
		Preload("TaxLines").
		Preload("ShippingLines").
		Preload("DiscountApplications").
		First(&cart, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get cart")
		return nil, err
	}
	return &cart, nil
}

// GetCartBySessionID retrieves a cart by session ID
func (r *CartRepository) GetCartBySessionID(ctx context.Context, sessionID string) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Preload("TaxLines").
		Preload("ShippingLines").
		Preload("DiscountApplications").
		First(&cart, "session_id = ? AND status = ?", sessionID, models.CartStatusActive).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get cart by session ID")
		return nil, err
	}
	return &cart, nil
}

// GetCartByCustomerID retrieves active cart for a customer
func (r *CartRepository) GetCartByCustomerID(ctx context.Context, customerID uuid.UUID, merchantID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Preload("TaxLines").
		Preload("ShippingLines").
		Preload("DiscountApplications").
		First(&cart, "customer_id = ? AND merchant_id = ? AND status = ?", customerID, merchantID, models.CartStatusActive).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get cart by customer ID")
		return nil, err
	}
	return &cart, nil
}

// UpdateCart updates a cart
func (r *CartRepository) UpdateCart(ctx context.Context, cart *models.Cart) error {
	if err := r.db.WithContext(ctx).Save(cart).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update cart")
		return err
	}
	return nil
}

// UpdateCartTotals recalculates and updates cart totals
func (r *CartRepository) UpdateCartTotals(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Calculate subtotal from line items
		var subtotal float64
		if err := tx.Model(&models.CartLineItem{}).
			Where("cart_id = ?", cartID).
			Select("COALESCE(SUM(line_price), 0)").
			Scan(&subtotal).Error; err != nil {
			return err
		}

		// Calculate total tax
		var totalTax float64
		if err := tx.Model(&models.CartTaxLine{}).
			Where("cart_id = ?", cartID).
			Select("COALESCE(SUM(price), 0)").
			Scan(&totalTax).Error; err != nil {
			return err
		}

		// Calculate total shipping
		var totalShipping float64
		if err := tx.Model(&models.CartShippingLine{}).
			Where("cart_id = ?", cartID).
			Select("COALESCE(SUM(discounted_price), 0)").
			Scan(&totalShipping).Error; err != nil {
			return err
		}

		// Calculate total discount
		var totalDiscount float64
		if err := tx.Model(&models.CartLineItem{}).
			Where("cart_id = ?", cartID).
			Select("COALESCE(SUM(total_discount), 0)").
			Scan(&totalDiscount).Error; err != nil {
			return err
		}

		// Update cart totals
		totalPrice := subtotal + totalTax + totalShipping - totalDiscount
		if err := tx.Model(&models.Cart{}).
			Where("id = ?", cartID).
			Updates(map[string]interface{}{
				"subtotal_price": subtotal,
				"total_tax":      totalTax,
				"total_shipping": totalShipping,
				"total_discount": totalDiscount,
				"total_price":    totalPrice,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteCart soft deletes a cart
func (r *CartRepository) DeleteCart(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Cart{}, "id = ?", id).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete cart")
		return err
	}
	return nil
}

// Line Item Operations

// AddLineItem adds a line item to a cart
func (r *CartRepository) AddLineItem(ctx context.Context, lineItem *models.CartLineItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if line item with same product/variant already exists
		var existingItem models.CartLineItem
		err := tx.Where("cart_id = ? AND product_id = ? AND product_variant_id = ?",
			lineItem.CartID, lineItem.ProductID, lineItem.ProductVariantID).
			First(&existingItem).Error

		if err == nil {
			// Update existing line item quantity
			existingItem.Quantity += lineItem.Quantity
			existingItem.LinePrice = float64(existingItem.Quantity) * existingItem.Price
			if err := tx.Save(&existingItem).Error; err != nil {
				return err
			}
		} else if err == gorm.ErrRecordNotFound {
			// Create new line item
			lineItem.LinePrice = float64(lineItem.Quantity) * lineItem.Price
			if err := tx.Create(lineItem).Error; err != nil {
				return err
			}
		} else {
			return err
		}

		// Update cart totals
		return r.recalculateCartTotals(tx, lineItem.CartID)
	})
}

// UpdateLineItem updates a cart line item
func (r *CartRepository) UpdateLineItem(ctx context.Context, lineItem *models.CartLineItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update line price
		lineItem.LinePrice = float64(lineItem.Quantity) * lineItem.Price

		// Save line item
		if err := tx.Save(lineItem).Error; err != nil {
			return err
		}

		// Update cart totals
		return r.recalculateCartTotals(tx, lineItem.CartID)
	})
}

// RemoveLineItem removes a line item from a cart
func (r *CartRepository) RemoveLineItem(ctx context.Context, lineItemID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get line item to get cart ID
		var lineItem models.CartLineItem
		if err := tx.First(&lineItem, "id = ?", lineItemID).Error; err != nil {
			return err
		}

		// Delete line item
		if err := tx.Delete(&lineItem).Error; err != nil {
			return err
		}

		// Update cart totals
		return r.recalculateCartTotals(tx, lineItem.CartID)
	})
}

// GetLineItem retrieves a line item by ID
func (r *CartRepository) GetLineItem(ctx context.Context, id uuid.UUID) (*models.CartLineItem, error) {
	var lineItem models.CartLineItem
	if err := r.db.WithContext(ctx).First(&lineItem, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get line item")
		return nil, err
	}
	return &lineItem, nil
}

// Checkout Operations

// CreateCheckout creates a new checkout session
func (r *CartRepository) CreateCheckout(ctx context.Context, checkout *models.Checkout) error {
	if err := r.db.WithContext(ctx).Create(checkout).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create checkout")
		return err
	}
	return nil
}

// GetCheckout retrieves a checkout by ID
func (r *CartRepository) GetCheckout(ctx context.Context, id uuid.UUID) (*models.Checkout, error) {
	var checkout models.Checkout
	if err := r.db.WithContext(ctx).
		Preload("Cart").
		Preload("Cart.LineItems").
		Preload("Events").
		First(&checkout, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get checkout")
		return nil, err
	}
	return &checkout, nil
}

// GetCheckoutByToken retrieves a checkout by token
func (r *CartRepository) GetCheckoutByToken(ctx context.Context, token string) (*models.Checkout, error) {
	var checkout models.Checkout
	if err := r.db.WithContext(ctx).
		Preload("Cart").
		Preload("Cart.LineItems").
		First(&checkout, "checkout_token = ?", token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get checkout by token")
		return nil, err
	}
	return &checkout, nil
}

// UpdateCheckout updates a checkout
func (r *CartRepository) UpdateCheckout(ctx context.Context, checkout *models.Checkout) error {
	if err := r.db.WithContext(ctx).Save(checkout).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update checkout")
		return err
	}
	return nil
}

// CompleteCheckout marks a checkout as completed
func (r *CartRepository) CompleteCheckout(ctx context.Context, checkoutID uuid.UUID, orderID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update checkout
		if err := tx.Model(&models.Checkout{}).
			Where("id = ?", checkoutID).
			Updates(map[string]interface{}{
				"status":         models.CheckoutStatusCompleted,
				"completed_at":   &now,
				"completed_step": models.CheckoutStepComplete,
			}).Error; err != nil {
			return err
		}

		// Get checkout to get cart ID
		var checkout models.Checkout
		if err := tx.First(&checkout, "id = ?", checkoutID).Error; err != nil {
			return err
		}

		// Update cart status
		if err := tx.Model(&models.Cart{}).
			Where("id = ?", checkout.CartID).
			Updates(map[string]interface{}{
				"status":       models.CartStatusCompleted,
				"completed_at": &now,
			}).Error; err != nil {
			return err
		}

		// Create completion event
		event := &models.CheckoutEvent{
			CheckoutID:  checkoutID,
			EventType:   models.CheckoutEventCompleted,
			Description: "Checkout completed",
			Metadata:    map[string]interface{}{"order_id": orderID},
		}
		if err := tx.Create(event).Error; err != nil {
			return err
		}

		return nil
	})
}

// Discount Operations

// ApplyDiscount applies a discount to a cart
func (r *CartRepository) ApplyDiscount(ctx context.Context, discount *models.CartDiscountApplication) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create discount application
		if err := tx.Create(discount).Error; err != nil {
			return err
		}

		// Recalculate totals
		return r.recalculateCartTotals(tx, discount.CartID)
	})
}

// RemoveDiscount removes a discount from a cart
func (r *CartRepository) RemoveDiscount(ctx context.Context, cartID uuid.UUID, discountCode string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete discount application
		if err := tx.Where("cart_id = ? AND code = ?", cartID, discountCode).
			Delete(&models.CartDiscountApplication{}).Error; err != nil {
			return err
		}

		// Recalculate totals
		return r.recalculateCartTotals(tx, cartID)
	})
}

// Shipping Operations

// AddShippingLine adds a shipping line to a cart
func (r *CartRepository) AddShippingLine(ctx context.Context, shippingLine *models.CartShippingLine) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove existing shipping lines
		if err := tx.Where("cart_id = ?", shippingLine.CartID).
			Delete(&models.CartShippingLine{}).Error; err != nil {
			return err
		}

		// Add new shipping line
		if err := tx.Create(shippingLine).Error; err != nil {
			return err
		}

		// Recalculate totals
		return r.recalculateCartTotals(tx, shippingLine.CartID)
	})
}

// Tax Operations

// UpdateTaxLines updates tax lines for a cart
func (r *CartRepository) UpdateTaxLines(ctx context.Context, cartID uuid.UUID, taxLines []models.CartTaxLine) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove existing tax lines
		if err := tx.Where("cart_id = ?", cartID).Delete(&models.CartTaxLine{}).Error; err != nil {
			return err
		}

		// Add new tax lines
		for _, taxLine := range taxLines {
			taxLine.CartID = cartID
			if err := tx.Create(&taxLine).Error; err != nil {
				return err
			}
		}

		// Recalculate totals
		return r.recalculateCartTotals(tx, cartID)
	})
}

// Event Operations

// CreateCheckoutEvent creates a checkout event
func (r *CartRepository) CreateCheckoutEvent(ctx context.Context, event *models.CheckoutEvent) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create checkout event")
		return err
	}
	return nil
}

// Abandonment Operations

// MarkCartAbandoned marks carts as abandoned based on criteria
func (r *CartRepository) MarkCartAbandoned(ctx context.Context, abandonmentThreshold time.Duration) error {
	cutoffTime := time.Now().Add(-abandonmentThreshold)

	return r.db.WithContext(ctx).Model(&models.Cart{}).
		Where("status = ? AND updated_at < ? AND abandoned_at IS NULL", models.CartStatusActive, cutoffTime).
		Updates(map[string]interface{}{
			"status":       models.CartStatusAbandoned,
			"abandoned_at": time.Now(),
		}).Error
}

// GetAbandonedCarts retrieves abandoned carts for recovery campaigns
func (r *CartRepository) GetAbandonedCarts(ctx context.Context, limit, offset int) ([]*models.Cart, error) {
	var carts []*models.Cart
	if err := r.db.WithContext(ctx).
		Preload("LineItems").
		Where("status = ? AND abandoned_at IS NOT NULL", models.CartStatusAbandoned).
		Order("abandoned_at DESC").
		Limit(limit).Offset(offset).
		Find(&carts).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get abandoned carts")
		return nil, err
	}
	return carts, nil
}

// Cleanup Operations

// CleanupExpiredCarts removes expired carts
func (r *CartRepository) CleanupExpiredCarts(ctx context.Context) error {
	now := time.Now()

	return r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Updates(map[string]interface{}{
			"status": models.CartStatusExpired,
		}).Error
}

// Helper Methods

// recalculateCartTotals recalculates cart totals within a transaction
func (r *CartRepository) recalculateCartTotals(tx *gorm.DB, cartID uuid.UUID) error {
	// Calculate subtotal from line items
	var subtotal float64
	if err := tx.Model(&models.CartLineItem{}).
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(line_price - total_discount), 0)").
		Scan(&subtotal).Error; err != nil {
		return err
	}

	// Calculate total tax
	var totalTax float64
	if err := tx.Model(&models.CartTaxLine{}).
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(price), 0)").
		Scan(&totalTax).Error; err != nil {
		return err
	}

	// Calculate total shipping
	var totalShipping float64
	if err := tx.Model(&models.CartShippingLine{}).
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(CASE WHEN discounted_price > 0 THEN discounted_price ELSE price END), 0)").
		Scan(&totalShipping).Error; err != nil {
		return err
	}

	// Calculate total discount from discount applications
	var totalDiscount float64
	if err := tx.Model(&models.CartLineItem{}).
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(total_discount), 0)").
		Scan(&totalDiscount).Error; err != nil {
		return err
	}

	// Update cart totals
	totalPrice := subtotal + totalTax + totalShipping
	return tx.Model(&models.Cart{}).
		Where("id = ?", cartID).
		Updates(map[string]interface{}{
			"subtotal_price": subtotal,
			"total_tax":      totalTax,
			"total_shipping": totalShipping,
			"total_discount": totalDiscount,
			"total_price":    totalPrice,
		}).Error
}

// GetShippingRates retrieves available shipping rates
func (r *CartRepository) GetShippingRates(ctx context.Context) ([]*models.ShippingRate, error) {
	var rates []*models.ShippingRate
	if err := r.db.WithContext(ctx).Find(&rates).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get shipping rates")
		return nil, err
	}
	return rates, nil
}

// GetPaymentMethods retrieves available payment methods
func (r *CartRepository) GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethod, error) {
	var methods []*models.PaymentMethod
	if err := r.db.WithContext(ctx).Where("enabled = ?", true).Find(&methods).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get payment methods")
		return nil, err
	}
	return methods, nil
}
