package models

import (
	"time"

	"gorm.io/gorm"
)

// Merchant represents a merchant business account
type Merchant struct {
	ID             string           `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BusinessName   string           `json:"business_name" gorm:"not null"`
	LegalName      string           `json:"legal_name"`
	BusinessType   string           `json:"business_type"` // "sole_proprietorship", "partnership", "corporation", "llc"
	Industry       string           `json:"industry"`
	TaxID          string           `json:"tax_id"`
	WebsiteURL     string           `json:"website_url"`
	Description    string           `json:"description"`
	LogoURL        string           `json:"logo_url"`
	PrimaryEmail   string           `json:"primary_email" gorm:"not null"`
	PrimaryPhone   string           `json:"primary_phone"`
	Status         string           `json:"status" gorm:"default:'pending'"` // "pending", "active", "suspended", "closed"
	IsVerified     bool             `json:"is_verified" gorm:"default:false"`
	VerifiedAt     *time.Time       `json:"verified_at"`
	OnboardingStep int              `json:"onboarding_step" gorm:"default:1"`
	Settings       MerchantSettings `json:"settings" gorm:"type:jsonb"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `json:"-" gorm:"index"`

	// Relationships
	Addresses     []MerchantAddress `json:"addresses" gorm:"foreignKey:MerchantID"`
	BankAccounts  []BankAccount     `json:"bank_accounts" gorm:"foreignKey:MerchantID"`
	Subscriptions []Subscription    `json:"subscriptions" gorm:"foreignKey:MerchantID"`
	Stores        []Store           `json:"stores" gorm:"foreignKey:MerchantID"`
	Members       []MerchantMember  `json:"members" gorm:"foreignKey:MerchantID"`
}

// IsEntity implements the gqlgen federation Entity interface
func (m *Merchant) IsEntity() {}

// MerchantSettings contains merchant configuration settings
type MerchantSettings struct {
	Currency          string                 `json:"currency"`
	Timezone          string                 `json:"timezone"`
	DateFormat        string                 `json:"date_format"`
	WeightUnit        string                 `json:"weight_unit"`
	DimensionUnit     string                 `json:"dimension_unit"`
	OrderIDPrefix     string                 `json:"order_id_prefix"`
	NotificationPrefs map[string]bool        `json:"notification_preferences"`
	BusinessHours     map[string]string      `json:"business_hours"`
	ShippingZones     []string               `json:"shipping_zones"`
	TaxSettings       map[string]interface{} `json:"tax_settings"`
}

// MerchantAddress represents merchant business addresses
type MerchantAddress struct {
	ID         string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID string    `json:"merchant_id" gorm:"not null;index"`
	Type       string    `json:"type" gorm:"not null"` // "business", "billing", "shipping", "warehouse"
	Label      string    `json:"label"`
	Company    string    `json:"company"`
	Address1   string    `json:"address1" gorm:"not null"`
	Address2   string    `json:"address2"`
	City       string    `json:"city" gorm:"not null"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code" gorm:"not null"`
	Country    string    `json:"country" gorm:"not null"`
	Phone      string    `json:"phone"`
	IsDefault  bool      `json:"is_default" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	Merchant Merchant `json:"merchant" gorm:"foreignKey:MerchantID"`
}

// BankAccount represents merchant banking information
type BankAccount struct {
	ID                string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID        string    `json:"merchant_id" gorm:"not null;index"`
	AccountHolderName string    `json:"account_holder_name" gorm:"not null"`
	BankName          string    `json:"bank_name" gorm:"not null"`
	AccountType       string    `json:"account_type"` // "checking", "savings"
	RoutingNumber     string    `json:"routing_number" gorm:"not null"`
	AccountNumber     string    `json:"account_number" gorm:"not null"` // Should be encrypted
	Currency          string    `json:"currency" gorm:"default:'USD'"`
	IsVerified        bool      `json:"is_verified" gorm:"default:false"`
	IsDefault         bool      `json:"is_default" gorm:"default:false"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationships
	Merchant Merchant `json:"merchant" gorm:"foreignKey:MerchantID"`
}

