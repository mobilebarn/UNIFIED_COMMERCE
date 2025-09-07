import { apolloClient } from '../lib/apollo';
import { GET_DASHBOARD_STATS } from '../lib/graphql';

// Interfaces for dashboard data structures
export interface RevenuePoint { month: string; amount: number }
export interface OrderSummary { id: string; customer: string; total: number; status: string; item: string }
export interface DashboardMetrics {
  totalRevenue: number;
  totalOrders: number;
  newCustomers: number;
  conversionRate: number;
  averageOrderValue: number;
  returningCustomerRate: number;
  topProducts: { name: string; sales: number }[];
}

// Caching mechanism to prevent redundant API calls
let dashboardDataPromise: Promise<any> | null = null;
let lastRange: string | null = null;

// Fetches data from GraphQL endpoint and caches the promise
async function fetchAndCacheDashboardData(range: '7d' | '30d' | '90d' | 'ytd') {
  if (!dashboardDataPromise || lastRange !== range) {
    lastRange = range;
    dashboardDataPromise = apolloClient.query({
      query: GET_DASHBOARD_STATS,
      variables: { period: range },
      fetchPolicy: 'network-only', // Ensures we always get fresh data
    });
  }
  return dashboardDataPromise;
}

// Fetches and transforms dashboard metrics
export async function fetchDashboardMetrics(range: '7d' | '30d' | '90d' | 'ytd'): Promise<DashboardMetrics> {
  const { data } = await fetchAndCacheDashboardData(range);
  const stats = data.dashboardStats;
  
  return {
    totalRevenue: stats.totalRevenue || 0,
    totalOrders: stats.totalOrders || 0,
    newCustomers: stats.totalCustomers || 0,
    averageOrderValue: stats.averageOrderValue || 0,
    // NOTE: The following metrics are not in the GraphQL schema. Defaulting to 0.
    conversionRate: stats.conversionRate || 0,
    returningCustomerRate: stats.returningCustomerRate || 0,
    topProducts: stats.topProducts?.map((p: any) => ({ name: p.title, sales: p.orderCount })) || [],
  };
}

// Fetches and transforms revenue series data
export async function fetchRevenueSeries(range: '7d' | '30d' | '90d' | 'ytd'): Promise<RevenuePoint[]> {
  await fetchAndCacheDashboardData(range);
  // NOTE: The GraphQL schema does not currently return time-series data for the revenue chart.
  // Returning an empty array to prevent UI errors. This can be implemented in the backend.
  return [];
}

// Fetches and transforms recent order data
export async function fetchRecentOrders(limit = 5): Promise<OrderSummary[]> {
  const { data } = await fetchAndCacheDashboardData('30d'); // Using a default range for recent orders
  const stats = data.dashboardStats;

  return stats.recentOrders?.slice(0, limit).map((o: any) => ({
    id: o.orderNumber,
    customer: o.customer ? `${o.customer.firstName} ${o.customer.lastName}` : 'N/A',
    total: o.total,
    status: o.status,
    // NOTE: The schema provides a list of items, not a single summary item.
    item: `${o.items?.length || 0} items`,
  })) || [];
}

// The following functions are not used by the dashboard but are kept for other potential uses.
// They still use mock data and should be updated to use GraphQL mutations.

const delay = (ms:number)=> new Promise(r=>setTimeout(r, ms));

export async function createProductDraft(data:{ name:string; price:number }) {
  await delay(200);
  return { id: 'temp-' + Date.now(), ...data };
}

export async function createOrderDraft(data:{ customer:string; total:number }) {
  await delay(200);
  return { id: '#NEW', ...data };
}

export async function createCustomerDraft(data:{ name:string; email:string }) {
  await delay(150);
  return { id: 'cust-' + Date.now(), ...data };
}
