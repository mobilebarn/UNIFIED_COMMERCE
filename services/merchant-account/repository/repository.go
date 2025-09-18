package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"unified-commerce/services/merchant-account/models"
	"unified-commerce/services/shared/database"
)

// MerchantRepository handles merchant data operations
type MerchantRepository struct {
	db *database.PostgresDB
}

// NewMerchantRepository creates a new merchant repository
func NewMerchantRepository(db *database.PostgresDB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

// Create creates a new merchant
func (r *MerchantRepository) Create(ctx context.Context, merchant *models.Merchant) error {
	return r.db.DB.WithContext(ctx).Create(merchant).Error
}

// GetByID retrieves a merchant by ID
func (r *MerchantRepository) GetByID(ctx context.Context, id string) (*models.Merchant, error) {
	var merchant models.Merchant
	err := r.db.DB.WithContext(ctx).
		Preload("Addresses").
		Preload("BankAccounts").
		Preload("Subscriptions.Plan").
		Preload("Stores").
		Preload("Members").
		First(&merchant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// GetByEmail retrieves a merchant by primary email
func (r *MerchantRepository) GetByEmail(ctx context.Context, email string) (*models.Merchant, error) {
	var merchant models.Merchant
	err := r.db.DB.WithContext(ctx).
		Preload("Addresses").
		Preload("Subscriptions.Plan").
		First(&merchant, "primary_email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// GetByUserID retrieves merchants where the user is a member
func (r *MerchantRepository) GetByUserID(ctx context.Context, userID string) ([]models.Merchant, error) {
	var merchants []models.Merchant
	err := r.db.DB.WithContext(ctx).
		Joins("JOIN merchant_members ON merchants.id = merchant_members.merchant_id").
		Where("merchant_members.user_id = ? AND merchant_members.status = ?", userID, "active").
		Preload("Members").
		Find(&merchants).Error
	return merchants, err
}

// Update updates a merchant
func (r *MerchantRepository) Update(ctx context.Context, merchant *models.Merchant) error {
	return r.db.DB.WithContext(ctx).Save(merchant).Error
}

// UpdateStatus updates merchant status
func (r *MerchantRepository) UpdateStatus(ctx context.Context, merchantID, status string) error {
	return r.db.DB.WithContext(ctx).
		Model(&models.Merchant{}).
		Where("id = ?", merchantID).
		Update("status", status).Error
}

// UpdateVerificationStatus updates merchant verification status
func (r *MerchantRepository) UpdateVerificationStatus(ctx context.Context, merchantID string, isVerified bool) error {
	updates := map[string]interface{}{
		"is_verified": isVerified,
	}
	if isVerified {
		now := time.Now()
		updates["verified_at"] = &now
	}
	
	return r.db.DB.WithContext(ctx).
		Model(&models.Merchant{}).
		Where("id = ?", merchantID).
		Updates(updates).Error
}

// List retrieves merchants with pagination
func (r *MerchantRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]models.Merchant, int64, error) {
	var merchants []models.Merchant
	var total int64

	query := r.db.DB.WithContext(ctx).Model(&models.Merchant{})

	// Apply filters
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if industry, ok := filters["industry"]; ok {
		query = query.Where("industry = ?", industry)
	}
	if verified, ok := filters["verified"]; ok {
		query = query.Where("is_verified = ?", verified)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.
		Preload("Subscriptions.Plan").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&merchants).Error

	return merchants, total, err
}

// Delete soft deletes a merchant
func (r *MerchantRepository) Delete(ctx context.Context, id string) error {
	return r.db.DB.WithContext(ctx).Delete(&models.Merchant{}, "id = ?", id).Error
}

// MerchantAddressRepository handles merchant address operations
type MerchantAddressRepository struct {
	db *database.PostgresDB
}

// NewMerchantAddressRepository creates a new merchant address repository
func NewMerchantAddressRepository(db *database.PostgresDB) *MerchantAddressRepository {
	return &MerchantAddressRepository{db: db}
}

// Create creates a new merchant address
func (r *MerchantAddressRepository) Create(ctx context.Context, address *models.MerchantAddress) error {
	return r.db.DB.WithContext(ctx).Create(address).Error
}

// GetByMerchantID retrieves addresses for a merchant
func (r *MerchantAddressRepository) GetByMerchantID(ctx context.Context, merchantID string) ([]models.MerchantAddress, error) {
	var addresses []models.MerchantAddress
	err := r.db.DB.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Order("is_default DESC, created_at ASC").
		Find(&addresses).Error
	return addresses, err
}

// Update updates a merchant address
func (r *MerchantAddressRepository) Update(ctx context.Context, address *models.MerchantAddress) error {
	return r.db.DB.WithContext(ctx).Save(address).Error
}

// SetDefault sets an address as the default for a specific type
func (r *MerchantAddressRepository) SetDefault(ctx context.Context, addressID, merchantID, addressType string) error {
	return r.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// First, unset all existing defaults for this type
		if err := tx.Model(&models.MerchantAddress{}).
			Where("merchant_id = ? AND type = ?", merchantID, addressType).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Then set the new default
		return tx.Model(&models.MerchantAddress{}).
			Where("id = ?", addressID).
			Update("is_default", true).Error
	})
}

// SubscriptionRepository handles subscription operations
type SubscriptionRepository struct {
	db *database.PostgresDB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *database.PostgresDB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new subscription
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	return r.db.DB.WithContext(ctx).Create(subscription).Error
}

// GetByMerchantID retrieves subscriptions for a merchant
func (r *SubscriptionRepository) GetByMerchantID(ctx context.Context, merchantID string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.DB.WithContext(ctx).
		Preload("Plan").
		Where("merchant_id = ?", merchantID).
		Order("created_at DESC").
		Find(&subscriptions).Error
	return subscriptions, err
}

// GetActiveByMerchantID retrieves the active subscription for a merchant
func (r *SubscriptionRepository) GetActiveByMerchantID(ctx context.Context, merchantID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.DB.WithContext(ctx).
		Preload("Plan").
		Where("merchant_id = ? AND status = ?", merchantID, "active").
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// UpdateStatus updates subscription status
func (r *SubscriptionRepository) UpdateStatus(ctx context.Context, subscriptionID, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == "cancelled" {
		now := time.Now()
		updates["cancelled_at"] = &now
	}

	return r.db.DB.WithContext(ctx).
		Model(&models.Subscription{}).
		Where("id = ?", subscriptionID).
		Updates(updates).Error
}

// PlanRepository handles plan operations
type PlanRepository struct {
	db *database.PostgresDB
}

// NewPlanRepository creates a new plan repository
func NewPlanRepository(db *database.PostgresDB) *PlanRepository {
	return &PlanRepository{db: db}
}

// Create creates a new plan
func (r *PlanRepository) Create(ctx context.Context, plan *models.Plan) error {
	return r.db.DB.WithContext(ctx).Create(plan).Error
}

// GetByID retrieves a plan by ID
func (r *PlanRepository) GetByID(ctx context.Context, id string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.DB.WithContext(ctx).First(&plan, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// GetByName retrieves a plan by name
func (r *PlanRepository) GetByName(ctx context.Context, name string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.DB.WithContext(ctx).First(&plan, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// ListActive retrieves all active plans
func (r *PlanRepository) ListActive(ctx context.Context) ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.DB.WithContext(ctx).
		Where("is_active = ?", true).
		Order("monthly_price ASC").
		Find(&plans).Error
	return plans, err
}

// MerchantMemberRepository handles merchant member operations
type MerchantMemberRepository struct {
	db *database.PostgresDB
}

// NewMerchantMemberRepository creates a new merchant member repository
func NewMerchantMemberRepository(db *database.PostgresDB) *MerchantMemberRepository {
	return &MerchantMemberRepository{db: db}
}

// Create creates a new merchant member
func (r *MerchantMemberRepository) Create(ctx context.Context, member *models.MerchantMember) error {
	return r.db.DB.WithContext(ctx).Create(member).Error
}

// GetByMerchantID retrieves members for a merchant
func (r *MerchantMemberRepository) GetByMerchantID(ctx context.Context, merchantID string) ([]models.MerchantMember, error) {
	var members []models.MerchantMember
	err := r.db.DB.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Order("created_at ASC").
		Find(&members).Error
	return members, err
}

// GetByUserAndMerchant retrieves a specific member relationship
func (r *MerchantMemberRepository) GetByUserAndMerchant(ctx context.Context, userID, merchantID string) (*models.MerchantMember, error) {
	var member models.MerchantMember
	err := r.db.DB.WithContext(ctx).
		Where("user_id = ? AND merchant_id = ?", userID, merchantID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpdateStatus updates member status
func (r *MerchantMemberRepository) UpdateStatus(ctx context.Context, memberID, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == "active" {
		now := time.Now()
		updates["joined_at"] = &now
	}

	return r.db.DB.WithContext(ctx).
		Model(&models.MerchantMember{}).
		Where("id = ?", memberID).
		Updates(updates).Error
}

// Repository aggregates all repositories
type Repository struct {
	Merchant        *MerchantRepository
	MerchantAddress *MerchantAddressRepository
	Subscription    *SubscriptionRepository
	Plan            *PlanRepository
	MerchantMember  *MerchantMemberRepository
}

// NewRepository creates a new repository with all sub-repositories
func NewRepository(db *database.PostgresDB) *Repository {
	return &Repository{
		Merchant:        NewMerchantRepository(db),
		MerchantAddress: NewMerchantAddressRepository(db),
		Subscription:    NewSubscriptionRepository(db),
		Plan:            NewPlanRepository(db),
		MerchantMember:  NewMerchantMemberRepository(db),
	}
}

// Migrate runs database migrations for all models
func (r *Repository) Migrate() error {
	// Configure GORM to disable foreign keys for this migration
	db := r.Merchant.db.DB.Set("gorm:auto_preload", false)
	db = db.Set("gorm:association_autoupdate", false)
	db = db.Set("gorm:association_autocreate", false)

	// Drop ALL foreign key constraints on merchant_members table if they exist
	db.Exec("ALTER TABLE merchant_members DROP CONSTRAINT IF EXISTS fk_users_merchant_members")
	db.Exec("ALTER TABLE merchant_members DROP CONSTRAINT IF EXISTS fk_merchant_members_user_id")
	db.Exec("ALTER TABLE merchant_members DROP CONSTRAINT IF EXISTS fk_merchant_members_invited_by")
	db.Exec("ALTER TABLE merchant_members DROP CONSTRAINT IF EXISTS merchant_members_user_id_fkey")
	db.Exec("ALTER TABLE merchant_members DROP CONSTRAINT IF EXISTS merchant_members_invited_by_fkey")
	
	// Check if table exists and handle it separately
	var tableExists bool
	db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'merchant_members')").Scan(&tableExists)
	
	if tableExists {
		// If table exists, drop and recreate it to avoid constraint issues
		db.Exec("DROP TABLE IF EXISTS merchant_members CASCADE")
	}
	
	// Run migrations with disabled constraints
	err := db.AutoMigrate(
		&models.Merchant{},
		&models.MerchantAddress{},
		&models.BankAccount{},
		&models.Subscription{},
		&models.Plan{},
		&models.Store{},
		&models.MerchantMember{},
		&models.Invoice{},
		&models.InvoiceLineItem{},
		&models.MerchantVerification{},
	)

	return err
}