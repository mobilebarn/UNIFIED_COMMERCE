package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Promotion represents a promotional campaign
type Promotion struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MerchantID    uuid.UUID      `json:"merchant_id" gorm:"type:uuid;not null;index"`
	Name          string         `json:"name" gorm:"not null"`
	Description   string         `json:"description"`
	Status        PromoStatus    `json:"status" gorm:"default:'active'"`
	Type          PromoType      `json:"type" gorm:"not null"`
	Priority      int            `json:"priority" gorm:"default:0"`
	StartDate     time.Time      `json:"start_date" gorm:"not null"`
	EndDate       *time.Time     `json:"end_date"`
	UsageLimit    *int           `json:"usage_limit"`
	UsedCount     int            `json:"used_count" gorm:"default:0"`
	AppliesTo     AppliesTo      `json:"applies_to" gorm:"type:jsonb"`
	Target        PromoTarget    `json:"target" gorm:"type:jsonb"`
	Allocation    Allocation     `json:"allocation" gorm:"type:jsonb"`
	Prerequisites []Prerequisite `json:"prerequisites,omitempty" gorm:"type:jsonb"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	DiscountCodes []DiscountCode `json:"discount_codes,omitempty" gorm:"foreignKey:PromotionID"`
}

// PromoStatus represents the status of a promotion
type PromoStatus string

const (
	PromoStatusActive    PromoStatus = "active"
	PromoStatusInactive  PromoStatus = "inactive"
	PromoStatusExpired   PromoStatus = "expired"
	PromoStatusScheduled PromoStatus = "scheduled"
)

// PromoType represents the type of promotion
type PromoType string

const (
	PromoTypeDiscount      PromoType = "discount"
	PromoTypeBuyXGetY      PromoType = "buy_x_get_y"
	PromoTypeFreeShipping  PromoType = "free_shipping"
	PromoTypeBOGO          PromoType = "bogo"
	PromoTypeVolume        PromoType = "volume"
	PromoTypeGiftCard      PromoType = "gift_card"
	PromoTypeLoyaltyPoints PromoType = "loyalty_points"
)

// AppliesTo represents what the promotion applies to
type AppliesTo struct {
	Products    []uuid.UUID `json:"products,omitempty"`
	Categories  []uuid.UUID `json:"categories,omitempty"`
	Collections []uuid.UUID `json:"collections,omitempty"`
	AllProducts bool        `json:"all_products"`
}

// PromoTarget represents the target of the promotion
type PromoTarget struct {
	Type      TargetType `json:"type"`
	Value     float64    `json:"value"`
	ValueType ValueType  `json:"value_type"`
}

// TargetType represents the type of target
type TargetType string

const (
	TargetTypeOrder    TargetType = "order"
	TargetTypeProduct  TargetType = "product"
	TargetTypeShipping TargetType = "shipping"
	TargetTypeCustomer TargetType = "customer"
)

// ValueType represents the type of value
type ValueType string

const (
	ValueTypePercentage ValueType = "percentage"
	ValueTypeFixed      ValueType = "fixed"
	ValueTypeFree       ValueType = "free"
)

// Allocation represents how the discount is allocated
type Allocation struct {
	Method AllocationMethod `json:"method"`
}

// AllocationMethod represents the allocation method
type AllocationMethod string

const (
	AllocationMethodEach     AllocationMethod = "each"
	AllocationMethodAcross   AllocationMethod = "across"
	AllocationMethodCustomer AllocationMethod = "customer"
)

// Prerequisite represents a prerequisite for a promotion
type Prerequisite struct {
	Type       PrerequisiteType `json:"type"`
	TargetType TargetType       `json:"target_type"`
	Value      float64          `json:"value"`
	Condition  Condition        `json:"condition"`
}

// PrerequisiteType represents the type of prerequisite
type PrerequisiteType string

const (
	PrerequisiteTypeMinimumOrderAmount PrerequisiteType = "minimum_order_amount"
	PrerequisiteTypeMinimumQuantity    PrerequisiteType = "minimum_quantity"
	PrerequisiteTypeCustomerGroup      PrerequisiteType = "customer_group"
	PrerequisiteTypeProductPurchase    PrerequisiteType = "product_purchase"
)

// Condition represents the condition for a prerequisite
type Condition string

const (
	ConditionGreaterOrEqual Condition = "greater_or_equal"
	ConditionLessOrEqual    Condition = "less_or_equal"
	ConditionEqual          Condition = "equal"
)

// DiscountCode represents a discount code
type DiscountCode struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PromotionID  uuid.UUID  `json:"promotion_id" gorm:"type:uuid;not null;index"`
	Code         string     `json:"code" gorm:"unique;not null;index"`
	Status       CodeStatus `json:"status" gorm:"default:'active'"`
	UsageLimit   *int       `json:"usage_limit"`
	UsedCount    int        `json:"used_count" gorm:"default:0"`
	CustomerUses int        `json:"customer_uses" gorm:"default:1"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Promotion Promotion   `json:"promotion,omitempty" gorm:"foreignKey:PromotionID"`
	Usages    []CodeUsage `json:"usages,omitempty" gorm:"foreignKey:DiscountCodeID"`
}

