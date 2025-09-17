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
        ? ['https://unified-commerce.vercel.app', 'https://admin-panel-igp522vr5-crypticogs-projects.vercel.app']
        : ['http://localhost:3000', 'http://localhost:3001', 'http://localhost:3003', 'http://localhost:4000'],
      credentials: true,
    }));

    // Create Apollo Gateway with ONLY working services
    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: 'http://localhost:8001/graphql' },
          { name: 'cart', url: 'http://localhost:8002/graphql' },
          { name: 'payment', url: 'http://localhost:8004/graphql' },
          { name: 'product', url: 'http://localhost:8006/graphql' },
          { name: 'promotions', url: 'http://localhost:8007/graphql' },
          { name: 'merchant', url: 'http://localhost:8008/graphql' }
        ],
      })
    });

    // Create Apollo Server with the gateway
    const server = new ApolloServer({
      gateway,
      context: getGraphQLContext,
      // Enable introspection and playground in development
      introspection: true,
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
          subgraphs: 6, // Only working services
          active: true
        },
        services: {
          identity: 'http://localhost:8001/graphql',
          cart: 'http://localhost:8002/graphql',
          payment: 'http://localhost:8004/graphql',
          product: 'http://localhost:8006/graphql',
          promotions: 'http://localhost:8007/graphql',
          merchant: 'http://localhost:8008/graphql'
        },
        disabled: {
          order: 'http://localhost:8003/graphql - Kafka connection failed',
          inventory: 'http://localhost:8005/graphql - Kafka connection failed'
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
      console.log('\n📊 Federated Services (6/8 active):');
      console.log('  ✅ Identity Service: http://localhost:8001/graphql');
      console.log('  ✅ Cart Service: http://localhost:8002/graphql');
      console.log('  ✅ Payment Service: http://localhost:8004/graphql');
      console.log('  ✅ Product Catalog Service: http://localhost:8006/graphql');
      console.log('  ✅ Promotions Service: http://localhost:8007/graphql');
      console.log('  ✅ Merchant Account Service: http://localhost:8008/graphql');
      console.log('\n⚠️  Temporarily disabled (Kafka issues):');
      console.log('  ❌ Order Service: http://localhost:8003/graphql');
      console.log('  ❌ Inventory Service: http://localhost:8005/graphql');
      console.log('\n🎉 GraphQL Federation Gateway is ready for testing!');
    });

  } catch (error) {
    console.error('❌ Failed to start GraphQL Federation Gateway:', error);
    console.error('Error details:', error.message);
    process.exit(1);
  }
}

// Start the gateway
startGateway();