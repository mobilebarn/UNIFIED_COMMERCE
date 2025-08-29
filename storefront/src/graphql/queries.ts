import { gql } from '@apollo/client';

// Product Queries
export const GET_PRODUCTS = gql`
  query GetProducts($limit: Int, $offset: Int, $search: String) {
    products(limit: $limit, offset: $offset, search: $search) {
      id
      name
      description
      price
      image
      category
      inventory {
        quantity
        inStock
      }
      createdAt
      updatedAt
    }
  }
`;

export const GET_PRODUCT = gql`
  query GetProduct($id: ID!) {
    product(id: $id) {
      id
      name
      description
      price
      image
      category
      inventory {
        quantity
        inStock
      }
      createdAt
      updatedAt
    }
  }
`;

// User Queries
export const GET_USER_PROFILE = gql`
  query GetUserProfile {
    userProfile {
      id
      email
      firstName
      lastName
      role
      createdAt
    }
  }
`;

export const GET_USER_ORDERS = gql`
  query GetUserOrders {
    userOrders {
      id
      status
      total
      items {
        id
        productId
        productName
        quantity
        price
      }
      createdAt
      updatedAt
    }
  }
`;

// Authentication Mutations
export const LOGIN_MUTATION = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      token
      user {
        id
        email
        firstName
        lastName
        role
      }
    }
  }
`;

export const REGISTER_MUTATION = gql`
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
      token
      user {
        id
        email
        firstName
        lastName
        role
      }
    }
  }
`;

// Order Mutations
export const CREATE_ORDER = gql`
  mutation CreateOrder($input: OrderInput!) {
    createOrder(input: $input) {
      id
      status
      total
      items {
        id
        productId
        productName
        quantity
        price
      }
      createdAt
    }
  }
`;

// Cart Mutations (if using server-side cart)
export const ADD_TO_CART = gql`
  mutation AddToCart($productId: ID!, $quantity: Int!) {
    addToCart(productId: $productId, quantity: $quantity) {
      id
      items {
        id
        productId
        productName
        quantity
        price
      }
      total
    }
  }
`;

export const REMOVE_FROM_CART = gql`
  mutation RemoveFromCart($productId: ID!) {
    removeFromCart(productId: $productId) {
      id
      items {
        id
        productId
        productName
        quantity
        price
      }
      total
    }
  }
`;
