# Retail OS Sample Applications

This directory contains sample applications demonstrating how to integrate with the Retail OS API using different programming languages and frameworks.

## Available Samples

1. [JavaScript/Node.js](#javascriptnodejs-sample)
2. [Python](#python-sample)
3. [Java](#java-sample)
4. [Go](#go-sample)
5. [PHP](#php-sample)

## JavaScript/Node.js Sample

A complete Node.js application demonstrating product browsing, cart management, and checkout.

### Features
- Product catalog browsing
- Shopping cart functionality
- User authentication
- Order creation
- Payment processing

### Setup
```bash
cd javascript-nodejs
npm install
npm start
```

### Configuration
Create a `.env` file with your credentials:
```
CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
API_ENDPOINT=https://api.unified-commerce.com
```

## Python Sample

A Python Flask application showing how to build an e-commerce storefront.

### Features
- Product search and filtering
- User registration and login
- Order history
- Payment integration

### Setup
```bash
cd python-flask
pip install -r requirements.txt
python app.py
```

### Configuration
Create a `.env` file with your credentials:
```
CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
API_ENDPOINT=https://api.unified-commerce.com
```

## Java Sample

A Spring Boot application demonstrating enterprise integration patterns.

### Features
- Product inventory management
- Order processing workflows
- Reporting and analytics
- Webhook handling

### Setup
```bash
cd java-springboot
./mvnw spring-boot:run
```

### Configuration
Update `src/main/resources/application.properties`:
```
unified-commerce.client-id=your_client_id
unified-commerce.client-secret=your_client_secret
unified-commerce.api-endpoint=https://api.unified-commerce.com
```

## Go Sample

A Go application showing high-performance API integration.

### Features
- Real-time inventory updates
- Concurrent request handling
- Efficient data processing
- Health monitoring

### Setup
```bash
cd go-app
go run main.go
```

### Configuration
Create a `config.json` file:
```json
{
  "client_id": "your_client_id",
  "client_secret": "your_client_secret",
  "api_endpoint": "https://api.unified-commerce.com"
}
```

## PHP Sample

A PHP application demonstrating traditional web application integration.

### Features
- Product catalog display
- Shopping cart management
- User account functionality
- Order processing

### Setup
```bash
cd php-app
composer install
php -S localhost:8000
```

### Configuration
Create a `.env` file:
```
CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
API_ENDPOINT=https://api.unified-commerce.com
```

## Common Patterns

All samples demonstrate these common integration patterns:

### Authentication
```javascript
// JavaScript example
const client = new UnifiedCommerceClient({
  clientId: process.env.CLIENT_ID,
  clientSecret: process.env.CLIENT_SECRET
});
```

### Error Handling
```python
# Python example
try:
    product = client.products.get(product_id)
except UnifiedCommerceError as e:
    if e.code == 'PRODUCT_NOT_FOUND':
        # Handle missing product
        pass
    else:
        # Log and handle other errors
        logger.error(f"API error: {e.message}")
```

### Pagination
```java
// Java example
List<Product> allProducts = new ArrayList<>();
String cursor = null;
do {
    ProductListResult result = client.products().list(ProductListOptions.builder()
        .limit(50)
        .cursor(cursor)
        .build());
    allProducts.addAll(result.getProducts());
    cursor = result.getNextCursor();
} while (cursor != null);
```

### Rate Limiting
```go
// Go example
// SDK automatically handles rate limiting with exponential backoff
products, err := client.Products.List(ctx, &ProductsListOptions{
    Limit: 100,
})
if err != nil {
    // Handle error
}
```

## Running the Samples

### Prerequisites
1. A Retail OS developer account
2. API credentials (Client ID and Client Secret)
3. The appropriate runtime environment for each sample

### Getting API Credentials
1. Sign up at https://developer.retail-os.com
2. Create a new application in your dashboard
3. Note your Client ID and Client Secret
4. Configure the samples with your credentials

### Environment Variables
Most samples use environment variables for configuration:
- `CLIENT_ID` - Your application's Client ID
- `CLIENT_SECRET` - Your application's Client Secret
- `API_ENDPOINT` - The API endpoint (defaults to production)

## Contributing

We welcome contributions to our sample applications:

1. Fork the repository
2. Create a new branch for your feature
3. Add your sample application in a new directory
4. Include a README with setup instructions
5. Submit a pull request

## Support

For help with the samples:
- Check the [API Documentation](../api-docs/)
- Join our [Developer Community](https://community.retail-os.com)
- Email us at developers@retail-os.com