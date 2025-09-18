#!/usr/bin/env node

/**
 * Unified Commerce OS - Sample Data Seeding Script
 * 
 * This script populates all services with comprehensive sample data for testing.
 * It uses GraphQL mutations to create realistic data across all microservices.
 */

const https = require('https');
const fs = require('fs');

// Configuration
const GATEWAY_URL = process.env.GATEWAY_URL || 'https://unified-commerce-gateway.onrender.com/graphql';
const ADMIN_EMAIL = 'admin@example.com';
const ADMIN_PASSWORD = 'Admin123!';

console.log('🌱 Unified Commerce OS - Sample Data Seeding');
console.log('=============================================');
console.log(`Gateway URL: ${GATEWAY_URL}`);
console.log('');

// GraphQL mutation helper
async function graphqlRequest(query, variables = {}) {
  const data = JSON.stringify({
    query,
    variables
  });

  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length
    }
  };

  return new Promise((resolve, reject) => {
    const req = https.request(GATEWAY_URL, options, (res) => {
      let body = '';
      res.on('data', chunk => body += chunk);
      res.on('end', () => {
        try {
          const response = JSON.parse(body);
          if (response.error || response.errors) {
            reject(new Error(JSON.stringify(response.errors || response.error)));
          } else {
            resolve(response.data);
          }
        } catch (e) {
          reject(e);
        }
      });
    });

    req.on('error', reject);
    req.write(data);
    req.end();
  });
}

// Sample data definitions
const sampleCategories = [
  { name: 'Electronics', description: 'Electronic devices and gadgets', slug: 'electronics' },
  { name: 'Fashion', description: 'Clothing and accessories', slug: 'fashion' },
  { name: 'Home & Kitchen', description: 'Home appliances and kitchen items', slug: 'home-kitchen' },
  { name: 'Sports & Outdoors', description: 'Sports equipment and outdoor gear', slug: 'sports-outdoors' },
  { name: 'Books', description: 'Books and educational materials', slug: 'books' }
];

const sampleProducts = [
  {
    name: 'Premium Wireless Headphones',
    description: 'High-quality noise-cancelling wireless headphones with premium sound quality and long battery life.',
    sku: 'PWH-001',
    price: 299.99,
    category: 'electronics',
    status: 'ACTIVE',
    tags: ['electronics', 'audio', 'wireless', 'premium']
  },
  {
    name: 'Smart Fitness Watch',
    description: 'Advanced fitness tracking watch with heart rate monitoring, GPS, and smartphone integration.',
    sku: 'SFW-002',
    price: 199.99,
    category: 'electronics',
    status: 'ACTIVE',
    tags: ['electronics', 'fitness', 'wearable', 'smart']
  },
  {
    name: 'Portable Bluetooth Speaker',
    description: 'Compact wireless speaker with powerful sound, waterproof design, and 12-hour battery life.',
    sku: 'PBS-003',
    price: 79.99,
    category: 'electronics',
    status: 'ACTIVE',
    tags: ['electronics', 'audio', 'portable', 'waterproof']
  },
  {
    name: 'Classic Cotton T-Shirt',
    description: 'Comfortable 100% cotton t-shirt available in multiple colors and sizes.',
    sku: 'CCT-004',
    price: 24.99,
    category: 'fashion',
    status: 'ACTIVE',
    tags: ['fashion', 'clothing', 'cotton', 'casual']
  },
  {
    name: 'Stainless Steel Water Bottle',
    description: 'Insulated stainless steel water bottle that keeps drinks cold for 24 hours or hot for 12 hours.',
    sku: 'SSWB-005',
    price: 34.99,
    category: 'home-kitchen',
    status: 'ACTIVE',
    tags: ['home-kitchen', 'drinkware', 'insulated', 'eco-friendly']
  },
  {
    name: 'Yoga Mat Pro',
    description: 'Professional-grade yoga mat with superior grip and cushioning. Perfect for all types of yoga practice.',
    sku: 'YMP-006',
    price: 89.99,
    category: 'sports-outdoors',
    status: 'ACTIVE',
    tags: ['sports', 'yoga', 'fitness', 'professional']
  },
  {
    name: 'JavaScript Programming Guide',
    description: 'Comprehensive guide to modern JavaScript programming with practical examples and best practices.',
    sku: 'JSG-007',
    price: 49.99,
    category: 'books',
    status: 'ACTIVE',
    tags: ['books', 'programming', 'javascript', 'education']
  },
  {
    name: 'Wireless Gaming Mouse',
    description: 'High-precision wireless gaming mouse with customizable RGB lighting and programmable buttons.',
    sku: 'WGM-008',
    price: 69.99,
    category: 'electronics',
    status: 'ACTIVE',
    tags: ['electronics', 'gaming', 'mouse', 'wireless']
  }
];

