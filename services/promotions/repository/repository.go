package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"unified-commerce/services/promotions/models"
	"unified-commerce/services/shared/logger"
)

// PromotionsRepository handles database operations for promotions management
type PromotionsRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewPromotionsRepository creates a new promotions repository
func NewPromotionsRepository(db *gorm.DB, logger *logger.Logger) *PromotionsRepository {
	return &PromotionsRepository{
		db:     db,
		logger: logger,
	}
}

// Promotion Operations

// CreatePromotion creates a new promotion
func (r *PromotionsRepository) CreatePromotion(ctx context.Context, promotion *models.Promotion) error {
	if err := r.db.WithContext(ctx).Create(promotion).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create promotion")
		return err
	}
	return nil
}

// GetPromotion retrieves a promotion by ID
func (r *PromotionsRepository) GetPromotion(ctx context.Context, id uuid.UUID) (*models.Promotion, error) {
	var promotion models.Promotion
	if err := r.db.WithContext(ctx).
		Preload("DiscountCodes").
		First(&promotion, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get promotion")
		return nil, err
	}
	return &promotion, nil
}

// GetPromotionsByMerchant retrieves promotions for a merchant
func (r *PromotionsRepository) GetPromotionsByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*models.Promotion, int64, error) {
	var promotions []*models.Promotion
	var total int64

	query := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID)

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Model(&models.Promotion{}).Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count promotions")
		return nil, 0, err
	}

	// Get promotions with pagination
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&promotions).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get promotions")
		return nil, 0, err
	}

	return promotions, total, nil
}

// UpdatePromotion updates a promotion
func (r *PromotionsRepository) UpdatePromotion(ctx context.Context, promotion *models.Promotion) error {
	if err := r.db.WithContext(ctx).Save(promotion).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update promotion")
		return err
	}
	return nil
}

// DeletePromotion soft deletes a promotion
func (r *PromotionsRepository) DeletePromotion(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Promotion{}, "id = ?", id).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete promotion")
		return err
	}
	return nil
}

// Discount Code Operations

// CreateDiscountCode creates a new discount code
func (r *PromotionsRepository) CreateDiscountCode(ctx context.Context, discountCode *models.DiscountCode) error {
	if err := r.db.WithContext(ctx).Create(discountCode).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create discount code")
		return err
	}
	return nil
}

// GetDiscountCode retrieves a discount code by ID
func (r *PromotionsRepository) GetDiscountCode(ctx context.Context, id uuid.UUID) (*models.DiscountCode, error) {
	var discountCode models.DiscountCode
	if err := r.db.WithContext(ctx).
		Preload("Promotion").
		Preload("Usages").
		First(&discountCode, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get discount code")
		return nil, err
	}
	return &discountCode, nil
}

// GetDiscountCodeByCode retrieves a discount code by code
func (r *PromotionsRepository) GetDiscountCodeByCode(ctx context.Context, code string) (*models.DiscountCode, error) {
	var discountCode models.DiscountCode
	if err := r.db.WithContext(ctx).
		Preload("Promotion").
		First(&discountCode, "code = ?", code).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get discount code by code")
		return nil, err
	}
	return &discountCode, nil
}

// GetDiscountCodesByPromotion retrieves discount codes for a promotion
func (r *PromotionsRepository) GetDiscountCodesByPromotion(ctx context.Context, promotionID uuid.UUID) ([]*models.DiscountCode, error) {
	var discountCodes []*models.DiscountCode
	if err := r.db.WithContext(ctx).
		Where("promotion_id = ?", promotionID).
		Find(&discountCodes).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get discount codes by promotion")
		return nil, err
	}
	return discountCodes, nil
}

// UpdateDiscountCode updates a discount code
func (r *PromotionsRepository) UpdateDiscountCode(ctx context.Context, discountCode *models.DiscountCode) error {
	if err := r.db.WithContext(ctx).Save(discountCode).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update discount code")
		return err
	}
	return nil
}

// DeleteDiscountCode deletes a discount code
func (r *PromotionsRepository) DeleteDiscountCode(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.DiscountCode{}, "id = ?", id).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete discount code")
		return err
	}
	return nil
}

