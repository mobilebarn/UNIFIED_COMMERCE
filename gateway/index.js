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
    // NOTE: Use environment variables for Railway deployment or fallback to localhost
    const getServiceUrl = (serviceName, defaultPort) => {
      const envVar = `${serviceName.toUpperCase()}_SERVICE_URL`;
      return process.env[envVar] || `http://localhost:${defaultPort}/graphql`;
    };

    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: getServiceUrl('identity', 8001) },
          { name: 'cart', url: getServiceUrl('cart', 8002) },
          { name: 'order', url: getServiceUrl('order', 8003) },
          { name: 'payment', url: getServiceUrl('payment', 8004) },
          { name: 'inventory', url: getServiceUrl('inventory', 8005) },
          { name: 'product', url: getServiceUrl('product', 8006) },
          { name: 'promotions', url: getServiceUrl('promotions', 8007) },
          { name: 'merchant', url: getServiceUrl('merchant', 8008) }
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
          subgraphs: 8, // All 8 services
          active: true
        },
        services: {
          identity: getServiceUrl('identity', 8001),
          cart: getServiceUrl('cart', 8002),
          order: getServiceUrl('order', 8003),
          payment: getServiceUrl('payment', 8004),
          inventory: getServiceUrl('inventory', 8005),
          product: getServiceUrl('product', 8006),
          promotions: getServiceUrl('promotions', 8007),
          merchant: getServiceUrl('merchant', 8008)
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
      console.log(`  ✅ Identity Service: ${getServiceUrl('identity', 8001)}`);
      console.log(`  ✅ Cart Service: ${getServiceUrl('cart', 8002)}`);
      console.log(`  ✅ Order Service: ${getServiceUrl('order', 8003)}`);
      console.log(`  ✅ Payment Service: ${getServiceUrl('payment', 8004)}`);
      console.log(`  ✅ Inventory Service: ${getServiceUrl('inventory', 8005)}`);
      console.log(`  ✅ Product Catalog Service: ${getServiceUrl('product', 8006)}`);
      console.log(`  ✅ Promotions Service: ${getServiceUrl('promotions', 8007)}`);
      console.log(`  ✅ Merchant Account Service: ${getServiceUrl('merchant', 8008)}`);
      console.log('\n🎉 All 8 services are now connected to the GraphQL Federation Gateway!');
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