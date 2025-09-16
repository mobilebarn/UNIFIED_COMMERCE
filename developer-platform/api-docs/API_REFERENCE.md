# Retail OS API Reference

## Introduction

This document provides detailed reference information for the Retail OS GraphQL API. The API is organized around core e-commerce entities and operations, with all services accessible through a single GraphQL endpoint.

## GraphQL Schema Overview

Retail OS uses GraphQL Federation to combine multiple microservices into a single, cohesive API. The schema includes the following main types:

### Core Types

#### User
Represents a customer or admin user in the system.

**Fields:**
- `id: ID!` - Unique identifier
- `email: String!` - User's email address
- `username: String!` - Username
- `firstName: String!` - First name
- `lastName: String!` - Last name
- `isActive: Boolean!` - Whether the user account is active
- `roles: [Role!]!` - Roles assigned to the user
- `createdAt: String!` - Creation timestamp
- `updatedAt: String` - Last update timestamp

#### Product
Represents a product in the catalog.

**Fields:**
- `id: ID!` - Unique identifier
- `merchantId: ID!` - ID of the merchant who owns this product
- `title: String!` - Product title
- `handle: String!` - URL-friendly product handle
- `description: String` - Product description
- `vendor: String` - Product vendor
- `productType: String` - Product type/category
- `status: ProductStatus!` - Product status (ACTIVE, ARCHIVED, DRAFT)
- `seoTitle: String` - SEO title
- `seoDescription: String` - SEO description
- `publishedAt: String` - Publication timestamp
- `publishedScope: PublishedScope!` - Publication scope
- `featuredImage: String` - URL to featured image
- `images: [ProductImage!]!` - Product images
- `priceRange: ProductPriceRange!` - Price range for variants
- `compareAtPriceRange: ProductPriceRange` - Compare at price range
- `totalInventory: Int!` - Total inventory count
- `hasOnlyDefaultVariant: Boolean!` - Whether product has only default variant
- `requiresSellingPlan: Boolean!` - Whether product requires selling plan
- `tags: [String!]!` - Product tags
- `options: [ProductOption!]!` - Product options
- `metafields: JSON` - Custom metadata
- `createdAt: String!` - Creation timestamp
- `updatedAt: String` - Last update timestamp
- `variants: [ProductVariant!]!` - Product variants
- `collections: [Collection!]!` - Collections this product belongs to
- `category: Category` - Product category
- `brand: Brand` - Product brand

#### Order
Represents a customer order.

**Fields:**
- `id: ID!` - Unique identifier
- `orderNumber: String!` - Order number
- `merchantId: ID!` - ID of the merchant
- `customerId: ID` - ID of the customer
- `locationId: ID` - ID of the location
- `status: OrderStatus!` - Order status
- `fulfillmentStatus: FulfillmentStatus!` - Fulfillment status
- `paymentStatus: PaymentStatus!` - Payment status
- `customer: CustomerInfo!` - Customer information
- `billingAddress: Address!` - Billing address
- `shippingAddress: Address!` - Shipping address
- `subtotalPrice: Float!` - Subtotal price
- `totalTax: Float!` - Total tax
- `totalShipping: Float!` - Total shipping cost
- `totalDiscount: Float!` - Total discount
- `totalPrice: Float!` - Total price
- `shippingMethod: String` - Shipping method
- `shippingRate: Float` - Shipping rate
- `trackingNumber: String` - Tracking number
- `trackingUrl: String` - Tracking URL
- `carrier: String` - Shipping carrier
- `source: OrderSource!` - Order source
- `channel: String` - Sales channel
- `currency: String!` - Currency code
- `tags: [String!]!` - Order tags
- `notes: String` - Customer notes
- `internalNotes: String` - Internal notes
- `processedAt: String` - Processing timestamp
- `cancelledAt: String` - Cancellation timestamp
- `fulfilledAt: String` - Fulfillment timestamp
- `closedAt: String` - Closure timestamp
- `createdAt: String!` - Creation timestamp
- `updatedAt: String` - Last update timestamp
- `lineItems: [OrderLineItem!]!` - Order line items
- `fulfillments: [Fulfillment!]!` - Order fulfillments
- `customerUser: User` - Associated user account
- `transactions: [Transaction!]!` - Payment transactions

#### Payment
Represents a payment transaction.

