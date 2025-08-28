package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"unified-commerce/services/product-catalog/service"
	httputil "unified-commerce/services/shared/http"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	service   *service.ProductService
	logger    *logger.Logger
	validator *validator.Validate
}

// NewProductHandler creates a new product handler
func NewProductHandler(service *service.ProductService, logger *logger.Logger) *ProductHandler {
	return &ProductHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// RegisterRoutes registers all product catalog routes
func (h *ProductHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required for storefront)
	public := v1.Group("/catalog")
	{
		public.GET("/products", h.GetPublicProducts)
		public.GET("/products/:slug", h.GetProductBySlug)
		public.GET("/categories", h.GetPublicCategories)
		public.GET("/categories/:slug", h.GetCategoryBySlug)
		public.GET("/categories/:slug/products", h.GetProductsByCategory)
		public.GET("/collections", h.GetPublicCollections)
		public.GET("/collections/:slug", h.GetCollectionBySlug)
		public.GET("/search", h.SearchProducts)
	}

	// Authentication required for merchant routes
	authConfig := middleware.DefaultAuthConfig()
	v1.Use(middleware.JWTAuth(authConfig))

	// Product management routes
	products := v1.Group("/products")
	{
		products.POST("/", h.CreateProduct)
		products.GET("/", h.ListProducts)
		products.GET("/:id", h.GetProduct)
		products.PUT("/:id", h.UpdateProduct)
		products.DELETE("/:id", h.DeleteProduct)
		products.POST("/:id/publish", h.PublishProduct)
		products.POST("/:id/unpublish", h.UnpublishProduct)

		// Product variants
		products.GET("/:id/variants", h.GetProductVariants)
		products.POST("/:id/variants", h.CreateProductVariant)
		products.PUT("/:id/variants/:variantId", h.UpdateProductVariant)
		products.DELETE("/:id/variants/:variantId", h.DeleteProductVariant)

		// Product images
		products.GET("/:id/images", h.GetProductImages)
		products.POST("/:id/images", h.UploadProductImage)
		products.PUT("/:id/images/:imageId", h.UpdateProductImage)
		products.DELETE("/:id/images/:imageId", h.DeleteProductImage)
	}

	// Category management routes
	categories := v1.Group("/categories")
	{
		categories.POST("/", h.CreateCategory)
		categories.GET("/", h.GetCategories)
		categories.GET("/:id", h.GetCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}

	// Collection management routes
	collections := v1.Group("/collections")
	{
		collections.POST("/", h.CreateCollection)
		collections.GET("/", h.GetCollections)
		collections.GET("/:id", h.GetCollection)
		collections.PUT("/:id", h.UpdateCollection)
		collections.DELETE("/:id", h.DeleteCollection)
	}

	// Admin routes (require admin role)
	admin := v1.Group("/admin")
	admin.Use(middleware.RequireRole("admin", "super_admin"))
	{
		admin.GET("/products", h.AdminListProducts)
		admin.GET("/products/:id", h.AdminGetProduct)
		admin.POST("/products/:id/approve", h.ApproveProduct)
		admin.POST("/products/:id/reject", h.RejectProduct)
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body", map[string]interface{}{"error": err.Error()})
		return
	}

	// Get merchant ID from authenticated user context
	// This would typically come from the user's merchant association
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}
	req.MerchantID = merchantID

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrProductExists:
			httputil.Conflict(c, "Product with this SKU already exists")
		default:
			h.logger.WithError(err).Error("Failed to create product")
			httputil.InternalServerError(c, "Failed to create product")
		}
		return
	}

	httputil.Created(c, product, "Product created successfully")
}

// GetProduct retrieves a specific product
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		httputil.BadRequest(c, "Product ID is required")
		return
	}

	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), productID, merchantID)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			httputil.NotFound(c, "Product not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			h.logger.WithError(err).Error("Failed to get product")
			httputil.InternalServerError(c, "Failed to retrieve product")
		}
		return
	}

	httputil.Success(c, product)
}

