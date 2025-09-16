import { apolloClient } from '../lib/apollo';
import { 
  GET_CUSTOMER_BEHAVIORS, 
  GET_PRODUCT_RECOMMENDATIONS, 
  GET_CUSTOMER_SEGMENTS, 
  GET_BUSINESS_METRICS 
} from '../lib/graphql';

// Interfaces for analytics data structures
export interface CustomerBehavior {
  id: string;
  customerId: string;
  action: string;
  entityId: string;
  entityType: string;
  timestamp: string;
  userAgent?: string;
  referrer?: string;
}

export interface ProductRecommendation {
  id: string;
  customerId: string;
  productId: string;
  product?: {
    id: string;
    title: string;
    price: number;
  };
  score: number;
  recommendationType: string;
  createdAt: string;
  expiresAt: string;
}

export interface CustomerSegment {
  id: string;
  name: string;
  description?: string;
  customerIds: string[];
  createdAt: string;
  updatedAt: string;
}

export interface BusinessMetric {
  id: string;
  name: string;
  value: number;
  dimension?: string;
  timestamp: string;
  tags?: Record<string, string>;
}

// Fetch customer behaviors
export async function fetchCustomerBehaviors(customerId: string, limit = 100): Promise<CustomerBehavior[]> {
  try {
    const { data } = await apolloClient.query({
      query: GET_CUSTOMER_BEHAVIORS,
      variables: { customerId, limit },
      fetchPolicy: 'network-only',
    });
    return data.customerBehaviors || [];
  } catch (error: any) {
    console.error('Error fetching customer behaviors:', error);
    throw new Error(error.message || 'Failed to fetch customer behaviors');
  }
}

// Fetch product recommendations
export async function fetchProductRecommendations(customerId: string, limit = 10): Promise<ProductRecommendation[]> {
  try {
    const { data } = await apolloClient.query({
      query: GET_PRODUCT_RECOMMENDATIONS,
      variables: { customerId, limit },
      fetchPolicy: 'network-only',
    });
    return data.productRecommendations || [];
  } catch (error: any) {
    console.error('Error fetching product recommendations:', error);
    throw new Error(error.message || 'Failed to fetch product recommendations');
  }
}

// Fetch customer segments
export async function fetchCustomerSegments(): Promise<CustomerSegment[]> {
  try {
    const { data } = await apolloClient.query({
      query: GET_CUSTOMER_SEGMENTS,
      fetchPolicy: 'network-only',
    });
    return data.customerSegments || [];
  } catch (error: any) {
    console.error('Error fetching customer segments:', error);
    throw new Error(error.message || 'Failed to fetch customer segments');
  }
}

// Fetch business metrics
export async function fetchBusinessMetrics(name: string, start?: string, end?: string): Promise<BusinessMetric[]> {
  try {
    const { data } = await apolloClient.query({
      query: GET_BUSINESS_METRICS,
      variables: { name, start, end },
      fetchPolicy: 'network-only',
    });
    return data.businessMetrics || [];
  } catch (error: any) {
    console.error('Error fetching business metrics:', error);
    throw new Error(error.message || 'Failed to fetch business metrics');
  }
}

// Generate product recommendations
export async function generateProductRecommendations(customerId: string, algorithm?: string, limit?: number): Promise<ProductRecommendation[]> {
  // This would be a mutation in a real implementation
  // For now, we'll just fetch existing recommendations
  return fetchProductRecommendations(customerId, limit);
}