// CreateCodeUsage creates a new code usage
func (r *PromotionsRepository) CreateCodeUsage(ctx context.Context, usage *models.CodeUsage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create usage
		if err := tx.Create(usage).Error; err != nil {
			return err
		}

		// Update discount code used count
		if err := tx.Model(&models.DiscountCode{}).
			Where("id = ?", usage.DiscountCodeID).
			Update("used_count", gorm.Expr("used_count + ?", 1)).Error; err != nil {
			return err
		}

		// Update promotion used count
		var discountCode models.DiscountCode
		if err := tx.First(&discountCode, "id = ?", usage.DiscountCodeID).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Promotion{}).
			Where("id = ?", discountCode.PromotionID).
			Update("used_count", gorm.Expr("used_count + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetCodeUsagesByDiscountCode retrieves usages for a discount code
func (r *PromotionsRepository) GetCodeUsagesByDiscountCode(ctx context.Context, discountCodeID uuid.UUID) ([]*models.CodeUsage, error) {
	var usages []*models.CodeUsage
	if err := r.db.WithContext(ctx).
		Where("discount_code_id = ?", discountCodeID).
		Order("created_at DESC").
		Find(&usages).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get code usages by discount code")
		return nil, err
	}
	return usages, nil
}

// Gift Card Operations

// CreateGiftCard creates a new gift card
func (r *PromotionsRepository) CreateGiftCard(ctx context.Context, giftCard *models.GiftCard) error {
	if err := r.db.WithContext(ctx).Create(giftCard).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create gift card")
		return err
	}
	return nil
}

// GetGiftCard retrieves a gift card by ID
func (r *PromotionsRepository) GetGiftCard(ctx context.Context, id uuid.UUID) (*models.GiftCard, error) {
	var giftCard models.GiftCard
	if err := r.db.WithContext(ctx).
		Preload("Transactions").
		First(&giftCard, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get gift card")
		return nil, err
	}
	return &giftCard, nil
}

// GetGiftCardByCode retrieves a gift card by code
func (r *PromotionsRepository) GetGiftCardByCode(ctx context.Context, code string) (*models.GiftCard, error) {
	var giftCard models.GiftCard
	if err := r.db.WithContext(ctx).
		Preload("Transactions").
		First(&giftCard, "code = ?", code).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get gift card by code")
		return nil, err
	}
	return &giftCard, nil
}

// GetGiftCardsByMerchant retrieves gift cards for a merchant
func (r *PromotionsRepository) GetGiftCardsByMerchant(ctx context.Context, merchantID uuid.UUID, filters map[string]interface{}, limit, offset int) ([]*models.GiftCard, int64, error) {
	var giftCards []*models.GiftCard
	var total int64

	query := r.db.WithContext(ctx).Where("merchant_id = ?", merchantID)

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Model(&models.GiftCard{}).Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count gift cards")
		return nil, 0, err
	}

	// Get gift cards with pagination
	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&giftCards).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get gift cards")
		return nil, 0, err
	}

	return giftCards, total, nil
}

// UpdateGiftCard updates a gift card
func (r *PromotionsRepository) UpdateGiftCard(ctx context.Context, giftCard *models.GiftCard) error {
	if err := r.db.WithContext(ctx).Save(giftCard).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update gift card")
		return err
	}
	return nil
}

// CreateGiftCardTransaction creates a new gift card transaction
func (r *PromotionsRepository) CreateGiftCardTransaction(ctx context.Context, transaction *models.GiftCardTransaction) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create transaction
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// Update gift card balance
		if err := tx.Model(&models.GiftCard{}).
			Where("id = ?", transaction.GiftCardID).
			Update("balance", transaction.Balance).Error; err != nil {
			return err
		}

		// Update used_at timestamp if this is a redemption
		if transaction.Type == models.TransactionTypeRedeem {
			if err := tx.Model(&models.GiftCard{}).
				Where("id = ?", transaction.GiftCardID).
				Update("used_at", time.Now()).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetGiftCardTransactions retrieves transactions for a gift card
func (r *PromotionsRepository) GetGiftCardTransactions(ctx context.Context, giftCardID uuid.UUID) ([]*models.GiftCardTransaction, error) {
	var transactions []*models.GiftCardTransaction
	if err := r.db.WithContext(ctx).
		Where("gift_card_id = ?", giftCardID).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get gift card transactions")
		return nil, err
	}
	return transactions, nil
}

// Loyalty Program Operations

// CreateLoyaltyProgram creates a new loyalty program
func (r *PromotionsRepository) CreateLoyaltyProgram(ctx context.Context, program *models.LoyaltyProgram) error {
	if err := r.db.WithContext(ctx).Create(program).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create loyalty program")
		return err
	}
	return nil
}

// GetLoyaltyProgram retrieves a loyalty program by ID
func (r *PromotionsRepository) GetLoyaltyProgram(ctx context.Context, id uuid.UUID) (*models.LoyaltyProgram, error) {
	var program models.LoyaltyProgram
	if err := r.db.WithContext(ctx).
		Preload("Members").
		First(&program, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get loyalty program")
		return nil, err
	}
	return &program, nil
}

