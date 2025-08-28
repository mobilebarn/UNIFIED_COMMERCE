package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment represents a payment transaction
type Payment struct {
	ID               uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID          uuid.UUID              `json:"order_id" gorm:"type:uuid;not null;index"`
	MerchantID       uuid.UUID              `json:"merchant_id" gorm:"type:uuid;not null;index"`
	CustomerID       *uuid.UUID             `json:"customer_id" gorm:"type:uuid;index"`
	PaymentMethodID  uuid.UUID              `json:"payment_method_id" gorm:"type:uuid;not null;index"`
	GatewayID        uuid.UUID              `json:"gateway_id" gorm:"type:uuid;not null;index"`
	Status           PaymentStatus          `json:"status" gorm:"default:'pending'"`
	Amount           float64                `json:"amount" gorm:"type:decimal(12,2);not null"`
	Currency         string                 `json:"currency" gorm:"default:'USD'"`
	GatewayReference string                 `json:"gateway_reference"`
	AuthorizationID  string                 `json:"authorization_id"`
	TransactionFee   float64                `json:"transaction_fee" gorm:"type:decimal(10,2);default:0"`
	NetAmount        float64                `json:"net_amount" gorm:"type:decimal(12,2);default:0"`
	Description      string                 `json:"description"`
	Metadata         map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	ProcessedAt      *time.Time             `json:"processed_at"`
	CompletedAt      *time.Time             `json:"completed_at"`
	FailedAt         *time.Time             `json:"failed_at"`
	RefundedAt       *time.Time             `json:"refunded_at"`
	CreatedAt        time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time              `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order         interface{}    `json:"order,omitempty" gorm:"-"` // Not stored in DB, populated by service
	PaymentMethod PaymentMethod  `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID"`
	Gateway       PaymentGateway `json:"gateway,omitempty" gorm:"foreignKey:GatewayID"`
	Refunds       []Refund       `json:"refunds,omitempty" gorm:"foreignKey:PaymentID"`
	Events        []PaymentEvent `json:"events,omitempty" gorm:"foreignKey:PaymentID"`
}

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "pending"
	PaymentStatusProcessing        PaymentStatus = "processing"
	PaymentStatusAuthorized        PaymentStatus = "authorized"
	PaymentStatusCaptured          PaymentStatus = "captured"
	PaymentStatusFailed            PaymentStatus = "failed"
	PaymentStatusCancelled         PaymentStatus = "cancelled"
	PaymentStatusRefunded          PaymentStatus = "refunded"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
)

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CustomerID  *uuid.UUID             `json:"customer_id" gorm:"type:uuid;index"`
	MerchantID  *uuid.UUID             `json:"merchant_id" gorm:"type:uuid;index"`
	Type        PaymentMethodType      `json:"type" gorm:"not null"`
	Provider    string                 `json:"provider"`
	Token       string                 `json:"token"`
	Last4       string                 `json:"last4"`
	ExpiryMonth int                    `json:"expiry_month"`
	ExpiryYear  int                    `json:"expiry_year"`
	Brand       string                 `json:"brand"`
	Name        string                 `json:"name"`
	Email       string                 `json:"email"`
	IsDefault   bool                   `json:"is_default" gorm:"default:false"`
	IsActive    bool                   `json:"is_active" gorm:"default:true"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// PaymentMethodType represents the type of payment method
type PaymentMethodType string

const (
	PaymentMethodTypeCreditCard   PaymentMethodType = "credit_card"
	PaymentMethodTypeDebitCard    PaymentMethodType = "debit_card"
	PaymentMethodTypePayPal       PaymentMethodType = "paypal"
	PaymentMethodTypeApplePay     PaymentMethodType = "apple_pay"
	PaymentMethodTypeGooglePay    PaymentMethodType = "google_pay"
	PaymentMethodTypeBankTransfer PaymentMethodType = "bank_transfer"
	PaymentMethodTypeCrypto       PaymentMethodType = "crypto"
	PaymentMethodTypeWallet       PaymentMethodType = "wallet"
)

// PaymentGateway represents a payment gateway
type PaymentGateway struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string                 `json:"name" gorm:"not null;unique"`
	Provider    string                 `json:"provider" gorm:"not null"`
	IsEnabled   bool                   `json:"is_enabled" gorm:"default:true"`
	IsSandbox   bool                   `json:"is_sandbox" gorm:"default:false"`
	Credentials map[string]interface{} `json:"credentials" gorm:"type:jsonb"`
	Settings    map[string]interface{} `json:"settings" gorm:"type:jsonb"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// Refund represents a refund transaction
