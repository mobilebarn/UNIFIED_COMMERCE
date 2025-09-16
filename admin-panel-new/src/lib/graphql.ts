import { gql } from '@apollo/client';

// Product Queries and Mutations
export const GET_PRODUCTS = gql`
  query GetProducts($filter: ProductFilter) {
    products(filter: $filter) {
      id
      title
      description
      price
      status
      createdAt
      updatedAt
      variants {
        id
        sku
        price
        compareAtPrice
        inventory {
          quantity
          location
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
      updatedAt
      variants {
        id
        sku
        price
        compareAtPrice
        inventory {
          quantity
          location
        }
      }
      categories {
        id
        name
      }
    }
  }
`;

export const CREATE_PRODUCT = gql`
  mutation CreateProduct($input: CreateProductInput!) {
    createProduct(input: $input) {
      id
      title
      description
      price
      status
    }
  }
`;

export const UPDATE_PRODUCT = gql`
  mutation UpdateProduct($id: ID!, $input: UpdateProductInput!) {
    updateProduct(id: $id, input: $input) {
      id
      title
      description
      price
      status
    }
  }
`;

export const DELETE_PRODUCT = gql`
  mutation DeleteProduct($id: ID!) {
    deleteProduct(id: $id)
  }
`;

// Order Queries and Mutations
export const GET_ORDERS = gql`
  query GetOrders($filter: OrderFilter) {
    orders(filter: $filter) {
      id
      orderNumber
      status
      total
      currency
      createdAt
      updatedAt
      customer {
        id
        firstName
        lastName
        email
      }
      items {
        id
        quantity
        price
        product {
          id
          title
        }
        variant {
          id
          sku
        }
      }
      payments {
        id
        status
        amount
        method
        transactionId
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

export const GET_ORDER_BY_ID = gql`
  query GetOrderById($id: ID!) {
    order(id: $id) {
      id
      orderNumber
      status
      total
      currency
      createdAt
      updatedAt
      customer {
        id
        firstName
        lastName
        email
        phone
      }
      shippingAddress {
        firstName
        lastName
        address1
        address2
        city
        province
        zip
        country
      }
      billingAddress {
        firstName
        lastName
        address1
        address2
        city
        province
        zip
        country
      }
      items {
        id
        quantity
        price
        product {
          id
          title
        }
        variant {
          id
          sku
        }
      }
      payments {
        id
        status
        amount
        method
        transactionId
        createdAt
      }
      fulfillments {
        id
        status
        trackingNumber
        carrier
        createdAt
      }
    }
  }
`;

export const UPDATE_ORDER_STATUS = gql`
  mutation UpdateOrderStatus($id: ID!, $status: OrderStatus!) {
    updateOrderStatus(id: $id, status: $status) {
      id
      status
    }
  }
`;

// Customer Queries
export const GET_CUSTOMERS = gql`
  query GetCustomers($filter: UserFilter) {
    users(filter: $filter) {
      id
      firstName
      lastName
      email
      phone
      status
      createdAt
      totalOrders
      totalSpent
    }
  }
`;

export const GET_CUSTOMER_BY_ID = gql`
  query GetCustomerById($id: ID!) {
    user(id: $id) {
      id
      firstName
      lastName
      email
      phone
      status
      createdAt
      orders {
        id
        orderNumber
        status
        total
        createdAt
      }
    }
  }
`;

// Inventory Queries
export const GET_INVENTORY_ITEMS = gql`
  query GetInventoryItems($filter: InventoryFilter) {
    inventory(filter: $filter) {
      id
      productId
      variantId
      location
      quantity
      reserved
      available
      product {
        id
        title
      }
      variant {
        id
        sku
      }
    }
  }
`;

export const UPDATE_INVENTORY = gql`
  mutation UpdateInventory($id: ID!, $input: UpdateInventoryInput!) {
    updateInventory(id: $id, input: $input) {
      id
      quantity
      reserved
      available
    }
  }
`;

// Promotions Queries
export const GET_PROMOTIONS = gql`
  query GetPromotions($filter: PromotionFilter) {
    promotions(filter: $filter) {
      id
      name
      description
      type
      value
      status
      startDate
      endDate
      usageLimit
      usedCount
      createdAt
    }
  }
`;

export const CREATE_PROMOTION = gql`
  mutation CreatePromotion($input: CreatePromotionInput!) {
    createPromotion(input: $input) {
      id
      name
      description
      type
      value
      status
    }
  }
