// Simple script to verify GraphQL endpoint is accessible
const fetch = require('node-fetch');

async function verifyGraphQLEndpoint() {
  const endpoint = 'http://localhost:4000/graphql';
  
  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query: `
          query {
            __typename
          }
        `,
      }),
    });
    
    const result = await response.json();
    
    if (response.ok && result.data) {
      console.log('✅ GraphQL endpoint is accessible');
      console.log('Endpoint:', endpoint);
      return true;
    } else {
      console.log('❌ GraphQL endpoint returned an error');
      console.log('Status:', response.status);
      console.log('Error:', result.errors);
      return false;
    }
  } catch (error) {
    console.log('❌ Failed to connect to GraphQL endpoint');
    console.log('Error:', error.message);
    return false;
  }
}

// Run the verification
verifyGraphQLEndpoint().then(success => {
  if (success) {
    console.log('\n🎉 GraphQL integration is ready for use!');
  } else {
    console.log('\n⚠️  Please ensure the GraphQL Federation Gateway is running on port 4000');
  }
});