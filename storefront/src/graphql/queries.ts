import { gql } from '@apollo/client';

// This file contains the GraphQL queries, mutations, and TypeScript interfaces
// for the storefront. It is aligned with the latest backend schema.

//---------------------------
// INTERFACES
//---------------------------

export interface Product {
  id: string;
  title: string;
  description?: string;
  price: number;
  status: 'ACTIVE' | 'DRAFT' | 'ARCHIVED';
  imageUrl?: string; // Assuming the schema might provide a primary image URL
  variants?: ProductVariant[];
  categories?: Category[];
  createdAt: string;
}

export interface ProductVariant {
  id: string;
  sku: string;
  price: number;
  inventory?: {
    quantity: number;
  };
}

export interface Category {
  id: string;
  name: string;
}

//---------------------------
// QUERIES
//---------------------------

export const GET_PRODUCTS = gql`
  query GetProducts($filter: ProductFilter) {
    products(filter: $filter) {
      id
      title
      description
      price
      status
      createdAt
      # Assuming a primary image URL might be available on the product itself
      # If not, we might need to get it from variants or another source
      # imageUrl 
      variants {
        id
        price
        inventory {
          quantity
        }
      }
      categories {
        id
        name
      }
    }
  }
`;

export const GET_PRODUCT_BY_ID = gql`
  query GetProductById($id: ID!) {
    product(id: $id) {
      id
      title
      description
      price
      status
      createdAt
      variants {
        id
        sku
        price
        inventory {
          quantity
        }
      }
      categories {
        id
        name
      }
    }
  }
`;

//---------------------------
// MUTATIONS
//---------------------------

// Cart-related mutations would go here
