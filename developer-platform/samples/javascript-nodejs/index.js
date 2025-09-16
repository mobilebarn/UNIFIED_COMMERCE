// Sample Node.js application for Retail OS
require('dotenv').config();
const express = require('express');
const axios = require('axios');

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(express.json());
app.use(express.static('public'));

// Retail OS API client
class RetailOSClient {
  constructor() {
    this.apiEndpoint = process.env.API_ENDPOINT || 'http://localhost:4000/graphql';
    this.clientId = process.env.CLIENT_ID;
    this.clientSecret = process.env.CLIENT_SECRET;
    this.accessToken = null;
  }

  async authenticate() {
    // In a real implementation, you would authenticate with the API
    // For this sample, we'll just return a mock token
    this.accessToken = 'mock-access-token';
    return this.accessToken;
  }

  async makeRequest(query, variables = {}) {
    if (!this.accessToken) {
      await this.authenticate();
    }

    try {
      const response = await axios.post(this.apiEndpoint, {
        query,
        variables
      }, {
        headers: {
          'Authorization': `Bearer ${this.accessToken}`,
          'Content-Type': 'application/json'
        }
      });

      return response.data;
    } catch (error) {
      console.error('API request failed:', error.response?.data || error.message);
      throw error;
    }
  }

  // Product methods
  async getProducts(limit = 10) {
    const query = `
      query GetProducts($limit: Int) {
        products(filter: { limit: $limit }) {
          id
          title
          handle
          description
          featuredImage
          priceRange {
            minVariantPrice
          }
          tags
        }
      }
    `;

    const response = await this.makeRequest(query, { limit });
    return response.data.products;
  }

  async getProductById(id) {
    const query = `
      query GetProduct($id: ID!) {
        product(id: $id) {
          id
          title
          handle
          description
          featuredImage
          images {
            src
            altText
          }
          priceRange {
            minVariantPrice
          }
          variants {
            id
            title
            price
            inventoryQuantity
          }
          tags
        }
      }
    `;

    const response = await this.makeRequest(query, { id });
    return response.data.product;
  }

  // User methods
  async getCurrentUser() {
    const query = `
      query {
        currentUser {
          id
          email
          firstName
          lastName
        }
      }
    `;

    const response = await this.makeRequest(query);
    return response.data.currentUser;
  }

  // Order methods
  async getOrders(limit = 10) {
    const query = `
      query GetOrders($limit: Int) {
        orders(filter: { limit: $limit }) {
          id
          orderNumber
          status
          totalPrice
          createdAt
          lineItems {
            id
            name
            quantity
            price
          }
        }
      }
    `;

    const response = await this.makeRequest(query, { limit });
    return response.data.orders;
  }
}

// Initialize client
const client = new UnifiedCommerceClient();

// Routes
app.get('/api/products', async (req, res) => {
  try {
    const limit = parseInt(req.query.limit) || 10;
    const products = await client.getProducts(limit);
    res.json(products);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch products' });
  }
});

app.get('/api/products/:id', async (req, res) => {
  try {
    const product = await client.getProductById(req.params.id);
    res.json(product);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch product' });
  }
});

app.get('/api/user', async (req, res) => {
  try {
    const user = await client.getCurrentUser();
    res.json(user);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch user' });
  }
});

app.get('/api/orders', async (req, res) => {
  try {
    const limit = parseInt(req.query.limit) || 10;
    const orders = await client.getOrders(limit);
    res.json(orders);
  } catch (error) {
    res.status(500).json({ error: 'Failed to fetch orders' });
  }
});

// Serve a simple frontend
app.get('/', (req, res) => {
  res.send(`
    <!DOCTYPE html>
    <html>
    <head>
        <title>Retail OS Sample</title>
        <style>
            body { font-family: Arial, sans-serif; margin: 40px; }
            .product { border: 1px solid #ddd; padding: 20px; margin: 10px 0; }
            .product img { max-width: 200px; }
        </style>
    </head>
    <body>
        <h1>Retail OS Sample Application</h1>
        <p>This sample demonstrates how to integrate with the Retail OS API.</p>
        
        <h2>API Endpoints</h2>
        <ul>
            <li><a href="/api/products">/api/products</a> - Get products</li>
            <li><a href="/api/user">/api/user</a> - Get current user</li>
            <li><a href="/api/orders">/api/orders</a> - Get orders</li>
        </ul>
        
        <h2>Documentation</h2>
        <p>See <a href="https://docs.retail-os.com">API Documentation</a> for more details.</p>
    </body>
    </html>
  `);
});

// Start server
app.listen(PORT, () => {
  console.log(`Retail OS sample app listening at http://localhost:${PORT}`);
  console.log(`API endpoint: ${client.apiEndpoint}`);
});

module.exports = app;