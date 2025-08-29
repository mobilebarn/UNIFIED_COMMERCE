package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"unified-commerce/services/promotions/models"
	"unified-commerce/services/promotions/repository"
	"unified-commerce/services/shared/logger"
)

// Service errors
var (
	ErrPromotionNotFound      = errors.New("promotion not found")
	ErrDiscountCodeNotFound   = errors.New("discount code not found")
	ErrGiftCardNotFound       = errors.New("gift card not found")
	ErrLoyaltyProgramNotFound = errors.New("loyalty program not found")
	ErrLoyaltyMemberNotFound  = errors.New("loyalty member not found")
	ErrInvalidCode            = errors.New("invalid discount code")
	ErrCodeExpired            = errors.New("discount code expired")
	ErrCodeInactive           = errors.New("discount code inactive")
	ErrUsageLimitExceeded     = errors.New("usage limit exceeded")
	ErrCustomerUsageLimit     = errors.New("customer usage limit exceeded")
	ErrGiftCardExpired        = errors.New("gift card expired")
	ErrInsufficientBalance    = errors.New("insufficient gift card balance")
	ErrInvalidAmount          = errors.New("invalid amount")
)

// PromotionsService handles business logic for promotions management
type PromotionsService struct {
	repo   *repository.PromotionsRepository
	logger *logger.Logger
}

// NewPromotionsService creates a new promotions service
func NewPromotionsService(repo *repository.PromotionsRepository, logger *logger.Logger) *PromotionsService {
	return &PromotionsService{
		repo:   repo,
		logger: logger,
	}
}

// Promotion Operations

// CreatePromotionRequest represents a request to create a promotion
type CreatePromotionRequest struct {
	MerchantID    uuid.UUID             `json:"merchant_id" validate:"required"`
	Name          string                `json:"name" validate:"required"`
	Description   string                `json:"description"`
	Type          models.PromoType      `json:"type" validate:"required"`
	StartDate     time.Time             `json:"start_date" validate:"required"`
	EndDate       *time.Time            `json:"end_date"`
	UsageLimit    *int                  `json:"usage_limit"`
	Priority      int                   `json:"priority"`
	AppliesTo     models.AppliesTo      `json:"applies_to"`
	Target        models.PromoTarget    `json:"target"`
	Allocation    models.Allocation     `json:"allocation"`
	Prerequisites []models.Prerequisite `json:"prerequisites"`
}

// CreatePromotion creates a new promotion
func (s *PromotionsService) CreatePromotion(ctx context.Context, req *CreatePromotionRequest) (*models.Promotion, error) {
	// Set default status based on dates
	status := models.PromoStatusActive
	if req.StartDate.After(time.Now()) {
		status = models.PromoStatusScheduled
	}

	promotion := &models.Promotion{
		MerchantID:    req.MerchantID,
		Name:          req.Name,
		Description:   req.Description,
		Status:        status,
		Type:          req.Type,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		UsageLimit:    req.UsageLimit,
		Priority:      req.Priority,
		AppliesTo:     req.AppliesTo,
		Target:        req.Target,
		Allocation:    req.Allocation,
		Prerequisites: req.Prerequisites,
	}

	if err := s.repo.CreatePromotion(ctx, promotion); err != nil {
		s.logger.WithError(err).Error("Failed to create promotion")
		return nil, err
	}

	s.logger.WithField("promotion_id", promotion.ID).Info("Promotion created successfully")
	return promotion, nil
}

// GetPromotion retrieves a promotion by ID
func (s *PromotionsService) GetPromotion(ctx context.Context, id uuid.UUID) (*models.Promotion, error) {
	promotion, err := s.repo.GetPromotion(ctx, id)
	if err != nil {
		return nil, err
	}
	if promotion == nil {
		return nil, ErrPromotionNotFound
	}
	return promotion, nil
}

// GetPromotionsByMerchant retrieves promotions for a merchant
func (s *PromotionsService) GetPromotionsByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, page, limit int) ([]*models.Promotion, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetPromotionsByMerchant(ctx, merchantID, filters, limit, offset)
}

