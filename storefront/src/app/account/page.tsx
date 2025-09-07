'use client';

import { useQuery } from '@apollo/client';
import { GET_CURRENT_USER, GET_ORDERS } from '@/graphql/queries';
import { useAuthStore } from '@/stores/auth';

export default function AccountDashboard() {
  const { user: localUser } = useAuthStore();
  const { data: userData, loading: userLoading, error: userError } = useQuery(GET_CURRENT_USER);
  const { data: ordersData, loading: ordersLoading, error: ordersError } = useQuery(GET_ORDERS, {
    variables: {
      filter: {
        limit: 3
      }
    }
  });
  
  const user = userData?.currentUser || localUser;
  const orders = ordersData?.orders || [];

  if (userLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (userError) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-700">Error loading account information: {userError.message}</p>
      </div>
    );
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'DELIVERED':
        return 'text-green-600';
      case 'PROCESSING':
      case 'CONFIRMED':
        return 'text-yellow-600';
      case 'SHIPPED':
        return 'text-blue-600';
      case 'CANCELLED':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

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
              <p className="font-medium">{user?.firstName} {user?.lastName}</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Email</p>
              <p className="font-medium">{user?.email}</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Member Since</p>
              <p className="font-medium">
                {user?.createdAt ? new Date(user.createdAt).toLocaleDateString() : 'N/A'}
              </p>
            </div>
          </div>
          <button className="mt-4 text-sm font-medium text-blue-600 hover:text-blue-500">
            Edit Information
          </button>
        </div>
        
        {/* Recent Orders */}
        <div className="bg-gray-50 rounded-lg p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-medium text-gray-900">Recent Orders</h2>
            <a href="/account/orders" className="text-sm font-medium text-blue-600 hover:text-blue-500">
              View All
            </a>
          </div>
          
          {ordersLoading ? (
            <div className="flex justify-center py-4">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
            </div>
          ) : ordersError ? (
            <div className="text-red-500 text-sm py-2">
              Error loading orders: {ordersError.message}
            </div>
          ) : orders.length > 0 ? (
            <div className="space-y-4">
              {orders.map((order: any) => (
                <div key={order.id} className="flex justify-between items-center">
                  <div>
                    <p className="font-medium">Order #{order.orderNumber}</p>
                    <p className="text-sm text-gray-600">
                      Placed on {new Date(order.createdAt).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="font-medium">${order.total.toFixed(2)}</p>
                    <p className={`text-sm ${getStatusColor(order.status)}`}>
                      {order.status.charAt(0) + order.status.slice(1).toLowerCase()}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-500 text-sm py-2">No orders yet</p>
          )}
        </div>
      </div>
      
      {/* Saved Items - This would be populated with real data in a full implementation */}
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