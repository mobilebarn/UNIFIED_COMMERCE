package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Cart represents a shopping cart
type Cart struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SessionID  string     `json:"session_id" gorm:"index"`
	CustomerID *uuid.UUID `json:"customer_id" gorm:"type:uuid;index"`
	MerchantID uuid.UUID  `json:"merchant_id" gorm:"type:uuid;not null;index"`
	Status     CartStatus `json:"status" gorm:"default:'active'"`
	Currency   string     `json:"currency" gorm:"default:'USD'"`

	// Customer Information (for guest checkouts)
	CustomerEmail     string `json:"customer_email"`
	CustomerPhone     string `json:"customer_phone"`
	CustomerFirstName string `json:"customer_first_name"`
	CustomerLastName  string `json:"customer_last_name"`

	// Addresses
	BillingAddress  Address `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
	ShippingAddress Address `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`

	// Pricing Information
	SubtotalPrice float64 `json:"subtotal_price" gorm:"type:decimal(12,2);default:0"`
	TotalTax      float64 `json:"total_tax" gorm:"type:decimal(12,2);default:0"`
	TotalShipping float64 `json:"total_shipping" gorm:"type:decimal(12,2);default:0"`
	TotalDiscount float64 `json:"total_discount" gorm:"type:decimal(12,2);default:0"`
	TotalPrice    float64 `json:"total_price" gorm:"type:decimal(12,2);default:0"`

	// Checkout Information
	CheckoutStep     CheckoutStep `json:"checkout_step" gorm:"default:'cart'"`
	PaymentMethodID  string       `json:"payment_method_id"`
	ShippingMethodID string       `json:"shipping_method_id"`

	// Marketing
	DiscountCodes []string   `json:"discount_codes" gorm:"type:text[]"`
	AbandonedAt   *time.Time `json:"abandoned_at"`
	RecoveredAt   *time.Time `json:"recovered_at"`

	// Metadata
	Notes      string                 `json:"notes"`
	Attributes map[string]interface{} `json:"attributes" gorm:"type:jsonb"`

	// Timestamps
	ExpiresAt   *time.Time `json:"expires_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	LineItems            []CartLineItem            `json:"line_items,omitempty" gorm:"foreignKey:CartID"`
	TaxLines             []CartTaxLine             `json:"tax_lines,omitempty" gorm:"foreignKey:CartID"`
	ShippingLines        []CartShippingLine        `json:"shipping_lines,omitempty" gorm:"foreignKey:CartID"`
	DiscountApplications []CartDiscountApplication `json:"discount_applications,omitempty" gorm:"foreignKey:CartID"`
}

// IsEntity marks Cart as a federation entity
func (c Cart) IsEntity() {}

// CartStatus represents the status of a cart
type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusAbandoned CartStatus = "abandoned"
	CartStatusCompleted CartStatus = "completed"
	CartStatusExpired   CartStatus = "expired"
)

// CheckoutStep represents the current step in the checkout process
type CheckoutStep string

const (
	CheckoutStepCart     CheckoutStep = "cart"
	CheckoutStepShipping CheckoutStep = "shipping"
	CheckoutStepPayment  CheckoutStep = "payment"
	CheckoutStepReview   CheckoutStep = "review"
	CheckoutStepComplete CheckoutStep = "complete"
)

// Address represents a billing or shipping address
type Address struct {
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Company    string  `json:"company"`
	Street1    string  `json:"street1"`
	Street2    string  `json:"street2"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	PostalCode string  `json:"postal_code"`
	Phone      string  `json:"phone"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

// IsEntity marks Address as a federation entity
func (a Address) IsEntity() {}

// CartLineItem represents an item in a shopping cart
type CartLineItem struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CartID           uuid.UUID  `json:"cart_id" gorm:"type:uuid;not null;index"`
	ProductID        uuid.UUID  `json:"product_id" gorm:"type:uuid;not null;index"`
	ProductVariantID *uuid.UUID `json:"product_variant_id" gorm:"type:uuid;index"`

	// Product Information (snapshot at time of adding to cart)
	Name         string `json:"name" gorm:"not null"`
	SKU          string `json:"sku" gorm:"not null;index"`
	Barcode      string `json:"barcode"`
	ProductTitle string `json:"product_title"`
	VariantTitle string `json:"variant_title"`
	Vendor       string `json:"vendor"`
	ProductImage string `json:"product_image"`

	// Quantity and Pricing
	Quantity       int     `json:"quantity" gorm:"not null"`
	Price          float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	CompareAtPrice float64 `json:"compare_at_price" gorm:"type:decimal(10,2)"`
	LinePrice      float64 `json:"line_price" gorm:"type:decimal(12,2)"`

	// Discounts and Tax
	TotalDiscount float64 `json:"total_discount" gorm:"type:decimal(10,2);default:0"`
	Taxable       bool    `json:"taxable" gorm:"default:true"`

	// Fulfillment
	RequiresShipping bool `json:"requires_shipping" gorm:"default:true"`
	IsGiftCard       bool `json:"is_gift_card" gorm:"default:false"`

	// Metadata
	Properties map[string]string `json:"properties,omitempty" gorm:"type:jsonb"`
	CreatedAt  time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Cart                Cart                             `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	DiscountAllocations []CartLineItemDiscountAllocation `json:"discount_allocations,omitempty" gorm:"foreignKey:LineItemID"`
}

// CartTaxLine represents tax information for a cart
type CartTaxLine struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CartID    uuid.UUID `json:"cart_id" gorm:"type:uuid;not null;index"`
	Title     string    `json:"title" gorm:"not null"`
	Rate      float64   `json:"rate" gorm:"type:decimal(8,4);not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Cart Cart `json:"cart,omitempty" gorm:"foreignKey:CartID"`
}