**Fields:**
- `id: ID!` - Unique identifier
- `orderId: ID` - Associated order ID
- `customerId: ID` - Customer ID
- `merchantId: ID!` - Merchant ID
- `amount: Float!` - Payment amount
- `currency: String!` - Currency code
- `status: PaymentStatus!` - Payment status
- `gateway: String!` - Payment gateway
- `gatewayTransactionId: String` - Gateway transaction ID
- `paymentMethodId: String` - Payment method ID
- `paymentMethodType: PaymentMethodType!` - Payment method type
- `billingAddress: Address` - Billing address
- `processorResponse: String` - Processor response
- `failureReason: String` - Failure reason
- `metadata: JSON` - Custom metadata
- `description: String` - Payment description
- `authorizedAt: String` - Authorization timestamp
- `capturedAt: String` - Capture timestamp
- `failedAt: String` - Failure timestamp
- `refundedAt: String` - Refund timestamp
- `voidedAt: String` - Void timestamp
- `createdAt: String!` - Creation timestamp
- `updatedAt: String` - Last update timestamp
- `refunds: [Refund!]!` - Refunds for this payment
- `order: Order` - Associated order
- `customer: User` - Associated customer

## Queries

### User Queries

#### user(id: ID!)
Fetch a user by ID.

**Example:**
```graphql
query {
  user(id: "123") {
    id
    email
    firstName
    lastName
    roles {
      name
      description
    }
  }
}
```

#### users(limit: Int, offset: Int)
Fetch a list of users with pagination.

**Example:**
```graphql
query {
  users(limit: 10, offset: 0) {
    id
    email
    firstName
    lastName
    createdAt
  }
}
```

#### currentUser
Fetch the currently authenticated user.

**Example:**
```graphql
query {
  currentUser {
    id
    email
    firstName
    lastName
  }
}
```

### Product Queries

#### product(id: ID!)
Fetch a product by ID.

**Example:**
```graphql
query {
  product(id: "456") {
    id
    title
    description
    priceRange {
      minVariantPrice
      maxVariantPrice
    }
    images {
      src
      altText
    }
    variants {
      id
      title
      price
      inventoryQuantity
    }
  }
}
```

#### products(filter: ProductFilter)
Fetch a list of products with filtering.

**Example:**
```graphql
query {
  products(filter: {
    status: ACTIVE
    limit: 20
  }) {
    id
    title
    handle
    featuredImage
    priceRange {
      minVariantPrice
    }
  }
}
```

#### productByHandle(handle: String!)
Fetch a product by its handle.

**Example:**
```graphql
query {
  productByHandle(handle: "awesome-product") {
    id
    title
    description
  }
}
```

#### categories(filter: CategoryFilter)
Fetch a list of categories.

**Example:**
```graphql
query {
  categories(filter: {
    isVisible: true
  }) {
    id
    name
    handle
    children {
      id
      name
    }
  }
}
```

### Order Queries

#### order(id: ID!)
Fetch an order by ID.

**Example:**
```graphql
query {
  order(id: "789") {
    id
    orderNumber
    status
    totalPrice
    lineItems {
      id
      name
      quantity
      price
    }
  }
}
```

#### orders(filter: OrderFilter)
Fetch a list of orders with filtering.

**Example:**
```graphql
query {
  orders(filter: {
    customerId: "123"
    limit: 10
  }) {
    id
    orderNumber
    status
    totalPrice
    createdAt
  }
}
```

### Payment Queries

#### payment(id: ID!)
Fetch a payment by ID.

**Example:**
```graphql
query {
  payment(id: "101") {
    id
    amount
    currency
    status
    gateway
  }
}
```

#### payments(filter: PaymentFilter)
Fetch a list of payments with filtering.

**Example:**
```graphql
query {
  payments(filter: {
    customerId: "123"
    limit: 10
  }) {
    id
    amount
    currency
    status
    createdAt
  }
}
```

## Mutations

### User Mutations

#### register(input: RegisterInput!)
Register a new user.

**Example:**
```graphql
mutation {
  register(input: {
    email: "newuser@example.com"
    username: "newuser"
    password: "securepassword"
    firstName: "New"
    lastName: "User"
  }) {
    user {
      id
      email
      firstName
      lastName
    }
    accessToken
    refreshToken
  }
}
```

#### login(input: LoginInput!)
Authenticate a user.

**Example:**
```graphql
mutation {
  login(input: {
    email: "user@example.com"
    password: "password123"
  }) {
    user {
      id
      email
    }
    accessToken
    refreshToken
  }
}
```

#### updateUser(id: ID!, input: UpdateUserInput!)
Update a user's information.

**Example:**
```graphql
mutation {
  updateUser(id: "123", input: {
    firstName: "Updated"
    lastName: "Name"
  }) {
    id
    firstName
    lastName
  }
}
```

### Product Mutations

#### createProduct(input: CreateProductInput!)
Create a new product.

