package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order represents a customer order
type Order struct {
	ID                uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderNumber       string            `json:"order_number" gorm:"unique;not null;index"`
	MerchantID        uuid.UUID         `json:"merchant_id" gorm:"type:uuid;not null;index"`
	CustomerID        *uuid.UUID        `json:"customer_id" gorm:"type:uuid;index"`
	LocationID        *uuid.UUID        `json:"location_id" gorm:"type:uuid;index"`
	Status            OrderStatus       `json:"status" gorm:"default:'pending'"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status" gorm:"default:'unfulfilled'"`
	PaymentStatus     PaymentStatus     `json:"payment_status" gorm:"default:'pending'"`

	// Customer Information
	Customer CustomerInfo `json:"customer" gorm:"embedded;embeddedPrefix:customer_"`

	// Addresses
	BillingAddress  Address `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
	ShippingAddress Address `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`

	// Financial Information
	SubtotalPrice float64 `json:"subtotal_price" gorm:"type:decimal(12,2);default:0"`
	TotalTax      float64 `json:"total_tax" gorm:"type:decimal(12,2);default:0"`
	TotalShipping float64 `json:"total_shipping" gorm:"type:decimal(12,2);default:0"`
	TotalDiscount float64 `json:"total_discount" gorm:"type:decimal(12,2);default:0"`
	TotalPrice    float64 `json:"total_price" gorm:"type:decimal(12,2);default:0"`

	// Shipping Information
	ShippingMethod string  `json:"shipping_method"`
	ShippingRate   float64 `json:"shipping_rate" gorm:"type:decimal(10,2)"`
	TrackingNumber string  `json:"tracking_number"`
	TrackingURL    string  `json:"tracking_url"`
	Carrier        string  `json:"carrier"`

	// Order Metadata
	Source        OrderSource `json:"source" gorm:"default:'online'"`
	Channel       string      `json:"channel"`
	Currency      string      `json:"currency" gorm:"default:'USD'"`
	Tags          []string    `json:"tags" gorm:"type:text[]"`
	Notes         string      `json:"notes"`
	InternalNotes string      `json:"internal_notes"`

	// Timestamps
	ProcessedAt *time.Time `json:"processed_at"`
	FulfilledAt *time.Time `json:"fulfilled_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	LineItems    []OrderLineItem `json:"line_items,omitempty" gorm:"foreignKey:OrderID"`
	Fulfillments []Fulfillment   `json:"fulfillments,omitempty" gorm:"foreignKey:OrderID"`
	Transactions []Transaction   `json:"transactions,omitempty" gorm:"foreignKey:OrderID"`
	Returns      []Return        `json:"returns,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderStatus represents the overall status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusReturned   OrderStatus = "returned"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// FulfillmentStatus represents the fulfillment status of an order
type FulfillmentStatus string

const (
	FulfillmentStatusUnfulfilled        FulfillmentStatus = "unfulfilled"
	FulfillmentStatusPartiallyFulfilled FulfillmentStatus = "partially_fulfilled"
	FulfillmentStatusFulfilled          FulfillmentStatus = "fulfilled"
	FulfillmentStatusRestocked          FulfillmentStatus = "restocked"
)

// PaymentStatus represents the payment status of an order
type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "pending"
	PaymentStatusAuthorized        PaymentStatus = "authorized"
	PaymentStatusPartiallyPaid     PaymentStatus = "partially_paid"
	PaymentStatusPaid              PaymentStatus = "paid"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
	PaymentStatusRefunded          PaymentStatus = "refunded"
	PaymentStatusVoided            PaymentStatus = "voided"
)

// OrderSource represents the source/channel of the order
type OrderSource string

const (
	OrderSourceOnline OrderSource = "online"
	OrderSourcePOS    OrderSource = "pos"
	OrderSourcePhone  OrderSource = "phone"
	OrderSourceEmail  OrderSource = "email"
	OrderSourceAPI    OrderSource = "api"
)

// CustomerInfo represents customer information for an order
type CustomerInfo struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Address represents a billing or shipping address
type Address struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Company   string  `json:"company"`
	Address1  string  `json:"address1"`
	Address2  string  `json:"address2"`
	City      string  `json:"city"`
	Province  string  `json:"province"`
	Country   string  `json:"country"`
	Zip       string  `json:"zip"`
	Phone     string  `json:"phone"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// OrderLineItem represents an item in an order
type OrderLineItem struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID          uuid.UUID  `json:"order_id" gorm:"type:uuid;not null;index"`
	ProductID        uuid.UUID  `json:"product_id" gorm:"type:uuid;not null;index"`
	ProductVariantID *uuid.UUID `json:"product_variant_id" gorm:"type:uuid;index"`

	// Product Information (snapshot at time of order)
	Name         string `json:"name" gorm:"not null"`
	SKU          string `json:"sku" gorm:"not null;index"`
	Barcode      string `json:"barcode"`
	ProductTitle string `json:"product_title"`
	VariantTitle string `json:"variant_title"`
	Vendor       string `json:"vendor"`

	// Quantity and Pricing
	Quantity       int     `json:"quantity" gorm:"not null"`
	Price          float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	CompareAtPrice float64 `json:"compare_at_price" gorm:"type:decimal(10,2)"`
	LinePrice      float64 `json:"line_price" gorm:"type:decimal(12,2)"`
	TotalDiscount  float64 `json:"total_discount" gorm:"type:decimal(10,2);default:0"`

	// Tax Information
	Taxable  bool      `json:"taxable" gorm:"default:true"`
	TaxLines []TaxLine `json:"tax_lines,omitempty" gorm:"type:jsonb"`

	// Fulfillment
	FulfillmentStatus FulfillmentStatus `json:"fulfillment_status" gorm:"default:'unfulfilled'"`
	FulfilledQuantity int               `json:"fulfilled_quantity" gorm:"default:0"`

	// Metadata
	Properties       map[string]string `json:"properties,omitempty" gorm:"type:jsonb"`
	RequiresShipping bool              `json:"requires_shipping" gorm:"default:true"`
	IsGiftCard       bool              `json:"is_gift_card" gorm:"default:false"`
	CreatedAt        time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time         `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// TaxLine represents tax information for a line item
type TaxLine struct {
	Title string  `json:"title"`
	Rate  float64 `json:"rate"`
	Price float64 `json:"price"`
}

// Fulfillment represents a shipment of order items
type Fulfillment struct {
	ID         uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID    uuid.UUID         `json:"order_id" gorm:"type:uuid;not null;index"`
	LocationID *uuid.UUID        `json:"location_id" gorm:"type:uuid;index"`
	Status     FulfillmentStatus `json:"status" gorm:"default:'pending'"`

	// Shipping Information
	TrackingNumber  string         `json:"tracking_number"`
	TrackingURL     string         `json:"tracking_url"`
	TrackingCompany string         `json:"tracking_company"`
	ShipmentStatus  ShipmentStatus `json:"shipment_status" gorm:"default:'pending'"`

	// Service Information
	Service string `json:"service"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	ShippedAt   *time.Time `json:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at"`

	// Relationships
	Order     Order                 `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	LineItems []FulfillmentLineItem `json:"line_items,omitempty" gorm:"foreignKey:FulfillmentID"`
}

// ShipmentStatus represents the status of a shipment
type ShipmentStatus string

const (
	ShipmentStatusPending        ShipmentStatus = "pending"
	ShipmentStatusConfirmed      ShipmentStatus = "confirmed"
	ShipmentStatusInTransit      ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered      ShipmentStatus = "delivered"
	ShipmentStatusException      ShipmentStatus = "exception"
	ShipmentStatusFailure        ShipmentStatus = "failure"
	ShipmentStatusCancelled      ShipmentStatus = "cancelled"
)

// FulfillmentLineItem represents items in a fulfillment
type FulfillmentLineItem struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FulfillmentID uuid.UUID `json:"fulfillment_id" gorm:"type:uuid;not null;index"`
	LineItemID    uuid.UUID `json:"line_item_id" gorm:"type:uuid;not null;index"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Fulfillment Fulfillment   `json:"fulfillment,omitempty" gorm:"foreignKey:FulfillmentID"`
	LineItem    OrderLineItem `json:"line_item,omitempty" gorm:"foreignKey:LineItemID"`
}

