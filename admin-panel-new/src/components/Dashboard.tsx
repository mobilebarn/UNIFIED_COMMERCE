export default function Dashboard() {
  return (
    <div>
      <h2 className="text-3xl font-bold text-gray-900 mb-8">Dashboard</h2>
      
      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold text-gray-700 mb-2">Total Orders</h3>
          <p className="text-3xl font-bold text-blue-600">1,234</p>
          <p className="text-sm text-green-600">+12% from last month</p>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold text-gray-700 mb-2">Revenue</h3>
          <p className="text-3xl font-bold text-green-600">$45,678</p>
          <p className="text-sm text-green-600">+8% from last month</p>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold text-gray-700 mb-2">Products</h3>
          <p className="text-3xl font-bold text-purple-600">567</p>
          <p className="text-sm text-blue-600">23 new this week</p>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold text-gray-700 mb-2">Customers</h3>
          <p className="text-3xl font-bold text-orange-600">8,901</p>
          <p className="text-sm text-green-600">+15% from last month</p>
        </div>
      </div>

      {/* Recent Orders */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b">
          <h3 className="text-lg font-semibold text-gray-900">Recent Orders</h3>
        </div>
        <div className="p-6">
          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2">Order ID</th>
                  <th className="text-left py-2">Customer</th>
                  <th className="text-left py-2">Amount</th>
                  <th className="text-left py-2">Status</th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-b">
                  <td className="py-2">#12345</td>
                  <td className="py-2">John Doe</td>
                  <td className="py-2">$299.99</td>
                  <td className="py-2">
                    <span className="bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs">
                      Completed
                    </span>
                  </td>
                </tr>
                <tr className="border-b">
                  <td className="py-2">#12346</td>
                  <td className="py-2">Jane Smith</td>
                  <td className="py-2">$199.99</td>
                  <td className="py-2">
                    <span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full text-xs">
                      Processing
                    </span>
                  </td>
                </tr>
                <tr>
                  <td className="py-2">#12347</td>
                  <td className="py-2">Bob Johnson</td>
                  <td className="py-2">$79.99</td>
                  <td className="py-2">
                    <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs">
                      Shipped
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  )
}
