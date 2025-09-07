package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in the catalog
type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MerchantID  string             `json:"merchant_id" bson:"merchant_id"`
	SKU         string             `json:"sku" bson:"sku"`
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status"`         // "active", "inactive", "archived"
	Type        string             `json:"type" bson:"type"`             // "simple", "variable", "grouped", "digital"
	Visibility  string             `json:"visibility" bson:"visibility"` // "public", "private", "hidden"

	// Pricing
	Price          float64 `json:"price" bson:"price"`
	CompareAtPrice float64 `json:"compare_at_price" bson:"compare_at_price"`
	CostPrice      float64 `json:"cost_price" bson:"cost_price"`
	Currency       string  `json:"currency" bson:"currency"`
	TaxClass       string  `json:"tax_class" bson:"tax_class"`

	// Physical properties
	Weight           float64           `json:"weight" bson:"weight"`
	WeightUnit       string            `json:"weight_unit" bson:"weight_unit"`
	Dimensions       ProductDimensions `json:"dimensions" bson:"dimensions"`
	RequiresShipping bool              `json:"requires_shipping" bson:"requires_shipping"`

	// Inventory tracking
	TrackInventory  bool   `json:"track_inventory" bson:"track_inventory"`
	InventoryPolicy string `json:"inventory_policy" bson:"inventory_policy"` // "deny", "continue"

	// Categories and tags
	Categories []string `json:"categories" bson:"categories"`
	Tags       []string `json:"tags" bson:"tags"`

	// Images and media
	Images []ProductImage `json:"images" bson:"images"`
	Videos []ProductVideo `json:"videos" bson:"videos"`

	// Variants (for variable products)
	Variants []ProductVariant `json:"variants" bson:"variants"`
	Options  []ProductOption  `json:"options" bson:"options"` // Size, Color, etc.

	// SEO and marketing
	SEO ProductSEO `json:"seo" bson:"seo"`

	// Custom attributes (flexible schema)
	Attributes map[string]interface{} `json:"attributes" bson:"attributes"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" bson:"updated_at"`
	PublishedAt *time.Time `json:"published_at" bson:"published_at"`
}

// IsEntity implements the gqlgen federation Entity interface
func (p *Product) IsEntity() {}

// ProductDimensions represents product physical dimensions
type ProductDimensions struct {
	Length float64 `json:"length" bson:"length"`
	Width  float64 `json:"width" bson:"width"`
	Height float64 `json:"height" bson:"height"`
	Unit   string  `json:"unit" bson:"unit"` // "cm", "in", etc.
}

