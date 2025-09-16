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
  description?: string;
  parentId?: string;
  parent?: Category;
  children?: Category[];
}

export interface User {
  id: string;
  email: string;
  username: string;
  firstName: string;
  lastName: string;
  isActive: boolean;
  createdAt: string;
}

export interface AuthPayload {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface Order {
  id: string;
  orderNumber: string;
  status: 'PENDING' | 'CONFIRMED' | 'PROCESSING' | 'SHIPPED' | 'DELIVERED' | 'CANCELLED';
  total: number;
  currency: string;
  createdAt: string;
  items?: OrderItem[];
  payments?: Payment[];
  fulfillments?: Fulfillment[];
}

export interface OrderItem {
  id: string;
  quantity: number;
  price: number;
  product?: Product;
  variant?: ProductVariant;
}

export interface Payment {
  id: string;
  status: 'PENDING' | 'COMPLETED' | 'FAILED';
  amount: number;
  method: string;
  transactionId?: string;
  createdAt: string;
}

export interface Fulfillment {
  id: string;
  status: 'PENDING' | 'SHIPPED' | 'DELIVERED';
  trackingNumber?: string;
  carrier?: string;
  createdAt: string;
}

export interface WishlistItem {
  id: string;
  productId: string;
  addedAt: string;
  product?: Product;
}

export interface Address {
  id: string;
  firstName: string;
  lastName: string;
  company?: string;
  address1: string;
  address2?: string;
  city: string;
  province: string;
  country: string;
  zip: string;
  phone?: string;
  isDefault: boolean;
}

export interface PaymentMethod {
  id: string;
  type: 'credit_card' | 'debit_card' | 'paypal' | 'bank_account';
  name: string;
  last4: string;
  expiryMonth?: number;
  expiryYear?: number;
  isDefault: boolean;
}

export interface Promotion {
  id: string;
  name: string;
  description?: string;
  type: string;
  status: string;
  discountType: string;
  discountValue: number;
  startDate: string;
  endDate?: string;
  applicableProducts?: string[];
  applicableCollections?: string[];
  applicableCustomers?: string[];
  priority: number;
  usageLimit?: number;
  usedCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface DiscountCode {
  id: string;
  promotionId: string;
  code: string;
  isActive: boolean;
  usageCount: number;
  usageLimit?: number;
  usageLimitPerCustomer?: number;
  startsAt?: string;
  endsAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Campaign {
  id: string;
  merchantId: string;
  name: string;
  description?: string;
  type: string;
  status: string;
  startDate: string;
  endDate?: string;
  budget?: number;
  goalType?: string;
  goalValue?: number;
  createdAt: string;
  updatedAt: string;
  promotions?: Promotion[];
}

// Add SearchSuggestion interface
export interface SearchSuggestion {
  id: string;
  title: string;
  type: 'PRODUCT' | 'CATEGORY' | 'BRAND';
  imageUrl?: string;
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

// Add a new query to get a category with its parent and children
export const GET_CATEGORY = gql`
  query GetCategory($id: ID!) {
    category(id: $id) {
      id
      name
      description
      parentId
      parent {
        id
        name
      }
      children {
        id
        name
        description
      }
    }
  }
`;

// Add a query to get categories with their children
export const GET_CATEGORIES = gql`
  query GetCategories($filter: CategoryFilter) {
    categories(filter: $filter) {
      id
      name
      description
      parentId
      children {
        id
        name
        description
      }
    }
  }
`;

export const GET_CURRENT_USER = gql`
  query GetCurrentUser {
    currentUser {
      id
      email
      username
      firstName
      lastName
      isActive
      createdAt
    }
  }
`;

export const GET_ORDERS = gql`
  query GetOrders($filter: OrderFilter) {
    orders(filter: $filter) {
      id
      orderNumber
      status
      total
      currency
      createdAt
      items {
        id
        quantity
        price
        product {
          id
          title
        }
      }
      payments {
        id
        status
        amount
        method
      }
      fulfillments {
        id
        status
        trackingNumber
        carrier
      }
    }
  }
`;

export const GET_WISHLIST = gql`
  query GetWishlist {
    wishlist {
      id
      productId
      addedAt
      product {
        id
        title
        price
        description
        imageUrl
        variants {
          id
          price
          inventory {
            quantity
          }
        }
      }
    }
  }
`;

export const GET_ADDRESSES = gql`
  query GetAddresses {
    addresses {
      id
      firstName
      lastName
      company
      address1
      address2
      city
      province
      country
      zip
      phone
      isDefault
    }
  }
`;

export const GET_PAYMENT_METHODS = gql`
  query GetPaymentMethods {
    paymentMethods {
      id
      type
      name
      last4
      expiryMonth
      expiryYear
      isDefault
    }
  }
`;

export const GET_SEARCH_SUGGESTIONS = gql`
  query GetSearchSuggestions($query: String!, $limit: Int) {
    searchSuggestions(query: $query, limit: $limit) {
      id
      title
      type
      imageUrl
    }
  }
`;

export const GET_ACTIVE_PROMOTIONS = gql`
  query GetActivePromotions($merchantId: ID!) {
    activePromotions(merchantId: $merchantId) {
      id
      name
      description
      type
      status
      discountType
      discountValue
      startDate
      endDate
      applicableProducts
      applicableCollections
      applicableCustomers
      priority
      usageLimit
      usedCount
      createdAt
      updatedAt
    }
  }
`;

export const GET_CAMPAIGNS = gql`
  query GetCampaigns($filter: CampaignFilter) {
    campaigns(filter: $filter) {
      id
      merchantId
      name
      description
      type
      status
      startDate
      endDate
      budget
      goalType
      goalValue
      createdAt
      updatedAt
      promotions {
        id
        name
        description
        type
        status
        discountType
        discountValue
        startDate
        endDate
      }
    }
  }
`;

//---------------------------
// MUTATIONS
//---------------------------

export const LOGIN = gql`
  mutation Login($input: LoginInput!) {
    login(input: $input) {
      user {
        id
        email
        username
        firstName
        lastName
        isActive
        createdAt
      }
      accessToken
      refreshToken
      expiresIn
    }
  }
`;

export const REGISTER = gql`
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
      user {
        id
        email
        username
        firstName
        lastName
        isActive
        createdAt
      }
      accessToken
      refreshToken
      expiresIn
    }
  }
`;

export const LOGOUT = gql`
  mutation Logout {
    logout
  }
`;

export const ADD_TO_WISHLIST = gql`
  mutation AddToWishlist($productId: ID!) {
    addToWishlist(productId: $productId) {
      id
      productId
      addedAt
    }
  }
`;

export const REMOVE_FROM_WISHLIST = gql`
  mutation RemoveFromWishlist($productId: ID!) {
    removeFromWishlist(productId: $productId)
  }
`;

export const ADD_ADDRESS = gql`
  mutation AddAddress($input: AddressInput!) {
    addAddress(input: $input) {
      id
      firstName
      lastName
      company
      address1
      address2
      city
      province
      country
      zip
      phone
      isDefault
    }
  }
`;

export const UPDATE_ADDRESS = gql`
  mutation UpdateAddress($id: ID!, $input: AddressInput!) {
    updateAddress(id: $id, input: $input) {
      id
      firstName
      lastName
      company
      address1
      address2
      city
      province
      country
      zip
      phone
      isDefault
    }
  }
`;

export const REMOVE_ADDRESS = gql`
  mutation RemoveAddress($id: ID!) {
    removeAddress(id: $id)
  }
`;

export const SET_DEFAULT_ADDRESS = gql`
  mutation SetDefaultAddress($id: ID!) {
    setDefaultAddress(id: $id) {
      id
      isDefault
    }
  }
`;

export const ADD_PAYMENT_METHOD = gql`
  mutation AddPaymentMethod($input: PaymentMethodInput!) {
    addPaymentMethod(input: $input) {
      id
      type
      name
      last4
      expiryMonth
      expiryYear
      isDefault
    }
  }
`;

export const REMOVE_PAYMENT_METHOD = gql`
  mutation RemovePaymentMethod($id: ID!) {
    removePaymentMethod(id: $id)
  }
`;

export const SET_DEFAULT_PAYMENT_METHOD = gql`
  mutation SetDefaultPaymentMethod($id: ID!) {
    setDefaultPaymentMethod(id: $id) {
      id
      isDefault
    }
  }
`;
