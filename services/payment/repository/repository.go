package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"unified-commerce/services/payment/models"
	"unified-commerce/services/shared/logger"
)

// PaymentRepository handles database operations for payment management
type PaymentRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB, logger *logger.Logger) *PaymentRepository {
	return &PaymentRepository{
		db:     db,
		logger: logger,
	}
}

// Payment Operations

// CreatePayment creates a new payment
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	if err := r.db.WithContext(ctx).Create(payment).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create payment")
		return err
	}
	return nil
}

// GetPayment retrieves a payment by ID with all related data
func (r *PaymentRepository) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Preload("Gateway").
		Preload("Refunds").
		Preload("Events").
		First(&payment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get payment")
		return nil, err
	}
	return &payment, nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Preload("Gateway").
		Preload("Refunds").
		Preload("Events").
		First(&payment, "order_id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get payment by order ID")
		return nil, err
	}
	return &payment, nil
}

// GetPaymentsByMerchant retrieves payments for a merchant
func (r *PaymentRepository) GetPaymentsByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*models.Payment, int64, error) {
	var payments []*models.Payment
	var total int64

	query := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Preload("Gateway").
		Where("merchant_id = ?", merchantID)

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Model(&models.Payment{}).Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count payments")
		return nil, 0, err
	}

	// Get payments with pagination
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get payments")
		return nil, 0, err
	}

	return payments, total, nil
}

// UpdatePayment updates a payment
func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	if err := r.db.WithContext(ctx).Save(payment).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update payment")
		return err
	}
	return nil
}

// UpdatePaymentStatus updates the status of a payment
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status models.PaymentStatus, processedAt *time.Time) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if processedAt != nil {
		updates["processed_at"] = processedAt
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update payment status")
		return err
	}

	return nil
}

// Payment Method Operations

// CreatePaymentMethod creates a new payment method
func (r *PaymentRepository) CreatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	if err := r.db.WithContext(ctx).Create(paymentMethod).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create payment method")
		return err
	}
	return nil
}

// GetPaymentMethod retrieves a payment method by ID
func (r *PaymentRepository) GetPaymentMethod(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	if err := r.db.WithContext(ctx).First(&paymentMethod, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get payment method")
		return nil, err
	}
	return &paymentMethod, nil
}

// GetPaymentMethodsByCustomer retrieves payment methods for a customer
func (r *PaymentRepository) GetPaymentMethodsByCustomer(ctx context.Context, customerID uuid.UUID) ([]*models.PaymentMethod, error) {
	var paymentMethods []*models.PaymentMethod
	if err := r.db.WithContext(ctx).
		Where("customer_id = ? AND is_active = ?", customerID, true).
		Order("is_default DESC, created_at DESC").
		Find(&paymentMethods).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get payment methods by customer")
		return nil, err
	}
	return paymentMethods, nil
}

// UpdatePaymentMethod updates a payment method
func (r *PaymentRepository) UpdatePaymentMethod(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	if err := r.db.WithContext(ctx).Save(paymentMethod).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update payment method")
		return err
	}
	return nil
}

// DeletePaymentMethod soft deletes a payment method
func (r *PaymentRepository) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&models.PaymentMethod{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": false,
		}).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete payment method")
		return err
	}
	return nil
}

// Gateway Operations

// CreateGateway creates a new payment gateway
func (r *PaymentRepository) CreateGateway(ctx context.Context, gateway *models.PaymentGateway) error {
	if err := r.db.WithContext(ctx).Create(gateway).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create payment gateway")
		return err
	}
	return nil
}

// GetGateway retrieves a payment gateway by ID
func (r *PaymentRepository) GetGateway(ctx context.Context, id uuid.UUID) (*models.PaymentGateway, error) {
	var gateway models.PaymentGateway
	if err := r.db.WithContext(ctx).First(&gateway, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get payment gateway")
		return nil, err
	}
	return &gateway, nil
}

// GetGatewayByName retrieves a payment gateway by name
func (r *PaymentRepository) GetGatewayByName(ctx context.Context, name string) (*models.PaymentGateway, error) {
	var gateway models.PaymentGateway
	if err := r.db.WithContext(ctx).First(&gateway, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get payment gateway by name")
		return nil, err
	}
	return &gateway, nil
}

// GetAllGateways retrieves all payment gateways
func (r *PaymentRepository) GetAllGateways(ctx context.Context) ([]*models.PaymentGateway, error) {
	var gateways []*models.PaymentGateway
	if err := r.db.WithContext(ctx).
		Where("is_enabled = ?", true).
		Find(&gateways).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get payment gateways")
		return nil, err
	}
	return gateways, nil
}

