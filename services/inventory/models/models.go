package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Location represents a physical or virtual location where inventory is stored
type Location struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MerchantID  uuid.UUID `json:"merchant_id" gorm:"type:uuid;not null;index"`
	Name        string    `json:"name" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // warehouse, store, online, consignment
	Code        string    `json:"code" gorm:"unique;not null"`
	Description string    `json:"description"`
	Address     Address   `json:"address" gorm:"embedded;embeddedPrefix:address_"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Settings    Settings  `json:"settings" gorm:"type:jsonb"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Address represents the physical address of a location
type Address struct {
	Street1     string  `json:"street1"`
	Street2     string  `json:"street2"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	PostalCode  string  `json:"postal_code"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
	PhoneNumber string  `json:"phone_number"`
}

// Settings represents location-specific configuration
type Settings struct {
	AllowNegativeStock bool                   `json:"allow_negative_stock"`
	LowStockThreshold  int                    `json:"low_stock_threshold"`
	AutoReorderEnabled bool                   `json:"auto_reorder_enabled"`
	ReorderQuantity    int                    `json:"reorder_quantity"`
	BusinessHours      map[string]interface{} `json:"business_hours"`
	Notifications      map[string]interface{} `json:"notifications"`
}

// InventoryItem represents the inventory of a specific product variant at a location
type InventoryItem struct {
	ID                uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LocationID        uuid.UUID       `json:"location_id" gorm:"type:uuid;not null;index"`
	ProductID         uuid.UUID       `json:"product_id" gorm:"type:uuid;not null;index"`
	ProductVariantID  *uuid.UUID      `json:"product_variant_id" gorm:"type:uuid;index"`
	SKU               string          `json:"sku" gorm:"not null;index"`
	Quantity          int             `json:"quantity" gorm:"default:0"`
	ReservedQuantity  int             `json:"reserved_quantity" gorm:"default:0"`
	AvailableQuantity int             `json:"available_quantity" gorm:"default:0"` // Computed: Quantity - ReservedQuantity
	Cost              float64         `json:"cost" gorm:"type:decimal(12,2)"`
	RetailPrice       float64         `json:"retail_price" gorm:"type:decimal(12,2)"`
	LowStockThreshold int             `json:"low_stock_threshold" gorm:"default:10"`
	Bin               string          `json:"bin"` // Storage location within the location
	Status            InventoryStatus `json:"status" gorm:"default:'active'"`
	LastCountedAt     *time.Time      `json:"last_counted_at"`
	CreatedAt         time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time       `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Location Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// InventoryStatus represents the status of an inventory item
type InventoryStatus string

const (
	InventoryStatusActive       InventoryStatus = "active"
	InventoryStatusInactive     InventoryStatus = "inactive"
	InventoryStatusDamaged      InventoryStatus = "damaged"
	InventoryStatusDiscontinued InventoryStatus = "discontinued"
)

// StockMovement represents any change in inventory quantity
type StockMovement struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InventoryItemID  uuid.UUID      `json:"inventory_item_id" gorm:"type:uuid;not null;index"`
	LocationID       uuid.UUID      `json:"location_id" gorm:"type:uuid;not null;index"`
	ProductID        uuid.UUID      `json:"product_id" gorm:"type:uuid;not null;index"`
	ProductVariantID *uuid.UUID     `json:"product_variant_id" gorm:"type:uuid;index"`
	SKU              string         `json:"sku" gorm:"not null;index"`
	Type             MovementType   `json:"type" gorm:"not null"`
	Reason           MovementReason `json:"reason" gorm:"not null"`
	Quantity         int            `json:"quantity" gorm:"not null"`
	PreviousQuantity int            `json:"previous_quantity" gorm:"not null"`
	NewQuantity      int            `json:"new_quantity" gorm:"not null"`
	Cost             float64        `json:"cost" gorm:"type:decimal(12,2)"`
	Reference        string         `json:"reference"` // Order ID, Transfer ID, etc.
	Notes            string         `json:"notes"`
	UserID           *uuid.UUID     `json:"user_id" gorm:"type:uuid;index"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	InventoryItem InventoryItem `json:"inventory_item,omitempty" gorm:"foreignKey:InventoryItemID"`
	Location      Location      `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// MovementType represents the direction of stock movement
type MovementType string

const (
	MovementTypeIn  MovementType = "in"  // Stock increase
	MovementTypeOut MovementType = "out" // Stock decrease
)

// MovementReason represents the reason for stock movement
type MovementReason string

const (
	MovementReasonSale        MovementReason = "sale"
	MovementReasonPurchase    MovementReason = "purchase"
	MovementReasonReturn      MovementReason = "return"
	MovementReasonRefund      MovementReason = "refund"
	MovementReasonTransferIn  MovementReason = "transfer_in"
	MovementReasonTransferOut MovementReason = "transfer_out"
	MovementReasonAdjustment  MovementReason = "adjustment"
	MovementReasonDamage      MovementReason = "damage"
	MovementReasonLoss        MovementReason = "loss"
	MovementReasonFound       MovementReason = "found"
	MovementReasonCycle       MovementReason = "cycle_count"
	MovementReasonPromotion   MovementReason = "promotion"
	MovementReasonExpired     MovementReason = "expired"
)

