package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"unified-commerce/services/inventory/models"
	"unified-commerce/services/shared/logger"
)

// InventoryRepository handles database operations for inventory management
type InventoryRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *gorm.DB, logger *logger.Logger) *InventoryRepository {
	return &InventoryRepository{
		db:     db,
		logger: logger,
	}
}

// Location Operations

// CreateLocation creates a new inventory location
func (r *InventoryRepository) CreateLocation(ctx context.Context, location *models.Location) error {
	if err := r.db.WithContext(ctx).Create(location).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create location")
		return err
	}
	return nil
}

// GetLocation retrieves a location by ID
func (r *InventoryRepository) GetLocation(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	var location models.Location
	if err := r.db.WithContext(ctx).First(&location, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get location")
		return nil, err
	}
	return &location, nil
}

// GetLocationsByMerchant retrieves all locations for a merchant
func (r *InventoryRepository) GetLocationsByMerchant(ctx context.Context, merchantID uuid.UUID, limit, offset int) ([]*models.Location, int64, error) {
	var locations []*models.Location
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Location{}).
		Where("merchant_id = ?", merchantID).Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count locations")
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&locations).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get locations")
		return nil, 0, err
	}

	return locations, total, nil
}

// UpdateLocation updates a location
func (r *InventoryRepository) UpdateLocation(ctx context.Context, location *models.Location) error {
	if err := r.db.WithContext(ctx).Save(location).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update location")
		return err
	}
	return nil
}

// DeleteLocation soft deletes a location
func (r *InventoryRepository) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Location{}, "id = ?", id).Error; err != nil {
		r.logger.WithError(err).Error("Failed to delete location")
		return err
	}
	return nil
}

// Inventory Item Operations

// CreateInventoryItem creates a new inventory item
func (r *InventoryRepository) CreateInventoryItem(ctx context.Context, item *models.InventoryItem) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create inventory item")
		return err
	}
	return nil
}

// GetInventoryItem retrieves an inventory item by ID
func (r *InventoryRepository) GetInventoryItem(ctx context.Context, id uuid.UUID) (*models.InventoryItem, error) {
	var item models.InventoryItem
	if err := r.db.WithContext(ctx).Preload("Location").First(&item, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get inventory item")
		return nil, err
	}
	return &item, nil
}

// GetInventoryItemBySKUAndLocation retrieves an inventory item by SKU and location
func (r *InventoryRepository) GetInventoryItemBySKUAndLocation(ctx context.Context, sku string, locationID uuid.UUID) (*models.InventoryItem, error) {
	var item models.InventoryItem
	if err := r.db.WithContext(ctx).Preload("Location").
		First(&item, "sku = ? AND location_id = ?", sku, locationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get inventory item by SKU and location")
		return nil, err
	}
	return &item, nil
}

// GetInventoryItemsByLocation retrieves all inventory items for a location
func (r *InventoryRepository) GetInventoryItemsByLocation(ctx context.Context, locationID uuid.UUID, limit, offset int) ([]*models.InventoryItem, int64, error) {
	var items []*models.InventoryItem
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.InventoryItem{}).
		Where("location_id = ?", locationID).Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count inventory items")
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.WithContext(ctx).Preload("Location").
		Where("location_id = ?", locationID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&items).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get inventory items")
		return nil, 0, err
	}

	return items, total, nil
}

// GetInventoryItemsByProduct retrieves all inventory items for a product across locations
func (r *InventoryRepository) GetInventoryItemsByProduct(ctx context.Context, productID uuid.UUID) ([]*models.InventoryItem, error) {
	var items []*models.InventoryItem
	if err := r.db.WithContext(ctx).Preload("Location").
		Where("product_id = ?", productID).
		Find(&items).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get inventory items by product")
		return nil, err
	}
	return items, nil
}

// UpdateInventoryItem updates an inventory item
func (r *InventoryRepository) UpdateInventoryItem(ctx context.Context, item *models.InventoryItem) error {
	if err := r.db.WithContext(ctx).Save(item).Error; err != nil {
		r.logger.WithError(err).Error("Failed to update inventory item")
		return err
	}
	return nil
}