// Subscription represents merchant subscription plans
type Subscription struct {
	ID                 string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID         string     `json:"merchant_id" gorm:"not null;index"`
	PlanID             string     `json:"plan_id" gorm:"not null"`
	Status             string     `json:"status" gorm:"not null"` // "active", "cancelled", "past_due", "unpaid"
	BillingCycle       string     `json:"billing_cycle"`          // "monthly", "annual"
	Amount             float64    `json:"amount" gorm:"not null"`
	Currency           string     `json:"currency" gorm:"default:'USD'"`
	TrialEndsAt        *time.Time `json:"trial_ends_at"`
	CurrentPeriodStart time.Time  `json:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end"`
	CancelledAt        *time.Time `json:"cancelled_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Relationships
	Merchant Merchant `json:"merchant" gorm:"foreignKey:MerchantID"`
	Plan     Plan     `json:"plan" gorm:"foreignKey:PlanID"`
}

// Plan represents subscription plans available to merchants
type Plan struct {
	ID           string                 `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string                 `json:"name" gorm:"not null;uniqueIndex"`
	DisplayName  string                 `json:"display_name" gorm:"not null"`
	Description  string                 `json:"description"`
	PlanType     string                 `json:"plan_type"` // "basic", "professional", "enterprise"
	MonthlyPrice float64                `json:"monthly_price"`
	AnnualPrice  float64                `json:"annual_price"`
	Currency     string                 `json:"currency" gorm:"default:'USD'"`
	Features     map[string]interface{} `json:"features" gorm:"type:jsonb"`
	Limits       map[string]int         `json:"limits" gorm:"type:jsonb"`
	IsActive     bool                   `json:"is_active" gorm:"default:true"`
	TrialDays    int                    `json:"trial_days" gorm:"default:14"`
	SetupFee     float64                `json:"setup_fee" gorm:"default:0"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`

	// Relationships
	Subscriptions []Subscription `json:"subscriptions" gorm:"foreignKey:PlanID"`
}

// Store represents physical or online store locations
type Store struct {
	ID         string                 `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID string                 `json:"merchant_id" gorm:"not null;index"`
	Name       string                 `json:"name" gorm:"not null"`
	StoreType  string                 `json:"store_type"` // "physical", "online", "popup", "warehouse"
	Address    string                 `json:"address"`
	City       string                 `json:"city"`
	State      string                 `json:"state"`
	PostalCode string                 `json:"postal_code"`
	Country    string                 `json:"country"`
	Phone      string                 `json:"phone"`
	Email      string                 `json:"email"`
	Website    string                 `json:"website"`
	Timezone   string                 `json:"timezone"`
	IsActive   bool                   `json:"is_active" gorm:"default:true"`
	Settings   map[string]interface{} `json:"settings" gorm:"type:jsonb"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`

	// Relationships
	Merchant Merchant `json:"merchant" gorm:"foreignKey:MerchantID"`
}

// MerchantMember represents users who have access to merchant accounts
type MerchantMember struct {
	ID          string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID  string     `json:"merchant_id" gorm:"not null;index"`
	UserID      string     `json:"user_id" gorm:"not null;index"`  // References Identity Service (no FK constraint)
	Role        string     `json:"role" gorm:"not null"`           // "owner", "admin", "manager", "staff", "viewer"
	Status      string     `json:"status" gorm:"default:'active'"` // "active", "invited", "suspended"
	Permissions []string   `json:"permissions" gorm:"type:jsonb"`
	InvitedBy   string     `json:"invited_by"` // UserID who sent invitation (no FK constraint)
	InvitedAt   *time.Time `json:"invited_at"`
	JoinedAt    *time.Time `json:"joined_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships (internal to this service only)
	Merchant Merchant `json:"merchant" gorm:"foreignKey:MerchantID"`
	// Note: No User relationship as it's in a different service/database
}

