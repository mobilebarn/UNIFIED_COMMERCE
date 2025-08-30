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
    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: 'http://localhost:8001/graphql' },
          { name: 'cart', url: 'http://localhost:8002/graphql' },
          { name: 'order', url: 'http://localhost:8003/graphql' },
          { name: 'payment', url: 'http://localhost:8004/graphql' },
          { name: 'inventory', url: 'http://localhost:8005/graphql' },
          { name: 'product-catalog', url: 'http://localhost:8006/graphql' },
          { name: 'promotions', url: 'http://localhost:8007/graphql' },
          { name: 'merchant-account', url: 'http://localhost:8008/graphql' }
        ],
      }),
      buildService({ url }) {
        return {
          url,
          // Forward authentication headers to subgraphs
          willSendRequest({ request, context }) {
            if (context.userId) {
              request.http.headers.set('user-id', context.userId);
              request.http.headers.set('user-email', context.email);
              request.http.headers.set('user-roles', JSON.stringify(context.roles));
            }
          }
        };
      }
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
        status: 'healthy', 
        service: 'graphql-federation-gateway', 
        timestamp: new Date().toISOString(),
        subgraphs: [
          'identity', 'cart', 'order', 'payment', 
          'inventory', 'product-catalog', 'promotions'
        ]
      });
    });

    // Apply the Apollo GraphQL middleware
    app.use('/graphql', 
      json(),
      expressMiddleware(server, {
        context: getGraphQLContext
      })
    );

    // Serve GraphQL Playground in development
    if (process.env.NODE_ENV !== 'production') {
      app.get('/playground', (req, res) => {
        res.send(`
          <!DOCTYPE html>
          <html>
          <head>
            <title>GraphQL Federation Gateway Playground</title>
            <style>
              body { 
                font-family: Arial, sans-serif; 
                margin: 0; 
                padding: 20px; 
                background: #f5f5f5; 
              }
              .container { 
                max-width: 800px; 
                margin: 0 auto; 
                background: white; 
                padding: 30px; 
                border-radius: 8px; 
                box-shadow: 0 2px 10px rgba(0,0,0,0.1);
              }
              h1 { color: #333; text-align: center; }
              .endpoint { 
                background: #f0f0f0; 
                padding: 10px; 
                border-radius: 4px; 
                margin: 10px 0; 
                font-family: monospace; 
              }
              .services { 
                display: grid; 
                grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); 
                gap: 15px; 
                margin: 20px 0; 
              }
              .service { 
                background: #e8f4f8; 
                padding: 15px; 
                border-radius: 6px; 
                text-align: center;
              }
              .service h3 { margin: 0 0 10px 0; color: #2c5282; }
              .service p { margin: 5px 0; font-size: 14px; color: #666; }
            </style>
          </head>
          <body>
            <div class="container">
              <h1>🚀 GraphQL Federation Gateway</h1>
              <p><strong>Unified Commerce Platform</strong> - GraphQL Federation Gateway providing a single endpoint for all microservices.</p>
              
              <h2>📍 GraphQL Endpoint</h2>
              <div class="endpoint">POST http://localhost:4000/graphql</div>
              
              <h2>🔧 Federated Services</h2>
              <div class="services">
                <div class="service">
                  <h3>Identity</h3>
                  <p>User authentication, roles & permissions</p>
                  <p><code>:8001/graphql</code></p>
                </div>
                <div class="service">
                  <h3>Cart</h3>
                  <p>Shopping cart management</p>
                  <p><code>:8002/graphql</code></p>
                </div>
                <div class="service">
                  <h3>Order</h3>
                  <p>Order processing & fulfillment</p>
                  <p><code>:8003/graphql</code></p>
                </div>
                <div class="service">
                  <h3>Payment</h3>
                  <p>Payment processing & transactions</p>
                  <p><code>:8004/graphql</code></p>
                </div>
                <div class="service">
                  <h3>Inventory</h3>
                  <p>Stock management & tracking</p>
                  <p><code>:8005/graphql</code></p>
                </div>
                <div class="service">
                  <h3>Product Catalog</h3>
                  <p>Products, variants & categories</p>
                  <p><code>:8006/graphql</code></p>
                </div>
                <div class="service">
                  <h3>Promotions</h3>
                  <p>Discounts, campaigns & loyalty</p>
                  <p><code>:8007/graphql</code></p>
                </div>
              </div>
              
              <h2>🔗 Example Query</h2>
              <div class="endpoint">
query {
  user(id: "1") {
    id
    email
    firstName
    lastName
  }
  
  products(filter: { limit: 5 }) {
    id
    title
    status
    variants {
      id
      sku
      price
    }
  }
}
              </div>
              
                            
              <p><em>Use GraphQL Playground, Apollo Studio, or any GraphQL client to interact with the federated schema.</em></p>
            </div>
          </body>
          </html>
        `);
      });
    }

    // Start server
    const PORT = process.env.PORT || 4000;
    
    app.listen(PORT, () => {
      console.log(`✅ GraphQL Federation Gateway running on http://localhost:${PORT}`);
      console.log(`📊 Health check: http://localhost:${PORT}/health`);
      console.log(`🎮 GraphQL endpoint: http://localhost:${PORT}/graphql`);
      if (process.env.NODE_ENV !== 'production') {
        console.log(`🎪 Playground: http://localhost:${PORT}/playground`);
      }
      console.log('🔗 Federated services:');
      console.log('  - Identity: http://localhost:8001/graphql');
      console.log('  - Cart: http://localhost:8002/graphql');
      console.log('  - Order: http://localhost:8003/graphql');
      console.log('  - Payment: http://localhost:8004/graphql');
      console.log('  - Inventory: http://localhost:8005/graphql');
      console.log('  - Product Catalog: http://localhost:8006/graphql');
      console.log('  - Promotions: http://localhost:8007/graphql');
      console.log('  - Merchant Account: http://localhost:8008/graphql');
    });

  } catch (error) {
    console.error('❌ Failed to start GraphQL Federation Gateway:', error);
    process.exit(1);
  }
}