// CodeStatus represents the status of a discount code
type CodeStatus string

const (
	CodeStatusActive   CodeStatus = "active"
	CodeStatusInactive CodeStatus = "inactive"
	CodeStatusExpired  CodeStatus = "expired"
)

// CodeUsage represents a usage of a discount code
type CodeUsage struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DiscountCodeID uuid.UUID  `json:"discount_code_id" gorm:"type:uuid;not null;index"`
	CustomerID     *uuid.UUID `json:"customer_id" gorm:"type:uuid;index"`
	OrderID        *uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	Amount         float64    `json:"amount" gorm:"type:decimal(12,2)"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	DiscountCode DiscountCode `json:"discount_code,omitempty" gorm:"foreignKey:DiscountCodeID"`
}

// GiftCard represents a gift card
type GiftCard struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MerchantID     uuid.UUID      `json:"merchant_id" gorm:"type:uuid;not null;index"`
	Code           string         `json:"code" gorm:"unique;not null;index"`
	Balance        float64        `json:"balance" gorm:"type:decimal(12,2);not null"`
	InitialBalance float64        `json:"initial_balance" gorm:"type:decimal(12,2);not null"`
	Currency       string         `json:"currency" gorm:"default:'USD'"`
	Status         GiftCardStatus `json:"status" gorm:"default:'active'"`
	ExpiresAt      *time.Time     `json:"expires_at"`
	IssuedAt       time.Time      `json:"issued_at" gorm:"autoCreateTime"`
	UsedAt         *time.Time     `json:"used_at"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Transactions []GiftCardTransaction `json:"transactions,omitempty" gorm:"foreignKey:GiftCardID"`
}

// GiftCardStatus represents the status of a gift card
type GiftCardStatus string

const (
	GiftCardStatusActive   GiftCardStatus = "active"
	GiftCardStatusInactive GiftCardStatus = "inactive"
	GiftCardStatusExpired  GiftCardStatus = "expired"
	GiftCardStatusUsed     GiftCardStatus = "used"
)

