const axios = require('axios');

// Test GraphQL query to fetch products (simplified)
const query = `
  query {
    products {
      id
      title
    }
  }
`;

async function testGraphQLData() {
  try {
    console.log('Testing GraphQL data fetch from Federation Gateway...');
    
    const response = await axios.post('http://localhost:4000/graphql', {
      query: query
    }, {
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    console.log('✅ GraphQL query successful');
    console.log('📋 Response data:');
    console.log(JSON.stringify(response.data, null, 2));
    
    if (response.data.data && response.data.data.products) {
      console.log(`\n🎉 Found ${response.data.data.products.length} products`);
      if (response.data.data.products.length > 0) {
        const product = response.data.data.products[0];
        console.log('Sample product:');
        console.log(`  ID: ${product.id}`);
        console.log(`  Title: ${product.title}`);
      }
    }
  } catch (error) {
    console.error('❌ Error testing GraphQL data:', error.message);
    if (error.response) {
      console.error('Response data:', JSON.stringify(error.response.data, null, 2));
    }
  }
}

testGraphQLData();