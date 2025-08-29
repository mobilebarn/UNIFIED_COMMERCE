package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"unified-commerce/services/inventory/models"
	"unified-commerce/services/inventory/repository"
	"unified-commerce/services/shared/logger"
)

// Service errors
var (
	ErrLocationNotFound      = errors.New("location not found")
	ErrInventoryItemNotFound = errors.New("inventory item not found")
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrReservationNotFound   = errors.New("reservation not found")
	ErrTransferNotFound      = errors.New("transfer not found")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrLocationInactive      = errors.New("location is inactive")
	ErrDuplicateLocation     = errors.New("location code already exists")
)

// StockAdjustment represents a single adjustment to stock levels.
type StockAdjustment struct {
	ProductVariantID uuid.UUID
	Quantity         int
}

// InventoryService handles business logic for inventory management
type InventoryService struct {
	repo   *repository.InventoryRepository
	logger *logger.Logger
}

// NewInventoryService creates a new inventory service
func NewInventoryService(repo *repository.InventoryRepository, logger *logger.Logger) *InventoryService {
	return &InventoryService{
		repo:   repo,
		logger: logger,
	}
}

// Location Service Methods

// CreateLocationRequest represents a request to create a location
type CreateLocationRequest struct {
	MerchantID  uuid.UUID       `json:"merchant_id" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	Type        string          `json:"type" validate:"required,oneof=warehouse store online consignment"`
	Code        string          `json:"code" validate:"required"`
	Description string          `json:"description"`
	Address     models.Address  `json:"address"`
	Settings    models.Settings `json:"settings"`
}

// CreateLocation creates a new inventory location
func (s *InventoryService) CreateLocation(ctx context.Context, req *CreateLocationRequest) (*models.Location, error) {
	location := &models.Location{
		MerchantID:  req.MerchantID,
		Name:        req.Name,
		Type:        req.Type,
		Code:        req.Code,
		Description: req.Description,
		Address:     req.Address,
		IsActive:    true,
		Settings:    req.Settings,
	}

	if err := s.repo.CreateLocation(ctx, location); err != nil {
		s.logger.WithError(err).Error("Failed to create location")
		return nil, err
	}

	s.logger.WithField("location_id", location.ID).Info("Location created successfully")
	return location, nil
}

// AdjustStock adjusts the stock levels for multiple inventory items at a specific location.
func (s *InventoryService) AdjustStock(ctx context.Context, locationID uuid.UUID, reason, reference string, adjustments []StockAdjustment) error {
	// In a real application, this should be a single database transaction.
	tx := s.repo.BeginTx(ctx)
	if tx == nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer s.repo.RollbackTx(tx)

	for _, adj := range adjustments {
		err := s.repo.WithTx(tx).AdjustInventoryLevel(ctx, locationID, adj.ProductVariantID, adj.Quantity, reason, reference)
		if err != nil {
			s.logger.WithError(err).
				WithField("location_id", locationID).
				WithField("product_variant_id", adj.ProductVariantID).
				Error("Failed to adjust inventory level")
			return err
		}
	}

	return s.repo.CommitTx(tx)
}

// GetLocation retrieves a location by ID
func (s *InventoryService) GetLocation(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	location, err := s.repo.GetLocation(ctx, id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, ErrLocationNotFound
	}
	return location, nil
}

// GetLocationsByMerchant retrieves all locations for a merchant
func (s *InventoryService) GetLocationsByMerchant(ctx context.Context, merchantID uuid.UUID, page, limit int) ([]*models.Location, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetLocationsByMerchant(ctx, merchantID, limit, offset)
}

// UpdateLocationRequest represents a request to update a location
type UpdateLocationRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Address     models.Address  `json:"address"`
	IsActive    *bool           `json:"is_active"`
	Settings    models.Settings `json:"settings"`
}

// UpdateLocation updates a location
func (s *InventoryService) UpdateLocation(ctx context.Context, id uuid.UUID, req *UpdateLocationRequest) (*models.Location, error) {
	location, err := s.repo.GetLocation(ctx, id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, ErrLocationNotFound
	}

	// Update fields
	if req.Name != "" {
		location.Name = req.Name
	}
	if req.Description != "" {
		location.Description = req.Description
	}
	if req.IsActive != nil {
		location.IsActive = *req.IsActive
	}
	location.Address = req.Address
	location.Settings = req.Settings

	if err := s.repo.UpdateLocation(ctx, location); err != nil {
		s.logger.WithError(err).Error("Failed to update location")
		return nil, err
	}

	s.logger.WithField("location_id", location.ID).Info("Location updated successfully")
	return location, nil
}

// DeleteLocation deletes a location
func (s *InventoryService) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	location, err := s.repo.GetLocation(ctx, id)
	if err != nil {
		return err
	}
	if location == nil {
		return ErrLocationNotFound
	}

	if err := s.repo.DeleteLocation(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to delete location")
		return err
	}

	s.logger.WithField("location_id", id).Info("Location deleted successfully")
	return nil
}

// Inventory Item Service Methods

// CreateInventoryItemRequest represents a request to create an inventory item
type CreateInventoryItemRequest struct {
	LocationID        uuid.UUID  `json:"location_id" validate:"required"`
	ProductID         uuid.UUID  `json:"product_id" validate:"required"`
	ProductVariantID  *uuid.UUID `json:"product_variant_id"`
	SKU               string     `json:"sku" validate:"required"`
	Quantity          int        `json:"quantity" validate:"min=0"`
	Cost              float64    `json:"cost" validate:"min=0"`
	RetailPrice       float64    `json:"retail_price" validate:"min=0"`
	LowStockThreshold int        `json:"low_stock_threshold" validate:"min=0"`
	Bin               string     `json:"bin"`
}

// CreateInventoryItem creates a new inventory item
func (s *InventoryService) CreateInventoryItem(ctx context.Context, req *CreateInventoryItemRequest) (*models.InventoryItem, error) {
	// Verify location exists and is active
	location, err := s.repo.GetLocation(ctx, req.LocationID)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, ErrLocationNotFound
	}
	if !location.IsActive {
		return nil, ErrLocationInactive
	}

	// Check if inventory item already exists for this SKU and location
	existingItem, err := s.repo.GetInventoryItemBySKUAndLocation(ctx, req.SKU, req.LocationID)
	if err != nil {
		return nil, err
	}
	if existingItem != nil {
		return nil, fmt.Errorf("inventory item already exists for SKU %s at location %s", req.SKU, req.LocationID)
	}

	item := &models.InventoryItem{
		LocationID:        req.LocationID,
		ProductID:         req.ProductID,
		ProductVariantID:  req.ProductVariantID,
		SKU:               req.SKU,
		Quantity:          req.Quantity,
		ReservedQuantity:  0,
		AvailableQuantity: req.Quantity,
		Cost:              req.Cost,
		RetailPrice:       req.RetailPrice,
		LowStockThreshold: req.LowStockThreshold,
		Bin:               req.Bin,
		Status:            models.InventoryStatusActive,
	}

	if err := s.repo.CreateInventoryItem(ctx, item); err != nil {
		s.logger.WithError(err).Error("Failed to create inventory item")
		return nil, err
	}

	// Create initial stock movement if quantity > 0
	if req.Quantity > 0 {
		if err := s.repo.UpdateInventoryQuantity(ctx, item.ID, req.Quantity, models.MovementTypeIn, models.MovementReasonPurchase, "initial_stock", "Initial inventory", nil); err != nil {
			s.logger.WithError(err).Error("Failed to create initial stock movement")
		}
	}

	s.logger.WithField("inventory_item_id", item.ID).Info("Inventory item created successfully")
	return item, nil
}

// GetInventoryItem retrieves an inventory item by ID
func (s *InventoryService) GetInventoryItem(ctx context.Context, id uuid.UUID) (*models.InventoryItem, error) {
	item, err := s.repo.GetInventoryItem(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrInventoryItemNotFound
	}
	return item, nil
}

// GetInventoryItemsByLocation retrieves all inventory items for a location
func (s *InventoryService) GetInventoryItemsByLocation(ctx context.Context, locationID uuid.UUID, page, limit int) ([]*models.InventoryItem, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetInventoryItemsByLocation(ctx, locationID, limit, offset)
}

// GetInventoryItemsByProduct retrieves all inventory items for a product across locations
func (s *InventoryService) GetInventoryItemsByProduct(ctx context.Context, productID uuid.UUID) ([]*models.InventoryItem, error) {
	return s.repo.GetInventoryItemsByProduct(ctx, productID)
}

// UpdateInventoryItemRequest represents a request to update an inventory item
type UpdateInventoryItemRequest struct {
	Cost              *float64               `json:"cost"`
	RetailPrice       *float64               `json:"retail_price"`
	LowStockThreshold *int                   `json:"low_stock_threshold"`
	Bin               string                 `json:"bin"`
	Status            models.InventoryStatus `json:"status"`
}

// UpdateInventoryItem updates an inventory item
func (s *InventoryService) UpdateInventoryItem(ctx context.Context, id uuid.UUID, req *UpdateInventoryItemRequest) (*models.InventoryItem, error) {
	item, err := s.repo.GetInventoryItem(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrInventoryItemNotFound
	}

	// Update fields
	if req.Cost != nil {
		item.Cost = *req.Cost
	}
	if req.RetailPrice != nil {
		item.RetailPrice = *req.RetailPrice
	}
	if req.LowStockThreshold != nil {
		item.LowStockThreshold = *req.LowStockThreshold
	}
	if req.Bin != "" {
		item.Bin = req.Bin
	}
	if req.Status != "" {
		item.Status = req.Status
	}

	if err := s.repo.UpdateInventoryItem(ctx, item); err != nil {
		s.logger.WithError(err).Error("Failed to update inventory item")
		return nil, err
	}

	s.logger.WithField("inventory_item_id", item.ID).Info("Inventory item updated successfully")
	return item, nil
}

// AdjustInventoryRequest represents a request to adjust inventory quantity
type AdjustInventoryRequest struct {
	Quantity  int                   `json:"quantity" validate:"required"`
	Reason    models.MovementReason `json:"reason" validate:"required"`
	Reference string                `json:"reference"`
	Notes     string                `json:"notes"`
	UserID    *uuid.UUID            `json:"user_id"`
}

// AdjustInventory adjusts the quantity of an inventory item
func (s *InventoryService) AdjustInventory(ctx context.Context, id uuid.UUID, req *AdjustInventoryRequest) (*models.InventoryItem, error) {
	item, err := s.repo.GetInventoryItem(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrInventoryItemNotFound
	}

	if req.Quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	movementType := models.MovementTypeIn
	if req.Quantity < item.Quantity {
		movementType = models.MovementTypeOut
	}

	if err := s.repo.UpdateInventoryQuantity(ctx, id, req.Quantity, movementType, req.Reason, req.Reference, req.Notes, req.UserID); err != nil {
		s.logger.WithError(err).Error("Failed to adjust inventory")
		return nil, err
	}

	// Get updated item
	updatedItem, err := s.repo.GetInventoryItem(ctx, id)
	if err != nil {
		return nil, err
	}

	s.logger.WithField("inventory_item_id", id).WithField("new_quantity", req.Quantity).Info("Inventory adjusted successfully")
	return updatedItem, nil
}

// GetLowStockItems retrieves inventory items with low stock
func (s *InventoryService) GetLowStockItems(ctx context.Context, locationID *uuid.UUID, page, limit int) ([]*models.InventoryItem, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetLowStockItems(ctx, locationID, limit, offset)
}

// Stock Movement Service Methods

// GetStockMovements retrieves stock movements with filters
func (s *InventoryService) GetStockMovements(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*models.StockMovement, int64, error) {
	offset := (page - 1) * limit
	return s.repo.GetStockMovements(ctx, filters, limit, offset)
}

// Stock Reservation Service Methods

// CreateReservationRequest represents a request to create a stock reservation
type CreateReservationRequest struct {
	InventoryItemID uuid.UUID              `json:"inventory_item_id" validate:"required"`
	Quantity        int                    `json:"quantity" validate:"required,min=1"`
	Type            models.ReservationType `json:"type" validate:"required"`
	Reference       string                 `json:"reference" validate:"required"`
	ExpiresAt       *time.Time             `json:"expires_at"`
	UserID          *uuid.UUID             `json:"user_id"`
	Notes           string                 `json:"notes"`
}

// CreateReservation creates a new stock reservation
func (s *InventoryService) CreateReservation(ctx context.Context, req *CreateReservationRequest) (*models.StockReservation, error) {
	// Verify inventory item exists
	item, err := s.repo.GetInventoryItem(ctx, req.InventoryItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrInventoryItemNotFound
	}

	reservation := &models.StockReservation{
		InventoryItemID:  req.InventoryItemID,
		LocationID:       item.LocationID,
		ProductID:        item.ProductID,
		ProductVariantID: item.ProductVariantID,
		SKU:              item.SKU,
		Quantity:         req.Quantity,
		Type:             req.Type,
		Reference:        req.Reference,
		Status:           models.ReservationStatusActive,
		ExpiresAt:        req.ExpiresAt,
		UserID:           req.UserID,
		Notes:            req.Notes,
	}

	if err := s.repo.CreateReservation(ctx, reservation); err != nil {
		s.logger.WithError(err).Error("Failed to create reservation")
		if err.Error() == "insufficient stock available" {
			return nil, ErrInsufficientStock
		}
		return nil, err
	}

	s.logger.WithField("reservation_id", reservation.ID).Info("Stock reservation created successfully")
	return reservation, nil
}

// GetReservation retrieves a reservation by ID
func (s *InventoryService) GetReservation(ctx context.Context, id uuid.UUID) (*models.StockReservation, error) {
	reservation, err := s.repo.GetReservation(ctx, id)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

// GetReservationsByReference retrieves reservations by reference
func (s *InventoryService) GetReservationsByReference(ctx context.Context, reference string) ([]*models.StockReservation, error) {
	return s.repo.GetReservationsByReference(ctx, reference)
}

// FulfillReservation fulfills a stock reservation
func (s *InventoryService) FulfillReservation(ctx context.Context, id uuid.UUID, actualQuantity int, userID *uuid.UUID) error {
	reservation, err := s.repo.GetReservation(ctx, id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return ErrReservationNotFound
	}

	if reservation.Status != models.ReservationStatusActive {
		return fmt.Errorf("reservation is not active")
	}

	if actualQuantity <= 0 {
		return ErrInvalidQuantity
	}

	if err := s.repo.FulfillReservation(ctx, id, actualQuantity, userID); err != nil {
		s.logger.WithError(err).Error("Failed to fulfill reservation")
		return err
	}

	s.logger.WithField("reservation_id", id).WithField("actual_quantity", actualQuantity).Info("Reservation fulfilled successfully")
	return nil
}

// CancelReservation cancels a stock reservation
func (s *InventoryService) CancelReservation(ctx context.Context, id uuid.UUID) error {
	reservation, err := s.repo.GetReservation(ctx, id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return ErrReservationNotFound
	}

	if reservation.Status != models.ReservationStatusActive {
		return fmt.Errorf("reservation is not active")
	}

	if err := s.repo.CancelReservation(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to cancel reservation")
		return err
	}

	s.logger.WithField("reservation_id", id).Info("Reservation cancelled successfully")
	return nil
}

// Utility Methods

// CheckStockAvailability checks if a product has sufficient stock at a location
func (s *InventoryService) CheckStockAvailability(ctx context.Context, sku string, locationID uuid.UUID, requiredQuantity int) (bool, int, error) {
	item, err := s.repo.GetInventoryItemBySKUAndLocation(ctx, sku, locationID)
	if err != nil {
		return false, 0, err
	}
	if item == nil {
		return false, 0, nil
	}

	availableQuantity := item.Quantity - item.ReservedQuantity
	return availableQuantity >= requiredQuantity, availableQuantity, nil
}

// GetInventorySummaryByLocation provides inventory summary for a location
func (s *InventoryService) GetInventorySummaryByLocation(ctx context.Context, locationID uuid.UUID) (map[string]interface{}, error) {
	items, _, err := s.repo.GetInventoryItemsByLocation(ctx, locationID, 1000, 0) // Get all items
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"total_items":        len(items),
		"total_quantity":     0,
		"total_value":        0.0,
		"low_stock_items":    0,
		"out_of_stock_items": 0,
		"inactive_items":     0,
	}

	for _, item := range items {
		summary["total_quantity"] = summary["total_quantity"].(int) + item.Quantity
		summary["total_value"] = summary["total_value"].(float64) + (float64(item.Quantity) * item.Cost)

		if item.Status != models.InventoryStatusActive {
			summary["inactive_items"] = summary["inactive_items"].(int) + 1
		} else if item.Quantity == 0 {
			summary["out_of_stock_items"] = summary["out_of_stock_items"].(int) + 1
		} else if item.Quantity <= item.LowStockThreshold {
			summary["low_stock_items"] = summary["low_stock_items"].(int) + 1
		}
	}

	return summary, nil
}

// ProcessExpiredReservations processes and expires reservations that have passed their expiry time
func (s *InventoryService) ProcessExpiredReservations(ctx context.Context) error {
	if err := s.repo.ExpireReservations(ctx); err != nil {
		s.logger.WithError(err).Error("Failed to process expired reservations")
		return err
	}

	s.logger.Info("Expired reservations processed successfully")
	return nil
}