// Handle graceful shutdown
process.on('SIGTERM', () => {
  console.log('🛑 Received SIGTERM, shutting down gracefully...');
  process.exit(0);
});

process.on('SIGINT', () => {
  console.log('🛑 Received SIGINT, shutting down gracefully...');
  process.exit(0);
});

// Start the gateway
startGateway();
            </div>
          </body>
          </html>
        `);
      });
    }
          if (req.userId) {
            proxyReq.setHeader('user-id', req.userId);
            proxyReq.setHeader('user-email', req.email);
            proxyReq.setHeader('user-roles', JSON.stringify(req.roles));
          }
        },
        onError: (err, req, res) => {
          console.error(`Proxy error for ${serviceName}:`, err.message);
          res.status(503).json({ 
            error: 'Service temporarily unavailable', 
            service: serviceName,
            message: err.message 
          });
        }
      }));
    });

    // GraphQL endpoint specifically for identity service
    app.use('/graphql', createProxyMiddleware({
      target: 'http://localhost:8001',
      changeOrigin: true,
      onProxyReq: (proxyReq, req, res) => {
        if (req.userId) {
          proxyReq.setHeader('user-id', req.userId);
          proxyReq.setHeader('user-email', req.email);
          proxyReq.setHeader('user-roles', JSON.stringify(req.roles));
        }
      }
    }));

    // API info endpoint
    app.get('/', (req, res) => {
      res.json({
        name: 'Unified Commerce API Gateway',
        version: '1.0.0',
        endpoints: {
          graphql: '/graphql (Identity service)',
          auth: '/auth (Identity service)',
          cart: '/api/cart',
          order: '/api/order',
          payment: '/api/payment',
          inventory: '/api/inventory',
          'product-catalog': '/api/product-catalog',
          promotions: '/api/promotions',
          health: '/health'
        },
        description: 'REST API Gateway with authentication proxy'
      });
    });

    // Catch-all route
    app.use('*', (req, res) => {
      res.status(404).json({ 
        error: 'Route not found', 
        path: req.originalUrl,
        availableRoutes: ['/health', '/auth', '/graphql', ...Object.keys(services).filter(s => s !== 'identity').map(s => `/api/${s}`)]
      });
    });

    // Start server
    const port = process.env.PORT || 4000;
    app.listen(port, () => {
      console.log(`✅ API Gateway running on http://localhost:${port}`);
      console.log(`� Available routes:`);
      console.log(`   🔐 Auth/GraphQL: /auth, /graphql`);
      Object.keys(services).forEach(service => {
        if (service !== 'identity') {
          console.log(`   📦 ${service}: /api/${service}`);
        }
      });
      console.log(`   📊 Health: /health`);
      console.log(`   📖 Info: /`);
    });

  } catch (error) {
    console.error('❌ Failed to start gateway:', error);
    process.exit(1);
  }
}

// Add global error handlers
process.on('unhandledRejection', (reason, promise) => {
  console.error('❌ Unhandled Promise Rejection:', reason);
});

process.on('uncaughtException', (error) => {
  console.error('❌ Uncaught Exception:', error);
  process.exit(1);
});

startGateway();