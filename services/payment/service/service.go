package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"unified-commerce/services/payment/models"
	"unified-commerce/services/payment/repository"
	"unified-commerce/services/shared/logger"
)

// Service errors
var (
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrPaymentMethodNotFound   = errors.New("payment method not found")
	ErrGatewayNotFound         = errors.New("payment gateway not found")
	ErrRefundNotFound          = errors.New("refund not found")
	ErrInvalidPaymentStatus    = errors.New("invalid payment status")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrPaymentAlreadyProcessed = errors.New("payment already processed")
	ErrRefundAmountExceeds     = errors.New("refund amount exceeds payment amount")
	ErrInvalidRefundStatus     = errors.New("invalid refund status")
)

// PaymentService handles business logic for payment management
type PaymentService struct {
	repo   *repository.PaymentRepository
	logger *logger.Logger
}

// NewPaymentService creates a new payment service
func NewPaymentService(repo *repository.PaymentRepository, logger *logger.Logger) *PaymentService {
	return &PaymentService{
		repo:   repo,
		logger: logger,
	}
}

// Payment Operations

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	OrderID         uuid.UUID  `json:"order_id" validate:"required"`
	MerchantID      uuid.UUID  `json:"merchant_id" validate:"required"`
	CustomerID      *uuid.UUID `json:"customer_id"`
	PaymentMethodID uuid.UUID  `json:"payment_method_id" validate:"required"`
	Amount          float64    `json:"amount" validate:"required,min=0"`
	Currency        string     `json:"currency"`
	Description     string     `json:"description"`
}

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*models.Payment, error) {
	// Validate payment method exists
	paymentMethod, err := s.repo.GetPaymentMethod(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, err
	}
	if paymentMethod == nil {
		return nil, ErrPaymentMethodNotFound
	}

	// Get gateway for payment method
	gateway, err := s.repo.GetGateway(ctx, paymentMethod.ID)
	if err != nil {
		return nil, err
	}
	if gateway == nil {
		// Use default gateway if payment method gateway not found
		gateways, err := s.repo.GetAllGateways(ctx)
		if err != nil {
			return nil, err
		}
		if len(gateways) == 0 {
			return nil, ErrGatewayNotFound
		}
		gateway = gateways[0]
	}

	// Set default currency
	if req.Currency == "" {
		req.Currency = "USD"
	}

	// Create payment
	payment := &models.Payment{
		OrderID:         req.OrderID,
		MerchantID:      req.MerchantID,
		CustomerID:      req.CustomerID,
		PaymentMethodID: req.PaymentMethodID,
		GatewayID:       gateway.ID,
		Status:          models.PaymentStatusPending,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Description:     req.Description,
		NetAmount:       req.Amount, // Will be updated after processing
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		s.logger.WithError(err).Error("Failed to create payment")
		return nil, err
	}

	// Create payment event
	event := &models.PaymentEvent{
		PaymentID:   payment.ID,
		EventType:   models.PaymentEventCreated,
		Description: "Payment created",
	}
	if err := s.repo.CreatePaymentEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create payment event")
	}

	s.logger.WithField("payment_id", payment.ID).Info("Payment created successfully")
	return payment, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	payment, err := s.repo.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*models.Payment, error) {
	payment, err := s.repo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

// GetPaymentsByMerchant retrieves payments for a merchant
func (s *PaymentService) GetPaymentsByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, page, limit int) ([]*models.Payment, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetPaymentsByMerchant(ctx, merchantID, filters, limit, offset)
}

