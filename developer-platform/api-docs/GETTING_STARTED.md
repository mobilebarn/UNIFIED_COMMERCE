# Getting Started with the Retail OS API

Welcome to Retail OS! This guide will help you quickly get up and running with our API.

## Prerequisites

Before you begin, you'll need:
1. A developer account (sign up at https://developer.retail-os.com)
2. Basic knowledge of GraphQL
3. An HTTP client (curl, Postman, or your preferred tool)
4. A programming language of your choice

## Step 1: Create a Developer Account

1. Visit https://developer.retail-os.com
2. Click "Sign Up" and complete the registration form
3. Verify your email address
4. Log in to your developer dashboard

## Step 2: Create an Application

1. In your developer dashboard, click "Create Application"
2. Provide an application name and description
3. Select the appropriate permissions for your use case
4. Note your Client ID and Client Secret (keep these secure)

## Step 3: Obtain an Access Token

Use your credentials to obtain an access token:

```bash
curl -X POST \
  https://api.retail-os.com/auth/token \
  -H 'Content-Type: application/json' \
  -d '{
    "client_id": "YOUR_CLIENT_ID",
    "client_secret": "YOUR_CLIENT_SECRET",
    "grant_type": "client_credentials"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

## Step 4: Make Your First API Call

Use the access token to make your first GraphQL query:

```bash
curl -X POST \
  https://api.retail-os.com/graphql \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "query { currentUser { id email firstName lastName } }"
  }'
```

Response:
```json
{
  "data": {
    "currentUser": {
      "id": "123",
      "email": "developer@example.com",
      "firstName": "Developer",
      "lastName": "User"
    }
  }
}
```

## Step 5: Explore the API

Try querying products:

```graphql
query {
  products(filter: { limit: 5 }) {
    id
    title
    handle
    priceRange {
      minVariantPrice
    }
    featuredImage
  }
}
```

## Authentication Methods

### OAuth 2.0 (Recommended)
For applications that need to act on behalf of users:

```bash
curl -X POST \
  https://api.retail-os.com/auth/token \
  -H 'Content-Type: application/json' \
  -d '{
    "client_id": "YOUR_CLIENT_ID",
    "client_secret": "YOUR_CLIENT_SECRET",
    "grant_type": "authorization_code",
    "code": "AUTHORIZATION_CODE",
    "redirect_uri": "YOUR_REDIRECT_URI"
  }'
```

### API Keys
For server-to-server communication:

```bash
curl -X POST \
  https://api.retail-os.com/graphql \
  -H 'X-API-Key: YOUR_API_KEY' \
  -H 'Content-Type: application/json' \
  -d '{
    "query": "{ products(filter: { limit: 5 }) { id title } }"
  }'
```

## SDKs

We provide official SDKs to make integration easier:

### JavaScript/TypeScript
```bash
npm install @retail-os/sdk
```

```javascript
import { RetailOSClient } from '@retail-os/sdk';

const client = new RetailOSClient({
  clientId: 'YOUR_CLIENT_ID',
  clientSecret: 'YOUR_CLIENT_SECRET'
});

// Fetch products
const products = await client.products.list({ limit: 10 });
```

### Python
```bash
pip install retail-os-sdk
```

```python
from retail_os import Client

client = Client(
    client_id='YOUR_CLIENT_ID',
    client_secret='YOUR_CLIENT_SECRET'
)

# Fetch products
products = client.products.list(limit=10)
```

## Common Use Cases

### Creating a Product
```graphql
mutation {
  createProduct(input: {
    merchantId: "merchant_123"
    title: "New Product"
    description: "A great new product"
    status: ACTIVE
  }) {
    id
    title
    handle
  }
}
```

### Creating an Order
```graphql
mutation {
  createOrder(input: {
    merchantId: "merchant_123"
    customer: {
      email: "customer@example.com"
      firstName: "John"
      lastName: "Doe"
    }
    billingAddress: {
      firstName: "John"
      lastName: "Doe"
      street1: "123 Main St"
      city: "Anytown"
      state: "CA"
      country: "US"
      postalCode: "12345"
    }
    shippingAddress: {
      firstName: "John"
      lastName: "Doe"
      street1: "123 Main St"
      city: "Anytown"
      state: "CA"
      country: "US"
      postalCode: "12345"
    }
    currency: "USD"
  }) {
    id
    orderNumber
    status
  }
}
```

### Processing a Payment
```graphql
mutation {
  createPayment(input: {
    orderId: "order_123"
    merchantId: "merchant_123"
    amount: 99.99
    currency: "USD"
    paymentMethodId: "pm_456"
  }) {
    id
    amount
    status
  }
}
```

## Best Practices

1. **Use Pagination**: Always implement pagination for list queries
2. **Handle Errors Gracefully**: Check for errors in every response
3. **Cache When Appropriate**: Cache frequently accessed data to improve performance
4. **Respect Rate Limits**: Implement retry logic with exponential backoff
5. **Secure Your Credentials**: Never expose API keys or tokens in client-side code

## Next Steps

1. Review the [API Reference](API_REFERENCE.md) for detailed information about all endpoints
2. Check out the [SDK Documentation](SDK_DOCS.md) for language-specific guides
3. Join our [Developer Community](https://community.retail-os.com) for support and updates
4. Explore our [Sample Applications](../samples/) for complete implementation examples

## Support

If you need help, you can:
- Email us at developers@retail-os.com
- Check our [FAQ](FAQ.md)
- Join our [Developer Community](https://community.retail-os.com)
- File a support ticket through your developer dashboard