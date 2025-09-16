import React from 'react';
import { render } from '@testing-library/react-native';
import POSScreen from '../app/pos/index';
import CheckoutScreen from '../app/pos/checkout';

// Mock Apollo Client
jest.mock('@apollo/client', () => ({
  useQuery: () => ({
    data: {
      products: [],
      cartBySession: {
        lineItems: []
      }
    },
    loading: false,
    error: undefined,
    refetch: jest.fn()
  }),
  useMutation: () => [jest.fn(), { loading: false, error: undefined }],
  useApolloClient: () => ({
    query: jest.fn().mockResolvedValue({ data: { cartBySession: null } })
  })
}));

// Mock localStorage
jest.mock('../utils/storage', () => ({
  localStorage: {
    getItem: jest.fn().mockReturnValue('test-session-id'),
    setItem: jest.fn(),
    removeItem: jest.fn()
  }
}));

describe('POS Screens', () => {
  it('should render POS screen without crashing', () => {
    const { getByText } = render(<POSScreen />);
    expect(getByText('Retail OS Point of Sale')).toBeTruthy();
  });

  it('should render Checkout screen without crashing', () => {
    const { getByText } = render(<CheckoutScreen />);
    expect(getByText('Checkout')).toBeTruthy();
  });
});