// UpdatePromotionRequest represents a request to update a promotion
type UpdatePromotionRequest struct {
	Name          *string                `json:"name"`
	Description   *string                `json:"description"`
	Status        *models.PromoStatus    `json:"status"`
	StartDate     *time.Time             `json:"start_date"`
	EndDate       *time.Time             `json:"end_date"`
	UsageLimit    *int                   `json:"usage_limit"`
	Priority      *int                   `json:"priority"`
	AppliesTo     *models.AppliesTo      `json:"applies_to"`
	Target        *models.PromoTarget    `json:"target"`
	Allocation    *models.Allocation     `json:"allocation"`
	Prerequisites *[]models.Prerequisite `json:"prerequisites"`
}

// UpdatePromotion updates a promotion
func (s *PromotionsService) UpdatePromotion(ctx context.Context, id uuid.UUID, req *UpdatePromotionRequest) (*models.Promotion, error) {
	promotion, err := s.GetPromotion(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		promotion.Name = *req.Name
	}
	if req.Description != nil {
		promotion.Description = *req.Description
	}
	if req.Status != nil {
		promotion.Status = *req.Status
	}
	if req.StartDate != nil {
		promotion.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		promotion.EndDate = req.EndDate
	}
	if req.UsageLimit != nil {
		promotion.UsageLimit = req.UsageLimit
	}
	if req.Priority != nil {
		promotion.Priority = *req.Priority
	}
	if req.AppliesTo != nil {
		promotion.AppliesTo = *req.AppliesTo
	}
	if req.Target != nil {
		promotion.Target = *req.Target
	}
	if req.Allocation != nil {
		promotion.Allocation = *req.Allocation
	}
	if req.Prerequisites != nil {
		promotion.Prerequisites = *req.Prerequisites
	}

	if err := s.repo.UpdatePromotion(ctx, promotion); err != nil {
		s.logger.WithError(err).Error("Failed to update promotion")
		return nil, err
	}

	s.logger.WithField("promotion_id", id).Info("Promotion updated successfully")
	return promotion, nil
}

// DeletePromotion deletes a promotion
func (s *PromotionsService) DeletePromotion(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeletePromotion(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to delete promotion")
		return err
	}

	s.logger.WithField("promotion_id", id).Info("Promotion deleted successfully")
	return nil
}

// Discount Code Operations

// CreateDiscountCodeRequest represents a request to create a discount code
type CreateDiscountCodeRequest struct {
	PromotionID  uuid.UUID         `json:"promotion_id" validate:"required"`
	Code         string            `json:"code"`
	Status       models.CodeStatus `json:"status"`
	UsageLimit   *int              `json:"usage_limit"`
	CustomerUses int               `json:"customer_uses"`
	ExpiresAt    *time.Time        `json:"expires_at"`
}

// CreateDiscountCode creates a new discount code
func (s *PromotionsService) CreateDiscountCode(ctx context.Context, req *CreateDiscountCodeRequest) (*models.DiscountCode, error) {
	// Validate promotion exists
	_, err := s.GetPromotion(ctx, req.PromotionID)
	if err != nil {
		return nil, err
	}

	// Generate code if not provided
	code := req.Code
	if code == "" {
		code = s.generateDiscountCode()
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = models.CodeStatusActive
	}

	discountCode := &models.DiscountCode{
		PromotionID:  req.PromotionID,
		Code:         code,
		Status:       status,
		UsageLimit:   req.UsageLimit,
		CustomerUses: req.CustomerUses,
		ExpiresAt:    req.ExpiresAt,
	}

	if err := s.repo.CreateDiscountCode(ctx, discountCode); err != nil {
		s.logger.WithError(err).Error("Failed to create discount code")
		return nil, err
	}

	s.logger.WithField("discount_code_id", discountCode.ID).Info("Discount code created successfully")
	return discountCode, nil
}

// GetDiscountCode retrieves a discount code by ID
func (s *PromotionsService) GetDiscountCode(ctx context.Context, id uuid.UUID) (*models.DiscountCode, error) {
	discountCode, err := s.repo.GetDiscountCode(ctx, id)
	if err != nil {
		return nil, err
	}
	if discountCode == nil {
		return nil, ErrDiscountCodeNotFound
	}
	return discountCode, nil
}

// GetDiscountCodeByCode retrieves a discount code by code
func (s *PromotionsService) GetDiscountCodeByCode(ctx context.Context, code string) (*models.DiscountCode, error) {
	discountCode, err := s.repo.GetDiscountCodeByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if discountCode == nil {
		return nil, ErrDiscountCodeNotFound
	}
	return discountCode, nil
}

