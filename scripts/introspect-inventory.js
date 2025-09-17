const https = require('https');
const http = require('http');

// Introspect the InventoryItem type to see what fields are available
const query = JSON.stringify({
  query: `
    query {
      __type(name: "InventoryItem") {
        name
        fields {
          name
          type {
            name
            kind
          }
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
    const response = JSON.parse(data);
    
    if (response.data && response.data.__type) {
      console.log('✅ InventoryItem type fields:');
      const fields = response.data.__type.fields;
      fields.forEach(field => {
        console.log(`  - ${field.name}: ${field.type.name || field.type.kind}`);
      });
    }
    
    if (response.errors) {
      console.log('\n❌ Errors during introspection:');
      response.errors.forEach(error => {
        console.log(`  - ${error.message}`);
      });
    }
  });
});

req.on('error', (error) => {
  console.error('❌ Error introspecting InventoryItem type:', error.message);
});

req.write(query);
req.end();