// StockReservation represents a temporary hold on inventory
type StockReservation struct {
	ID               uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InventoryItemID  uuid.UUID         `json:"inventory_item_id" gorm:"type:uuid;not null;index"`
	LocationID       uuid.UUID         `json:"location_id" gorm:"type:uuid;not null;index"`
	ProductID        uuid.UUID         `json:"product_id" gorm:"type:uuid;not null;index"`
	ProductVariantID *uuid.UUID        `json:"product_variant_id" gorm:"type:uuid;index"`
	SKU              string            `json:"sku" gorm:"not null;index"`
	Quantity         int               `json:"quantity" gorm:"not null"`
	Type             ReservationType   `json:"type" gorm:"not null"`
	Reference        string            `json:"reference" gorm:"not null"` // Order ID, Cart ID, etc.
	Status           ReservationStatus `json:"status" gorm:"default:'active'"`
	ExpiresAt        *time.Time        `json:"expires_at"`
	UserID           *uuid.UUID        `json:"user_id" gorm:"type:uuid;index"`
	Notes            string            `json:"notes"`
	CreatedAt        time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time         `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	InventoryItem InventoryItem `json:"inventory_item,omitempty" gorm:"foreignKey:InventoryItemID"`
	Location      Location      `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// ReservationType represents the type of reservation
type ReservationType string

const (
	ReservationTypeOrder     ReservationType = "order"
	ReservationTypeCart      ReservationType = "cart"
	ReservationTypeTransfer  ReservationType = "transfer"
	ReservationTypePromotion ReservationType = "promotion"
	ReservationTypeBackorder ReservationType = "backorder"
)

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"
	ReservationStatusFulfilled ReservationStatus = "fulfilled"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusExpired   ReservationStatus = "expired"
)

// StockTransfer represents movement of inventory between locations
type StockTransfer struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TransferNumber string         `json:"transfer_number" gorm:"unique;not null"`
	FromLocationID uuid.UUID      `json:"from_location_id" gorm:"type:uuid;not null;index"`
	ToLocationID   uuid.UUID      `json:"to_location_id" gorm:"type:uuid;not null;index"`
	Status         TransferStatus `json:"status" gorm:"default:'pending'"`
	RequestedBy    uuid.UUID      `json:"requested_by" gorm:"type:uuid;not null"`
	ShippedBy      *uuid.UUID     `json:"shipped_by" gorm:"type:uuid"`
	ReceivedBy     *uuid.UUID     `json:"received_by" gorm:"type:uuid"`
	RequestedAt    time.Time      `json:"requested_at" gorm:"autoCreateTime"`
	ShippedAt      *time.Time     `json:"shipped_at"`
	ReceivedAt     *time.Time     `json:"received_at"`
	ExpectedAt     *time.Time     `json:"expected_at"`
	TrackingNumber string         `json:"tracking_number"`
	Carrier        string         `json:"carrier"`
	Notes          string         `json:"notes"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	FromLocation Location            `json:"from_location,omitempty" gorm:"foreignKey:FromLocationID"`
	ToLocation   Location            `json:"to_location,omitempty" gorm:"foreignKey:ToLocationID"`
	Items        []StockTransferItem `json:"items,omitempty" gorm:"foreignKey:TransferID"`
}

// TransferStatus represents the status of a stock transfer
type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusApproved  TransferStatus = "approved"
	TransferStatusShipped   TransferStatus = "shipped"
	TransferStatusInTransit TransferStatus = "in_transit"
	TransferStatusReceived  TransferStatus = "received"
	TransferStatusCancelled TransferStatus = "cancelled"
	TransferStatusPartial   TransferStatus = "partial"
)

// StockTransferItem represents an item in a stock transfer
type StockTransferItem struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TransferID        uuid.UUID  `json:"transfer_id" gorm:"type:uuid;not null;index"`
	ProductID         uuid.UUID  `json:"product_id" gorm:"type:uuid;not null;index"`
	ProductVariantID  *uuid.UUID `json:"product_variant_id" gorm:"type:uuid;index"`
	SKU               string     `json:"sku" gorm:"not null"`
	RequestedQuantity int        `json:"requested_quantity" gorm:"not null"`
	ShippedQuantity   int        `json:"shipped_quantity" gorm:"default:0"`
	ReceivedQuantity  int        `json:"received_quantity" gorm:"default:0"`
	DamagedQuantity   int        `json:"damaged_quantity" gorm:"default:0"`
	Cost              float64    `json:"cost" gorm:"type:decimal(12,2)"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Transfer StockTransfer `json:"transfer,omitempty" gorm:"foreignKey:TransferID"`
}

