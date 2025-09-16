import React from 'react';
import { ApolloProvider } from '@apollo/client';
import client from '../graphql/client';

interface ApolloProviderWrapperProps {
  children: React.ReactNode;
}

export default function ApolloProviderWrapper({ children }: ApolloProviderWrapperProps) {
  return (
    <ApolloProvider client={client}>
      {children}
    </ApolloProvider>
  );
}