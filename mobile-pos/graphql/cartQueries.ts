import { gql } from '@apollo/client';

// Query to get a cart by ID
export const GET_CART = gql`
  query GetCart($id: ID!) {
    cart(id: $id) {
      id
      sessionId
      customerId
      merchantId
      status
      currency
      subtotalPrice
      totalTax
      totalShipping
      totalDiscount
      totalPrice
      lineItems {
        id
        productId
        productVariantId
        name
        sku
        quantity
        price
        linePrice
        productImage
      }
      createdAt
      updatedAt
    }
  }
`;

// Query to get a cart by session ID
export const GET_CART_BY_SESSION = gql`
  query GetCartBySession($sessionId: String!) {
    cartBySession(sessionId: $sessionId) {
      id
      sessionId
      customerId
      merchantId
      status
      currency
      subtotalPrice
      totalTax
      totalShipping
      totalDiscount
      totalPrice
      lineItems {
        id
        productId
        productVariantId
        name
        sku
        quantity
        price
        linePrice
        productImage
      }
      createdAt
      updatedAt
    }
  }
`;

// Mutation to create a new cart
export const CREATE_CART = gql`
  mutation CreateCart($input: CreateCartInput!) {
    createCart(input: $input) {
      id
      sessionId
      customerId
      merchantId
      status
      currency
      subtotalPrice
      totalTax
      totalShipping
      totalDiscount
      totalPrice
      lineItems {
        id
        productId
        productVariantId
        name
        sku
        quantity
        price
        linePrice
        productImage
      }
      createdAt
      updatedAt
    }
  }
`;

// Mutation to add an item to the cart
export const ADD_CART_LINE_ITEM = gql`
  mutation AddCartLineItem($input: AddLineItemInput!) {
    addCartLineItem(input: $input) {
      id
      cartId
      productId
      productVariantId
      name
      sku
      quantity
      price
      linePrice
      productImage
    }
  }
`;

// Mutation to update a cart line item
export const UPDATE_CART_LINE_ITEM = gql`
  mutation UpdateCartLineItem($id: ID!, $input: UpdateLineItemInput!) {
    updateCartLineItem(id: $id, input: $input) {
      id
      cartId
      productId
      productVariantId
      name
      sku
      quantity
      price
      linePrice
      productImage
    }
  }
`;

// Mutation to remove an item from the cart
export const REMOVE_CART_LINE_ITEM = gql`
  mutation RemoveCartLineItem($id: ID!) {
    removeCartLineItem(id: $id)
  }
`;

// Mutation to clear the cart
export const CLEAR_CART = gql`
  mutation ClearCart($cartId: ID!) {
    clearCart(cartId: $cartId)
  }
`;