// GetDiscountCodesByPromotion retrieves discount codes for a promotion
func (s *PromotionsService) GetDiscountCodesByPromotion(ctx context.Context, promotionID uuid.UUID) ([]*models.DiscountCode, error) {
	return s.repo.GetDiscountCodesByPromotion(ctx, promotionID)
}

// UpdateDiscountCodeRequest represents a request to update a discount code
type UpdateDiscountCodeRequest struct {
	Status     *models.CodeStatus `json:"status"`
	UsageLimit *int               `json:"usage_limit"`
	ExpiresAt  *time.Time         `json:"expires_at"`
}

// UpdateDiscountCode updates a discount code
func (s *PromotionsService) UpdateDiscountCode(ctx context.Context, id uuid.UUID, req *UpdateDiscountCodeRequest) (*models.DiscountCode, error) {
	discountCode, err := s.GetDiscountCode(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Status != nil {
		discountCode.Status = *req.Status
	}
	if req.UsageLimit != nil {
		discountCode.UsageLimit = req.UsageLimit
	}
	if req.ExpiresAt != nil {
		discountCode.ExpiresAt = req.ExpiresAt
	}

	if err := s.repo.UpdateDiscountCode(ctx, discountCode); err != nil {
		s.logger.WithError(err).Error("Failed to update discount code")
		return nil, err
	}

	s.logger.WithField("discount_code_id", id).Info("Discount code updated successfully")
	return discountCode, nil
}

// DeleteDiscountCode deletes a discount code
func (s *PromotionsService) DeleteDiscountCode(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteDiscountCode(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to delete discount code")
		return err
	}

	s.logger.WithField("discount_code_id", id).Info("Discount code deleted successfully")
	return nil
}

// ValidateDiscountCode validates a discount code and returns the applicable discount
func (s *PromotionsService) ValidateDiscountCode(ctx context.Context, code string, orderAmount float64, customerID *uuid.UUID) (*models.DiscountCode, float64, error) {
	// Get discount code
	discountCode, err := s.GetDiscountCodeByCode(ctx, code)
	if err != nil {
		return nil, 0, ErrInvalidCode
	}

	// Check if code is active
	if discountCode.Status != models.CodeStatusActive {
		return nil, 0, ErrCodeInactive
	}

	// Check if code has expired
	if discountCode.ExpiresAt != nil && discountCode.ExpiresAt.Before(time.Now()) {
		return nil, 0, ErrCodeExpired
	}

	// Check if promotion is active
	promotion := discountCode.Promotion
	if promotion.Status != models.PromoStatusActive {
		return nil, 0, ErrCodeInactive
	}

	// Check if promotion is within date range
	if promotion.StartDate.After(time.Now()) {
		return nil, 0, ErrCodeInactive
	}
	if promotion.EndDate != nil && promotion.EndDate.Before(time.Now()) {
		return nil, 0, ErrCodeExpired
	}

	// Check usage limit
	if discountCode.UsageLimit != nil && discountCode.UsedCount >= *discountCode.UsageLimit {
		return nil, 0, ErrUsageLimitExceeded
	}

	// Check customer usage limit
	if customerID != nil {
		usages, err := s.repo.GetCodeUsagesByDiscountCode(ctx, discountCode.ID)
		if err != nil {
			return nil, 0, err
		}

		customerUsageCount := 0
		for _, usage := range usages {
			if usage.CustomerID != nil && *usage.CustomerID == *customerID {
				customerUsageCount++
			}
		}

		if customerUsageCount >= discountCode.CustomerUses {
			return nil, 0, ErrCustomerUsageLimit
		}
	}

	// Calculate discount amount based on promotion target
	var discountAmount float64
	switch promotion.Target.ValueType {
	case models.ValueTypePercentage:
		discountAmount = orderAmount * (promotion.Target.Value / 100)
	case models.ValueTypeFixed:
		discountAmount = promotion.Target.Value
	case models.ValueTypeFree:
		discountAmount = orderAmount
	}

	return discountCode, discountAmount, nil
}

// RecordDiscountCodeUsage records the usage of a discount code
func (s *PromotionsService) RecordDiscountCodeUsage(ctx context.Context, discountCodeID uuid.UUID, customerID *uuid.UUID, orderID *uuid.UUID, amount float64) error {
	usage := &models.CodeUsage{
		DiscountCodeID: discountCodeID,
		CustomerID:     customerID,
		OrderID:        orderID,
		Amount:         amount,
	}

	if err := s.repo.CreateCodeUsage(ctx, usage); err != nil {
		s.logger.WithError(err).Error("Failed to record discount code usage")
		return err
	}

	s.logger.WithField("usage_id", usage.ID).Info("Discount code usage recorded successfully")
	return nil
}

// Gift Card Operations

// CreateGiftCardRequest represents a request to create a gift card
type CreateGiftCardRequest struct {
	MerchantID uuid.UUID  `json:"merchant_id" validate:"required"`
	Code       string     `json:"code"`
	Balance    float64    `json:"balance" validate:"required,min=0"`
	Currency   string     `json:"currency"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

// CreateGiftCard creates a new gift card
func (s *PromotionsService) CreateGiftCard(ctx context.Context, req *CreateGiftCardRequest) (*models.GiftCard, error) {
	// Generate code if not provided
	code := req.Code
	if code == "" {
		code = s.generateGiftCardCode()
	}

	// Set default currency
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	giftCard := &models.GiftCard{
		MerchantID:     req.MerchantID,
		Code:           code,
		Balance:        req.Balance,
		InitialBalance: req.Balance,
		Currency:       currency,
		Status:         models.GiftCardStatusActive,
		ExpiresAt:      req.ExpiresAt,
		IssuedAt:       time.Now(),
	}

	if err := s.repo.CreateGiftCard(ctx, giftCard); err != nil {
		s.logger.WithError(err).Error("Failed to create gift card")
		return nil, err
	}

	// Create initial transaction
	transaction := &models.GiftCardTransaction{
		GiftCardID:  giftCard.ID,
		Type:        models.TransactionTypeIssue,
		Amount:      req.Balance,
		Balance:     req.Balance,
		Description: "Gift card issued",
	}

	if err := s.repo.CreateGiftCardTransaction(ctx, transaction); err != nil {
		s.logger.WithError(err).Error("Failed to create gift card transaction")
		return nil, err
	}

	s.logger.WithField("gift_card_id", giftCard.ID).Info("Gift card created successfully")
	return giftCard, nil
}

// GetGiftCard retrieves a gift card by ID
func (s *PromotionsService) GetGiftCard(ctx context.Context, id uuid.UUID) (*models.GiftCard, error) {
	giftCard, err := s.repo.GetGiftCard(ctx, id)
	if err != nil {
		return nil, err
	}
	if giftCard == nil {
		return nil, ErrGiftCardNotFound
	}
	return giftCard, nil
}

// GetGiftCardByCode retrieves a gift card by code
func (s *PromotionsService) GetGiftCardByCode(ctx context.Context, code string) (*models.GiftCard, error) {
	giftCard, err := s.repo.GetGiftCardByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if giftCard == nil {
		return nil, ErrGiftCardNotFound
	}
	return giftCard, nil
}

// GetGiftCardsByMerchant retrieves gift cards for a merchant
func (s *PromotionsService) GetGiftCardsByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, page, limit int) ([]*models.GiftCard, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetGiftCardsByMerchant(ctx, merchantID, filters, limit, offset)
}

// UpdateGiftCardRequest represents a request to update a gift card
type UpdateGiftCardRequest struct {
	Status    *models.GiftCardStatus `json:"status"`
	ExpiresAt *time.Time             `json:"expires_at"`
}

// UpdateGiftCard updates a gift card
func (s *PromotionsService) UpdateGiftCard(ctx context.Context, id uuid.UUID, req *UpdateGiftCardRequest) (*models.GiftCard, error) {
	giftCard, err := s.GetGiftCard(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Status != nil {
		giftCard.Status = *req.Status
	}
	if req.ExpiresAt != nil {
		giftCard.ExpiresAt = req.ExpiresAt
	}

	if err := s.repo.UpdateGiftCard(ctx, giftCard); err != nil {
		s.logger.WithError(err).Error("Failed to update gift card")
		return nil, err
	}

	s.logger.WithField("gift_card_id", id).Info("Gift card updated successfully")
	return giftCard, nil
}

// RedeemGiftCard redeems a gift card for an order
func (s *PromotionsService) RedeemGiftCard(ctx context.Context, giftCardID uuid.UUID, orderID *uuid.UUID, amount float64) (*models.GiftCardTransaction, error) {
	// Get gift card
	giftCard, err := s.GetGiftCard(ctx, giftCardID)
	if err != nil {
		return nil, err
	}

	// Check if gift card is active
	if giftCard.Status != models.GiftCardStatusActive {
		return nil, errors.New("gift card is not active")
	}

	// Check if gift card has expired
	if giftCard.ExpiresAt != nil && giftCard.ExpiresAt.Before(time.Now()) {
		return nil, ErrGiftCardExpired
	}

	// Check if amount is valid
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Check if sufficient balance
	if giftCard.Balance < amount {
		return nil, ErrInsufficientBalance
	}

	// Calculate new balance
	newBalance := giftCard.Balance - amount

	// Create transaction
	transaction := &models.GiftCardTransaction{
		GiftCardID:  giftCardID,
		OrderID:     orderID,
		Type:        models.TransactionTypeRedeem,
		Amount:      amount,
		Balance:     newBalance,
		Description: fmt.Sprintf("Redeemed for order %s", orderID.String()),
	}

	if err := s.repo.CreateGiftCardTransaction(ctx, transaction); err != nil {
		s.logger.WithError(err).Error("Failed to create gift card transaction")
		return nil, err
	}

	// Update gift card status if balance is zero
	if newBalance == 0 {
		giftCard.Status = models.GiftCardStatusUsed
		giftCard.UsedAt = &transaction.CreatedAt
		if err := s.repo.UpdateGiftCard(ctx, giftCard); err != nil {
			s.logger.WithError(err).Error("Failed to update gift card status")
			return nil, err
		}
	}

	s.logger.WithField("transaction_id", transaction.ID).Info("Gift card redeemed successfully")
	return transaction, nil
}

// AdjustGiftCardBalance adjusts the balance of a gift card
func (s *PromotionsService) AdjustGiftCardBalance(ctx context.Context, giftCardID uuid.UUID, amount float64, description string) (*models.GiftCardTransaction, error) {
	// Get gift card
	giftCard, err := s.GetGiftCard(ctx, giftCardID)
	if err != nil {
		return nil, err
	}

	// Check if gift card is active
	if giftCard.Status != models.GiftCardStatusActive {
		return nil, errors.New("gift card is not active")
	}

	// Check if amount is valid
	if amount == 0 {
		return nil, ErrInvalidAmount
	}

	// Calculate new balance
	newBalance := giftCard.Balance + amount

	// Determine transaction type
	transactionType := models.TransactionTypeAdjust
	if amount > 0 {
		transactionType = models.TransactionTypeIssue
	} else {
		// Check if sufficient balance for negative adjustment
		if giftCard.Balance < -amount {
			return nil, ErrInsufficientBalance
		}
		transactionType = models.TransactionTypeRedeem
	}

	// Create transaction
	transaction := &models.GiftCardTransaction{
		GiftCardID:  giftCardID,
		Type:        transactionType,
		Amount:      amount,
		Balance:     newBalance,
		Description: description,
	}

	if err := s.repo.CreateGiftCardTransaction(ctx, transaction); err != nil {
		s.logger.WithError(err).Error("Failed to create gift card transaction")
		return nil, err
	}

	s.logger.WithField("transaction_id", transaction.ID).Info("Gift card balance adjusted successfully")
	return transaction, nil
}

// Loyalty Program Operations

// CreateLoyaltyProgramRequest represents a request to create a loyalty program
type CreateLoyaltyProgramRequest struct {
	MerchantID  uuid.UUID              `json:"merchant_id" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	PointValue  float64                `json:"point_value"`
	RewardRatio float64                `json:"reward_ratio"`
	Settings    models.LoyaltySettings `json:"settings"`
}

// CreateLoyaltyProgram creates a new loyalty program
func (s *PromotionsService) CreateLoyaltyProgram(ctx context.Context, req *CreateLoyaltyProgramRequest) (*models.LoyaltyProgram, error) {
	// Set default values
	pointValue := req.PointValue
	if pointValue == 0 {
		pointValue = 1.00
	}

	rewardRatio := req.RewardRatio
	if rewardRatio == 0 {
		rewardRatio = 100.00
	}

	program := &models.LoyaltyProgram{
		MerchantID:  req.MerchantID,
		Name:        req.Name,
		Description: req.Description,
		Status:      models.ProgramStatusActive,
		PointValue:  pointValue,
		RewardRatio: rewardRatio,
		Settings:    req.Settings,
	}

	if err := s.repo.CreateLoyaltyProgram(ctx, program); err != nil {
		s.logger.WithError(err).Error("Failed to create loyalty program")
		return nil, err
	}

	s.logger.WithField("program_id", program.ID).Info("Loyalty program created successfully")
	return program, nil
}

// GetLoyaltyProgram retrieves a loyalty program by ID
func (s *PromotionsService) GetLoyaltyProgram(ctx context.Context, id uuid.UUID) (*models.LoyaltyProgram, error) {
	program, err := s.repo.GetLoyaltyProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, ErrLoyaltyProgramNotFound
	}
	return program, nil
}

