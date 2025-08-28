// MongoDB initialization script for Product Catalog Service

// Switch to the product catalog database
db = db.getSiblingDB('product_catalog');

// Create a user for the Product Catalog service
db.createUser({
  user: 'catalog_user',
  pwd: 'catalog_pass',
  roles: [
    {
      role: 'readWrite',
      db: 'product_catalog'
    }
  ]
});

// Create collections with validation schemas
db.createCollection('products', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['sku', 'name', 'status', 'created_at', 'updated_at'],
      properties: {
        sku: {
          bsonType: 'string',
          description: 'Product SKU - required and unique'
        },
        name: {
          bsonType: 'string',
          description: 'Product name - required'
        },
        description: {
          bsonType: 'string',
          description: 'Product description'
        },
        status: {
          bsonType: 'string',
          enum: ['active', 'inactive', 'archived'],
          description: 'Product status - required'
        },
        variants: {
          bsonType: 'array',
          description: 'Product variants array',
          items: {
            bsonType: 'object',
            required: ['variant_id', 'price'],
            properties: {
              variant_id: {
                bsonType: 'string',
                description: 'Unique variant identifier'
              },
              price: {
                bsonType: 'number',
                minimum: 0,
                description: 'Variant price'
              },
              attributes: {
                bsonType: 'object',
                description: 'Flexible variant attributes (size, color, etc.)'
              }
            }
          }
        },
        categories: {
          bsonType: 'array',
          description: 'Product categories',
          items: {
            bsonType: 'string'
          }
        },
        tags: {
          bsonType: 'array',
          description: 'Product tags',
          items: {
            bsonType: 'string'
          }
        },
        images: {
          bsonType: 'array',
          description: 'Product images',
          items: {
            bsonType: 'object',
            properties: {
              url: {
                bsonType: 'string',
                description: 'Image URL'
              },
              alt: {
                bsonType: 'string',
                description: 'Alt text for accessibility'
              },
              is_primary: {
                bsonType: 'bool',
                description: 'Whether this is the primary image'
              }
            }
          }
        },
        seo: {
          bsonType: 'object',
          description: 'SEO optimization fields',
          properties: {
            title: {
              bsonType: 'string',
              description: 'SEO title'
            },
            description: {
              bsonType: 'string',
              description: 'SEO description'
            },
            keywords: {
              bsonType: 'array',
              items: {
                bsonType: 'string'
              }
            }
          }
        },
        created_at: {
          bsonType: 'date',
          description: 'Creation timestamp - required'
        },
        updated_at: {
          bsonType: 'date',
          description: 'Last update timestamp - required'
        }
      }
    }
  }
});

// Create unique index on SKU
db.products.createIndex({ 'sku': 1 }, { unique: true });

// Create indexes for common queries
db.products.createIndex({ 'status': 1 });
db.products.createIndex({ 'categories': 1 });
db.products.createIndex({ 'tags': 1 });
db.products.createIndex({ 'name': 'text', 'description': 'text' });

// Create collection for product categories
db.createCollection('categories', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'slug', 'created_at'],
      properties: {
        name: {
          bsonType: 'string',
          description: 'Category name - required'
        },
        slug: {
          bsonType: 'string',
          description: 'URL-friendly category slug - required and unique'
        },
        description: {
          bsonType: 'string',
          description: 'Category description'
        },
        parent_id: {
          bsonType: 'string',
          description: 'Parent category ID for hierarchical categories'
        },
        sort_order: {
          bsonType: 'int',
          description: 'Sort order for category display'
        },
        is_active: {
          bsonType: 'bool',
          description: 'Whether category is active'
        },
        created_at: {
          bsonType: 'date',
          description: 'Creation timestamp - required'
        }
      }
    }
  }
});

// Create unique index on category slug
db.categories.createIndex({ 'slug': 1 }, { unique: true });

print('Product Catalog database initialized successfully!');
print('Collections created: products, categories');
print('User created: catalog_user');