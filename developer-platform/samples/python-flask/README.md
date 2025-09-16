# Retail OS Python/Flask Sample

This sample application demonstrates how to integrate with the Retail OS API using Python and Flask.

## Features

- Product catalog browsing
- User authentication (mock implementation)
- Order history
- Simple web interface

## Setup

1. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Create a `.env` file with your API credentials:
   ```
   CLIENT_ID=your_client_id
   CLIENT_SECRET=your_client_secret
   API_ENDPOINT=http://localhost:4000/graphql
   ```

3. Start the application:
   ```bash
   python app.py
   ```

4. Visit `http://localhost:3000` in your browser

## API Integration

The sample includes a `UnifiedCommerceClient` class that demonstrates:

### Authentication
```python
client = UnifiedCommerceClient()
client.authenticate()
```

### Product Queries
```python
# Get a list of products
products = client.get_products(10)

# Get a specific product
product = client.get_product_by_id('product_123')
```

### User Queries
```python
# Get current user information
user = client.get_current_user()
```

### Order Queries
```python
# Get order history
orders = client.get_orders(10)
```

## Endpoints

- `GET /api/products` - Get products
- `GET /api/products/<id>` - Get a specific product
- `GET /api/user` - Get current user
- `GET /api/orders` - Get orders

## Customization

To connect to a real Retail OS instance:

1. Update the `.env` file with your actual credentials
2. Modify the authentication method in `UnifiedCommerceClient`
3. Update the API endpoint if needed

## Error Handling

The client includes basic error handling:
```python
try:
    products = client.get_products()
except Exception as e:
    print(f'Failed to fetch products: {e}')
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