// Transaction represents a payment transaction
type Transaction struct {
	ID      uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID uuid.UUID         `json:"order_id" gorm:"type:uuid;not null;index"`
	Kind    TransactionKind   `json:"kind" gorm:"not null"`
	Gateway string            `json:"gateway" gorm:"not null"`
	Status  TransactionStatus `json:"status" gorm:"not null"`
	Message string            `json:"message"`

	// Financial Information
	Amount   float64 `json:"amount" gorm:"type:decimal(12,2);not null"`
	Currency string  `json:"currency" gorm:"default:'USD'"`

	// Gateway Information
	GatewayTransactionID string `json:"gateway_transaction_id"`
	PaymentMethodID      string `json:"payment_method_id"`

	// Authorization Information
	AuthorizationCode string `json:"authorization_code"`
	AVSResultCode     string `json:"avs_result_code"`
	CVVResultCode     string `json:"cvv_result_code"`

	// Metadata
	ProcessedAt *time.Time `json:"processed_at"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid;index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order    Order         `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Parent   *Transaction  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Transaction `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// TransactionKind represents the type of transaction
type TransactionKind string

const (
	TransactionKindSale          TransactionKind = "sale"
	TransactionKindAuthorization TransactionKind = "authorization"
	TransactionKindCapture       TransactionKind = "capture"
	TransactionKindVoid          TransactionKind = "void"
	TransactionKindRefund        TransactionKind = "refund"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailure TransactionStatus = "failure"
	TransactionStatusError   TransactionStatus = "error"
)

// Return represents a product return
type Return struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID      uuid.UUID    `json:"order_id" gorm:"type:uuid;not null;index"`
	ReturnNumber string       `json:"return_number" gorm:"unique;not null"`
	Status       ReturnStatus `json:"status" gorm:"default:'pending'"`

	// Return Information
	Reason        ReturnReason `json:"reason" gorm:"not null"`
	CustomerNotes string       `json:"customer_notes"`
	InternalNotes string       `json:"internal_notes"`

	// Financial Information
	RefundAmount  float64 `json:"refund_amount" gorm:"type:decimal(12,2)"`
	RestockingFee float64 `json:"restocking_fee" gorm:"type:decimal(10,2);default:0"`

	// Processing Information
	RestockItems   bool `json:"restock_items" gorm:"default:true"`
	NotifyCustomer bool `json:"notify_customer" gorm:"default:true"`
	RefundShipping bool `json:"refund_shipping" gorm:"default:false"`

	// Timestamps
	RequestedAt time.Time  `json:"requested_at" gorm:"autoCreateTime"`
	ProcessedAt *time.Time `json:"processed_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order     Order            `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	LineItems []ReturnLineItem `json:"line_items,omitempty" gorm:"foreignKey:ReturnID"`
}