// ListProducts retrieves products with filtering and pagination
func (h *ProductHandler) ListProducts(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	pagination := httputil.GetPaginationParams(c)

	filters := service.ProductListFilters{
		MerchantID: merchantID,
		Status:     c.Query("status"),
		Category:   c.Query("category"),
		Tag:        c.Query("tag"),
		Search:     c.Query("search"),
		Sort:       c.Query("sort"),
	}

	if priceMin := c.Query("price_min"); priceMin != "" {
		if price, err := strconv.ParseFloat(priceMin, 64); err == nil {
			filters.PriceMin = price
		}
	}

	if priceMax := c.Query("price_max"); priceMax != "" {
		if price, err := strconv.ParseFloat(priceMax, 64); err == nil {
			filters.PriceMax = price
		}
	}

	products, total, err := h.service.ListProducts(c.Request.Context(), filters, pagination.CalculateOffset(), pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list products")
		httputil.InternalServerError(c, "Failed to retrieve products")
		return
	}

	meta := &httputil.MetaInfo{
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total),
	}

	httputil.SuccessWithMeta(c, products, meta)
}

// UpdateProduct handles product updates
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		httputil.BadRequest(c, "Product ID is required")
		return
	}

	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), productID, merchantID, &req)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			httputil.NotFound(c, "Product not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			h.logger.WithError(err).Error("Failed to update product")
			httputil.InternalServerError(c, "Failed to update product")
		}
		return
	}

	httputil.Success(c, product, "Product updated successfully")
}

// DeleteProduct handles product deletion
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		httputil.BadRequest(c, "Product ID is required")
		return
	}

	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	err := h.service.DeleteProduct(c.Request.Context(), productID, merchantID)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			httputil.NotFound(c, "Product not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			h.logger.WithError(err).Error("Failed to delete product")
			httputil.InternalServerError(c, "Failed to delete product")
		}
		return
	}

	httputil.Success(c, nil, "Product deleted successfully")
}

// PublishProduct handles product publishing
func (h *ProductHandler) PublishProduct(c *gin.Context) {
	productID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	err := h.service.PublishProduct(c.Request.Context(), productID, merchantID)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			httputil.NotFound(c, "Product not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			httputil.InternalServerError(c, "Failed to publish product")
		}
		return
	}

	httputil.Success(c, nil, "Product published successfully")
}

// UnpublishProduct handles product unpublishing
func (h *ProductHandler) UnpublishProduct(c *gin.Context) {
	productID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	err := h.service.UnpublishProduct(c.Request.Context(), productID, merchantID)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			httputil.NotFound(c, "Product not found")
		case service.ErrUnauthorized:
			httputil.Forbidden(c, "Access denied")
		default:
			httputil.InternalServerError(c, "Failed to unpublish product")
		}
		return
	}

	httputil.Success(c, nil, "Product unpublished successfully")
}

// Public API endpoints (for storefront)

// GetPublicProducts retrieves published products for public access
func (h *ProductHandler) GetPublicProducts(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	pagination := httputil.GetPaginationParams(c)

	filters := service.ProductListFilters{
		MerchantID: merchantID,
		Status:     "active", // Only published products
		Category:   c.Query("category"),
		Tag:        c.Query("tag"),
		Search:     c.Query("search"),
		Sort:       c.Query("sort"),
	}

	if priceMin := c.Query("price_min"); priceMin != "" {
		if price, err := strconv.ParseFloat(priceMin, 64); err == nil {
			filters.PriceMin = price
		}
	}

	if priceMax := c.Query("price_max"); priceMax != "" {
		if price, err := strconv.ParseFloat(priceMax, 64); err == nil {
			filters.PriceMax = price
		}
	}

	products, total, err := h.service.ListProducts(c.Request.Context(), filters, pagination.CalculateOffset(), pagination.PerPage)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list public products")
		httputil.InternalServerError(c, "Failed to retrieve products")
		return
	}

	meta := &httputil.MetaInfo{
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		Total:      total,
		TotalPages: pagination.CalculateTotalPages(total),
	}

	httputil.SuccessWithMeta(c, products, meta)
}

// GetProductBySlug retrieves a product by slug (public access)
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		httputil.BadRequest(c, "Product slug is required")
		return
	}

	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	product, err := h.service.GetProductBySlug(c.Request.Context(), merchantID, slug)
	if err != nil {
		switch err {
		case service.ErrProductNotFound:
			httputil.NotFound(c, "Product not found")
		default:
			h.logger.WithError(err).Error("Failed to get product by slug")
			httputil.InternalServerError(c, "Failed to retrieve product")
		}
		return
	}

	httputil.Success(c, product)
}

