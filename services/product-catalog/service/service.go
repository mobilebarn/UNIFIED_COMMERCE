package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"unified-commerce/services/product-catalog/models"
	"unified-commerce/services/product-catalog/repository"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/utils"
)

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrProductExists      = errors.New("product already exists")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrInvalidSKU         = errors.New("invalid SKU format")
	ErrInvalidSlug        = errors.New("invalid slug format")
)

// ProductService handles product catalog business logic
type ProductService struct {
	repo   *repository.Repository
	logger *logger.Logger
}

// NewProductService creates a new product service
func NewProductService(repo *repository.Repository, logger *logger.Logger) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logger,
	}
}

// CreateProductRequest represents a product creation request
type CreateProductRequest struct {
	MerchantID      string                  `json:"merchant_id" validate:"required"`
	Name            string                  `json:"name" validate:"required"`
	Description     string                  `json:"description"`
	SKU             string                  `json:"sku"`
	Type            string                  `json:"type" validate:"required"`
	Price           float64                 `json:"price" validate:"min=0"`
	CompareAtPrice  float64                 `json:"compare_at_price" validate:"min=0"`
	CostPrice       float64                 `json:"cost_price" validate:"min=0"`
	Weight          float64                 `json:"weight" validate:"min=0"`
	WeightUnit      string                  `json:"weight_unit"`
	TrackInventory  bool                    `json:"track_inventory"`
	InventoryPolicy string                  `json:"inventory_policy"`
	Categories      []string                `json:"categories"`
	Tags            []string                `json:"tags"`
	Images          []models.ProductImage   `json:"images"`
	Variants        []ProductVariantRequest `json:"variants"`
	Options         []models.ProductOption  `json:"options"`
	Attributes      map[string]interface{}  `json:"attributes"`
	SEO             models.ProductSEO       `json:"seo"`
}

// ProductVariantRequest represents a variant creation request
type ProductVariantRequest struct {
	SKU            string                 `json:"sku"`
	Price          float64                `json:"price"`
	CompareAtPrice float64                `json:"compare_at_price"`
	CostPrice      float64                `json:"cost_price"`
	Weight         float64                `json:"weight"`
	OptionValues   map[string]string      `json:"option_values"`
	Attributes     map[string]interface{} `json:"attributes"`
}

// UpdateProductRequest represents a product update request
type UpdateProductRequest struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Price           *float64               `json:"price"`
	CompareAtPrice  *float64               `json:"compare_at_price"`
	CostPrice       *float64               `json:"cost_price"`
	Weight          *float64               `json:"weight"`
	TrackInventory  *bool                  `json:"track_inventory"`
	InventoryPolicy string                 `json:"inventory_policy"`
	Categories      []string               `json:"categories"`
	Tags            []string               `json:"tags"`
	Images          []models.ProductImage  `json:"images"`
	Attributes      map[string]interface{} `json:"attributes"`
	SEO             models.ProductSEO      `json:"seo"`
}