// GetLoyaltyProgramsByMerchant retrieves loyalty programs for a merchant
func (s *PromotionsService) GetLoyaltyProgramsByMerchant(ctx context.Context, merchantID uuid.UUID) ([]*models.LoyaltyProgram, error) {
	return s.repo.GetLoyaltyProgramsByMerchant(ctx, merchantID)
}

// UpdateLoyaltyProgramRequest represents a request to update a loyalty program
type UpdateLoyaltyProgramRequest struct {
	Name        *string                 `json:"name"`
	Description *string                 `json:"description"`
	Status      *models.ProgramStatus   `json:"status"`
	PointValue  *float64                `json:"point_value"`
	RewardRatio *float64                `json:"reward_ratio"`
	Settings    *models.LoyaltySettings `json:"settings"`
}

// UpdateLoyaltyProgram updates a loyalty program
func (s *PromotionsService) UpdateLoyaltyProgram(ctx context.Context, id uuid.UUID, req *UpdateLoyaltyProgramRequest) (*models.LoyaltyProgram, error) {
	program, err := s.GetLoyaltyProgram(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		program.Name = *req.Name
	}
	if req.Description != nil {
		program.Description = *req.Description
	}
	if req.Status != nil {
		program.Status = *req.Status
	}
	if req.PointValue != nil {
		program.PointValue = *req.PointValue
	}
	if req.RewardRatio != nil {
		program.RewardRatio = *req.RewardRatio
	}
	if req.Settings != nil {
		program.Settings = *req.Settings
	}

	if err := s.repo.UpdateLoyaltyProgram(ctx, program); err != nil {
		s.logger.WithError(err).Error("Failed to update loyalty program")
		return nil, err
	}

	s.logger.WithField("program_id", id).Info("Loyalty program updated successfully")
	return program, nil
}

