#!/usr/bin/env node

const fetch = require('node-fetch');

async function testGateway() {
  console.log('🧪 Testing GraphQL Federation Gateway (Proxy Mode)\n');

  try {
    // Test health endpoint
    console.log('📊 Testing health endpoint...');
    const healthResponse = await fetch('http://localhost:4000/health');
    const healthData = await healthResponse.json();
    console.log('✅ Health check:', healthData);

    // Test GraphQL introspection
    console.log('\n🔍 Testing GraphQL introspection...');
    const introspectionQuery = `
      query {
        __schema {
          queryType {
            name
          }
          types {
            name
            kind
          }
        }
      }
    `;

    const graphqlResponse = await fetch('http://localhost:4000/graphql', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ query: introspectionQuery }),
    });

    const graphqlData = await graphqlResponse.json();

    if (graphqlData.errors) {
      console.log('❌ GraphQL Errors:', graphqlData.errors);
    } else {
      console.log('✅ GraphQL Schema introspection successful');
      console.log('📋 Query Type:', graphqlData.data?.__schema?.queryType?.name);
      console.log('📊 Available Types:', graphqlData.data?.__schema?.types?.length || 0);
    }

    // Test a simple query
    console.log('\n🔍 Testing simple query...');
    const simpleQuery = `
      query {
        __typename
      }
    `;

    const simpleResponse = await fetch('http://localhost:4000/graphql', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ query: simpleQuery }),
    });

    const simpleData = await simpleResponse.json();
    console.log('✅ Simple query result:', simpleData);

  } catch (error) {
    console.error('❌ Test failed:', error.message);
  }
}

testGateway();
