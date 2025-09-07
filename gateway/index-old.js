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

    // Create Apollo Gateway with service definitions (starting with just identity)
    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: 'http://localhost:8001/graphql' }
          // Temporarily disabled other services due to schema issues
          // { name: 'cart', url: 'http://localhost:8002/graphql' },
          // { name: 'order', url: 'http://localhost:8003/graphql' },
          // { name: 'payment', url: 'http://localhost:8004/graphql' },
          // { name: 'inventory', url: 'http://localhost:8005/graphql' },
          // { name: 'product-catalog', url: 'http://localhost:8006/graphql' },
          // { name: 'promotions', url: 'http://localhost:8007/graphql' }
        ],
      }),
    });

    // Create Apollo Server with the gateway
    const server = new ApolloServer({
      gateway,
      context: getGraphQLContext,
      introspection: true,
      playground: true,
    });

    // Start the gateway
    await server.start();

    // Apply middleware
    app.use('/graphql', 
      json(),
      expressMiddleware(server, {
        context: async ({ req }) => getGraphQLContext({ req })
      })
    );

    // Health check endpoint
    app.get('/health', (req, res) => {
      res.json({
        service: 'graphql-gateway',
        status: 'healthy',
        time: new Date().toISOString(),
        federation: {
          subgraphs: 1, // Starting with just identity
          active: true
        }
      });
    });

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
      console.log('\n📊 Federated Services (Phase 1):');
      console.log('  ✅ Identity Service: http://localhost:8001/graphql');
      console.log('  ⏳ Cart Service: [SCHEMA ISSUES - needs gqlgen generation]');
      console.log('  ⏳ Order Service: [SCHEMA ISSUES - needs gqlgen generation]');
      console.log('  ⏳ Payment Service: [SCHEMA ISSUES - needs gqlgen generation]');
      console.log('  ⏳ Inventory Service: [SCHEMA ISSUES - needs gqlgen generation]');
      console.log('  ⏳ Product Catalog Service: [SCHEMA ISSUES - needs gqlgen generation]');
      console.log('  ⏳ Promotions Service: [SCHEMA ISSUES - needs gqlgen generation]');
      console.log('  ⏳ Merchant Account Service: [OFFLINE - compilation errors]');
    });

  } catch (error) {
    console.error('❌ Failed to start GraphQL Federation Gateway:', error);
    process.exit(1);
  }
}

// Start the gateway
startGateway();
