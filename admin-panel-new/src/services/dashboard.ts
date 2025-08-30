// Simple front-end dashboard service layer.
// In production, replace mock implementations with real API calls including auth headers.
import { api } from './apiClient'

export interface RevenuePoint { month: string; amount: number }
export interface OrderSummary { id: string; customer: string; total: number; status: string; item: string }
export interface DashboardMetrics {
  totalRevenue: number
  totalOrders: number
  newCustomers: number
  conversionRate: number
  averageOrderValue: number
  returningCustomerRate: number
  topProducts: { name: string; sales: number }[]
}

// Mock delay helper
const delay = (ms:number)=> new Promise(r=>setTimeout(r, ms))

// Example base URL (adjust to gateway / API aggregation once available)
const BASE = '/api/v1'

export async function fetchDashboardMetrics(range:'7d'|'30d'|'90d'|'ytd'): Promise<DashboardMetrics> {
  try {
    // Attempt real analytics endpoint (to be implemented). Example REST path.
    // Expecting server response shape to match DashboardMetrics; adapt if different.
    const { data } = await api.get(`/api/v1/analytics/dashboard`, { params: { range } })
    if(data && typeof data === 'object') return data as DashboardMetrics
    throw new Error('Unexpected dashboard metrics response')
  } catch (e){
    // Fallback mock
    await delay(150)
    return {
      totalRevenue: range==='7d'?45678: range==='30d'?181234: 512345,
      totalOrders: range==='7d'?1284: range==='30d'?5120: 16890,
      newCustomers: range==='7d'?347: range==='30d'?1290: 4980,
      conversionRate: 3.4,
      averageOrderValue: 71.2,
      returningCustomerRate: 42.5,
      topProducts: [
        { name: 'Premium Wireless Headphones', sales: 320 },
        { name: 'Smart Fitness Watch', sales: 280 },
        { name: 'Portable Bluetooth Speaker', sales: 210 }
      ]
    }
  }
}

export async function fetchRevenueSeries(range:'7d'|'30d'|'90d'|'ytd'): Promise<RevenuePoint[]> {
  try {
    const { data } = await api.get(`/api/v1/analytics/revenue`, { params: { range } })
    if(Array.isArray(data)) return data as RevenuePoint[]
    throw new Error('Unexpected revenue response')
  } catch (e){
    await delay(120)
    if(range==='7d') {
      return ['Mon','Tue','Wed','Thu','Fri','Sat','Sun'].map((d,i)=>({month: d, amount: 12000 + i*1500}))
    }
    const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
    return months.slice(0, range==='30d'?3: range==='90d'?6: 8).map((m,i)=>({month:m, amount: 28000 + i*3500}))
  }
}

export async function fetchRecentOrders(limit=5): Promise<OrderSummary[]> {
  try {
    const { data } = await api.get(`/api/v1/orders`, { params: { per_page: limit, page: 1 } })
    // Expect an array or object with data list; adapt if metadata wrapper present
  const list = Array.isArray(data) ? data : (Array.isArray((data as any)?.data) ? (data as any).data : [])
    if(list.length){
      return list.map((o:any)=> ({
        id: o.order_number || o.id,
        customer: o.customer_name || o.customer?.name || '—',
        total: o.total_amount || o.total || 0,
        status: o.status || '—',
        item: o.items?.[0]?.name || `${o.items?.length||0} items`
      }))
    }
    throw new Error('Empty orders response')
  } catch (e){
    await delay(200)
    return [
      { id:'#1234', customer:'John Doe', total:299.99, status:'Completed', item:'Premium Headphones' },
      { id:'#1235', customer:'Jane Smith', total:199.99, status:'Processing', item:'Smart Watch' },
      { id:'#1236', customer:'Bob Johnson', total:79.99, status:'Shipped', item:'Bluetooth Speaker' },
      { id:'#1237', customer:'Emily Wilson', total:129.99, status:'Completed', item:'Mechanical Keyboard' },
      { id:'#1238', customer:'Robert Miles', total:449.99, status:'Completed', item:'Gaming Tablet' }
    ].slice(0,limit)
  }
}

export async function createProductDraft(data:{ name:string; price:number }) {
  try {
    const { data:resp } = await api.post(`/api/v1/products`, { ...data, sku: 'TEMP-'+Date.now(), currency: 'USD' }, { headers: { 'X-Merchant-ID':'demo-merchant' }})
    return resp
  } catch {
    await delay(200)
    return { id: 'temp-'+Date.now(), ...data }
  }
}

export async function createOrderDraft(data:{ customer:string; total:number }) {
  try {
    const { data:resp } = await api.post(`/api/v1/orders`, { customer_name: data.customer, total: data.total, merchant_id: 'demo-merchant' })
    return resp
  } catch {
    await delay(200)
    return { id:'#NEW', ...data }
  }
}

export async function createCustomerDraft(data:{ name:string; email:string }) {
  await delay(150)
  return { id:'cust-'+Date.now(), ...data } // TODO: call customer service when available
}