// GiftCardTransaction represents a gift card transaction
type GiftCardTransaction struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GiftCardID  uuid.UUID       `json:"gift_card_id" gorm:"type:uuid;not null;index"`
	OrderID     *uuid.UUID      `json:"order_id" gorm:"type:uuid;index"`
	Type        TransactionType `json:"type" gorm:"not null"`
	Amount      float64         `json:"amount" gorm:"type:decimal(12,2);not null"`
	Balance     float64         `json:"balance" gorm:"type:decimal(12,2);not null"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	GiftCard GiftCard `json:"gift_card,omitempty" gorm:"foreignKey:GiftCardID"`
}

// TransactionType represents the type of gift card transaction
type TransactionType string

const (
	TransactionTypeIssue  TransactionType = "issue"
	TransactionTypeRedeem TransactionType = "redeem"
	TransactionTypeRefund TransactionType = "refund"
	TransactionTypeAdjust TransactionType = "adjust"
)

// LoyaltyProgram represents a loyalty program
type LoyaltyProgram struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MerchantID  uuid.UUID       `json:"merchant_id" gorm:"type:uuid;not null;index"`
	Name        string          `json:"name" gorm:"not null"`
	Description string          `json:"description"`
	Status      ProgramStatus   `json:"status" gorm:"default:'active'"`
	PointValue  float64         `json:"point_value" gorm:"type:decimal(10,2);default:1.00"`
	RewardRatio float64         `json:"reward_ratio" gorm:"type:decimal(10,2);default:100.00"`
	Settings    LoyaltySettings `json:"settings" gorm:"type:jsonb"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Members []LoyaltyMember `json:"members,omitempty" gorm:"foreignKey:LoyaltyProgramID"`
}

// ProgramStatus represents the status of a loyalty program
type ProgramStatus string

const (
	ProgramStatusActive   ProgramStatus = "active"
	ProgramStatusInactive ProgramStatus = "inactive"
)

// LoyaltySettings represents settings for a loyalty program
type LoyaltySettings struct {
	EarnOnPurchase          bool     `json:"earn_on_purchase"`
	EarnOnReferral          bool     `json:"earn_on_referral"`
	EarnOnReview            bool     `json:"earn_on_review"`
	MinimumPurchaseAmount   float64  `json:"minimum_purchase_amount" gorm:"type:decimal(12,2)"`
	PointsExpirationDays    *int     `json:"points_expiration_days"`
	PointsRounding          Rounding `json:"points_rounding"`
	RedemptionEnabled       bool     `json:"redemption_enabled"`
	RedemptionRate          float64  `json:"redemption_rate" gorm:"type:decimal(10,2);default:1.00"`
	MinimumRedemptionPoints int      `json:"minimum_redemption_points" gorm:"default:100"`
}

// Rounding represents rounding options for loyalty points
type Rounding string

const (
	RoundingUp    Rounding = "up"
	RoundingDown  Rounding = "down"
	RoundingHalf  Rounding = "half"
	RoundingExact Rounding = "exact"
)

// LoyaltyMember represents a member of a loyalty program
type LoyaltyMember struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LoyaltyProgramID uuid.UUID    `json:"loyalty_program_id" gorm:"type:uuid;not null;index"`
	CustomerID       uuid.UUID    `json:"customer_id" gorm:"type:uuid;not null;index"`
	Points           int          `json:"points" gorm:"default:0"`
	LifetimePoints   int          `json:"lifetime_points" gorm:"default:0"`
	TierID           *uuid.UUID   `json:"tier_id" gorm:"type:uuid;index"`
	Status           MemberStatus `json:"status" gorm:"default:'active'"`
	EnrolledAt       time.Time    `json:"enrolled_at" gorm:"autoCreateTime"`
	LastActivityAt   *time.Time   `json:"last_activity_at"`
	CreatedAt        time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time    `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	LoyaltyProgram LoyaltyProgram    `json:"loyalty_program,omitempty" gorm:"foreignKey:LoyaltyProgramID"`
	Tier           *LoyaltyTier      `json:"tier,omitempty" gorm:"foreignKey:TierID"`
	Activities     []LoyaltyActivity `json:"activities,omitempty" gorm:"foreignKey:LoyaltyMemberID"`
}

// MemberStatus represents the status of a loyalty member
type MemberStatus string

const (
	MemberStatusActive    MemberStatus = "active"
	MemberStatusInactive  MemberStatus = "inactive"
	MemberStatusSuspended MemberStatus = "suspended"
)

// LoyaltyTier represents a tier in a loyalty program
type LoyaltyTier struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LoyaltyProgramID uuid.UUID `json:"loyalty_program_id" gorm:"type:uuid;not null;index"`
	Name             string    `json:"name" gorm:"not null"`
	Description      string    `json:"description"`
	MinimumPoints    int       `json:"minimum_points" gorm:"not null"`
	Benefits         []Benefit `json:"benefits" gorm:"type:jsonb"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	LoyaltyProgram LoyaltyProgram `json:"loyalty_program,omitempty" gorm:"foreignKey:LoyaltyProgramID"`
}

