export default function Orders() {
  return (
    <div>
      <h2 className="text-3xl font-bold text-gray-900 mb-8">Orders</h2>
      
      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b">
          <h3 className="text-lg font-semibold text-gray-900">Order Management</h3>
        </div>
        <div className="p-6">
          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2">Order ID</th>
                  <th className="text-left py-2">Customer</th>
                  <th className="text-left py-2">Date</th>
                  <th className="text-left py-2">Amount</th>
                  <th className="text-left py-2">Status</th>
                  <th className="text-left py-2">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-b">
                  <td className="py-2">#12345</td>
                  <td className="py-2">John Doe</td>
                  <td className="py-2">2024-01-15</td>
                  <td className="py-2">$299.99</td>
                  <td className="py-2">
                    <span className="bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs">
                      Completed
                    </span>
                  </td>
                  <td className="py-2">
                    <button className="text-blue-600 hover:text-blue-800 mr-2">View</button>
                    <button className="text-gray-600 hover:text-gray-800">Refund</button>
                  </td>
                </tr>
                <tr className="border-b">
                  <td className="py-2">#12346</td>
                  <td className="py-2">Jane Smith</td>
                  <td className="py-2">2024-01-16</td>
                  <td className="py-2">$199.99</td>
                  <td className="py-2">
                    <span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded-full text-xs">
                      Processing
                    </span>
                  </td>
                  <td className="py-2">
                    <button className="text-blue-600 hover:text-blue-800 mr-2">View</button>
                    <button className="text-green-600 hover:text-green-800">Ship</button>
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