// GetLoyaltyProgramsByMerchant retrieves loyalty programs for a merchant
func (r *PromotionsRepository) GetLoyaltyProgramsByMerchant(ctx context.Context, merchantID uuid.UUID) ([]*models.LoyaltyProgram, error) {
	var programs []*models.LoyaltyProgram
	if err := r.db.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Find(&programs).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get loyalty programs by merchant")
		return nil, err
	}
	return programs, nil
}

// UpdateLoyaltyProgram updates a loyalty program
func (r *PromotionsRepository) UpdateLoyaltyProgram(ctx context.Context, program *models.LoyaltyProgram) error {
	if err := r.db.WithContext(ctx).Save(program).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update loyalty program")
		return err
	}
	return nil
}

// Loyalty Member Operations

// CreateLoyaltyMember creates a new loyalty member
func (r *PromotionsRepository) CreateLoyaltyMember(ctx context.Context, member *models.LoyaltyMember) error {
	if err := r.db.WithContext(ctx).Create(member).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create loyalty member")
		return err
	}
	return nil
}

// GetLoyaltyMember retrieves a loyalty member by ID
func (r *PromotionsRepository) GetLoyaltyMember(ctx context.Context, id uuid.UUID) (*models.LoyaltyMember, error) {
	var member models.LoyaltyMember
	if err := r.db.WithContext(ctx).
		Preload("LoyaltyProgram").
		Preload("Tier").
		Preload("Activities").
		First(&member, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get loyalty member")
		return nil, err
	}
	return &member, nil
}

// GetLoyaltyMemberByCustomer retrieves a loyalty member by customer ID and program ID
func (r *PromotionsRepository) GetLoyaltyMemberByCustomer(ctx context.Context, customerID, programID uuid.UUID) (*models.LoyaltyMember, error) {
	var member models.LoyaltyMember
	if err := r.db.WithContext(ctx).
		Preload("LoyaltyProgram").
		Preload("Tier").
		Preload("Activities").
		First(&member, "customer_id = ? AND loyalty_program_id = ?", customerID, programID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get loyalty member by customer")
		return nil, err
	}
	return &member, nil
}

// UpdateLoyaltyMember updates a loyalty member
func (r *PromotionsRepository) UpdateLoyaltyMember(ctx context.Context, member *models.LoyaltyMember) error {
	if err := r.db.WithContext(ctx).Save(member).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update loyalty member")
		return err
	}
	return nil
}

// CreateLoyaltyActivity creates a new loyalty activity
func (r *PromotionsRepository) CreateLoyaltyActivity(ctx context.Context, activity *models.LoyaltyActivity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create activity
		if err := tx.Create(activity).Error; err != nil {
			return err
		}

		// Update member points if this is an earn or redeem activity
		if activity.Type == models.ActivityTypeEarned || activity.Type == models.ActivityTypeRedeemed {
			points := activity.Points
			if activity.Type == models.ActivityTypeRedeemed {
				points = -points
			}

			if err := tx.Model(&models.LoyaltyMember{}).
				Where("id = ?", activity.LoyaltyMemberID).
				Updates(map[string]interface{}{
					"points":           gorm.Expr("points + ?", points),
					"lifetime_points":  gorm.Expr("lifetime_points + ?", points),
					"last_activity_at": time.Now(),
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetLoyaltyActivities retrieves activities for a loyalty member
func (r *PromotionsRepository) GetLoyaltyActivities(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*models.LoyaltyActivity, error) {
	var activities []*models.LoyaltyActivity
	if err := r.db.WithContext(ctx).
		Where("loyalty_member_id = ?", memberID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get loyalty activities")
		return nil, err
	}
	return activities, nil
}

// Loyalty Tier Operations

// CreateLoyaltyTier creates a new loyalty tier
func (r *PromotionsRepository) CreateLoyaltyTier(ctx context.Context, tier *models.LoyaltyTier) error {
	if err := r.db.WithContext(ctx).Create(tier).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create loyalty tier")
		return err
	}
	return nil
}

// GetLoyaltyTiersByProgram retrieves loyalty tiers for a program
func (r *PromotionsRepository) GetLoyaltyTiersByProgram(ctx context.Context, programID uuid.UUID) ([]*models.LoyaltyTier, error) {
	var tiers []*models.LoyaltyTier
	if err := r.db.WithContext(ctx).
		Where("loyalty_program_id = ?", programID).
		Find(&tiers).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get loyalty tiers by program")
		return nil, err
	}
	return tiers, nil
}

// UpdateLoyaltyTier updates a loyalty tier
func (r *PromotionsRepository) UpdateLoyaltyTier(ctx context.Context, tier *models.LoyaltyTier) error {
	if err := r.db.WithContext(ctx).Save(tier).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update loyalty tier")
		return err
	}
	return nil
}