// Benefit represents a benefit in a loyalty tier
type Benefit struct {
	Type        BenefitType `json:"type"`
	Value       float64     `json:"value"`
	Description string      `json:"description"`
}

// BenefitType represents the type of benefit
type BenefitType string

const (
	BenefitTypeDiscount       BenefitType = "discount"
	BenefitTypeFreeShipping   BenefitType = "free_shipping"
	BenefitTypeBonusPoints    BenefitType = "bonus_points"
	BenefitTypeExclusiveItems BenefitType = "exclusive_items"
)

// LoyaltyActivity represents an activity in a loyalty program
type LoyaltyActivity struct {
	ID              uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LoyaltyMemberID uuid.UUID              `json:"loyalty_member_id" gorm:"type:uuid;not null;index"`
	Type            ActivityType           `json:"type" gorm:"not null"`
	Points          int                    `json:"points"`
	Description     string                 `json:"description"`
	ReferenceID     *uuid.UUID             `json:"reference_id" gorm:"type:uuid"`
	Metadata        map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	LoyaltyMember LoyaltyMember `json:"loyalty_member,omitempty" gorm:"foreignKey:LoyaltyMemberID"`
}

// ActivityType represents the type of loyalty activity
type ActivityType string

const (
	ActivityTypeEarned      ActivityType = "earned"
	ActivityTypeRedeemed    ActivityType = "redeemed"
	ActivityTypeExpired     ActivityType = "expired"
	ActivityTypeAdjusted    ActivityType = "adjusted"
	ActivityTypeTierChanged ActivityType = "tier_changed"
)

// BeforeCreate sets up UUID for new records
func (p *Promotion) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (dc *DiscountCode) BeforeCreate(tx *gorm.DB) error {
	if dc.ID == uuid.Nil {
		dc.ID = uuid.New()
	}
	return nil
}

func (cu *CodeUsage) BeforeCreate(tx *gorm.DB) error {
	if cu.ID == uuid.Nil {
		cu.ID = uuid.New()
	}
	return nil
}

func (gc *GiftCard) BeforeCreate(tx *gorm.DB) error {
	if gc.ID == uuid.Nil {
		gc.ID = uuid.New()
	}
	return nil
}

func (gct *GiftCardTransaction) BeforeCreate(tx *gorm.DB) error {
	if gct.ID == uuid.Nil {
		gct.ID = uuid.New()
	}
	return nil
}

func (lp *LoyaltyProgram) BeforeCreate(tx *gorm.DB) error {
	if lp.ID == uuid.Nil {
		lp.ID = uuid.New()
	}
	return nil
}

func (lm *LoyaltyMember) BeforeCreate(tx *gorm.DB) error {
	if lm.ID == uuid.Nil {
		lm.ID = uuid.New()
	}
	return nil
}

func (lt *LoyaltyTier) BeforeCreate(tx *gorm.DB) error {
	if lt.ID == uuid.Nil {
		lt.ID = uuid.New()
	}
	return nil
}

func (la *LoyaltyActivity) BeforeCreate(tx *gorm.DB) error {
	if la.ID == uuid.Nil {
		la.ID = uuid.New()
	}
	return nil
}
