'use client';

export default function AccountDashboard() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Account Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        {/* Account Information */}
        <div className="bg-gray-50 rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-900 mb-4">Account Information</h2>
          <div className="space-y-3">
            <div>
              <p className="text-sm text-gray-600">Name</p>
              <p className="font-medium">John Doe</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Email</p>
              <p className="font-medium">john.doe@example.com</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Member Since</p>
              <p className="font-medium">January 15, 2023</p>
            </div>
          </div>
          <button className="mt-4 text-sm font-medium text-blue-600 hover:text-blue-500">
            Edit Information
          </button>
        </div>
        
        {/* Recent Orders */}
        <div className="bg-gray-50 rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-900 mb-4">Recent Orders</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <div>
                <p className="font-medium">Order #ORD-12345</p>
                <p className="text-sm text-gray-600">Placed on Jan 10, 2023</p>
              </div>
              <div className="text-right">
                <p className="font-medium">$89.99</p>
                <p className="text-sm text-green-600">Delivered</p>
              </div>
            </div>
            <div className="flex justify-between items-center">
              <div>
                <p className="font-medium">Order #ORD-12344</p>
                <p className="text-sm text-gray-600">Placed on Jan 5, 2023</p>
              </div>
              <div className="text-right">
                <p className="font-medium">$129.99</p>
                <p className="text-sm text-green-600">Delivered</p>
              </div>
            </div>
          </div>
          <button className="mt-4 text-sm font-medium text-blue-600 hover:text-blue-500">
            View All Orders
          </button>
        </div>
      </div>
      
      {/* Saved Items */}
      <div className="mb-8">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-medium text-gray-900">Saved Items</h2>
          <button className="text-sm font-medium text-blue-600 hover:text-blue-500">
            View All
          </button>
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
          {[1, 2, 3, 4].map((item) => (
            <div key={item} className="border rounded-lg overflow-hidden">
              <div className="h-32 bg-gray-200 flex items-center justify-center">
                <span className="text-gray-500">Product Image</span>
              </div>
              <div className="p-3">
                <p className="text-sm font-medium text-gray-900 line-clamp-1">Premium Wireless Headphones</p>
                <p className="text-sm text-gray-600">$89.99</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}