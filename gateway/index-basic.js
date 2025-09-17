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
    console.log('🚀 Starting GraphQL Federation Gateway (Basic Services)...');

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

    // CORS configuration for live applications
    app.use(cors({
      origin: [
        'https://unified-commerce.vercel.app',
        'https://admin-panel-igp522vr5-crypticogs-projects.vercel.app',
        'http://localhost:3000',
        'http://localhost:3001', 
        'http://localhost:3003',
        'http://localhost:4000',
        'http://localhost:5173'
      ],
      credentials: true,
    }));

    // Create Apollo Gateway with core services (no cross-dependencies)
    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: 'http://localhost:8001/graphql' },
          { name: 'cart', url: 'http://localhost:8002/graphql' },
          { name: 'payment', url: 'http://localhost:8004/graphql' },
          { name: 'product', url: 'http://localhost:8006/graphql' },
          { name: 'merchant', url: 'http://localhost:8008/graphql' }
        ],
      })
    });

    // Create Apollo Server with the gateway
    const server = new ApolloServer({
      gateway,
      context: getGraphQLContext,
      introspection: true,
      plugins: [
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
          subgraphs: 5,
          active: true
        },
        services: {
          identity: 'http://localhost:8001/graphql',
          cart: 'http://localhost:8002/graphql',
          payment: 'http://localhost:8004/graphql',
          product: 'http://localhost:8006/graphql',
          merchant: 'http://localhost:8008/graphql'
        },
        frontend_apps: {
          storefront: 'https://unified-commerce.vercel.app',
          admin: 'https://admin-panel-igp522vr5-crypticogs-projects.vercel.app'
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
      console.log('\n📊 Federated Services (5/8 active):');
      console.log('  ✅ Identity Service: http://localhost:8001/graphql');
      console.log('  ✅ Cart Service: http://localhost:8002/graphql');
      console.log('  ✅ Payment Service: http://localhost:8004/graphql');
      console.log('  ✅ Product Catalog Service: http://localhost:8006/graphql');
      console.log('  ✅ Merchant Account Service: http://localhost:8008/graphql');
      console.log('\n🌐 Live Frontend Applications:');
      console.log('  📱 Storefront: https://unified-commerce.vercel.app');
      console.log('  🏢 Admin Panel: https://admin-panel-igp522vr5-crypticogs-projects.vercel.app');
      console.log('\n🎉 Ready to test frontend-backend integration!');
    });

  } catch (error) {
    console.error('❌ Failed to start GraphQL Federation Gateway:', error);
    console.error('Error details:', error.message);
    process.exit(1);
  }
}

startGateway();