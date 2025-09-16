import { create } from 'zustand';
import { 
  fetchCustomerBehaviors, 
  fetchProductRecommendations, 
  fetchCustomerSegments, 
  fetchBusinessMetrics,
  CustomerBehavior,
  ProductRecommendation,
  CustomerSegment,
  BusinessMetric
} from '../services/analytics';

interface AnalyticsState {
  // Customer behaviors
  behaviors: CustomerBehavior[];
  behaviorsLoading: boolean;
  behaviorsError?: string;
  
  // Product recommendations
  recommendations: ProductRecommendation[];
  recommendationsLoading: boolean;
  recommendationsError?: string;
  
  // Customer segments
  segments: CustomerSegment[];
  segmentsLoading: boolean;
  segmentsError?: string;
  
  // Business metrics
  metrics: BusinessMetric[];
  metricsLoading: boolean;
  metricsError?: string;
  
  // Actions
  fetchBehaviors: (customerId: string, limit?: number) => Promise<void>;
  fetchRecommendations: (customerId: string, limit?: number) => Promise<void>;
  fetchSegments: () => Promise<void>;
  fetchMetrics: (name: string, start?: string, end?: string) => Promise<void>;
  clearErrors: () => void;
}

export const useAnalyticsStore = create<AnalyticsState>((set, get) => ({
  // Initial state
  behaviors: [],
  behaviorsLoading: false,
  behaviorsError: undefined,
  
  recommendations: [],
  recommendationsLoading: false,
  recommendationsError: undefined,
  
  segments: [],
  segmentsLoading: false,
  segmentsError: undefined,
  
  metrics: [],
  metricsLoading: false,
  metricsError: undefined,
  
  // Actions
  fetchBehaviors: async (customerId: string, limit = 100) => {
    try {
      set({ behaviorsLoading: true, behaviorsError: undefined });
      const behaviors = await fetchCustomerBehaviors(customerId, limit);
      set({ behaviors, behaviorsLoading: false });
    } catch (error: any) {
      set({ behaviorsError: error.message, behaviorsLoading: false });
    }
  },
  
  fetchRecommendations: async (customerId: string, limit = 10) => {
    try {
      set({ recommendationsLoading: true, recommendationsError: undefined });
      const recommendations = await fetchProductRecommendations(customerId, limit);
      set({ recommendations, recommendationsLoading: false });
    } catch (error: any) {
      set({ recommendationsError: error.message, recommendationsLoading: false });
    }
  },
  
  fetchSegments: async () => {
    try {
      set({ segmentsLoading: true, segmentsError: undefined });
      const segments = await fetchCustomerSegments();
      set({ segments, segmentsLoading: false });
    } catch (error: any) {
      set({ segmentsError: error.message, segmentsLoading: false });
    }
  },
  
  fetchMetrics: async (name: string, start?: string, end?: string) => {
    try {
      set({ metricsLoading: true, metricsError: undefined });
      const metrics = await fetchBusinessMetrics(name, start, end);
      set({ metrics, metricsLoading: false });
    } catch (error: any) {
      set({ metricsError: error.message, metricsLoading: false });
    }
  },
  
  clearErrors: () => {
    set({ 
      behaviorsError: undefined, 
      recommendationsError: undefined, 
      segmentsError: undefined, 
      metricsError: undefined 
    });
  }
}));