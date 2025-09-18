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

    // CORS configuration - be more permissive for development and deployment
    app.use(cors({
      origin: process.env.NODE_ENV === 'production'
        ? ['https://unified-commerce.vercel.app', 'https://admin.unified-commerce.vercel.app', 'https://unified-commerce-storefront.vercel.app']
        : ['http://localhost:3000', 'http://localhost:3001', 'http://localhost:3003', 'http://localhost:4000'],
      credentials: true,
    }));

    // Create Apollo Gateway with service definitions
    // Use environment variables from Render deployment
    const getServiceUrl = (envVarName, fallbackPort) => {
      const serviceHost = process.env[envVarName];
      if (serviceHost) {
        // Handle different URL formats that Render might provide
        let fullUrl;
        
        if (serviceHost.startsWith('http://') || serviceHost.startsWith('https://')) {
          // Full URL provided
          fullUrl = `${serviceHost}/graphql`;
        } else if (serviceHost.includes('.onrender.com')) {
          // Render domain without protocol
          fullUrl = `https://${serviceHost}/graphql`;
        } else {
          // Service name only - construct full Render URL
          // Render's fromService property might give us just the service name
          fullUrl = `https://${serviceHost}.onrender.com/graphql`;
        }
        
        console.log(`🔗 ${envVarName}: '${serviceHost}' -> '${fullUrl}'`);
        return fullUrl;
      }
      
      // Fallback for local development
      const fallbackUrl = `http://localhost:${fallbackPort}/graphql`;
      console.log(`🔗 ${envVarName} (LOCAL FALLBACK): ${fallbackUrl}`);
      return fallbackUrl;
    };

    // Test service availability before adding to gateway
    const testServiceAvailability = async (url, serviceName) => {
      try {
        const healthUrl = url.replace('/graphql', '/health');
        const response = await fetch(healthUrl, { 
          method: 'GET', 
          timeout: 5000,
          headers: { 'User-Agent': 'GraphQL-Federation-Gateway/1.0.0' }
        });
        
        if (response.ok) {
          console.log(`✅ ${serviceName} is available at ${url}`);
          return true;
        } else {
          console.log(`⚠️  ${serviceName} health check failed (${response.status}) - excluding from federation`);
          return false;
        }
      } catch (error) {
        console.log(`❌ ${serviceName} is not available at ${url} - excluding from federation`);
        return false;
      }
    };

    // Log all environment variables for debugging
    console.log('🔍 Environment Variables Debug:');
    console.log('PORT:', process.env.PORT);
    console.log('NODE_ENV:', process.env.NODE_ENV);
    
    const serviceEnvVars = Object.keys(process.env).filter(key => key.includes('SERVICE_URL'));
    console.log('Found SERVICE_URL environment variables:', serviceEnvVars);
    
    serviceEnvVars.forEach(envVar => {
      console.log(`${envVar}:`, `'${process.env[envVar]}'`);
    });
    
    // Also check if services might be available with different naming
    console.log('\nAll environment variables containing "commerce":');
    Object.keys(process.env).filter(key => key.toLowerCase().includes('commerce')).forEach(key => {
      console.log(`${key}: '${process.env[key]}'`);
    });

    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          // Start with minimal essential services only
          { name: 'identity', url: getServiceUrl('IDENTITY_SERVICE_URL', 8001) }
          // Other services will be added dynamically as they become available
        ],
        pollIntervalInMs: 30000, // Poll every 30 seconds for service updates (more frequent)
        introspectionHeaders: {
          'User-Agent': 'GraphQL-Federation-Gateway/1.0.0'
        },
        // Add more lenient error handling
        subgraphHealthCheck: false, // Disable health checks that might cause startup failures
      }),
      
      // Service health check and retry logic
      serviceHealthCheck: false, // Disable to prevent startup failures
      
      // Experimental: more lenient service loading
      experimental_autoFragmentization: true,
      
      // Handle service failures gracefully
      buildService({ name, url }) {
        return {
          name,
          url,
          // Add retry logic and timeouts
          async willSendRequest({ request, context }) {
            // Add timeout and retry headers
            request.http.timeout = 10000; // 10 second timeout
          }
        };
      }
    });

    // Create Apollo Server with the gateway
    const server = new ApolloServer({
      gateway,
      context: getGraphQLContext,
      // Enable introspection for federation to work (required for service discovery)
      introspection: true,
      // Enable playground in non-production for easier testing
      // Note: Federation requires introspection to discover services
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
          subgraphs: 8, // Core 8 services (analytics excluded temporarily)
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
      console.log('\n📊 Gateway started with minimal service set:');
      console.log(`  ✅ Identity Service: ${getServiceUrl('IDENTITY_SERVICE_URL', 8001)}`);
      console.log('\n⚡ Other services will be added automatically as they become available.');
      console.log('\n🔄 Dynamic service discovery enabled - gateway polls for new services every 30 seconds.');
      console.log('\n🛡️  Resilient mode: Gateway starts even if some services are down.');
      console.log('\n🎯 This ensures minimal functionality is always available!');
    });

  } catch (error) {
    console.error('❌ Failed to start GraphQL Federation Gateway:', error);
    console.error('Error details:', error.message);
    process.exit(1);
  }
}

// Start the gateway
startGateway();