// Category handlers

// CreateCategory handles category creation
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "Invalid request body")
		return
	}

	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}
	req.MerchantID = merchantID

	if err := h.validator.Struct(&req); err != nil {
		httputil.ValidationError(c, map[string]interface{}{"validation": err.Error()})
		return
	}

	category, err := h.service.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create category")
		httputil.InternalServerError(c, "Failed to create category")
		return
	}

	httputil.Created(c, category, "Category created successfully")
}

// GetCategories retrieves categories for a merchant
func (h *ProductHandler) GetCategories(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		httputil.BadRequest(c, "Merchant ID is required")
		return
	}

	categories, err := h.service.GetCategories(c.Request.Context(), merchantID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get categories")
		httputil.InternalServerError(c, "Failed to retrieve categories")
		return
	}

	httputil.Success(c, categories)
}

// GetPublicCategories retrieves categories for public access
func (h *ProductHandler) GetPublicCategories(c *gin.Context) {
	h.GetCategories(c) // Same as private for now
}

// GetCategoryBySlug retrieves a category by slug
func (h *ProductHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	merchantID := c.GetHeader("X-Merchant-ID")

	category, err := h.service.GetCategoryBySlug(c.Request.Context(), merchantID, slug)
	if err != nil {
		httputil.NotFound(c, "Category not found")
		return
	}

	httputil.Success(c, category)
}

// GetProductsByCategory retrieves products by category
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	slug := c.Param("slug")
	merchantID := c.GetHeader("X-Merchant-ID")
	pagination := httputil.GetPaginationParams(c)

	products, err := h.service.GetProductsByCategory(c.Request.Context(), merchantID, slug, pagination.CalculateOffset(), pagination.PerPage)
	if err != nil {
		httputil.InternalServerError(c, "Failed to retrieve products")
		return
	}

	httputil.Success(c, products)
}

// Placeholder implementations for remaining endpoints

func (h *ProductHandler) GetProductVariants(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get product variants endpoint - to be implemented"})
}

func (h *ProductHandler) CreateProductVariant(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Create product variant endpoint - to be implemented"})
}

func (h *ProductHandler) UpdateProductVariant(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Update product variant endpoint - to be implemented"})
}

func (h *ProductHandler) DeleteProductVariant(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Delete product variant endpoint - to be implemented"})
}

func (h *ProductHandler) GetProductImages(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get product images endpoint - to be implemented"})
}

func (h *ProductHandler) UploadProductImage(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Upload product image endpoint - to be implemented"})
}

func (h *ProductHandler) UpdateProductImage(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Update product image endpoint - to be implemented"})
}

func (h *ProductHandler) DeleteProductImage(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Delete product image endpoint - to be implemented"})
}

func (h *ProductHandler) GetCategory(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get category endpoint - to be implemented"})
}

func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Update category endpoint - to be implemented"})
}

func (h *ProductHandler) DeleteCategory(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Delete category endpoint - to be implemented"})
}

func (h *ProductHandler) CreateCollection(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Create collection endpoint - to be implemented"})
}

func (h *ProductHandler) GetCollections(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get collections endpoint - to be implemented"})
}

func (h *ProductHandler) GetCollection(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get collection endpoint - to be implemented"})
}

func (h *ProductHandler) UpdateCollection(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Update collection endpoint - to be implemented"})
}

func (h *ProductHandler) DeleteCollection(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Delete collection endpoint - to be implemented"})
}

func (h *ProductHandler) GetPublicCollections(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get public collections endpoint - to be implemented"})
}

func (h *ProductHandler) GetCollectionBySlug(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Get collection by slug endpoint - to be implemented"})
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Search products endpoint - to be implemented"})
}

func (h *ProductHandler) AdminListProducts(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Admin list products endpoint - to be implemented"})
}

func (h *ProductHandler) AdminGetProduct(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Admin get product endpoint - to be implemented"})
}

func (h *ProductHandler) ApproveProduct(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Approve product endpoint - to be implemented"})
}

func (h *ProductHandler) RejectProduct(c *gin.Context) {
	httputil.Success(c, gin.H{"message": "Reject product endpoint - to be implemented"})
}