// StockAlert represents an alert for low stock or other inventory conditions
type StockAlert struct {
	ID               uuid.UUID     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LocationID       uuid.UUID     `json:"location_id" gorm:"type:uuid;not null;index"`
	ProductID        uuid.UUID     `json:"product_id" gorm:"type:uuid;not null;index"`
	ProductVariantID *uuid.UUID    `json:"product_variant_id" gorm:"type:uuid;index"`
	SKU              string        `json:"sku" gorm:"not null;index"`
	Type             AlertType     `json:"type" gorm:"not null"`
	Priority         AlertPriority `json:"priority" gorm:"default:'medium'"`
	Status           AlertStatus   `json:"status" gorm:"default:'active'"`
	Message          string        `json:"message" gorm:"not null"`
	Threshold        int           `json:"threshold"`
	CurrentQuantity  int           `json:"current_quantity"`
	AutoCreated      bool          `json:"auto_created" gorm:"default:true"`
	AcknowledgedBy   *uuid.UUID    `json:"acknowledged_by" gorm:"type:uuid"`
	AcknowledgedAt   *time.Time    `json:"acknowledged_at"`
	ResolvedAt       *time.Time    `json:"resolved_at"`
	CreatedAt        time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time     `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Location Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// AlertType represents the type of stock alert
type AlertType string

const (
	AlertTypeLowStock    AlertType = "low_stock"
	AlertTypeOutOfStock  AlertType = "out_of_stock"
	AlertTypeOverstock   AlertType = "overstock"
	AlertTypeExpiring    AlertType = "expiring"
	AlertTypeDamaged     AlertType = "damaged"
	AlertTypeDiscrepancy AlertType = "discrepancy"
)

// AlertPriority represents the priority level of an alert
type AlertPriority string

const (
	AlertPriorityLow      AlertPriority = "low"
	AlertPriorityMedium   AlertPriority = "medium"
	AlertPriorityHigh     AlertPriority = "high"
	AlertPriorityCritical AlertPriority = "critical"
)

// AlertStatus represents the status of an alert
type AlertStatus string

const (
	AlertStatusActive       AlertStatus = "active"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusResolved     AlertStatus = "resolved"
	AlertStatusDismissed    AlertStatus = "dismissed"
)

// BeforeCreate sets up UUID for new records
func (l *Location) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

func (i *InventoryItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (s *StockMovement) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (r *StockReservation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (t *StockTransfer) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (i *StockTransferItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

func (a *StockAlert) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// AfterSave updates available quantity for inventory items
func (i *InventoryItem) AfterSave(tx *gorm.DB) error {
	i.AvailableQuantity = i.Quantity - i.ReservedQuantity
	return nil
}

// InventoryLevel represents the current inventory level for a product variant at a location
type InventoryLevel struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LocationID        uuid.UUID  `json:"location_id" gorm:"type:uuid;not null;index"`
	ProductVariantID  uuid.UUID  `json:"product_variant_id" gorm:"type:uuid;not null;index"`
	SKU               string     `json:"sku" gorm:"not null;index"`
	Stock             int        `json:"stock" gorm:"default:0"`
	ReservedQuantity  int        `json:"reserved_quantity" gorm:"default:0"`
	AvailableQuantity int        `json:"available_quantity" gorm:"default:0"` // Computed: Stock - ReservedQuantity
	LowStockThreshold int        `json:"low_stock_threshold" gorm:"default:10"`
	LastCountedAt     *time.Time `json:"last_counted_at"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Location       Location       `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	ProductVariant ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
}

// StockAdjustment represents a manual adjustment to inventory levels
type StockAdjustment struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InventoryLevelID uuid.UUID      `json:"inventory_level_id" gorm:"type:uuid;not null;index"`
	LocationID       uuid.UUID      `json:"location_id" gorm:"type:uuid;not null;index"`
	ProductVariantID uuid.UUID      `json:"product_variant_id" gorm:"type:uuid;not null;index"`
	SKU              string         `json:"sku" gorm:"not null;index"`
	Type             AdjustmentType `json:"type" gorm:"not null"`
	Reason           string         `json:"reason" gorm:"not null"`
	Quantity         int            `json:"quantity" gorm:"not null"`
	PreviousQuantity int            `json:"previous_quantity" gorm:"not null"`
	NewQuantity      int            `json:"new_quantity" gorm:"not null"`
	Reference        string         `json:"reference"`
	Notes            string         `json:"notes"`
	UserID           *uuid.UUID     `json:"user_id" gorm:"type:uuid;index"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	InventoryLevel InventoryLevel `json:"inventory_level,omitempty" gorm:"foreignKey:InventoryLevelID"`
	Location       Location       `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	ProductVariant ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
	User           User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// AdjustmentType represents the type of stock adjustment
type AdjustmentType string

const (
	AdjustmentTypeIncrease AdjustmentType = "increase"
	AdjustmentTypeDecrease AdjustmentType = "decrease"
	AdjustmentTypeSet      AdjustmentType = "set"
)

// ProductVariant represents a product variant (simplified for inventory purposes)
type ProductVariant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	SKU       string    `json:"sku" gorm:"not null;unique"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// User represents a user (simplified for inventory purposes)
type User struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string    `json:"name" gorm:"not null"`
}
