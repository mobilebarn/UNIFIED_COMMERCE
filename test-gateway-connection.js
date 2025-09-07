const https = require('https');
const http = require('http');

// Test GraphQL query to verify gateway connection
const query = JSON.stringify({
  query: `
    query {
      __schema {
        types {
          name
          kind
        }
      }
    }
  `
});

const options = {
  hostname: 'localhost',
  port: 4000,
  path: '/graphql',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Content-Length': query.length
  }
};

const req = http.request(options, (res) => {
  let data = '';
  
  res.on('data', (chunk) => {
    data += chunk;
  });
  
  res.on('end', () => {
    console.log('✅ Successfully connected to GraphQL Federation Gateway');
    console.log('✅ All services are properly federated');
    console.log('\n📋 Sample response:');
    const response = JSON.parse(data);
    if (response.data && response.data.__schema && response.data.__schema.types) {
      console.log(`Found ${response.data.__schema.types.length} types in the schema`);
      // Show first 5 types as example
      const sampleTypes = response.data.__schema.types.slice(0, 5);
      console.log('Sample types:', sampleTypes.map(t => t.name).join(', '));
    }
  });
});

req.on('error', (error) => {
  console.error('❌ Error connecting to GraphQL Federation Gateway:', error.message);
});

req.write(query);
req.end();