`;

// Add Customer/User Mutations
export const REGISTER_USER = gql`
  mutation RegisterUser($input: RegisterInput!) {
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

// Add Order Mutations
export const CREATE_ORDER = gql`
  mutation CreateOrder($input: CreateOrderInput!) {
    createOrder(input: $input) {
      id
      orderNumber
      status
      total
      currency
      createdAt
      customer {
        firstName
        lastName
        email
      }
      items {
        id
        quantity
        price
        product {
          id
          title
        }
      }
    }
  }
`;

// Merchant Account Queries
export const GET_MERCHANT_PROFILE = gql`
  query GetMerchantProfile {
    merchant {
      id
      businessName
      email
      phone
      status
      subscription {
        id
        plan
        status
        currentPeriodEnd
      }
      stores {
        id
        name
        address
        city
        province
        country
        status
      }
    }
  }
`;

// Dashboard Analytics
export const GET_DASHBOARD_STATS = gql`
  query GetDashboardStats($period: String!) {
    dashboardStats(period: $period) {
      totalRevenue
      totalOrders
      totalCustomers
      averageOrderValue
      revenueGrowth
      orderGrowth
      customerGrowth
      topProducts {
        id
        title
        revenue
        orderCount
      }
      recentOrders {
        id
        orderNumber
        customer {
          firstName
          lastName
        }
        total
        status
        createdAt
      }
    }
  }
`;

// Analytics Queries
export const GET_CUSTOMER_BEHAVIORS = gql`
  query GetCustomerBehaviors($customerId: ID!, $limit: Int) {
    customerBehaviors(customerId: $customerId, limit: $limit) {
      id
      action
      entityId
      entityType
      timestamp
      userAgent
      referrer
    }
  }
`;

export const GET_PRODUCT_RECOMMENDATIONS = gql`
  query GetProductRecommendations($customerId: ID!, $limit: Int) {
    productRecommendations(customerId: $customerId, limit: $limit) {
      id
      productId
      score
      recommendationType
      createdAt
      expiresAt
      product {
        id
        title
        price
      }
    }
  }
`;

export const GET_CUSTOMER_SEGMENTS = gql`
  query GetCustomerSegments {
    customerSegments {
      id
      name
      description
      customerIds
      createdAt
      updatedAt
    }
  }
`;

export const GET_BUSINESS_METRICS = gql`
  query GetBusinessMetrics($name: String!, $start: String, $end: String) {
    businessMetrics(name: $name, start: $start, end: $end) {
      id
      name
      value
      dimension
      timestamp
      tags
    }
  }
`;

// TypeScript interfaces for type safety
export interface Product {
  id: string;
  title: string;
  description?: string;
  price: number;
  status: 'active' | 'draft' | 'archived';
  createdAt: string;
  updatedAt: string;
  variants?: ProductVariant[];
  categories?: Category[];
}

export interface ProductVariant {
  id: string;
  sku: string;
  price: number;
  compareAtPrice?: number;
  inventory?: InventoryItem[];
}

export interface Category {
  id: string;
  name: string;
}

export interface Order {
  id: string;
  orderNumber: string;
  status: 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled';
  total: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
  customer?: Customer;
  items?: OrderItem[];
  payments?: Payment[];
  fulfillments?: Fulfillment[];
}

export interface Customer {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone?: string;
  status: 'active' | 'inactive';
  createdAt: string;
  totalOrders?: number;
  totalSpent?: number;
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
  status: 'pending' | 'completed' | 'failed';
  amount: number;
  method: string;
  transactionId?: string;
  createdAt: string;
}

export interface Fulfillment {
  id: string;
  status: 'pending' | 'shipped' | 'delivered';
  trackingNumber?: string;
  carrier?: string;
  createdAt: string;
}

export interface InventoryItem {
  id: string;
  productId: string;
  variantId?: string;
  location: string;
  quantity: number;
  reserved: number;
  available: number;
  product?: Product;
  variant?: ProductVariant;
}

export interface Promotion {
  id: string;
  name: string;
  description?: string;
  type: 'percentage' | 'fixed_amount' | 'free_shipping';
  value: number;
  status: 'active' | 'inactive' | 'expired';
  startDate: string;
  endDate: string;
  usageLimit?: number;
  usedCount: number;
  createdAt: string;
}
