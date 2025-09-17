const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { ApolloServer } = require('apollo-server-express');
const express = require('express');

const app = express();

// Enable CORS for all origins
app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization');
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  if (req.method === 'OPTIONS') {
    return res.sendStatus(200);
  }
  next();
});

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: 'identity', url: 'http://localhost:8001/graphql' },
      { name: 'merchant', url: 'http://localhost:8002/graphql' },
      { name: 'order', url: 'http://localhost:8004/graphql' },
      { name: 'product', url: 'http://localhost:8006/graphql' },
      { name: 'promotions', url: 'http://localhost:8007/graphql' },
      { name: 'analytics', url: 'http://localhost:8008/graphql' },
      { name: 'cart', url: 'http://localhost:8080/graphql' }
    ],
  }),
  debug: true,
});

async function startServer() {
  try {
    console.log('Starting GraphQL Federation Gateway...');
    console.log('Available services:');
    console.log('- Identity Service: http://localhost:8001/graphql');
    console.log('- Merchant Service: http://localhost:8002/graphql');
    console.log('- Order Service: http://localhost:8004/graphql');
    console.log('- Product Service: http://localhost:8006/graphql');
    console.log('- Promotions Service: http://localhost:8007/graphql');
    console.log('- Analytics Service: http://localhost:8008/graphql');
    console.log('- Cart Service: http://localhost:8080/graphql');

    const server = new ApolloServer({ 
      gateway,
      introspection: true,
      playground: true,
      cors: {
        origin: '*',
        credentials: true,
      }
    });

    await server.start();
    server.applyMiddleware({ app, cors: false });

    const PORT = process.env.PORT || 4000;
    
    app.listen(PORT, () => {
      console.log(`🚀 GraphQL Federation Gateway ready at http://localhost:${PORT}${server.graphqlPath}`);
      console.log(`🎮 GraphQL Playground available at http://localhost:${PORT}${server.graphqlPath}`);
      console.log('\nFederated Services Status:');
      console.log('✅ Identity Service (8001)');
      console.log('✅ Merchant Service (8002)');
      console.log('✅ Order Service (8004)');
      console.log('✅ Product Service (8006)');
      console.log('✅ Promotions Service (8007)');
      console.log('✅ Analytics Service (8008)');
      console.log('✅ Cart Service (8080)');
    });

  } catch (error) {
    console.error('Failed to start GraphQL Federation Gateway:', error);
    console.error('Error details:', error.message);
    if (error.errors) {
      error.errors.forEach((err, index) => {
        console.error(`Error ${index + 1}:`, err.message);
      });
    }
    process.exit(1);
  }
}

startServer().catch(console.error);