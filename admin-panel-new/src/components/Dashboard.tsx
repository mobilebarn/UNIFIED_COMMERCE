import React, { useEffect, useState } from 'react';
import { useDashboardStore } from '../store/dashboardStore';
import { exportToCsv } from '../utils/exportCsv';
import { AddProductModal } from './modals/AddProductModal';
import { CreateOrderModal } from './modals/CreateOrderModal';
import { AddCustomerModal } from './modals/AddCustomerModal';

export default function Dashboard() {
  const { metrics, revenue, recent, loading, range, setRange, refresh, error } = useDashboardStore()
  const [showAddProduct, setShowAddProduct] = useState(false)
  const [showCreateOrder, setShowCreateOrder] = useState(false)
  const [showAddCustomer, setShowAddCustomer] = useState(false)

  useEffect(()=>{ refresh() }, [])

  const activeBtn = 'px-3 py-1 text-sm bg-blue-100 text-blue-600 rounded-md font-medium'
  const inactiveBtn = 'px-3 py-1 text-sm text-gray-500 hover:bg-gray-100 rounded-md'

  function downloadReport(){
    if(!metrics) return
    exportToCsv(`dashboard-${range}.csv`, [{
      range,
      totalRevenue: metrics.totalRevenue,
      totalOrders: metrics.totalOrders,
      newCustomers: metrics.newCustomers,
      conversionRate: metrics.conversionRate,
      averageOrderValue: metrics.averageOrderValue,
      returningCustomerRate: metrics.returningCustomerRate
    }])
  }

  return (
    <div className="h-full">
      <div className="p-6">
        {/* Header */}
        <div className="mb-8 flex flex-col lg:flex-row lg:justify-between lg:items-center gap-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
            <p className="text-gray-600 mt-1">Real-time snapshot of store performance. Range: {range.toUpperCase()}.</p>
          </div>
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
            <div className="relative">
              <select value={range} onChange={e=> setRange(e.target.value as any)} className="pl-4 pr-10 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 appearance-none bg-white min-w-[140px]">
                <option value="7d">Last 7 days</option>
                <option value="30d">Last 30 days</option>
                <option value="90d">Last 90 days</option>
                <option value="ytd">Year to date</option>
              </select>
              <svg className="absolute right-3 top-3 w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd"></path>
              </svg>
            </div>
            <button onClick={downloadReport} className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm flex items-center gap-2 transition-colors whitespace-nowrap">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4"></path>
              </svg>
              Export CSV
            </button>
          </div>
        </div>
        
        {/* Stats Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 lg:gap-6 mb-8">
          {/* Total Revenue */}
          <div className="bg-white p-4 lg:p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow relative">
            {loading && <div className="absolute inset-0 bg-white/70 flex items-center justify-center text-xs text-gray-500">Loading...</div>}
            <div className="flex items-center justify-between mb-4">
              <div className="p-2 bg-green-100 rounded-lg">
                <svg className="w-5 h-5 lg:w-6 lg:h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
                </svg>
              </div>
              <span className="text-xs lg:text-sm text-green-600 bg-green-50 px-2 py-1 rounded-full font-medium">+12.5%</span>
            </div>
            <h3 className="text-xl lg:text-2xl font-bold text-gray-900 mb-1">${metrics?.totalRevenue?.toLocaleString() || '—'}</h3>
            <p className="text-sm text-gray-600">Total Revenue</p>
          </div>

          {/* Total Orders */}
          <div className="bg-white p-4 lg:p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow relative">
            <div className="flex items-center justify-between mb-4">
              <div className="p-2 bg-blue-100 rounded-lg">
                <svg className="w-5 h-5 lg:w-6 lg:h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"></path>
                </svg>
              </div>
              <span className="text-xs lg:text-sm text-blue-600 bg-blue-50 px-2 py-1 rounded-full font-medium">+8.2%</span>
            </div>
            <h3 className="text-xl lg:text-2xl font-bold text-gray-900 mb-1">{metrics?.totalOrders?.toLocaleString() || '—'}</h3>
            <p className="text-sm text-gray-600">Total Orders</p>
          </div>

          {/* New Customers */}
          <div className="bg-white p-4 lg:p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow relative">
            <div className="flex items-center justify-between mb-4">
              <div className="p-2 bg-purple-100 rounded-lg">
                <svg className="w-5 h-5 lg:w-6 lg:h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
                </svg>
              </div>
              <span className="text-xs lg:text-sm text-purple-600 bg-purple-50 px-2 py-1 rounded-full font-medium">+15.3%</span>
            </div>
            <h3 className="text-xl lg:text-2xl font-bold text-gray-900 mb-1">{metrics?.newCustomers?.toLocaleString() || '—'}</h3>
            <p className="text-sm text-gray-600">New Customers</p>
          </div>

          {/* Conversion Rate */}
          <div className="bg-white p-4 lg:p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow relative">
            <div className="flex items-center justify-between mb-4">
              <div className="p-2 bg-orange-100 rounded-lg">
                <svg className="w-5 h-5 lg:w-6 lg:h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
                </svg>
              </div>
              <span className="text-xs lg:text-sm text-orange-600 bg-orange-50 px-2 py-1 rounded-full font-medium">+2.1%</span>
            </div>
            <h3 className="text-xl lg:text-2xl font-bold text-gray-900 mb-1">{metrics? metrics.conversionRate+'%':'—'}</h3>
            <p className="text-sm text-gray-600">Conversion Rate</p>
          </div>
        </div>

        {/* Secondary Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 lg:gap-6 mb-8">
          <div className="bg-white p-4 rounded-xl border border-gray-100">
            <p className="text-xs uppercase font-medium text-gray-500 mb-1">Avg Order Value</p>
            <p className="text-lg font-semibold text-gray-900">${metrics?.averageOrderValue?.toFixed(2) || '—'}</p>
          </div>
          <div className="bg-white p-4 rounded-xl border border-gray-100">
            <p className="text-xs uppercase font-medium text-gray-500 mb-1">Returning Customer Rate</p>
            <p className="text-lg font-semibold text-gray-900">{metrics? metrics.returningCustomerRate+'%':'—'}</p>
          </div>
          <div className="bg-white p-4 rounded-xl border border-gray-100">
            <p className="text-xs uppercase font-medium text-gray-500 mb-1">Top Product</p>
            <p className="text-lg font-semibold text-gray-900 truncate">{metrics?.topProducts?.[0]?.name || '—'}</p>
          </div>
        </div>

        {/* Charts and Recent Activity */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 lg:gap-8 mb-8">
          {/* Revenue Chart */}
          <div className="bg-white p-4 lg:p-6 rounded-xl shadow-sm border border-gray-100">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-6 gap-4">
              <h3 className="text-lg font-semibold text-gray-900">Revenue Overview</h3>
              <div className="flex space-x-2">
                <button onClick={()=> setRange('7d')} className={range==='7d'?activeBtn:inactiveBtn}>7D</button>
                <button onClick={()=> setRange('30d')} className={range==='30d'?activeBtn:inactiveBtn}>30D</button>
                <button onClick={()=> setRange('90d')} className={range==='90d'?activeBtn:inactiveBtn}>90D</button>
                <button onClick={()=> setRange('ytd')} className={range==='ytd'?activeBtn:inactiveBtn}>YTD</button>
              </div>
            </div>
            <div className="h-48 lg:h-64 flex items-end space-x-1 lg:space-x-2">
              {revenue.map((item, index) => (
                <div key={index} className="flex-1 flex flex-col items-center">
                  <div
                    className="w-full bg-blue-500 rounded-t-md transition-all hover:bg-blue-600"
                    style={{ height: revenue.length? `${(item.amount / Math.max(...revenue.map(r=>r.amount))) * 100}%` : '0%' }}
                  />
                  <span className="text-xs text-gray-500 mt-2">{item.month}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Recent Orders */}
          <div className="bg-white p-4 lg:p-6 rounded-xl shadow-sm border border-gray-100">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-lg font-semibold text-gray-900">Recent Orders</h3>
              <a href="/orders" className="text-sm text-blue-600 hover:text-blue-800 font-medium">View all</a>
            </div>
            <div className="space-y-4">
              {recent.map(o=> (
                <div key={o.id} className="flex items-center justify-between p-3 lg:p-4 bg-gray-50 rounded-lg">
                  <div className="flex items-center space-x-3">
                    <div className="w-8 h-8 lg:w-10 lg:h-10 bg-blue-100 rounded-full flex items-center justify-center">
                      <span className="text-blue-600 font-semibold text-xs lg:text-sm">{o.id}</span>
                    </div>
                    <div>
                      <p className="font-medium text-gray-900 text-sm lg:text-base">{o.customer}</p>
                      <p className="text-xs lg:text-sm text-gray-500">{o.item}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-semibold text-gray-900 text-sm lg:text-base">${o.total.toFixed(2)}</p>
                    <span className="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full">{o.status}</span>
                  </div>
                </div>
              ))}
              {loading && !recent.length && <div className="text-xs text-gray-500">Loading orders...</div>}
            </div>
          </div>
        </div>

        {/* Quick Actions */}
    <div className="bg-white p-4 lg:p-6 rounded-xl shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold text-gray-900 mb-6">Quick Actions</h3>
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
      <button onClick={()=> setShowAddProduct(true)} className="flex flex-col items-center p-3 lg:p-4 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors">
              <svg className="w-6 h-6 lg:w-8 lg:h-8 text-blue-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4"></path>
              </svg>
              <span className="text-xs lg:text-sm font-medium text-gray-900 text-center">Add Product</span>
            </button>
      <button onClick={()=> setShowCreateOrder(true)} className="flex flex-col items-center p-3 lg:p-4 bg-green-50 hover:bg-green-100 rounded-lg transition-colors">
              <svg className="w-6 h-6 lg:w-8 lg:h-8 text-green-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
              </svg>
              <span className="text-xs lg:text-sm font-medium text-gray-900 text-center">Create Order</span>
            </button>
      <button onClick={()=> setShowAddCustomer(true)} className="flex flex-col items-center p-3 lg:p-4 bg-purple-50 hover:bg-purple-100 rounded-lg transition-colors">
              <svg className="w-6 h-6 lg:w-8 lg:h-8 text-purple-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
              </svg>
              <span className="text-xs lg:text-sm font-medium text-gray-900 text-center">Add Customer</span>
            </button>
      <a href="/analytics" className="flex flex-col items-center p-3 lg:p-4 bg-orange-50 hover:bg-orange-100 rounded-lg transition-colors">
              <svg className="w-6 h-6 lg:w-8 lg:h-8 text-orange-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
              </svg>
              <span className="text-xs lg:text-sm font-medium text-gray-900 text-center">View Analytics</span>
      </a>
          </div>
        </div>
    {error && <div className="mt-6 text-sm text-red-600">{error}</div>}
    <AddProductModal open={showAddProduct} onClose={()=> setShowAddProduct(false)} />
    <CreateOrderModal open={showCreateOrder} onClose={()=> setShowCreateOrder(false)} />
    <AddCustomerModal open={showAddCustomer} onClose={()=> setShowAddCustomer(false)} />
      </div>
    </div>
  );
}
