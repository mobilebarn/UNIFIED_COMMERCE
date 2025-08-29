import { useState } from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import Dashboard from './components/Dashboard'
import Products from './components/Products'
import Orders from './components/Orders'
import './App.css'

function App() {
  const [activeTab, setActiveTab] = useState('dashboard')

  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        {/* Header */}
        <header className="bg-white shadow-sm border-b">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex items-center justify-between h-16">
              <h1 className="text-2xl font-bold text-gray-900">
                Unified Commerce Admin
              </h1>
              <div className="flex items-center space-x-4">
                <span className="text-sm text-gray-500">Administrator</span>
                <button className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700">
                  Logout
                </button>
              </div>
            </div>
          </div>
        </header>

        <div className="flex">
          {/* Sidebar */}
          <nav className="w-64 bg-white shadow-sm h-screen">
            <div className="p-4">
              <ul className="space-y-2">
                <li>
                  <Link
                    to="/"
                    className={`block px-4 py-2 rounded-lg transition-colors ${
                      activeTab === 'dashboard'
                        ? 'bg-blue-600 text-white'
                        : 'text-gray-700 hover:bg-gray-100'
                    }`}
                    onClick={() => setActiveTab('dashboard')}
                  >
                    Dashboard
                  </Link>
                </li>
                <li>
                  <Link
                    to="/products"
                    className={`block px-4 py-2 rounded-lg transition-colors ${
                      activeTab === 'products'
                        ? 'bg-blue-600 text-white'
                        : 'text-gray-700 hover:bg-gray-100'
                    }`}
                    onClick={() => setActiveTab('products')}
                  >
                    Products
                  </Link>
                </li>
                <li>
                  <Link
                    to="/orders"
                    className={`block px-4 py-2 rounded-lg transition-colors ${
                      activeTab === 'orders'
                        ? 'bg-blue-600 text-white'
                        : 'text-gray-700 hover:bg-gray-100'
                    }`}
                    onClick={() => setActiveTab('orders')}
                  >
                    Orders
                  </Link>
                </li>
              </ul>
            </div>
          </nav>

          {/* Main Content */}
          <main className="flex-1 p-8">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/products" element={<Products />} />
              <Route path="/orders" element={<Orders />} />
            </Routes>
          </main>
        </div>
      </div>
    </Router>
  )
}

export default App
