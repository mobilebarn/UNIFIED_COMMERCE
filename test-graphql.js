const { ApolloServer } = require('@apollo/server');
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { readFileSync } = require('fs');

async function testFederation() {
  console.log('Testing GraphQL Federation Gateway...');
  
  try {
    // Create Apollo Gateway with service definitions
    const gateway = new ApolloGateway({
      supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
          { name: 'identity', url: 'http://localhost:8001/graphql' },
          { name: 'inventory', url: 'http://localhost:8005/graphql' },
          { name: 'order', url: 'http://localhost:8003/graphql' },
          { name: 'payment', url: 'http://localhost:8004/graphql' }
        ],
      })
    });

    // Try to initialize the gateway
    const { supergraphSdl } = await gateway.load();
    console.log('✅ GraphQL Federation Gateway initialized successfully');
    console.log('✅ All services are properly federated');
    
    // Show a portion of the schema to verify it's working
    const schemaLines = supergraphSdl.split('\n');
    console.log('\n📋 Sample of the federated schema:');
    console.log(schemaLines.slice(0, 20).join('\n'));
    
    process.exit(0);
  } catch (error) {
    console.error('❌ Error testing GraphQL Federation Gateway:', error.message);
    process.exit(1);
  }
}

testFederation();