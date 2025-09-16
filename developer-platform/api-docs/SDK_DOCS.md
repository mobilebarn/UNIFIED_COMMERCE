# Retail OS SDK Documentation

## Overview

Retail OS provides official SDKs for several popular programming languages to make integration easier and more reliable. Each SDK handles authentication, request signing, error handling, and pagination automatically.

## Available SDKs

1. [JavaScript/TypeScript](#javascripttypescript-sdk)
2. [Python](#python-sdk)
3. [Java](#java-sdk)
4. [Go](#go-sdk)
5. [PHP](#php-sdk)

## JavaScript/TypeScript SDK

### Installation
```bash
npm install @retail-os/sdk
```

### Initialization
```javascript
import { RetailOSClient } from '@retail-os/sdk';

const client = new RetailOSClient({
  clientId: 'YOUR_CLIENT_ID',
  clientSecret: 'YOUR_CLIENT_SECRET',
  environment: 'production' // or 'development'
});
```

### Usage Examples

#### Fetch Products
```javascript
// Fetch a list of products
const products = await client.products.list({
  limit: 10,
  status: 'ACTIVE'
});

// Fetch a specific product
const product = await client.products.get('product_123');
```

#### Create a Product
```javascript
const newProduct = await client.products.create({
  merchantId: 'merchant_123',
  title: 'New Product',
  description: 'A great new product',
  status: 'ACTIVE'
});
```

#### Manage Orders
```javascript
// Create an order
const order = await client.orders.create({
  merchantId: 'merchant_123',
  customer: {
    email: 'customer@example.com',
    firstName: 'John',
    lastName: 'Doe'
  },
  // ... other order details
});

// Fetch orders
const orders = await client.orders.list({
  customerId: 'customer_123',
  limit: 20
});
```

#### Handle Payments
```javascript
// Process a payment
const payment = await client.payments.create({
  orderId: 'order_123',
  amount: 99.99,
  currency: 'USD',
  paymentMethodId: 'pm_456'
});
```

### Error Handling
```javascript
try {
  const product = await client.products.get('nonexistent_product');
} catch (error) {
  if (error.code === 'PRODUCT_NOT_FOUND') {
    console.log('Product not found');
  } else {
    console.error('An error occurred:', error.message);
  }
}
```

### Configuration Options
```javascript
const client = new RetailOSClient({
  clientId: 'YOUR_CLIENT_ID',
  clientSecret: 'YOUR_CLIENT_SECRET',
  environment: 'production', // 'production' or 'development'
  timeout: 30000, // Request timeout in milliseconds
  retries: 3, // Number of retry attempts
  retryDelay: 1000 // Delay between retries in milliseconds
});
```

## Python SDK

### Installation
```bash
pip install retail-os-sdk
```

### Initialization
```python
from retail_os import Client

client = Client(
    client_id='YOUR_CLIENT_ID',
    client_secret='YOUR_CLIENT_SECRET',
    environment='production'  # or 'development'
)
```

### Usage Examples

#### Fetch Products
```python
# Fetch a list of products
products = client.products.list(limit=10, status='ACTIVE')

# Fetch a specific product
product = client.products.get('product_123')
```

#### Create a Product
```python
new_product = client.products.create(
    merchant_id='merchant_123',
    title='New Product',
    description='A great new product',
    status='ACTIVE'
)
```

#### Manage Orders
```python
# Create an order
order = client.orders.create(
    merchant_id='merchant_123',
    customer={
        'email': 'customer@example.com',
        'first_name': 'John',
        'last_name': 'Doe'
    }
    # ... other order details
)

# Fetch orders
orders = client.orders.list(customer_id='customer_123', limit=20)
```

### Error Handling
```python
try:
    product = client.products.get('nonexistent_product')
except RetailOSError as e:
    if e.code == 'PRODUCT_NOT_FOUND':
        print('Product not found')
    else:
        print(f'An error occurred: {e.message}')
```

## Java SDK

### Installation
Add the following dependency to your `pom.xml`:
```xml
<dependency>
    <groupId>com.retailos</groupId>
    <artifactId>sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

### Initialization
```java
import com.retailos.sdk.Client;

Client client = new Client.Builder()
    .clientId("YOUR_CLIENT_ID")
    .clientSecret("YOUR_CLIENT_SECRET")
    .environment("production") // or "development"
    .build();
```

### Usage Examples

#### Fetch Products
```java
// Fetch a list of products
List<Product> products = client.products().list(10, ProductStatus.ACTIVE);

// Fetch a specific product
Product product = client.products().get("product_123");
```

#### Create a Product
```java
ProductInput input = ProductInput.builder()
    .merchantId("merchant_123")
    .title("New Product")
    .description("A great new product")
    .status(ProductStatus.ACTIVE)
    .build();

Product newProduct = client.products().create(input);
```

### Error Handling
```java
try {
    Product product = client.products().get("nonexistent_product");
} catch (RetailOSException e) {
    if ("PRODUCT_NOT_FOUND".equals(e.getCode())) {
        System.out.println("Product not found");
    } else {
        System.err.println("An error occurred: " + e.getMessage());
    }
}
```

## Go SDK

### Installation
```bash
go get github.com/retail-os/sdk-go
```

### Initialization
```go
import "github.com/retail-os/sdk-go/client"

cfg := client.Config{
    ClientID:     "YOUR_CLIENT_ID",
    ClientSecret: "YOUR_CLIENT_SECRET",
    Environment:  "production", // or "development"
}

client, err := client.New(cfg)
if err != nil {
    log.Fatal(err)
}
```

### Usage Examples

#### Fetch Products
```go
// Fetch a list of products
products, err := client.Products.List(context.Background(), &ProductsListOptions{
    Limit: 10,
    Status: "ACTIVE",
})
if err != nil {
    log.Fatal(err)
}

// Fetch a specific product
product, err := client.Products.Get(context.Background(), "product_123")
if err != nil {
    log.Fatal(err)
}
```

#### Create a Product
```go
input := &ProductCreateInput{
    MerchantID:  "merchant_123",
    Title:       "New Product",
    Description: "A great new product",
    Status:      "ACTIVE",
}

newProduct, err := client.Products.Create(context.Background(), input)
if err != nil {
    log.Fatal(err)
}
```

### Error Handling
```go
product, err := client.Products.Get(context.Background(), "nonexistent_product")
if err != nil {
    if retailos.IsErrorCode(err, "PRODUCT_NOT_FOUND") {
        fmt.Println("Product not found")
    } else {
        fmt.Printf("An error occurred: %v\n", err)
    }
}
```

## PHP SDK

### Installation
```bash
composer require retail-os/sdk
```

### Initialization
```php
use RetailOS\Client;

$client = new Client([
    'client_id' => 'YOUR_CLIENT_ID',
    'client_secret' => 'YOUR_CLIENT_SECRET',
    'environment' => 'production' // or 'development'
]);
```

### Usage Examples

#### Fetch Products
```php
// Fetch a list of products
$products = $client->products()->list([
    'limit' => 10,
    'status' => 'ACTIVE'
]);

// Fetch a specific product
$product = $client->products()->get('product_123');
```

#### Create a Product
```php
$newProduct = $client->products()->create([
    'merchant_id' => 'merchant_123',
    'title' => 'New Product',
    'description' => 'A great new product',
    'status' => 'ACTIVE'
]);
```

### Error Handling
```php
try {
    $product = $client->products()->get('nonexistent_product');
} catch (RetailOSException $e) {
    if ($e->getCode() === 'PRODUCT_NOT_FOUND') {
        echo 'Product not found';
    } else {
        echo 'An error occurred: ' . $e->getMessage();
    }
}
```

## SDK Features

All SDKs provide the following features:

### Authentication
- Automatic token management
- Refresh token handling
- Multiple authentication methods

### Error Handling
- Consistent error types across languages
- Detailed error information
- Automatic retry logic for transient errors

### Pagination
- Automatic pagination handling
- Iterator support where applicable
- Configurable page sizes

### Rate Limiting
- Automatic rate limit handling
- Exponential backoff
- Configurable retry settings

### Logging
- Configurable logging levels
- Request/response logging
- Performance metrics

### Security
- Secure credential storage
- Transport security (HTTPS)
- Input validation

## Best Practices

### Initialization
- Initialize the client once and reuse it
- Store credentials securely (environment variables, key management services)
- Configure appropriate timeouts and retry settings

### Resource Management
- Handle errors appropriately
- Use pagination for large data sets
- Close connections when done (if applicable)

### Performance
- Use connection pooling
- Cache frequently accessed data
- Batch requests when possible

### Security
- Never log credentials
- Validate input data
- Use the principle of least privilege

## Support

For SDK-related issues:
- Check the [SDK GitHub repositories](https://github.com/retail-os) for issues and updates
- File bug reports and feature requests
- Join our [Developer Community](https://community.retail-os.com) for support