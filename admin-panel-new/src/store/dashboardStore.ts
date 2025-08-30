import { create } from 'zustand'
import { fetchDashboardMetrics, fetchRevenueSeries, fetchRecentOrders, DashboardMetrics, RevenuePoint, OrderSummary } from '../services/dashboard'

interface User {
  id: string
  email: string
  username: string
  firstName: string
  lastName: string
  role: string
}

interface DashboardState {
  range: '7d'|'30d'|'90d'|'ytd'
  loading: boolean
  metrics?: DashboardMetrics
  revenue: RevenuePoint[]
  recent: OrderSummary[]
  error?: string
  user?: User
  isAuthenticated: boolean
  setRange: (r:DashboardState['range'])=>Promise<void>
  refresh: ()=>Promise<void>
  setUser: (user: User) => void
  logout: () => void
}

export const useDashboardStore = create<DashboardState>((set,get)=>({
  range: '7d',
  loading: false,
  revenue: [],
  recent: [],
  isAuthenticated: false,
  async setRange(r){
    set({ range: r })
    await get().refresh()
  },
  async refresh(){
    const r = get().range
    try {
      set({ loading: true, error: undefined })
      const [metrics, revenue, recent] = await Promise.all([
        fetchDashboardMetrics(r),
        fetchRevenueSeries(r),
        fetchRecentOrders(5)
      ])
      set({ metrics, revenue, recent, loading:false })
    } catch(e:any){
      set({ error: e.message || 'Failed to load dashboard', loading:false })
    }
  },
  setUser: (user: User) => {
    set({ user, isAuthenticated: true })
  },
  logout: () => {
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    set({ user: undefined, isAuthenticated: false })
  }
}))
