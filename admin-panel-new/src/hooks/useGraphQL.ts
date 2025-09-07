import { useQuery, useMutation } from '@apollo/client';
import { useApolloClient } from '@apollo/client';
import {
  GET_PRODUCTS,
  GET_PRODUCT_BY_ID,
  CREATE_PRODUCT,
  UPDATE_PRODUCT,
  DELETE_PRODUCT,
  GET_ORDERS,
  GET_ORDER_BY_ID,
  UPDATE_ORDER_STATUS,
  GET_CUSTOMERS,
  GET_CUSTOMER_BY_ID,
  GET_INVENTORY_ITEMS,
  UPDATE_INVENTORY,
  GET_PROMOTIONS,
  CREATE_PROMOTION,
  GET_MERCHANT_PROFILE,
  GET_DASHBOARD_STATS,
  type Product,
  type Order,
  type Customer,
  type InventoryItem,
  type Promotion,
} from '../lib/graphql';

// Product Hooks
export const useProducts = (filter?: any) => {
  return useQuery(GET_PRODUCTS, {
    variables: { filter },
    errorPolicy: 'all',
  });
};

export const useProduct = (id: string) => {
  return useQuery(GET_PRODUCT_BY_ID, {
    variables: { id },
    skip: !id,
    errorPolicy: 'all',
  });
};

export const useCreateProduct = () => {
  const client = useApolloClient();
  
  return useMutation(CREATE_PRODUCT, {
    onCompleted: () => {
      client.refetchQueries({ include: [GET_PRODUCTS] });
    },
    onError: (error) => {
      console.error('Error creating product:', error);
    },
  });
};

export const useUpdateProduct = () => {
  const client = useApolloClient();
  
  return useMutation(UPDATE_PRODUCT, {
    onCompleted: () => {
      client.refetchQueries({ include: [GET_PRODUCTS, GET_PRODUCT_BY_ID] });
    },
    onError: (error) => {
      console.error('Error updating product:', error);
    },
  });
};

export const useDeleteProduct = () => {
  const client = useApolloClient();
  
  return useMutation(DELETE_PRODUCT, {
    onCompleted: () => {
      client.refetchQueries({ include: [GET_PRODUCTS] });
    },
    onError: (error) => {
      console.error('Error deleting product:', error);
    },
  });
};

// Order Hooks
export const useOrders = (filter?: any) => {
  return useQuery(GET_ORDERS, {
    variables: { filter },
    errorPolicy: 'all',
  });
};

export const useOrder = (id: string) => {
  return useQuery(GET_ORDER_BY_ID, {
    variables: { id },
    skip: !id,
    errorPolicy: 'all',
  });
};

export const useUpdateOrderStatus = () => {
  const client = useApolloClient();
  
  return useMutation(UPDATE_ORDER_STATUS, {
    onCompleted: () => {
      client.refetchQueries({ include: [GET_ORDERS, GET_ORDER_BY_ID] });
    },
    onError: (error) => {
      console.error('Error updating order status:', error);
    },
  });
};

// Customer Hooks
export const useCustomers = (filter?: any) => {
  return useQuery(GET_CUSTOMERS, {
    variables: { filter },
    errorPolicy: 'all',
  });
};

export const useCustomer = (id: string) => {
  return useQuery(GET_CUSTOMER_BY_ID, {
    variables: { id },
    skip: !id,
    errorPolicy: 'all',
  });
};

// Inventory Hooks
export const useInventory = (filter?: any) => {
  return useQuery(GET_INVENTORY_ITEMS, {
    variables: { filter },
    errorPolicy: 'all',
  });
};

export const useUpdateInventory = () => {
  const client = useApolloClient();
  
  return useMutation(UPDATE_INVENTORY, {
    onCompleted: () => {
      client.refetchQueries({ include: [GET_INVENTORY_ITEMS] });
    },
    onError: (error) => {
      console.error('Error updating inventory:', error);
    },
  });
};

// Promotion Hooks
export const usePromotions = (filter?: any) => {
  return useQuery(GET_PROMOTIONS, {
    variables: { filter },
    errorPolicy: 'all',
  });
};

export const useCreatePromotion = () => {
  const client = useApolloClient();
  
  return useMutation(CREATE_PROMOTION, {
    onCompleted: () => {
      client.refetchQueries({ include: [GET_PROMOTIONS] });
    },
    onError: (error) => {
      console.error('Error creating promotion:', error);
    },
  });
};

// Merchant Account Hooks
export const useMerchantProfile = () => {
  return useQuery(GET_MERCHANT_PROFILE, {
    errorPolicy: 'all',
  });
};

// Dashboard Hooks
export const useDashboardStats = (period: string = '30d') => {
  return useQuery(GET_DASHBOARD_STATS, {
    variables: { period },
    errorPolicy: 'all',
  });
};

// Utility hooks for common patterns
export const useRefreshQueries = () => {
  const client = useApolloClient();
  
  return {
    refreshProducts: () => client.refetchQueries({ include: [GET_PRODUCTS] }),
    refreshOrders: () => client.refetchQueries({ include: [GET_ORDERS] }),
    refreshCustomers: () => client.refetchQueries({ include: [GET_CUSTOMERS] }),
    refreshInventory: () => client.refetchQueries({ include: [GET_INVENTORY_ITEMS] }),
    refreshPromotions: () => client.refetchQueries({ include: [GET_PROMOTIONS] }),
    refreshDashboard: () => client.refetchQueries({ include: [GET_DASHBOARD_STATS] }),
    refreshAll: () => client.refetchQueries({ include: 'all' }),
  };
};
