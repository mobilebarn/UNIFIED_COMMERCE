# Retail OS Point of Sale (POS) Application

This is a React Native mobile application for the Retail OS Point of Sale system, built with Expo and integrated with the Retail OS GraphQL Federation Gateway.

## Features

- Product browsing and search
- Shopping cart management
- Checkout process with multiple payment methods
- Integration with Retail OS GraphQL Federation Gateway
- Real-time inventory updates
- Payment processing

## Prerequisites

- Node.js (version 16 or higher)
- Expo CLI
- Android Studio or Xcode for mobile development
- Retail OS GraphQL Federation Gateway running locally

## Installation

1. Clone the repository
2. Navigate to the `mobile-pos` directory
3. Install dependencies:
   ```bash
   npm install
   ```

## Running the Application

To start the development server:
```bash
npm start
```

Then follow the instructions to run on:
- iOS: Press `i`
- Android: Press `a`
- Web: Press `w`

## Architecture

The application follows a clean architecture pattern with the following components:

### GraphQL Integration

The app connects to the Retail OS GraphQL Federation Gateway which unifies all microservice APIs into a single endpoint. This allows the POS application to access data from multiple services including:

- Product Catalog Service
- Inventory Service
- Cart & Checkout Service
- Payment Service
- Order Service
- Merchant Account Service

### Data Flow

1. **Product Browsing**: Users can search and browse products from the Product Catalog Service
2. **Cart Management**: Items are added to a cart managed by the Cart & Checkout Service
3. **Checkout**: Payment processing is handled through the Payment Service
4. **Order Creation**: Orders are created and managed by the Order Service

### Folder Structure

```
mobile-pos/
├── app/                 # Application screens and navigation
│   ├── auth/           # Authentication screens
│   ├── pos/            # Point of Sale screens
│   └── (tabs)/         # Tab navigation
├── components/         # Reusable UI components
├── graphql/            # GraphQL queries and mutations
├── utils/              # Utility functions
└── ...
```

## GraphQL Schema

The application uses the following GraphQL queries and mutations:

### Product Queries
- `GET_PRODUCTS`: Fetch products with filtering
- `GET_PRODUCT`: Get a single product by ID
- `SEARCH_PRODUCTS`: Search for products by name or barcode

### Cart Queries and Mutations
- `GET_CART`: Get a cart by ID
- `GET_CART_BY_SESSION`: Get a cart by session ID
- `CREATE_CART`: Create a new cart
- `ADD_CART_LINE_ITEM`: Add an item to the cart
- `UPDATE_CART_LINE_ITEM`: Update a cart line item
- `REMOVE_CART_LINE_ITEM`: Remove an item from the cart
- `CLEAR_CART`: Clear the cart

### Payment Queries and Mutations
- `GET_PAYMENT_METHODS`: Get payment methods for a customer
- `CREATE_PAYMENT`: Create a payment
- `CAPTURE_PAYMENT`: Capture a payment transaction
- `REFUND_PAYMENT`: Refund a payment
- `VOID_PAYMENT`: Void a payment

## Development

### Adding New Features

1. Create new GraphQL queries/mutations in the `graphql/` directory
2. Add new screens in the `app/` directory
3. Use the Apollo Client hooks (`useQuery`, `useMutation`) to interact with the GraphQL API

### State Management

The application uses:
- Apollo Client for GraphQL state management
- React Context API for global state (authentication, etc.)
- Local component state for UI interactions

## Testing

To run tests:
```bash
npm test
```

## Deployment

To build for production:
```bash
# For Android
expo build:android

# For iOS
expo build:ios
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a pull request

## License

This project is licensed under the MIT License.