// ProcessPayment processes a payment through the gateway
func (s *PaymentService) ProcessPayment(ctx context.Context, paymentID uuid.UUID) error {
	payment, err := s.GetPayment(ctx, paymentID)
	if err != nil {
		return err
	}

	// Check if payment is already processed
	if payment.Status != models.PaymentStatusPending {
		return ErrPaymentAlreadyProcessed
	}

	// In a real implementation, this would integrate with the actual payment gateway
	// For now, we'll simulate the process

	// Update payment status to processing
	now := time.Now()
	if err := s.repo.UpdatePaymentStatus(ctx, paymentID, models.PaymentStatusPending, &now); err != nil {
		s.logger.WithError(err).Error("Failed to update payment status to processing")
		return err
	}

	// Create payment event
	event := &models.PaymentEvent{
		PaymentID:   payment.ID,
		EventType:   models.PaymentEventProcessing,
		Description: "Payment processing started",
	}
	if err := s.repo.CreatePaymentEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create payment event")
	}

	// Simulate gateway processing delay
	time.Sleep(2 * time.Second)

	// Simulate successful payment processing
	// In a real implementation, this would depend on the gateway response
	success := true // This would come from the gateway response

	if success {
		// Update payment status to captured
		if err := s.repo.UpdatePaymentStatus(ctx, paymentID, models.PaymentStatusCaptured, &now); err != nil {
			s.logger.WithError(err).Error("Failed to update payment status to captured")
			return err
		}

		// Update completed timestamp
		payment.CompletedAt = &now
		payment.ProcessedAt = &now
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			s.logger.WithError(err).Error("Failed to update payment completion time")
			return err
		}

		// Create payment event
		event := &models.PaymentEvent{
			PaymentID:   payment.ID,
			EventType:   models.PaymentEventCaptured,
			Description: "Payment captured successfully",
		}
		if err := s.repo.CreatePaymentEvent(ctx, event); err != nil {
			s.logger.WithError(err).Error("Failed to create payment event")
		}

		s.logger.WithField("payment_id", paymentID).Info("Payment processed successfully")
	} else {
		// Update payment status to failed
		if err := s.repo.UpdatePaymentStatus(ctx, paymentID, models.PaymentStatusFailed, &now); err != nil {
			s.logger.WithError(err).Error("Failed to update payment status to failed")
			return err
		}

		// Update failed timestamp
		payment.FailedAt = &now
		if err := s.repo.UpdatePayment(ctx, payment); err != nil {
			s.logger.WithError(err).Error("Failed to update payment failure time")
			return err
		}

		// Create payment event
		event := &models.PaymentEvent{
			PaymentID:   payment.ID,
			EventType:   models.PaymentEventFailed,
			Description: "Payment processing failed",
		}
		if err := s.repo.CreatePaymentEvent(ctx, event); err != nil {
			s.logger.WithError(err).Error("Failed to create payment event")
		}

		s.logger.WithField("payment_id", paymentID).Info("Payment processing failed")
		return errors.New("payment processing failed")
	}

	return nil
}

// CancelPayment cancels a pending payment
func (s *PaymentService) CancelPayment(ctx context.Context, paymentID uuid.UUID) error {
	payment, err := s.GetPayment(ctx, paymentID)
	if err != nil {
		return err
	}

	// Check if payment can be cancelled
	if payment.Status != models.PaymentStatusPending && payment.Status != models.PaymentStatusPending {
		return ErrInvalidPaymentStatus
	}

	now := time.Now()
	if err := s.repo.UpdatePaymentStatus(ctx, paymentID, models.PaymentStatusCancelled, &now); err != nil {
		s.logger.WithError(err).Error("Failed to cancel payment")
		return err
	}

	// Update cancelled timestamp
	payment.CompletedAt = &now
	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		s.logger.WithError(err).Error("Failed to update payment cancellation time")
		return err
	}

	// Create payment event
	event := &models.PaymentEvent{
		PaymentID:   payment.ID,
		EventType:   models.PaymentEventCancelled,
		Description: "Payment cancelled",
	}
	if err := s.repo.CreatePaymentEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create payment event")
	}

	s.logger.WithField("payment_id", paymentID).Info("Payment cancelled successfully")
	return nil
}

// Payment Method Operations

// CreatePaymentMethodRequest represents a request to create a payment method
type CreatePaymentMethodRequest struct {
	CustomerID  *uuid.UUID               `json:"customer_id"`
	MerchantID  *uuid.UUID               `json:"merchant_id"`
	Type        models.PaymentMethodType `json:"type" validate:"required"`
	Provider    string                   `json:"provider"`
	Token       string                   `json:"token"`
	Last4       string                   `json:"last4"`
	ExpiryMonth int                      `json:"expiry_month"`
	ExpiryYear  int                      `json:"expiry_year"`
	Brand       string                   `json:"brand"`
	Name        string                   `json:"name"`
	Email       string                   `json:"email"`
	IsDefault   bool                     `json:"is_default"`
	Metadata    map[string]interface{}   `json:"metadata"`
}

