package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"unified-commerce/services/inventory/service"
	"unified-commerce/services/shared/logger"
)

// OrderEventHandler handles events related to orders.
type OrderEventHandler struct {
	inventoryService *service.InventoryService
	log              *logger.Logger
}

// NewOrderEventHandler creates a new OrderEventHandler.
func NewOrderEventHandler(inventoryService *service.InventoryService, log *logger.Logger) *OrderEventHandler {
	return &OrderEventHandler{
		inventoryService: inventoryService,
		log:              log,
	}
}

// OrderPlacedEvent represents the structure of an order placed event.
// We define it here to avoid a direct dependency on the order service's models.
type OrderPlacedEvent struct {
	ID        uuid.UUID          `json:"id"`
	LineItems []OrderLineItemDTO `json:"line_items"`
}

type OrderLineItemDTO struct {
	ProductVariantID *uuid.UUID `json:"product_variant_id"`
	Quantity         int        `json:"quantity"`
}

// HandleOrderPlacedEvent handles the order placed event from raw message bytes
func (h *OrderEventHandler) HandleOrderPlacedEvent(messageBytes []byte) error {
	h.log.Info("Received OrderPlaced event")

	var event OrderPlacedEvent
	if err := json.Unmarshal(messageBytes, &event); err != nil {
		h.log.WithError(err).Error("Failed to unmarshal OrderPlaced event")
		// Return nil to commit the message and prevent reprocessing of a malformed event.
		return nil
	}

	return h.processOrderPlacedEvent(event)
}

// HandleOrderPlaced is the legacy handler for the "orders.placed" topic (for backward compatibility)
func (h *OrderEventHandler) HandleOrderPlaced(messageBytes []byte) error {
	return h.HandleOrderPlacedEvent(messageBytes)
}

// processOrderPlacedEvent processes the actual order placed event
func (h *OrderEventHandler) processOrderPlacedEvent(event OrderPlacedEvent) error {

	h.log.WithField("order_id", event.ID).Info("Processing OrderPlaced event")

	// In a real application, you would determine the location from the order.
	// For now, we'll assume a default location or a simple lookup.
	// This is a placeholder for more complex logic.
	const defaultLocationID = "00000000-0000-0000-0000-000000000001" // Example UUID
	locationID, _ := uuid.Parse(defaultLocationID)

	adjustments := make([]service.StockAdjustment, len(event.LineItems))
	for i, item := range event.LineItems {
		if item.ProductVariantID != nil {
			adjustments[i] = service.StockAdjustment{
				ProductVariantID: *item.ProductVariantID,
				Quantity:         -item.Quantity, // Decrement stock
			}
		}
	}

	ctx := context.Background()
	if err := h.inventoryService.AdjustStock(ctx, locationID, "sale", fmt.Sprintf("Order %s", event.ID), adjustments); err != nil {
		h.log.WithError(err).WithField("order_id", event.ID).Error("Failed to adjust stock for order")
		// Return the error to signal that the message should be re-processed.
		return err
	}

	h.log.WithField("order_id", event.ID).Info("Successfully adjusted stock for order")
	return nil
}