// ReturnStatus represents the status of a return
type ReturnStatus string

const (
	ReturnStatusPending    ReturnStatus = "pending"
	ReturnStatusAuthorized ReturnStatus = "authorized"
	ReturnStatusReceived   ReturnStatus = "received"
	ReturnStatusProcessed  ReturnStatus = "processed"
	ReturnStatusCompleted  ReturnStatus = "completed"
	ReturnStatusDenied     ReturnStatus = "denied"
	ReturnStatusCancelled  ReturnStatus = "cancelled"
)

// ReturnReason represents the reason for a return
type ReturnReason string

const (
	ReturnReasonDefective      ReturnReason = "defective"
	ReturnReasonWrongItem      ReturnReason = "wrong_item"
	ReturnReasonNotAsDescribed ReturnReason = "not_as_described"
	ReturnReasonChangedMind    ReturnReason = "changed_mind"
	ReturnReasonTooLarge       ReturnReason = "too_large"
	ReturnReasonTooSmall       ReturnReason = "too_small"
	ReturnReasonDamaged        ReturnReason = "damaged"
	ReturnReasonOther          ReturnReason = "other"
)

// ReturnLineItem represents items in a return
type ReturnLineItem struct {
	ID              uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ReturnID        uuid.UUID    `json:"return_id" gorm:"type:uuid;not null;index"`
	LineItemID      uuid.UUID    `json:"line_item_id" gorm:"type:uuid;not null;index"`
	Quantity        int          `json:"quantity" gorm:"not null"`
	Reason          ReturnReason `json:"reason"`
	Notes           string       `json:"notes"`
	RefundAmount    float64      `json:"refund_amount" gorm:"type:decimal(10,2)"`
	RestockQuantity int          `json:"restock_quantity" gorm:"default:0"`
	CreatedAt       time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Return   Return        `json:"return,omitempty" gorm:"foreignKey:ReturnID"`
	LineItem OrderLineItem `json:"line_item,omitempty" gorm:"foreignKey:LineItemID"`
}

// OrderEvent represents events in the order lifecycle
type OrderEvent struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID     uuid.UUID              `json:"order_id" gorm:"type:uuid;not null;index"`
	EventType   OrderEventType         `json:"event_type" gorm:"not null"`
	Description string                 `json:"description"`
	UserID      *uuid.UUID             `json:"user_id" gorm:"type:uuid;index"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	Order Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderEventType represents the type of order event
type OrderEventType string

const (
	OrderEventCreated           OrderEventType = "created"
	OrderEventConfirmed         OrderEventType = "confirmed"
	OrderEventPaymentAuthorized OrderEventType = "payment_authorized"
	OrderEventPaymentCaptured   OrderEventType = "payment_captured"
	OrderEventFulfilled         OrderEventType = "fulfilled"
	OrderEventShipped           OrderEventType = "shipped"
	OrderEventDelivered         OrderEventType = "delivered"
	OrderEventCancelled         OrderEventType = "cancelled"
	OrderEventReturned          OrderEventType = "returned"
	OrderEventRefunded          OrderEventType = "refunded"
	OrderEventNoteAdded         OrderEventType = "note_added"
	OrderEventTagAdded          OrderEventType = "tag_added"
	OrderEventTagRemoved        OrderEventType = "tag_removed"
)

// BeforeCreate sets up UUID for new records
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (oli *OrderLineItem) BeforeCreate(tx *gorm.DB) error {
	if oli.ID == uuid.Nil {
		oli.ID = uuid.New()
	}
	return nil
}

func (f *Fulfillment) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

func (fli *FulfillmentLineItem) BeforeCreate(tx *gorm.DB) error {
	if fli.ID == uuid.Nil {
		fli.ID = uuid.New()
	}
	return nil
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (r *Return) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (rli *ReturnLineItem) BeforeCreate(tx *gorm.DB) error {
	if rli.ID == uuid.Nil {
		rli.ID = uuid.New()
	}
	return nil
}

func (oe *OrderEvent) BeforeCreate(tx *gorm.DB) error {
	if oe.ID == uuid.Nil {
		oe.ID = uuid.New()
	}
	return nil
}