// CreatePaymentMethod creates a new payment method
func (s *PaymentService) CreatePaymentMethod(ctx context.Context, req *CreatePaymentMethodRequest) (*models.PaymentMethod, error) {
	paymentMethod := &models.PaymentMethod{
		CustomerID:  req.CustomerID,
		MerchantID:  req.MerchantID,
		Type:        req.Type,
		Provider:    req.Provider,
		Token:       req.Token,
		Last4:       req.Last4,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		Brand:       req.Brand,
		Name:        req.Name,
		Email:       req.Email,
		IsDefault:   req.IsDefault,
		Metadata:    req.Metadata,
	}

	if err := s.repo.CreatePaymentMethod(ctx, paymentMethod); err != nil {
		s.logger.WithError(err).Error("Failed to create payment method")
		return nil, err
	}

	s.logger.WithField("payment_method_id", paymentMethod.ID).Info("Payment method created successfully")
	return paymentMethod, nil
}

// GetPaymentMethod retrieves a payment method by ID
func (s *PaymentService) GetPaymentMethod(ctx context.Context, id uuid.UUID) (*models.PaymentMethod, error) {
	paymentMethod, err := s.repo.GetPaymentMethod(ctx, id)
	if err != nil {
		return nil, err
	}
	if paymentMethod == nil {
		return nil, ErrPaymentMethodNotFound
	}
	return paymentMethod, nil
}

// GetPaymentMethodsByCustomer retrieves payment methods for a customer
func (s *PaymentService) GetPaymentMethodsByCustomer(ctx context.Context, customerID uuid.UUID) ([]*models.PaymentMethod, error) {
	return s.repo.GetPaymentMethodsByCustomer(ctx, customerID)
}

