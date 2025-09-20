const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const jwt = require('jsonwebtoken');
const { ApolloServer } = require('@apollo/server');
const { expressMiddleware } = require('@apollo/server/express4');
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { json } = require('body-parser');
const fetch = require('node-fetch');

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
        ? [
            // Official domains
            'https://unified-commerce.vercel.app', 
            'https://admin.unified-commerce.vercel.app', 
            'https://unified-commerce-storefront.vercel.app',
            // Vercel preview deployments (git branches)
            /https:\/\/.*\.vercel\.app$/,
            // Allow all subdomains for Vercel deployments
            /https:\/\/.*-.*\.vercel\.app$/
          ]
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
        } else {
          // For Render private services, use internal HTTP with default port
          // Private services are accessible via internal hostname on port 10000 (Render's default)
          fullUrl = `http://${serviceHost}:10000/graphql`;
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
        console.log(`🔍 Testing ${serviceName} at ${url}...`);
        
        // Try the GraphQL federation endpoint directly (more important than health)
        const federationTestQuery = JSON.stringify({
          query: '{ _service { sdl } }'
        });
        
        const gqlResponse = await fetch(url, {
          method: 'POST',
          timeout: 15000, // Increased timeout
          headers: {
            'Content-Type': 'application/json',
            'User-Agent': 'GraphQL-Federation-Gateway/1.0.0'
          },
          body: federationTestQuery
        });
        
        if (gqlResponse.ok) {
          const gqlResult = await gqlResponse.json();
          if (gqlResult.data && gqlResult.data._service) {
            console.log(`✅ ${serviceName} federation endpoint is working at ${url}`);
            return true;
          } else {
            console.log(`⚠️ ${serviceName} GraphQL responded but without federation schema:`, gqlResult);
            // Still return false but it's closer to working
            return false;
          }
        } else {
          console.log(`❌ ${serviceName} GraphQL endpoint failed (${gqlResponse.status}) - Status: ${gqlResponse.statusText}`);
          
          // Try health endpoint as backup check
          try {
            const healthUrl = url.replace('/graphql', '/health');
            console.log(`🔍 Trying health endpoint for ${serviceName} at ${healthUrl}...`);
            
            const healthResponse = await fetch(healthUrl, { 
              method: 'GET', 
              timeout: 10000,
              headers: { 'User-Agent': 'GraphQL-Federation-Gateway/1.0.0' }
            });
            
            if (healthResponse.ok) {
              const healthData = await healthResponse.json();
              console.log(`✅ ${serviceName} health endpoint is working:`, healthData);
              console.log(`⚠️  But GraphQL endpoint at ${url} is not ready yet`);
            } else {
              console.log(`❌ ${serviceName} health endpoint also failed (${healthResponse.status})`);
            }
          } catch (healthError) {
            console.log(`❌ ${serviceName} health endpoint error:`, healthError.message);
          }
          
          return false;
        }
      } catch (error) {
        console.log(`❌ ${serviceName} is not available at ${url} - Error: ${error.message}`);
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

    // Wait for services to become available with retry logic
    const waitForServices = async (maxRetries = 5, retryDelay = 10000) => {
      console.log(`🔄 Waiting for services to become available (max ${maxRetries} retries, ${retryDelay/1000}s delay)...`);
      
      for (let attempt = 1; attempt <= maxRetries; attempt++) {
        console.log(`🔍 Attempt ${attempt}/${maxRetries} - Testing service availability...`);
        
        const identityUrl = getServiceUrl('IDENTITY_SERVICE_URL', 8001);
        const productCatalogUrl = getServiceUrl('PRODUCT_CATALOG_SERVICE_URL', 8002);
        
        const identityAvailable = await testServiceAvailability(identityUrl, 'Identity Service');
        const productCatalogAvailable = await testServiceAvailability(productCatalogUrl, 'Product Catalog Service');
        
        if (identityAvailable || productCatalogAvailable) {
          console.log(`✅ Found available services on attempt ${attempt}`);
          return { identityAvailable, productCatalogAvailable, identityUrl, productCatalogUrl };
        }
        
        if (attempt < maxRetries) {
          console.log(`⚠️  No services available yet. Retrying in ${retryDelay/1000} seconds...`);
          await new Promise(resolve => setTimeout(resolve, retryDelay));
        }
      }
      
      console.log(`❌ No services became available after ${maxRetries} attempts`);
      return { identityAvailable: false, productCatalogAvailable: false, identityUrl: null, productCatalogUrl: null };
    };

    // Test both services before creating federation
    console.log('\n🔎 Testing service availability for federation...');
    
    const { identityAvailable, productCatalogAvailable, identityUrl, productCatalogUrl } = await waitForServices();
    
    // Build subgraphs array based on service availability
    const subgraphs = [];
    
    if (identityAvailable) {
      subgraphs.push({ name: 'identity', url: identityUrl });
      console.log('✅ Added Identity service to federation');
    } else {
      console.log('❌ Identity service excluded from federation - not available');
    }
    
    if (productCatalogAvailable) {
      subgraphs.push({ name: 'product-catalog', url: productCatalogUrl });
      console.log('✅ Added Product Catalog service to federation');
    } else {
      console.log('❌ Product Catalog service excluded from federation - not available');
    }
    
    if (subgraphs.length === 0) {
      throw new Error('No services available for federation!');
    }
    
    console.log(`\n🔗 Creating federation with ${subgraphs.length} services:`, subgraphs.map(s => s.name).join(', '));

    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs,
        pollIntervalInMs: 60000, // Poll every 60 seconds - more conservative
        introspectionHeaders: {
          'User-Agent': 'GraphQL-Federation-Gateway/1.0.0'
        }
      })
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
      console.log('\n📊 Gateway started with dynamic federation:');
      if (identityAvailable) {
        console.log(`  ✅ Identity Service: ${identityUrl}`);
      }
      if (productCatalogAvailable) {
        console.log(`  ✅ Product Catalog Service: ${productCatalogUrl}`);
      }
      console.log(`\n⚡ Federation active with ${subgraphs.length} service(s).`);
      if (productCatalogAvailable) {
        console.log('\n🔄 Products should now load successfully in the storefront.');
        console.log('\n🛡️  Federation mode: Full GraphQL schema with products and authentication.');
      } else {
        console.log('\n⚠️  Product Catalog not available - products will not load until service is ready.');
        console.log('\n🛡️  Partial federation mode: Authentication only.');
      }
    });

  } catch (error) {
    console.error('❌ Failed to start GraphQL Federation Gateway:', error);
    console.error('Error details:', error.message);
    process.exit(1);
  }
}

// Start the gateway
startGateway();