type Refund struct {
	ID               uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PaymentID        uuid.UUID    `json:"payment_id" gorm:"type:uuid;not null;index"`
	OrderID          uuid.UUID    `json:"order_id" gorm:"type:uuid;not null;index"`
	Amount           float64      `json:"amount" gorm:"type:decimal(12,2);not null"`
	Currency         string       `json:"currency" gorm:"default:'USD'"`
	Reason           RefundReason `json:"reason"`
	Status           RefundStatus `json:"status" gorm:"default:'pending'"`
	GatewayReference string       `json:"gateway_reference"`
	ProcessedAt      *time.Time   `json:"processed_at"`
	CompletedAt      *time.Time   `json:"completed_at"`
	FailedAt         *time.Time   `json:"failed_at"`
	CreatedAt        time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time    `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Payment Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
}

// RefundStatus represents the status of a refund
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusProcessed RefundStatus = "processed"
	RefundStatusCompleted RefundStatus = "completed"
	RefundStatusFailed    RefundStatus = "failed"
	RefundStatusCancelled RefundStatus = "cancelled"
)

// RefundReason represents the reason for a refund
type RefundReason string

const (
	RefundReasonCustomerRequest       RefundReason = "customer_request"
	RefundReasonProductNotAsDescribed RefundReason = "product_not_as_described"
	RefundReasonProductDamaged        RefundReason = "product_damaged"
	RefundReasonProductDefective      RefundReason = "product_defective"
	RefundReasonDuplicateOrder        RefundReason = "duplicate_order"
	RefundReasonBillingError          RefundReason = "billing_error"
	RefundReasonOther                 RefundReason = "other"
)

// PaymentEvent represents events in the payment lifecycle
type PaymentEvent struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PaymentID   uuid.UUID              `json:"payment_id" gorm:"type:uuid;not null;index"`
	EventType   PaymentEventType       `json:"event_type" gorm:"not null"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Payment Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
}

// PaymentEventType represents the type of payment event
type PaymentEventType string

const (
	PaymentEventCreated           PaymentEventType = "created"
	PaymentEventProcessing        PaymentEventType = "processing"
	PaymentEventAuthorized        PaymentEventType = "authorized"
	PaymentEventCaptured          PaymentEventType = "captured"
	PaymentEventFailed            PaymentEventType = "failed"
	PaymentEventCancelled         PaymentEventType = "cancelled"
	PaymentEventRefunded          PaymentEventType = "refunded"
	PaymentEventPartiallyRefunded PaymentEventType = "partially_refunded"
	PaymentEventUpdated           PaymentEventType = "updated"
)

// Settlement represents a payment settlement
type Settlement struct {
	ID          uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GatewayID   uuid.UUID        `json:"gateway_id" gorm:"type:uuid;not null;index"`
	Reference   string           `json:"reference" gorm:"unique;not null"`
	Status      SettlementStatus `json:"status" gorm:"default:'pending'"`
	Amount      float64          `json:"amount" gorm:"type:decimal(12,2);not null"`
	Currency    string           `json:"currency" gorm:"default:'USD'"`
	DepositedAt *time.Time       `json:"deposited_at"`
	ProcessedAt *time.Time       `json:"processed_at"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Gateway PaymentGateway `json:"gateway,omitempty" gorm:"foreignKey:GatewayID"`
}

// SettlementStatus represents the status of a settlement
type SettlementStatus string

const (
	SettlementStatusPending   SettlementStatus = "pending"
	SettlementStatusProcessed SettlementStatus = "processed"
	SettlementStatusCompleted SettlementStatus = "completed"
	SettlementStatusFailed    SettlementStatus = "failed"
)

// BeforeCreate sets up UUID for new records
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) error {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return nil
}

func (pg *PaymentGateway) BeforeCreate(tx *gorm.DB) error {
	if pg.ID == uuid.Nil {
		pg.ID = uuid.New()
	}
	return nil
}

func (r *Refund) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (pe *PaymentEvent) BeforeCreate(tx *gorm.DB) error {
	if pe.ID == uuid.Nil {
		pe.ID = uuid.New()
	}
	return nil
}

func (s *Settlement) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
