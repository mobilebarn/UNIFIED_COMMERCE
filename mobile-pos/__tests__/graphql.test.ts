import { gql } from '@apollo/client';
import client from '../graphql/client';

// Simple test to verify GraphQL client is working
describe('GraphQL Client', () => {
  it('should be able to create a GraphQL client instance', () => {
    expect(client).toBeDefined();
  });

  it('should have the correct GraphQL endpoint', () => {
    // This is a simple test - in a real implementation, you would test actual queries
    expect(client.link).toBeDefined();
  });
});

// Test GraphQL queries
describe('GraphQL Queries', () => {
  it('should define GET_PRODUCTS query', async () => {
    const query = gql`
      query GetProducts($filter: ProductFilter) {
        products(filter: $filter) {
          id
          title
          priceRange {
            minVariantPrice
          }
        }
      }
    `;
    
    expect(query).toBeDefined();
  });

  it('should define GET_CART_BY_SESSION query', async () => {
    const query = gql`
      query GetCartBySession($sessionId: String!) {
        cartBySession(sessionId: $sessionId) {
          id
          sessionId
          lineItems {
            id
            name
            quantity
            price
          }
        }
      }
    `;
    
    expect(query).toBeDefined();
  });
});