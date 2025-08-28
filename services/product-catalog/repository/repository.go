package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"unified-commerce/services/product-catalog/models"
	"unified-commerce/services/shared/database"
)

// ProductRepository handles product data operations
type ProductRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *database.MongoDB) *ProductRepository {
	return &ProductRepository{
		db:         db,
		collection: db.Collection("products"),
	}
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	var product models.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductRepository) GetBySKU(ctx context.Context, merchantID, sku string) (*models.Product, error) {
	var product models.Product
	filter := bson.M{
		"merchant_id": merchantID,
		"sku":         sku,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetBySlug retrieves a product by slug
func (r *ProductRepository) GetBySlug(ctx context.Context, merchantID, slug string) (*models.Product, error) {
	var product models.Product
	filter := bson.M{
		"merchant_id": merchantID,
		"slug":        slug,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()

	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete soft deletes a product by setting status to archived
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     "archived",
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// List retrieves products with filtering, sorting, and pagination
func (r *ProductRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Product, int64, error) {
	// Build filter
	filter := bson.M{}

	if merchantID, ok := filters["merchant_id"]; ok {
		filter["merchant_id"] = merchantID
	}

	if status, ok := filters["status"]; ok {
		filter["status"] = status
	}

	if category, ok := filters["category"]; ok {
		filter["categories"] = bson.M{"$in": []string{category.(string)}}
	}

	if tag, ok := filters["tag"]; ok {
		filter["tags"] = bson.M{"$in": []string{tag.(string)}}
	}

	if searchTerm, ok := filters["search"]; ok {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": searchTerm, "$options": "i"}},
			{"description": bson.M{"$regex": searchTerm, "$options": "i"}},
			{"sku": bson.M{"$regex": searchTerm, "$options": "i"}},
		}
	}

	if priceMin, ok := filters["price_min"]; ok {
		if filter["price"] == nil {
			filter["price"] = bson.M{}
		}
		filter["price"].(bson.M)["$gte"] = priceMin
	}

	if priceMax, ok := filters["price_max"]; ok {
		if filter["price"] == nil {
			filter["price"] = bson.M{}
		}
		filter["price"].(bson.M)["$lte"] = priceMax
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort options
	sortOptions := options.Find()
	if sortBy, ok := filters["sort"]; ok {
		switch sortBy {
		case "name_asc":
			sortOptions.SetSort(bson.D{{"name", 1}})
		case "name_desc":
			sortOptions.SetSort(bson.D{{"name", -1}})
		case "price_asc":
			sortOptions.SetSort(bson.D{{"price", 1}})
		case "price_desc":
			sortOptions.SetSort(bson.D{{"price", -1}})
		case "created_asc":
			sortOptions.SetSort(bson.D{{"created_at", 1}})
		case "created_desc":
			sortOptions.SetSort(bson.D{{"created_at", -1}})
		default:
			sortOptions.SetSort(bson.D{{"created_at", -1}})
		}
	} else {
		sortOptions.SetSort(bson.D{{"created_at", -1}})
	}

	// Apply pagination
	sortOptions.SetSkip(int64(offset))
	sortOptions.SetLimit(int64(limit))

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateStatus updates product status
func (r *ProductRepository) UpdateStatus(ctx context.Context, id, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetByCategory retrieves products by category
func (r *ProductRepository) GetByCategory(ctx context.Context, merchantID, categorySlug string, offset, limit int) ([]models.Product, error) {
	filter := bson.M{
		"merchant_id": merchantID,
		"categories":  bson.M{"$in": []string{categorySlug}},
		"status":      "active",
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"created_at", -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// CategoryRepository handles category data operations
type CategoryRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *database.MongoDB) *CategoryRepository {
	return &CategoryRepository{
		db:         db,
		collection: db.Collection("categories"),
	}
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return err
	}

	category.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*models.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID: %w", err)
	}

	var category models.Category
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, merchantID, slug string) (*models.Category, error) {
	var category models.Category
	filter := bson.M{
		"merchant_id": merchantID,
		"slug":        slug,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// ListByMerchant retrieves categories for a merchant
func (r *CategoryRepository) ListByMerchant(ctx context.Context, merchantID string) ([]models.Category, error) {
	filter := bson.M{
		"merchant_id": merchantID,
		"is_active":   true,
	}

	opts := options.Find().SetSort(bson.D{{"level", 1}, {"sort_order", 1}, {"name", 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// Update updates a category
func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	category.UpdatedAt = time.Now()

	filter := bson.M{"_id": category.ID}
	update := bson.M{"$set": category}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// CollectionRepository handles collection data operations
type CollectionRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

// NewCollectionRepository creates a new collection repository
func NewCollectionRepository(db *database.MongoDB) *CollectionRepository {
	return &CollectionRepository{
		db:         db,
		collection: db.Collection("collections"),
	}
}

// Create creates a new collection
func (r *CollectionRepository) Create(ctx context.Context, collection *models.Collection) error {
	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, collection)
	if err != nil {
		return err
	}

	collection.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetBySlug retrieves a collection by slug
func (r *CollectionRepository) GetBySlug(ctx context.Context, merchantID, slug string) (*models.Collection, error) {
	var collection models.Collection
	filter := bson.M{
		"merchant_id": merchantID,
		"slug":        slug,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// ListByMerchant retrieves collections for a merchant
func (r *CollectionRepository) ListByMerchant(ctx context.Context, merchantID string) ([]models.Collection, error) {
	filter := bson.M{
		"merchant_id": merchantID,
		"is_active":   true,
	}

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var collections []models.Collection
	if err := cursor.All(ctx, &collections); err != nil {
		return nil, err
	}

	return collections, nil
}

// Repository aggregates all repositories
type Repository struct {
	Product    *ProductRepository
	Category   *CategoryRepository
	Collection *CollectionRepository
}

// NewRepository creates a new repository with all sub-repositories
func NewRepository(db *database.MongoDB) *Repository {
	return &Repository{
		Product:    NewProductRepository(db),
		Category:   NewCategoryRepository(db),
		Collection: NewCollectionRepository(db),
	}
}

// CreateIndexes creates necessary database indexes for optimal performance
func (r *Repository) CreateIndexes(ctx context.Context) error {
	// Product indexes
	productIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"merchant_id", 1}, {"sku", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{"merchant_id", 1}, {"slug", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"merchant_id", 1}, {"status", 1}},
		},
		{
			Keys: bson.D{{"categories", 1}},
		},
		{
			Keys: bson.D{{"tags", 1}},
		},
		{
			Keys: bson.D{{"price", 1}},
		},
		{
			Keys: bson.D{{"created_at", -1}},
		},
		{
			Keys: bson.D{{"name", "text"}, {"description", "text"}, {"sku", "text"}},
		},
	}

	_, err := r.Product.collection.Indexes().CreateMany(ctx, productIndexes)
	if err != nil {
		return fmt.Errorf("failed to create product indexes: %w", err)
	}

	// Category indexes
	categoryIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"merchant_id", 1}, {"slug", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"merchant_id", 1}, {"parent_id", 1}},
		},
		{
			Keys: bson.D{{"level", 1}, {"sort_order", 1}},
		},
	}

	_, err = r.Category.collection.Indexes().CreateMany(ctx, categoryIndexes)
	if err != nil {
		return fmt.Errorf("failed to create category indexes: %w", err)
	}

	// Collection indexes
	collectionIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"merchant_id", 1}, {"slug", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"merchant_id", 1}, {"is_active", 1}},
		},
	}

	_, err = r.Collection.collection.Indexes().CreateMany(ctx, collectionIndexes)
	if err != nil {
		return fmt.Errorf("failed to create collection indexes: %w", err)
	}

	return nil
}