// UpdateInventoryQuantity updates the quantity of an inventory item and creates a stock movement
func (r *InventoryRepository) UpdateInventoryQuantity(ctx context.Context, itemID uuid.UUID, newQuantity int, movementType models.MovementType, reason models.MovementReason, reference, notes string, userID *uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current inventory item
		var item models.InventoryItem
		if err := tx.First(&item, "id = ?", itemID).Error; err != nil {
			return err
		}

		previousQuantity := item.Quantity
		quantityChange := newQuantity - previousQuantity

		// Update inventory item
		item.Quantity = newQuantity
		item.AvailableQuantity = newQuantity - item.ReservedQuantity
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		// Create stock movement record
		movement := &models.StockMovement{
			InventoryItemID:  itemID,
			LocationID:       item.LocationID,
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			SKU:              item.SKU,
			Type:             movementType,
			Reason:           reason,
			Quantity:         quantityChange,
			PreviousQuantity: previousQuantity,
			NewQuantity:      newQuantity,
			Cost:             item.Cost,
			Reference:        reference,
			Notes:            notes,
			UserID:           userID,
		}
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetLowStockItems retrieves inventory items with low stock
func (r *InventoryRepository) GetLowStockItems(ctx context.Context, locationID *uuid.UUID, limit, offset int) ([]*models.InventoryItem, int64, error) {
	var items []*models.InventoryItem
	var total int64

	query := r.db.WithContext(ctx).Model(&models.InventoryItem{}).
		Where("quantity <= low_stock_threshold AND status = ?", models.InventoryStatusActive)

	if locationID != nil {
		query = query.Where("location_id = ?", *locationID)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count low stock items")
		return nil, 0, err
	}

	// Get paginated results
	query = r.db.WithContext(ctx).Preload("Location").
		Where("quantity <= low_stock_threshold AND status = ?", models.InventoryStatusActive)

	if locationID != nil {
		query = query.Where("location_id = ?", *locationID)
	}

	if err := query.Order("quantity ASC").
		Limit(limit).Offset(offset).
		Find(&items).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get low stock items")
		return nil, 0, err
	}

	return items, total, nil
}

// Stock Movement Operations

// CreateStockMovement creates a new stock movement record
func (r *InventoryRepository) CreateStockMovement(ctx context.Context, movement *models.StockMovement) error {
	if err := r.db.WithContext(ctx).Create(movement).Error; err != nil {
		r.logger.WithError(err).Error("Failed to create stock movement")
		return err
	}
	return nil
}

// GetStockMovements retrieves stock movements with filters
func (r *InventoryRepository) GetStockMovements(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*models.StockMovement, int64, error) {
	var movements []*models.StockMovement
	var total int64

	query := r.db.WithContext(ctx).Model(&models.StockMovement{})

	// Apply filters
	for key, value := range filters {
		switch key {
		case "location_id":
			query = query.Where("location_id = ?", value)
		case "product_id":
			query = query.Where("product_id = ?", value)
		case "sku":
			query = query.Where("sku = ?", value)
		case "type":
			query = query.Where("type = ?", value)
		case "reason":
			query = query.Where("reason = ?", value)
		case "reference":
			query = query.Where("reference = ?", value)
		case "date_from":
			query = query.Where("created_at >= ?", value)
		case "date_to":
			query = query.Where("created_at <= ?", value)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		r.logger.WithError(err).Error("Failed to count stock movements")
		return nil, 0, err
	}

	// Get paginated results
	query = r.db.WithContext(ctx).Preload("InventoryItem").Preload("Location")

	// Apply filters again
	for key, value := range filters {
		switch key {
		case "location_id":
			query = query.Where("location_id = ?", value)
		case "product_id":
			query = query.Where("product_id = ?", value)
		case "sku":
			query = query.Where("sku = ?", value)
		case "type":
			query = query.Where("type = ?", value)
		case "reason":
			query = query.Where("reason = ?", value)
		case "reference":
			query = query.Where("reference = ?", value)
		case "date_from":
			query = query.Where("created_at >= ?", value)
		case "date_to":
			query = query.Where("created_at <= ?", value)
		}
	}

	if err := query.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&movements).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get stock movements")
		return nil, 0, err
	}

	return movements, total, nil
}

// Stock Reservation Operations