// UpdateGateway updates a payment gateway
func (r *PaymentRepository) UpdateGateway(ctx context.Context, gateway *models.PaymentGateway) error {
	if err := r.db.WithContext(ctx).Save(gateway).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update payment gateway")
		return err
	}
	return nil
}

// Refund Operations

// CreateRefund creates a new refund
func (r *PaymentRepository) CreateRefund(ctx context.Context, refund *models.Refund) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create refund
		if err := tx.Create(refund).Error; err != nil {
			return err
		}

		// Update payment status if fully refunded
		var payment models.Payment
		if err := tx.First(&payment, "id = ?", refund.PaymentID).Error; err != nil {
			return err
		}

		// Calculate total refunded amount
		var totalRefunded float64
		if err := tx.Model(&models.Refund{}).
			Where("payment_id = ? AND status IN ?", refund.PaymentID, []models.RefundStatus{models.RefundStatusCompleted, models.RefundStatusProcessed}).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&totalRefunded).Error; err != nil {
			return err
		}

		// Update payment status based on refund amount
		if totalRefunded >= payment.Amount {
			payment.Status = models.PaymentStatusRefunded
		} else if totalRefunded > 0 {
			payment.Status = models.PaymentStatusPartiallyRefunded
		}

		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetRefund retrieves a refund by ID
func (r *PaymentRepository) GetRefund(ctx context.Context, id uuid.UUID) (*models.Refund, error) {
	var refund models.Refund
	if err := r.db.WithContext(ctx).
		Preload("Payment").
		First(&refund, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get refund")
		return nil, err
	}
	return &refund, nil
}

// GetRefundsByPayment retrieves refunds for a payment
func (r *PaymentRepository) GetRefundsByPayment(ctx context.Context, paymentID uuid.UUID) ([]*models.Refund, error) {
	var refunds []*models.Refund
	if err := r.db.WithContext(ctx).
		Preload("Payment").
		Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&refunds).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get refunds by payment")
		return nil, err
	}
	return refunds, nil
}

// UpdateRefund updates a refund
func (r *PaymentRepository) UpdateRefund(ctx context.Context, refund *models.Refund) error {
	if err := r.db.WithContext(ctx).Save(refund).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update refund")
		return err
	}
	return nil
}

// Event Operations

// CreatePaymentEvent creates a payment event
func (r *PaymentRepository) CreatePaymentEvent(ctx context.Context, event *models.PaymentEvent) error {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create payment event")
		return err
	}
	return nil
}

// GetPaymentEvents retrieves events for a payment
func (r *PaymentRepository) GetPaymentEvents(ctx context.Context, paymentID uuid.UUID) ([]*models.PaymentEvent, error) {
	var events []*models.PaymentEvent
	if err := r.db.WithContext(ctx).
		Where("payment_id = ?", paymentID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get payment events")
		return nil, err
	}
	return events, nil
}

// Settlement Operations

// CreateSettlement creates a new settlement
func (r *PaymentRepository) CreateSettlement(ctx context.Context, settlement *models.Settlement) error {
	if err := r.db.WithContext(ctx).Create(settlement).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create settlement")
		return err
	}
	return nil
}

// GetSettlement retrieves a settlement by ID
func (r *PaymentRepository) GetSettlement(ctx context.Context, id uuid.UUID) (*models.Settlement, error) {
	var settlement models.Settlement
	if err := r.db.WithContext(ctx).
		Preload("Gateway").
		First(&settlement, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get settlement")
		return nil, err
	}
	return &settlement, nil
}

// GetSettlementsByGateway retrieves settlements for a gateway
func (r *PaymentRepository) GetSettlementsByGateway(ctx context.Context, gatewayID uuid.UUID, limit, offset int) ([]*models.Settlement, int64, error) {
	var settlements []*models.Settlement
	var total int64

	query := r.db.WithContext(ctx).Where("gateway_id = ?", gatewayID)

	// Get total count
	if err := query.Model(&models.Settlement{}).Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count settlements")
		return nil, 0, err
	}

	// Get settlements with pagination
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&settlements).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get settlements")
		return nil, 0, err
	}

	return settlements, total, nil
}

// UpdateSettlement updates a settlement
func (r *PaymentRepository) UpdateSettlement(ctx context.Context, settlement *models.Settlement) error {
	if err := r.db.WithContext(ctx).Save(settlement).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update settlement")
		return err
	}
	return nil
}