const sampleCustomers = [
  {
    firstName: 'John',
    lastName: 'Doe',
    email: 'john.doe@example.com',
    phone: '+1-555-0101'
  },
  {
    firstName: 'Jane',
    lastName: 'Smith',
    email: 'jane.smith@example.com',
    phone: '+1-555-0102'
  },
  {
    firstName: 'Bob',
    lastName: 'Johnson',
    email: 'bob.johnson@example.com',
    phone: '+1-555-0103'
  },
  {
    firstName: 'Alice',
    lastName: 'Williams',
    email: 'alice.williams@example.com',
    phone: '+1-555-0104'
  },
  {
    firstName: 'Charlie',
    lastName: 'Brown',
    email: 'charlie.brown@example.com',
    phone: '+1-555-0105'
  }
];

// Seeding functions
async function seedCategories() {
  console.log('📂 Seeding categories...');
  
  for (const category of sampleCategories) {
    try {
      const query = `
        mutation CreateCategory($input: CategoryInput!) {
          createCategory(input: $input) {
            id
            name
            slug
          }
        }
      `;
      
      const result = await graphqlRequest(query, { input: category });
      console.log(`  ✅ Created category: ${category.name}`);
    } catch (error) {
      console.log(`  ⚠️  Category ${category.name} might already exist or service unavailable`);
    }
  }
}

async function seedProducts() {
  console.log('📦 Seeding products...');
  
  for (const product of sampleProducts) {
    try {
      const query = `
        mutation CreateProduct($input: ProductInput!) {
          createProduct(input: $input) {
            id
            title
            sku
          }
        }
      `;
      
      const result = await graphqlRequest(query, { input: product });
      console.log(`  ✅ Created product: ${product.name}`);
    } catch (error) {
      console.log(`  ⚠️  Product ${product.name} might already exist or service unavailable`);
    }
  }
}

async function seedCustomers() {
  console.log('👥 Seeding customers...');
  
  for (const customer of sampleCustomers) {
    try {
      const query = `
        mutation CreateCustomer($input: CustomerInput!) {
          createCustomer(input: $input) {
            id
            firstName
            lastName
            email
          }
        }
      `;
      
      const result = await graphqlRequest(query, { input: customer });
      console.log(`  ✅ Created customer: ${customer.firstName} ${customer.lastName}`);
    } catch (error) {
      console.log(`  ⚠️  Customer ${customer.firstName} ${customer.lastName} might already exist or service unavailable`);
    }
  }
}

async function seedPromotions() {
  console.log('🎯 Seeding promotions...');
  
  const promotions = [
    {
      title: 'Welcome Discount',
      description: '10% off your first order',
      discountType: 'PERCENTAGE',
      discountValue: 10,
      code: 'WELCOME10',
      isActive: true,
      validFrom: new Date().toISOString(),
      validTo: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString() // 30 days from now
    },
    {
      title: 'Summer Sale',
      description: '$25 off orders over $100',
      discountType: 'FIXED',
      discountValue: 25,
      code: 'SUMMER25',
      minimumOrderValue: 100,
      isActive: true,
      validFrom: new Date().toISOString(),
      validTo: new Date(Date.now() + 60 * 24 * 60 * 60 * 1000).toISOString() // 60 days from now
    },
    {
      title: 'Electronics Deal',
      description: '15% off all electronics',
      discountType: 'PERCENTAGE',
      discountValue: 15,
      code: 'TECH15',
      isActive: true,
      validFrom: new Date().toISOString(),
      validTo: new Date(Date.now() + 45 * 24 * 60 * 60 * 1000).toISOString() // 45 days from now
    }
  ];
  
  for (const promotion of promotions) {
    try {
      const query = `
        mutation CreatePromotion($input: PromotionInput!) {
          createPromotion(input: $input) {
            id
            title
            code
          }
        }
      `;
      
      const result = await graphqlRequest(query, { input: promotion });
      console.log(`  ✅ Created promotion: ${promotion.title} (${promotion.code})`);
    } catch (error) {
      console.log(`  ⚠️  Promotion ${promotion.title} might already exist or service unavailable`);
    }
  }
}

// Main seeding function
async function seedAllData() {
  try {
    console.log('🚀 Starting comprehensive data seeding...\n');
    
    // Note: The order matters - some entities depend on others
    await seedCategories();
    console.log('');
    
    await seedProducts();
    console.log('');
    
    await seedCustomers();
    console.log('');
    
    await seedPromotions();
    console.log('');
    
    console.log('✨ Sample data seeding completed!');
    console.log('');
    console.log('📊 Summary:');
    console.log(`  - ${sampleCategories.length} categories`);
    console.log(`  - ${sampleProducts.length} products`);
    console.log(`  - ${sampleCustomers.length} customers`);
    console.log('  - 3 promotional campaigns');
    console.log('');
    console.log('🎉 Your Unified Commerce OS is now ready for testing!');
    console.log('');
    console.log('Next steps:');
    console.log('  1. Visit your storefront to see the products');
    console.log('  2. Use the admin panel to manage data');
    console.log('  3. Test GraphQL queries through the gateway');
    
  } catch (error) {
    console.error('❌ Error during seeding:', error.message);
    process.exit(1);
  }
}

// Run the seeding script
if (require.main === module) {
  seedAllData();
}

module.exports = {
  seedAllData,
  seedCategories,
  seedProducts,
  seedCustomers,
  seedPromotions
};