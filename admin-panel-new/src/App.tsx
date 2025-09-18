import { useState, useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Link, Navigate } from 'react-router-dom'
import { ApolloProvider } from '@apollo/client'
import { useDashboardStore } from './store/dashboardStore'
import apolloClient from './lib/apollo'
import Dashboard from './components/Dashboard'
import Products from './components/Products'
import Orders from './components/Orders'
import Customers from './components/Customers'
import Analytics from './components/Analytics'
import Login from './components/Login'
import './App.css'

// Icons
const DashboardIcon = () => (
  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
    <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
  </svg>
)

const ProductIcon = () => (
  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
    <path fillRule="evenodd" d="M10 2a4 4 0 00-4 4v1H5a1 1 0 00-.994.89l-1 9A1 1 0 004 18h12a1 1 0 00.994-1.11l-1-9A1 1 0 0015 7h-1V6a4 4 0 00-4-4zm2 5V6a2 2 0 10-4 0v1h4zm-6 3a1 1 0 112 0 1 1 0 01-2 0zm7-1a1 1 0 100 2 1 1 0 000-2z" clipRule="evenodd"></path>
  </svg>
)

const OrderIcon = () => (
  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
    <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z"></path>
    <path fillRule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z" clipRule="evenodd"></path>
  </svg>
)

const CustomerIcon = () => (
  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
    <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z"></path>
  </svg>
)

// Protected Route Component
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated, user } = useDashboardStore()
  console.log('ProtectedRoute: isAuthenticated =', isAuthenticated, 'user =', user)
  
  if (isAuthenticated) {
    console.log('User is authenticated, showing protected content')
    return <>{children}</>
  } else {
    console.log('User not authenticated, redirecting to login')
    return <Navigate to="/login" replace />
  }
}

