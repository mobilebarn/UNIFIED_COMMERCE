const { ApolloServer } = require('@apollo/server');
const { expressMiddleware } = require('@apollo/server/express4');
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const jwt = require('jsonwebtoken');
const fetch = require('node-fetch');

// Authentication context function
function getAuthContext({ req }) {
  const token = req.headers.authorization?.replace('Bearer ', '');
  
  if (!token) {
    return { user: null, isAuthenticated: false };
  }

  try {
    // Use same JWT secret as Identity service
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
    return { user: null, isAuthenticated: false };
  }
}

// Custom data source to forward authentication headers
class AuthenticatedDataSource extends require('@apollo/gateway').RemoteGraphQLDataSource {
  willSendRequest({ request, context }) {
    if (context.user) {
      request.http.headers.set('user-id', context.userId);
      request.http.headers.set('user-email', context.email);
      request.http.headers.set('user-roles', JSON.stringify(context.roles));
    }
  }
}

async function startGateway() {
  try {
    // First, verify that the Identity service is available
    console.log('🔍 Checking Identity service availability...');
    
    const checkService = async (url, retries = 5) => {
      for (let i = 0; i < retries; i++) {
        try {
          const response = await fetch(`${url.replace('/graphql', '')}/health`);
          if (response.ok) {
            console.log(`✅ Service at ${url} is available`);
            return true;
          }
        } catch (error) {
          console.log(`⏳ Attempt ${i + 1}/${retries}: Service at ${url} not ready, retrying in 2s...`);
          if (i < retries - 1) await new Promise(resolve => setTimeout(resolve, 2000));
        }
      }
      throw new Error(`Service at ${url} is not available after ${retries} attempts`);
    };

    await checkService('http://localhost:8080/graphql');

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
      : ['http://localhost:3000', 'http://localhost:3001', 'http://localhost:4000'],
    credentials: true,
  }));

  // Create the gateway with retry logic and better error handling
  console.log('🔧 Setting up Apollo Federation Gateway...');
  
  let gateway;
  try {
    gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: 'http://localhost:8080/graphql' },
        ],
        introspectionHeaders: {
          'User-Agent': 'Apollo-Gateway/2.0',
        },
      }),
      buildService({ url }) {
        return new AuthenticatedDataSource({ url });
      },
    });
  } catch (gatewayError) {
    console.error('❌ Gateway creation failed:', gatewayError);
    throw gatewayError;
  }

  // Create the server with enhanced configuration
  console.log('🚀 Starting Apollo Server...');
  let server;
  try {
    server = new ApolloServer({
      gateway,
      introspection: true,
      formatError: (error) => {
        // Log errors for monitoring
        console.error('GraphQL Error:', {
          message: error.message,
          path: error.path,
          locations: error.locations,
          timestamp: new Date().toISOString(),
        });

        // Don't expose internal errors in production
        if (process.env.NODE_ENV === 'production' && !error.message.startsWith('Authentication')) {
          return new Error('Internal server error');
        }

        return error;
      },
    });
  } catch (serverError) {
    console.error('❌ Server creation failed:', serverError);
    throw serverError;
  }

  // Start the server
  console.log('⚡ Starting Apollo Server...');
  try {
    await server.start();
    console.log('✅ Apollo Server started successfully');
  } catch (startError) {
    console.error('❌ Server start failed:', startError);
    throw startError;
  }

  // Apply the Apollo GraphQL middleware
  app.use('/graphql', expressMiddleware(server, {
    context: getAuthContext,
  }));

  // Health check endpoint
  app.get('/health', (req, res) => {
    res.json({ 
      status: 'healthy', 
      timestamp: new Date().toISOString(),
      services: {
        gateway: 'running',
        identity: 'connected',
      }
    });
  });

  // API info endpoint
  app.get('/', (req, res) => {
    res.json({
      name: 'Unified Commerce GraphQL Gateway',
      version: '1.0.0',
      endpoints: {
        graphql: '/graphql',
        health: '/health',
      },
      documentation: 'Visit /graphql for GraphQL Playground'
    });
  });

  // Start the HTTP server
  const PORT = process.env.PORT || 4000;
  app.listen(PORT, () => {
    console.log(`🚀 GraphQL Gateway ready at http://localhost:${PORT}/graphql`);
    console.log(`📊 Health check available at http://localhost:${PORT}/health`);
    console.log(`🎮 GraphQL Explorer at http://localhost:${PORT}/graphql`);
    console.log(`📖 API Info at http://localhost:${PORT}/`);
  });

  } catch (error) {
    console.error('Failed to start gateway:', error);
    process.exit(1);
  }
}

startGateway();

// Add global error handlers to catch unhandled rejections
process.on('unhandledRejection', (reason, promise) => {
  console.error('❌ Unhandled Promise Rejection:', reason);
  console.error('Promise:', promise);
});

process.on('uncaughtException', (error) => {
  console.error('❌ Uncaught Exception:', error);
  process.exit(1);
});