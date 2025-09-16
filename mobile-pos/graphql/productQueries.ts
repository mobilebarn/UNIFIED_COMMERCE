import { gql } from '@apollo/client';

// Query to get products with filtering
export const GET_PRODUCTS = gql`
  query GetProducts($filter: ProductFilter) {
    products(filter: $filter) {
      id
      title
      handle
      description
      featuredImage
      priceRange {
        minVariantPrice
      }
      tags
      variants {
        id
        title
        price
        sku
        barcode
        inventoryQuantity
      }
    }
  }
`;

// Query to get a single product by ID
export const GET_PRODUCT = gql`
  query GetProduct($id: ID!) {
    product(id: $id) {
      id
      title
      handle
      description
      featuredImage
      images {
        src
        altText
      }
      priceRange {
        minVariantPrice
      }
      variants {
        id
        title
        price
        sku
        barcode
        inventoryQuantity
      }
      tags
    }
  }
`;

// Query to search for products by name or barcode
export const SEARCH_PRODUCTS = gql`
  query SearchProducts($query: String!) {
    products(filter: { search: $query }) {
      id
      title
      handle
      featuredImage
      priceRange {
        minVariantPrice
      }
      variants {
        id
        title
        price
        sku
        barcode
        inventoryQuantity
      }
    }
  }
`;