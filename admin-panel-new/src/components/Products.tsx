export default function Products() {
  return (
    <div>
      <h2 className="text-3xl font-bold text-gray-900 mb-8">Products</h2>
      
      <div className="bg-white rounded-lg shadow">
        <div className="p-6 border-b flex justify-between items-center">
          <h3 className="text-lg font-semibold text-gray-900">Product List</h3>
          <button className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
            Add Product
          </button>
        </div>
        <div className="p-6">
          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2">Product</th>
                  <th className="text-left py-2">Price</th>
                  <th className="text-left py-2">Stock</th>
                  <th className="text-left py-2">Status</th>
                  <th className="text-left py-2">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr className="border-b">
                  <td className="py-2">Premium Wireless Headphones</td>
                  <td className="py-2">$299.99</td>
                  <td className="py-2">45</td>
                  <td className="py-2">
                    <span className="bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs">
                      Active
                    </span>
                  </td>
                  <td className="py-2">
                    <button className="text-blue-600 hover:text-blue-800 mr-2">Edit</button>
                    <button className="text-red-600 hover:text-red-800">Delete</button>
                  </td>
                </tr>
                <tr className="border-b">
                  <td className="py-2">Smart Fitness Watch</td>
                  <td className="py-2">$199.99</td>
                  <td className="py-2">23</td>
                  <td className="py-2">
                    <span className="bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs">
                      Active
                    </span>
                  </td>
                  <td className="py-2">
                    <button className="text-blue-600 hover:text-blue-800 mr-2">Edit</button>
                    <button className="text-red-600 hover:text-red-800">Delete</button>
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