// ProductListFilters represents filters for product listing
type ProductListFilters struct {
	MerchantID string  `json:"merchant_id"`
	Status     string  `json:"status"`
	Category   string  `json:"category"`
	Tag        string  `json:"tag"`
	Search     string  `json:"search"`
	PriceMin   float64 `json:"price_min"`
	PriceMax   float64 `json:"price_max"`
	Sort       string  `json:"sort"`
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*models.Product, error) {
	// Check if repository is available
	if s.repo.Product == nil {
		return nil, fmt.Errorf("product catalog service is running in degraded mode - database unavailable")
	}
	
	// Generate SKU if not provided
	if req.SKU == "" {
		req.SKU = s.generateSKU(req.Name)
	}

	// Check if SKU already exists
	existing, err := s.repo.Product.GetBySKU(ctx, req.MerchantID, req.SKU)
	if err == nil && existing != nil {
		return nil, ErrProductExists
	}

	// Generate slug from name
	slug := utils.SlugifyString(req.Name)

	// Ensure slug is unique
	slugCounter := 1
	originalSlug := slug
	for {
		_, err := s.repo.Product.GetBySlug(ctx, req.MerchantID, slug)
		if err != nil {
			break // Slug is available
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, slugCounter)
		slugCounter++
	}

	// Create product model
	product := &models.Product{
		MerchantID:       req.MerchantID,
		SKU:              req.SKU,
		Name:             req.Name,
		Slug:             slug,
		Description:      req.Description,
		Status:           "active",
		Type:             req.Type,
		Visibility:       "public",
		Price:            req.Price,
		CompareAtPrice:   req.CompareAtPrice,
		CostPrice:        req.CostPrice,
		Currency:         "USD", // Default, should come from merchant settings
		Weight:           req.Weight,
		WeightUnit:       req.WeightUnit,
		TrackInventory:   req.TrackInventory,
		InventoryPolicy:  req.InventoryPolicy,
		Categories:       req.Categories,
		Tags:             req.Tags,
		Images:           req.Images,
		Options:          req.Options,
		SEO:              req.SEO,
		Attributes:       req.Attributes,
		RequiresShipping: true, // Default
	}

	// Set default values
	if product.WeightUnit == "" {
		product.WeightUnit = "kg"
	}
	if product.InventoryPolicy == "" {
		product.InventoryPolicy = "deny"
	}

	// Process variants if provided
	if len(req.Variants) > 0 {
		for i, variantReq := range req.Variants {
			variant := models.ProductVariant{
				ID:             utils.GenerateID(),
				SKU:            variantReq.SKU,
				Price:          variantReq.Price,
				CompareAtPrice: variantReq.CompareAtPrice,
				CostPrice:      variantReq.CostPrice,
				Weight:         variantReq.Weight,
				OptionValues:   variantReq.OptionValues,
				Attributes:     variantReq.Attributes,
				Position:       i,
				IsDefault:      i == 0, // First variant is default
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// Generate variant SKU if not provided
			if variant.SKU == "" {
				variant.SKU = fmt.Sprintf("%s-VAR-%d", product.SKU, i+1)
			}

			product.Variants = append(product.Variants, variant)
		}
	}

	// Set primary image
	if len(product.Images) > 0 {
		product.Images[0].IsPrimary = true
	}

	// Create product in database
	if err := s.repo.Product.Create(ctx, product); err != nil {
		s.logger.WithError(err).Error("Failed to create product")
		return nil, fmt.Errorf("failed to create product")
	}

	s.logger.WithFields(map[string]interface{}{
		"product_id":  product.ID.Hex(),
		"merchant_id": product.MerchantID,
		"sku":         product.SKU,
	}).Info("Product created successfully")

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, productID, merchantID string) (*models.Product, error) {
	// Check if repository is available
	if s.repo.Product == nil {
		return nil, fmt.Errorf("product catalog service is running in degraded mode - database unavailable")
	}
	
	product, err := s.repo.Product.GetByID(ctx, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Check merchant access
	if product.MerchantID != merchantID {
		return nil, ErrUnauthorized
	}

	return product, nil
}

// GetProductBySlug retrieves a product by slug (public access)
func (s *ProductService) GetProductBySlug(ctx context.Context, merchantID, slug string) (*models.Product, error) {
	// Check if repository is available
	if s.repo.Product == nil {
		return nil, fmt.Errorf("product catalog service is running in degraded mode - database unavailable")
	}
	
	product, err := s.repo.Product.GetBySlug(ctx, merchantID, slug)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Only return published products for public access
	if !product.IsPublished() {
		return nil, ErrProductNotFound
	}

	return product, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, productID, merchantID string, req *UpdateProductRequest) (*models.Product, error) {
	product, err := s.repo.Product.GetByID(ctx, productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Check merchant access
	if product.MerchantID != merchantID {
		return nil, ErrUnauthorized
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
		// Regenerate slug if name changed
		slug := utils.SlugifyString(req.Name)
		if slug != product.Slug {
			product.Slug = s.ensureUniqueSlug(ctx, merchantID, slug, productID)
		}
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.CompareAtPrice != nil {
		product.CompareAtPrice = *req.CompareAtPrice
	}

	if req.CostPrice != nil {
		product.CostPrice = *req.CostPrice
	}

	if req.Weight != nil {
		product.Weight = *req.Weight
	}

	if req.TrackInventory != nil {
		product.TrackInventory = *req.TrackInventory
	}

	if req.InventoryPolicy != "" {
		product.InventoryPolicy = req.InventoryPolicy
	}

	if req.Categories != nil {
		product.Categories = req.Categories
	}

	if req.Tags != nil {
		product.Tags = req.Tags
	}

	if req.Images != nil {
		product.Images = req.Images
		// Ensure at least one primary image
		if len(product.Images) > 0 {
			hasPrimary := false
			for _, img := range product.Images {
				if img.IsPrimary {
					hasPrimary = true
					break
				}
			}
			if !hasPrimary {
				product.Images[0].IsPrimary = true
			}
		}
	}

	if req.Attributes != nil {
		product.Attributes = req.Attributes
	}

	if req.SEO.Title != "" || req.SEO.Description != "" {
		product.SEO = req.SEO
	}

	// Save changes
	if err := s.repo.Product.Update(ctx, product); err != nil {
		s.logger.WithError(err).Error("Failed to update product")
		return nil, fmt.Errorf("failed to update product")
	}

	return product, nil
}

// ListProducts retrieves products with filtering and pagination
func (s *ProductService) ListProducts(ctx context.Context, filters ProductListFilters, offset, limit int) ([]models.Product, int64, error) {
	// Check if repository is available
	if s.repo.Product == nil {
		return nil, 0, fmt.Errorf("product catalog service is running in degraded mode - database unavailable")
	}
	
	filterMap := map[string]interface{}{
		"merchant_id": filters.MerchantID,
	}

	if filters.Status != "" {
		filterMap["status"] = filters.Status
	}

	if filters.Category != "" {
		filterMap["category"] = filters.Category
	}

	if filters.Tag != "" {
		filterMap["tag"] = filters.Tag
	}

	if filters.Search != "" {
		filterMap["search"] = filters.Search
	}

	if filters.PriceMin > 0 {
		filterMap["price_min"] = filters.PriceMin
	}

	if filters.PriceMax > 0 {
		filterMap["price_max"] = filters.PriceMax
	}

	if filters.Sort != "" {
		filterMap["sort"] = filters.Sort
	}

	return s.repo.Product.List(ctx, filterMap, offset, limit)
}

// PublishProduct publishes a product
func (s *ProductService) PublishProduct(ctx context.Context, productID, merchantID string) error {
	product, err := s.repo.Product.GetByID(ctx, productID)
	if err != nil {
		return ErrProductNotFound
	}

	if product.MerchantID != merchantID {
		return ErrUnauthorized
	}

	product.Status = "active"
	now := time.Now()
	product.PublishedAt = &now

	return s.repo.Product.Update(ctx, product)
}

// UnpublishProduct unpublishes a product
func (s *ProductService) UnpublishProduct(ctx context.Context, productID, merchantID string) error {
	product, err := s.repo.Product.GetByID(ctx, productID)
	if err != nil {
		return ErrProductNotFound
	}

	if product.MerchantID != merchantID {
		return ErrUnauthorized
	}

	product.Status = "inactive"
	product.PublishedAt = nil

	return s.repo.Product.Update(ctx, product)
}

// DeleteProduct soft deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, productID, merchantID string) error {
	product, err := s.repo.Product.GetByID(ctx, productID)
	if err != nil {
		return ErrProductNotFound
	}

	if product.MerchantID != merchantID {
		return ErrUnauthorized
	}

	return s.repo.Product.Delete(ctx, productID)
}

// GetProductsByCategory retrieves products by category
func (s *ProductService) GetProductsByCategory(ctx context.Context, merchantID, categorySlug string, offset, limit int) ([]models.Product, error) {
	return s.repo.Product.GetByCategory(ctx, merchantID, categorySlug, offset, limit)
}

// Helper methods

// generateSKU generates a SKU from product name
func (s *ProductService) generateSKU(name string) string {
	// Remove special characters and convert to uppercase
	sku := strings.ToUpper(utils.SlugifyString(name))
	sku = strings.ReplaceAll(sku, "-", "")

	// Add random suffix to ensure uniqueness
	suffix := utils.GenerateShortID()
	return fmt.Sprintf("%s-%s", sku, suffix)
}

// ensureUniqueSlug ensures the slug is unique for the merchant
func (s *ProductService) ensureUniqueSlug(ctx context.Context, merchantID, slug, excludeProductID string) string {
	slugCounter := 1
	originalSlug := slug

	for {
		existing, err := s.repo.Product.GetBySlug(ctx, merchantID, slug)
		if err != nil || (existing != nil && existing.ID.Hex() == excludeProductID) {
			break // Slug is available or belongs to the same product
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, slugCounter)
		slugCounter++
	}

	return slug
}

// Category Service Methods

// CreateCategoryRequest represents a category creation request
type CreateCategoryRequest struct {
	MerchantID  string               `json:"merchant_id" validate:"required"`
	Name        string               `json:"name" validate:"required"`
	Description string               `json:"description"`
	ParentID    string               `json:"parent_id"`
	Image       *models.ProductImage `json:"image"`
	SEO         models.ProductSEO    `json:"seo"`
}

// CreateCategory creates a new category
func (s *ProductService) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*models.Category, error) {
	// Generate slug from name
	slug := utils.SlugifyString(req.Name)

	// Ensure slug is unique
	slugCounter := 1
	originalSlug := slug
	for {
		_, err := s.repo.Category.GetBySlug(ctx, req.MerchantID, slug)
		if err != nil {
			break // Slug is available
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, slugCounter)
		slugCounter++
	}

	// Calculate level and path
	level := 0
	path := slug
	if req.ParentID != "" {
		parent, err := s.repo.Category.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, ErrCategoryNotFound
		}
		level = parent.Level + 1
		path = fmt.Sprintf("%s/%s", parent.Path, slug)
	}

	category := &models.Category{
		MerchantID:  req.MerchantID,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		Level:       level,
		Path:        path,
		Image:       req.Image,
		IsActive:    true,
		SEO:         req.SEO,
	}

	if err := s.repo.Category.Create(ctx, category); err != nil {
		s.logger.WithError(err).Error("Failed to create category")
		return nil, fmt.Errorf("failed to create category")
	}

	return category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *ProductService) GetCategoryBySlug(ctx context.Context, merchantID, slug string) (*models.Category, error) {
	return s.repo.Category.GetBySlug(ctx, merchantID, slug)
}

// GetCategoryChildren retrieves child categories for a parent category
func (s *ProductService) GetCategoryChildren(ctx context.Context, parentID string) ([]models.Category, error) {
	return s.repo.Category.GetChildren(ctx, parentID)
}

// GetCategoryParent retrieves the parent category for a category
func (s *ProductService) GetCategoryParent(ctx context.Context, parentID string) (*models.Category, error) {
	return s.repo.Category.GetParent(ctx, parentID)
}

// GetCategories retrieves categories for a merchant
func (s *ProductService) GetCategories(ctx context.Context, merchantID string) ([]models.Category, error) {
	return s.repo.Category.ListByMerchant(ctx, merchantID)
}
