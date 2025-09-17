const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const jwt = require('jsonwebtoken');
const { ApolloServer } = require('@apollo/server');
const { expressMiddleware } = require('@apollo/server/express4');
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { json } = require('body-parser');

async function startGateway() {
  try {
    console.log('🚀 Starting GraphQL Federation Gateway (Minimal Setup)...');

    const app = express();

    app.use(cors({
      origin: [
        'https://unified-commerce.vercel.app',
        'https://admin-panel-igp522vr5-crypticogs-projects.vercel.app',
        'http://localhost:3000',
        'http://localhost:5173'
      ],
      credentials: true,
    }));

    // Add Order Service to the working minimal setup
    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: 'http://localhost:8001/graphql' },
          { name: 'product', url: 'http://localhost:8006/graphql' },
          { name: 'order', url: 'http://localhost:8004/graphql' }
        ],
      })
    });

    const server = new ApolloServer({
      gateway,
      introspection: true,
    });

    await server.start();

    app.get('/health', (req, res) => {
      res.json({
        service: 'graphql-federation-gateway',
        status: 'healthy',
        time: new Date().toISOString(),
        federation: {
          subgraphs: 3,
          active: true
        },
        services: {
          identity: 'http://localhost:8001/graphql',
          product: 'http://localhost:8006/graphql',
          order: 'http://localhost:8004/graphql'
        },
        frontend_apps: {
          storefront: 'https://unified-commerce.vercel.app',
          admin: 'https://admin-panel-igp522vr5-crypticogs-projects.vercel.app'
        }
      });
    });

    app.use('/graphql', json(), expressMiddleware(server));
    app.get('/', (req, res) => res.redirect('/graphql'));

    const PORT = 4000;
    app.listen(PORT, () => {
      console.log(`✅ GraphQL Federation Gateway running at http://localhost:${PORT}/graphql`);
      console.log(`🎮 GraphQL Playground: http://localhost:${PORT}/graphql`);
      console.log('\n📊 Active Services:');
      console.log('  ✅ Identity Service (Users/Auth)');
      console.log('  ✅ Product Catalog Service (Products)');
      console.log('  ✅ Order Service (Orders/Fulfillment)');
      console.log('\n🌐 Test Integration:');
      console.log('  📱 Storefront: https://unified-commerce.vercel.app');
      console.log('  🏢 Admin Panel: https://admin-panel-igp522vr5-crypticogs-projects.vercel.app');
      console.log('\n🎯 Ready for basic frontend-backend testing!');
    });

  } catch (error) {
    console.error('❌ Gateway failed:', error.message);
  }
}

startGateway();