// UpdatePaymentMethodRequest represents a request to update a payment method
type UpdatePaymentMethodRequest struct {
	IsDefault bool                   `json:"is_default"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// UpdatePaymentMethod updates a payment method
func (s *PaymentService) UpdatePaymentMethod(ctx context.Context, id uuid.UUID, req *UpdatePaymentMethodRequest) (*models.PaymentMethod, error) {
	paymentMethod, err := s.GetPaymentMethod(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.IsDefault {
		paymentMethod.IsDefault = req.IsDefault
	}
	if req.Metadata != nil {
		paymentMethod.Metadata = req.Metadata
	}

	if err := s.repo.UpdatePaymentMethod(ctx, paymentMethod); err != nil {
		s.logger.WithError(err).Error("Failed to update payment method")
		return nil, err
	}

	s.logger.WithField("payment_method_id", id).Info("Payment method updated successfully")
	return paymentMethod, nil
}

// DeletePaymentMethod deletes a payment method
func (s *PaymentService) DeletePaymentMethod(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeletePaymentMethod(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to delete payment method")
		return err
	}

	s.logger.WithField("payment_method_id", id).Info("Payment method deleted successfully")
	return nil
}

// Gateway Operations

// CreateGatewayRequest represents a request to create a payment gateway
type CreateGatewayRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Provider    string                 `json:"provider" validate:"required"`
	IsEnabled   bool                   `json:"is_enabled"`
	IsSandbox   bool                   `json:"is_sandbox"`
	Credentials map[string]interface{} `json:"credentials"`
	Settings    map[string]interface{} `json:"settings"`
}

// CreateGateway creates a new payment gateway
func (s *PaymentService) CreateGateway(ctx context.Context, req *CreateGatewayRequest) (*models.PaymentGateway, error) {
	gateway := &models.PaymentGateway{
		Name:        req.Name,
		Provider:    req.Provider,
		IsEnabled:   req.IsEnabled,
		IsSandbox:   req.IsSandbox,
		Credentials: req.Credentials,
		Settings:    req.Settings,
	}

	if err := s.repo.CreateGateway(ctx, gateway); err != nil {
		s.logger.WithError(err).Error("Failed to create payment gateway")
		return nil, err
	}

	s.logger.WithField("gateway_id", gateway.ID).Info("Payment gateway created successfully")
	return gateway, nil
}

// GetGateway retrieves a payment gateway by ID
func (s *PaymentService) GetGateway(ctx context.Context, id uuid.UUID) (*models.PaymentGateway, error) {
	gateway, err := s.repo.GetGateway(ctx, id)
	if err != nil {
		return nil, err
	}
	if gateway == nil {
		return nil, ErrGatewayNotFound
	}
	return gateway, nil
}

// GetGatewayByName retrieves a payment gateway by name
func (s *PaymentService) GetGatewayByName(ctx context.Context, name string) (*models.PaymentGateway, error) {
	gateway, err := s.repo.GetGatewayByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if gateway == nil {
		return nil, ErrGatewayNotFound
	}
	return gateway, nil
}

// GetAllGateways retrieves all payment gateways
func (s *PaymentService) GetAllGateways(ctx context.Context) ([]*models.PaymentGateway, error) {
	return s.repo.GetAllGateways(ctx)
}

// UpdateGatewayRequest represents a request to update a payment gateway
type UpdateGatewayRequest struct {
	IsEnabled   *bool                  `json:"is_enabled"`
	IsSandbox   *bool                  `json:"is_sandbox"`
	Credentials map[string]interface{} `json:"credentials"`
	Settings    map[string]interface{} `json:"settings"`
}

// UpdateGateway updates a payment gateway
func (s *PaymentService) UpdateGateway(ctx context.Context, id uuid.UUID, req *UpdateGatewayRequest) (*models.PaymentGateway, error) {
	gateway, err := s.GetGateway(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.IsEnabled != nil {
		gateway.IsEnabled = *req.IsEnabled
	}
	if req.IsSandbox != nil {
		gateway.IsSandbox = *req.IsSandbox
	}
	if req.Credentials != nil {
		gateway.Credentials = req.Credentials
	}
	if req.Settings != nil {
		gateway.Settings = req.Settings
	}

	if err := s.repo.UpdateGateway(ctx, gateway); err != nil {
		s.logger.WithError(err).Error("Failed to update payment gateway")
		return nil, err
	}

	s.logger.WithField("gateway_id", id).Info("Payment gateway updated successfully")
	return gateway, nil
}

// Refund Operations

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	PaymentID uuid.UUID           `json:"payment_id" validate:"required"`
	Amount    float64             `json:"amount" validate:"required,min=0"`
	Reason    models.RefundReason `json:"reason"`
}

// CreateRefund creates a new refund
func (s *PaymentService) CreateRefund(ctx context.Context, req *CreateRefundRequest) (*models.Refund, error) {
	// Get payment
	payment, err := s.GetPayment(ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}

	// Check if payment is captured
	if payment.Status != models.PaymentStatusCaptured {
		return nil, ErrInvalidPaymentStatus
	}

	// Check refund amount
	if req.Amount <= 0 {
		return nil, errors.New("refund amount must be greater than 0")
	}

	// Check if refund amount exceeds payment amount
	if req.Amount > payment.Amount {
		return nil, ErrRefundAmountExceeds
	}

	// Calculate total refunded amount
	refunds, err := s.repo.GetRefundsByPayment(ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}

	var totalRefunded float64
	for _, refund := range refunds {
		if refund.Status == models.RefundStatusCompleted || refund.Status == models.RefundStatusProcessed {
			totalRefunded += refund.Amount
		}
	}

	// Check if total refunds exceed payment amount
	if totalRefunded+req.Amount > payment.Amount {
		return nil, ErrRefundAmountExceeds
	}

	// Create refund
	refund := &models.Refund{
		PaymentID: req.PaymentID,
		OrderID:   payment.OrderID,
		Amount:    req.Amount,
		Currency:  payment.Currency,
		Reason:    req.Reason,
		Status:    models.RefundStatusPending,
	}

	if err := s.repo.CreateRefund(ctx, refund); err != nil {
		s.logger.WithError(err).Error("Failed to create refund")
		return nil, err
	}

	// Create payment event
	event := &models.PaymentEvent{
		PaymentID:   payment.ID,
		EventType:   models.PaymentEventRefunded,
		Description: fmt.Sprintf("Refund initiated for amount: %.2f", req.Amount),
	}
	if err := s.repo.CreatePaymentEvent(ctx, event); err != nil {
		s.logger.WithError(err).Error("Failed to create payment event")
	}

	s.logger.WithField("refund_id", refund.ID).WithField("payment_id", req.PaymentID).Info("Refund created successfully")
	return refund, nil
}

// GetRefund retrieves a refund by ID
func (s *PaymentService) GetRefund(ctx context.Context, id uuid.UUID) (*models.Refund, error) {
	refund, err := s.repo.GetRefund(ctx, id)
	if err != nil {
		return nil, err
	}
	if refund == nil {
		return nil, ErrRefundNotFound
	}
	return refund, nil
}

// GetRefundsByPayment retrieves refunds for a payment
func (s *PaymentService) GetRefundsByPayment(ctx context.Context, paymentID uuid.UUID) ([]*models.Refund, error) {
	return s.repo.GetRefundsByPayment(ctx, paymentID)
}

// ProcessRefund processes a refund through the gateway
func (s *PaymentService) ProcessRefund(ctx context.Context, refundID uuid.UUID) error {
	refund, err := s.GetRefund(ctx, refundID)
	if err != nil {
		return err
	}

	// Check if refund is already processed
	if refund.Status != models.RefundStatusPending {
		return ErrInvalidRefundStatus
	}

	// In a real implementation, this would integrate with the actual payment gateway
	// For now, we'll simulate the process

	// Update refund status to processing
	now := time.Now()
	refund.Status = models.RefundStatusProcessed
	refund.ProcessedAt = &now
	if err := s.repo.UpdateRefund(ctx, refund); err != nil {
		s.logger.WithError(err).Error("Failed to update refund status to processed")
		return err
	}

	// Simulate gateway processing delay
	time.Sleep(2 * time.Second)

	// Simulate successful refund processing
	// In a real implementation, this would depend on the gateway response
	success := true // This would come from the gateway response

	if success {
		// Update refund status to completed
		refund.Status = models.RefundStatusCompleted
		refund.CompletedAt = &now
		if err := s.repo.UpdateRefund(ctx, refund); err != nil {
			s.logger.WithError(err).Error("Failed to update refund status to completed")
			return err
		}

		s.logger.WithField("refund_id", refundID).Info("Refund processed successfully")
	} else {
		// Update refund status to failed
		refund.Status = models.RefundStatusFailed
		refund.FailedAt = &now
		if err := s.repo.UpdateRefund(ctx, refund); err != nil {
			s.logger.WithError(err).Error("Failed to update refund status to failed")
			return err
		}

		s.logger.WithField("refund_id", refundID).Info("Refund processing failed")
		return errors.New("refund processing failed")
	}

	return nil
}

// Event Operations

// GetPaymentEvents retrieves events for a payment
func (s *PaymentService) GetPaymentEvents(ctx context.Context, paymentID uuid.UUID) ([]*models.PaymentEvent, error) {
	return s.repo.GetPaymentEvents(ctx, paymentID)
}

// Settlement Operations

// CreateSettlementRequest represents a request to create a settlement
type CreateSettlementRequest struct {
	GatewayID uuid.UUID `json:"gateway_id" validate:"required"`
	Reference string    `json:"reference" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,min=0"`
	Currency  string    `json:"currency"`
}

