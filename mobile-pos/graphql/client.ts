import { ApolloClient, InMemoryCache, HttpLink, ApolloLink, from } from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { localStorage } from '../utils/storage';

// Error handling link
const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    graphQLErrors.forEach(({ message, locations, path }) => {
      console.error(
        `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
      );
    });
  }
  
  if (networkError) {
    console.error(`[Network error]: ${networkError}`);
  }
});

// HTTP link to GraphQL Federation Gateway
const httpLink = new HttpLink({
  uri: 'https://unified-commerce-gateway.onrender.com/graphql', // Production GraphQL Federation Gateway
  credentials: 'include', // Include cookies/session tokens
});

// Authentication middleware
const authLink = new ApolloLink((operation, forward) => {
  // Get token from storage
  const token = localStorage.getItem('authToken');
  
  // Add authorization header if token exists
  if (token) {
    operation.setContext({
      headers: {
        authorization: `Bearer ${token}`,
      },
    });
  }
  
  return forward(operation);
});

// Create Apollo Client
const client = new ApolloClient({
  link: from([errorLink, authLink, httpLink]),
  cache: new InMemoryCache({
    typePolicies: {
      Product: {
        keyFields: ['id'],
      },
      Cart: {
        keyFields: ['id'],
      },
      Payment: {
        keyFields: ['id'],
      },
      // Add more type policies as needed
    },
  }),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
    },
  },
});

export default client;