// Invoice represents billing invoices for merchants
type Invoice struct {
	ID             string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID     string     `json:"merchant_id" gorm:"not null;index"`
	SubscriptionID string     `json:"subscription_id" gorm:"index"`
	InvoiceNumber  string     `json:"invoice_number" gorm:"uniqueIndex;not null"`
	Status         string     `json:"status"` // "draft", "sent", "paid", "overdue", "cancelled"
	Currency       string     `json:"currency" gorm:"default:'USD'"`
	Subtotal       float64    `json:"subtotal"`
	TaxAmount      float64    `json:"tax_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	TotalAmount    float64    `json:"total_amount"`
	AmountPaid     float64    `json:"amount_paid" gorm:"default:0"`
	AmountDue      float64    `json:"amount_due"`
	DueDate        time.Time  `json:"due_date"`
	PaidAt         *time.Time `json:"paid_at"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	Merchant     Merchant          `json:"merchant" gorm:"foreignKey:MerchantID"`
	Subscription *Subscription     `json:"subscription" gorm:"foreignKey:SubscriptionID"`
	LineItems    []InvoiceLineItem `json:"line_items" gorm:"foreignKey:InvoiceID"`
}

// InvoiceLineItem represents individual items on an invoice
type InvoiceLineItem struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InvoiceID   string    `json:"invoice_id" gorm:"not null;index"`
	Description string    `json:"description" gorm:"not null"`
	Quantity    float64   `json:"quantity" gorm:"default:1"`
	UnitPrice   float64   `json:"unit_price"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Invoice Invoice `json:"invoice" gorm:"foreignKey:InvoiceID"`
}

// MerchantVerification represents the verification process for merchants
type MerchantVerification struct {
	ID               string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MerchantID       string     `json:"merchant_id" gorm:"not null;index"`
	VerificationType string     `json:"verification_type"` // "identity", "business", "bank_account"
	Status           string     `json:"status"`            // "pending", "approved", "rejected", "requires_action"
	DocumentType     string     `json:"document_type"`     // "drivers_license", "passport", "business_license"
	DocumentURL      string     `json:"document_url"`
	RejectionReason  string     `json:"rejection_reason"`
	VerifiedBy       string     `json:"verified_by"` // Admin user ID
	VerifiedAt       *time.Time `json:"verified_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relationships
	Merchant Merchant `json:"merchant" gorm:"foreignKey:MerchantID"`
}

// BeforeCreate hook for Merchant model
func (m *Merchant) BeforeCreate(tx *gorm.DB) error {
	if m.Settings.Currency == "" {
		m.Settings.Currency = "USD"
	}
	if m.Settings.Timezone == "" {
		m.Settings.Timezone = "UTC"
	}
	if m.Settings.WeightUnit == "" {
		m.Settings.WeightUnit = "kg"
	}
	if m.Settings.DimensionUnit == "" {
		m.Settings.DimensionUnit = "cm"
	}
	return nil
}

// IsOwner checks if a user is the owner of the merchant account
func (m *Merchant) IsOwner(userID string) bool {
	for _, member := range m.Members {
		if member.UserID == userID && member.Role == "owner" {
			return true
		}
	}
	return false
}

// HasMember checks if a user is a member of the merchant account
func (m *Merchant) HasMember(userID string) bool {
	for _, member := range m.Members {
		if member.UserID == userID && member.Status == "active" {
			return true
		}
	}
	return false
}

// GetMemberRole returns the role of a specific user in the merchant account
func (m *Merchant) GetMemberRole(userID string) string {
	for _, member := range m.Members {
		if member.UserID == userID && member.Status == "active" {
			return member.Role
		}
	}
	return ""
}

// IsActive checks if the merchant account is active
func (m *Merchant) IsActive() bool {
	return m.Status == "active"
}

// GetActiveSubscription returns the active subscription for the merchant
func (m *Merchant) GetActiveSubscription() *Subscription {
	for _, sub := range m.Subscriptions {
		if sub.Status == "active" {
			return &sub
		}
	}
	return nil
}
