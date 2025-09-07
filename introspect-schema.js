const https = require('https');
const http = require('http');

// Introspect the schema to see what's actually available
const query = JSON.stringify({
  query: `
    query {
      __schema {
        queryType {
          name
          fields {
            name
            type {
              name
              kind
              ofType {
                name
                kind
              }
            }
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
    
    if (response.data && response.data.__schema && response.data.__schema.queryType) {
      console.log('✅ Available query fields in the federated schema:');
      const fields = response.data.__schema.queryType.fields;
      fields.forEach(field => {
        const typeName = field.type.name || (field.type.ofType ? field.type.ofType.name : 'Unknown');
        console.log(`  - ${field.name}: ${typeName}`);
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
  console.error('❌ Error introspecting schema:', error.message);
});

req.write(query);
req.end();