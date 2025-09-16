# Retail OS JavaScript/Node.js Sample

This sample application demonstrates how to integrate with the Retail OS API using Node.js and Express.

## Features

- Product catalog browsing
- User authentication (mock implementation)
- Order history
- Simple web interface

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create a `.env` file with your API credentials:
   ```
   CLIENT_ID=your_client_id
   CLIENT_SECRET=your_client_secret
   API_ENDPOINT=http://localhost:4000/graphql
   ```

3. Start the application:
   ```bash
   npm start
   ```

4. Visit `http://localhost:3000` in your browser

## API Integration

The sample includes a `UnifiedCommerceClient` class that demonstrates:

### Authentication
```javascript
const client = new UnifiedCommerceClient();
await client.authenticate();
```

### Product Queries
```javascript
// Get a list of products
const products = await client.getProducts(10);

// Get a specific product
const product = await client.getProductById('product_123');
```

### User Queries
```javascript
// Get current user information
const user = await client.getCurrentUser();
```

### Order Queries
```javascript
// Get order history
const orders = await client.getOrders(10);
```

## Endpoints

- `GET /api/products` - Get products
- `GET /api/products/:id` - Get a specific product
- `GET /api/user` - Get current user
- `GET /api/orders` - Get orders

## Customization

To connect to a real Retail OS instance:

1. Update the `.env` file with your actual credentials
2. Modify the authentication method in `UnifiedCommerceClient`
3. Update the API endpoint if needed

## Error Handling

The client includes basic error handling:
```javascript
try {
  const products = await client.getProducts();
} catch (error) {
  console.error('Failed to fetch products:', error.message);
}
```

## Next Steps

1. Implement real authentication with OAuth 2.0
2. Add cart functionality
3. Implement checkout flow
4. Add webhook handlers for real-time updates
5. Enhance error handling and logging

## Support

For help with this sample:
- Check the [API Documentation](../../api-docs/)
- Join our [Developer Community](https://community.unified-commerce.com)
- Email us at developers@unified-commerce.com