// CartShippingLine represents shipping information for a cart
type CartShippingLine struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CartID           uuid.UUID `json:"cart_id" gorm:"type:uuid;not null;index"`
	Title            string    `json:"title" gorm:"not null"`
	Code             string    `json:"code"`
	Source           string    `json:"source"`
	Price            float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	DiscountedPrice  float64   `json:"discounted_price" gorm:"type:decimal(10,2)"`
	CarrierService   string    `json:"carrier_service"`
	DeliveryCategory string    `json:"delivery_category"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Cart Cart `json:"cart,omitempty" gorm:"foreignKey:CartID"`
}

// CartDiscountApplication represents a discount applied to a cart
type CartDiscountApplication struct {
	ID               uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CartID           uuid.UUID         `json:"cart_id" gorm:"type:uuid;not null;index"`
	Type             DiscountType      `json:"type" gorm:"not null"`
	Code             string            `json:"code"`
	Title            string            `json:"title" gorm:"not null"`
	Description      string            `json:"description"`
	Value            float64           `json:"value" gorm:"type:decimal(10,2);not null"`
	ValueType        DiscountValueType `json:"value_type" gorm:"not null"`
	AllocationMethod AllocationMethod  `json:"allocation_method" gorm:"not null"`
	TargetSelection  TargetSelection   `json:"target_selection" gorm:"not null"`
	TargetType       TargetType        `json:"target_type" gorm:"not null"`
	CreatedAt        time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time         `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Cart Cart `json:"cart,omitempty" gorm:"foreignKey:CartID"`
}

// DiscountType represents the type of discount
type DiscountType string

const (
	DiscountTypeManual    DiscountType = "manual"
	DiscountTypeCode      DiscountType = "code"
	DiscountTypeAutomatic DiscountType = "automatic"
)

// DiscountValueType represents how the discount value is applied
type DiscountValueType string

const (
	DiscountValueTypeFixed      DiscountValueType = "fixed"
	DiscountValueTypePercentage DiscountValueType = "percentage"
)

// AllocationMethod represents how the discount is allocated
type AllocationMethod string

const (
	AllocationMethodAcross AllocationMethod = "across"
	AllocationMethodEach   AllocationMethod = "each"
)

// TargetSelection represents what the discount targets
type TargetSelection string

const (
	TargetSelectionAll      TargetSelection = "all"
	TargetSelectionEntitled TargetSelection = "entitled"
	TargetSelectionExplicit TargetSelection = "explicit"
)

// TargetType represents the type of target for discounts
type TargetType string

const (
	TargetTypeLineItem     TargetType = "line_item"
	TargetTypeShippingLine TargetType = "shipping_line"
)

// CartLineItemDiscountAllocation represents discount allocation to a line item
type CartLineItemDiscountAllocation struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LineItemID            uuid.UUID `json:"line_item_id" gorm:"type:uuid;not null;index"`
	DiscountApplicationID uuid.UUID `json:"discount_application_id" gorm:"type:uuid;not null;index"`
	Amount                float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	LineItem            CartLineItem            `json:"line_item,omitempty" gorm:"foreignKey:LineItemID"`
	DiscountApplication CartDiscountApplication `json:"discount_application,omitempty" gorm:"foreignKey:DiscountApplicationID"`
}

// Checkout represents a checkout session
type Checkout struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CartID        uuid.UUID      `json:"cart_id" gorm:"type:uuid;not null;index"`
	CheckoutToken string         `json:"checkout_token" gorm:"unique;not null"`
	Status        CheckoutStatus `json:"status" gorm:"default:'pending'"`

	// Customer Information
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	AcceptsMarketing bool   `json:"accepts_marketing" gorm:"default:false"`

	// Checkout Flow
	CompletedStep    CheckoutStep `json:"completed_step" gorm:"default:'cart'"`
	RequiresShipping bool         `json:"requires_shipping" gorm:"default:true"`

	// Payment Information
	PaymentGateway   string `json:"payment_gateway"`
	PaymentMethodID  string `json:"payment_method_id"`
	PaymentSessionID string `json:"payment_session_id"`

	// Timestamps
	AbandonedAt *time.Time `json:"abandoned_at"`
	CompletedAt *time.Time `json:"completed_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Cart   Cart            `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	Events []CheckoutEvent `json:"events,omitempty" gorm:"foreignKey:CheckoutID"`
}

