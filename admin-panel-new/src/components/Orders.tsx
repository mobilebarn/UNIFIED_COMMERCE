export default function Orders() {
  // Sample filters and orders data
  const orderStatuses = ["All", "Processing", "Shipped", "Delivered", "Cancelled", "Refunded"];
  const timeFilters = ["Last 7 days", "Last 30 days", "Last 3 months", "This year"];
  
  const orders = [
    {
      id: "#12345",
      customer: { name: "John Doe", email: "john.doe@example.com", avatar: "JD" },
      date: "2024-01-15",
      amount: 299.99,
      items: 2,
      status: "Completed",
      paymentMethod: "Credit Card",
      shippingMethod: "Express"
    },
    {
      id: "#12346",
      customer: { name: "Jane Smith", email: "jane.smith@example.com", avatar: "JS" },
      date: "2024-01-16",
      amount: 199.99,
      items: 1,
      status: "Processing",
      paymentMethod: "PayPal",
      shippingMethod: "Standard"
    },
    {
      id: "#12347",
      customer: { name: "Robert Johnson", email: "robert.j@example.com", avatar: "RJ" },
      date: "2024-01-17",
      amount: 149.50,
      items: 3,
      status: "Shipped",
      paymentMethod: "Credit Card",
      shippingMethod: "Standard"
    },
    {
      id: "#12348",
      customer: { name: "Emily Wilson", email: "emily.w@example.com", avatar: "EW" },
      date: "2024-01-18",
      amount: 399.99,
      items: 2,
      status: "Delivered",
      paymentMethod: "Apple Pay",
      shippingMethod: "Express"
    },
    {
      id: "#12349",
      customer: { name: "Michael Brown", email: "michael.b@example.com", avatar: "MB" },
      date: "2024-01-19",
      amount: 89.99,
      items: 1,
      status: "Refunded",
      paymentMethod: "Credit Card",
      shippingMethod: "Standard"
    }
  ];
  
  const getStatusColor = (status) => {
    switch(status) {
      case "Completed":
        return "bg-green-100 text-green-800";
      case "Processing":
        return "bg-yellow-100 text-yellow-800";
      case "Shipped":
        return "bg-blue-100 text-blue-800";
      case "Delivered":
        return "bg-emerald-100 text-emerald-800";
      case "Cancelled":
        return "bg-gray-100 text-gray-800";
      case "Refunded":
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };
  
  return (
    <div className="h-full">
      <div className="p-4 lg:p-6">
        <div className="mb-8 flex flex-col lg:flex-row lg:justify-between lg:items-center gap-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Orders</h1>
            <p className="text-gray-600 mt-1">Manage and fulfill your customer orders</p>
          </div>
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3">
            <button className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 text-sm flex items-center gap-2 transition-colors whitespace-nowrap">
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd"></path>
              </svg>
              Export
            </button>
            <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm flex items-center gap-2 transition-colors whitespace-nowrap">
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd"></path>
              </svg>
              Create Order
            </button>
          </div>
        </div>
      
      {/* Filters */}
      <div className="mb-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div className="relative">
          <input 
            type="text" 
            placeholder="Search orders..." 
            className="pl-10 pr-4 py-2 w-full border border-gray-300 rounded-md text-sm focus:ring-emerald-500 focus:border-emerald-500"
          />
          <svg className="absolute left-3 top-2.5 w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
            <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd"></path>
          </svg>
        </div>
        
        <div className="flex gap-4">
          <div className="relative w-full">
            <select className="pl-4 pr-10 py-2 w-full border border-gray-300 rounded-md text-sm focus:ring-emerald-500 focus:border-emerald-500 appearance-none">
              {orderStatuses.map((status, index) => (
                <option key={index}>{status}</option>
              ))}
            </select>
            <svg className="absolute right-3 top-3 w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
              <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd"></path>
            </svg>
          </div>
        </div>
        
        <div className="relative">
          <select className="pl-4 pr-10 py-2 w-full border border-gray-300 rounded-md text-sm focus:ring-emerald-500 focus:border-emerald-500 appearance-none">
            {timeFilters.map((filter, index) => (
              <option key={index}>{filter}</option>
            ))}
          </select>
          <svg className="absolute right-3 top-3 w-4 h-4 text-gray-500" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
            <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd"></path>
          </svg>
        </div>
      </div>
      
      {/* Orders Table */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Order
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Customer
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Date
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Total
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Payment
                </th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {orders.map((order, index) => (
                <tr key={index} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">{order.id}</div>
                    <div className="text-xs text-gray-500">{order.items} items</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="h-8 w-8 rounded-full bg-gray-200 flex items-center justify-center text-sm font-medium text-gray-600">
                        {order.customer.avatar}
                      </div>
                      <div className="ml-3">
                        <div className="text-sm font-medium text-gray-900">{order.customer.name}</div>
                        <div className="text-xs text-gray-500">{order.customer.email}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">{order.date}</div>
                    <div className="text-xs text-gray-500">{order.shippingMethod}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">${order.amount.toFixed(2)}</div>
                    <div className="text-xs text-gray-500">USD</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold leading-5 rounded-full ${getStatusColor(order.status)}`}>
                      {order.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {order.paymentMethod}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button className="text-emerald-600 hover:text-emerald-900 mr-3">View</button>
                    <button className="text-gray-500 hover:text-gray-700">Edit</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        
        {/* Pagination */}
        <div className="px-6 py-3 flex items-center justify-between border-t border-gray-200">
          <div className="text-sm text-gray-500">
            Showing 1 to 5 of 125 orders
          </div>
          <div className="flex items-center space-x-2">
            <button className="px-3 py-1 rounded-md bg-white border border-gray-300 text-gray-500 hover:bg-gray-50 text-sm">Previous</button>
            <button className="px-3 py-1 rounded-md bg-emerald-50 border border-emerald-500 text-emerald-600 text-sm">1</button>
            <button className="px-3 py-1 rounded-md bg-white border border-gray-300 text-gray-500 hover:bg-gray-50 text-sm">2</button>
            <button className="px-3 py-1 rounded-md bg-white border border-gray-300 text-gray-500 hover:bg-gray-50 text-sm">3</button>
            <button className="px-3 py-1 rounded-md bg-white border border-gray-300 text-gray-500 hover:bg-gray-50 text-sm">Next</button>
          </div>
        </div>
      </div>
    </div>
    </div>
  )
}