// Loyalty Member Operations

// CreateLoyaltyMemberRequest represents a request to create a loyalty member
type CreateLoyaltyMemberRequest struct {
	LoyaltyProgramID uuid.UUID `json:"loyalty_program_id" validate:"required"`
	CustomerID       uuid.UUID `json:"customer_id" validate:"required"`
}

// CreateLoyaltyMember creates a new loyalty member
func (s *PromotionsService) CreateLoyaltyMember(ctx context.Context, req *CreateLoyaltyMemberRequest) (*models.LoyaltyMember, error) {
	member := &models.LoyaltyMember{
		LoyaltyProgramID: req.LoyaltyProgramID,
		CustomerID:       req.CustomerID,
		Status:           models.MemberStatusActive,
		EnrolledAt:       time.Now(),
	}

	if err := s.repo.CreateLoyaltyMember(ctx, member); err != nil {
		s.logger.WithError(err).Error("Failed to create loyalty member")
		return nil, err
	}

	// Record enrollment activity
	activity := &models.LoyaltyActivity{
		LoyaltyMemberID: member.ID,
		Type:            models.ActivityTypeEarned,
		Points:          0,
		Description:     "Joined loyalty program",
	}

	if err := s.repo.CreateLoyaltyActivity(ctx, activity); err != nil {
		s.logger.WithError(err).Error("Failed to create loyalty activity")
		return nil, err
	}

	s.logger.WithField("member_id", member.ID).Info("Loyalty member created successfully")
	return member, nil
}