// CheckoutStatus represents the status of a checkout
type CheckoutStatus string

const (
	CheckoutStatusPending   CheckoutStatus = "pending"
	CheckoutStatusCompleted CheckoutStatus = "completed"
	CheckoutStatusAbandoned CheckoutStatus = "abandoned"
	CheckoutStatusExpired   CheckoutStatus = "expired"
)

// CheckoutEvent represents events in the checkout process
type CheckoutEvent struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CheckoutID  uuid.UUID              `json:"checkout_id" gorm:"type:uuid;not null;index"`
	EventType   CheckoutEventType      `json:"event_type" gorm:"not null"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Checkout Checkout `json:"checkout,omitempty" gorm:"foreignKey:CheckoutID"`
}

// CheckoutEventType represents the type of checkout event
type CheckoutEventType string

const (
	CheckoutEventCreated           CheckoutEventType = "created"
	CheckoutEventCustomerInfoAdded CheckoutEventType = "customer_info_added"
	CheckoutEventShippingAdded     CheckoutEventType = "shipping_added"
	CheckoutEventPaymentAdded      CheckoutEventType = "payment_added"
	CheckoutEventDiscountApplied   CheckoutEventType = "discount_applied"
	CheckoutEventDiscountRemoved   CheckoutEventType = "discount_removed"
	CheckoutEventCompleted         CheckoutEventType = "completed"
	CheckoutEventAbandoned         CheckoutEventType = "abandoned"
	CheckoutEventRecovered         CheckoutEventType = "recovered"
)

// ShippingRate represents available shipping rates
type ShippingRate struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title          string    `json:"title" gorm:"not null"`
	Code           string    `json:"code" gorm:"not null"`
	Price          float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Source         string    `json:"source"`
	CarrierService string    `json:"carrier_service"`
	ServiceCode    string    `json:"service_code"`
	DeliveryRange  string    `json:"delivery_range"`
	DeliveryDays   int       `json:"delivery_days"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// PaymentMethod represents available payment methods
type PaymentMethod struct {
	ID            uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Gateway       string                 `json:"gateway" gorm:"not null"`
	Type          PaymentType            `json:"type" gorm:"not null"`
	Title         string                 `json:"title" gorm:"not null"`
	Description   string                 `json:"description"`
	Enabled       bool                   `json:"enabled" gorm:"default:true"`
	Configuration map[string]interface{} `json:"configuration" gorm:"type:jsonb"`
	CreatedAt     time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// PaymentType represents the type of payment method
type PaymentType string

const (
	PaymentTypeCreditCard   PaymentType = "credit_card"
	PaymentTypeDebitCard    PaymentType = "debit_card"
	PaymentTypePayPal       PaymentType = "paypal"
	PaymentTypeApplePay     PaymentType = "apple_pay"
	PaymentTypeGooglePay    PaymentType = "google_pay"
	PaymentTypeBankTransfer PaymentType = "bank_transfer"
	PaymentTypeCrypto       PaymentType = "crypto"
	PaymentTypeWallet       PaymentType = "wallet"
)

// BeforeCreate sets up UUID for new records
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (cli *CartLineItem) BeforeCreate(tx *gorm.DB) error {
	if cli.ID == uuid.Nil {
		cli.ID = uuid.New()
	}
	return nil
}

func (ctl *CartTaxLine) BeforeCreate(tx *gorm.DB) error {
	if ctl.ID == uuid.Nil {
		ctl.ID = uuid.New()
	}
	return nil
}

func (csl *CartShippingLine) BeforeCreate(tx *gorm.DB) error {
	if csl.ID == uuid.Nil {
		csl.ID = uuid.New()
	}
	return nil
}

func (cda *CartDiscountApplication) BeforeCreate(tx *gorm.DB) error {
	if cda.ID == uuid.Nil {
		cda.ID = uuid.New()
	}
	return nil
}

func (clida *CartLineItemDiscountAllocation) BeforeCreate(tx *gorm.DB) error {
	if clida.ID == uuid.Nil {
		clida.ID = uuid.New()
	}
	return nil
}

func (ch *Checkout) BeforeCreate(tx *gorm.DB) error {
	if ch.ID == uuid.Nil {
		ch.ID = uuid.New()
	}
	return nil
}

func (ce *CheckoutEvent) BeforeCreate(tx *gorm.DB) error {
	if ce.ID == uuid.Nil {
		ce.ID = uuid.New()
	}
	return nil
}

func (sr *ShippingRate) BeforeCreate(tx *gorm.DB) error {
	if sr.ID == uuid.Nil {
		sr.ID = uuid.New()
	}
	return nil
}

func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}

// AfterSave recalculates cart totals after line item changes
func (cli *CartLineItem) AfterSave(tx *gorm.DB) error {
	cli.LinePrice = float64(cli.Quantity) * cli.Price
	return nil
}
