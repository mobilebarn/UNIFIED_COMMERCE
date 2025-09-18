const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const jwt = require('jsonwebtoken');
const { ApolloServer } = require('@apollo/server');
const { expressMiddleware } = require('@apollo/server/express4');
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { json } = require('body-parser');

// Authentication middleware for GraphQL context
function getGraphQLContext({ req }) {
  const token = req.headers.authorization?.replace('Bearer ', '');

  if (!token) {
    return {
      user: null,
      isAuthenticated: false
    };
  }

  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET || 'your-secret-key');
    return {
      user: decoded,
      isAuthenticated: true,
      userId: decoded.user_id,
      email: decoded.email,
      roles: decoded.roles || []
    };
  } catch (error) {
    console.warn('Invalid JWT token:', error.message);
    return {
      user: null,
      isAuthenticated: false
    };
  }
}

async function startGateway() {
  try {
    console.log('🚀 Starting GraphQL Federation Gateway...');

    // Create Express app
    const app = express();

    // Security middleware
    app.use(helmet({
      contentSecurityPolicy: {
        directives: {
          defaultSrc: ["'self'"],
          scriptSrc: ["'self'", "'unsafe-inline'", "https://unpkg.com"],
          styleSrc: ["'self'", "'unsafe-inline'", "https://fonts.googleapis.com"],
          fontSrc: ["'self'", "https://fonts.gstatic.com"],
          imgSrc: ["'self'", "data:", "https:"],
        },
      },
    }));

    // CORS configuration
    app.use(cors({
      origin: process.env.NODE_ENV === 'production'
        ? ['https://yourdomain.com', 'https://admin.yourdomain.com']
        : ['http://localhost:3000', 'http://localhost:3001', 'http://localhost:3003', 'http://localhost:4000'],
      credentials: true,
    }));

    // Create Apollo Gateway with service definitions
    // Use environment variables from Render deployment
    const getServiceUrl = (envVarName, fallbackPort) => {
      const baseUrl = process.env[envVarName];
      if (baseUrl) {
        // Render provides the host URL, append /graphql path
        return `${baseUrl}/graphql`;
      }
      // Fallback for local development
      return `http://localhost:${fallbackPort}/graphql`;
    };

    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: getServiceUrl('IDENTITY_SERVICE_URL', 8001) },
          { name: 'product-catalog', url: getServiceUrl('PRODUCT_CATALOG_SERVICE_URL', 8002) },
          { name: 'order', url: getServiceUrl('ORDER_SERVICE_URL', 8003) },
          { name: 'payment', url: getServiceUrl('PAYMENT_SERVICE_URL', 8004) },
          { name: 'inventory', url: getServiceUrl('INVENTORY_SERVICE_URL', 8005) },
          { name: 'merchant-account', url: getServiceUrl('MERCHANT_ACCOUNT_SERVICE_URL', 8006) },
          { name: 'cart', url: getServiceUrl('CART_SERVICE_URL', 8007) },
          { name: 'promotions', url: getServiceUrl('PROMOTIONS_SERVICE_URL', 8008) },
          { name: 'analytics', url: getServiceUrl('ANALYTICS_SERVICE_URL', 8009) }
        ],
      })
    });

    // Create Apollo Server with the gateway
    const server = new ApolloServer({
      gateway,
      context: getGraphQLContext,
      // Enable introspection and playground in development
      introspection: process.env.NODE_ENV !== 'production',
      plugins: [
        // Custom plugin for logging
        {
          requestDidStart() {
            return {
              didResolveOperation(requestContext) {
                console.log(`🔍 GraphQL Operation: ${requestContext.request.operationName || 'Anonymous'}`);
              },
              didEncounterErrors(requestContext) {
                console.error('❌ GraphQL Errors:', requestContext.errors);
              }
            };
          }
        }
      ]
    });

    // Start the server
    await server.start();

    // Health check endpoint
    app.get('/health', (req, res) => {
      res.json({
        service: 'graphql-federation-gateway',
        status: 'healthy',
        time: new Date().toISOString(),
        federation: {
          subgraphs: 9, // All 9 services including analytics
          active: true
        },
        services: {
          identity: getServiceUrl('IDENTITY_SERVICE_URL', 8001),
          'product-catalog': getServiceUrl('PRODUCT_CATALOG_SERVICE_URL', 8002),
          order: getServiceUrl('ORDER_SERVICE_URL', 8003),
          payment: getServiceUrl('PAYMENT_SERVICE_URL', 8004),
          inventory: getServiceUrl('INVENTORY_SERVICE_URL', 8005),
          'merchant-account': getServiceUrl('MERCHANT_ACCOUNT_SERVICE_URL', 8006),
          cart: getServiceUrl('CART_SERVICE_URL', 8007),
          promotions: getServiceUrl('PROMOTIONS_SERVICE_URL', 8008),
          analytics: getServiceUrl('ANALYTICS_SERVICE_URL', 8009)
        }
      });
    });

    // Apply the Apollo GraphQL middleware
    app.use('/graphql',
      json(),
      expressMiddleware(server, {
        context: getGraphQLContext
      })
    );

    // GraphQL playground redirect
    app.get('/', (req, res) => {
      res.redirect('/graphql');
    });

    // Start server
    const PORT = process.env.PORT || 4000;
    app.listen(PORT, () => {
      console.log(`✅ GraphQL Federation Gateway running at http://localhost:${PORT}/graphql`);
      console.log(`🎮 GraphQL Playground available at http://localhost:${PORT}/graphql`);
      console.log(`🔍 Health check available at http://localhost:${PORT}/health`);
      console.log('\n📊 Federated Services:');
      console.log(`  ✅ Identity Service: ${getServiceUrl('IDENTITY_SERVICE_URL', 8001)}`);
      console.log(`  ✅ Product Catalog Service: ${getServiceUrl('PRODUCT_CATALOG_SERVICE_URL', 8002)}`);
      console.log(`  ✅ Order Service: ${getServiceUrl('ORDER_SERVICE_URL', 8003)}`);
      console.log(`  ✅ Payment Service: ${getServiceUrl('PAYMENT_SERVICE_URL', 8004)}`);
      console.log(`  ✅ Inventory Service: ${getServiceUrl('INVENTORY_SERVICE_URL', 8005)}`);
      console.log(`  ✅ Merchant Account Service: ${getServiceUrl('MERCHANT_ACCOUNT_SERVICE_URL', 8006)}`);
      console.log(`  ✅ Cart Service: ${getServiceUrl('CART_SERVICE_URL', 8007)}`);
      console.log(`  ✅ Promotions Service: ${getServiceUrl('PROMOTIONS_SERVICE_URL', 8008)}`);
      console.log(`  ✅ Analytics Service: ${getServiceUrl('ANALYTICS_SERVICE_URL', 8009)}`);
      console.log('\n🎉 All 9 services are now connected to the GraphQL Federation Gateway!');
      console.log('\n🔧 Next Steps:');
      console.log('  1. Test unified GraphQL queries across all connected services');
      console.log('  2. Verify cross-service relationships and federated queries');
      console.log('  3. Monitor performance and optimize as needed');
    });

  } catch (error) {
    console.error('❌ Failed to start GraphQL Federation Gateway:', error);
    console.error('Error details:', error.message);
    process.exit(1);
  }
}

// Start the gateway
startGateway();