// CreateReservation creates a new stock reservation
func (r *InventoryRepository) CreateReservation(ctx context.Context, reservation *models.StockReservation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get inventory item
		var item models.InventoryItem
		if err := tx.First(&item, "id = ?", reservation.InventoryItemID).Error; err != nil {
			return err
		}

		// Check if enough stock is available
		availableQuantity := item.Quantity - item.ReservedQuantity
		if availableQuantity < reservation.Quantity {
			return fmt.Errorf("insufficient stock available: %d requested, %d available", reservation.Quantity, availableQuantity)
		}

		// Create reservation
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		// Update reserved quantity
		item.ReservedQuantity += reservation.Quantity
		item.AvailableQuantity = item.Quantity - item.ReservedQuantity
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetReservation retrieves a reservation by ID
func (r *InventoryRepository) GetReservation(ctx context.Context, id uuid.UUID) (*models.StockReservation, error) {
	var reservation models.StockReservation
	if err := r.db.WithContext(ctx).Preload("InventoryItem").Preload("Location").
		First(&reservation, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.WithError(err).Error("Failed to get reservation")
		return nil, err
	}
	return &reservation, nil
}

// GetReservationsByReference retrieves reservations by reference (order ID, cart ID, etc.)
func (r *InventoryRepository) GetReservationsByReference(ctx context.Context, reference string) ([]*models.StockReservation, error) {
	var reservations []*models.StockReservation
	if err := r.db.WithContext(ctx).Preload("InventoryItem").Preload("Location").
		Where("reference = ?", reference).
		Find(&reservations).Error; err != nil {
		r.logger.WithError(err).Error("Failed to get reservations by reference")
		return nil, err
	}
	return reservations, nil
}

// UpdateReservation updates a reservation
func (r *InventoryRepository) UpdateReservation(ctx context.Context, reservation *models.StockReservation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current reservation
		var currentReservation models.StockReservation
		if err := tx.First(&currentReservation, "id = ?", reservation.ID).Error; err != nil {
			return err
		}

		// Get inventory item
		var item models.InventoryItem
		if err := tx.First(&item, "id = ?", currentReservation.InventoryItemID).Error; err != nil {
			return err
		}

		// Calculate quantity difference
		quantityDiff := reservation.Quantity - currentReservation.Quantity

		// Update reservation
		if err := tx.Save(reservation).Error; err != nil {
			return err
		}

		// Update reserved quantity in inventory item
		item.ReservedQuantity += quantityDiff
		item.AvailableQuantity = item.Quantity - item.ReservedQuantity
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		return nil
	})
}

// FulfillReservation fulfills a reservation and updates inventory
func (r *InventoryRepository) FulfillReservation(ctx context.Context, reservationID uuid.UUID, actualQuantity int, userID *uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get reservation
		var reservation models.StockReservation
		if err := tx.First(&reservation, "id = ?", reservationID).Error; err != nil {
			return err
		}

		// Get inventory item
		var item models.InventoryItem
		if err := tx.First(&item, "id = ?", reservation.InventoryItemID).Error; err != nil {
			return err
		}

		// Update reservation status
		reservation.Status = models.ReservationStatusFulfilled
		if err := tx.Save(&reservation).Error; err != nil {
			return err
		}

		// Update inventory quantities
		item.Quantity -= actualQuantity
		item.ReservedQuantity -= reservation.Quantity
		item.AvailableQuantity = item.Quantity - item.ReservedQuantity
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		// Create stock movement
		movement := &models.StockMovement{
			InventoryItemID:  reservation.InventoryItemID,
			LocationID:       reservation.LocationID,
			ProductID:        reservation.ProductID,
			ProductVariantID: reservation.ProductVariantID,
			SKU:              reservation.SKU,
			Type:             models.MovementTypeOut,
			Reason:           models.MovementReasonSale,
			Quantity:         actualQuantity,
			PreviousQuantity: item.Quantity + actualQuantity,
			NewQuantity:      item.Quantity,
			Cost:             item.Cost,
			Reference:        reservation.Reference,
			Notes:            fmt.Sprintf("Fulfilled reservation %s", reservationID),
			UserID:           userID,
		}
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})
}

// CancelReservation cancels a reservation and releases inventory
func (r *InventoryRepository) CancelReservation(ctx context.Context, reservationID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get reservation
		var reservation models.StockReservation
		if err := tx.First(&reservation, "id = ?", reservationID).Error; err != nil {
			return err
		}

		// Get inventory item
		var item models.InventoryItem
		if err := tx.First(&item, "id = ?", reservation.InventoryItemID).Error; err != nil {
			return err
		}

		// Update reservation status
		reservation.Status = models.ReservationStatusCancelled
		if err := tx.Save(&reservation).Error; err != nil {
			return err
		}

		// Release reserved quantity
		item.ReservedQuantity -= reservation.Quantity
		item.AvailableQuantity = item.Quantity - item.ReservedQuantity
		if err := tx.Save(&item).Error; err != nil {
			return err
		}

		return nil
	})
}

// ExpireReservations expires reservations that have passed their expiry time
func (r *InventoryRepository) ExpireReservations(ctx context.Context) error {
	now := time.Now()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get expired reservations
		var expiredReservations []*models.StockReservation
		if err := tx.Where("status = ? AND expires_at < ?", models.ReservationStatusActive, now).
			Find(&expiredReservations).Error; err != nil {
			return err
		}

		for _, reservation := range expiredReservations {
			// Update reservation status
			reservation.Status = models.ReservationStatusExpired
			if err := tx.Save(reservation).Error; err != nil {
				return err
			}

			// Release reserved quantity
			var item models.InventoryItem
			if err := tx.First(&item, "id = ?", reservation.InventoryItemID).Error; err != nil {
				return err
			}

			item.ReservedQuantity -= reservation.Quantity
			item.AvailableQuantity = item.Quantity - item.ReservedQuantity
			if err := tx.Save(&item).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