// CreateSettlement creates a new settlement
func (s *PaymentService) CreateSettlement(ctx context.Context, req *CreateSettlementRequest) (*models.Settlement, error) {
	// Set default currency
	if req.Currency == "" {
		req.Currency = "USD"
	}

	settlement := &models.Settlement{
		GatewayID: req.GatewayID,
		Reference: req.Reference,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    models.SettlementStatusPending,
	}

	if err := s.repo.CreateSettlement(ctx, settlement); err != nil {
		s.logger.WithError(err).Error("Failed to create settlement")
		return nil, err
	}

	s.logger.WithField("settlement_id", settlement.ID).Info("Settlement created successfully")
	return settlement, nil
}

// GetSettlement retrieves a settlement by ID
func (s *PaymentService) GetSettlement(ctx context.Context, id uuid.UUID) (*models.Settlement, error) {
	settlement, err := s.repo.GetSettlement(ctx, id)
	if err != nil {
		return nil, err
	}
	if settlement == nil {
		return nil, errors.New("settlement not found")
	}
	return settlement, nil
}

// GetSettlementsByGateway retrieves settlements for a gateway
func (s *PaymentService) GetSettlementsByGateway(ctx context.Context, gatewayID uuid.UUID, page, limit int) ([]*models.Settlement, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetSettlementsByGateway(ctx, gatewayID, limit, offset)
}

// UpdateSettlementRequest represents a request to update a settlement
type UpdateSettlementRequest struct {
	Status      *models.SettlementStatus `json:"status"`
	DepositedAt *time.Time               `json:"deposited_at"`
}

// UpdateSettlement updates a settlement
func (s *PaymentService) UpdateSettlement(ctx context.Context, id uuid.UUID, req *UpdateSettlementRequest) (*models.Settlement, error) {
	settlement, err := s.GetSettlement(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Status != nil {
		settlement.Status = *req.Status
	}
	if req.DepositedAt != nil {
		settlement.DepositedAt = req.DepositedAt
	}

	now := time.Now()
	settlement.ProcessedAt = &now

	if err := s.repo.UpdateSettlement(ctx, settlement); err != nil {
		s.logger.WithError(err).Error("Failed to update settlement")
		return nil, err
	}

	s.logger.WithField("settlement_id", id).Info("Settlement updated successfully")
	return settlement, nil
}