**Example:**
```graphql
mutation {
  createProduct(input: {
    merchantId: "456"
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

#### updateProduct(id: ID!, input: UpdateProductInput!)
Update an existing product.

**Example:**
```graphql
mutation {
  updateProduct(id: "789", input: {
    title: "Updated Product Title"
    description: "An even better product"
  }) {
    id
    title
    description
  }
}
```

### Order Mutations

#### createOrder(input: CreateOrderInput!)
Create a new order.

**Example:**
```graphql
mutation {
  createOrder(input: {
    merchantId: "456"
    customer: {
      email: "customer@example.com"
      firstName: "Customer"
      lastName: "Name"
    }
    billingAddress: {
      firstName: "Customer"
      lastName: "Name"
      street1: "123 Main St"
      city: "Anytown"
      state: "CA"
      country: "US"
      postalCode: "12345"
    }
    shippingAddress: {
      firstName: "Customer"
      lastName: "Name"
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

### Payment Mutations

#### createPayment(input: CreatePaymentInput!)
Create a new payment.

**Example:**
```graphql
mutation {
  createPayment(input: {
    orderId: "789"
    merchantId: "456"
    amount: 99.99
    currency: "USD"
    paymentMethodId: "pm_123"
  }) {
    id
    amount
    status
  }
}
```

## Enums

### ProductStatus
- `ACTIVE` - Product is active and visible
- `ARCHIVED` - Product is archived
- `DRAFT` - Product is a draft

### OrderStatus
- `PENDING` - Order is pending
- `CONFIRMED` - Order is confirmed
- `PROCESSING` - Order is being processed
- `SHIPPED` - Order has been shipped
- `DELIVERED` - Order has been delivered
- `CANCELLED` - Order has been cancelled
- `REFUNDED` - Order has been refunded
- `RETURNED` - Order has been returned

### PaymentStatus
- `PENDING` - Payment is pending
- `AUTHORIZED` - Payment is authorized
- `CAPTURED` - Payment is captured
- `FAILED` - Payment failed
- `CANCELLED` - Payment is cancelled
- `REFUNDED` - Payment is refunded
- `PARTIALLY_REFUNDED` - Payment is partially refunded
- `VOIDED` - Payment is voided
- `PARTIALLY_PAID` - Payment is partially paid
- `PAID` - Payment is paid

## Input Objects

### RegisterInput
- `email: String!` - User's email
- `username: String!` - Username
- `password: String!` - Password
- `firstName: String!` - First name
- `lastName: String!` - Last name
- `phone: String` - Phone number

### CreateProductInput
- `merchantId: ID!` - Merchant ID
- `title: String!` - Product title
- `description: String` - Product description
- `vendor: String` - Vendor
- `productType: String` - Product type
- `tags: [String!]` - Product tags
- `status: ProductStatus` - Product status
- `publishedScope: PublishedScope` - Publication scope
- `seoTitle: String` - SEO title
- `seoDescription: String` - SEO description
- `categoryId: ID` - Category ID
- `brandId: ID` - Brand ID
- `metafields: JSON` - Custom metadata

### CreateOrderInput
- `merchantId: ID!` - Merchant ID
- `customerId: ID` - Customer ID
- `customer: CustomerInfoInput!` - Customer information
- `billingAddress: AddressInput!` - Billing address
- `shippingAddress: AddressInput!` - Shipping address
- `currency: String` - Currency code
- `source: OrderSource` - Order source
- `channel: String` - Sales channel
- `notes: String` - Customer notes

## Error Codes

The API uses standardized error codes to help developers understand and handle errors:

- `USER_NOT_FOUND` - The specified user could not be found
- `INVALID_CREDENTIALS` - The provided credentials are invalid
- `PRODUCT_NOT_FOUND` - The specified product could not be found
- `ORDER_NOT_FOUND` - The specified order could not be found
- `PAYMENT_FAILED` - The payment could not be processed
- `INSUFFICIENT_INVENTORY` - There is not enough inventory for the requested items
- `VALIDATION_ERROR` - The input data failed validation
- `PERMISSION_DENIED` - The user does not have permission to perform this action
- `RATE_LIMIT_EXCEEDED` - The rate limit has been exceeded

## Best Practices

### Query Optimization
1. Only request the fields you need
2. Use aliases for clarity when querying the same field multiple times
3. Use fragments to avoid repetition
4. Batch related queries when possible

### Error Handling
1. Always check for errors in the response
2. Implement retry logic for transient errors
3. Log errors for debugging purposes
4. Provide user-friendly error messages

### Security
1. Never expose sensitive data in client-side code
2. Use HTTPS for all API requests
3. Implement proper authentication and authorization
4. Validate all input data
5. Use rate limiting to prevent abuse

### Performance
1. Implement caching for frequently accessed data
2. Use pagination for large data sets
3. Minimize the number of API requests
4. Use connection pooling for database connections