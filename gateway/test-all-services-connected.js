const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');

async function testAllServicesConnected() {
  console.log('🔍 Testing GraphQL Federation Gateway service connections...');
  
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
    })
  });

  try {
    console.log('🔄 Attempting to compose schema from all services...');
    const { supergraphSdl } = await gateway.load();
    
    if (supergraphSdl) {
      console.log('✅ SUCCESS: All services successfully connected to GraphQL Federation Gateway!');
      console.log('   Schema composition completed successfully.');
      
      // Count the number of services by looking for service definitions
      const serviceCount = (supergraphSdl.match(/schema.*?subgraph.*?name/g) || []).length;
      console.log(`📊 Services connected: ${serviceCount}/8`);
      
      if (serviceCount === 8) {
        console.log('🎉 ALL SERVICES SUCCESSFULLY CONNECTED!');
      } else {
        console.log(`⚠️  Only ${serviceCount}/8 services connected. Some services may be offline.`);
      }
    } else {
      console.log('❌ FAILED: Schema composition returned empty result');
    }
  } catch (error) {
    console.log('❌ FAILED: Error connecting to services');
    console.log('   Error details:', error.message);
    
    // Try to identify which service is failing
    if (error.message.includes('8001')) {
      console.log('   🔍 Issue with Identity Service (port 8001)');
    } else if (error.message.includes('8002')) {
      console.log('   🔍 Issue with Cart Service (port 8002)');
    } else if (error.message.includes('8003')) {
      console.log('   🔍 Issue with Order Service (port 8003)');
    } else if (error.message.includes('8004')) {
      console.log('   🔍 Issue with Payment Service (port 8004)');
    } else if (error.message.includes('8005')) {
      console.log('   🔍 Issue with Inventory Service (port 8005)');
    } else if (error.message.includes('8006')) {
      console.log('   🔍 Issue with Product Catalog Service (port 8006)');
    } else if (error.message.includes('8007')) {
      console.log('   🔍 Issue with Promotions Service (port 8007)');
    } else if (error.message.includes('8008')) {
      console.log('   🔍 Issue with Merchant Account Service (port 8008)');
    }
  }
  
  // Clean up
  await gateway.stop();
}

// Run the test
testAllServicesConnected();