// ProductImage represents product images
type ProductImage struct {
	ID        string    `json:"id" bson:"id"`
	URL       string    `json:"url" bson:"url"`
	Alt       string    `json:"alt" bson:"alt"`
	IsPrimary bool      `json:"is_primary" bson:"is_primary"`
	SortOrder int       `json:"sort_order" bson:"sort_order"`
	Width     int       `json:"width" bson:"width"`
	Height    int       `json:"height" bson:"height"`
	FileSize  int64     `json:"file_size" bson:"file_size"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// ProductVideo represents product videos
type ProductVideo struct {
	ID           string    `json:"id" bson:"id"`
	URL          string    `json:"url" bson:"url"`
	ThumbnailURL string    `json:"thumbnail_url" bson:"thumbnail_url"`
	Title        string    `json:"title" bson:"title"`
	Duration     int       `json:"duration" bson:"duration"` // in seconds
	SortOrder    int       `json:"sort_order" bson:"sort_order"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

// ProductVariant represents product variants
type ProductVariant struct {
	ID             string                 `json:"id" bson:"id"`
	SKU            string                 `json:"sku" bson:"sku"`
	Barcode        string                 `json:"barcode" bson:"barcode"`
	Price          float64                `json:"price" bson:"price"`
	CompareAtPrice float64                `json:"compare_at_price" bson:"compare_at_price"`
	CostPrice      float64                `json:"cost_price" bson:"cost_price"`
	Weight         float64                `json:"weight" bson:"weight"`
	Image          *ProductImage          `json:"image" bson:"image"`
	Position       int                    `json:"position" bson:"position"`
	OptionValues   map[string]string      `json:"option_values" bson:"option_values"` // {"Color": "Red", "Size": "Large"}
	Attributes     map[string]interface{} `json:"attributes" bson:"attributes"`
	IsDefault      bool                   `json:"is_default" bson:"is_default"`
	CreatedAt      time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" bson:"updated_at"`
}

// IsEntity implements the gqlgen federation Entity interface
func (p *ProductVariant) IsEntity() {}

// ProductOption represents product options (like Color, Size)
type ProductOption struct {
	ID       string   `json:"id" bson:"id"`
	Name     string   `json:"name" bson:"name"`
	Position int      `json:"position" bson:"position"`
	Values   []string `json:"values" bson:"values"`
}

// ProductSEO represents SEO-related fields
type ProductSEO struct {
	Title          string   `json:"title" bson:"title"`
	Description    string   `json:"description" bson:"description"`
	Keywords       []string `json:"keywords" bson:"keywords"`
	CanonicalURL   string   `json:"canonical_url" bson:"canonical_url"`
	OpenGraphTitle string   `json:"og_title" bson:"og_title"`
	OpenGraphDesc  string   `json:"og_description" bson:"og_description"`
	OpenGraphImage string   `json:"og_image" bson:"og_image"`
}

// Category represents product categories
type Category struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MerchantID  string             `json:"merchant_id" bson:"merchant_id"`
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description" bson:"description"`
	ParentID    string             `json:"parent_id" bson:"parent_id"`
	Level       int                `json:"level" bson:"level"`
	Path        string             `json:"path" bson:"path"` // Hierarchical path like "electronics/computers/laptops"
	Image       *ProductImage      `json:"image" bson:"image"`
	SortOrder   int                `json:"sort_order" bson:"sort_order"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	SEO         ProductSEO         `json:"seo" bson:"seo"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsEntity implements the gqlgen federation Entity interface
func (c *Category) IsEntity() {}

// Collection represents curated product collections
type Collection struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MerchantID  string             `json:"merchant_id" bson:"merchant_id"`
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description" bson:"description"`
	Type        string             `json:"type" bson:"type"` // "manual", "auto", "smart"
	Rules       []CollectionRule   `json:"rules" bson:"rules"`
	ProductIDs  []string           `json:"product_ids" bson:"product_ids"`
	Image       *ProductImage      `json:"image" bson:"image"`
	SortOrder   string             `json:"sort_order" bson:"sort_order"` // "manual", "created", "price_asc", "price_desc", "name_asc", "name_desc"
	IsActive    bool               `json:"is_active" bson:"is_active"`
	IsPublished bool               `json:"is_published" bson:"is_published"`
	PublishedAt *time.Time         `json:"published_at" bson:"published_at"`
	SEO         ProductSEO         `json:"seo" bson:"seo"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsEntity implements the gqlgen federation Entity interface
func (c *Collection) IsEntity() {}

// CollectionRule represents automated collection rules
type CollectionRule struct {
	Field    string      `json:"field" bson:"field"`       // "price", "category", "tag", "title", etc.
	Operator string      `json:"operator" bson:"operator"` // "equals", "not_equals", "contains", "gt", "lt", etc.
	Value    interface{} `json:"value" bson:"value"`
}

// Brand represents product brands
type Brand struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MerchantID  string             `json:"merchant_id" bson:"merchant_id"`
	Name        string             `json:"name" bson:"name"`
	Slug        string             `json:"slug" bson:"slug"`
	Description string             `json:"description" bson:"description"`
	Logo        *ProductImage      `json:"logo" bson:"logo"`
	Website     string             `json:"website" bson:"website"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	SEO         ProductSEO         `json:"seo" bson:"seo"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsEntity implements the gqlgen federation Entity interface
func (b *Brand) IsEntity() {}

// ProductReview represents customer product reviews
type ProductReview struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductID    string             `json:"product_id" bson:"product_id"`
	CustomerID   string             `json:"customer_id" bson:"customer_id"`
	OrderID      string             `json:"order_id" bson:"order_id"`
	Rating       int                `json:"rating" bson:"rating"` // 1-5
	Title        string             `json:"title" bson:"title"`
	Content      string             `json:"content" bson:"content"`
	Images       []ProductImage     `json:"images" bson:"images"`
	Status       string             `json:"status" bson:"status"` // "pending", "approved", "rejected"
	IsVerified   bool               `json:"is_verified" bson:"is_verified"`
	HelpfulCount int                `json:"helpful_count" bson:"helpful_count"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

// ProductPriceRange represents the price range for a product
type ProductPriceRange struct {
	MinVariantPrice float64 `json:"min_variant_price" bson:"min_variant_price"`
	MaxVariantPrice float64 `json:"max_variant_price" bson:"max_variant_price"`
}

// Helper methods for Product model

// GetPrimaryImage returns the primary product image
func (p *Product) GetPrimaryImage() *ProductImage {
	for _, image := range p.Images {
		if image.IsPrimary {
			return &image
		}
	}
	if len(p.Images) > 0 {
		return &p.Images[0]
	}
	return nil
}

// GetDefaultVariant returns the default product variant
func (p *Product) GetDefaultVariant() *ProductVariant {
	for _, variant := range p.Variants {
		if variant.IsDefault {
			return &variant
		}
	}
	if len(p.Variants) > 0 {
		return &p.Variants[0]
	}
	return nil
}

// IsPublished checks if the product is published
func (p *Product) IsPublished() bool {
	return p.Status == "active" && p.PublishedAt != nil
}

// IsInStock checks if the product is in stock (basic check)
func (p *Product) IsInStock() bool {
	// This is a simplified check - in reality, you'd check against inventory service
	return p.Status == "active"
}

// GetDisplayPrice returns the price to display (considering compare_at_price)
func (p *Product) GetDisplayPrice() float64 {
	if p.CompareAtPrice > 0 && p.CompareAtPrice > p.Price {
		return p.Price // On sale
	}
	return p.Price
}

// IsOnSale checks if the product is on sale
func (p *Product) IsOnSale() bool {
	return p.CompareAtPrice > 0 && p.CompareAtPrice > p.Price
}

// GetSavingsAmount returns the savings amount if on sale
func (p *Product) GetSavingsAmount() float64 {
	if p.IsOnSale() {
		return p.CompareAtPrice - p.Price
	}
	return 0
}

// GetSavingsPercentage returns the savings percentage if on sale
func (p *Product) GetSavingsPercentage() float64 {
	if p.IsOnSale() {
		return ((p.CompareAtPrice - p.Price) / p.CompareAtPrice) * 100
	}
	return 0
}

// AddTag adds a tag to the product if it doesn't exist
func (p *Product) AddTag(tag string) {
	for _, existingTag := range p.Tags {
		if existingTag == tag {
			return
		}
	}
	p.Tags = append(p.Tags, tag)
}

// RemoveTag removes a tag from the product
func (p *Product) RemoveTag(tag string) {
	for i, existingTag := range p.Tags {
		if existingTag == tag {
			p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
			return
		}
	}
}

// HasTag checks if the product has a specific tag
func (p *Product) HasTag(tag string) bool {
	for _, existingTag := range p.Tags {
		if existingTag == tag {
			return true
		}
	}
	return false
}

// SearchSuggestion represents a search suggestion for autocomplete
type SearchSuggestion struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Type     string  `json:"type"` // "PRODUCT", "CATEGORY", "BRAND"
	ImageURL *string `json:"image_url,omitempty"`
}