// Main App Layout Component
const AppLayout = ({ children, activeTab, setActiveTab }: { 
  children: React.ReactNode
  activeTab: string
  setActiveTab: (tab: string) => void 
}) => {
  const { user, logout } = useDashboardStore()

  const handleLogout = () => {
    logout()
  }

  return (
    <div className="min-h-screen bg-[#f9fafb]">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-10">
        <div className="max-w-full px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <span className="text-2xl font-bold bg-gradient-to-r from-emerald-600 to-teal-500 bg-clip-text text-transparent">
                  Unified Commerce OS
                </span>
              </div>
            </div>
            <div className="flex items-center gap-6">
              <div className="relative">
                <button className="text-gray-500 hover:text-gray-700">
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path>
                  </svg>
                </button>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-sm font-medium text-gray-600">
                  {user?.firstName?.[0] || user?.username?.[0] || 'A'}
                </div>
                <span className="text-sm font-medium text-gray-700">
                  {user?.firstName || user?.username || 'Admin'}
                </span>
                <button
                  onClick={handleLogout}
                  className="ml-2 text-gray-500 hover:text-gray-700"
                  title="Logout"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      <div className="flex">
        {/* Sidebar */}
        <nav className="w-64 bg-white border-r border-gray-200 h-[calc(100vh-4rem)] sticky top-16">
          <div className="p-4">
            <div className="mb-6">
              <div className="px-4 py-2">
                <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">Core</span>
              </div>
            </div>
            <ul className="space-y-1">
              <li>
                <Link
                  to="/"
                  className={`flex items-center gap-3 px-4 py-2.5 rounded-md transition-colors ${
                    activeTab === 'dashboard'
                      ? 'bg-emerald-50 text-emerald-700 font-medium'
                      : 'text-gray-700 hover:bg-gray-50'
                  }`}
                  onClick={() => setActiveTab('dashboard')}
                >
                  <DashboardIcon />
                  Dashboard
                </Link>
              </li>
              <li>
                <Link
                  to="/products"
                  className={`flex items-center gap-3 px-4 py-2.5 rounded-md transition-colors ${
                    activeTab === 'products'
                      ? 'bg-emerald-50 text-emerald-700 font-medium'
                      : 'text-gray-700 hover:bg-gray-50'
                  }`}
                  onClick={() => setActiveTab('products')}
                >
                  <ProductIcon />
                  Products
                </Link>
              </li>
              <li>
                <Link
                  to="/orders"
                  className={`flex items-center gap-3 px-4 py-2.5 rounded-md transition-colors ${
                    activeTab === 'orders'
                      ? 'bg-emerald-50 text-emerald-700 font-medium'
                      : 'text-gray-700 hover:bg-gray-50'
                  }`}
                  onClick={() => setActiveTab('orders')}
                >
                  <OrderIcon />
                  Orders
                </Link>
              </li>
              <li>
                <Link
                  to="/customers"
                  className={`flex items-center gap-3 px-4 py-2.5 rounded-md transition-colors ${
                    activeTab === 'customers'
                      ? 'bg-emerald-50 text-emerald-700 font-medium'
                      : 'text-gray-700 hover:bg-gray-50'
                  }`}
                  onClick={() => setActiveTab('customers')}
                >
                  <CustomerIcon />
                  Customers
                </Link>
              </li>
              <li>
                <Link
                  to="/analytics"
                  className={`flex items-center gap-3 px-4 py-2.5 rounded-md transition-colors ${
                    activeTab === 'analytics'
                      ? 'bg-emerald-50 text-emerald-700 font-medium'
                      : 'text-gray-700 hover:bg-gray-50'
                  }`}
                  onClick={() => setActiveTab('analytics')}
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 3a1 1 0 00-1 1v16a1 1 0 002 0V4a1 1 0 00-1-1zM5 9a1 1 0 00-1 1v10a1 1 0 002 0V10a1 1 0 00-1-1zm12-6a1 1 0 00-1 1v20a1 1 0 002 0V4a1 1 0 00-1-1zm6 8a1 1 0 00-1 1v12a1 1 0 002 0V12a1 1 0 00-1-1z"/></svg>
                  Analytics
                </Link>
              </li>

              <div className="mt-6 mb-2">
                <div className="px-4 py-2">
                  <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">Sales Channels</span>
                </div>
              </div>
              <li>
                <Link
                  to="/online-store"
                  className="flex items-center gap-3 px-4 py-2.5 rounded-md text-gray-700 hover:bg-gray-50"
                  onClick={() => setActiveTab('online-store')}
                >
                  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                    <path d="M11 17a1 1 0 001.447.894l4-2A1 1 0 0017 15V9.236a1 1 0 00-1.447-.894l-4 2a1 1 0 00-.553.894V17zM15.211 6.276a1 1 0 000-1.788l-4.764-2.382a1 1 0 00-.894 0L4.789 4.488a1 1 0 000 1.788l4.764 2.382a1 1 0 00.894 0l4.764-2.382z"></path>
                    <path d="M4.447 8.342A1 1 0 003 9.236V15a1 1 0 00.553.894l4 2A1 1 0 009 17v-5.764a1 1 0 00-.553-.894l-4-2z"></path>
                  </svg>
                  Online Store
                </Link>
              </li>
              <li>
                <Link
                  to="/point-of-sale"
                  className="flex items-center gap-3 px-4 py-2.5 rounded-md text-gray-700 hover:bg-gray-50"
                  onClick={() => setActiveTab('point-of-sale')}
                >
                  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                    <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4z"></path>
                    <path fillRule="evenodd" d="M18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z" clipRule="evenodd"></path>
                  </svg>
                  Point of Sale
                </Link>
              </li>
            </ul>
          </div>
        </nav>

        {/* Main Content */}
        <main className="flex-1 bg-gray-50 overflow-auto">
          {children}
        </main>
      </div>
    </div>
  )
}

function App() {
  const [activeTab, setActiveTab] = useState('dashboard')
  const { isAuthenticated, setUser } = useDashboardStore()

  // Check for existing auth token on app load
  useEffect(() => {
    const token = localStorage.getItem('accessToken')
    if (token && !isAuthenticated) {
      // For now, we'll just set a basic user object
      // In a real app, you'd validate the token with the server
      setUser({
        id: 'current-user',
        email: 'admin@example.com',
        username: 'admin',
        firstName: 'Admin',
        lastName: 'User',
        role: 'admin'
      })
    }
  }, [isAuthenticated, setUser])

  return (
    <ApolloProvider client={apolloClient}>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/*" element={
            <ProtectedRoute>
              <AppLayout activeTab={activeTab} setActiveTab={setActiveTab}>
                <Routes>
                  <Route path="/" element={<Dashboard />} />
                  <Route path="/products" element={<Products />} />
                  <Route path="/orders" element={<Orders />} />
                  <Route path="/customers" element={<Customers />} />
                  <Route path="/analytics" element={<Analytics />} />
                </Routes>
              </AppLayout>
            </ProtectedRoute>
          } />
        </Routes>
      </Router>
    </ApolloProvider>
  )
}

export default App