// GetLoyaltyMember retrieves a loyalty member by ID
func (s *PromotionsService) GetLoyaltyMember(ctx context.Context, id uuid.UUID) (*models.LoyaltyMember, error) {
	member, err := s.repo.GetLoyaltyMember(ctx, id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrLoyaltyMemberNotFound
	}
	return member, nil
}

// GetLoyaltyMemberByCustomer retrieves a loyalty member by customer ID and program ID
func (s *PromotionsService) GetLoyaltyMemberByCustomer(ctx context.Context, customerID, programID uuid.UUID) (*models.LoyaltyMember, error) {
	member, err := s.repo.GetLoyaltyMemberByCustomer(ctx, customerID, programID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrLoyaltyMemberNotFound
	}
	return member, nil
}

// UpdateLoyaltyMemberRequest represents a request to update a loyalty member
type UpdateLoyaltyMemberRequest struct {
	Status *models.MemberStatus `json:"status"`
	TierID *uuid.UUID           `json:"tier_id"`
}

// UpdateLoyaltyMember updates a loyalty member
func (s *PromotionsService) UpdateLoyaltyMember(ctx context.Context, id uuid.UUID, req *UpdateLoyaltyMemberRequest) (*models.LoyaltyMember, error) {
	member, err := s.GetLoyaltyMember(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Status != nil {
		member.Status = *req.Status
	}
	if req.TierID != nil {
		member.TierID = req.TierID
	}

	if err := s.repo.UpdateLoyaltyMember(ctx, member); err != nil {
		s.logger.WithError(err).Error("Failed to update loyalty member")
		return nil, err
	}

	s.logger.WithField("member_id", id).Info("Loyalty member updated successfully")
	return member, nil
}

// EarnLoyaltyPoints earns points for a loyalty member
func (s *PromotionsService) EarnLoyaltyPoints(ctx context.Context, memberID uuid.UUID, points int, description string, referenceID *uuid.UUID) error {
	// Get member
	member, err := s.GetLoyaltyMember(ctx, memberID)
	if err != nil {
		return err
	}

	// Check if member is active
	if member.Status != models.MemberStatusActive {
		return errors.New("loyalty member is not active")
	}

	// Record activity
	activity := &models.LoyaltyActivity{
		LoyaltyMemberID: memberID,
		Type:            models.ActivityTypeEarned,
		Points:          points,
		Description:     description,
		ReferenceID:     referenceID,
	}

	if err := s.repo.CreateLoyaltyActivity(ctx, activity); err != nil {
		s.logger.WithError(err).Error("Failed to create loyalty activity")
		return err
	}

	s.logger.WithField("activity_id", activity.ID).WithField("points", points).Info("Loyalty points earned successfully")
	return nil
}

// RedeemLoyaltyPoints redeems points for a loyalty member
func (s *PromotionsService) RedeemLoyaltyPoints(ctx context.Context, memberID uuid.UUID, points int, description string, referenceID *uuid.UUID) error {
	// Get member
	member, err := s.GetLoyaltyMember(ctx, memberID)
	if err != nil {
		return err
	}

	// Check if member is active
	if member.Status != models.MemberStatusActive {
		return errors.New("loyalty member is not active")
	}

	// Check if sufficient points
	if member.Points < points {
		return errors.New("insufficient loyalty points")
	}

	// Record activity
	activity := &models.LoyaltyActivity{
		LoyaltyMemberID: memberID,
		Type:            models.ActivityTypeRedeemed,
		Points:          points,
		Description:     description,
		ReferenceID:     referenceID,
	}

	if err := s.repo.CreateLoyaltyActivity(ctx, activity); err != nil {
		s.logger.WithError(err).Error("Failed to create loyalty activity")
		return err
	}

	s.logger.WithField("activity_id", activity.ID).WithField("points", points).Info("Loyalty points redeemed successfully")
	return nil
}

// GetLoyaltyActivities retrieves activities for a loyalty member
func (s *PromotionsService) GetLoyaltyActivities(ctx context.Context, memberID uuid.UUID, page, limit int) ([]*models.LoyaltyActivity, error) {
	return s.repo.GetLoyaltyActivities(ctx, memberID, limit, (page-1)*limit)
}

// Loyalty Tier Operations

// CreateLoyaltyTierRequest represents a request to create a loyalty tier
type CreateLoyaltyTierRequest struct {
	LoyaltyProgramID uuid.UUID        `json:"loyalty_program_id" validate:"required"`
	Name             string           `json:"name" validate:"required"`
	Description      string           `json:"description"`
	MinimumPoints    int              `json:"minimum_points" validate:"required"`
	Benefits         []models.Benefit `json:"benefits"`
}

// CreateLoyaltyTier creates a new loyalty tier
func (s *PromotionsService) CreateLoyaltyTier(ctx context.Context, req *CreateLoyaltyTierRequest) (*models.LoyaltyTier, error) {
	tier := &models.LoyaltyTier{
		LoyaltyProgramID: req.LoyaltyProgramID,
		Name:             req.Name,
		Description:      req.Description,
		MinimumPoints:    req.MinimumPoints,
		Benefits:         req.Benefits,
	}

	if err := s.repo.CreateLoyaltyTier(ctx, tier); err != nil {
		s.logger.WithError(err).Error("Failed to create loyalty tier")
		return nil, err
	}

	s.logger.WithField("tier_id", tier.ID).Info("Loyalty tier created successfully")
	return tier, nil
}

// GetLoyaltyTiersByProgram retrieves loyalty tiers for a program
func (s *PromotionsService) GetLoyaltyTiersByProgram(ctx context.Context, programID uuid.UUID) ([]*models.LoyaltyTier, error) {
	return s.repo.GetLoyaltyTiersByProgram(ctx, programID)
}

// UpdateLoyaltyTierRequest represents a request to update a loyalty tier
type UpdateLoyaltyTierRequest struct {
	Name          *string           `json:"name"`
	Description   *string           `json:"description"`
	MinimumPoints *int              `json:"minimum_points"`
	Benefits      *[]models.Benefit `json:"benefits"`
}

// UpdateLoyaltyTier updates a loyalty tier
func (s *PromotionsService) UpdateLoyaltyTier(ctx context.Context, id uuid.UUID, req *UpdateLoyaltyTierRequest) (*models.LoyaltyTier, error) {
	// Get tier
	tiers, err := s.repo.GetLoyaltyTiersByProgram(ctx, id)
	if err != nil {
		return nil, err
	}

	var tier *models.LoyaltyTier
	for _, t := range tiers {
		if t.ID == id {
			tier = t
			break
		}
	}

	if tier == nil {
		return nil, errors.New("loyalty tier not found")
	}

	// Update fields
	if req.Name != nil {
		tier.Name = *req.Name
	}
	if req.Description != nil {
		tier.Description = *req.Description
	}
	if req.MinimumPoints != nil {
		tier.MinimumPoints = *req.MinimumPoints
	}
	if req.Benefits != nil {
		tier.Benefits = *req.Benefits
	}

	if err := s.repo.UpdateLoyaltyTier(ctx, tier); err != nil {
		s.logger.WithError(err).Error("Failed to update loyalty tier")
		return nil, err
	}

	s.logger.WithField("tier_id", id).Info("Loyalty tier updated successfully")
	return tier, nil
}

// Utility Methods

// generateDiscountCode generates a unique discount code
func (s *PromotionsService) generateDiscountCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// generateGiftCardCode generates a unique gift card code
func (s *PromotionsService) generateGiftCardCode() string {
	const prefix = "GC"
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 10)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return prefix + string(code)
}

// CalculateLoyaltyPoints calculates loyalty points based on purchase amount
func (s *PromotionsService) CalculateLoyaltyPoints(program *models.LoyaltyProgram, amount float64) int {
	if !program.Settings.EarnOnPurchase {
		return 0
	}

	if amount < program.Settings.MinimumPurchaseAmount {
		return 0
	}

	points := int(amount / program.PointValue * program.RewardRatio / 100)

	// Apply rounding
	switch program.Settings.PointsRounding {
	case models.RoundingUp:
		// Already handled by int conversion
	case models.RoundingDown:
		points = int(amount / program.PointValue * program.RewardRatio / 100)
	case models.RoundingHalf:
		points = int(amount / program.PointValue * program.RewardRatio / 100)
	case models.RoundingExact:
		// Already handled
	}

	return points
}
